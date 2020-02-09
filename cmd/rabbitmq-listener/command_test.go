package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
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
	var logger Logger = &LoggerMock{}

	for _, item := range commands {
		cmd, err1 := PrepareCommand(item.command, &logger)
		if err1 != nil {
			t.Fatal(err1)
		}
		stdout, stderr, err := cmd.Execute()
		assert.Contains(t, stdout, item.output)
		assert.Equal(t, cmd.command, item.command)
		assert.Greater(t, cmd.id, uint32(0))
		assert.Empty(t, stderr)
		assert.Nil(t, err)
	}
    

}