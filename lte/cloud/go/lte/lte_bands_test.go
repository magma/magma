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

package lte_test

import (
	"testing"

	"magma/lte/cloud/go/lte"

	"github.com/stretchr/testify/assert"
)

func TestGetBand(t *testing.T) {
	expected := map[uint32]uint32{
		0:     1,
		599:   1,
		600:   2,
		749:   2,
		38650: 40,
		43590: 43,
		45589: 43,
	}

	for earfcndl, bandExpected := range expected {
		band, err := lte.GetBand(earfcndl)
		assert.NoError(t, err)
		assert.Equal(t, bandExpected, band.ID)
	}
}

func TestGetBandError(t *testing.T) {
	expectedErr := [...]uint32{60140, 60255}

	for _, earfcndl := range expectedErr {
		_, err := lte.GetBand(earfcndl)
		assert.Error(t, err, "Invalid EARFCNDL: no matching band")
	}
}
