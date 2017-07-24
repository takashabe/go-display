package display

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
	localAddr chan string
}

// Protocol return supported protocol
func (p *PrintUDP) Protocol() string { return "udp" }

// Listen are to wait receive a UDP packet and display that
func (p *PrintUDP) Listen(ctx context.Context, addr string) error {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	if p.localAddr != nil {
		p.localAddr <- conn.LocalAddr().String()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		go p.readConn(conn)
		<-ctx.Done()
		wg.Done()
	}()

	fmt.Println("Listening UDP packet at " + conn.LocalAddr().String())
	wg.Wait()
	return nil
}

func (p *PrintUDP) readConn(l net.Conn) error {
	for {
		var (
			errCnt = 0
			maxCnt = 20
		)

		_, err := io.Copy(p.outStream, l)
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
	}
}
