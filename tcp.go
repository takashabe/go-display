package display

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
		go p.readListener(l)
		<-ctx.Done()
		wg.Done()
	}()

	fmt.Println("Listening TCP connection at " + l.Addr().String())
	wg.Wait()
	return nil
}

func (p *PrintTCP) readListener(l net.Listener) error {
	for {
		var (
			errCnt = 0
			maxCnt = 20
		)

		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if errCnt < maxCnt {
					time.Sleep(10 * time.Millisecond)
					continue
				}
			}
			fmt.Fprintf(p.errStream, "%s\n", err.Error())
			return err
		}
		io.Copy(p.outStream, conn)
		conn.Close()
	}
}
