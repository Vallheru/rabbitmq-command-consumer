package main

import (
	"testing"
	"net/url"
	"bytes"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPrepareCommand(t *testing.T) {
	commands := []string{
		"test",
		"[ -e /test.log ] && echo \"OK\";",
		"",
	}

	for _, command := range commands {
		cmd, err := PrepareCommand(command, nil)
		assert.Equal(t, cmd.command, command)
		assert.Greater(t, cmd.id, uint32(0))
		assert.Nil(t, err)
	}
}

type MemorySink struct {
    *bytes.Buffer
}
func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }


func getLoggerMock(sink *MemorySink) (*zap.Logger, error) {
    zap.RegisterSink("memory", func(*url.URL) (zap.Sink, error) {
        return sink, nil
	})
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{"memory://"}
	
	return conf.Build()
}

func TestCommandLog(t *testing.T) {
	
	commands := []struct{
		command string
		output string
	}{
		{
			"echo 'ABCD'",
			"ABC",
		},
		{
			"echo 'OK'",
			"OK",
		},
	}

	sink := &MemorySink{new(bytes.Buffer)}
	for _, item := range commands {

		logger, err := getLoggerMock(sink)
		if err != nil {
			t.Fatal(err)
		}

		cmd, err1 := PrepareCommand(item.command, logger)
		if err1 != nil {
			t.Fatal(err1)
		}
		cmd.Execute()
		output := sink.String()
		assert.Contains(t, output, item.output)
		assert.Equal(t, cmd.command, item.command)
		assert.Greater(t, cmd.id, uint32(0))
		assert.Nil(t, err)
	}
    

}