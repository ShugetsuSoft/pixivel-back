package models

type MessageQueue interface {
	Publish(string, MQMessage) error
	PublishToExchange(string, MQMessage) error
	Consume(string) (<-chan MQMessage, error)
	Ack(uint64) error
	Reject(uint64) error
}

type Filter interface {
	Add(string, string) (bool, error)
	Exists(string, string) (bool, error)
}
