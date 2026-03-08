package app

import (
	"context"
	"log"
	"time"

	"github.com/justEstif/todo-open/internal/core"
)

// StartLeaseSweeper starts a background goroutine that periodically sweeps
// expired agent leases. It stops when ctx is cancelled.
func StartLeaseSweeper(ctx context.Context, svc core.TaskService, interval time.Duration) {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				n, err := svc.SweepExpiredLeases(ctx)
				if err != nil {
					log.Printf("lease sweeper error: %v", err)
					continue
				}
				if n > 0 {
					log.Printf("lease sweeper: expired %d task lease(s), transitioned back to open", n)
				}
			}
		}
	}()
}
