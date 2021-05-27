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

package relay

import (
	"fmt"
	"strings"

	"google.golang.org/grpc"

	"magma/feg/cloud/go/services/feg_relay"
	"magma/gateway/service_registry"
	"magma/lte/cloud/go/protos"
)

type CloseableSessionProxyResponderClient struct {
	protos.SessionProxyResponderClient
	conn *grpc.ClientConn
}

func (client *CloseableSessionProxyResponderClient) Close() {
	client.conn.Close()
}

// GetSessionProxyResponderClient Get a client to the local session manager client.
// To avoid leaking connections, defer Close() on the returned client.
func GetSessionProxyResponderClient(
	cloudRegistry service_registry.GatewayRegistry) (*CloseableSessionProxyResponderClient, error) {

	conn, err := cloudRegistry.GetCloudConnection(feg_relay.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to gw relay: %s", err)
	}
	return &CloseableSessionProxyResponderClient{
		SessionProxyResponderClient: protos.NewSessionProxyResponderClient(conn),
		conn:                        conn,
	}, nil
}

type CloseableAbortSessionResponderClient struct {
	protos.AbortSessionResponderClient
	conn *grpc.ClientConn
}

func (client *CloseableAbortSessionResponderClient) Close() {
	client.conn.Close()
}

// GetAbortSessionResponderClient returns a client to the local abort session client. To avoid leaking
// connections, defer Close() on the returned client.
func GetAbortSessionResponderClient(
	cloudRegistry service_registry.GatewayRegistry) (*CloseableAbortSessionResponderClient, error) {

	conn, err := cloudRegistry.GetCloudConnection(feg_relay.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to gw relay: %s", err)
	}
	return &CloseableAbortSessionResponderClient{
		AbortSessionResponderClient: protos.NewAbortSessionResponderClient(conn),
		conn:                        conn,
	}, nil
}

func GetIMSIFromSessionID(sessionID string) (string, error) {
	split := strings.Split(sessionID, "-")
	if len(split) < 2 {
		return "", fmt.Errorf("Session ID %s does not match format 'IMSI-RandNum'", sessionID)
	}
	return split[0], nil
}
