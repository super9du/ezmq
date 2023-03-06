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
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func ExampleConsumer_Receive() {
	conn := getConnection()

	consumer := conn.Consumer()
	consumer.Receive(
		"queue.direct",
		NewReceiveOptsBuilder().SetAutoAck(false).Build(),
		&AbsReceiveListener{
			ConsumerMethod: func(delivery *amqp.Delivery) (brk bool) {
				log.Println("queue.direct ", delivery.DeliveryTag, " ", string(delivery.Body))
				err := delivery.Ack(false)
				if err != nil {
					log.Println(err)
				}
				return
			},
			FinishMethod: func(err error) {
				if err != nil {
					// 处理错误
					return
				}
				// defer xxx.close() // 关闭资源操作等
			},
		})

	time.Sleep(time.Minute) // 由于 Consumer.Receive() 内部采用了异步方式处理，因此 Receive 方法不会阻塞等待
}
