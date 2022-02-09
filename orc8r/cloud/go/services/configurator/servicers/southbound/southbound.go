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

	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/services/configurator/storage"
	orc8r_storage "magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
	cfgExt_protos "magma/orc8r/lib/go/protos"
)

type sbExternalConfiguratorServicer struct {
	factory storage.ConfiguratorStorageFactory
}

func NewSouthboundExternalConfiguratorServicer(factory storage.ConfiguratorStorageFactory) (cfgExt_protos.SouthboundExternalConfiguratorServer, error) {
	if factory == nil {
		return nil, fmt.Errorf("storage factory is nil")
	}
	return &sbExternalConfiguratorServicer{factory}, nil
}

func (srv *sbExternalConfiguratorServicer) GetMconfig(ctx context.Context, void *protos.Void) (*protos.GatewayConfigs, error) {
	gw := protos.GetClientGateway(ctx)
	if gw == nil {
		return nil, status.Errorf(codes.PermissionDenied, "missing gateway identity")
	}
	if !gw.Registered() {
		return nil, status.Errorf(codes.PermissionDenied, "gateway not registered")
	}
	return srv.getExternalMconfigImpl(gw.NetworkId, gw.LogicalId)
}

func (srv *sbExternalConfiguratorServicer) getExternalMconfigImpl(networkID string, gatewayID string) (*protos.GatewayConfigs, error) {
	store, err := srv.factory.StartTransaction(context.Background(), &orc8r_storage.TxOptions{ReadOnly: true})
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Aborted, "failed to start transaction: %s", err)
	}

	graph, err := store.LoadGraphForEntity(
		networkID,
		storage.EntityID{Type: orc8r.MagmadGatewayType, Key: gatewayID},
		storage.FullEntityLoadCriteria,
	)
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Internal, "failed to load entity graph: %s", err)
	}

	nwLoad, err := store.LoadNetworks(storage.NetworkLoadFilter{Ids: []string{networkID}}, storage.FullNetworkLoadCriteria)
	if err != nil {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Internal, "failed to load network: %s", err)
	}
	if !funk.IsEmpty(nwLoad.NetworkIDsNotFound) || funk.IsEmpty(nwLoad.Networks) {
		storage.RollbackLogOnError(store)
		return nil, status.Errorf(codes.Internal, "network %s not found: %s", networkID, err)
	}

	// Error on commit is fine for a readonly tx
	storage.CommitLogOnError(store)

	ret, err := mconfig.CreateMconfigJSON(nwLoad.Networks[0], &graph, gatewayID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to build mconfig: %s", err)
	}
	return ret, nil
}
