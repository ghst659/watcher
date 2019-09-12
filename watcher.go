// Package watcher provides a watch and report a subject.
package watcher

import (
	"context"
	"time"
)

// Reporter is a function that builds a report of strings.
type Reporter func(context.Context) ([]string, error)

// A Ticker is a mockable interface to a ticker.
type Ticker interface {
	Chan() <-chan time.Time
	Stop()
}

// A Logger is the interface required to output messages.
type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

// Watch continuously runs a reporting loop until cancelled.
func Watch(ctx context.Context, ticker Ticker, reporter Reporter, out Logger) {
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.Chan():
			report, err := reporter(ctx)
			if err != nil {
				out.Printf("%v: %v", now, err)
				return
			}
			for _, line := range report {
				out.Print(line)
			}
		}
	}
}
