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

package internal_servicer

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tenant_protos "magma/orc8r/cloud/go/services/tenants/protos"
	"magma/orc8r/cloud/go/services/tenants/servicers/storage"
	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
)

type tenantsServicer struct {
	store storage.Store
}

// NewTenantsServicer returns a state server backed by storage passed in
func NewTenantsServicer(store storage.Store) (tenant_protos.TenantsServiceServer, error) {
	if store == nil {
		return nil, fmt.Errorf("Storage store is nil")
	}
	return &tenantsServicer{store}, nil
}

func (s *tenantsServicer) GetAllTenants(c context.Context, _ *protos.Void) (*tenant_protos.TenantList, error) {
	tenants, err := s.store.GetAllTenants()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting tenants: %v", err)
	}
	return tenants, nil
}

func (s *tenantsServicer) CreateTenant(c context.Context, request *tenant_protos.IDAndTenant) (*protos.Void, error) {
	_, err := s.store.GetTenant(request.Id)
	switch {
	case err == nil:
		return nil, status.Errorf(codes.AlreadyExists, "Tenant with Id %d already exists", request.Id)
	case err != errors.ErrNotFound:
		return nil, status.Errorf(codes.Internal, "Error getting existing tenants: %v", err)
	}

	err = s.store.CreateTenant(request.Id, *request.Tenant)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error creating tenant: %v", err)
	}
	return &protos.Void{}, nil
}

func (s *tenantsServicer) GetTenant(c context.Context, request *tenant_protos.GetTenantRequest) (*tenant_protos.Tenant, error) {
	tenant, err := s.store.GetTenant(request.Id)
	if err != nil {
		return nil, mapErrForGet(err, request.Id, "tenant")
	}
	return tenant, nil

}

func (s *tenantsServicer) SetTenant(c context.Context, request *tenant_protos.IDAndTenant) (*protos.Void, error) {
	err := s.validateTenantExists(request.Id)
	if err != nil {
		return nil, err
	}

	err = s.store.SetTenant(request.Id, *request.Tenant)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error setting tenant %d: %v", request.Id, err)
	}

	return &protos.Void{}, nil
}

func (s *tenantsServicer) DeleteTenant(c context.Context, request *tenant_protos.GetTenantRequest) (*protos.Void, error) {
	err := s.store.DeleteTenant(request.Id)
	switch {
	case err == errors.ErrNotFound:
		return nil, status.Errorf(codes.NotFound, "Tenant %d not found", request.Id)
	case err != nil:
		return nil, status.Errorf(codes.Internal, "Error deleting tenant %d: %v", request.Id, err)
	}
	return &protos.Void{}, nil
}

func (s *tenantsServicer) GetControlProxy(c context.Context, request *tenant_protos.GetControlProxyRequest) (*tenant_protos.GetControlProxyResponse, error) {
	err := s.validateTenantExists(request.Id)
	if err != nil {
		return nil, err
	}

	controlProxy, err := s.store.GetControlProxy(request.Id)
	if err != nil {
		return nil, mapErrForGet(err, request.Id, "controlProxy")
	}

	return &tenant_protos.GetControlProxyResponse{Id: request.Id, ControlProxy: controlProxy}, nil
}

func (s *tenantsServicer) CreateOrUpdateControlProxy(c context.Context, request *tenant_protos.CreateOrUpdateControlProxyRequest) (*tenant_protos.CreateOrUpdateControlProxyResponse, error) {
	err := s.validateTenantExists(request.Id)
	if err != nil {
		return nil, err
	}

	err = s.store.CreateOrUpdateControlProxy(request.Id, request.ControlProxy)
	if err != nil {
		request.GetControlProxy()
		return nil, status.Errorf(codes.Internal, "Error setting control proxy %d: %v", request.Id, err)
	}

	return &tenant_protos.CreateOrUpdateControlProxyResponse{}, nil
}

func mapErrForGet(err error, id int64, getType string) error {
	switch {
	case err == errors.ErrNotFound:
		return status.Errorf(codes.NotFound, "%s %d not found", getType, id)
	case err != nil:
		return status.Errorf(codes.Internal, "Error %s getting %d: %v", getType, id, err)
	}
	return nil
}

func (s *tenantsServicer) validateTenantExists(tenantID int64) error {
	_, err := s.store.GetTenant(tenantID)
	if err != nil {
		return mapErrForGet(err, tenantID, "tenant")
	}
	return nil
}
