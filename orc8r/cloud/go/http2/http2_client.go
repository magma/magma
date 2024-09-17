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

// Package http2 contains a minimal implementation of non-TLS http/2 server
// and client
package http2

import (
	"crypto/tls"
	"net"
	"net/http"

	"golang.org/x/net/http2"
)

// H2CClient is a http2 client supports non-SSL only
type H2CClient struct {
	*http.Client
}

// NewH2CClient creates a new h2cclient.
func NewH2CClient() *H2CClient {
	return &H2CClient{&http.Client{
		// Skip TLS dial
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
				// dial the addr from url,
				// or :80 if no legitimate ip retrieved from url
				return net.Dial(netw, addr)
			},
		},
	}}
}
