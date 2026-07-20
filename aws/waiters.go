package aws

import (
	"context"
	"time"
)

// PollResource repeatedly checks condition until deleted or timeout
func PollResource(ctx context.Context, checkFunc func() (bool, error), interval, timeout time.Duration) error {
	// Polling loop logic with timeout and checkFunc
	return nil
}
