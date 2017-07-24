package main

import (
	"os"

	"github.com/takashabe/go-display"
)

func main() {
	cli := &display.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}

	os.Exit(cli.Run(os.Args))
}
