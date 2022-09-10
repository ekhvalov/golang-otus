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
		panic(err)
	}
	RunCmd(args[2:], env)
}
