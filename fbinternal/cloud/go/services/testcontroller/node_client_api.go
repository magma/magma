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

package testcontroller

import (
	"context"

	"magma/fbinternal/cloud/go/services/testcontroller/protos"
	"magma/fbinternal/cloud/go/services/testcontroller/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func getNodeClient() (protos.NodeLeasorClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewNodeLeasorClient(conn), nil
}

func GetNodes(ctx context.Context, ids []string, tag *string) (map[string]*storage.CINode, error) {
	client, err := getNodeClient()
	if err != nil {
		return nil, err
	}
	res, err := client.GetNodes(ctx, &protos.GetNodesRequest{Ids: ids, Tag: asStringValue(tag)})
	if err != nil {
		return nil, err
	}
	return res.Nodes, nil
}

func CreateOrUpdateNode(ctx context.Context, node *storage.MutableCINode) error {
	client, err := getNodeClient()
	if err != nil {
		return err
	}
	_, err = client.CreateOrUpdateNode(ctx, &protos.CreateOrUpdateNodeRequest{Node: node})
	return err
}

func DeleteNode(ctx context.Context, id string) error {
	client, err := getNodeClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteNode(ctx, &protos.DeleteNodeRequest{Id: id})
	return err
}

func ReserveNode(ctx context.Context, id string) (*storage.NodeLease, error) {
	client, err := getNodeClient()
	if err != nil {
		return nil, err
	}
	res, err := client.ReserveNode(ctx, &protos.ReserveNodeRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return res.Lease, nil
}

func LeaseNode(ctx context.Context, tag string) (*storage.NodeLease, error) {
	client, err := getNodeClient()
	if err != nil {
		return nil, err
	}
	res, err := client.LeaseNode(ctx, &protos.LeaseNodeRequest{Tag: tag})
	if err != nil {
		return nil, err
	}
	return res.Lease, nil
}

func ReleaseNode(ctx context.Context, id string, leaseID string) error {
	client, err := getNodeClient()
	if err != nil {
		return err
	}
	_, err = client.ReleaseNode(ctx, &protos.ReleaseNodeRequest{NodeID: id, LeaseID: leaseID})
	return err
}

func asStringValue(s *string) *wrappers.StringValue {
	if s == nil {
		return nil
	}
	return &wrappers.StringValue{Value: *s}
}
