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

package registration

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/lib/go/protos"
)

// NotImplementedWarning is pulled out of getDomainName for ease of testing
const NotImplementedWarning = "warning: not implemented"

type cloudRegistrationServicer struct {
	store   Store
	rootCA  string
	timeout time.Duration
}

func NewCloudRegistrationServicer(store Store, rootCA string, timeout time.Duration) (protos.CloudRegistrationServer, error) {
	if store == nil {
		return nil, fmt.Errorf("storage store is nil")
	}
	return &cloudRegistrationServicer{store: store, rootCA: rootCA, timeout: timeout}, nil
}

func (c *cloudRegistrationServicer) GetToken(ctx context.Context, request *protos.GetTokenRequest) (*protos.GetTokenResponse, error) {
	networkId := request.GatewayDeviceInfo.NetworkId
	logicalId := request.GatewayDeviceInfo.LogicalId

	tokenInfo, err := c.store.GetTokenInfoFromLogicalID(networkId, logicalId)
	if err != nil {
		// Error is not bubbled up since the tokenInfo is only important if the token exists and is expired
		// If GetTokenInfoFromLogicalID fails, we can just ignore and continue
		glog.V(2).Infof("could not get tokenInfo for networkID %v and logicalID %v: %v", networkId, logicalId, err)
	}

	refresh := request.Refresh || tokenInfo == nil || IsExpired(tokenInfo)
	if refresh {
		tokenInfo, err = c.generateAndSaveTokenInfo(networkId, logicalId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not generate and save tokenInfo for networkID %v and logicalID %v: %v", networkId, logicalId, err)
		}
		glog.V(2).Infof("generated new token for networkID %v and logicalID %v: %v", networkId, logicalId, err)
	}

	res := &protos.GetTokenResponse{Timeout: tokenInfo.Timeout, Token: NonceToToken(tokenInfo.Nonce)}
	return res, nil
}

func (c *cloudRegistrationServicer) GetGatewayRegistrationInfo(ctx context.Context, request *protos.GetGatewayRegistrationInfoRequest) (*protos.GetGatewayRegistrationInfoResponse, error) {
	domainName := getDomainName()
	res := &protos.GetGatewayRegistrationInfoResponse{
		RootCa:     c.rootCA,
		DomainName: domainName,
	}
	return res, nil
}

func (c *cloudRegistrationServicer) GetGatewayDeviceInfo(ctx context.Context, request *protos.GetGatewayDeviceInfoRequest) (*protos.GetGatewayDeviceInfoResponse, error) {
	nonce, err := NonceFromToken(request.Token)
	if err != nil {
		res := &protos.GetGatewayDeviceInfoResponse{
			Response: &protos.GetGatewayDeviceInfoResponse_Error{
				Error: err.Error(),
			},
		}
		return res, nil
	}

	tokenInfo, err := c.store.GetTokenInfoFromNonce(nonce)
	if err != nil {
		res := &protos.GetGatewayDeviceInfoResponse{
			Response: &protos.GetGatewayDeviceInfoResponse_Error{
				Error: fmt.Sprintf("could not get token info from token %v: %v", request.Token, err),
			},
		}
		return res, nil
	}
	if IsExpired(tokenInfo) {
		res := &protos.GetGatewayDeviceInfoResponse{
			Response: &protos.GetGatewayDeviceInfoResponse_Error{
				Error: fmt.Sprintf("token %v has expired", request.Token),
			},
		}
		return res, nil
	}

	res := &protos.GetGatewayDeviceInfoResponse{
		Response: &protos.GetGatewayDeviceInfoResponse_GatewayDeviceInfo{
			GatewayDeviceInfo: &protos.GatewayDeviceInfo{
				NetworkId: tokenInfo.GatewayDeviceInfo.NetworkId,
				LogicalId: tokenInfo.GatewayDeviceInfo.LogicalId,
			},
		},
	}
	return res, nil
}

func (c *cloudRegistrationServicer) generateAndSaveTokenInfo(networkID string, logicalID string) (*protos.TokenInfo, error) {
	nonce := GenerateNonce(NonceLength)
	timeout := clock.Now().Add(c.timeout)

	tokenInfo := &protos.TokenInfo{
		GatewayDeviceInfo: &protos.GatewayDeviceInfo{
			NetworkId: networkID,
			LogicalId: logicalID,
		},
		Nonce:   nonce,
		Timeout: GetTimestamp(timeout),
	}

	err := c.store.SetTokenInfo(tokenInfo)
	if err != nil {
		return nil, err
	}

	return tokenInfo, nil
}

// TODO(#10437)
func getDomainName() string {
	return NotImplementedWarning
}
