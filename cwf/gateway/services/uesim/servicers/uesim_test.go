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

package servicers_test

import (
	"context"
	"testing"

	"magma/cwf/cloud/go/protos"
	"magma/cwf/gateway/services/uesim/servicers"
	"magma/orc8r/cloud/go/test_utils"

	"github.com/stretchr/testify/assert"
)

func TestUESimulator_AddUE(t *testing.T) {
	store := test_utils.NewSQLBlobstore(t, "uesim_uesim_test_blobstore")

	server, err := servicers.NewUESimServer(store)
	assert.NoError(t, err)

	expectedIMSI1 := "1234567890"
	expectedIMSI2 := "2345678901"
	ue1 := &protos.UEConfig{Imsi: expectedIMSI1, AuthKey: make([]byte, 16), AuthOpc: make([]byte, 16), Seq: 0}
	ue2 := &protos.UEConfig{Imsi: expectedIMSI2, AuthKey: make([]byte, 16), AuthOpc: make([]byte, 16), Seq: 0}

	_, err = server.AddUE(context.Background(), ue1)
	assert.NoError(t, err)

	_, err = server.AddUE(context.Background(), ue2)
	assert.NoError(t, err)
}
