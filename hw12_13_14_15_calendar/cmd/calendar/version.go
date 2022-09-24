package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	release    = "UNKNOWN"
	buildDate  = "UNKNOWN"
	gitHash    = "UNKNOWN"
	versionCmd = &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}
)

func init() {
	calendarCmd.AddCommand(versionCmd)
}

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
