// Package alive implements a probe to show that the server is alive.
//
// Loading the probe returns a read-only string channel.
package alive

import (
	"fmt"
	"time"
)

func Load() <-chan string {
	ch := make(chan string)

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				s := fmt.Sprintf("data: current time is %s\n\n", t.Format(time.RFC3339))
				ch <- s
			}
		}
	}()

	return ch
}
