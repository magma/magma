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
	"net/http/httptest"
	"testing"

	"fbc/lib/go/http/server"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCensusConfig(t *testing.T) {
	config, err := NewConfig(`{"jaeger":{"AgentEndpoint":"localhost:12345"},"xray":{"region": "eu-west-1"},"prometheus":{}}`)
	require.NotNil(t, config)
	assert.NoError(t, err)
	opts := config.ServerOptions()
	assert.Len(t, opts, 4)
	srv, err := server.New(server.Config{}, opts...)
	assert.NotNil(t, srv)
	assert.NoError(t, err)
	_, pattern := srv.Mux.Handler(httptest.NewRequest(http.MethodGet, "/metrics", nil))
	assert.NotEmpty(t, pattern)
}

func TestCensusConfigBadConfig(t *testing.T) {
	_, err := NewConfig("")
	assert.Error(t, err)
	config, err := NewConfig(`{"unknown": false}`)
	assert.NoError(t, err)
	assert.Nil(t, config.XRay)
	assert.Nil(t, config.Jaeger)
	assert.Nil(t, config.Prometheus)
}

func TestCensusConfigWithService(t *testing.T) {
	config, err := NewConfig(`{"jaeger":{},"prometheus":{}}`)
	require.NotNil(t, config)
	assert.NoError(t, err)
	assert.Empty(t, config.Jaeger.Process.ServiceName)
	assert.Empty(t, config.Prometheus.Namespace)
	service := "test"
	config = config.WithService(service)
	assert.Equal(t, service, config.Jaeger.Process.ServiceName)
	assert.Equal(t, service, config.Prometheus.Namespace)
}
