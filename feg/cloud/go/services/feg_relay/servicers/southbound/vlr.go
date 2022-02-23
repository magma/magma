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

	"github.com/golang/glog"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/orc8r/lib/go/protos"
)

// VLRResetAck relays the ResetAck sent from VLR->FeG->All Access Gateway
func (srv *FegToGwRelayServer) VLRResetAck(
	ctx context.Context,
	req *fegprotos.ResetAck,
) (*protos.Void, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, err
	}
	return srv.VLRResetAckUnverified(ctx, req)
}

func (srv *FegToGwRelayServer) VLRResetAckUnverified(
	ctx context.Context,
	req *fegprotos.ResetAck,
) (*protos.Void, error) {
	connList, ctxList, err := getAllGWSGSServiceConnCtx(ctx)
	if err != nil {
		return &protos.Void{}, err
	}

	// forward the message to all gateways
	for idx, conn := range connList {
		client := fegprotos.NewCSFBGatewayServiceClient(conn)
		_, err = client.VLRResetAck(ctxList[idx], req)
		if err != nil {
			glog.V(2).Infof("Failed to send SGsAP-RESET-ACK to a gateway, continuing")
		}
	}

	return &protos.Void{}, nil
}

// VLRResetIndication relays the ResetIndication sent from VLR->FeG->All Access Gateway
func (srv *FegToGwRelayServer) VLRResetIndication(
	ctx context.Context,
	req *fegprotos.ResetIndication,
) (*protos.Void, error) {
	if err := validateFegContext(ctx); err != nil {
		glog.Errorf("Failed to validate FeG context: %s", err)
		return nil, err
	}
	return srv.VLRResetIndicationUnverified(ctx, req)
}

func (srv *FegToGwRelayServer) VLRResetIndicationUnverified(
	ctx context.Context,
	req *fegprotos.ResetIndication,
) (*protos.Void, error) {
	connList, ctxList, err := getAllGWSGSServiceConnCtx(ctx)
	if err != nil {
		glog.Errorf("Failed to getAllGWSGSServiceConnCtx(ctx): %s", err)
		return &protos.Void{}, err
	}

	// forward the message to all gateways
	for idx, conn := range connList {
		client := fegprotos.NewCSFBGatewayServiceClient(conn)
		_, err = client.VLRResetIndication(ctxList[idx], req)
		if err != nil {
			glog.V(2).Infof("Failed to send SGsAP-RESET-INDICATION to a gateway, continuing")
		}
	}

	return &protos.Void{}, nil
}

// VLRStatus relays the Status sent from VLR->FeG->All Access Gateway
func (srv *FegToGwRelayServer) VLRStatus(
	ctx context.Context,
	req *fegprotos.Status,
) (*protos.Void, error) {
	if err := validateFegContext(ctx); err != nil {
		return nil, err
	}
	return srv.VLRStatusUnverified(ctx, req)
}

func (srv *FegToGwRelayServer) VLRStatusUnverified(
	ctx context.Context,
	req *fegprotos.Status,
) (*protos.Void, error) {
	connList, ctxList, err := getAllGWSGSServiceConnCtx(ctx)
	if err != nil {
		return &protos.Void{}, err
	}

	// forward the message to all gateways
	for idx, conn := range connList {
		client := fegprotos.NewCSFBGatewayServiceClient(conn)
		_, err = client.VLRStatus(ctxList[idx], req)
		if err != nil {
			glog.V(2).Infof("Failed to send SGsAP-STATUS to a gateway, continuing")
		}
	}

	return &protos.Void{}, nil
}
