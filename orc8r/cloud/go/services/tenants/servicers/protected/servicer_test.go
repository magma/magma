package servicers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/tenants"
	tenant_protos "magma/orc8r/cloud/go/services/tenants/protos"
	servicers "magma/orc8r/cloud/go/services/tenants/servicers/protected"
	"magma/orc8r/cloud/go/services/tenants/servicers/storage"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

var (
	sampleTenant = tenant_protos.Tenant{
		Name:     "test",
		Networks: []string{"network_1", "network_2"},
	}
	sampleTenant2 = tenant_protos.Tenant{
		Name:     "test2",
		Networks: []string{"network_1"},
	}
	sampleTenantID      int64 = 3
	sampleControlProxy        = "{ otherInfo }"
	sampleControlProxy2       = "{ otherInfo2 }"

	sampleCreateControlProxyReq = tenant_protos.CreateOrUpdateControlProxyRequest{
		Id:           sampleTenantID,
		ControlProxy: sampleControlProxy,
	}
	sampleCreateControlProxyReq2 = tenant_protos.CreateOrUpdateControlProxyRequest{
		Id:           sampleTenantID,
		ControlProxy: sampleControlProxy2,
	}
	sampleGetControlProxyRes = tenant_protos.GetControlProxyResponse{
		Id:           sampleTenantID,
		ControlProxy: sampleControlProxy,
	}
	sampleGetControlProxyRes2 = tenant_protos.GetControlProxyResponse{
		Id:           sampleTenantID,
		ControlProxy: sampleControlProxy2,
	}
)

func TestTenantsServicer(t *testing.T) {
	srv, err := newTestService(t)
	assert.NoError(t, err)

	// Empty db, no tenant found
	_, err = srv.GetTenant(context.Background(), &tenant_protos.GetTenantRequest{Id: 1})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, "tenant 1 not found", status.Convert(err).Message())

	// Create "test" tenant
	createResp, err := srv.CreateTenant(context.Background(), &tenant_protos.IDAndTenant{
		Id:     1,
		Tenant: &sampleTenant,
	})
	assert.NoError(t, err)
	assert.Equal(t, &protos.Void{}, createResp)

	// Get "test" tenant
	getResp, err := srv.GetTenant(context.Background(), &tenant_protos.GetTenantRequest{Id: 1})
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, &sampleTenant, getResp)

	// Get "other" tenant
	_, err = srv.GetTenant(context.Background(), &tenant_protos.GetTenantRequest{Id: 2})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, "tenant 2 not found", status.Convert(err).Message())

	// Update "test" tenant
	setResp, err := srv.SetTenant(context.Background(), &tenant_protos.IDAndTenant{
		Id:     1,
		Tenant: &sampleTenant2,
	})
	assert.NoError(t, err)
	assert.Equal(t, &protos.Void{}, setResp)
	// get updated tenant
	getResp, err = srv.GetTenant(context.Background(), &tenant_protos.GetTenantRequest{Id: 1})
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, &sampleTenant2, getResp)

	// Update nonexistent tenant
	_, err = srv.SetTenant(context.Background(), &tenant_protos.IDAndTenant{
		Id:     3,
		Tenant: &sampleTenant2,
	})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, "tenant 3 not found", status.Convert(err).Message())

	// Create second tenant
	_, err = srv.CreateTenant(context.Background(), &tenant_protos.IDAndTenant{
		Id:     2,
		Tenant: &sampleTenant,
	})
	assert.NoError(t, err)

	// Get all tenants
	getAllResp, err := srv.GetAllTenants(context.Background(), &protos.Void{})
	assert.NoError(t, err)
	assert.Len(t, getAllResp.Tenants, 2)

	// Delete "other" tenant
	delResp, err := srv.DeleteTenant(context.Background(), &tenant_protos.GetTenantRequest{Id: 2})
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, &protos.Void{}, delResp)

	_, err = srv.GetTenant(context.Background(), &tenant_protos.GetTenantRequest{Id: 2})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, "tenant 2 not found", status.Convert(err).Message())
}

