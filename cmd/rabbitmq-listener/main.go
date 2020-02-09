package main

import (
	"io/ioutil"
	"github.com/streadway/amqp"
)

var (
	logger Logger
)

func main() {
	var (
		data []byte
		err error
		config Config
		content string
		zapLogger *ZapLogger
	)

	data, err = ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	content = string(data)	
	config, err = parseConfigString(&content)
	
	zapLogger, err = NewZapLogger(config.Program.LogFilePath)
	logger = zapLogger

	if err != nil {
		panic(err)
	}
	defer logger.Destroy()

	forever := make(chan bool)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer conn.Close()

	exchangeListener := RabbitMQExchangeListener{
		connection: conn,
		config: &config,
	}
	exchangeListener.Init()
	exchangeListener.DeclareResources()
	exchangeListener.Listen(&config)
	defer exchangeListener.Destroy()
	
	<-forever
}
