package servicers_test

import (
	"context"
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/tenants"
	"magma/orc8r/cloud/go/services/tenants/servicers"
	"magma/orc8r/cloud/go/services/tenants/servicers/storage"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	sampleTenant = protos.Tenant{
		Name:     "test",
		Networks: []string{"network_1", "network_2"},
	}
	sampleTenant2 = protos.Tenant{
		Name:     "test2",
		Networks: []string{"network_1"},
	}
)

func TestTenantsServicer(t *testing.T) {
	srv, err := newTestService(t)
	assert.NoError(t, err)

	// Empty db, no tenant found
	_, err = srv.GetTenant(context.Background(), &protos.GetTenantRequest{Id: 1})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, "Tenant 1 not found", status.Convert(err).Message())

	// Create "test" tenant
	createResp, err := srv.CreateTenant(context.Background(), &protos.IDAndTenant{
		Id:     1,
		Tenant: &sampleTenant,
	})
	assert.NoError(t, err)
	assert.Equal(t, &protos.Void{}, createResp)

	// Get "test" tenant
	getResp, err := srv.GetTenant(context.Background(), &protos.GetTenantRequest{Id: 1})
	assert.NoError(t, err)
	assert.Equal(t, &sampleTenant, getResp)

	// Get "other" tenant
	_, err = srv.GetTenant(context.Background(), &protos.GetTenantRequest{Id: 2})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, "Tenant 2 not found", status.Convert(err).Message())

	// Update "test" tenant
	setResp, err := srv.SetTenant(context.Background(), &protos.IDAndTenant{
		Id:     1,
		Tenant: &sampleTenant2,
	})
	assert.NoError(t, err)
	assert.Equal(t, &protos.Void{}, setResp)
	// get updated tenant
	getResp, err = srv.GetTenant(context.Background(), &protos.GetTenantRequest{Id: 1})
	assert.NoError(t, err)
	assert.Equal(t, sampleTenant2, *getResp)

	// Update nonexistent tenant
	_, err = srv.SetTenant(context.Background(), &protos.IDAndTenant{
		Id:     3,
		Tenant: &sampleTenant2,
	})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, "Tenant 3 not found", status.Convert(err).Message())

	// Create second tenant
	_, err = srv.CreateTenant(context.Background(), &protos.IDAndTenant{
		Id:     2,
		Tenant: &sampleTenant,
	})
	assert.NoError(t, err)

	// Get all tenants
	getAllResp, err := srv.GetAllTenants(context.Background(), &protos.Void{})
	assert.NoError(t, err)
	assert.Len(t, getAllResp.Tenants, 2)

	// Delete "other" tenant
	delResp, err := srv.DeleteTenant(context.Background(), &protos.GetTenantRequest{Id: 2})
	assert.NoError(t, err)
	assert.Equal(t, protos.Void{}, *delResp)

	_, err = srv.GetTenant(context.Background(), &protos.GetTenantRequest{Id: 2})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, "Tenant 2 not found", status.Convert(err).Message())
}

func newTestService(t *testing.T) (protos.TenantsServiceServer, error) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, tenants.ServiceName)
	factory := test_utils.NewSQLBlobstore(t, "tenants_servicer_test_blobstore")
	store := storage.NewBlobstoreStore(factory)
	servicer, err := servicers.NewTenantsServicer(store)
	assert.NoError(t, err)
	protos.RegisterTenantsServiceServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
	return servicer, nil
}
