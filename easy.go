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

type Consumer struct {
	c *Connection
}

// Receive 持续接收消息并消费。如果期望只接收一次消息，可以使用 Get 方法。
// 此方法是异步方法，内部使用了 go routine 执行接收操作，因此即便没有消息
// 可以接收时，该方法也不会阻塞。
//
// 详见 Channel.ReceiveOpts
func (c *Consumer) Receive(queue string, opts *ReceiveOpts, lis ReceiveListener) {
	c.c.RegisterAndExec(func(key string, ch *Channel) {
		err := ch.ReceiveOpts(queue, lis.Consumer, opts)
		if err == nil {
			lis.Remove(key, ch)
		}
		lis.Finish(err)
	})
}

// TO FIX: 当前不能直接使用 consumer.Get()，会导致关闭错误。
//// When autoAck is true, the server will automatically acknowledge this message so you don't have to.
//// But if you are unable to fully process this message before the channel or connection is closed,
//// the message will not get requeued
//func (c *Consumer) Get(queue string, autoAck bool) (*amqp.Delivery, bool, error) {
//	ch, err := c.c.Channel()
//	if err != nil {
//		return nil, false, err
//	}
//	defer ch.Close()
//	msg, ok, err := ch.Get(queue, autoAck)
//	if err != nil {
//		return nil, false, err
//	}
//	return &msg, ok, err
//}

type Producer struct {
	c *Connection
}

// Send 发送消息。
//
// 参数 body 即需要发送的消息。
//
// 参数 opts 即发送消息需要配置的选项。如果 opts 为 nil，则表示使用默认配置。可以通过配置 SendOpts.retryable
// 启用消息重发的能力。请注意，由于消息重发使用的是同步的方式处理 ack，因此启用消息重发会极大降低 QPS。
func (p *Producer) Send(exchange string, routingKey string, body []byte, opts *SendOpts) error {
	ch, err := p.c.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.SendOpts(exchange, routingKey, body, opts)
}
