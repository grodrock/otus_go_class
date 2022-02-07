package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Not enough args")
		os.Exit(1)
	}
	envDir := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ret := RunCmd(cmd, env)
	os.Exit(ret)
}