func TestControlProxyTenantsServicer(t *testing.T) {
	srv, err := newTestService(t)
	assert.NoError(t, err)

	// Get control_proxy not set
	_, err = srv.GetControlProxy(context.Background(), &tenant_protos.GetControlProxyRequest{Id: sampleTenantID})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, fmt.Sprintf("tenant %d not found", sampleTenantID), status.Convert(err).Message())

	// Create control_proxy when tenant not created yet
	_, err = srv.CreateOrUpdateControlProxy(context.Background(), &sampleCreateControlProxyReq)
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, fmt.Sprintf("tenant %d not found", sampleTenantID), status.Convert(err).Message())

	// Create "test" tenant
	createResp, err := srv.CreateTenant(context.Background(), &tenant_protos.IDAndTenant{
		Id:     sampleTenantID,
		Tenant: &sampleTenant,
	})
	assert.NoError(t, err)
	assert.Equal(t, &protos.Void{}, createResp)

	// Get control_proxy not set
	_, err = srv.GetControlProxy(context.Background(), &tenant_protos.GetControlProxyRequest{Id: sampleTenantID})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, fmt.Sprintf("controlProxy %d not found", sampleTenantID), status.Convert(err).Message())

	// Get control_proxy not set
	_, err = srv.GetControlProxyFromNetworkID(context.Background(), &tenant_protos.GetControlProxyFromNetworkIDRequest{NetworkID: sampleTenant.Networks[0]})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, fmt.Sprintf("no control-proxy found for tenant %d", sampleTenantID), status.Convert(err).Message())

	// Create control_proxy
	_, err = srv.CreateOrUpdateControlProxy(context.Background(), &sampleCreateControlProxyReq)
	assert.NoError(t, err)
	// get updated control_proxy
	controlProxy, err := srv.GetControlProxy(context.Background(), &tenant_protos.GetControlProxyRequest{Id: sampleTenantID})
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, controlProxy, &sampleGetControlProxyRes)

	// Update control_proxy
	_, err = srv.CreateOrUpdateControlProxy(context.Background(), &sampleCreateControlProxyReq2)
	assert.NoError(t, err)
	// get updated control_proxy
	controlProxy, err = srv.GetControlProxy(context.Background(), &tenant_protos.GetControlProxyRequest{Id: sampleTenantID})
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, controlProxy, &sampleGetControlProxyRes2)

	// get control_proxy from network ID
	controlProxy, err = srv.GetControlProxyFromNetworkID(context.Background(), &tenant_protos.GetControlProxyFromNetworkIDRequest{NetworkID: sampleTenant.Networks[0]})
	assert.NoError(t, err)
	test_utils.AssertMessagesEqual(t, controlProxy, &sampleGetControlProxyRes2)

	// get control_proxy from network ID, no tenant
	_, err = srv.GetControlProxyFromNetworkID(context.Background(), &tenant_protos.GetControlProxyFromNetworkIDRequest{NetworkID: "network_nonexistent"})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, "tenantID for current NetworkID network_nonexistent not found", status.Convert(err).Message())

	// Get control_proxy not set
	_, err = srv.GetControlProxy(context.Background(), &tenant_protos.GetControlProxyRequest{Id: sampleTenantID + 1})
	assert.Equal(t, codes.NotFound, status.Convert(err).Code())
	assert.Equal(t, fmt.Sprintf("tenant %d not found", sampleTenantID+1), status.Convert(err).Message())
}

func newTestService(t *testing.T) (tenant_protos.TenantsServiceServer, error) {
	srv, lis, _ := test_utils.NewTestService(t, orc8r.ModuleName, tenants.ServiceName)
	factory := test_utils.NewSQLBlobstore(t, "tenants_servicer_test_blobstore")
	store := storage.NewBlobstoreStore(factory)
	servicer, err := servicers.NewTenantsServicer(store)
	assert.NoError(t, err)
	tenant_protos.RegisterTenantsServiceServer(srv.ProtectedGrpcServer, servicer)
	go srv.RunTest(lis, nil)
	return servicer, nil
}
