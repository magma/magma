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

package storage_test

import (
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/sqorc"
	merrors "magma/orc8r/lib/go/errors"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestDirectorydStorageBlobstore_Integation(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	fact := blobstore.NewEntStorage(storage.DirectorydTableBlobstore, db, sqorc.GetSqlBuilder())
	err = fact.InitializeFactory()
	assert.NoError(t, err)
	store := storage.NewDirectorydBlobstore(fact)
	testDirectorydStorageImpl(t, store)
}

func testDirectorydStorageImpl(t *testing.T, store storage.DirectorydStorage) {
	hwid0 := "some_hwid_0"
	hwid1 := "some_hwid_1"
	hwid2 := "some_hwid_2"
	hwid3 := "some_hwid_3"
	hostname0 := "some_hostname_0"
	hostname1 := "some_hostname_1"
	hostname2 := "some_hostname_2"
	hostname3 := "some_hostname_3"

	nid0 := "some_networkid_0"
	nid1 := "some_networkid_1"
	sid0 := "some_sessionid_0"
	sid1 := "some_sessionid_1"
	imsi0 := "some_imsi_0"
	imsi1 := "some_imsi_1"

	teid0 := "10"
	teid1 := "20"

	//////////////////////////////
	// Hostname -> HWID
	//////////////////////////////

	// Empty initially
	_, err := store.GetHostnameForHWID(hwid0)
	assert.Exactly(t, err, merrors.ErrNotFound)
	_, err = store.GetHostnameForHWID(hwid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get hwid0->hostname1
	err = store.MapHWIDsToHostnames(map[string]string{hwid0: hostname1})
	assert.NoError(t, err)
	recvd, err := store.GetHostnameForHWID(hwid0)
	assert.NoError(t, err)
	assert.Equal(t, hostname1, recvd)
	_, err = store.GetHostnameForHWID(hwid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get hwid0->hostname0
	err = store.MapHWIDsToHostnames(map[string]string{hwid0: hostname0})
	assert.NoError(t, err)
	recvd, err = store.GetHostnameForHWID(hwid0)
	assert.NoError(t, err)
	assert.Equal(t, hostname0, recvd)
	_, err = store.GetHostnameForHWID(hwid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get hwid1->hostname1
	err = store.MapHWIDsToHostnames(map[string]string{hwid1: hostname1})
	assert.NoError(t, err)
	recvd, err = store.GetHostnameForHWID(hwid0)
	assert.NoError(t, err)
	assert.Equal(t, hostname0, recvd)
	recvd, err = store.GetHostnameForHWID(hwid1)
	assert.NoError(t, err)
	assert.Equal(t, hostname1, recvd)

	// Multi-put: Put and Get hwid2->hostname2, hwid3->hostname3
	err = store.MapHWIDsToHostnames(map[string]string{hwid2: hostname2, hwid3: hostname3})
	assert.NoError(t, err)
	recvd, err = store.GetHostnameForHWID(hwid2)
	assert.NoError(t, err)
	assert.Equal(t, hostname2, recvd)
	recvd, err = store.GetHostnameForHWID(hwid3)
	assert.NoError(t, err)
	assert.Equal(t, hostname3, recvd)

	//////////////////////////////
	// Session ID -> IMSI
	//////////////////////////////

	// Empty initially
	_, err = store.GetIMSIForSessionID(nid0, sid0)
	assert.Exactly(t, err, merrors.ErrNotFound)
	_, err = store.GetIMSIForSessionID(nid0, sid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get sid0->imsi1
	err = store.MapSessionIDsToIMSIs(nid0, map[string]string{sid0: imsi1})
	assert.NoError(t, err)
	recvd, err = store.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, recvd)
	_, err = store.GetIMSIForSessionID(nid0, sid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get sid0->imsi0
	err = store.MapSessionIDsToIMSIs(nid0, map[string]string{sid0: imsi0})
	assert.NoError(t, err)
	recvd, err = store.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, recvd)
	_, err = store.GetIMSIForSessionID(nid0, sid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get sid1->imsi1
	err = store.MapSessionIDsToIMSIs(nid0, map[string]string{sid1: imsi1})
	assert.NoError(t, err)
	recvd, err = store.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, recvd)
	recvd, err = store.GetIMSIForSessionID(nid0, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, recvd)

	// Multi-put: Put and Get sid0->imsi0, sid1->imsi1 for nid1
	err = store.MapSessionIDsToIMSIs(nid1, map[string]string{sid0: imsi0, sid1: imsi1})
	assert.NoError(t, err)
	recvd, err = store.GetIMSIForSessionID(nid1, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, recvd)
	recvd, err = store.GetIMSIForSessionID(nid1, sid1)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, recvd)

	// Correctly network-partitioned: {nid0: sid0->imsi0, nid1: sid0->imsi1}
	err = store.MapSessionIDsToIMSIs(nid0, map[string]string{sid0: imsi0})
	assert.NoError(t, err)
	err = store.MapSessionIDsToIMSIs(nid1, map[string]string{sid0: imsi1})
	assert.NoError(t, err)
	recvd, err = store.GetIMSIForSessionID(nid0, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi0, recvd)
	recvd, err = store.GetIMSIForSessionID(nid1, sid0)
	assert.NoError(t, err)
	assert.Equal(t, imsi1, recvd)

	//////////////////////////////
	// Teid -> HwId
	//////////////////////////////

	// Empty initially
	_, err = store.GetHWIDForSgwCTeid(nid0, teid0)
	assert.Exactly(t, err, merrors.ErrNotFound)
	_, err = store.GetHWIDForSgwCTeid(nid0, teid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get teid0->HwId1
	err = store.MapSgwCTeidToHWID(nid0, map[string]string{teid0: hwid1})
	assert.NoError(t, err)
	recvd, err = store.GetHWIDForSgwCTeid(nid0, teid0)
	assert.NoError(t, err)
	assert.Equal(t, hwid1, recvd)
	_, err = store.GetHWIDForSgwCTeid(nid0, teid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get teid0->HwId0
	err = store.MapSgwCTeidToHWID(nid0, map[string]string{teid0: hwid0})
	assert.NoError(t, err)
	recvd, err = store.GetHWIDForSgwCTeid(nid0, teid0)
	assert.NoError(t, err)
	assert.Equal(t, hwid0, recvd)
	_, err = store.GetHWIDForSgwCTeid(nid0, teid1)
	assert.Exactly(t, err, merrors.ErrNotFound)

	// Put and Get teid1->HwId1
	err = store.MapSgwCTeidToHWID(nid0, map[string]string{teid1: hwid1})
	assert.NoError(t, err)
	recvd, err = store.GetHWIDForSgwCTeid(nid0, teid0)
	assert.NoError(t, err)
	assert.Equal(t, hwid0, recvd)
	recvd, err = store.GetHWIDForSgwCTeid(nid0, teid1)
	assert.NoError(t, err)
	assert.Equal(t, hwid1, recvd)

	// Multi-put: Put and Get teid0->HwId0, teid1->HwId1 for nid1
	err = store.MapSgwCTeidToHWID(nid1, map[string]string{teid0: hwid0, teid1: hwid1})
	assert.NoError(t, err)
	recvd, err = store.GetHWIDForSgwCTeid(nid1, teid0)
	assert.NoError(t, err)
	assert.Equal(t, hwid0, recvd)
	recvd, err = store.GetHWIDForSgwCTeid(nid1, teid1)
	assert.NoError(t, err)
	assert.Equal(t, hwid1, recvd)

	// Correctly network-partitioned: {nid0: teid0->HwId0, nid1: teid0->HwId1}
	err = store.MapSgwCTeidToHWID(nid0, map[string]string{teid0: hwid0})
	assert.NoError(t, err)
	err = store.MapSgwCTeidToHWID(nid1, map[string]string{teid0: hwid1})
	assert.NoError(t, err)
	recvd, err = store.GetHWIDForSgwCTeid(nid0, teid0)
	assert.NoError(t, err)
	assert.Equal(t, hwid0, recvd)
	recvd, err = store.GetHWIDForSgwCTeid(nid1, teid0)
	assert.NoError(t, err)
	assert.Equal(t, hwid1, recvd)

}
