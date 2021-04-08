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
	cstorage "magma/orc8r/cloud/go/services/ctraced/storage"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	placeholderNetworkID = "placeholder_network"
)

func TestCtracedBlobstoreStorage_GetCallTrace(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	ctid := "some_call_trace"
	tk := storage.TypeAndKey{Type: cstorage.CtracedBlobType, Key: ctid}

	ctData := "abcdefghijklmnopqrstuvwxyz"
	blob := blobstore.Blob{
		Type:  tk.Type,
		Key:   tk.Key,
		Value: []byte(ctData),
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCtracedBlobstore(blobFactMock)

	_, err := store.GetCallTrace(placeholderNetworkID, ctid)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", placeholderNetworkID, tk).Return(blobstore.Blob{}, merrors.ErrNotFound).Once()
	store = cstorage.NewCtracedBlobstore(blobFactMock)

	_, err = store.GetCallTrace(placeholderNetworkID, ctid)
	assert.Exactly(t, merrors.ErrNotFound, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", placeholderNetworkID, tk).Return(blobstore.Blob{}, someErr).Once()
	store = cstorage.NewCtracedBlobstore(blobFactMock)

	_, err = store.GetCallTrace(placeholderNetworkID, ctid)
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
	store = cstorage.NewCtracedBlobstore(blobFactMock)

	callTraceRecvd, err := store.GetCallTrace(placeholderNetworkID, ctid)
	assert.NoError(t, err)
	assert.Equal(t, ctData, string(callTraceRecvd))
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCtracedBlobstoreStorage_StoreCallTrace(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	ctid := "some_call_trace"
	tk := storage.TypeAndKey{Type: cstorage.CtracedBlobType, Key: ctid}

	ctData := []byte("abcdefghijklmnopqrstuvwxyz")
	blob := blobstore.Blob{
		Type:  tk.Type,
		Key:   tk.Key,
		Value: ctData,
	}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCtracedBlobstore(blobFactMock)

	err := store.StoreCallTrace(placeholderNetworkID, ctid, ctData)
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
	store = cstorage.NewCtracedBlobstore(blobFactMock)

	err = store.StoreCallTrace(placeholderNetworkID, ctid, ctData)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", placeholderNetworkID, blobstore.Blobs{blob}).
		Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCtracedBlobstore(blobFactMock)

	err = store.StoreCallTrace(placeholderNetworkID, ctid, ctData)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCtracedBlobstoreStorage_DeleteCallTrace(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage
	someErr := errors.New("generic error")

	ctid := "some_call_trace"
	tk := storage.TypeAndKey{Type: cstorage.CtracedBlobType, Key: ctid}
	tkSet := []storage.TypeAndKey{tk}

	// Fail to start transaction
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCtracedBlobstore(blobFactMock)

	err := store.DeleteCallTrace(placeholderNetworkID, ctid)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Delete", placeholderNetworkID, tkSet).Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCtracedBlobstore(blobFactMock)

	err = store.DeleteCallTrace(placeholderNetworkID, ctid)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}
