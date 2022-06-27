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
