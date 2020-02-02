package main

const (
	RESOURCE_QUEUE = "queue"
	RESOURCE_EXCHANGE = "exchange"
)

type RabbitMQListener interface {
	Init() error
	DeclareResources() error
	Listen() error
	Destroy() error
}

