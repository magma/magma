package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/blobstore/mocks"
	"magma/orc8r/cloud/go/services/tenants"
	tenant_protos "magma/orc8r/cloud/go/services/tenants/protos"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

var (
	sampleTenant0        = tenant_protos.Tenant{Name: "tenant0", Networks: []string{"net1", "net2"}}
	sampleTenant0Blob, _ = tenantToBlob(0, &sampleTenant0)

	sampleTenant1        = tenant_protos.Tenant{Name: "tenant1", Networks: []string{"net3", "net4"}}
	sampleTenant1Blob, _ = tenantToBlob(1, &sampleTenant1)

	marshaledTenant0, _ = protos.Marshal(&sampleTenant0)
	invalidBlob         = blobstore.Blob{
		Type:  tenants.TenantInfoType,
		Key:   "word",
		Value: marshaledTenant0,
	}

	sampleTenantID      int64 = 0
	sampleTenantID2     int64 = 3
	sampleControlProxy        = "{ info }"
	sampleControlProxy2       = "{ info2 }"
)

func setupTestStore() (*mocks.Store, Store) {
	store := &mocks.Store{}
	store.On("Rollback").Return(nil)
	store.On("Commit").Return(nil)

	factory := &mocks.StoreFactory{}
	factory.On("StartTransaction", mock.Anything).Return(store, nil)

	return store, NewBlobstoreStore(factory)
}

func TestBlobstoreStore_CreateTenant(t *testing.T) {
	txStore, s := setupTestStore()
	txStore.On("Write", networkWildcard, blobstore.Blobs{sampleTenant0Blob}).Return(nil)
	err := s.CreateTenant(0, &sampleTenant0)
	assert.NoError(t, err)

	txStore, s = setupTestStore()
	txStore.On("Write", networkWildcard, blobstore.Blobs{sampleTenant0Blob}).Return(errors.New("error"))
	err = s.CreateTenant(0, &sampleTenant0)
	assert.EqualError(t, err, "error")
}

func TestBlobstoreStore_GetTenant(t *testing.T) {
	txStore, s := setupTestStore()
	txStore.On("Get", networkWildcard, storage.TK{Type: tenants.TenantInfoType, Key: "0"}).Return(sampleTenant0Blob, nil)
	tenant, err := s.GetTenant(0)
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, &sampleTenant0, tenant)

	txStore, s = setupTestStore()
	txStore.On("Get", networkWildcard, storage.TK{Type: tenants.TenantInfoType, Key: "0"}).Return(blobstore.Blob{}, errors.New("error"))
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
	txStore.On("GetMany", networkWildcard, storage.TKs{
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
	test_utils.AssertMessagesEqual(t, &tenant_protos.TenantList{Tenants: []*tenant_protos.IDAndTenant{
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
	txStore.On("GetMany", networkWildcard, storage.TKs{
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
	txStore.On("GetMany", networkWildcard, storage.TKs{
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
	txStore.On("Write", networkWildcard, blobstore.Blobs{sampleTenant0Blob}).Return(nil)
	err := s.SetTenant(0, &sampleTenant0)
	assert.NoError(t, err)

	txStore, s = setupTestStore()
	txStore.On("Write", networkWildcard, blobstore.Blobs{sampleTenant0Blob}).Return(errors.New("error"))
	err = s.SetTenant(0, &sampleTenant0)
	assert.EqualError(t, err, "error")
}

func TestBlobstoreStore_DeleteTenant(t *testing.T) {
	txStore, s := setupTestStore()
	txStore.On("Delete", networkWildcard, storage.TKs{{Type: tenants.TenantInfoType, Key: "0"}}).Return(nil)
	err := s.DeleteTenant(0)
	assert.NoError(t, err)

	txStore, s = setupTestStore()
	txStore.On("Delete", networkWildcard, storage.TKs{{Type: tenants.TenantInfoType, Key: "0"}}).Return(errors.New("error"))
	err = s.DeleteTenant(0)
	assert.EqualError(t, err, "error")
}

func TestBlobstoreStore_ControlProxy(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	factory := blobstore.NewSQLStoreFactory(tenants.DBTableName, db, sqorc.GetSqlBuilder())
	assert.NoError(t, factory.InitializeFactory())
	s := NewBlobstoreStore(factory)

	_, err = s.GetControlProxy(sampleTenantID)
	assert.EqualError(t, err, "Not found")

	_, err = s.GetControlProxy(sampleTenantID2)
	assert.EqualError(t, err, "Not found")

	err = s.CreateOrUpdateControlProxy(sampleTenantID2, sampleControlProxy)
	assert.NoError(t, err)

	controlProxy, err := s.GetControlProxy(sampleTenantID2)
	assert.NoError(t, err)
	assert.Equal(t, sampleControlProxy, controlProxy)

	_, err = s.GetControlProxy(sampleTenantID)
	assert.EqualError(t, err, "Not found")

	err = s.CreateOrUpdateControlProxy(sampleTenantID, sampleControlProxy2)
	assert.NoError(t, err)

	controlProxy, err = s.GetControlProxy(sampleTenantID)
	assert.NoError(t, err)
	assert.Equal(t, sampleControlProxy2, controlProxy)
}
