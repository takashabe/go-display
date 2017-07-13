package printserver

import (
	"bytes"
	"context"
	"os"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestWriteTCP(t *testing.T) {
	var buf bytes.Buffer
	p := &PrintTCP{
		outStream: &buf,
		errStream: &buf,
		interval:  10 * time.Millisecond,
		localAddr: make(chan string, 1),
	}

	// write to server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go p.Listen(ctx, ":0")

	// wait to write in server
	addr := <-p.localAddr
	sendMessage := []byte("test")
	time.Sleep(10 * time.Millisecond)
	writeClient(t, p.Protocol(), addr, sendMessage)
	time.Sleep(10 * time.Millisecond)

	receivedMessage := make([]byte, len(sendMessage))
	_, err := buf.Read(receivedMessage)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
	if !reflect.DeepEqual(receivedMessage, sendMessage) {
		t.Errorf("want message %s, got %s", sendMessage, receivedMessage)
	}
}

func TestCancelTCP(t *testing.T) {
	p := &PrintTCP{
		outStream: os.Stdout,
		errStream: os.Stderr,
		interval:  10 * time.Millisecond,
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
