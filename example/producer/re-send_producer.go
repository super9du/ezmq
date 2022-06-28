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

// 断线重连，消息确认，失败重发
func main() {
	onErr := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", ezmq.DefaultTimesRetry())
	onErr(err)
	defer conn.Close()

	ch, _ := conn.Channel()
	defer ch.Close()
	err = ch.SendOpts(
		"amq.direct",
		"key.direct",
		[]byte("reSendSync() | "+time.Now().Format("2006-01-02 15:04:05")),
		ezmq.NewSendOptsBuilder().
			SetRetryable(ezmq.DefaultTimesRetry()).
			SetMessageFactory(ezmq.MessageJsonPersistent).
			Build(),
	)
	onErr(err)
}
