package rabbitmq

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Item struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	OrderID    string  `json:"order_id"`
	UserID     string  `json:"user_id"`
	Items      []Item  `json:"items"`
	TotalPrice float64 `json:"total_price"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Producer() {
	// 连接到 RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 声明队列
	q, err := ch.QueueDeclare(
		"order_queue", // 队列名称
		true,          // 是否持久化
		false,         // 是否自动删除
		false,         // 是否独占
		false,         // 是否阻塞
		nil,           // 额外参数
	)
	failOnError(err, "Failed to declare a queue")

	// 创建订单
	order := Order{
		OrderID:    "123456",
		UserID:     "987",
		Items:      []Item{{ItemID: "A001", Quantity: 2}, {ItemID: "B002", Quantity: 1}},
		TotalPrice: 150.0,
	}

	// 序列化订单信息
	body, err := json.Marshal(order)
	failOnError(err, "Failed to marshal order")

	// 发送消息到队列
	err = ch.Publish(
		"",     // 交换器名称
		q.Name, // 路由键，即队列名称
		false,  // 如果没有队列绑定是否返回消息
		false,  // 如果队列没有消费者是否返回消息
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent order: %s", body)
}

func Consumer() {
	// 连接到 RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// 声明队列（确保队列存在）
	q, err := ch.QueueDeclare(
		"order_queue", // 队列名称
		true,          // 是否持久化
		false,         // 是否自动删除
		false,         // 是否独占
		false,         // 是否阻塞
		nil,           // 额外参数
	)
	failOnError(err, "Failed to declare a queue")

	// 注册消费者
	msgs, err := ch.Consume(
		q.Name, // 队列名称
		"",     // 消费者名称
		true,   // 自动确认
		false,  // 是否独占
		false,  // 是否不等待
		false,  // 是否阻塞
		nil,    // 额外参数
	)
	failOnError(err, "Failed to register a consumer")

	// 消费消息
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var order Order
			err := json.Unmarshal(d.Body, &order)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}

			log.Printf(" [x] Received order: %+v", order)

			// 模拟扣减库存
			for _, item := range order.Items {
				log.Printf(" - Deducting %d of item %s", item.Quantity, item.ItemID)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
