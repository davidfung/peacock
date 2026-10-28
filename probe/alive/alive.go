// Package alive implements a probe to show that the server is
// alive.

package alive

import (
	"fmt"
	"time"
)

func Load(out chan<- string) {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				s := fmt.Sprintf("current time is %s", t.Format(time.RFC3339))
				out <- s
			}
		}
	}()
}
