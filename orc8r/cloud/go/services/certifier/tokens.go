package certifier

import (
	"crypto/rand"
	"errors"
	"fmt"
	"hash/crc32"
	"strings"

	"github.com/jxskiss/base62"
)

type TokenType string

const (
	Personal TokenType = "op"
)

// Length of checksum of token's byte array
const checksumLength = 4

// GenerateToken generates a random 32-byte token with a checksum
func GenerateToken(tokenType TokenType) (string, error) {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Generate checksum bytes
	checksum := crc32.ChecksumIEEE(bytes)
	checksumBytes := i32tob(checksum)

	// Combine to form final token
	finalBytes := append(bytes, checksumBytes...)
	token := base62.StdEncoding.EncodeToString(finalBytes)
	return fmt.Sprintf("%v_%s", tokenType, token), nil
}

// ValidateToken makes sure the token has the appropriate header and
// that the checksum is correct
func ValidateToken(token string) error {
	value, err := stripTokenHeader(token)
	if err != nil {
		return err
	}
	if err := validateTokenChecksum(value); err != nil {
		return err
	}
	return nil
}

func stripTokenHeader(token string) (string, error) {
	s := strings.Split(token, "_")
	if len(s) != 2 {
		return "", errors.New("missing token type")
	}
	tokenType, value := s[0], s[1]

	// Validate token tokenType
	switch tokenType {
	case string(Personal):
		return value, nil
	}
	return "", errors.New("invalid token type")
}

func validateTokenChecksum(token string) error {
	bytes, err := base62.DecodeString(token)
	if err != nil {
		return err
	}
	bytesLen := len(bytes)
	if bytesLen < checksumLength {
		return errors.New("token not long enough")
	}
	claimedChecksum := btoi32(bytes[bytesLen-checksumLength:])
	calculatedChecksum := crc32.ChecksumIEEE(bytes[:bytesLen-checksumLength])
	if claimedChecksum != calculatedChecksum {
		return errors.New("invalid token checksum")
	}
	return nil
}

// Implementation taken from https://gist.github.com/chiro-hiro/2674626cebbcb5a676355b7aaac4972d
func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

// Implementation taken from https://gist.github.com/chiro-hiro/2674626cebbcb5a676355b7aaac4972d
func btoi32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 4; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}
