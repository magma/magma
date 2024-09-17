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

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMPPEKeyGeneration(t *testing.T) {
	// Arrange
	ExpectedSendKey := []byte{0x9b, 0x87, 0x83, 0x49, 0x6a, 0x78, 0xcc, 0xaa, 0x34, 0x4e, 0x45, 0x51, 0x7f, 0x15,
		0x37, 0xf9, 0x30, 0x94, 0x26, 0x07, 0x60, 0x68, 0x97, 0xf0, 0xb5, 0x69, 0xab, 0x1d, 0x61, 0x9d, 0x8b, 0xa9,
		0x85, 0x3c, 0xc8, 0xaf, 0x68, 0x4b, 0xaa, 0x8f, 0x8f, 0x77, 0x5f, 0x68, 0x94, 0xf0, 0xcd, 0xc6, 0xc9, 0x2f}
	SendSalt := []byte{0x9b, 0x87}
	RecvSalt := []byte{0x95, 0x63}
	R := []byte{
		0x9f, 0xe8, 0xff, 0xcb, 0xc9, 0xd4, 0x85, 0x97, 0xb9, 0x5b, 0x79, 0x7c, 0x2d, 0xf5, 0x43, 0x31}
	S := []byte("1qaz2wsx")
	MSK := []byte{193, 234, 134, 135, 10, 96, 201, 35, 253, 234, 74, 134, 57, 109, 36, 48, 83, 53, 156, 151, 81, 200, 113, 239, 83, 114, 220, 94, 215, 140, 166, 66, 151, 22, 173, 197, 45, 71, 82, 41, 199, 171, 3, 151, 215, 154, 88, 129, 151, 70, 199, 172, 10, 81, 110, 110, 238, 69, 23, 111, 218, 231, 247, 34}

	// Act
	sendKey := GenerateMPPEKey(MSK[32:], S, R, SendSalt)
	sendKey = append(SendSalt, sendKey...)
	recvKey := GenerateMPPEKey(MSK[:32], S, R, RecvSalt)
	recvKey = append(RecvSalt, recvKey...)

	// Assert
	require.NotNil(t, recvKey)
	require.NotNil(t, sendKey)
	require.Equal(t, ExpectedSendKey, sendKey)
}

func TestSplit(t *testing.T) {
	// Arrange
	arr := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}

	// Act
	chunks := split(arr, 2)

	// Assert
	require.Equal(t, 3, len(chunks))
	require.Equal(t, []byte{0x01, 0x02}, chunks[0])
	require.Equal(t, []byte{0x03, 0x04}, chunks[1])
	require.Equal(t, []byte{0x05, 0x06}, chunks[2])
}

func TestSplitFailure(t *testing.T) {
	// Arrange
	arr := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}

	// Act
	chunks := split(arr, 5)

	// Assert
	require.Nil(t, chunks)
}
