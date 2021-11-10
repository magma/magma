/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers_test

import (
	"magma/feg/gateway/services/n7_proxy/servicers"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	INVALID_URL        = "magam-feg/api/root"
	TEST_URL           = "https://magam-feg/api/root"
	SRV_STR            = "https://magam-feg"
	TEST_URL_WITH_PORT = "https://magma-feg:8080/api/root"
	SRV_STR_WITH_PORT  = "https://magma-feg:8080"
)

func TestGetServerStringFromUrl(t *testing.T) {
	_, err := servicers.GetServerStringFromUrl(INVALID_URL)
	assert.Error(t, err)

	testUrls := [][]string{
		{"https://magam-feg/api/root", "https://magam-feg"},
		{"https://magma-feg:8080/api/root", "https://magma-feg:8080"},
	}

	for _, testUrl := range testUrls {
		srvStr, err := servicers.GetServerStringFromUrl(testUrl[0])
		assert.NoError(t, err)
		assert.Equal(t, testUrl[1], srvStr)
	}
}
