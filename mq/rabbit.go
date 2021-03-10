package mq

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"telegram-parser/db"
)

/*---------------------------------STRUCTURES----------------------------------------*/

type Rabbit struct {
	Connection   *amqp.Connection
	Channel      *amqp.Channel
	UpdatesQueue amqp.Queue
}

/*-----------------------------------METHODS-----------------------------------------*/

// RabbitInit returns a message broker instance with the required queue connections.
func RabbitInit() *Rabbit {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@localhost:%v/", os.Getenv("RABBITPORT")))
	if err != nil {
		logrus.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatal(err)
	}

	updatesQ, err := ch.QueueDeclare(
		"updates",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logrus.Fatal(err)
	}

	rabbit := &Rabbit{
		Connection:   conn,
		Channel:      ch,
		UpdatesQueue: updatesQ,
	}

	return rabbit
}

// CloseConn cleanly closes the RabbitMQ Channel and Connection.
func (r *Rabbit) CloseConn() {
	r.Channel.Close()
	r.Connection.Close()
}

// Publish sends a Publishing from the client to an exchange on the server.
func (r *Rabbit) Publish(msg *db.Message) error {
	body, err := marshalMessage(msg)
	if err != nil {
		return err
	}

	err = r.Channel.Publish(
		"",
		r.UpdatesQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	return err
}

// Consume immediately starts delivering queued messages
func (r *Rabbit) Consume() <-chan amqp.Delivery {
	msgs, err := r.Channel.Consume(
		r.UpdatesQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		logrus.Fatal(err)
	}

	return msgs
}

/*-----------------------------------HELPERS-----------------------------------------*/

// marshalMessage returns the JSON encoding of *db.Message.
func marshalMessage(msg *db.Message) (result []byte, err error) {
	bytes, err := json.Marshal(msg)
	return bytes, err
}

// UnmarshalRabbitBody  parses the JSON-encoded data into *db.Message.
func UnmarshalRabbitBody(bytes []byte) (message *db.Message, err error) {
	var msg db.Message
	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		return nil, err
	}

	return &msg, nil
}
