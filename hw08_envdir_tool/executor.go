package main

import (
	"fmt"
	"os"
	"os/exec"
)

const (
	RetOk  = 0
	RetErr = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if err := env.UpdateOsEnv(); err != nil {
		return RetErr
	}

	args := make([]string, 0)
	if len(cmd) > 1 {
		args = cmd[1:]
	}
	command := exec.Command(cmd[0], args...) //nolint:gosec

	command.Stdout = os.Stdout
	if err := command.Run(); err != nil {
		fmt.Println(err)
		return RetErr
	}
	return RetOk
}
