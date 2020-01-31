package storage

import (
	"encoding/json"
	"fmt"
	"strconv"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/cloud/go/storage"
)

const networkWildcard = "*"

type Store interface {
	CreateTenant(tenantID uint64, tenant protos.Tenant) error
	GetTenant(tenantID uint64) (*protos.Tenant, error)
	GetAllTenants() (*protos.TenantList, error)
	SetTenant(tenantID uint64, tenant protos.Tenant) error
	DeleteTenant(tenantID uint64) error
}

type blobstoreStore struct {
	factory blobstore.BlobStorageFactory
}

func NewBlobstoreStore(factory blobstore.BlobStorageFactory) Store {
	return &blobstoreStore{factory}
}

func (b *blobstoreStore) CreateTenant(tenantID uint64, tenant protos.Tenant) error {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return err
	}
	tenantBlob, err := tenantToBlob(tenantID, tenant)
	err = store.CreateOrUpdate(networkWildcard, []blobstore.Blob{tenantBlob})
	if err != nil {
		store.Rollback()
		return err
	}
	return store.Commit()
}

func (b *blobstoreStore) GetTenant(tenantID uint64) (*protos.Tenant, error) {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	tenantTypeAndKey := storage.TypeAndKey{
		Type: tenants.TenantInfoType,
		Key:  strconv.FormatUint(tenantID, 10),
	}
	tenantBlob, err := store.Get(networkWildcard, tenantTypeAndKey)
	if err != nil {
		store.Rollback()
		return nil, err
	}
	retTenant := &protos.Tenant{}
	err = protos.Unmarshal(tenantBlob.Value, retTenant)
	if err != nil {
		store.Rollback()
		return nil, err
	}
	return retTenant, store.Commit()
}

func (b *blobstoreStore) GetAllTenants() (*protos.TenantList, error) {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	keys, err := store.ListKeys(networkWildcard, tenants.TenantInfoType)
	if err != nil {
		store.Rollback()
		return nil, err
	}

	keysAndTypes := make([]storage.TypeAndKey, 0)
	for _, key := range keys {
		keysAndTypes = append(keysAndTypes, storage.TypeAndKey{Key: key, Type: tenants.TenantInfoType})
	}

	tenantBlobs, err := store.GetMany(networkWildcard, keysAndTypes)
	retTenants := &protos.TenantList{}
	for _, blob := range tenantBlobs {
		tenant := protos.Tenant{}
		err = json.Unmarshal(blob.Value, &tenant)
		if err != nil {
			store.Rollback()
			return nil, err
		}
		intID, err := strconv.ParseUint(blob.Key, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("non-integer key: %v\n", err)
		}
		idAndTenant := &protos.IDAndTenant{
			Id:     intID,
			Tenant: &tenant,
		}
		retTenants.Tenants = append(retTenants.Tenants, idAndTenant)
	}
	return retTenants, nil
}

func (b *blobstoreStore) SetTenant(tenantID uint64, tenant protos.Tenant) error {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return err
	}
	tenantBlob, err := tenantToBlob(tenantID, tenant)
	if err != nil {
		return err
	}
	err = store.CreateOrUpdate(networkWildcard, []blobstore.Blob{tenantBlob})
	if err != nil {
		store.Rollback()
		return err
	}
	return store.Commit()
}

func (b *blobstoreStore) DeleteTenant(tenantID uint64) error {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return err
	}
	tenantTypeAndKey := []storage.TypeAndKey{{
		Type: tenants.TenantInfoType,
		Key:  strconv.FormatUint(tenantID, 10),
	}}
	err = store.Delete(networkWildcard, tenantTypeAndKey)
	if err != nil {
		store.Rollback()
		return err
	}
	return store.Commit()
}

func tenantToBlob(tenantID uint64, tenant protos.Tenant) (blobstore.Blob, error) {
	marshaledTenant, err := protos.Marshal(&tenant)
	if err != nil {
		return blobstore.Blob{}, err
	}
	return blobstore.Blob{
		Type:  tenants.TenantInfoType,
		Key:   strconv.FormatUint(tenantID, 10),
		Value: marshaledTenant,
	}, nil
}
