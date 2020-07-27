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

// Package service_health encapsulates service functionality related to health
// that service303 services can extend themselves with
package service_health

import (
	"context"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/orc8r/lib/go/errors"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

// getClient is a utility function to get an RPC connection to
// ServiceHealth
func getClient(service string) (protos.ServiceHealthClient, error) {
	conn, err := registry.GetConnection(service)
	if err != nil {
		initErr := errors.NewInitError(err, service)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewServiceHealthClient(conn), nil
}

// Disable disables service functionality for the period of time
// specified in the DisableMessage for the service provided
func Disable(service string, req *protos.DisableMessage) error {
	client, err := getClient(service)
	if err != nil {
		return err
	}
	_, err = client.Disable(context.Background(), req)
	return err
}

// Enable enables service functionality for the service provided
func Enable(service string) error {
	client, err := getClient(service)
	if err != nil {
		return err
	}
	_, err = client.Enable(context.Background(), &orcprotos.Void{})
	return err
}

// GetHealthStatus returns a HealthStatus object that indicates the current health of
// the service provided
func GetHealthStatus(service string) (*protos.HealthStatus, error) {
	client, err := getClient(service)
	if err != nil {
		return nil, err
	}
	healthStatus, err := client.GetHealthStatus(context.Background(), &orcprotos.Void{})
	return healthStatus, err
}
