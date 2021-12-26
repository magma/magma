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
	"errors"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/blobstore/mocks"
	"magma/orc8r/cloud/go/services/certifier/constants"
	"magma/orc8r/cloud/go/services/certifier/protos"
	cstorage "magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
)

func TestCertifierBlobstore_GetCertInfo(t *testing.T) {
	var blobFactMock *mocks.StoreFactory
	var blobStoreMock *mocks.Store
	someErr := errors.New("generic error")

	serialNumber := "serial_number"
	info := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xdead, Nanos: 0xbeef},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xbeef, Nanos: 0xfeed},
		CertType:  1,
	}
	marshaledInfo, err := proto.Marshal(info)
	assert.NoError(t, err)
	tks := storage.TKs{{Type: constants.CertInfoType, Key: serialNumber}}
	blobs := blobstore.Blobs{{Type: constants.CertInfoType, Key: serialNumber, Value: marshaledInfo}}

	// Fail to start transaction
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetCertInfo(serialNumber)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with ErrNotFound
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobstore.Blobs{}, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetCertInfo(serialNumber)
	assert.Exactly(t, err, merrors.ErrNotFound)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Get fails with error other than ErrNotFound
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobstore.Blobs{}, someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetCertInfo(serialNumber)
	assert.Error(t, err)
	assert.NotEqual(t, err, merrors.ErrNotFound)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobs, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	infoRecvd, err := store.GetCertInfo(serialNumber)
	assert.NoError(t, err)
	assert.True(t, proto.Equal(info, infoRecvd))
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCertifierBlobstore_GetManyCertInfo(t *testing.T) {
	var blobFactMock *mocks.StoreFactory
	var blobStoreMock *mocks.Store
	someErr := errors.New("generic error")

	info0 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xdead, Nanos: 0xbeef},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xbeef, Nanos: 0xfeed},
		CertType:  1,
	}
	info1 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xaaaa, Nanos: 0xbbbb},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xcccc, Nanos: 0xdddd},
		CertType:  2,
	}
	marshaledInfo0, err := proto.Marshal(info0)
	assert.NoError(t, err)
	marshaledInfo1, err := proto.Marshal(info1)
	assert.NoError(t, err)

	serialNumbers := []string{"serial_number_0", "serial_number_1"}
	tks := storage.TKs{
		{Type: constants.CertInfoType, Key: serialNumbers[0]},
		{Type: constants.CertInfoType, Key: serialNumbers[1]},
	}
	blobs := blobstore.Blobs{
		{Type: constants.CertInfoType, Key: serialNumbers[0], Value: marshaledInfo0},
		{Type: constants.CertInfoType, Key: serialNumbers[1], Value: marshaledInfo1},
	}
	infos := map[string]*protos.CertificateInfo{
		serialNumbers[0]: info0,
		serialNumbers[1]: info1,
	}

	// Fail to start transaction
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetManyCertInfo(serialNumbers)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.GetMany fails
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobstore.Blobs{}, someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetManyCertInfo(serialNumbers)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.GetMany succeeds with empty return
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobstore.Blobs{}, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	infosRecvds, err := store.GetManyCertInfo(serialNumbers)
	assert.NoError(t, err)
	assert.Equal(t, map[string]*protos.CertificateInfo{}, infosRecvds)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobs, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	infosRecvds, err = store.GetManyCertInfo(serialNumbers)
	assert.NoError(t, err)
	assert.Len(t, infosRecvds, 2)
	for k := range infos {
		a, b := infos[k], infosRecvds[k]
		assert.True(t, proto.Equal(a, b))
	}
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCertifierBlobstore_GetAllCertInfo(t *testing.T) {
	var blobFactMock *mocks.StoreFactory
	var blobStoreMock *mocks.Store
	someErr := errors.New("generic error")

	info0 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xdead, Nanos: 0xbeef},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xbeef, Nanos: 0xfeed},
		CertType:  1,
	}
	info1 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xaaaa, Nanos: 0xbbbb},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xcccc, Nanos: 0xdddd},
		CertType:  2,
	}
	marshaledInfo0, err := proto.Marshal(info0)
	assert.NoError(t, err)
	marshaledInfo1, err := proto.Marshal(info1)
	assert.NoError(t, err)

	serialNumbers := []string{"serial_number_0", "serial_number_1"}
	tks := storage.TKs{
		{Type: constants.CertInfoType, Key: serialNumbers[0]},
		{Type: constants.CertInfoType, Key: serialNumbers[1]},
	}
	blobs := blobstore.Blobs{
		{Type: constants.CertInfoType, Key: serialNumbers[0], Value: marshaledInfo0},
		{Type: constants.CertInfoType, Key: serialNumbers[1], Value: marshaledInfo1},
	}
	infos := map[string]*protos.CertificateInfo{
		serialNumbers[0]: info0,
		serialNumbers[1]: info1,
	}
	placeHolderNetwork := "placeholder_network"
	searchResult := map[string]blobstore.Blobs{placeHolderNetwork: blobs}

	// Fail to start transaction
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetAllCertInfo()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.ListKeys fails
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On(
		"Search",
		blobstore.CreateSearchFilter(&placeHolderNetwork, []string{constants.CertInfoType}, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(map[string]blobstore.Blobs{}, someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetAllCertInfo()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.ListKeys succeeds with empty return
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On(
		"Search",
		blobstore.CreateSearchFilter(&placeHolderNetwork, []string{constants.CertInfoType}, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(map[string]blobstore.Blobs{}, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	infosRecvds, err := store.GetAllCertInfo()
	assert.NoError(t, err)
	assert.Equal(t, map[string]*protos.CertificateInfo{}, infosRecvds)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.GetMany fails
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On(
		"Search",
		blobstore.CreateSearchFilter(&placeHolderNetwork, []string{constants.CertInfoType}, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(searchResult, nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobstore.Blobs{}, someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.GetAllCertInfo()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On(
		"Search",
		blobstore.CreateSearchFilter(&placeHolderNetwork, []string{constants.CertInfoType}, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(searchResult, nil).Once()
	blobStoreMock.On("GetMany", mock.Anything, tks).Return(blobs, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	infosRecvds, err = store.GetAllCertInfo()
	assert.NoError(t, err)
	assert.Len(t, infosRecvds, 2)
	for k := range infos {
		a, b := infos[k], infosRecvds[k]
		assert.True(t, proto.Equal(a, b))
	}
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCertifierBlobstore_PutCertInfo(t *testing.T) {
	var blobFactMock *mocks.StoreFactory
	var blobStoreMock *mocks.Store

	someErr := errors.New("generic error")

	serialNumber := "serial_number"
	info := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xdead, Nanos: 0xbeef},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xbeef, Nanos: 0xfeed},
		CertType:  1,
	}
	marshaledInfos, err := proto.Marshal(info)
	assert.NoError(t, err)
	blob := blobstore.Blob{Type: constants.CertInfoType, Key: serialNumber, Value: marshaledInfos}

	// Fail to start transaction
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	err = store.PutCertInfo(serialNumber, info)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Put fails
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Write", mock.Anything, blobstore.Blobs{blob}).Return(someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	err = store.PutCertInfo(serialNumber, info)
	assert.Error(t, err)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Write", mock.Anything, blobstore.Blobs{blob}).
		Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	err = store.PutCertInfo(serialNumber, info)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCertifierBlobstore_DeleteCertInfo(t *testing.T) {
	var blobFactMock *mocks.StoreFactory
	var blobStoreMock *mocks.Store

	someErr := errors.New("generic error")

	serialNumber := "serial_number"
	tks := storage.TKs{{Type: constants.CertInfoType, Key: serialNumber}}

	// Fail to start transaction
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	err := store.DeleteCertInfo(serialNumber)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.Delete fails
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Delete", mock.Anything, tks).Return(someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	err = store.DeleteCertInfo(serialNumber)
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Delete", mock.Anything, tks).Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	err = store.DeleteCertInfo(serialNumber)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}

func TestCertifierBlobstore_ListSerialNumbers(t *testing.T) {
	var blobFactMock *mocks.StoreFactory
	var blobStoreMock *mocks.Store

	someErr := errors.New("generic error")

	info0 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xdead, Nanos: 0xbeef},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xbeef, Nanos: 0xfeed},
		CertType:  1,
	}
	info1 := &protos.CertificateInfo{
		Id:        nil,
		NotBefore: &timestamp.Timestamp{Seconds: 0xaaaa, Nanos: 0xbbbb},
		NotAfter:  &timestamp.Timestamp{Seconds: 0xcccc, Nanos: 0xdddd},
		CertType:  2,
	}
	marshaledInfo0, err := proto.Marshal(info0)
	assert.NoError(t, err)
	marshaledInfo1, err := proto.Marshal(info1)
	assert.NoError(t, err)

	blobs := blobstore.Blobs{
		{Type: constants.CertInfoType, Key: "serial_number_0", Value: marshaledInfo0},
		{Type: constants.CertInfoType, Key: "serial_number_1", Value: marshaledInfo1},
	}
	placeHolderNetwork := "placeholder_network"
	searchResult := map[string]blobstore.Blobs{placeHolderNetwork: blobs}

	serialNumbers := blobs.Keys()
	assert.Equal(t, serialNumbers, []string{"serial_number_0", "serial_number_1"})

	// Fail to start transaction
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(nil, someErr).Once()
	store := cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.ListSerialNumbers()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// store.ListKeys fails
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On(
		"Search",
		blobstore.CreateSearchFilter(&placeHolderNetwork, []string{constants.CertInfoType}, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(map[string]blobstore.Blobs{}, someErr).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	_, err = store.ListSerialNumbers()
	assert.Error(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Success
	blobFactMock = &mocks.StoreFactory{}
	blobStoreMock = &mocks.Store{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On(
		"Search",
		blobstore.CreateSearchFilter(&placeHolderNetwork, []string{constants.CertInfoType}, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(searchResult, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = cstorage.NewCertifierBlobstore(blobFactMock)

	serialNumbersRecvd, err := store.ListSerialNumbers()
	assert.NoError(t, err)
	assert.Equal(t, serialNumbers, serialNumbersRecvd)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}
