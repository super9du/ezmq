// ezmq: An easy golang amqp client.
// Copyright (C) 2022  super9du
//
// This library is free software; you can redistribute it and/or
// modify it under the terms of the GNU Lesser General Public
// License as published by the Free Software Foundation; either
// version 2.1 of the License, or (at your option) any later version.
//
// This library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public
// License along with this library; If not, see <https://www.gnu.org/licenses/>.

package ezmq

import (
	"fmt"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestChannel_Send(t *testing.T) {
	channel, conn := getChannel()
	defer conn.Close()

	type args struct {
		exchange   string
		routingKey string
		body       []byte
	}
	tests := []struct {
		name    string
		Channel *Channel
		args    args
		wantErr bool
	}{
		{name: "send", Channel: channel, args: args{
			exchange:   "amq.direct",
			routingKey: "key.direct",
			body:       []byte("send | " + time.Now().Format("2006-01-02 15:04:05")),
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.Channel
			if err := c.Send(tt.args.exchange, tt.args.routingKey, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChannel_Send_reSend(t *testing.T) {
	channel, conn := getChannel()
	defer conn.Close()

	type args struct {
		exchange   string
		routingKey string
		body       []byte
		retryable  Retryable
	}
	tests := []struct {
		name    string
		Channel *Channel
		args    args
		wantErr bool
	}{
		{name: "reSendSync BYTIMES", Channel: channel, args: args{
			exchange:   "amq.direct",
			routingKey: "key.direct",
			body:       []byte("reSendSync BYTIMES | " + time.Now().Format("2006-01-02 15:04:05")),
			retryable:  &TimesRetry{RetryTimes: 10, Interval: 3 * time.Second},
		}, wantErr: false},
		{name: "reSendSync ALWAYS", Channel: channel, args: args{
			exchange:   "amq.direct",
			routingKey: "key.direct",
			body:       []byte("reSendSync ALWAYS | " + time.Now().Format("2006-01-02 15:04:05")),
			retryable:  &TimesRetry{Always: true, Interval: 3 * time.Second},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.Channel
			if err := c.SendOpts(
				tt.args.exchange, tt.args.routingKey, tt.args.body,
				NewSendOptsBuilder().SetRetryable(tt.args.retryable).Build(),
			); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChannel_Receive(t *testing.T) {
	channel, conn := getChannel()
	defer conn.Close()

	type args struct {
		queue string
		fn    ConsumerFunc
	}
	tests := []struct {
		name    string
		Channel *Channel
		args    args
		wantErr bool
	}{
		{name: "receive", Channel: channel, args: args{
			queue: "queue.direct",
			fn: func(d *amqp.Delivery) (brk bool) {
				fmt.Println("DeliveryTag:", d.DeliveryTag, "| Receive:", string(d.Body))
				return
			},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				c := tt.Channel
				if err := c.Receive(tt.args.queue, tt.args.fn); (err != nil) != tt.wantErr {
					t.Errorf("Receive() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
		})
	}
	time.Sleep(time.Minute * 3)
}

func TestChannel_Receive_With_Context(t *testing.T) {
	channel, cancel, conn := getChannelWithContext()
	// 正常情况应该立刻使用 defer 语句调用 cancel 函数。
	// 这里为了测试，在下面使用了，所以这里没有调用 `defer cancel()`。
	defer conn.Close()

	type args struct {
		queue string
		fn    ConsumerFunc
	}
	tests := []struct {
		name    string
		Channel *Channel
		args    args
		wantErr bool
	}{
		{name: "receive", Channel: channel, args: args{
			queue: "queue.direct",
			fn: func(d *amqp.Delivery) (brk bool) {
				fmt.Println("DeliveryTag:", d.DeliveryTag, "| Receive:", string(d.Body))
				return
			},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				c := tt.Channel
				if err := c.Receive(tt.args.queue, tt.args.fn); (err != nil) != tt.wantErr {
					t.Errorf("Receive() error = %v, wantErr %v", err, tt.wantErr)
				}
			}()
		})
	}
	time.Sleep(time.Second * 30)
	cancel()
	time.Sleep(time.Second * 3)
	info("finish!")
}

func TestChannel_Receive_limit_time(t *testing.T) {
	ch, conn := getChannel()
	defer conn.Close()

	go func() {
		err := ch.Receive("queue.direct", func(d *amqp.Delivery) (brk bool) {
			fmt.Println("DeliveryTag:", d.DeliveryTag, " | Receive:", string(d.Body))
			return
		})
		onErr(err)
	}()
	time.Sleep(time.Minute * 3)
}
