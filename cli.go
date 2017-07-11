package printserver

import (
	"flag"
	"fmt"
	"io"
)

// Exit codes. used only in Run()
const (
	ExitCodeOK = 0

	// Specific error codes. begin 10-
	ExitCodeError = 10 + iota
	ExitCodeParseError
	ExitCodeInvalidArgsError
)

// default parameters
const (
	defaultProto = "http"
	defaultAddr  = "localhost:0"
)

type param struct {
	proto string
	addr  string
}

// CLI is the command line interface object
type CLI struct {
	OutStream io.Writer
	ErrStream io.Writer
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	// parse args
	param := &param{}
	err := c.parseArgs(args[1:], param)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "args parse error: %v", err)
		return ExitCodeParseError
	}

	switch param.proto {
	case "http", "tcp":
		fmt.Fprintf(c.ErrStream, "yet not support protocol")
	case "udp":
		runUDP(param.addr)
	default:
		fmt.Fprintf(c.ErrStream, "unknown protocol")
		return ExitCodeInvalidArgsError
	}

	return ExitCodeOK
}

func (c *CLI) parseArgs(args []string, p *param) error {
	flags := flag.NewFlagSet("param", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)

	flags.StringVar(&p.proto, "proto", defaultProto, "Network protocol. By default, http. Choose from http, tcp and udp.")
	flags.StringVar(&p.addr, "addr", defaultAddr, "Listen port. By default, localhost:0 (:0 means, return a free port).")

	return flags.Parse(args)
}
