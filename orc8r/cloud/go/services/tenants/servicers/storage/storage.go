package storage

import (
	"fmt"
	"strconv"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/tenants"
	tenant_protos "magma/orc8r/cloud/go/services/tenants/protos"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
)

const networkWildcard = "*"

type Store interface {
	CreateTenant(tenantID int64, tenant *tenant_protos.Tenant) error
	GetTenant(tenantID int64) (*tenant_protos.Tenant, error)
	GetAllTenants() (*tenant_protos.TenantList, error)
	SetTenant(tenantID int64, tenant *tenant_protos.Tenant) error
	DeleteTenant(tenantID int64) error
	GetControlProxy(tenantID int64) (string, error)
	CreateOrUpdateControlProxy(tenantID int64, controlProxy string) error
}

type blobstoreStore struct {
	factory blobstore.StoreFactory
}

func NewBlobstoreStore(factory blobstore.StoreFactory) Store {
	return &blobstoreStore{factory}
}

func (b *blobstoreStore) CreateTenant(tenantID int64, tenant *tenant_protos.Tenant) error {
	return b.SetTenant(tenantID, tenant)
}

func (b *blobstoreStore) GetTenant(tenantID int64) (*tenant_protos.Tenant, error) {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	defer store.Rollback()

	tenantTK := storage.TK{
		Type: tenants.TenantInfoType,
		Key:  strconv.FormatInt(tenantID, 10),
	}
	tenantBlob, err := store.Get(networkWildcard, tenantTK)
	if err != nil {
		return nil, err
	}
	retTenant, err := tenantFromBlob(tenantBlob)
	if err != nil {
		return nil, err
	}
	return retTenant, store.Commit()
}

func (b *blobstoreStore) GetAllTenants() (*tenant_protos.TenantList, error) {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	defer store.Rollback()

	keys, err := blobstore.ListKeys(store, networkWildcard, tenants.TenantInfoType)
	if err != nil {
		return nil, err
	}

	keysAndTypes := make(storage.TKs, 0)
	for _, key := range keys {
		keysAndTypes = append(keysAndTypes, storage.TK{Key: key, Type: tenants.TenantInfoType})
	}

	tenantBlobs, err := store.GetMany(networkWildcard, keysAndTypes)
	if err != nil {
		return nil, err
	}

	retTenants := &tenant_protos.TenantList{}
	for _, blob := range tenantBlobs {
		tenant, err := tenantFromBlob(blob)
		if err != nil {
			return nil, err
		}
		intID, err := strconv.ParseInt(blob.Key, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("non-integer key: %v", err)
		}
		idAndTenant := &tenant_protos.IDAndTenant{
			Id:     intID,
			Tenant: tenant,
		}
		retTenants.Tenants = append(retTenants.Tenants, idAndTenant)
	}
	return retTenants, nil
}

func (b *blobstoreStore) SetTenant(tenantID int64, tenant *tenant_protos.Tenant) error {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return err
	}
	defer store.Rollback()

	tenantBlob, err := tenantToBlob(tenantID, tenant)
	if err != nil {
		return err
	}
	err = store.Write(networkWildcard, blobstore.Blobs{tenantBlob})
	if err != nil {
		return err
	}
	return store.Commit()
}

func (b *blobstoreStore) DeleteTenant(tenantID int64) error {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return err
	}
	defer store.Rollback()

	tenantTK := storage.TKs{{
		Type: tenants.TenantInfoType,
		Key:  strconv.FormatInt(tenantID, 10),
	}}
	err = store.Delete(networkWildcard, tenantTK)
	if err != nil {
		return err
	}
	return store.Commit()
}

func (b *blobstoreStore) GetControlProxy(tenantID int64) (string, error) {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return "", err
	}
	defer store.Rollback()

	controlProxyTK := storage.TK{
		Type: tenants.ControlProxyInfoType,
		Key:  strconv.FormatInt(tenantID, 10),
	}
	controlProxyBlob, err := store.Get(networkWildcard, controlProxyTK)
	if err != nil {
		return "", err
	}

	return string(controlProxyBlob.Value), store.Commit()
}

func (b *blobstoreStore) CreateOrUpdateControlProxy(tenantID int64, controlProxy string) error {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return err
	}
	defer store.Rollback()

	controlProxyBlob := blobstore.Blob{
		Type:  tenants.ControlProxyInfoType,
		Key:   strconv.FormatInt(tenantID, 10),
		Value: []byte(controlProxy),
	}
	err = store.Write(networkWildcard, blobstore.Blobs{controlProxyBlob})
	if err != nil {
		return err
	}

	return store.Commit()
}

func tenantToBlob(tenantID int64, tenant *tenant_protos.Tenant) (blobstore.Blob, error) {
	marshaledTenant, err := protos.Marshal(tenant)
	if err != nil {
		return blobstore.Blob{}, fmt.Errorf("Error marshaling protobuf: %w", err)
	}
	return blobstore.Blob{
		Type:  tenants.TenantInfoType,
		Key:   strconv.FormatInt(tenantID, 10),
		Value: marshaledTenant,
	}, nil
}

func tenantFromBlob(blob blobstore.Blob) (*tenant_protos.Tenant, error) {
	tenant := tenant_protos.Tenant{}
	err := protos.Unmarshal(blob.Value, &tenant)
	if err != nil {
		return &tenant_protos.Tenant{}, fmt.Errorf("Error unmarshaling protobuf: %w", err)
	}
	return &tenant, nil
}
