package main

import (
	"os"

	"github.com/hanjm/natsbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
