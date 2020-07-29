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

// package status implements magmad status collector & reporter
package status

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/golang/glog"

	"magma/gateway/service_registry"
	"magma/gateway/status"
	"magma/orc8r/lib/go/protos"
)

type serviceInfoClientState struct {
	queryInFlight bool
	states        []*protos.State
}

var (
	pollerMu     sync.Mutex
	clientStates map[string]*serviceInfoClientState = map[string]*serviceInfoClientState{}
)

func startServiceQuery(service string, queryMetrics bool, maxMetricsQueueSz int) error {
	conn, err := service_registry.Get().GetConnection(strings.ToUpper(service))
	if err != nil {
		return err
	}
	client := protos.NewService303Client(conn)

	pollerMu.Lock()
	defer pollerMu.Unlock()

	clState, ok := clientStates[service]
	if !ok || clState == nil {
		clState = &serviceInfoClientState{}
		clientStates[service] = clState
	}
	if clState.queryInFlight {
		return fmt.Errorf("'%s' fb303 client is still in use", service)
	}
	clState.queryInFlight = true

	go func() {
		var (
			serviceMetrics *protos.MetricsContainer
			merr           error
		)

		ctx := context.Background()
		statesResp, err := client.GetOperationalStates(ctx, &protos.Void{})
		if err != nil {
			glog.Errorf("service '%s' GetServiceInfo error: %v", service, err)
		}
		if queryMetrics {
			serviceMetrics, merr = client.GetMetrics(ctx, &protos.Void{})
			if merr != nil {
				glog.Errorf("service '%s' GetMetrics error: %v", service, err)
			}
		}
		pollerMu.Lock()
		clState.queryInFlight = false
		if err == nil {
			clState.states = statesResp.States
		}
		pollerMu.Unlock()

		if queryMetrics && merr == nil {
			if qlen := enqueueMetrics(service, serviceMetrics); qlen > maxMetricsQueueSz && glog.V(1) {
				glog.Warningf("metrics queue length %d exceeds max %d", qlen, maxMetricsQueueSz)
			}
		}
	}()
	return nil
}

func collect() *protos.ReportStatesRequest {
	marshaledGwState, err := json.Marshal(status.GetGatewayStatus())
	if err != nil {
		marshaledGwState, _ = json.Marshal(&status.GatewayStatus{
			Meta: map[string]string{"error": err.Error()},
		})
		glog.Errorf("failed to marshal gw_state: %v", err)
	}
	states := []*protos.State{{
		Type:     "gw_state",
		DeviceID: status.GetHwId(),
		Value:    marshaledGwState,
	}}

	pollerMu.Lock()
	for service, clientState := range clientStates {
		if clientState != nil && len(clientState.states) > 0 {
			states = append(states, clientState.states...)
			clientState.states = nil
		} else {
			glog.V(2).Infof("no states to report from '%s'", service)
		}
	}
	pollerMu.Unlock()

	return &protos.ReportStatesRequest{States: states}
}
