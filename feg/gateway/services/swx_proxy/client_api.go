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

// Package swx_proxy provides a thin client for using swx proxy service.
// This can be used by apps to discover and contact the service, without knowing about
// the RPC implementation.
package swx_proxy

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/orc8r/lib/go/util"
)

// Wrapper for GRPC Client
// functionality
type swxProxyClient struct {
	protos.SwxProxyClient
	cc *grpc.ClientConn
}

// getSwxProxyClient is a utility function to get a RPC connection to the
// Swx Proxy service
func getSwxProxyClient() (*swxProxyClient, error) {
	var conn *grpc.ClientConn
	var err error
	if util.GetEnvBool("USE_REMOTE_SWX_PROXY", true) {
		conn, err = registry.Get().GetSharedCloudConnection(strings.ToLower(registry.SWX_PROXY))
	} else {
		conn, err = registry.GetConnection(registry.SWX_PROXY)
	}
	if err != nil {
		errMsg := fmt.Sprintf("Swx Proxy client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &swxProxyClient{
		protos.NewSwxProxyClient(conn),
		conn,
	}, err
}

// Authenticate sends MAR (code 303) over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func Authenticate(req *protos.AuthenticationRequest) (*protos.AuthenticationAnswer, error) {
	err := verifyAuthenticationRequest(req)
	if err != nil {
		errMsg := fmt.Errorf("Invalid AuthenticationRequest provided: %s", err)
		return nil, errors.New(errMsg.Error())
	}
	cli, err := getSwxProxyClient()
	if err != nil {
		return nil, err
	}
	return cli.Authenticate(context.Background(), req)
}

// Register sends SAR (Code 301) over diameter connection with ServerAssignmentType
// set to REGISTRATION, waits (blocks) for SAA & returns its RPC representation
func Register(req *protos.RegistrationRequest) (*protos.RegistrationAnswer, error) {
	err := verifyRegistrationRequest(req)
	if err != nil {
		errMsg := fmt.Errorf("Invalid RegistrationRequest provided: %s", err)
		return nil, errors.New(errMsg.Error())
	}
	cli, err := getSwxProxyClient()
	if err != nil {
		return nil, err
	}
	return cli.Register(context.Background(), req)
}

// Deregister sends SAR (Code 301) over diameter connection with ServerAssignmentType
// set to USER_DEREGISTRATION, waits (blocks) for SAA & returns its RPC representation
func Deregister(req *protos.RegistrationRequest) (*protos.RegistrationAnswer, error) {
	err := verifyRegistrationRequest(req)
	if err != nil {
		errMsg := fmt.Errorf("Invalid RegistrationRequest provided: %s", err)
		return nil, errors.New(errMsg.Error())
	}
	cli, err := getSwxProxyClient()
	if err != nil {
		return nil, err
	}
	return cli.Deregister(context.Background(), req)
}

func verifyAuthenticationRequest(req *protos.AuthenticationRequest) error {
	if req == nil {
		return fmt.Errorf("request is nil")
	}
	return verifyUsername(req.GetUserName())
}

func verifyRegistrationRequest(req *protos.RegistrationRequest) error {
	if req == nil {
		return fmt.Errorf("request is nil")
	}
	return verifyUsername(req.GetUserName())
}

func verifyUsername(username string) error {
	if len(username) == 0 {
		return fmt.Errorf("no username provided")
	} else if len(username) > 16 {
		return fmt.Errorf("username is too long (must be 16 digits or less)")
	}
	return nil
}
