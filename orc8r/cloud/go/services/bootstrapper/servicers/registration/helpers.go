package registration

import (
	"math/rand"
	"time"
	"unsafe"

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
const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func generateSecureNonce(length int) string {
	src := rand.NewSource(time.Now().UnixNano())

	b := make([]byte, length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
