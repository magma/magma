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

// Package aka implements EAP-AKA provider
package provider

import (
	"context"
	"errors"
	"fmt"

	"magma/feg/gateway/services/eap/providers"
	"magma/feg/gateway/services/eap/providers/aka/servicers"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
	eapp "magma/feg/gateway/services/eap/protos"
	_ "magma/feg/gateway/services/eap/providers/aka/servicers/handlers"
)

// Wrapper to provide a wrapper for GRPC Client to extend it with Cleanup
// functionality
type akaClient struct {
	eapp.EapServiceClient
	cc *grpc.ClientConn
}

func (cl *akaClient) Cleanup() {
	if cl != nil && cl.cc != nil {
		cl.cc.Close()
	}
}

// getAKAClient is a utility function to get a RPC connection to the EAP service
func getAKAClient() (*akaClient, error) {
	conn, err := registry.GetConnection(registry.EAP_AKA)
	if err != nil {
		errMsg := fmt.Sprintf("EAP client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &akaClient{
		eapp.NewEapServiceClient(conn),
		conn,
	}, err
}

// Handle handles passed EAP-AKA payload & returns corresponding result
// this Handle implementation is using GRPC based AKA provider service
func (*providerImpl) Handle(msg *protos.Eap) (*protos.Eap, error) {
	if msg == nil {
		return nil, errors.New("Invalid EAP AKA Message")
	}
	cli, err := getAKAClient()
	if err != nil {
		return nil, err
	}
	return cli.Handle(context.Background(), msg)
}

func NewService(_ *servicers.EapAkaSrv) providers.Method {
	return New()
}
