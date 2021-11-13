package certifier

import (
	"crypto/rand"
	"errors"
	"hash/crc32"

	"github.com/jxskiss/base62"
)

type TokenPrefix string

const (
	Personal TokenPrefix = "op"
)

func GenerateToken(prefix TokenPrefix) (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	checksum := crc32.ChecksumIEEE(bytes)
	checksumBytes := i32tob(checksum)

	finalBytes := append(bytes, checksumBytes...)
	token := base62.StdEncoding.EncodeToString(finalBytes)

	return string(prefix) + "_" + token, nil
}

func ValidateTokenChecksum(token string) error {
	token = stripTokenHeader(token)
	bytes, err := base62.DecodeString(token)
	if err != nil {
		return err
	}
	orignalChecksum := btoi32(bytes[len(bytes)-4:])
	newChecksum := crc32.ChecksumIEEE(bytes[:len(bytes)-4])
	if orignalChecksum != newChecksum {
		return errors.New("invalid token")
	}
	return nil
}

func stripTokenHeader(token string) string {
	return token[3:]
}

func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

func btoi32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 4; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}
