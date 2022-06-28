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
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", ezmq.DefaultTimesRetry())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queueBuilder := conn.QueueBuilder()
	queue := queueBuilder.
		SetQueueDeclareOpts(func(builder *ezmq.QueueDeclareOptsBuilder) *ezmq.QueueDeclareOpts {
			return builder.SetNowait(false).SetDurable(true).SetArgs(nil).Build()
		}).
		SetQueueBindOpts(func(builder *ezmq.QueueBindOptsBuilder) *ezmq.QueueBindOpts {
			return builder.SetNoWait(false).SetArgs(&amqp.Table{}).Build()
		}).
		SetRetryable(ezmq.DefaultTimesRetry()).
		Build()

	err = queue.DeclareAndBind("test", "key.test", "amq.direct")
	if err != nil {
		log.Fatal(err)
	}
}
