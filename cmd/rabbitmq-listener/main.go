package main

import (
	"io/ioutil"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func newLogger(path string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	
	if path != "" {
		cfg.OutputPaths = []string{
			path,
		}
	}
	return cfg.Build()
}

var (
	logger *zap.Logger
)

func main() {
	var (
		data []byte
		err error
		config Config
		content string
	)

	data, err = ioutil.ReadFile("./config.yml")
	if err != nil {
		panic(err)
	}
	content = string(data)	
	config, err = parseConfigString(&content)

	
	logger, err = newLogger(config.Program.LogFilePath)
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

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
