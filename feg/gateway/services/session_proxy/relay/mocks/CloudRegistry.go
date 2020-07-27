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

package mocks

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	platform_registry "magma/orc8r/lib/go/registry"
)

type MockCloudRegistry struct {
	*platform_registry.ServiceRegistry
	ServerAddr string
}

// Mocked implementation which returns a grpc connection to the `ServerAddr`
// field in the struct.
func (m *MockCloudRegistry) GetCloudConnection(service string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, m.ServerAddr, grpc.WithInsecure())
	if err != nil {
		log.Printf("Err: %s", err)
		return nil, fmt.Errorf("Address: %s GRPC Dial error: %s", m.ServerAddr, err)
	} else if ctx.Err() != nil {
		log.Printf("Err: %s", ctx.Err())
		return nil, ctx.Err()
	}
	return conn, nil
}
func (m *MockCloudRegistry) GetConnection(service string) (*grpc.ClientConn, error) {
	return m.GetCloudConnection(service)
}
