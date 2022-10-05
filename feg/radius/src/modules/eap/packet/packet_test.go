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

package packet

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacketConstruction(t *testing.T) {
	// Act
	packet, err := NewPacketFromRaw([]byte{
		byte(CodeRESPONSE),
		0x17,
		0x00,
		0x07,
		byte(EAPTypeEAPMSCHAPV2),
		0x01,
		0x02,
	})

	// Assert
	assert.True(t, err == nil)
	assert.True(t, packet != nil)
	assert.Equal(t, packet.Length(), 7)
	assert.Equal(t, packet.Code, CodeRESPONSE)
	assert.Equal(t, packet.EAPType, EAPTypeEAPMSCHAPV2)
}

func TestInvalidEAPTypeFails(t *testing.T) {
	// Act
	packet, err := NewPacketFromRaw([]byte{0x00, 0x17, 0x00, 0x07, 0x00, 0x01})

	// Assert
	assert.True(t, packet == nil)
	assert.True(t, err != nil)
	assert.Equal(t, err.Error(), "invalid eap packet code '0'")
}

func TestInvalidLengthFails(t *testing.T) {
	// Act
	packet, err := NewPacketFromRaw([]byte{
		byte(CodeREQUEST),
		0x17,
		0x00,
		0x07,
		byte(EAPTypeIDENTITY),
		0x01,
	})

	// Assert
	assert.True(t, packet == nil)
	assert.True(t, err != nil)
	assert.Equal(t, err.Error(), "length mismatch (packet header indicates 7, but packet contains 6 data bytes)")
}

func TestPacketTooShortFails(t *testing.T) {
	// Act
	packet, err := NewPacketFromRaw([]byte{
		byte(CodeREQUEST),
		0x17,
	})

	// Assert
	assert.True(t, packet == nil)
	assert.True(t, err != nil)
	assert.Equal(t, err.Error(), "packet length must be at least 4 bytes, got 2 bytes")
}

func TestToBytes(t *testing.T) {
	// Arrange
	originalBytes := []byte{byte(CodeREQUEST), 0x17, 0x00, 0x07, byte(EAPTypeAKA), 0x01, 0x02}
	packet, err := NewPacketFromRaw(originalBytes)
	assert.True(t, err == nil)

	// Act
	bytes, err := packet.Bytes()
	assert.True(t, err == nil)

	// Assert
	assert.True(t, reflect.DeepEqual(originalBytes, bytes))
}
