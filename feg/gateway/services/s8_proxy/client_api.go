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

package s8_proxy

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

type s8ProxyClient struct {
	protos.S8ProxyClient
}

func CreateSession(req *protos.CreateSessionRequestPgw) (*protos.CreateSessionResponsePgw, error) {
	if req == nil {
		return nil, errors.New("Invalid CreateSessionRequestPgw")
	}
	cli, err := getS8ProxyClient()
	if err != nil {
		return nil, err
	}
	return cli.CreateSession(context.Background(), req)
}

func DeleteSession(req *protos.DeleteSessionRequestPgw) (*protos.DeleteSessionResponsePgw, error) {
	if req == nil {
		return nil, errors.New("Invalid CreateSessionRequestPgw")
	}
	cli, err := getS8ProxyClient()
	if err != nil {
		return nil, err
	}
	return cli.DeleteSession(context.Background(), req)
}

func SendEcho(req *protos.EchoRequest) (*protos.EchoResponse, error) {
	if req == nil {
		return nil, errors.New("Invalid CreateSessionRequestPgw")
	}
	cli, err := getS8ProxyClient()
	if err != nil {
		return nil, err
	}
	return cli.SendEcho(context.Background(), req)
}

func getS8ProxyClient() (*s8ProxyClient, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	if util.GetEnvBool("USE_REMOTE_S8_PROXY") {
		conn, err = registry.Get().GetSharedCloudConnection(strings.ToLower(registry.S8_PROXY))
	} else {
		conn, err = registry.GetConnection(registry.S8_PROXY)
	}
	if err != nil {
		errMsg := fmt.Sprintf("S8 Proxy client initialization error: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return &s8ProxyClient{protos.NewS8ProxyClient(conn)}, nil
}
