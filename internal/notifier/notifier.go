package notifier

import "context"

// Sender defines a minimal notification interface.
type Sender interface {
	Send(ctx context.Context, text string) error
}
