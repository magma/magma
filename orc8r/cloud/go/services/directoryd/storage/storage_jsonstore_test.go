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

package storage_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/JsonStore"
	"magma/orc8r/cloud/go/JsonStore/mocks"
	
	dstorage "magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/merrors"
)

const (
	placeholderNetworkID = "placeholder_network"
)

func TestDirectorydJsonstoreStorage_GetHostnameForHWID(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	hwid := "some_hwid"
	tk := storage.TK{Type: dstorage.DirectorydTypeHWIDToHostname, Key: hwid}

	hostname := "some_hostname"
	json := JsonStore.Json{
		Type:  tk.Type,
		Key:   tk.Key,
		Value: hostname,
	}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err := store.GetHostnameForHWID(hwid)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", placeholderNetworkID, tk).Return(JsonStore.Json{}, merrors.ErrNotFound).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err = store.GetHostnameForHWID(hwid)
	assert.Exactly(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", placeholderNetworkID, tk).Return(JsonStore.Json{}, someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err = store.GetHostnameForHWID(hwid)
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", placeholderNetworkID, tk).Return(json, nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	hostnameRecvd, err := store.GetHostnameForHWID(hwid)
	assert.NoError(t, err)
	assert.Equal(t, hostname, hostnameRecvd)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstoreStorage_MapHWIDToHostname(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	hwids := []string{"some_hwid_0", "some_hwid_1"}
	hostnames := []string{"some_hostname_0", "some_hostname_1"}
	hwidToHostname := map[string]string{
		hwids[0]: hostnames[0],
		hwids[1]: hostnames[1],
	}

	tks := storage.TKs{
		{Type: dstorage.DirectorydTypeHWIDToHostname, Key: hwids[0]},
		{Type: dstorage.DirectorydTypeHWIDToHostname, Key: hwids[1]},
	}

	jsons := JsonStore.Jsons{
		{
			Type:  tks[0].Type,
			Key:   tks[0].Key,
			Value: hostnames[0],
		},
		{
			Type:  tks[1].Type,
			Key:   tks[1].Key,
			Value: hostnames[1],
		},
	}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	err := store.MapHWIDsToHostnames(hwidToHostname)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.PutRecord fails
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Write", placeholderNetworkID, mock.Anything, mock.Anything).
		Return(someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.MapHWIDsToHostnames(hwidToHostname)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Write", placeholderNetworkID, jsons).
		Return(nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.MapHWIDsToHostnames(hwidToHostname)
	assert.NoError(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstoreStorage_DeMapHWIDToHostname(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	hwid0 := "some_hwid_0"

	tks := storage.TKs{
		{Type: dstorage.DirectorydTypeHWIDToHostname, Key: hwid0},
	}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	err := store.UnmapHWIDsToHostnames([]string{hwid0})
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with error
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Delete", placeholderNetworkID, tks).Return(someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.UnmapHWIDsToHostnames([]string{hwid0})
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Delete", placeholderNetworkID, tks).Return(nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.UnmapHWIDsToHostnames([]string{hwid0})
	assert.NoError(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstoreStorage_GetIMSIForSessionID(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	nid := "some_networkid"

	sid := "some_sessionid"
	tk := storage.TK{Type: dstorage.DirectorydTypeSessionIDToIMSI, Key: sid}

	imsi := "some_imsi"
	json := JsonStore.Json{
		Type:  tk.Type,
		Key:   tk.Key,
		Value: imsi,
	}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err := store.GetIMSIForSessionID(nid, sid)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", nid, tk).Return(JsonStore.Json{}, merrors.ErrNotFound).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err = store.GetIMSIForSessionID(nid, sid)
	assert.Exactly(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with error
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", nid, tk).Return(JsonStore.Json{}, someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err = store.GetIMSIForSessionID(nid, sid)
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", nid, tk).Return(json, nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	imsiRecvd, err := store.GetIMSIForSessionID(nid, sid)
	assert.NoError(t, err)
	assert.Equal(t, imsi, imsiRecvd)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstoreStorage_MapSessionIDToIMSI(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	nid := "some_networkid"

	sids := []string{"some_sessionid_0", "some_sessionid_1"}
	imsis := []string{"some_imsi_0", "some_imsi_1"}
	sidToIMSI := map[string]string{
		sids[0]: imsis[0],
		sids[1]: imsis[1],
	}

	tks := storage.TKs{
		{Type: dstorage.DirectorydTypeSessionIDToIMSI, Key: sids[0]},
		{Type: dstorage.DirectorydTypeSessionIDToIMSI, Key: sids[1]},
	}

	jsons := JsonStore.Jsons{
		{
			Type:  tks[0].Type,
			Key:   tks[0].Key,
			Value: imsis[0],
		},
		{
			Type:  tks[1].Type,
			Key:   tks[1].Key,
			Value: imsis[1],
		},
	}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	err := store.MapSessionIDsToIMSIs(nid, sidToIMSI)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.PutRecord fails
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Write", nid, mock.Anything, mock.Anything).
		Return(someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.MapSessionIDsToIMSIs(nid, sidToIMSI)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Write", nid, jsons).
		Return(nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.MapSessionIDsToIMSIs(nid, sidToIMSI)
	assert.NoError(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstoreStorage_DeMapSessionIDToIMSI(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	nid := "some_networkid"

	sid := "some_sessionid"
	tks := storage.TKs{{Type: dstorage.DirectorydTypeSessionIDToIMSI, Key: sid}}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	err := store.UnmapSessionIDsToIMSIs(nid, []string{sid})
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with error
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Delete", nid, tks).Return(someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.UnmapSessionIDsToIMSIs(nid, []string{sid})
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Delete", nid, tks).Return(nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.UnmapSessionIDsToIMSIs(nid, []string{sid})
	assert.NoError(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstoreStorage_GetHWIDForSgwCTeid(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	nid := "some_networkid"

	teid := "1"
	tk := storage.TK{Type: dstorage.DirectorydTypeSgwCteidToHwid, Key: teid}

	hwId := "hwId_1"
	json := JsonStore.Json{
		Type:  tk.Type,
		Key:   tk.Key,
		Value: hwId,
	}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err := store.GetHWIDForSgwCTeid(nid, teid)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", nid, tk).Return(JsonStore.Json{}, merrors.ErrNotFound).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err = store.GetHWIDForSgwCTeid(nid, teid)
	assert.Exactly(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", nid, tk).Return(JsonStore.Json{}, someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err = store.GetHWIDForSgwCTeid(nid, teid)
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", nid, tk).Return(json, nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	hwIdRecvd, err := store.GetHWIDForSgwCTeid(nid, teid)
	assert.NoError(t, err)
	assert.Equal(t, hwId, hwIdRecvd)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstoreStorage_MapSgwCTeidToHWID(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	nid := "some_networkid"

	teids := []string{"1", "2"}
	hwIds := []string{"hwid_1", "hwid_2"}
	teidsToHwIds := map[string]string{
		teids[0]: hwIds[0],
		teids[1]: hwIds[1],
	}

	tks := storage.TKs{
		{Type: dstorage.DirectorydTypeSgwCteidToHwid, Key: teids[0]},
		{Type: dstorage.DirectorydTypeSgwCteidToHwid, Key: teids[1]},
	}

	jsons := JsonStore.Jsons{
		{
			Type:  tks[0].Type,
			Key:   tks[0].Key,
			Value: hwIds[0],
		},
		{
			Type:  tks[1].Type,
			Key:   tks[1].Key,
			Value: hwIds[1],
		},
	}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	err := store.MapSgwCTeidToHWID(nid, teidsToHwIds)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.PutRecord fails
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Write", nid, mock.Anything, mock.Anything).
		Return(someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.MapSgwCTeidToHWID(nid, teidsToHwIds)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Write", nid, jsons).
		Return(nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.MapSgwCTeidToHWID(nid, teidsToHwIds)
	assert.NoError(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstoreStorage_UnmapSgwCTeidToHWID(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	nid := "some_networkid"

	teid := "some_sessionid"
	tks := storage.TKs{{Type: dstorage.DirectorydTypeSgwCteidToHwid, Key: teid}}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	err := store.UnmapSgwCTeidToHWID(nid, []string{teid})
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with error
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Delete", nid, tks).Return(someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.UnmapSgwCTeidToHWID(nid, []string{teid})
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Delete", nid, tks).Return(nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.UnmapSgwCTeidToHWID(nid, []string{teid})
	assert.NoError(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstore_GetHWIDForSgwUTeid(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	nid := "some_networkid"

	teid := "1"
	tk := storage.TK{Type: dstorage.DirectorydTypeSgwUteidToHwid, Key: teid}

	hwId := "hwId_1"
	json := JsonStore.Json{
		Type:  tk.Type,
		Key:   tk.Key,
		Value: hwId,
	}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err := store.GetHWIDForSgwUTeid(nid, teid)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", nid, tk).Return(JsonStore.Json{}, merrors.ErrNotFound).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err = store.GetHWIDForSgwUTeid(nid, teid)
	assert.Exactly(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", nid, tk).Return(JsonStore.Json{}, someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	_, err = store.GetHWIDForSgwUTeid(nid, teid)
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Get", nid, tk).Return(json, nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	hwIdRecvd, err := store.GetHWIDForSgwUTeid(nid, teid)
	assert.NoError(t, err)
	assert.Equal(t, hwId, hwIdRecvd)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstore_MapSgwUTeidToHWID(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	nid := "some_networkid"

	teids := []string{"1", "2"}
	hwIds := []string{"hwid_1", "hwid_2"}
	teidsToHwIds := map[string]string{
		teids[0]: hwIds[0],
		teids[1]: hwIds[1],
	}

	tks := storage.TKs{
		{Type: dstorage.DirectorydTypeSgwUteidToHwid, Key: teids[0]},
		{Type: dstorage.DirectorydTypeSgwUteidToHwid, Key: teids[1]},
	}

	jsons := JsonStore.Jsons{
		{
			Type:  tks[0].Type,
			Key:   tks[0].Key,
			Value: hwIds[0],
		},
		{
			Type:  tks[1].Type,
			Key:   tks[1].Key,
			Value: hwIds[1],
		},
	}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	err := store.MapSgwUTeidToHWID(nid, teidsToHwIds)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.PutRecord fails
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Write", nid, mock.Anything, mock.Anything).
		Return(someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.MapSgwUTeidToHWID(nid, teidsToHwIds)
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Write", nid, jsons).
		Return(nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.MapSgwUTeidToHWID(nid, teidsToHwIds)
	assert.NoError(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}

func TestDirectorydJsonstore_UnmapSgwUTeidToHWID(t *testing.T) {
	var jsonFactMock *mocks.StoreFactory
	var jsonStoreMock *mocks.Store
	someErr := errors.New("generic error")

	nid := "some_networkid"

	teid := "some_sessionid"
	tks := storage.TKs{{Type: dstorage.DirectorydTypeSgwUteidToHwid, Key: teid}}

	// Fail to start transaction
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydJsonstore(jsonFactMock)

	err := store.UnmapSgwUTeidToHWID(nid, []string{teid})
	assert.Error(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// store.Get fails with error
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Delete", nid, tks).Return(someErr).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.UnmapSgwUTeidToHWID(nid, []string{teid})
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)

	// Success
	jsonFactMock = &mocks.StoreFactory{}
	jsonStoreMock = &mocks.Store{}
	jsonFactMock.On("StartTransaction", mock.Anything).Return(jsonStoreMock, nil).Once()
	jsonStoreMock.On("Rollback").Return(nil).Once()
	jsonStoreMock.On("Delete", nid, tks).Return(nil).Once()
	jsonStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydJsonstore(jsonFactMock)

	err = store.UnmapSgwUTeidToHWID(nid, []string{teid})
	assert.NoError(t, err)
	jsonFactMock.AssertExpectations(t)
	jsonStoreMock.AssertExpectations(t)
}
