/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	orc8rregistry "magma/orc8r/lib/go/registry"

	"google.golang.org/grpc"
)

const (
	GrpcMaxDelaySec   = 10
	GrpcMaxTimeoutSec = 10
)

// ConnectionRegistry defines an interface to get a gRPC connection to a
// a service.
type ConnectionRegistry interface {
	GetConnection(addr string, port int) (*grpc.ClientConn, error)
}

type k8sConnectionRegistry struct {
	sync.RWMutex
	*orc8rregistry.ServiceRegistry
}

// NewK8sConnectionRegistry creates and initializes a connection registry.
func NewK8sConnectionRegistry() *k8sConnectionRegistry {
	return &k8sConnectionRegistry{
		sync.RWMutex{},
		orc8rregistry.New(),
	}
}

// GetConnection gets a connection to a kubernetes service at service:port.
// The connection implementation uses orc8r/lib's service registry, which will
// reuse a gRPC connection if it already exists.
func (r *k8sConnectionRegistry) GetConnection(service string, port int) (*grpc.ClientConn, error) {
	serviceAddr := fmt.Sprintf("%s:%d", service, port)
	allServices, err := r.ListAllServices()
	if err != nil {
		return nil, err
	}
	exists := doesServiceExist(allServices, serviceAddr)
	if !exists {
		// Kubernetes services can be reached at svc:port
		// Here we map svc:port -> addr (svc:port) in the registry to ensure
		// that the connections will work properly even if service port changes
		loc := orc8rregistry.ServiceLocation{Name: serviceAddr, Host: service, Port: port}
		r.AddService(loc)
	}
	ctx, cancel := context.WithTimeout(context.Background(), GrpcMaxTimeoutSec*time.Second)
	defer cancel()

	return r.GetConnectionImpl(ctx, serviceAddr)
}

func doesServiceExist(serviceList []string, service string) bool {
	for _, svc := range serviceList {
		if svc == service {
			return true
		}
	}
	return false
}
