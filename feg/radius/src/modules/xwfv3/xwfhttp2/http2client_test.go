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

package xwfhttp2

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/http2"
)

func TestHttp2PostJson(t *testing.T) {

	accessToken := "test_token"
	// Starting our mock http2 server
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Checking that we have accessToken
		accessTokenQueryParam := r.URL.Query().Get("access_token")
		require.NotEmpty(t, accessTokenQueryParam, "access_token query param wasn't supplied")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"fake response string"}`)
	}))
	// Needed to enable http2 server
	ts.TLS = &tls.Config{
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		NextProtos:   []string{http2.NextProtoTLS},
	}
	ts.StartTLS()
	defer ts.Close()

	mockServer := ts.URL

	// Starting the client
	client := NewClient(accessToken)
	res, err := client.PostJSON(mockServer+"/test", map[string]string{
		"raw_data": "test",
	}, map[string]string{})

	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res)

}

func TestHttp2PostJsonFailure(t *testing.T) {

	// Starting our mock http2 server
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}))
	// Needed to enable http2 server
	ts.TLS = &tls.Config{
		CipherSuites: []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
		NextProtos:   []string{http2.NextProtoTLS},
	}
	ts.StartTLS()
	defer ts.Close()

	mockServer := ts.URL

	// Starting the client
	client := NewClient("test")
	res, err := client.PostJSON(mockServer+"/test", map[string]string{
		"raw_data": "test",
	}, map[string]string{})

	require.Error(t, err)
	require.Nil(t, res)

}
