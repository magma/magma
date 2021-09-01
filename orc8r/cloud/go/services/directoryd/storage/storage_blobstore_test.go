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

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/blobstore/mocks"
	dstorage "magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	placeholderNetworkID = "placeholder_network"
)

func TestDirectorydBlobstoreStorage_GetHostnameForHWID(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	hwid := "some_hwid"
	tk := storage.TypeAndKey{Type: dstorage.DirectorydTypeHWIDToHostname, Key: hwid}

	hostname := "some_hostname"
	blob := blobstore.Blob{
		Type:  tk.Type,
		Key:   tk.Key,
		Value: []byte(hostname),
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err := store.GetHostnameForHWID(hwid)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", placeholderNetworkID, tk).Return(blobstore.Blob{}, merrors.ErrNotFound).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err = store.GetHostnameForHWID(hwid)
	assert.Exactly(t, merrors.ErrNotFound, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", placeholderNetworkID, tk).Return(blobstore.Blob{}, someErr).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err = store.GetHostnameForHWID(hwid)
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", placeholderNetworkID, tk).Return(blob, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	hostnameRecvd, err := store.GetHostnameForHWID(hwid)
	assert.NoError(t, err)
	assert.Equal(t, hostname, hostnameRecvd)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestDirectorydBlobstoreStorage_MapHWIDToHostname(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	hwids := []string{"some_hwid_0", "some_hwid_1"}
	hostnames := []string{"some_hostname_0", "some_hostname_1"}
	hwidToHostname := map[string]string{
		hwids[0]: hostnames[0],
		hwids[1]: hostnames[1],
	}

	tks := []storage.TypeAndKey{
		{Type: dstorage.DirectorydTypeHWIDToHostname, Key: hwids[0]},
		{Type: dstorage.DirectorydTypeHWIDToHostname, Key: hwids[1]},
	}

	blobs := blobstore.Blobs{
		{
			Type:  tks[0].Type,
			Key:   tks[0].Key,
			Value: []byte(hostnames[0]),
		},
		{
			Type:  tks[1].Type,
			Key:   tks[1].Key,
			Value: []byte(hostnames[1]),
		},
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydBlobstore(blobFactMock)

	err := store.MapHWIDsToHostnames(hwidToHostname)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.PutRecord fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", placeholderNetworkID, mock.Anything, mock.Anything).
		Return(someErr).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	err = store.MapHWIDsToHostnames(hwidToHostname)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", placeholderNetworkID, blobs).
		Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	err = store.MapHWIDsToHostnames(hwidToHostname)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestDirectorydBlobstore_GetIMSIForSessionID(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	nid := "some_networkid"

	sid := "some_sessionid"
	tk := storage.TypeAndKey{Type: dstorage.DirectorydTypeSessionIDToIMSI, Key: sid}

	imsi := "some_imsi"
	blob := blobstore.Blob{
		Type:  tk.Type,
		Key:   tk.Key,
		Value: []byte(imsi),
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err := store.GetIMSIForSessionID(nid, sid)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", nid, tk).Return(blobstore.Blob{}, merrors.ErrNotFound).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err = store.GetIMSIForSessionID(nid, sid)
	assert.Exactly(t, merrors.ErrNotFound, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", nid, tk).Return(blobstore.Blob{}, someErr).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err = store.GetIMSIForSessionID(nid, sid)
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", nid, tk).Return(blob, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	imsiRecvd, err := store.GetIMSIForSessionID(nid, sid)
	assert.NoError(t, err)
	assert.Equal(t, imsi, imsiRecvd)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestDirectorydBlobstore_MapSessionIDToIMSI(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	nid := "some_networkid"

	sids := []string{"some_sessionid_0", "some_sessionid_1"}
	imsis := []string{"some_imsi_0", "some_imsi_1"}
	sidToIMSI := map[string]string{
		sids[0]: imsis[0],
		sids[1]: imsis[1],
	}

	tks := []storage.TypeAndKey{
		{Type: dstorage.DirectorydTypeSessionIDToIMSI, Key: sids[0]},
		{Type: dstorage.DirectorydTypeSessionIDToIMSI, Key: sids[1]},
	}

	blobs := blobstore.Blobs{
		{
			Type:  tks[0].Type,
			Key:   tks[0].Key,
			Value: []byte(imsis[0]),
		},
		{
			Type:  tks[1].Type,
			Key:   tks[1].Key,
			Value: []byte(imsis[1]),
		},
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydBlobstore(blobFactMock)

	err := store.MapSessionIDsToIMSIs(nid, sidToIMSI)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.PutRecord fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", nid, mock.Anything, mock.Anything).
		Return(someErr).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	err = store.MapSessionIDsToIMSIs(nid, sidToIMSI)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", nid, blobs).
		Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	err = store.MapSessionIDsToIMSIs(nid, sidToIMSI)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestDirectorydBlobstore_GetHWIDForSgwCTeid(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	nid := "some_networkid"

	teid := "1"
	tk := storage.TypeAndKey{Type: dstorage.DirectorydTypeSgwCteidToHwid, Key: teid}

	hwId := "hwId_1"
	blob := blobstore.Blob{
		Type:  tk.Type,
		Key:   tk.Key,
		Value: []byte(hwId),
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err := store.GetHWIDForSgwCTeid(nid, teid)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", nid, tk).Return(blobstore.Blob{}, merrors.ErrNotFound).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err = store.GetHWIDForSgwCTeid(nid, teid)
	assert.Exactly(t, merrors.ErrNotFound, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", nid, tk).Return(blobstore.Blob{}, someErr).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	_, err = store.GetHWIDForSgwCTeid(nid, teid)
	assert.Error(t, err)
	assert.NotEqual(t, merrors.ErrNotFound, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", nid, tk).Return(blob, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	hwIdRecvd, err := store.GetHWIDForSgwCTeid(nid, teid)
	assert.NoError(t, err)
	assert.Equal(t, hwId, hwIdRecvd)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestDirectorydBlobstore_MapSgwCTeidToHWID(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	nid := "some_networkid"

	teids := []string{"1", "2"}
	hwIds := []string{"hwid_1", "hwid_2"}
	teidsToHwIds := map[string]string{
		teids[0]: hwIds[0],
		teids[1]: hwIds[1],
	}

	tks := []storage.TypeAndKey{
		{Type: dstorage.DirectorydTypeSgwCteidToHwid, Key: teids[0]},
		{Type: dstorage.DirectorydTypeSgwCteidToHwid, Key: teids[1]},
	}

	blobs := blobstore.Blobs{
		{
			Type:  tks[0].Type,
			Key:   tks[0].Key,
			Value: []byte(hwIds[0]),
		},
		{
			Type:  tks[1].Type,
			Key:   tks[1].Key,
			Value: []byte(hwIds[1]),
		},
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := dstorage.NewDirectorydBlobstore(blobFactMock)

	err := store.MapSgwCTeidToHWID(nid, teidsToHwIds)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.PutRecord fails
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", nid, mock.Anything, mock.Anything).
		Return(someErr).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	err = store.MapSgwCTeidToHWID(nid, teidsToHwIds)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", nid, blobs).
		Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = dstorage.NewDirectorydBlobstore(blobFactMock)

	err = store.MapSgwCTeidToHWID(nid, teidsToHwIds)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}
