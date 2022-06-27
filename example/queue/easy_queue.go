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
