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

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/orc8r"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test_utils "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/directoryd/servicers"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

const (
	hn0   = "some_hostname_0"
	hn1   = "some_hostname_1"
	hn2   = "some_hostname_2"
	hwid0 = "some_hardwareid_0"
	hwid1 = "some_hardwareid_1"
	hwid2 = "some_hardwareid_2"
	imsi0 = "some_imsi_0"
	imsi1 = "some_imsi_1"
	imsi2 = "some_imsi_2"
	imsi3 = "some_imsi_3"
	nid0  = "some_networkid_0"
	nid1  = "some_networkid_1"
	sid0  = "some_sessionid_0"
	sid1  = "some_sessionid_1"
	sid2  = "some_sessionid_2"
)

func newTestDirectoryLookupServicer(t *testing.T) protos.DirectoryLookupServer {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)

	fact := blobstore.NewEntStorage(storage.DirectorydTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	assert.NoError(t, err)

	store := storage.NewDirectorydBlobstore(fact)
	srv, err := servicers.NewDirectoryLookupServicer(store)
	assert.NoError(t, err)

	return srv
}

func TestDirectoryLookupServicer_HostnameToHWID(t *testing.T) {
	srv := newTestDirectoryLookupServicer(t)
	stateTestInit.StartTestService(t)
	ctx := context.Background()

	// Empty initially
	get := &protos.GetHostnameForHWIDRequest{Hwid: hwid0}
	_, err := srv.GetHostnameForHWID(ctx, get)
	assert.Error(t, err)

	// Put and get hwid0->hostname0
	put := &protos.MapHWIDToHostnameRequest{HwidToHostname: map[string]string{hwid0: hn0}}
	_, err = srv.MapHWIDsToHostnames(ctx, put)
	assert.NoError(t, err)
	get = &protos.GetHostnameForHWIDRequest{Hwid: hwid0}
	res, err := srv.GetHostnameForHWID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, hn0, res.Hostname)

	// Put and get hwid1->hostname1, hwid2->hostname2
	put = &protos.MapHWIDToHostnameRequest{HwidToHostname: map[string]string{hwid1: hn1, hwid2: hn2}}
	_, err = srv.MapHWIDsToHostnames(ctx, put)
	assert.NoError(t, err)
	get = &protos.GetHostnameForHWIDRequest{Hwid: hwid1}
	res, err = srv.GetHostnameForHWID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, hn1, res.Hostname)
	get = &protos.GetHostnameForHWIDRequest{Hwid: hwid2}
	res, err = srv.GetHostnameForHWID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, hn2, res.Hostname)

	// hwid0->hostname0 still intact
	get = &protos.GetHostnameForHWIDRequest{Hwid: hwid0}
	res, err = srv.GetHostnameForHWID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, hn0, res.Hostname)
}

