package census

import (
	"net/http"
)

// Census defines opencensus runtime settings.
type Census struct {
	// handler for /metrics endpoint.
	StatsHandler http.Handler

	// Set of closers to run on Close.
	closers []func()
}

// Close close the census
func (c *Census) Close() {
	for _, closer := range c.closers {
		closer()
	}
}
