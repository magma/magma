// +build !link_local_service

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

// Package sim implements EAP-SIM provider
package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	eapp "magma/feg/gateway/services/eap/protos"
	"magma/feg/gateway/services/eap/providers"
	"magma/feg/gateway/services/eap/providers/sim/servicers"
	_ "magma/feg/gateway/services/eap/providers/sim/servicers/handlers"
)

// Wrapper to provide a wrapper for GRPC Client to extend it with Cleanup
// functionality
type simClient struct {
	eapp.EapServiceClient
	cc *grpc.ClientConn
}

func (cl *simClient) Cleanup() {
	if cl != nil && cl.cc != nil {
		cl.cc.Close()
	}
}

// getSIMClient is a utility function to get a RPC connection to the EAP service
func getSIMClient() (*simClient, error) {
	conn, err := registry.GetConnection(registry.EAP_SIM)
	if err != nil {
		errMsg := fmt.Sprintf("EAP SIM client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &simClient{
		eapp.NewEapServiceClient(conn),
		conn,
	}, err
}

// Handle handles passed EAP-SIM payload & returns corresponding result
// this Handle implementation is using GRPC based SIM provider service
func (*providerImpl) Handle(msg *protos.Eap) (*protos.Eap, error) {
	if msg == nil {
		return nil, errors.New("Invalid EAP SIM Message")
	}
	cli, err := getSIMClient()
	if err != nil {
		return nil, err
	}
	return cli.Handle(context.Background(), msg)
}

func NewService(_ *servicers.EapSimSrv) providers.Method {
	return New()
}
