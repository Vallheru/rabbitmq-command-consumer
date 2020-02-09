package main

import (
	"os/exec"
	"bytes"
	"errors"

	"github.com/google/uuid"
)

// The CommandExec keeps pieces of information about command that is
// being executed in the system. During the execution, a few entries 
// are created in the logs file. To connect all rows struct keeps and 
// adds the ID field.
type CommandExec struct {
	id uint32
	command string
	logger *Logger
}

// PrepareCommand creates the CommandExec structure.
func PrepareCommand(cmdStr string, logger *Logger) (*CommandExec, error) {
	id, err := uuid.NewUUID()
    if err !=nil {
        return nil, errors.New("Cannot generat UUID")
    }

	cmd := CommandExec{
		id: id.ID(),
		command: cmdStr,
		logger: logger,
	}

	return &cmd, nil
}

// Log function is used to create Info entry in the log file
func (cmd *CommandExec) Log(message string, args ...interface{}) {
	params := []interface{}{
		"command_id", cmd.id,
		"command", cmd.command,
	}
	params = append(params, args...)

	(*cmd.logger).Info("CommandExec", message, 
		params...,
	)
}

// Execute function executes given command.
func (cmd *CommandExec) Execute() (string, string, error) {	
	cmd.Log("Start")
	
	
	cmdExec := exec.Command("/bin/bash", "-c", cmd.command)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmdExec.Stdout = &stdout
    cmdExec.Stderr = &stderr
	err := cmdExec.Start()
	if err != nil {
		(*cmd.logger).Error("ERROR", "error_msg", err.Error())
		return "", "", err
	}
	cmd.Log("Waiting to finish")
	err = cmdExec.Wait()

	cmd.Log("Finished", 
		"stdout", stdout.String(),
		"stderr", stderr.String(),
	)
	return stdout.String(), stderr.String(), err
}
