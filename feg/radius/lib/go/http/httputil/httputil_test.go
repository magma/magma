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

package httputil

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCloneRequest(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	assert.NoError(t, err)

	clone := CloneRequest(req)
	assert.Equal(t, req, clone)
	assert.False(t, &req.Header == &clone.Header)
}

func TestCloneHeader(t *testing.T) {
	header := make(http.Header)
	header.Set("Content-Length", "123")
	header.Set("Content-Type", "text/plain")
	header.Set("Date", time.Now().Format(time.RFC3339))

	clone := CloneHeader(header)
	assert.False(t, &header == &clone)
	assert.Equal(t, header, clone)

	clone.Set("Content-Language", "en")
	assert.NotEqual(t, header, clone)
	assert.Len(t, clone, len(header)+1)
}
