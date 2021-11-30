package protos

import "time"

func (ti *TokenInfo) IsExpired() bool {
	return time.Now().Before(time.Unix(0, int64(ti.Timeout.Nanos)))
}
