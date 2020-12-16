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

package lte

import (
	"context"
	"fmt"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/lte/protos"
	"magma/orc8r/cloud/go/serde"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

func GetEnodebState(networkID string, gatewayID string, enodebSN string) (*lte_models.EnodebState, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	res, err := client.GetEnodebState(
		context.Background(),
		&protos.GetEnodebStateRequest{
			NetworkId: networkID,
			GatewayId: gatewayID,
			EnodebSn:  enodebSN,
		},
	)
	if err != nil {
		return nil, err
	}
	enodebStateI, err := serde.Deserialize(res.SerializedState, lte.EnodebStateType, serdes.State)
	if err != nil {
		return nil, fmt.Errorf("error deserializing enodeb state for enodeb %s", enodebSN)
	}
	enodebState, ok := enodebStateI.(*lte_models.EnodebState)
	if !ok {
		return nil, fmt.Errorf("error converting returned state to EnodebState model for enodeb %s", enodebSN)
	}
	return enodebState, nil
}

func SetEnodebState(networkID string, gatewayID string, enodebSN string, serializedState []byte) error {
	client, err := getClient()
	if err != nil {
		return err
	}
	_, err = client.SetEnodebState(
		context.Background(),
		&protos.SetEnodebStateRequest{
			NetworkId:       networkID,
			GatewayId:       gatewayID,
			EnodebSn:        enodebSN,
			SerializedState: serializedState,
		},
	)
	return err
}

func getClient() (protos.EnodebStateLookupClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewEnodebStateLookupClient(conn), nil
}
