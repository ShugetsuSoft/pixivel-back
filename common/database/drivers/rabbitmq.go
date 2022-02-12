package drivers

import (
	"log"
	"sync"

	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn            *amqp.Connection
	chann           *amqp.Channel
	closesign       chan error
	reconnectsignal *sync.Cond
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	rbmq := &RabbitMQ{}
	reconnectsignal := sync.NewCond(&sync.Mutex{})
	rbmq.reconnectsignal = reconnectsignal

	connect := func() (chan *amqp.Error, error) {
		conn, err := amqp.Dial(uri)
		if err != nil {
			return nil, err
		}
		rbmq.conn = conn
		ch, err := conn.Channel()
		ch.Qos(100, 0, false)
		if err != nil {
			return nil, err
		}
		rbmq.chann = ch
		reconnectsignal.Broadcast()
		closesignal := make(chan *amqp.Error, 1)
		return conn.NotifyClose(closesignal), nil
	}

	closechan, err := connect()
	if err != nil {
		return nil, err
	}
	rbmq.closesign = make(chan error)

	go func() {
		for {
			select {
			case e := <-closechan:
				telemetry.Log(telemetry.Label{"pos": "rabbitmq"}, e.Error())
				closechan, err = connect()
				if err != nil {
					log.Fatal(err)
				}
			case <-rbmq.closesign:
				rbmq.conn.Close()
				rbmq.chann.Close()
				close(rbmq.closesign)
				close(closechan)
				return
			}
		}
	}()

	return rbmq, nil
}

func (mq *RabbitMQ) SetQos(prefetchCount int, prefetchSize int, global bool) error {
	return mq.chann.Qos(prefetchCount, prefetchSize, global)
}

func (mq *RabbitMQ) QueueDeclare(name string) (amqp.Queue, error) {
	queue, err := mq.chann.QueueDeclare(
		name,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return queue, err
}

func (mq *RabbitMQ) ExchangeDeclare(name string, kind string) error {
	return mq.chann.ExchangeDeclare(
		name,
		kind,
		true, false, false, false, nil)
}

func (mq *RabbitMQ) QueueBindExchange(name, key, exchange string) {
	mq.chann.QueueBind(name, key, exchange, false, nil)
}

func (mq *RabbitMQ) Publish(name string, message models.MQMessage) error {
	return mq.chann.Publish(
		"",    // exchange
		name,  // routing key
		false, // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Body:         message.Data,
			Priority:     message.Priority,
		},
	)
}

func (mq *RabbitMQ) PublishToExchange(name string, message models.MQMessage) error {
	return mq.chann.Publish(
		name,  // exchange
		name,  // routing key
		false, // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			Body:         message.Data,
			Priority:     message.Priority,
		},
	)
}

func (mq *RabbitMQ) Consume(name string) (<-chan models.MQMessage, error) {
	reschan := make(chan models.MQMessage, 50)
	mq.reconnectsignal.L.Lock()
	if mq.chann.IsClosed() {
		mq.reconnectsignal.Wait()
	}
	mq.reconnectsignal.L.Unlock()

	oniichan, err := mq.chann.Consume(
		name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	go func() {
		defer close(reschan)
		for onii := range oniichan {
			reschan <- models.MQMessage{
				Data:     onii.Body,
				Tag:      onii.DeliveryTag,
				Priority: onii.Priority,
			}
		}
	}()
	return reschan, err
}

func (mq *RabbitMQ) Ack(tag uint64) error {
	return mq.chann.Ack(tag, false)
}

func (mq *RabbitMQ) Reject(tag uint64) error {
	return mq.chann.Reject(tag, false)
}

func (mq *RabbitMQ) Close() {
	mq.closesign <- nil
}
