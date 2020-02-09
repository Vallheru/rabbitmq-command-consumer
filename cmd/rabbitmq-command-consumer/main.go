package main

import (
	"flag"
	"io/ioutil"
	"github.com/streadway/amqp"
	"os"
	"fmt"
)

var (
	logger Logger
)

func main() {
	configPath := flag.String("config", "/etc/rabbitmq-command-consumer.yml", "Path to the config file")
	flag.Parse()

	var (
		data []byte
		err error
		config Config
		content string
		zapLogger *ZapLogger
	)

	data, err = ioutil.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}
	content = string(data)
	content = os.ExpandEnv(content)
	config, err = parseConfigString(&content)
	
	zapLogger, err = NewZapLogger(config.Program.LogFilePath)
	logger = zapLogger

	if err != nil {
		panic(err)
	}
	defer logger.Destroy()

	forever := make(chan bool)
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s",
		config.RabbitMQ.User,
		config.RabbitMQ.Pass,
		config.RabbitMQ.Host,
		config.RabbitMQ.Port,
	))

	if err != nil {
		logger.Error("main", "Cannot connect to the RabbitMQ: ", "error", err.Error())
		os.Exit(1)
	}


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
