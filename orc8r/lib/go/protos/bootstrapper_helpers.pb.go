package protos

import (
	"time"

	"magma/orc8r/cloud/go/clock"
)

func (t *TokenInfo) IsExpired() bool {
	expirationTime := time.Unix(t.Timeout.Seconds, int64(t.Timeout.Nanos))
	return clock.Now().After(expirationTime)
}
