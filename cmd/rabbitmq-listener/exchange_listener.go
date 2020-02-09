package main

import (
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQExchangeListener struct {
	RabbitMQListener
	connection *amqp.Connection
	config *Config
	channels map[string]*amqp.Channel
}

func (v *RabbitMQExchangeListener) Init() error {
	v.channels = map[string]*amqp.Channel{}

	return nil
}

func (v *RabbitMQExchangeListener) DeclareResources() error {
	for _, item := range v.config.Resources.Exchanges {
		ch, err := v.connection.Channel()
		if err != nil {
			return err
		}

		err = ch.ExchangeDeclare(
			item.Name, 			// name
			item.ResourceType,	// type
			item.Durable,       // durable
			item.AutoDeleted,   // auto-deleted
			item.Internal,      // internal
			item.NoWait,        // no-wait
			nil,          		// arguments
		)

		if err != nil {
			return err
		}

		v.channels[item.Name] = ch
	}

	return nil
}


func (v *RabbitMQExchangeListener) Listen(config *Config) error {
	for commandKey, item := range v.config.Commands {
		resource, err := v.config.GetResource(item.Resource)

		if err != nil {
			continue
		}

		ch, ok := v.channels[item.Resource]

		if !ok {
			return fmt.Errorf("Channel for %s does not exists", item.Resource)
		}

		if resource.Category != RESOURCE_EXCHANGE {
			continue
		}

		q, err := ch.QueueDeclare(
			"",    // name
			false, // durable
			false, // delete when unused
			true,  // exclusive
			false, // no-wait
			nil,   // arguments
		)

		if err != nil {
			return err
		}
		
		err = ch.QueueBind(
			q.Name,       		// queue name
			item.RoutingKey,    // routing key
			item.Resource, 		// exchange
			false,
			nil)

		if err != nil {
			return err
		}
		
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto ack
			false,  // exclusive
			false,  // no local
			false,  // no wait
			nil,    // args
		)

		if err != nil {
			return err
		}

		go func(commandKey string, msgs <-chan amqp.Delivery) {
			var cmdExec *CommandExec

			fmt.Printf(" [%s] Start listening\n", commandKey)
			for d := range msgs {
				fmt.Printf(" [x][%s] %s\n", commandKey, d.Body)
				cmd, err := config.GetCommand(commandKey)
				
				if err != nil {
					break
				}

				if cmd.CommandPre != "" {
					cmdExec, _ = PrepareCommand(cmd.CommandPre, logger)
					cmdExec.Execute()
				}
				
				if cmd.Command != "" {
					cmdExec, _ = PrepareCommand(cmd.Command, logger)
					cmdExec.Execute()
				}

				if cmd.CommandPost != "" {
					cmdExec, _ = PrepareCommand(cmd.CommandPost, logger)
					cmdExec.Execute()
				}
			}
			fmt.Printf(" [%s] Stop listening\n", commandKey)
		}(commandKey, msgs)
	}

	return nil
}

func (v *RabbitMQExchangeListener) Destroy() error {
	for _, item := range v.channels {
		item.Close()
	}

	return nil
}