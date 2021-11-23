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

package bootstrapper

import (
	"context"

	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"
)

func GetToken(ctx context.Context, networkID string, logicalID string, refresh bool) (string, error) {
	client, err := getCloudRegistrationClient()
	if err != nil {
		return "", err
	}

	req := &protos.GetTokenRequest{
		GatewayPreregisterInfo: &protos.GatewayPreregisterInfo{
			NetworkId:            networkID,
			LogicalId:            logicalID,
		},
		Refresh:                refresh,
	}

	res, err := client.GetToken(ctx, req)
	return res.Token, err
}

func GetGatewayPreregisterInfo(ctx context.Context, token string) (*protos.GatewayPreregisterInfo, error) {
	client, err := getCloudRegistrationClient()
	if err != nil {
		return nil, err
	}

	req := &protos.GetGatewayPreregisterInfoRequest{
		Token:                token,
	}

	res, err := client.GetGatewayPreregisterInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	clientErr := res.Response.(*protos.GetGatewayPreregisterInfoResponse_Error)
	if clientErr != nil{
		// TODO(reginawang3495): Be more precise based on the different errors?
		return nil, status.Error(codes.Unauthenticated, clientErr.Error)
	}

	gatewayInfo := res.Response.(*protos.GetGatewayPreregisterInfoResponse_GatewayPreregisterInfo)
	return gatewayInfo.GatewayPreregisterInfo, nil
}

func GetInfoForGatewayRegistration(ctx context.Context, token string) (*protos.GetInfoForGatewayRegistrationResponse, error) {
	client, err := getCloudRegistrationClient()
	if err != nil {
		return nil, err
	}

	req := &protos.GetInfoForGatewayRegistrationRequest{}

	res, err := client.GetInfoForGatewayRegistration(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getCloudRegistrationClient() (protos.CloudRegistrationClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewCloudRegistrationClient(conn), err
}
