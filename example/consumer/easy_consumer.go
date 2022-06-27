package main

import (
	"ezmq"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", ezmq.DefaultTimesRetry())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	consumer := conn.Consumer()
	consumer.Receive(
		"queue.direct",
		ezmq.NewReceiveOptsBuilder().SetAutoAck(false).Build(),
		&ezmq.AbsReceiveListener{
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
					log.Fatal(err)
				}
				// defer xxx.close() // 关闭资源操作等
			},
		})

	time.Sleep(time.Minute)
}
