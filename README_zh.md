# ezmq：一个易于使用的AMQP客户端

介绍
---
* 💥扩展了 [rabbitmq/amqp091-go](https://github.com/rabbitmq/amqp091-go) 包，除此之外零依赖
* 💪支持网络断线后自动重连
* 💎支持消息发送失败后自动重试
* 🎈易于使用

实现方式简介
---
由于 AMQP 消息的发送和接收，依赖 Channel 的创建，Channel 的创建依赖与 Connection 的连接。如果断开网络连接，Connection 就会断开。
但是 Channel 并不知道是谁创建了自己。所以通常情况下，我们有两种断线重连的方案：

* 一种是 Connection 断开后，Channel 自己重新获取 Connection。此种方式是最简单的实现方案。
  但此法的问题是，如果 Channel 数量极其庞大，
  每个 Channel 都会创建 Connection 并连接。当网络状况不大好的时候，可能会有数量极其庞大的 Channel 反复尝试重连，
  导致服务器资源占用会暴增，甚至加剧网络的阻塞。但其实我们只需要一个 Connection 去判断网络是否能连接、已连接。
* 一种是 Connection 断开后，Connection 自己重连。如果重连不成功，Channel 不做任何操作。如果重连成功，自动重跑注册的操作。
  而如果使用断线重连发送消息，将在 Connection 连接成功后继续发送需要发送的消息。此种方式实现比较复杂，
  但避免了服务器资源占用以及加剧网络阻塞的问题。

ezmq 采用的是后者。

快速上手
---

### 发送消息

```go
package main

import (
  "ezmq"
  "log"
  "time"
)

func main() {
  // 创建 Connection 并连接服务器
  conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", ezmq.DefaultTimesRetry())
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  // 创建 Producer
  producer := conn.Producer()
  // 发送消息
  err = producer.Send("amq.direct", "key.direct", []byte("producer.Send() | "+time.Now().Format("2006-01-02 15:04:05")),
    ezmq.DefaultSendOpts())
  if err != nil {
    log.Fatal(err)
  }
}

```

### 接收消息

```go
package main

import (
  "ezmq"
  amqp "github.com/rabbitmq/amqp091-go"
  "log"
)

func main() {
  // 创建 Connection 并连接服务器
  conn, err := ezmq.Dial("amqp://guest:guest@localhost:5672/", ezmq.DefaultTimesRetry())
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  // 创建 Consumer
  consumer := conn.Consumer()
  // 接收消息
  consumer.Receive(
    "queue.direct",                                         // 队列名
    ezmq.NewReceiveOptsBuilder().SetAutoAck(false).Build(), // 接收选项
    &ezmq.AbsReceiveListener{
      ConsumerMethod: func(d *amqp.Delivery) (brk bool) { // 消费者方法
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

> 小建议：
> 
> 如无特殊需求，我们在全局只需创建一个 Connection，所有消息的发送和接收都使用 Producer 和 Consumer 处理。

待实现
---

* 消息异步发送与确认
* 重复消费的处理
* 网络分区状态下，找不到队列时消息重发（ReturnListener）
* 消费者 ack 重试

License
---

LGPL