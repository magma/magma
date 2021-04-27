package radius_test

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
	rad "layeh.com/radius"

	"magma/feg/gateway/services/aaa/radius"
)

// Test data
var (
	response = []byte{
		0x02, 0x01, 0x00, 0x38, 0x15, 0xef, 0xbc, 0x7d, 0xab, 0x26, 0xcf, 0xa3, 0xdc, 0x34, 0xd9, 0xc0,
		0x3c, 0x86, 0x01, 0xa4, 0x06, 0x06, 0x00, 0x00, 0x00, 0x02, 0x07, 0x06, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0xff, 0xff, 0xff, 0xfe, 0x0a, 0x06, 0x00, 0x00, 0x00, 0x00, 0x0d, 0x06, 0x00, 0x00,
		0x00, 0x01, 0x0c, 0x06, 0x00, 0x00, 0x05, 0xdc,
	}
	secret = []byte("xyzzy5461")
)

func TestMessageAuth(t *testing.T) {
	p, err := rad.Parse(response, secret)
	assert.NoError(t, err)

	enc1 := encodeWithMessageAuthenticator(p)
	err = radius.AddMessageAuthenticatorAttr(p)
	assert.NoError(t, err)
	enc2, err := p.Encode()
	assert.NoError(t, err)
	assert.Equal(t, enc1, enc2)
}

// encodeWithMessageAuthenticator - in place radius packet MA addition & encoding for testing
func encodeWithMessageAuthenticator(p *rad.Packet) []byte {
	// Fix the size
	encoded, err := p.Encode()
	if err != nil {
		panic(err)
	}
	size := binary.BigEndian.Uint16(encoded[2:4]) + 18
	binary.BigEndian.PutUint16(encoded[2:4], size)

	// Add Message Authenticator Attribute, 0 padded, and flatten
	zeroedOutMsgAuthenticator := [16]byte{}
	allBytes := [][]byte{
		encoded[:4],
		p.Authenticator[:],
		encoded[20:],
		{80, 18},
		zeroedOutMsgAuthenticator[:],
	}
	var radiusMsg []byte
	for _, b := range allBytes {
		radiusMsg = append(radiusMsg[:], b[:]...)
	}
	// Calculate Message Authenticator & Overwrite
	hash := hmac.New(md5.New, p.Secret)
	hash.Write(radiusMsg)
	encoded = hash.Sum(radiusMsg[:len(radiusMsg)-16])

	// Re-calc the Response Authenticator
	resAuth := md5.New()
	resAuth.Write(encoded[:4])
	resAuth.Write(p.Authenticator[:])
	resAuth.Write(encoded[20:])
	resAuth.Write(p.Secret)
	resAuth.Sum(encoded[4:4:20])
	return encoded
}
