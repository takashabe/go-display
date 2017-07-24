# go-display

Display is a server that displays messages sent via several protocols.

## Installation

```
go get -u github.com/takashabe/go-display/cmd/display
```

## Usage

```
display [options]
```

Options

* `-addr="localhost:8080"`: Listen address. By default, localhost:0 (:0 means, return a free port)
* `-protocol="udp"`: Network protocol. By default, http. Choose from http, tcp and udp

Usage example

```
display -addr="localhost:8080" -protocol="udp"
```
