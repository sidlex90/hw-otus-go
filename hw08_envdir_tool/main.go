package main

import (
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		os.Stderr.WriteString("Usage: envdir <path_to_env> <command> <arg1> <arg2>....\n")
		os.Exit(1)
	}

	pathToEnv := args[0]
	commandData := args[1:]

	env, err := ReadDir(pathToEnv)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	RunCmd(commandData, env)

	return
}
