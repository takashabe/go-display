package printserver

import (
	"context"
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
	protocol string
	addr     string
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

	var printer Printer

	switch param.protocol {
	case "http":
		printer = &PrintHTTP{
			outStream: c.OutStream,
			errStream: c.ErrStream,
			localAddr: make(chan string, 1),
		}
	case "tcp":
		printer = &PrintTCP{
			outStream: c.OutStream,
			errStream: c.ErrStream,
			localAddr: make(chan string, 1),
		}
	case "udp":
		printer = &PrintUDP{
			outStream: c.OutStream,
			errStream: c.ErrStream,
			localAddr: make(chan string, 1),
		}
	default:
		fmt.Fprintf(c.ErrStream, "unknown protocol")
		return ExitCodeInvalidArgsError
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listenSignal(cancel)
	if err := printer.Listen(ctx, param.addr); err != nil {
		fmt.Fprintf(c.ErrStream, "failed to Listen at %s protocol", printer.Protocol())
		return ExitCodeError
	}

	return ExitCodeOK
}

func (c *CLI) parseArgs(args []string, p *param) error {
	flags := flag.NewFlagSet("param", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)

	flags.StringVar(&p.protocol, "protocol", defaultProto, "Network protocol. By default, http. Choose from http, tcp and udp.")
	flags.StringVar(&p.addr, "addr", defaultAddr, "Listen port. By default, localhost:0 (:0 means, return a free port).")

	return flags.Parse(args)
}
