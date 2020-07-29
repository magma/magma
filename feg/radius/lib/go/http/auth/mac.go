/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fbc/lib/go/http/httputil"

	"github.com/pkg/errors"
)

const (
	// AuthorizationHeader is the header that holds the signature described below.
	AuthorizationHeader = "Authorization"
	// TimestampHeader is the header that holds the timestamp that used to generate the Authorization.
	// Read more about it below.
	TimestampHeader = "X-Authorization-Timestamp"
)

// ErrMissingHeader returns when one of the headers above is missing.
type ErrMissingHeader struct {
	header string
}

func (e ErrMissingHeader) Error() string {
	return fmt.Sprintf("auth: missing header %q", e.header)
}

// MACConfig is configuration for the MAC middleware.
type MACConfig struct {
	// KeyGetter is the function for getting the public key.
	KeyGetter func() (*rsa.PublicKey, error)
	// Hash is a crypto.Hash strategy (sha1, sha256, ..).
	Hash crypto.Hash
	// TimeSkew means how much time-skew we allow between the client and server. Defaults to 1 minute,
	// which in practice translates to a maximum of 2 minutes as the skew can be positive or negative.
	TimeSkew time.Duration
	// MaxBodySize prevents clients from accidentally or maliciously sending a large requests. It defaults to 1MB.
	MaxBodySize int64
	// Log is an optional logging function.
	Log func(string, ...interface{})
}

// DefaultMACConfig is the default configuration for the MAC middleware.
var DefaultMACConfig = MACConfig{
	Log:         func(string, ...interface{}) {},
	Hash:        crypto.SHA256,
	TimeSkew:    time.Minute,
	MaxBodySize: 1 << 20,
}

func (c *MACConfig) defaults() {
	if c.Log == nil {
		c.Log = DefaultMACConfig.Log
	}
	if c.Hash == 0 {
		c.Hash = DefaultMACConfig.Hash
	}
	if c.TimeSkew == 0 {
		c.TimeSkew = DefaultMACConfig.TimeSkew
	}
	if c.MaxBodySize == 0 {
		c.MaxBodySize = DefaultMACConfig.MaxBodySize
	}
}

// MACMiddleware is middleware for HTTP authentication using a message authentication code (MAC).
// Note, that it does not support multi-tenant, which means you can only receive requests (under
// this middleware) only from one "site". In order to make it make multi-tenant, you need to add
// an "id" parameter to the header format above and choose the public key (or shared secret) based
// on it. You can read more about it here:
// 1. https://tools.ietf.org/html/draft-hammer-oauth-v2-mac-token-05
// 2. https://github.com/hueniverse/hawk
//
// For future readers: we decided to use key-pair and not shared-secret, because we prefer to manage
// secrets only in one site (WWW).
type MACMiddleware struct {
	MACConfig
}

// NewMACMiddleware creates a new MAC handler.
func NewMACMiddleware(c MACConfig) func(http.Handler) http.Handler {
	c.defaults()
	m := &MACMiddleware{
		MACConfig: c,
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, m.MaxBodySize)
			if err := m.verify(r); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				m.Log("auth: failed to authorize request: %v", err)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

// verify decided if we pass or reject the request. Here's an example flow of the authentication:
// 1. Client calculates a checksum for the concatenation of the given values: method (POST/PUT/),
//    resource (/resource/1?a=1), unix timestamp (1353832234) and request payload.
// 2. Then, takes the hash created in #1 calculate the signature using RSASSA-PKCS1-V1_5-SIGN (with his private key),
//    encode it into a base64, store it in the Authorization header and the timestamp used in the X-Authorization-Timestamp header.
// 3. Client exposes its public key (it will be resolved using the `KeyGetter` function).
// 4. When the server gets the request, it first makes sure that the time skew is not over the limit (avoid replays).
// 5. Decodes the Authorization header and verifies the signature using the public key and SHA1 digest.
func (m *MACMiddleware) verify(r *http.Request) error {
	var (
		ts = r.Header.Get(TimestampHeader)
		sg = r.Header.Get(AuthorizationHeader)
	)
	if ts == "" {
		return ErrMissingHeader{TimestampHeader}
	}
	if sg == "" {
		return ErrMissingHeader{AuthorizationHeader}
	}
	n, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return errors.Wrapf(err, "parse timestamp: %v", ts)
	}
	if ts := time.Unix(n, 0); abs(time.Since(ts)) > m.TimeSkew {
		return errors.Errorf("stable timestamp: %v", ts)
	}
	h := m.Hash.New()
	if err := hashRequest(h, r, ts); err != nil {
		return err
	}
	pk, err := m.KeyGetter()
	if err != nil {
		return errors.Wrap(err, "failed to get public key")
	}
	sig, err := base64.StdEncoding.DecodeString(sg)
	if err != nil {
		return errors.Wrap(err, "failed to decode signature header")
	}
	if err := rsa.VerifyPKCS1v15(pk, m.Hash, h.Sum(nil), sig); err != nil {
		return errors.New("unauthorized")
	}
	return nil
}

// MACTransport is the RoundTripper for signing ourgoing requests.
type MACTransport struct {
	Transport  http.RoundTripper
	Hash       crypto.Hash
	PrivateKey *rsa.PrivateKey
}

// RoundTrip signs and changes the outgoing request and add the authorization headers to it.
func (m *MACTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	var (
		h  = m.Hash.New()
		ts = strconv.Itoa(int(time.Now().Unix()))
	)
	r = httputil.CloneRequest(r)
	if err := hashRequest(h, r, ts); err != nil {
		return nil, err
	}
	sig, err := rsa.SignPKCS1v15(rand.Reader, m.PrivateKey, m.Hash, h.Sum(nil))
	if err != nil {
		return nil, err
	}
	r.Header.Add(TimestampHeader, ts)
	r.Header.Add(AuthorizationHeader, base64.StdEncoding.EncodeToString(sig))
	return m.Transport.RoundTrip(r)
}

// hashRequest copies the request content into the Hash interface.
func hashRequest(h hash.Hash, r *http.Request, ts string) error {
	b := new(bytes.Buffer)
	b.WriteString(strings.ToUpper(r.Method))
	b.WriteString(r.URL.RequestURI())
	b.WriteString(ts)
	if _, err := io.Copy(h, b); err != nil {
		return errors.Wrap(err, "failed to copy request info")
	}
	if r.Body != nil {
		b.Reset()
		if _, err := io.Copy(io.MultiWriter(h, b), r.Body); err != nil {
			return errors.Wrap(err, "failed to copy request body")
		}
		r.Body = ioutil.NopCloser(b)
	}
	return nil
}

func abs(d time.Duration) time.Duration {
	if d < 0 {
		return -d
	}
	return d
}
