// Package alive implements a probe to show that the server is alive.

package alive

import (
	"time"
)

const interval = 1 * time.Hour

func send(out chan<- string) {
	out <- "Alive"
}

func Load(out chan<- string) {
	send(out)
	go func() {
		const interval = 1 * time.Hour
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				send(out)
			}
		}
	}()
}
