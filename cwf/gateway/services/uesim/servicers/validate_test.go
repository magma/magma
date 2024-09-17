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

package servicers

import (
	"errors"
	"testing"

	"magma/cwf/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestValidateUEData(t *testing.T) {
	err := validateUEData(nil)
	assert.Exactly(t, errors.New("Invalid Argument: UE data cannot be nil"), err)

	ue := &protos.UEConfig{Imsi: "0123456789", AuthKey: make([]byte, 16), AuthOpc: make([]byte, 16), Seq: 0}
	err = validateUEData(ue)
	assert.NoError(t, err)
}

func TestValidateUEIMSI(t *testing.T) {
	err := validateUEIMSI("")
	assert.Exactly(t, errors.New("Invalid Argument: IMSI must be between 5 and 15 digits long"), err)

	err = validateUEIMSI("0123")
	assert.Exactly(t, errors.New("Invalid Argument: IMSI must be between 5 and 15 digits long"), err)

	err = validateUEIMSI("0123456789012345")
	assert.Exactly(t, errors.New("Invalid Argument: IMSI must be between 5 and 15 digits long"), err)

	err = validateUEIMSI("0ABCDEF")
	assert.Exactly(t, errors.New("Invalid Argument: IMSI must only be digits"), err)

	err = validateUEIMSI("ABCDEF0")
	assert.Exactly(t, errors.New("Invalid Argument: IMSI must only be digits"), err)

	err = validateUEIMSI("0123456789")
	assert.NoError(t, err)
}

func TestValidateUEKey(t *testing.T) {
	err := validateUEKey(nil)
	assert.Exactly(t, errors.New("Invalid Argument: key cannot be nil"), err)

	err = validateUEKey(make([]byte, 5))
	assert.Exactly(t, errors.New("Invalid Argument: key must be 16 bytes"), err)

	err = validateUEKey(make([]byte, 16))
	assert.NoError(t, err)
}

func TestValidateUEOpc(t *testing.T) {
	err := validateUEOpc(nil)
	assert.Exactly(t, errors.New("Invalid Argument: opc cannot be nil"), err)

	err = validateUEOpc(make([]byte, 5))
	assert.Exactly(t, errors.New("Invalid Argument: opc must be 16 bytes"), err)

	err = validateUEOpc(make([]byte, 16))
	assert.NoError(t, err)
}
