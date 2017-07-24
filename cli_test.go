package display

import (
	"bytes"
	"reflect"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"
)

func DummyCLI(t *testing.T) *CLI {
	return &CLI{
		OutStream: new(bytes.Buffer),
		ErrStream: new(bytes.Buffer),
	}
}

func TestParseArgs(t *testing.T) {
	cases := []struct {
		arg    string
		expect *param
	}{
		{
			"display -addr=:0 -protocol=tcp",
			&param{
				protocol: "tcp",
				addr:     ":0",
			},
		},
		{
			"display",
			&param{
				protocol: "http",
				addr:     "localhost:0",
			},
		},
	}
	for i, c := range cases {
		cli := DummyCLI(t)
		param := &param{}
		args := strings.Split(c.arg, " ")
		err := cli.parseArgs(args[1:], param)
		if err != nil {
			t.Fatalf("#%d: want no error, got %v", i, err)
		}

		if !reflect.DeepEqual(c.expect, param) {
			t.Errorf("#%d: want %v, got %v", i, c.expect, param)
		}
	}
}

func TestRun(t *testing.T) {
	cases := []struct {
		arg    string
		expect int
	}{
		{"display -addr=:0 -protocol=tcp", ExitCodeOK},
		{"display -addr=:0 -protocol=udp", ExitCodeOK},
		{"display ", ExitCodeOK},
		{"display -unknown", ExitCodeParseError},
		{"display -addr=unknown", ExitCodeError},
		{"display -addr=localhost:80", ExitCodeError},
		{"display -protocol=unknown", ExitCodeInvalidArgsError},
	}
	for i, c := range cases {
		cli := DummyCLI(t)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			code := cli.Run(strings.Split(c.arg, " "))
			if code != c.expect {
				t.Errorf("#%d: want code %d, got %d", i, c.expect, code)
			}
			wg.Done()
		}()

		time.Sleep(100 * time.Millisecond)
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGINT); err != nil {
			t.Fatalf("#%d: want no error, got %v", i, err)
		}
		wg.Wait()
	}
}