func TestDirectoryLookupServicer_SessionIDToIMSI(t *testing.T) {
	srv := newTestDirectoryLookupServicer(t)
	stateTestInit.StartTestService(t)
	ctx := context.Background()

	// Empty initially
	get := &protos.GetIMSIForSessionIDRequest{NetworkID: nid0, SessionID: sid0}
	_, err := srv.GetIMSIForSessionID(ctx, get)
	assert.Error(t, err)

	// Put and get sid0->imsi0
	put := &protos.MapSessionIDToIMSIRequest{NetworkID: nid0, SessionIDToIMSI: map[string]string{sid0: imsi0}}
	_, err = srv.MapSessionIDsToIMSIs(ctx, put)
	assert.NoError(t, err)
	get = &protos.GetIMSIForSessionIDRequest{NetworkID: nid0, SessionID: sid0}
	res, err := srv.GetIMSIForSessionID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, res.Imsi)

	// Put and get sid1->imsi1, sid2->imsi2
	put = &protos.MapSessionIDToIMSIRequest{NetworkID: nid0, SessionIDToIMSI: map[string]string{sid1: imsi1, sid2: imsi2}}
	_, err = srv.MapSessionIDsToIMSIs(ctx, put)
	assert.NoError(t, err)
	get = &protos.GetIMSIForSessionIDRequest{NetworkID: nid0, SessionID: sid1}
	res, err = srv.GetIMSIForSessionID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, res.Imsi)
	get = &protos.GetIMSIForSessionIDRequest{NetworkID: nid0, SessionID: sid2}
	res, err = srv.GetIMSIForSessionID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, imsi2, res.Imsi)

	// sid0->imsi0 still intact
	get = &protos.GetIMSIForSessionIDRequest{NetworkID: nid0, SessionID: sid0}
	res, err = srv.GetIMSIForSessionID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, res.Imsi)

	// Correctly network-partitioned: {nid0: sid0->imsi0, nid1: sid0->imsi1}
	put = &protos.MapSessionIDToIMSIRequest{NetworkID: nid0, SessionIDToIMSI: map[string]string{sid0: imsi0}}
	_, err = srv.MapSessionIDsToIMSIs(ctx, put)
	assert.NoError(t, err)
	put = &protos.MapSessionIDToIMSIRequest{NetworkID: nid1, SessionIDToIMSI: map[string]string{sid0: imsi1}}
	_, err = srv.MapSessionIDsToIMSIs(ctx, put)
	assert.NoError(t, err)
	get = &protos.GetIMSIForSessionIDRequest{NetworkID: nid0, SessionID: sid0}
	res, err = srv.GetIMSIForSessionID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, res.Imsi)
	get = &protos.GetIMSIForSessionIDRequest{NetworkID: nid1, SessionID: sid0}
	res, err = srv.GetIMSIForSessionID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, res.Imsi)

	// Fail with empty network ID
	get = &protos.GetIMSIForSessionIDRequest{SessionID: sid0}
	_, err = srv.GetIMSIForSessionID(ctx, get)
	assert.Error(t, err)
	put = &protos.MapSessionIDToIMSIRequest{SessionIDToIMSI: map[string]string{sid0: imsi0}}
	_, err = srv.MapSessionIDsToIMSIs(ctx, put)
	assert.Error(t, err)
}

