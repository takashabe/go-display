.PHONY: build test

build: cmd/display/main.go
	cd cmd/display && go build -a

test:
	go test -v ./...
