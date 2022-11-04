package main

import (
	"fmt"
	"os"
)

func main() {
	if err := migrateCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
