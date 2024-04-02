# ezmq: An Easy-to-Use AMQP Client

[中文文档](README_zh.md)

Introduction
---
* 💥 Extends the [rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go) package with zero dependencies beyond it
* 💪 Implements network disconnection reconnection
* 💎 Implements retry on message send failure
* 🎈 User-friendly

Brief on Implementation
---
Since the sending and receiving of AMQP messages depend on the creation of Channels, which in turn depends on the connection of Connections, a network disconnection would drop the Connection. However, Channels are not aware of their creator. Typically, there are two solutions for reconnecting after a disconnection:

* One is for Channels to reacquire a Connection after a Connection drops. This approach is the simplest to implement. However, if there are a vast number of Channels, each would create and connect to a Connection. In poor network conditions, an enormous number of Channels might attempt to reconnect repeatedly, leading to a surge in server resource consumption and potentially exacerbating network congestion. Yet, we only need one Connection to determine if the network is connectable or connected.
* Another approach is for the Connection itself to reconnect after a disconnection. Channels do not perform any operations if the reconnection fails. If the reconnection succeeds, it automatically reruns registered operations. Using disconnection reconnection to send messages continues to send the necessary messages after the Connection successfully reconnects. This method is more complex but avoids the issues of server resource consumption and network congestion.

ezmq adopts the latter approach.

Quick Start
---

### Sending Messages

```go
package main

import (
  "ezmq"
  "log"
  "time"
)

func main() {
  // Create and connect to a Connection
  conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", ezmq.DefaultTimesRetry())
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  // Create a Producer
  producer := conn.Producer()
  // Send a message
  err = producer.Send("amq.direct", "key.direct", []byte("producer.Send() | "+time.Now().Format("2006-01-02 15:04:05")),
    ezmq.DefaultSendOpts())
  if err != nil {
    log.Fatal(err)
  }
}

```

### Receiving Messages

```go
package main

import (
  "ezmq"
  amqp "github.com/rabbitmq/amqp091-go"
  "log"
)

func main() {
  // Create and connect to a Connection
  conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", ezmq.DefaultTimesRetry())
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  // Create a Consumer
  consumer := conn.Consumer()
  // Receive a message
  consumer.Receive(
    "queue.direct",                                         // Queue name
    ezmq.NewReceiveOptsBuilder().SetAutoAck(false).Build(), // Receiving options
    &ezmq.AbsReceiveListener{
      ConsumerMethod: func(d *amqp.Delivery) (brk bool) { // Consumer method
        log.Println("queue.direct ", d.DeliveryTag, " ", string(d.Body))
        err := d.Ack(false)
        if err != nil {
          log.Println(err)
        }
        return
      },
    })
}

```

> Pro Tip:
>
> Unless there's a specific need, it's advisable to create only one Connection globally, with all message sending and receiving handled by the Producer and Consumer.

To Be Implemented
---

* Asynchronous message sending and confirmation
* Handling of duplicate consumption
* Message resend when a queue is not found under network partition conditions (ReturnListener)
* Consumer ack retry

License
---

LGPL