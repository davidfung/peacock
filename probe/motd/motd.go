// Package motd implements a probe to display a message of the day.

package motd

import (
	"fmt"
	"time"

	"github.com/davidfung/motd"
)

const interval = 24 * time.Hour

func send(out chan<- string) {
	msg := motd.Next()
	s := fmt.Sprintf("%s", msg)
	out <- s
}

func Load(out chan<- string) {
	send(out)
	go func() {
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
