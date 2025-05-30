package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1 // Return error code if no command is provided
	}

	var command *exec.Cmd
	if len(cmd) == 1 {
		//nolint
		command = exec.Command(cmd[0])
	} else {
		//nolint
		command = exec.Command(cmd[0], cmd[1:]...)
	}

	// Set up environment variables
	command.Env = os.Environ()
	for key, value := range env {
		if value.NeedRemove {
			command.Env = removeEnvVar(command.Env, key)
		} else {
			command.Env = append(command.Env, key+"="+value.Value)
		}
	}

	// Set the command's standard input, output, and error to the current process
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Run the comman
	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode() // Return the exit code from the command
		}

		log.Fatal(err)
		return 1 // Return a generic error code for other errors
	}

	return 0 // Return 0 if the command runs successfully
}

func removeEnvVar(env []string, key string) []string {
	var result []string
	prefix := key + "="
	for _, e := range env {
		if !strings.HasPrefix(e, prefix) {
			result = append(result, e)
		}
	}
	return result
}
