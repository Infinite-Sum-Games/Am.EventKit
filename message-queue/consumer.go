package messagequeue

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *MsgBroker) StartConsumer(queueName string, handler func(d amqp.Delivery)) error {

	return nil
}
