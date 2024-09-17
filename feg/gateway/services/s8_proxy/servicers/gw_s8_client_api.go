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

package servicers

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/glog"
	"google.golang.org/grpc"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/feg_relay"
	"magma/feg/gateway/registry"
	orc8r_protos "magma/orc8r/lib/go/protos"
)

// GWS8ProxyCreateBearerRequest forwards Create Bearer Request to FegRelay and
// FegRelay then to AGW
func GWS8ProxyCreateBearerRequest(in *protos.CreateBearerRequestPgw) (*orc8r_protos.Void, error) {
	conn, err := getCloudConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := protos.NewS8ProxyResponderClient(conn)
	return client.CreateBearer(context.Background(), in)
}

// GWS8ProxyDeleteBearerRequest forwards Delete Bearer Request to FegRelay and
// FegRelay then to AGW
func GWS8ProxyDeleteBearerRequest(in *protos.DeleteBearerRequestPgw) (*orc8r_protos.Void, error) {
	conn, err := getCloudConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := protos.NewS8ProxyResponderClient(conn)
	return client.DeleteBearerRequest(context.Background(), in)
}

func getCloudConn() (*grpc.ClientConn, error) {
	conn, err := registry.Get().GetCloudConnection(feg_relay.ServiceName)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to establish connection to cloud FegToGwRelayClient: %s", err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return conn, nil
}
