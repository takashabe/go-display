package main

import (
	"os"

	printserver "github.com/takashabe/go-printserver"
)

func main() {
	cli := &printserver.CLI{
		OutStream: os.Stdout,
		ErrStream: os.Stderr,
	}

	os.Exit(cli.Run(os.Args))
}
