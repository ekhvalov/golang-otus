package main

import (
	"fmt"
	"os"
)

func main() {
	if err := calendarCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
