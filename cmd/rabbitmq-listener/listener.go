package main

import (
	"github.com/streadway/amqp"
)

const (
	// ResourceQueue ... 
	ResourceQueue = "queue"
	// ResourceExchange ...
	ResourceExchange = "exchange"
)

// RabbitMQListener ...
type RabbitMQListener interface {
	Init() error
	DeclareResources() error
	Listen() error
	Destroy() error
}

func consumeResource(cmd *Command, msgs <-chan amqp.Delivery) {
	var cmdExec *CommandExec
	for d := range msgs {
		logger.Info("consumer", "Got message", "body", d.Body)
		if cmd.CommandPre != "" {
			cmdExec, _ = PrepareCommand(cmd.CommandPre, &logger)
			cmdExec.Execute()
		}
		
		if cmd.Command != "" {
			cmdExec, _ = PrepareCommand(cmd.Command, &logger)
			cmdExec.Execute()
		}

		if cmd.CommandPost != "" {
			cmdExec, _ = PrepareCommand(cmd.CommandPost, &logger)
			cmdExec.Execute()
		}
	}
}