package printserver

import (
	"reflect"
	"strings"
	"testing"
)

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
		cli := &CLI{}
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
