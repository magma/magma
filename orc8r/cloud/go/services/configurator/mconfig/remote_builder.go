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

package mconfig

import (
	"context"
	"strings"

	"magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	"magma/orc8r/cloud/go/services/configurator/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

// remoteBuilder identifies a remote mconfig builder.
type remoteBuilder struct {
	// service name of the builder
	// should always be lowercase to match service registry convention
	service string
}

func NewRemoteBuilder(serviceName string) Builder {
	return &remoteBuilder{service: strings.ToLower(serviceName)}
}

func (r *remoteBuilder) Build(network *storage.Network, graph *storage.EntityGraph, gatewayID string) (ConfigsByKey, error) {
	c, err := r.getBuilderClient()
	if err != nil {
		return nil, err
	}

	res, err := c.Build(context.Background(), &protos.BuildRequest{Network: network, Graph: graph, GatewayId: gatewayID})
	if err != nil {
		return nil, err
	}

	return res.ConfigsByKey, nil
}

func (r *remoteBuilder) getBuilderClient() (protos.MconfigBuilderClient, error) {
	conn, err := registry.GetConnection(r.service)
	if err != nil {
		initErr := merrors.NewInitError(err, r.service)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewMconfigBuilderClient(conn), nil
}
