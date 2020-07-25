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

package servicers

import (
	"context"
	"fmt"

	"magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"magma/orc8r/cloud/go/services/tenants/servicers/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	tenants, err := s.store.GetAllTenants()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error getting tenants: %v", err)
	}
	return tenants, nil
}

func (s *tenantsServicer) CreateTenant(c context.Context, request *protos.IDAndTenant) (*protos.Void, error) {
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

func (s *tenantsServicer) GetTenant(c context.Context, request *protos.GetTenantRequest) (*protos.Tenant, error) {
	tenant, err := s.store.GetTenant(request.Id)
	switch {
	case err == errors.ErrNotFound:
		return nil, status.Errorf(codes.NotFound, "Tenant %d not found", request.Id)
	case err != nil:
		return nil, status.Errorf(codes.Internal, "Error getting tenant %d: %v", request.Id, err)
	}
	return tenant, nil

}

func (s *tenantsServicer) SetTenant(c context.Context, request *protos.IDAndTenant) (*protos.Void, error) {
	_, err := s.store.GetTenant(request.Id)
	switch {
	case err == errors.ErrNotFound:
		return nil, status.Errorf(codes.NotFound, "Tenant %d not found", request.Id)
	case err != nil:
		return nil, status.Errorf(codes.Internal, "Error getting tenant %d: %v", request.Id, err)
	}
	err = s.store.SetTenant(request.Id, *request.Tenant)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error setting tenant %d: %v", request.Id, err)
	}
	return &protos.Void{}, nil
}

func (s *tenantsServicer) DeleteTenant(c context.Context, request *protos.GetTenantRequest) (*protos.Void, error) {
	err := s.store.DeleteTenant(request.Id)
	switch {
	case err == errors.ErrNotFound:
		return nil, status.Errorf(codes.NotFound, "Tenant %d not found", request.Id)
	case err != nil:
		return nil, status.Errorf(codes.Internal, "Error deleting tenant %d: %v", request.Id, err)
	}
	return &protos.Void{}, nil
}
