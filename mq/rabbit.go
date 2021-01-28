package mq

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"telegram-parser/db"
)

//    --------------------------------------------------------------------------------
//                                    STRUCTS
//    --------------------------------------------------------------------------------

type Rabbit struct {
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	UpdatesQueue amqp.Queue
}

//    --------------------------------------------------------------------------------
//                                     METHODS
//    --------------------------------------------------------------------------------

// RabbitInit returns a message broker instance with the required queue connections
func RabbitInit() *Rabbit {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logrus.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		"updates", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		logrus.Fatal(err)
	}

	rabbit := &Rabbit{
		Connection:   conn,
		Channel:      ch,
		UpdatesQueue: q,
	}

	return rabbit
}

func (r *Rabbit) CloseConn() {
	r.Channel.Close()
	r.Connection.Close()
}

func (r *Rabbit) Publish(msg *db.Message) error {
	body := marshalMessage(msg)

	err := r.Channel.Publish(
		"",                  // exchange
		r.UpdatesQueue.Name, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	return err
}

func (r *Rabbit) Consume() <-chan amqp.Delivery {
	msgs, err := r.Channel.Consume(
		r.UpdatesQueue.Name, // queue
		"",                  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)

	if err != nil {
		logrus.Fatal(err)
	}

	return msgs
}

//    --------------------------------------------------------------------------------
//                                     HELPERS
//    --------------------------------------------------------------------------------

func marshalMessage(msg *db.Message) []byte {
	bytes, err := json.Marshal(msg)
	if err != nil {
		logrus.Fatal(err)
	}

	return bytes
}

func UnmarshalRabbitBody(bytes []byte) *db.Message {
	var msg db.Message
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		logrus.Fatal(err)
	}

	return &msg
}
