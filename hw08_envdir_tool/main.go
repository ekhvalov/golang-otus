package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Printf("Usage: %s <path_to_env_dir> <command>\n", args[0])
		return
	}
	env, err := ReadDir(args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "directory '%s' read error: %v", args[1], err)
		os.Exit(1)
	}
	RunCmd(args[2:], env)
}
