package printserver

import (
	"bytes"
	"net"
	"net/http"
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

func writeHTTPClient(t *testing.T, url string, message []byte) {
	req, err := http.NewRequest("GET", url, bytes.NewReader(message))
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("want no error, got %v", err)
	}
}
