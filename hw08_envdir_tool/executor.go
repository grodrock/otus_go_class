package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {

	command := exec.Command(cmd[0])
	command.Env = append(os.Environ())
	fmt.Println("command", command)
	if len(cmd) > 1 {
		command.Args = cmd[1:]
		fmt.Println("command.Args", command.Args)
	}
	command.Stdout = os.Stdout
	if err := command.Run(); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}