func TestDirectoryLookupServicer_HWIDToRecordIDs(t *testing.T) {
	srv := newTestDirectoryLookupServicer(t)
	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	stateTestInit.StartTestService(t)
	ctx := context.Background()

	recordIDs0 := &protos.DirectoryRecordIDs{Ids: []string{imsi0, imsi1}}
	recordIDs1 := &protos.DirectoryRecordIDs{Ids: []string{imsi2}}

	hwIDsToRecordIDs0 := map[string]*protos.DirectoryRecordIDs{
		hwid0: recordIDs0,
	}
	multiHWIDToRecordIds := map[string]*protos.DirectoryRecordIDs{
		hwid0: recordIDs0,
		hwid1: recordIDs1,
	}

	// Empty initially
	get := &protos.GetDirectoryRecordIDsForHWIDRequest{NetworkID: nid0, Hwid: hwid0}
	_, err := srv.GetDirectoryRecordIDsForHWID(ctx, get)
	assert.Error(t, err)

	configurator_test_utils.RegisterNetwork(t, nid0, "DirectoryD Service Test")
	configurator_test_utils.RegisterGateway(t, nid0, hwid0, &models.GatewayDevice{HardwareID: hwid0})
	configurator_test_utils.RegisterGateway(t, nid0, hwid1, &models.GatewayDevice{HardwareID: hwid1})

	// Report hwid0 records
	reportDirectorydState(t, hwid0, imsi0)
	reportDirectorydState(t, hwid0, imsi1)

	// Put and get hwid0->recordIDs0
	put := &protos.MapHWIDToDirectoryRecordIDsRequest{NetworkID: nid0, HwidToRecordIDs: hwIDsToRecordIDs0}
	_, err = srv.MapHWIDToDirectoryRecordIDs(ctx, put)
	assert.NoError(t, err)
	res, err := srv.GetDirectoryRecordIDsForHWID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, recordIDs0, res.Ids)

	// Report new records
	reportDirectorydState(t, hwid0, imsi3)
	reportDirectorydState(t, hwid1, imsi2)

	// Update recordIDs0, then put and get hwid0->recordIDs0, hwid1->recordIDs1
	hwIDsToRecordIDs0[hwid0].Ids = append(hwIDsToRecordIDs0[hwid0].Ids, imsi3)
	put = &protos.MapHWIDToDirectoryRecordIDsRequest{NetworkID: nid0, HwidToRecordIDs: multiHWIDToRecordIds}
	_, err = srv.MapHWIDToDirectoryRecordIDs(ctx, put)
	assert.NoError(t, err)

	res, err = srv.GetDirectoryRecordIDsForHWID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, recordIDs0, res.Ids)

	get = &protos.GetDirectoryRecordIDsForHWIDRequest{NetworkID: nid0, Hwid: hwid1}
	res, err = srv.GetDirectoryRecordIDsForHWID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, recordIDs1, res.Ids)

	// Ensure deleted state doesn't get returned
	deleteDirectorydState(t, hwid0, imsi3)
	get = &protos.GetDirectoryRecordIDsForHWIDRequest{NetworkID: nid0, Hwid: hwid0}
	expectedIDs := &protos.DirectoryRecordIDs{Ids: []string{imsi0, imsi1}}
	res, err = srv.GetDirectoryRecordIDsForHWID(ctx, get)
	assert.NoError(t, err)
	assert.Equal(t, expectedIDs, res.Ids)

	// Fail get with empty network ID
	get = &protos.GetDirectoryRecordIDsForHWIDRequest{Hwid: hwid0}
	_, err = srv.GetDirectoryRecordIDsForHWID(ctx, get)
	assert.Error(t, err)

	// Fail get with empty hardware ID
	get = &protos.GetDirectoryRecordIDsForHWIDRequest{NetworkID: nid0}
	_, err = srv.GetDirectoryRecordIDsForHWID(ctx, get)
	assert.Error(t, err)

	// Fail put with empty network ID
	put = &protos.MapHWIDToDirectoryRecordIDsRequest{HwidToRecordIDs: multiHWIDToRecordIds}
	_, err = srv.MapHWIDToDirectoryRecordIDs(ctx, put)
	assert.Error(t, err)

	// Fail put with empty map
	put = &protos.MapHWIDToDirectoryRecordIDsRequest{NetworkID: nid0}
	_, err = srv.MapHWIDToDirectoryRecordIDs(ctx, put)
	assert.Error(t, err)
}

func reportDirectorydState(t *testing.T, hwid string, imsi string) {
	stateClient, err := getStateServiceClient(t)
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, hwid)
	record := &types.DirectoryRecord{
		LocationHistory: []string{hwid},
		Identifiers: map[string]interface{}{
			types.RecordKeySessionID: sid0,
		},
	}
	serializedRecord, err := record.MarshalBinary()
	assert.NoError(t, err)

	st := &protos.State{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: imsi,
		Value:    serializedRecord,
	}
	reqReport := &protos.ReportStatesRequest{States: []*protos.State{st}}
	res, err := stateClient.ReportStates(ctx, reqReport)
	assert.NoError(t, err)
	assert.Empty(t, res.UnreportedStates)
}

func deleteDirectorydState(t *testing.T, hwid string, imsi string) {
	stateClient, err := getStateServiceClient(t)
	assert.NoError(t, err)
	ctx := test_utils.GetContextWithCertificate(t, hwid0)
	stId := &protos.StateID{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: imsi,
	}
	reqReport := &protos.DeleteStatesRequest{NetworkID: nid0, Ids: []*protos.StateID{stId}}
	_, err = stateClient.DeleteStates(ctx, reqReport)
	assert.NoError(t, err)
}

func getStateServiceClient(t *testing.T) (protos.StateServiceClient, error) {
	conn, err := registry.GetConnection(state.ServiceName)
	assert.NoError(t, err)
	return protos.NewStateServiceClient(conn), err
}
