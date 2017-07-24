package display

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Printer represents a protocol and listen methods
type Printer interface {
	Protocol() string
	Listen(ctx context.Context, addr string) error
}

func listenSignal(cancel context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()
}
