package printserver

import (
	"net"
	"testing"
)

func writeClient(t *testing.T, network, addr string, message []byte) {
	c, err := net.Dial(network, addr)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}

	_, err = c.Write(message)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
}
