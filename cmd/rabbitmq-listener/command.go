package main

import (
	"os/exec"
	"bytes"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CommandExec struct {
	id uint32
	command string
	logger *zap.Logger
}

func PrepareCommand(cmdStr string, logger *zap.Logger) (*CommandExec, error) {
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

func (cmd *CommandExec) Log(message string, args ...interface{}) {
	params := []interface{}{
		"command_id", cmd.id,
		"command", cmd.command,
	}
	params = append(params, args...)

	cmd.logger.Sugar().Infow("[ CommandExec ] " + message, 
		params... ,
	)
}

func (cmd *CommandExec) Execute() (error, string, string) {	
	cmd.Log("Start")
	
	
	cmdExec := exec.Command("/bin/bash", "-c", cmd.command)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmdExec.Stdout = &stdout
    cmdExec.Stderr = &stderr
	err := cmdExec.Start()
	if err != nil {
		cmd.Log("ERROR", "error_msg", err.Error())
		return err, "", ""
	}
	cmd.Log("Waiting to finish")
	err = cmdExec.Wait()

	cmd.Log("Finished", 
		"stdout", stdout.String(),
		"stderr", stderr.String(),
	)
	return nil, stdout.String(), stderr.String()
}
