package storage

import (
	"errors"
	"testing"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/blobstore/mocks"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/cloud/go/storage"

	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	sampleTenant0        = protos.Tenant{Name: "tenant0", Networks: []string{"net1", "net2"}}
	sampleTenant0Blob, _ = tenantToBlob(0, sampleTenant0)

	sampleTenant1        = protos.Tenant{Name: "tenant1", Networks: []string{"net3", "net4"}}
	sampleTenant1Blob, _ = tenantToBlob(1, sampleTenant1)

	marshaledTenant0, _ = protos.Marshal(&sampleTenant0)
	invalidBlob         = blobstore.Blob{
		Type:  tenants.TenantInfoType,
		Key:   "word",
		Value: marshaledTenant0,
	}
)

func setupTestStore() (*mocks.TransactionalBlobStorage, Store) {
	store := &mocks.TransactionalBlobStorage{}
	store.On("Rollback").Return(nil)
	store.On("Commit").Return(nil)

	factory := &mocks.BlobStorageFactory{}
	factory.On("StartTransaction", mock.Anything).Return(store, nil)

	return store, NewBlobstoreStore(factory)
}

func TestBlobstoreStore_CreateTenant(t *testing.T) {
	txStore, s := setupTestStore()
	txStore.On("CreateOrUpdate", networkWildcard, blobstore.Blobs{sampleTenant0Blob}).Return(nil)
	err := s.CreateTenant(0, sampleTenant0)
	assert.NoError(t, err)

	txStore, s = setupTestStore()
	txStore.On("CreateOrUpdate", networkWildcard, blobstore.Blobs{sampleTenant0Blob}).Return(errors.New("error"))
	err = s.CreateTenant(0, sampleTenant0)
	assert.EqualError(t, err, "error")
}

func TestBlobstoreStore_GetTenant(t *testing.T) {
	txStore, s := setupTestStore()
	txStore.On("Get", networkWildcard, storage.TypeAndKey{Type: tenants.TenantInfoType, Key: "0"}).Return(sampleTenant0Blob, nil)
	tenant, err := s.GetTenant(0)
	assert.NoError(t, err)
	assert.Equal(t, sampleTenant0, *tenant)

	txStore, s = setupTestStore()
	txStore.On("Get", networkWildcard, storage.TypeAndKey{Type: tenants.TenantInfoType, Key: "0"}).Return(blobstore.Blob{}, errors.New("error"))
	_, err = s.GetTenant(0)
	assert.EqualError(t, err, "error")
}

func TestBlobstoreStore_GetAllTenants(t *testing.T) {
	networkWildCard := "*"
	completeSearchResult := map[string]blobstore.Blobs{networkWildCard: {sampleTenant0Blob, sampleTenant1Blob}}
	partialSearchResult := map[string]blobstore.Blobs{networkWildCard: {sampleTenant0Blob}}

	// Successful GetAll
	txStore, s := setupTestStore()

	txStore.On(
		"Search",
		blobstore.CreateSearchFilter(&networkWildCard, []string{tenants.TenantInfoType}, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(completeSearchResult, nil)
	txStore.On("GetMany", networkWildcard, []storage.TypeAndKey{
		{
			Type: tenants.TenantInfoType,
			Key:  "0",
		}, {
			Type: tenants.TenantInfoType,
			Key:  "1",
		},
	}).Return(blobstore.Blobs{sampleTenant0Blob, sampleTenant1Blob}, nil)

	retTenants, err := s.GetAllTenants()
	assert.NoError(t, err)
	assert.Equal(t, &protos.TenantList{Tenants: []*protos.IDAndTenant{
		{Id: 0, Tenant: &sampleTenant0},
		{Id: 1, Tenant: &sampleTenant1},
	}}, retTenants)

	// Error in ListKeys
	txStore, s = setupTestStore()
	txStore.On(
		"Search",
		blobstore.CreateSearchFilter(&networkWildCard, []string{tenants.TenantInfoType}, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(map[string]blobstore.Blobs{}, errors.New("error"))
	_, err = s.GetAllTenants()
	assert.EqualError(t, err, "error")

	// Error in GetMany
	txStore, s = setupTestStore()
	txStore.On(
		"Search",
		blobstore.CreateSearchFilter(&networkWildCard, []string{tenants.TenantInfoType}, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(partialSearchResult, nil)
	txStore.On("GetMany", networkWildcard, []storage.TypeAndKey{
		{
			Type: tenants.TenantInfoType,
			Key:  "0",
		},
	}).Return(blobstore.Blobs{}, errors.New("error"))
	_, err = s.GetAllTenants()
	assert.EqualError(t, err, "error")

	// Non-integer key in tenant
	txStore, s = setupTestStore()
	txStore.On(
		"Search",
		blobstore.CreateSearchFilter(&networkWildCard, []string{tenants.TenantInfoType}, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	).Return(partialSearchResult, nil)
	txStore.On("GetMany", networkWildcard, []storage.TypeAndKey{
		{
			Type: tenants.TenantInfoType,
			Key:  "0",
		},
	}).Return(blobstore.Blobs{invalidBlob}, nil)
	_, err = s.GetAllTenants()
	assert.EqualError(t, err, `non-integer key: strconv.ParseInt: parsing "word": invalid syntax`)
}

func TestBlobstoreStore_SetTenant(t *testing.T) {
	txStore, s := setupTestStore()
	txStore.On("CreateOrUpdate", networkWildcard, blobstore.Blobs{sampleTenant0Blob}).Return(nil)
	err := s.SetTenant(0, sampleTenant0)
	assert.NoError(t, err)

	txStore, s = setupTestStore()
	txStore.On("CreateOrUpdate", networkWildcard, blobstore.Blobs{sampleTenant0Blob}).Return(errors.New("error"))
	err = s.SetTenant(0, sampleTenant0)
	assert.EqualError(t, err, "error")
}

func TestBlobstoreStore_DeleteTenant(t *testing.T) {
	txStore, s := setupTestStore()
	txStore.On("Delete", networkWildcard, []storage.TypeAndKey{{Type: tenants.TenantInfoType, Key: "0"}}).Return(nil)
	err := s.DeleteTenant(0)
	assert.NoError(t, err)

	txStore, s = setupTestStore()
	txStore.On("Delete", networkWildcard, []storage.TypeAndKey{{Type: tenants.TenantInfoType, Key: "0"}}).Return(errors.New("error"))
	err = s.DeleteTenant(0)
	assert.EqualError(t, err, "error")
}
