package printserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// PrintUDP represents printer for UDP protocol
type PrintUDP struct {
	outStream io.Writer
	errStream io.Writer
	interval  time.Duration
	localAddr chan string
}

// Protocol return supported protocol
func (p *PrintUDP) Protocol() string { return "udp" }

// Listen are to wait receive a UDP packet and display that
func (p *PrintUDP) Listen(ctx context.Context, addr string) error {
	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	if p.localAddr != nil {
		p.localAddr <- l.LocalAddr().String()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		t := time.NewTicker(p.interval)
		for {
			select {
			case <-t.C:
				l.SetDeadline(time.Now().Add(10 * time.Millisecond))
				bytes := make([]byte, 4096)
				_, _, err := l.ReadFrom(bytes)
				if err != nil {
					// skip when timeout
					if op, ok := err.(*net.OpError); ok && op.Timeout() {
						continue
					}
					fmt.Fprintf(p.errStream, "failed to read packet: %v\n", err)
					continue
				}
				fmt.Fprintf(p.outStream, "%s\n", bytes)
			case <-ctx.Done():
				t.Stop()
				wg.Done()
				return
			}
		}
	}()

	fmt.Println("Listening UDP packet at " + l.LocalAddr().String())
	wg.Wait()
	return nil
}
