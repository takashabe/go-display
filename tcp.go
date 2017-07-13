package printserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// PrintTCP represents printer for TCP protocol
type PrintTCP struct {
	outStream io.Writer
	errStream io.Writer
	interval  time.Duration
	localAddr chan string
}

// Protocol return supported protocol
func (p *PrintTCP) Protocol() string { return "tcp" }

// Listen are to wait receive a TCP packet and display that
func (p *PrintTCP) Listen(ctx context.Context, addr string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	defer l.Close()

	if p.localAddr != nil {
		p.localAddr <- l.Addr().String()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		t := time.NewTicker(p.interval)
		for {
			select {
			case <-t.C:
				l.SetDeadline(time.Now().Add(10 * time.Millisecond))
				conn, err := l.Accept()
				if err != nil {
					// skip when timeout
					if op, ok := err.(*net.OpError); ok && op.Timeout() {
						continue
					}
					fmt.Fprintf(p.errStream, "failed to connection: %v\n", err)
					continue
				}

				bytes := make([]byte, 4096)
				_, err = conn.Read(bytes)
				if err != nil {
					fmt.Fprintf(p.errStream, "failed to read packet: %v\n", err)
				} else {
					fmt.Fprintf(p.outStream, "%s\n", bytes)
				}
			case <-ctx.Done():
				t.Stop()
				wg.Done()
				return
			}
		}
	}()

	fmt.Println("Listening TCP connection at " + l.Addr().String())
	wg.Wait()
	return nil
}
