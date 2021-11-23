package registration

import (
	"math/rand"
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

func isTokenExpired(info *protos.TokenInfo) bool {
	return time.Now().Before(time.Unix(0, int64(info.Timeout.Nanos)))
}

// ========================================================================= //
// Sourced from https://stackoverflow.com/a/31832326
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateNonce(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
