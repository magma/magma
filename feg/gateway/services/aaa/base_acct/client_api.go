/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package base_acct provides a client API for interacting with the
// base_acct cloud service
package base_acct

import (
	"context"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"magma/feg/cloud/go/protos"
	"magma/orc8r/lib/go/errors"
	platformregistry "magma/orc8r/lib/go/registry"
)

const (
	ServiceName = "BASE_ACCT"
)

// Get a thin RPC client to the gateway base_acct service.
func getBaseAcctClient() (protos.AccountingClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	conn, err = platformregistry.Get().GetSharedCloudConnection(strings.ToLower(ServiceName))
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewAccountingClient(conn), nil
}

// Start will be called at the end of every new user session creation start is responsible for verification & initiation
// of an accounting contract between the user identity provider/MNO and service provider (ISP/WISP/PLTE) A non-error
// return will indicate successful contract establishment and will result in the beginning of service for the user
func Start(request *protos.AcctSession) (*protos.AcctSessionResp, error) {
	client, err := getBaseAcctClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.Start(context.Background(), request)
	if err != nil {
		glog.Error(err)
	}
	return resp, err
}

// Update should be continuously called for every ongoing service session to update the user bandwidth usage as well as
// current quality of provided service. If update returns error the session should be terminated and the user
// disconnected, In the case of unsuccessful update completion, service provider is suppose to follow up with
// final Stop call
func Update(request *protos.AcctUpdateReq) (*protos.AcctSessionResp, error) {
	client, err := getBaseAcctClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.Update(context.Background(), request)
	if err != nil {
		glog.Error(err)
	}
	return resp, err
}

// Stop is a notification call to communicate to identity provider user/network initiated service termination.
// stop will provide final used bandwidth count. stop call is issued after the user session was terminated.
func Stop(request *protos.AcctUpdateReq) (*protos.AcctStopResp, error) {
	client, err := getBaseAcctClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.Stop(context.Background(), request)
	if err != nil {
		glog.Error(err)
	}
	return resp, err
}
