package main

import (
	"fmt"
	"os"
)

func main() {
	if err := senderCmd.Execute(); err != nil {
		_, errPrint := fmt.Fprintf(os.Stderr, "%v\n", err)
		if errPrint != nil {
			panic(fmt.Sprintf("error '%v' occurred while write error '%v' into stderr", errPrint, err))
		}
		os.Exit(1)
	}
}
