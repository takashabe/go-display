package printserver

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func runUDP(addr string) {
	l, err := net.ListenPacket("udp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				bytes := make([]byte, 4096)
				_, _, err := l.ReadFrom(bytes)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s\n", bytes)
			}
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		sig := <-sigCh
		fmt.Println("receive signal: ", sig.String())
		wg.Done()
	}()

	fmt.Println("Listening UDP packet at " + l.LocalAddr().String())
	wg.Wait()
}
