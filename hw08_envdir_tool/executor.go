package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// nolint
	// suppress gosec G204
	c := exec.Command(cmd[0], cmd[1:]...)
	for k, e := range env {
		if e.NeedRemove {
			_ = os.Unsetenv(k)
		} else {
			_ = os.Setenv(k, e.Value)
		}
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Start()
	if err != nil {
		panic(err)
	}
	err = c.Wait()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}
		panic(err)
	}
	return 0
}
