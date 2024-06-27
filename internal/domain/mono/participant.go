package mono

import (
	"context"
)

type Participant interface {
	Prepare(ctx context.Context, transactionID string) error
	Commit(ctx context.Context, transactionID string) error
	Rollback(ctx context.Context, transactionID string) error
}
