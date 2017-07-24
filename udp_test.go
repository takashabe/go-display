package display

import (
	"context"
	"io"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestWriteUDP(t *testing.T) {
	pr, pw := io.Pipe()
	defer pr.Close()
	defer pw.Close()

	p := &PrintUDP{
		outStream: pw,
		errStream: pw,
		localAddr: make(chan string, 1),
	}

	// write to udp server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.Listen(ctx, ":0")

	// wait to write in udp server
	addr := <-p.localAddr
	sendMessage := []byte("test")
	time.Sleep(10 * time.Millisecond)
	writeClient(t, p.Protocol(), addr, sendMessage)
	time.Sleep(10 * time.Millisecond)

	receivedMessage := make([]byte, len(sendMessage))
	_, err := pr.Read(receivedMessage)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
	if !reflect.DeepEqual(receivedMessage, sendMessage) {
		t.Errorf("want message %s, got %s", sendMessage, receivedMessage)
	}
}

func TestCancelUDP(t *testing.T) {
	p := &PrintUDP{
		outStream: os.Stdout,
		errStream: os.Stderr,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go p.Listen(ctx, ":0")
	time.Sleep(50 * time.Millisecond)
	beforeCancel := runtime.NumGoroutine()

	cancel()
	time.Sleep(50 * time.Millisecond)
	afterCancel := runtime.NumGoroutine()

	if beforeCancel <= afterCancel {
		t.Errorf("want num goroutine less than %d, got %d", beforeCancel, afterCancel)
	}
}
