package main

import (
	"fmt"
	"os"
)

func main() {
	err := calendarCmd.Execute()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
