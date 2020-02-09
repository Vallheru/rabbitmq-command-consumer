package main

import (
	"testing"
	"net/url"
	"bytes"

	"go.uber.org/zap"
	"github.com/stretchr/testify/assert"
)

type MemorySink struct {
    *bytes.Buffer
}
func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

func getZapLoggerMock(sink *MemorySink) *ZapLogger {
    zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
        return sink, nil
	})
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{"memory://"}
	
	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}

	return &ZapLogger{log: logger}
}

func TestZapLoggerWrapper(t *testing.T) {
	logs := []struct{
		module string
		message string
		args []interface{}
	}{
		{
			"logger_test",
			"Message 1",
			[]interface{}{
				"param1", "val1",
				"param2", "val2",
			},
		},
		{
			"logger_test",
			"Message 2",
			[]interface{}{},
		},
	}

	sink := &MemorySink{new(bytes.Buffer)}
	logger := getZapLoggerMock(sink)
	for _, item := range logs {
		logger.Info(item.module, item.message, item.args...)

		output := sink.String()
		assert.Contains(t, output, item.message)
		assert.Contains(t, output, "\"level\":\"info\"")
		sink.Reset()


		logger.Warn(item.module, item.message, item.args...)
		output = sink.String()
		assert.Contains(t, output, item.message)
		assert.Contains(t, output, "\"level\":\"warn\"")
		sink.Reset()


		logger.Error(item.module, item.message, item.args...)
		output = sink.String()
		assert.Contains(t, output, item.message)
		assert.Contains(t, output, "\"level\":\"error\"")
		assert.Contains(t, output, "stacktrace")
		sink.Reset()
	}
}