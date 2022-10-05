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
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestMACMiddleware(t *testing.T) {
	tests := []struct {
		name   string
		do     func(*require.Assertions, *rsa.PrivateKey, string) *http.Response
		expect func(*require.Assertions, int, string)
	}{
		{
			name: "POST/with valid cert",
			do: func(require *require.Assertions, pk *rsa.PrivateKey, u string) *http.Response {
				r, err := http.NewRequest(http.MethodPost, u, strings.NewReader("a8m"))
				require.NoError(err)
				res, err := testClient(pk).Do(r)
				require.NoError(err)
				return res
			},
			expect: func(require *require.Assertions, code int, body string) {
				require.Equal(http.StatusOK, code)
				require.Equal("a8m", body)
			},
		},
		{
			name: "GET/with no body",
			do: func(require *require.Assertions, pk *rsa.PrivateKey, u string) *http.Response {
				r, err := http.NewRequest(http.MethodGet, u, nil)
				require.NoError(err)
				res, err := testClient(pk).Do(r)
				require.NoError(err)
				return res
			},
			expect: func(require *require.Assertions, code int, body string) {
				require.Equal(http.StatusOK, code)
				require.Empty(body)
			},
		},
		{
			name: "GET/with query string",
			do: func(require *require.Assertions, pk *rsa.PrivateKey, u string) *http.Response {
				r, err := http.NewRequest(http.MethodGet, u+"?foo=bar&baz=qux", nil)
				require.NoError(err)
				res, err := testClient(pk).Do(r)
				require.NoError(err)
				return res
			},
			expect: func(require *require.Assertions, code int, body string) {
				require.Equal(http.StatusOK, code)
				require.Empty(body)
			},
		},
		{
			name: "invalid cert",
			do: func(require *require.Assertions, pk *rsa.PrivateKey, u string) *http.Response {
				r, err := http.NewRequest(http.MethodGet, u, nil)
				require.NoError(err)
				res, err := http.DefaultClient.Do(r)
				require.NoError(err)
				return res
			},
			expect: func(require *require.Assertions, code int, body string) {
				require.Equal(http.StatusUnauthorized, code)
				require.Equal("auth: missing header \"X-Authorization-Timestamp\"\n", body)
			},
		},
		{
			name: "stale timestamp",
			do: func(require *require.Assertions, pk *rsa.PrivateKey, u string) *http.Response {
				r, err := http.NewRequest(http.MethodGet, u, nil)
				require.NoError(err)
				r.Header.Add(AuthorizationHeader, "boring")
				r.Header.Add(TimestampHeader, fmt.Sprint(time.Now().Add(-time.Hour).Unix()))
				res, err := http.DefaultClient.Do(r)
				require.NoError(err)
				return res
			},
			expect: func(require *require.Assertions, code int, body string) {
				require.Equal(http.StatusUnauthorized, code)
				require.Contains(body, "stable timestamp")
			},
		},
		{
			name: "invalid auth header",
			do: func(require *require.Assertions, pk *rsa.PrivateKey, u string) *http.Response {
				r, err := http.NewRequest(http.MethodGet, u, nil)
				require.NoError(err)
				r.Header.Add(AuthorizationHeader, "boring")
				r.Header.Add(TimestampHeader, fmt.Sprint(time.Now().Unix()))
				res, err := http.DefaultClient.Do(r)
				require.NoError(err)
				return res
			},
			expect: func(require *require.Assertions, code int, body string) {
				require.Equal(http.StatusUnauthorized, code)
				require.Equal("failed to decode signature header: illegal base64 data at input byte 4\n", body)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			pk, err := rsa.GenerateKey(rand.Reader, 2048)
			require.NoError(err)
			mw := NewMACMiddleware(MACConfig{
				Log: t.Logf,
				KeyGetter: func() (*rsa.PublicKey, error) {
					return &pk.PublicKey, nil
				},
			})
			ts := httptest.NewServer(mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, err := io.Copy(w, r.Body)
				require.NoError(err)
			})))
			res := tt.do(require, pk, ts.URL)
			buf, err := ioutil.ReadAll(res.Body)
			require.NoError(err)
			require.NoError(res.Body.Close())
			tt.expect(require, res.StatusCode, string(buf))
		})
	}
}

// testClient returns an http client with the MACTransport.
func testClient(pk *rsa.PrivateKey) *http.Client {
	return &http.Client{
		Transport: &MACTransport{
			Hash:       crypto.SHA256,
			PrivateKey: pk,
			Transport:  http.DefaultTransport,
		},
	}
}
