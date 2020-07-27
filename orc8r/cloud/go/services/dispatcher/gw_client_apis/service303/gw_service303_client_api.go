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

package service303

import (
	"errors"
	"fmt"

	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

func getGWClient(service gateway_registry.GwServiceType, hwId string) (protos.Service303Client, context.Context, error) {
	conn, ctx, err := gateway_registry.GetGatewayConnection(service, hwId)
	if err != nil {
		errMsg := fmt.Sprintf("service303 gwClient initialization error: %s", err)
		glog.Error(errMsg)
		return nil, nil, errors.New(errMsg)
	}
	return protos.NewService303Client(conn), ctx, nil

}

func GWService303GetServiceInfo(service gateway_registry.GwServiceType, hwId string) (*protos.ServiceInfo, error) {
	client, ctx, err := getGWClient(service, hwId)
	if err != nil {
		return nil, err
	}
	return client.GetServiceInfo(ctx, new(protos.Void))
}

func GWService303GetMetrics(service gateway_registry.GwServiceType, hwId string) (*protos.MetricsContainer, error) {
	client, ctx, err := getGWClient(service, hwId)
	if err != nil {
		return nil, err
	}
	return client.GetMetrics(ctx, new(protos.Void))
}

func GWService303StopService(service gateway_registry.GwServiceType, hwId string) error {
	client, ctx, err := getGWClient(service, hwId)
	if err != nil {
		return err
	}
	_, err = client.StopService(ctx, new(protos.Void))
	return err
}

func GWService303SetLogLevel(service gateway_registry.GwServiceType, hwId string, in *protos.LogLevelMessage) error {
	client, ctx, err := getGWClient(service, hwId)
	if err != nil {
		return err
	}
	_, err = client.SetLogLevel(ctx, in)
	return err
}

func GWService303SetLogVerbosity(service gateway_registry.GwServiceType, hwId string, in *protos.LogVerbosity) error {
	client, ctx, err := getGWClient(service, hwId)
	if err != nil {
		return err
	}
	_, err = client.SetLogVerbosity(ctx, in)
	return err
}
