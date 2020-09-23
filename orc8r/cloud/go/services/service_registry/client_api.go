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

package service_registry

import (
	"context"

	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

func getServiceRegistryClient() (protos.ServiceRegistryClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewServiceRegistryClient(conn), err
}

// ListAllServices returns the service name of all services in the registry.
func ListAllServices() ([]string, error) {
	client, err := getServiceRegistryClient()
	if err != nil {
		return []string{}, err
	}
	req := &protos.Void{}
	res, err := client.ListAllServices(context.Background(), req)
	if err != nil {
		return []string{}, err
	}
	return res.GetServices(), nil
}

// FindServices returns all services in that have the provided label.
func FindServices(label string) ([]string, error) {
	client, err := getServiceRegistryClient()
	if err != nil {
		return []string{}, err
	}
	req := &protos.FindServicesRequest{
		Label: label,
	}
	res, err := client.FindServices(context.Background(), req)
	if err != nil {
		return []string{}, err
	}
	return res.GetServices(), nil
}

// GetServiceAddress return the address of the gRPC server for the provided
// service.
func GetServiceAddress(service string) (string, error) {
	client, err := getServiceRegistryClient()
	if err != nil {
		return "", err
	}
	req := &protos.GetServiceAddressRequest{
		Service: service,
	}
	res, err := client.GetServiceAddress(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.GetAddress(), nil
}

// GetHttpServerAddress returns the address of the HTTP server for the provided
// service.
func GetHttpServerAddress(service string) (string, error) {
	client, err := getServiceRegistryClient()
	if err != nil {
		return "", err
	}
	req := &protos.GetHttpServerAddressRequest{
		Service: service,
	}
	res, err := client.GetHttpServerAddress(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.GetAddress(), nil
}

// GetAnnotation returns the annotation value for the provided service and
// annotation.
func GetAnnotation(service string, annotation string) (string, error) {
	client, err := getServiceRegistryClient()
	if err != nil {
		return "", err
	}
	req := &protos.GetAnnotationRequest{
		Service:    service,
		Annotation: annotation,
	}
	res, err := client.GetAnnotation(context.Background(), req)
	if err != nil {
		return "", err
	}
	return res.GetAnnotationValue(), nil
}
