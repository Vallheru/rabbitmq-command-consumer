package main

import (
	// "fmt"
	// "os"
	"io/ioutil"
	"github.com/streadway/amqp"
)

func main() {
	data, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	content := string(data)

	forever := make(chan bool)
	
	config, err := parseConfigString(&content)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()

	exchangeListener := RabbitMQExchangeListener{
		connection: conn,
		config: &config,
	}
	exchangeListener.Init()
	exchangeListener.DeclareResources()
	exchangeListener.Listen()
	defer exchangeListener.Destroy()

	<-forever
}
