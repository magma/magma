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

package protos_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/accessd/protos"
)

func TestAccessControlDefinitions(t *testing.T) {
	for n, v := range protos.AccessControl_Permission_value {
		if v&(v-1) != 0 {
			t.Fatalf(
				"Invalid AccessControl Permission definition: %s = %d (B%b). "+
					"AccessControl Permissions must be powers of 2.", n, v, v)
		}
	}

	assert.Equal(t, "READ", protos.AccessControl_READ.ToString())
	assert.Equal(t, "WRITE", protos.AccessControl_WRITE.ToString())
	assert.Equal(t, "NONE", protos.AccessControl_NONE.ToString())

	rwStr := (protos.AccessControl_READ | protos.AccessControl_WRITE).ToString()
	assert.True(t, "READ|WRITE" == rwStr || "WRITE|READ" == rwStr)

	assert.Equal(t, "NONE", protos.AccessControl_Permission(16).ToString())
}
