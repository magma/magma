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

package servicers_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	models "magma/orc8r/cloud/go/services/ctraced/obsidian/models"
	"magma/orc8r/cloud/go/services/ctraced/servicers"
	"magma/orc8r/cloud/go/services/ctraced/storage"
	deviceTestInit "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestCallTraceServicer(t *testing.T) {
	test_init.StartTestService(t)
	deviceTestInit.StartTestService(t)

	testNetworkId := "n1"
	testGwHwId := "hw1"
	testGwLogicalId := "g1"

	// Initialize network
	err := configurator.CreateNetwork(configurator.Network{ID: testNetworkId}, serdes.Network)
	assert.NoError(t, err)

	// Create a call trace
	testTraceCfg := &models.CallTraceConfig{
		TraceID:   "CallTrace1",
		GatewayID: "test_gateway_id",
		Timeout:   300,
		TraceType: models.CallTraceConfigTraceTypeGATEWAY,
	}
	testTrace := &models.CallTrace{
		Config: testTraceCfg,
		State: &models.CallTraceState{
			CallTraceAvailable: false,
			CallTraceEnding:    false,
		},
	}
	_, err = configurator.CreateEntity(
		testNetworkId,
		configurator.NetworkEntity{
			Type:   orc8r.CallTraceEntityType,
			Key:    "CallTrace1",
			Config: testTrace,
		},
		serdes.Entity,
	)
	assert.NoError(t, err)

	// Create an identity and context for sending requests as gateway
	id := protos.Identity{}
	idgw := protos.Identity_Gateway{HardwareId: testGwHwId, NetworkId: testNetworkId, LogicalId: testGwLogicalId}
	id.SetGateway(&idgw)
	ctx := id.NewContextWithIdentity(context.Background())

	fact := test_utils.NewSQLBlobstore(t, "ctraced_trace_servicer_test_blobstore")
	blobstore := storage.NewCtracedBlobstore(fact)
	srv := servicers.NewCallTraceServicer(blobstore)

	// Missing subscriber ID
	req := &protos.ReportEndedTraceRequest{TraceId: "CallTrace0", Success: true, TraceContent: []byte("abcdefghijklmnopqrstuvwxyz\n")}
	_, err = srv.ReportEndedCallTrace(ctx, req)
	assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = Call trace not found")

	// Successfully ending the call trace
	req = &protos.ReportEndedTraceRequest{TraceId: "CallTrace1", Success: true, TraceContent: []byte("abcdefghijklmnopqrstuvwxyz\n")}
	_, err = srv.ReportEndedCallTrace(ctx, req)
	assert.NoError(t, err)

	// Verify that the call trace has ended
	ent, err := configurator.LoadEntity(
		testNetworkId, orc8r.CallTraceEntityType, "CallTrace1",
		configurator.FullEntityLoadCriteria(),
		serdes.Entity,
	)
	testCallTrace := (&models.CallTrace{}).FromEntity(ent)
	assert.NoError(t, err)
	assert.Equal(t, true, testCallTrace.State.CallTraceAvailable)
	assert.Equal(t, true, testCallTrace.State.CallTraceEnding)
}
