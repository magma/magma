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

	fegprotos "magma/feg/cloud/go/protos"
	"magma/orc8r/lib/go/protos"
)

// LocationUpdateAcc relays the LocationUpdateAccept sent from VLR->FeG->Access Gateway
func (srv *FegToGwRelayServer) LocationUpdateAcc(
	ctx context.Context,
	req *fegprotos.LocationUpdateAccept,
) (*protos.Void, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, err
	}
	return srv.LocationUpdateAcceptUnverified(ctx, req)
}

func (srv *FegToGwRelayServer) LocationUpdateAcceptUnverified(
	ctx context.Context,
	req *fegprotos.LocationUpdateAccept,
) (*protos.Void, error) {
	conn, ctx, err := getGWSGSServiceConnCtx(ctx, req.Imsi)
	if err != nil {
		return &protos.Void{}, err
	}
	client := fegprotos.NewCSFBGatewayServiceClient(conn)
	return client.LocationUpdateAcc(ctx, req)
}

// LocationUpdateRej relays the LocationUpdateReject sent from VLR->FeG->Access Gateway
func (srv *FegToGwRelayServer) LocationUpdateRej(
	ctx context.Context,
	req *fegprotos.LocationUpdateReject,
) (*protos.Void, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, err
	}
	return srv.LocationUpdateRejectUnverified(ctx, req)
}

func (srv *FegToGwRelayServer) LocationUpdateRejectUnverified(
	ctx context.Context,
	req *fegprotos.LocationUpdateReject,
) (*protos.Void, error) {
	conn, ctx, err := getGWSGSServiceConnCtx(ctx, req.Imsi)
	if err != nil {
		return &protos.Void{}, err
	}
	client := fegprotos.NewCSFBGatewayServiceClient(conn)
	return client.LocationUpdateRej(ctx, req)
}
