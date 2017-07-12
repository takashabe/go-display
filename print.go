package printserver

import (
	"context"
)

// Printer represents a protocol and listen methods
type Printer interface {
	Protocol() string
	Listen(ctx context.Context, addr string) error
}
