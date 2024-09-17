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

package oc

import (
	"net/http"
	"sync"

	"go.opencensus.io/plugin/ochttp"
)

// ClientCache caches per operation tracing http clients.
type ClientCache struct {
	clients sync.Map
}

var (
	// DefaultTransport is the default tracing transport and is used by DefaultClient.
	DefaultTransport http.RoundTripper = &ochttp.Transport{}

	// DefaultClient is the default tracing http client.
	DefaultClient = &http.Client{
		Transport: DefaultTransport,
	}

	// Per operation global tracing client cache.
	clientCache = &ClientCache{}
)

// ClientFor returns a tracing http client from global client cache.
func ClientFor(operation string) *http.Client {
	return clientCache.ClientFor(operation)
}

// ClientFor returns a tracing http client for operation.
func (cc *ClientCache) ClientFor(operation string) *http.Client {
	if client, ok := cc.clients.Load(operation); ok {
		return client.(*http.Client)
	}

	transport := &ochttp.Transport{
		FormatSpanName: func(*http.Request) string { return operation },
	}
	client, _ := cc.clients.LoadOrStore(operation, &http.Client{Transport: transport})
	return client.(*http.Client)
}
