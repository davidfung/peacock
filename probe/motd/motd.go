// Package motd implements a probe to display a message of the day.

package motd

import (
	"fmt"
	"time"

	"github.com/davidfung/motd"
)

func Load() <-chan string {
	ch := make(chan string)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				msg := motd.Next()
				s := fmt.Sprintf("%s", msg)
				ch <- s
			}
		}
	}()

	return ch
}
