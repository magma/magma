package registration

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/lib/go/protos"
)

func nonceToToken(nonce string) string {
	return bootstrapper.TokenPrepend + nonce
}

func nonceFromToken(token string) string {
	return token[len(bootstrapper.TokenPrepend):]
}

func generateSecureNonce(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func tokenTimedOut(info *protos.TokenInfo) bool {
	return time.Now().Before(time.Unix(0, int64(info.Timeout.Nanos)))
}
