/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tenants

import (
	"context"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tenant_protos "magma/orc8r/cloud/go/services/tenants/protos"
	"magma/orc8r/lib/go/merrors"
	"magma/orc8r/lib/go/protos"
	srvRegistry "magma/orc8r/lib/go/registry"
)

// getTenantsClient is a utility function to get a RPC connection to the
// tenants service
func getTenantsClient() (tenant_protos.TenantsServiceClient, error) {
	conn, err := srvRegistry.GetConnection(ServiceName, protos.ServiceType_PROTECTED)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return tenant_protos.NewTenantsServiceClient(conn), nil
}

func GetAllTenants(ctx context.Context) (*tenant_protos.TenantList, error) {
	oc, err := getTenantsClient()
	if err != nil {
		return nil, err
	}

	tenants, err := oc.GetAllTenants(ctx, &protos.Void{})
	if err != nil {
		return nil, err
	}

	return tenants, nil
}

func CreateTenant(ctx context.Context, tenantID int64, tenant *tenant_protos.Tenant) (*tenant_protos.Tenant, error) {
	oc, err := getTenantsClient()
	if err != nil {
		return nil, err
	}

	_, err = oc.CreateTenant(
		ctx,
		&tenant_protos.IDAndTenant{
			Id:     tenantID,
			Tenant: tenant,
		},
	)
	if err != nil {
		return nil, err
	}

	return tenant, err
}

func GetTenant(ctx context.Context, tenantID int64) (*tenant_protos.Tenant, error) {
	oc, err := getTenantsClient()
	if err != nil {
		return nil, err
	}

	tenant, err := oc.GetTenant(ctx, &tenant_protos.GetTenantRequest{Id: tenantID})
	if err != nil {
		return nil, mapErr(err)
	}
	return tenant, nil
}

func SetTenant(ctx context.Context, tenantID int64, tenant *tenant_protos.Tenant) error {
	oc, err := getTenantsClient()
	if err != nil {
		return err
	}

	_, err = oc.SetTenant(
		ctx,
		&tenant_protos.IDAndTenant{
			Id:     tenantID,
			Tenant: tenant,
		},
	)
	if err != nil {
		return mapErr(err)
	}

	return nil
}

func DeleteTenant(ctx context.Context, tenantID int64) error {
	oc, err := getTenantsClient()
	if err != nil {
		return err
	}

	_, err = oc.DeleteTenant(ctx, &tenant_protos.GetTenantRequest{Id: tenantID})
	if err != nil {
		return mapErr(err)
	}

	return nil
}

func GetControlProxy(ctx context.Context, tenantID int64) (*tenant_protos.GetControlProxyResponse, error) {
	oc, err := getTenantsClient()
	if err != nil {
		return nil, err
	}

	controlProxy, err := oc.GetControlProxy(ctx, &tenant_protos.GetControlProxyRequest{Id: tenantID})
	if err != nil {
		return nil, mapErr(err)
	}

	return controlProxy, nil
}

func GetControlProxyFromNetworkID(ctx context.Context, networkID string) (*tenant_protos.GetControlProxyResponse, error) {
	oc, err := getTenantsClient()
	if err != nil {
		return nil, err
	}

	controlProxy, err := oc.GetControlProxyFromNetworkID(ctx, &tenant_protos.GetControlProxyFromNetworkIDRequest{NetworkID: networkID})
	if err != nil {
		return nil, mapErr(err)
	}

	return controlProxy, nil
}

func CreateOrUpdateControlProxy(ctx context.Context, controlProxy *tenant_protos.CreateOrUpdateControlProxyRequest) error {
	oc, err := getTenantsClient()
	if err != nil {
		return err
	}

	_, err = oc.CreateOrUpdateControlProxy(ctx, controlProxy)
	if err != nil {
		return mapErr(err)
	}

	return nil
}

func mapErr(err error) error {
	switch {
	case status.Convert(err).Code() == codes.NotFound:
		return merrors.ErrNotFound
	default:
		return err
	}
}
