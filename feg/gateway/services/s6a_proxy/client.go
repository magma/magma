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

// Package s6a_proxy provides a thin client for using s6a proxy service.
// This can be used by apps to discover and contact the service, without knowing about
// the RPC implementation.
package s6a_proxy

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/orc8r/lib/go/util"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

type s6aProxyClient struct {
	protos.S6AProxyClient
	protos.ServiceHealthClient
}

// getS6aProxyClient is a utility function to get a RPC connection to the
// S6a Proxy service
func getS6aProxyClient() (*s6aProxyClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	if util.GetEnvBool("USE_REMOTE_S6A_PROXY") {
		conn, err = registry.Get().GetSharedCloudConnection(strings.ToLower(registry.S6A_PROXY))
	} else {
		conn, err = registry.GetConnection(registry.S6A_PROXY)
	}
	if err != nil {
		errMsg := fmt.Sprintf("S6a Proxy client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &s6aProxyClient{
		protos.NewS6AProxyClient(conn),
		protos.NewServiceHealthClient(conn),
	}, err
}

// AuthenticationInformation sends AIR over diameter connection,
// waits (blocks) for AIA & returns its RPC representation
func AuthenticationInformation(req *protos.AuthenticationInformationRequest) (*protos.AuthenticationInformationAnswer, error) {
	if req == nil {
		return nil, errors.New("Invalid AuthenticationInformationRequest")
	}
	cli, err := getS6aProxyClient()
	if err != nil {
		return nil, err
	}
	return cli.AuthenticationInformation(context.Background(), req)
}

// UpdateLocation sends ULR (Code 316) over diameter connection,
// waits (blocks) for ULA & returns its RPC representation
func UpdateLocation(req *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error) {
	if req == nil {
		return nil, errors.New("Invalid UpdateLocation")
	}
	cli, err := getS6aProxyClient()
	if err != nil {
		return nil, err
	}
	return cli.UpdateLocation(context.Background(), req)
}

// PurgeUE sends PUR (Code 321) over diameter connection,
// waits (blocks) for PUA & returns its RPC representation
func PurgeUE(req *protos.PurgeUERequest) (*protos.PurgeUEAnswer, error) {
	if req == nil {
		return nil, errors.New("Invalid PurgeUE Request")
	}
	cli, err := getS6aProxyClient()
	if err != nil {
		return nil, err
	}
	return cli.PurgeUE(context.Background(), req)
}
