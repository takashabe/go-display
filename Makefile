.PHONY: build test clean

build: cmd/display/main.go
	cd cmd/display && go build -a

test:
	go test -v ./...

clean:
	rm cmd/display/display
