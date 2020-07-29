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

package diameter_test

import (
	"testing"

	"magma/feg/gateway/diameter"

	"github.com/stretchr/testify/assert"
)

func TestEncodePLMN(t *testing.T) {
	var plmnBytes []byte
	var err error
	// 5 digit
	plmnBytes, err = diameter.EncodePLMNID("72207")
	assert.Equal(t, plmnBytes, []byte{0x27, 0xf2, 0x70})
	assert.Nil(t, err)
	// 6 digit
	plmnBytes, err = diameter.EncodePLMNID("722070")
	assert.Equal(t, plmnBytes, []byte{0x27, 0x02, 0x70})
	assert.Nil(t, err)
	// 6 digit leading 0s
	plmnBytes, err = diameter.EncodePLMNID("001001")
	assert.Equal(t, plmnBytes, []byte{0x00, 0x11, 0x00})
	assert.Nil(t, err)
	// 6 digit leading 0s
	plmnBytes, err = diameter.EncodePLMNID("00101")
	assert.Equal(t, plmnBytes, []byte{0x00, 0xf1, 0x10})
	assert.Nil(t, err)

	// errors
	_, err = diameter.EncodePLMNID("0010100")
	assert.NotNil(t, err)
	_, err = diameter.EncodePLMNID("0010")
	assert.NotNil(t, err)
}

func TestEncodeUserLocation(t *testing.T) {
	// values and expected bytes obtained from example Gx trace
	plmn := "72207"
	var tai uint16 = 2461
	var ecgi uint32 = 134477315
	expectedBytes := []byte{0x82, 0x27, 0xf2, 0x70, 0x09, 0x9d, 0x27, 0xf2, 0x70, 0x08, 0x03, 0xf6, 0x03}

	userLocBytes, err := diameter.EncodeUserLocation(plmn, tai, ecgi)
	assert.Equal(t, userLocBytes, expectedBytes)
	assert.Nil(t, err)
}
