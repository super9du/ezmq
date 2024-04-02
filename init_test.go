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
	"context"
	"log"
	"testing"
)

var defaultURL = "amqp://guest:guest@localhost:5672/"

func TestInitBase(t *testing.T) {
	c, err := Dial(defaultURL, DefaultTimesRetry())
	onErr(err)
	defer c.Close()

	ch, err := c.Channel()
	onErr(err)

	// declare queue
	_, err = ch.QueueDeclare("queue.direct", true, false, false, false, nil)
	onErr(err)
	err = ch.QueueBind("queue.direct", "key.direct", "amq.direct", false, nil)
	onErr(err)
}

// ---- utils ----

func onErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getConnection() *Connection {
	c, err := Dial(defaultURL, DefaultTimesRetry())
	onErr(err)
	return c
}

func getChannel() (*Channel, *Connection) {
	var err error

	conn := getConnection()
	channel, err := conn.Channel()
	onErr(err)

	return channel, conn
}

func getChannelWithContext() (*Channel, context.CancelFunc, *Connection) {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	conn, err := Dial(defaultURL, DefaultCtxRetry(ctx))
	onErr(err)
	channel, err := conn.Channel()
	onErr(err)

	return channel, cancel, conn
}
