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

	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetToken(ctx context.Context, networkID string, logicalID string, refresh bool) (string, error) {
	client, err := getCloudRegistrationClient()
	if err != nil {
		return "", err
	}

	req := &protos.GetTokenRequest{
		GatewayDeviceInfo: &protos.GatewayDeviceInfo{
			NetworkId:            networkID,
			LogicalId:            logicalID,
		},
		Refresh:                refresh,
	}

	res, err := client.GetToken(ctx, req)
	return res.Token, err
}

func GetGatewayRegistrationInfo(ctx context.Context, token string) (*protos.GetGatewayRegistrationInfoResponse, error) {
	client, err := getCloudRegistrationClient()
	if err != nil {
		return nil, err
	}

	req := &protos.GetGatewayRegistrationInfoRequest{}

	res, err := client.GetGatewayRegistrationInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetGatewayDeviceInfo(ctx context.Context, token string) (*protos.GatewayDeviceInfo, error) {
	client, err := getCloudRegistrationClient()
	if err != nil {
		return nil, err
	}

	req := &protos.GetGatewayDeviceInfoRequest{
		Token:                token,
	}

	res, err := client.GetGatewayDeviceInfo(ctx, req)
	if err != nil {
		return nil, err
	}

	clientErr := res.Response.(*protos.GetGatewayDeviceInfoResponse_Error)
	if clientErr != nil{
		// TODO(reginawang3495): Be more precise based on the different errors?
		return nil, status.Error(codes.Unauthenticated, clientErr.Error)
	}

	gatewayInfo := res.Response.(*protos.GetGatewayDeviceInfoResponse_GatewayDeviceInfo)
	return gatewayInfo.GatewayDeviceInfo, nil
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
