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

package eirp_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/dp/cloud/go/active_mode_controller/internal/message_generator/sas/eirp"
	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func TestCalcLowerBound(t *testing.T) {
	ec := &active_mode.EirpCapabilities{
		MinPower:      0,
		MaxPower:      20,
		NumberOfPorts: 2,
	}
	c := eirp.NewCalculator(15, ec)

	actual := c.CalcLowerBound(10 * 1e6)
	assert.Equal(t, 9.0, actual)
}

func TestCalcUpperBound(t *testing.T) {
	ec := &active_mode.EirpCapabilities{
		MinPower:      0,
		MaxPower:      20,
		NumberOfPorts: 2,
	}
	c := eirp.NewCalculator(15, ec)

	actual := c.CalcUpperBound(10 * 1e6)
	assert.Equal(t, 28.0, actual)
}
