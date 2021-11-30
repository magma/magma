package registration

import (
	"math/rand"

	"magma/orc8r/cloud/go/services/bootstrapper"
)

func nonceToToken(nonce string) string {
	return bootstrapper.TokenPrefix + nonce
}

func nonceFromToken(token string) string {
	return token[len(bootstrapper.TokenPrefix):]
}

// generateNonce is sourced from https://stackoverflow.com/a/31832326
func generateNonce(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
