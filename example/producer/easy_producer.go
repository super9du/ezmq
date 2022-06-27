package main

import (
	"ezmq"
	"log"
	"time"
)

func main() {
	conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", ezmq.DefaultTimesRetry())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	producer := conn.Producer()
	err = producer.Send("amq.direct", "key.direct",
		[]byte("producer.Send() | "+time.Now().Format("2006-01-02 15:04:05")),
		ezmq.DefaultSendOpts())
	if err != nil {
		log.Fatal(err)
	}
}
