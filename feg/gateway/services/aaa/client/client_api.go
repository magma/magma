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

// Package client provides a thin API client for communicating with AAA Server.
// This can be used by apps to discover and contact the service, without knowing about
// the underlying RPC implementation.
package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/glog"

	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/aaa/protos"
)

type aaaClient struct {
	protos.AuthenticatorClient
	protos.AccountingClient
}

// getAaaClient is a utility function to get a RPC connection to the AAA service providing
// Authenticator & Accounting RPCs
func getAaaClient() (*aaaClient, error) {
	conn, err := registry.GetConnection(registry.AAA_SERVER)
	if err != nil {
		errMsg := fmt.Sprintf("AAA Server client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &aaaClient{
		protos.NewAuthenticatorClient(conn),
		protos.NewAccountingClient(conn),
	}, err
}

// SupportedMethods returns sorted list (ascending, by type) of registered EAP Provider Methods
func SupportedMethods() (*protos.EapMethodList, error) {
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.SupportedMethods(context.Background(), &protos.Void{})
}

// HandleIdentity passes Identity EAP payload to corresponding method provider & returns corresponding
// EAP result
// NOTE: Identity Request is handled by APs & does not involve EAP Authenticator's support
func HandleIdentity(in *protos.EapIdentity) (*protos.Eap, error) {
	if in == nil {
		return nil, errors.New("Nil EapIdentity Parameter")
	}
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.HandleIdentity(context.Background(), in)
}

// Handle handles passed EAP payload & returns corresponding EAP result
func Handle(in *protos.Eap) (*protos.Eap, error) {
	if in == nil {
		return nil, errors.New("Nil Eap Parameter")
	}
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.Handle(context.Background(), in)
}

// Start implements Radius Acct-Status-Type: Start endpoint
func Start(aaaCtx *protos.Context) (*protos.AcctResp, error) {
	if aaaCtx == nil {
		return nil, errors.New("Nil AAA Ctx")
	}
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.Start(context.Background(), aaaCtx)
}

// Acct-Status-Type Stop
func InterimUpdate(ur *protos.UpdateRequest) (*protos.AcctResp, error) {
	if ur == nil {
		return nil, errors.New("Nil Interim Update Request")
	}
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.InterimUpdate(context.Background(), ur)
}

// Stop implements Radius Acct-Status-Type: Stop endpoint
func Stop(req *protos.StopRequest) (*protos.AcctResp, error) {
	if req == nil {
		return nil, errors.New("Nil Stop Request")
	}
	cli, err := getAaaClient()
	if err != nil {
		return nil, err
	}
	return cli.Stop(context.Background(), req)
}
