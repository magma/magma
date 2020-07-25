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

	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	srvRegistry "magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// getTenantsClient is a utility function to get a RPC connection to the
// tenants service
func getTenantsClient() (protos.TenantsServiceClient, error) {
	conn, err := srvRegistry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewTenantsServiceClient(conn), nil
}

func GetAllTenants() (*protos.TenantList, error) {
	oc, err := getTenantsClient()
	if err != nil {
		return nil, err
	}
	tenants, err := oc.GetAllTenants(context.Background(), &protos.Void{})
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

func CreateTenant(tenantID int64, tenant *protos.Tenant) (*protos.Tenant, error) {
	oc, err := getTenantsClient()
	if err != nil {
		return nil, err
	}
	_, err = oc.CreateTenant(context.Background(), &protos.IDAndTenant{
		Id:     tenantID,
		Tenant: tenant,
	})
	if err != nil {
		return nil, err
	}
	return tenant, err
}

func GetTenant(tenantID int64) (*protos.Tenant, error) {
	oc, err := getTenantsClient()
	if err != nil {
		return nil, err
	}

	tenant, err := oc.GetTenant(context.Background(), &protos.GetTenantRequest{Id: tenantID})
	if err != nil {
		switch {
		case status.Convert(err).Code() == codes.NotFound:
			return tenant, merrors.ErrNotFound
		default:
			return tenant, err
		}
	}
	return tenant, nil
}

func SetTenant(tenantID int64, tenant protos.Tenant) error {
	oc, err := getTenantsClient()
	if err != nil {
		return err
	}

	_, err = oc.SetTenant(context.Background(), &protos.IDAndTenant{
		Id:     tenantID,
		Tenant: &tenant,
	})
	if err != nil {
		switch {
		case status.Convert(err).Code() == codes.NotFound:
			return merrors.ErrNotFound
		default:
			return err
		}
	}
	return err
}

func DeleteTenant(tenantID int64) error {
	oc, err := getTenantsClient()
	if err != nil {
		return err
	}

	_, err = oc.DeleteTenant(context.Background(), &protos.GetTenantRequest{Id: tenantID})
	if err != nil {
		switch {
		case status.Convert(err).Code() == codes.NotFound:
			return merrors.ErrNotFound
		default:
			return err
		}
	}
	return err
}
