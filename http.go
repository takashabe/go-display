package printserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	router "github.com/takashabe/go-router"
)

// PrintHTTP represents printer for HTTP protocol
type PrintHTTP struct {
	outStream io.Writer
	errStream io.Writer
	localAddr chan string
}

// Protocol return supported protocol
func (p *PrintHTTP) Protocol() string { return "http" }

// Listen are to wait receive a HTTP request and display that
func (p *PrintHTTP) Listen(ctx context.Context, addr string) error {
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
		p.localAddr <- "http://" + l.Addr().String()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		s := &http.Server{
			Handler: p.routes(),
		}
		defer s.Close()

		go s.Serve(l)
		<-ctx.Done()
		wg.Done()
	}()

	fmt.Println("Listening HTTP request at " + l.Addr().String())
	wg.Wait()
	return nil
}

func (p *PrintHTTP) routes() *router.Router {
	r := router.NewRouter()
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		io.Copy(p.outStream, req.Body)
		w.WriteHeader(200)
	})
	return r
}
