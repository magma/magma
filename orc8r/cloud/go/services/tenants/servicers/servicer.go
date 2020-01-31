/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"context"
	"fmt"

	"magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"

	"magma/orc8r/cloud/go/services/tenants/servicers/storage"
)

type tenantsServicer struct {
	store storage.Store
}

// NewTenantsServicer returns a state server backed by storage passed in
func NewTenantsServicer(store storage.Store) (protos.TenantsServiceServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Storage store is nil")
	}
	return &tenantsServicer{store}, nil
}

func (s *tenantsServicer) GetAllTenants(c context.Context, _ *protos.Void) (*protos.TenantList, error) {
	return s.store.GetAllTenants()
}

func (s *tenantsServicer) CreateTenant(c context.Context, request *protos.IDAndTenant) (*protos.Void, error) {
	if _, err := s.store.GetTenant(request.Id); err != errors.ErrNotFound {
		return nil, errors.ErrAlreadyExists
	}
	return &protos.Void{}, s.store.CreateTenant(request.Id, *request.Tenant)
}

func (s *tenantsServicer) GetTenant(c context.Context, request *protos.GetTenantRequest) (*protos.Tenant, error) {
	return s.store.GetTenant(request.Id)

}

func (s *tenantsServicer) SetTenant(c context.Context, request *protos.IDAndTenant) (*protos.Void, error) {
	if _, err := s.store.GetTenant(request.Id); err == errors.ErrNotFound {
		return nil, err
	}
	return &protos.Void{}, s.store.SetTenant(request.Id, *request.Tenant)
}

func (s *tenantsServicer) DeleteTenant(c context.Context, request *protos.GetTenantRequest) (*protos.Void, error) {
	return &protos.Void{}, s.store.DeleteTenant(request.Id)
}
