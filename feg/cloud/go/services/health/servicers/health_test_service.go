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
	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health/storage"
	"magma/feg/cloud/go/services/health/test_utils"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/lib/go/protos"

	"golang.org/x/net/context"
)

// A little Go "polymorphism" magic for testing
type TestHealthServer struct {
	HealthServer
	Feg1 bool //boolean to simulate requests coming from more than 1 FeG
}

// Health receiver for testHealthServer injects GW Identity into CTX if it's
// missing for testing without heavy mock of Certifier & certificate addition
func (srv *TestHealthServer) UpdateHealth(
	ctx context.Context,
	req *fegprotos.HealthRequest,
) (*fegprotos.HealthResponse, error) {

	gw := protos.GetClientGateway(ctx)
	if gw == nil {
		if srv.Feg1 {
			ctx = protos.NewGatewayIdentity(test_utils.TestFegHwId1, test_utils.TestFegNetwork, test_utils.TestFegLogicalId1).NewContextWithIdentity(ctx)
		} else {
			ctx = protos.NewGatewayIdentity(test_utils.TestFegHwId2, test_utils.TestFegNetwork, test_utils.TestFegLogicalId2).NewContextWithIdentity(ctx)
		}
	}
	return srv.HealthServer.UpdateHealth(ctx, req)
}

func NewTestHealthServer(mockFactory blobstore.BlobStorageFactory) (*TestHealthServer, error) {
	store, err := storage.NewHealthBlobstore(mockFactory)
	if err != nil {
		return nil, err
	}
	return &TestHealthServer{
		HealthServer: HealthServer{
			store: store,
		},
		Feg1: true,
	}, nil
}
