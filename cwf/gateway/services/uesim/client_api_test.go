/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package uesim_test

import (
	"testing"

	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/services/uesim"
	"magma/cwf/gateway/services/uesim/test_init"
	"magma/lte/cloud/go/crypto"

	"github.com/stretchr/testify/assert"
)

// todo use a config
const (
	Op = "\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"
)

func TestUESimClient(t *testing.T) {
	test_init.StartTestService(t)
	imsi := "001010000000001"
	key := make([]byte, 16)
	opc, err := crypto.GenerateOpc(key, []byte(Op))
	assert.NoError(t, err)
	seq := uint64(0)

	ue := &cwfprotos.UEConfig{Imsi: imsi, AuthKey: key, AuthOpc: opc[:], Seq: seq}
	err = uesim.AddUE(ue)
	assert.NoError(t, err)
}
