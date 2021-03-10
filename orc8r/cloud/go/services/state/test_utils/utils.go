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

package test_utils

import (
	"testing"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/service/middleware/unary/test_utils"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func ReportGatewayStatus(t *testing.T, ctx context.Context, req *models.GatewayStatus) {
	client, err := state.GetStateClient()
	assert.NoError(t, err)

	serializedGWStatus, err := serde.Serialize(req, orc8r.GatewayStateType, serdes.State)
	assert.NoError(t, err)
	states := []*protos.State{
		{
			Type:     orc8r.GatewayStateType,
			DeviceID: req.HardwareID,
			Value:    serializedGWStatus,
		},
	}
	_, err = client.ReportStates(
		ctx,
		&protos.ReportStatesRequest{States: states},
	)
	assert.NoError(t, err)
}

func ReportState(t *testing.T, ctx context.Context, stateType string, stateKey string, stateVal interface{}, serdes serde.Registry) {
	client, err := state.GetStateClient()
	assert.NoError(t, err)
	serializedState, err := serde.Serialize(stateVal, stateType, serdes)
	assert.NoError(t, err)
	states := []*protos.State{
		{
			Type:     stateType,
			DeviceID: stateKey,
			Value:    serializedState,
		},
	}
	res, err := client.ReportStates(ctx, &protos.ReportStatesRequest{States: states})
	assert.NoError(t, err)
	assert.Empty(t, res.UnreportedStates)
}

func GetContextWithCertificate(t *testing.T, hwID string) context.Context {
	csn := test_utils.StartMockGwAccessControl(t, []string{hwID})
	return metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs(identity.CLIENT_CERT_SN_KEY, csn[0]))
}
