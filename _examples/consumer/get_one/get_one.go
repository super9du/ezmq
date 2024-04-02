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

package main

import (
	"ezmq"
	"log"
	"time"
)

func main() {
	onErr := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	var err error

	conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", nil)
	onErr(err)
	defer conn.Close()

	ch, err := conn.Channel()
	onErr(err)
	defer ch.Close()

	// 设置为 Plain Persistent （无格式/持久化）的形式发送消息
	err = ch.SendOpts(
		"amq.direct", "key.direct", []byte("Send() | "+time.Now().Format("2006-01-02 15:04:05")),
		ezmq.NewSendOptsBuilder().SetMessageFactory(ezmq.MessagePlainPersistent).Build(),
	)
	onErr(err)
}

// TO FIX: 当前不能直接使用 consumer.Get()，会导致关闭错误。
//func main() {
//	onErr := func(err error) {
//		if err != nil {
//			log.Fatal(err)
//		}
//	}
//
//	conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", nil)
//	onErr(err)
//	defer conn.Close()
//
//	// 如果想要消费一次后立刻退出，可以使用如下方式
//	delivery, ok, err := conn.Consumer().Get("queue.direct", false)
//	onErr(err)
//	if ok {
//		log.Println("queue.direct-get ", delivery.DeliveryTag, " ", string(delivery.Body))
//		err := delivery.Ack(false) // Exception (504) Reason: "channel/connection is not open"
//		onErr(err)
//	}
//}
