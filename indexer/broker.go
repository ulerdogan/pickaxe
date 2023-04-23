package indexer

import (
	"github.com/streadway/amqp"
	config "github.com/ulerdogan/pickaxe/utils/config"
	logger "github.com/ulerdogan/pickaxe/utils/logger"
)

func SetupRabbitMQ(cnfg config.Config) (*amqp.Channel, error) {
	conn, err := amqp.Dial(cnfg.RMQUrl)
	if err != nil {
		logger.Error(err, "cannot connect to the rabbitmq")
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Error(err, "cannot create the rabbitmq channel")
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"EventsQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Error(err, "cannot declare the rabbitmq queue")
		return nil, err
	}

	logger.Info("rabbitmq succesfully initialized")
	return ch, nil
}

func (ix *indexer) PublishRmqMsg(msg []byte) error {
	err := ix.RabbitMQ.Publish(
		"",
		"EventsQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		},
	)

	return err
}
