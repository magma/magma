/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package grant_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas/grant"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestGetFrequencyGrantMapping(t *testing.T) {
	grants := []*active_mode.Grant{{
		Id:              "some_id",
		LowFrequencyHz:  3580e6,
		HighFrequencyHz: 3590e6,
	}, {
		Id:              "other_id",
		LowFrequencyHz:  3590e6,
		HighFrequencyHz: 3610e6,
	}}
	actual := grant.GetFrequencyGrantMapping(grants)
	expected := map[int64]*active_mode.Grant{
		3585e6: grants[0],
		3600e6: grants[1],
	}
	assert.Equal(t, expected, actual)
}
