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
)

func main() {
	onErr := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", nil)
	onErr(err)
	defer conn.Close()

	// 如果想要消费一次后立刻退出，可以使用如下方式
	delivery, ok, err := conn.Consumer().Get("queue.direct", false)
	onErr(err)
	if ok {
		log.Println("queue.direct-get ", delivery.DeliveryTag, " ", string(delivery.Body))
		err := delivery.Ack(false)
		onErr(err)
	}
}
