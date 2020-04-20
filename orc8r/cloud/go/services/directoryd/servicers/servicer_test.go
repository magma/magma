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
	"magma/orc8r/cloud/go/services/directoryd/servicers"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	stateTestInit "magma/orc8r/cloud/go/services/state/test_init"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"

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
