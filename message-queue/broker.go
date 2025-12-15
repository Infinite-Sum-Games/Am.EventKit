package messagequeue

import (
	"context"
	"fmt"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Add in queue names here and update the allQueues array
const (
	QueueHackathonRegistrations = "ai-hackathon-registrations"
	QueueWocRegistrations       = "woc-registrations"
)

var allQueues = []string{
	QueueWocRegistrations,
	QueueHackathonRegistrations,
}

var Rabbit *MsgBroker

// Using RabbitMQ as a message broker
type MsgBroker struct {
	conn    *amqp.Connection
	connURL string
	channel *amqp.Channel
}

func NewBroker(connStr string) (*MsgBroker, error) {
	client := &MsgBroker{
		connURL: connStr,
	}

	if err := client.connect(); err != nil {
		pkg.Log.Error("[CRASH]: Message broker failed to initialize.", err)
		return nil, err
	}

	go client.handleReconnect()

	return client, nil
}

func (r *MsgBroker) connect() error {
	var err error
	r.conn, err = amqp.Dial(r.connURL)
	if err != nil {
		return err
	}

	// Once connection has been established, a channel is created. This is
	// used to receive all commands to RabbitMQ in the AMQP format.
	r.channel, err = r.conn.Channel()
	if err != nil {
		// If this too fails, it's fine as the parent process will be killed
		_ = r.conn.Close()
		return err
	}

	// Once channele establishment is completed, initialize all queues
	if err := r.declareQueues(); err != nil {
		return err
	}

	return nil
}

func (r *MsgBroker) declareQueues() error {
	for _, queueName := range allQueues {
		_, err := r.channel.QueueDeclare(
			queueName,
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // args
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
		pkg.Log.Info(fmt.Sprintf("Successfully declared queue: %s", queueName))
	}

	pkg.Log.Info("All queues declared successfully.")
	return nil
}

// This listens for connection errors and attempts to reconnect. It is done
// by creating a receiver channel called "errChan". This channel is registered as
// the receiver channel in NotifyClose function. In case of a graceful shutdown
// no error would be propagated but in-case of ungraceful ones, errors would
// be sent to this channel.
func (r *MsgBroker) handleReconnect() {
	errChan := make(chan *amqp.Error)
	r.conn.NotifyClose(errChan)

	// We setup an event-driven loop which is blocking in nature. It will wait
	// indefinitely until a value is sent into errChan.
	for err := range errChan {
		pkg.Log.Error("Broker connection lost. Attempting to reconnect...", err)

		// Try repeatedly to re-establish connection INFINITELY with 5s backoff
		for {
			time.Sleep(5 * time.Second)

			// Handles reconnection and redeclaration of queues.
			if err := r.connect(); err == nil {
				pkg.Log.Info("Successfully reconnected to message broker")
				r.conn.NotifyClose(errChan)

				// Re-establish the error channel. The c.conn object is a new connection
				// as the old connection object has been garbage collected and the
				// error channel is now dangling in an unregistered state and is waiting
				// to be released and garbage collected, so we replace it with a new
				// channel whose memory has been freshly allocated and references the
				// new connection.

				// If we don't do this then the new connection, if it fails; would do so
				// silently as there are no channels listening to its errors and then
				// reconnection would not happen. In short, your handleReconnect()
				// goroutine would be stuck forever waiting for message on a channel
				// whose publisher does not exist anymore. Dead rabbit :)
				errChan = make(chan *amqp.Error)
				r.conn.NotifyClose(errChan)

				break
			}

			pkg.Log.Error("Failed to reconnect, retrying...", err)
		}
	}
	// End of handleReconnect
}

// Errors will be thrown in case of wrong queue names
func (r *MsgBroker) Publish(ctx context.Context, queueName string, body []byte) error {

	// The first parameter is an exchange which is an internal routing agent for
	// RabbitMQ. it receives messages and then send them to apt queues. An
	// empty string is a default excahange which is "direct" by default. It
	// means that it sends the message to exactly one queue whose name matches
	// the routing key

	// The Routing Key is the name of the queue to which the msg is published

	// Mandatory bool is a flag that tells the broker to return the message to
	// the producer if set to True or silently discard the message if set to False
	// I am setting this to false because I don't expect to go wrong with queue
	// names in my application code as I will be using constants that are earlier
	// declared in this file.

	// The Immediate flag tells what to do if there are no consumers actively
	// listening on the other end of the queue. If set to True then the message
	// would be returned back otherwise it will be left in the queue. I am setting
	// it to False because there might be scenarios where the consumer crashes
	// and restarts.

	err := r.channel.Publish(
		"",        // default exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json", // easy to deal with, although expensive
			Body:         body,               // bytes
			DeliveryMode: amqp.Persistent,    // survive broker restarts
		})

	return err
}

func (r *MsgBroker) Close() error {
	if r.channel == nil {
		return fmt.Errorf("no channels to rabbit-mq")
	}
	if r.conn == nil {
		return fmt.Errorf("no connection to rabbit-mq")
	}

	if err := r.channel.Close(); err != nil {
		return err
	}
	if err := r.conn.Close(); err != nil {
		return err
	}

	return nil
}
