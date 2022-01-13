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
	"fmt"
	"net/http"
	"sync"

	"github.com/golang/glog"

	"magma/feg/gateway/services/n7_n40_proxy/metrics"
	"magma/feg/gateway/services/n7_n40_proxy/n7"
	"magma/lte/cloud/go/protos"
)

// sendMutlipleSmPolicyUpdateRequests sends multiple parallel update requests to PCF and
// returns the accumulated the responses from PCF
func (srv *CentralSessionController) sendMutlipleSmPolicyUpdateRequests(
	reqCtxs []*n7.SmPolicyUpdateReqCtx,
) []*protos.UsageMonitoringUpdateResponse {
	var wg sync.WaitGroup
	respChan := make(chan []*protos.UsageMonitoringUpdateResponse)
	ctx, cancel := context.WithTimeout(context.Background(), srv.config.RequestTimeout)
	defer cancel()

	accResponses := []*protos.UsageMonitoringUpdateResponse{}
	for _, reqCtx := range reqCtxs {
		tmpReqCtx := reqCtx // don't use loop variable in func closure
		wg.Add(1)
		// Send updates in parallel and accumulate results when all done
		go func() {
			defer wg.Done()
			responses := srv.sendSingleSmPolicyUpdate(ctx, tmpReqCtx)
			respChan <- responses
		}()
	}

	// goroutine that waits for all updates to complete and close the response channel
	go func() {
		wg.Wait()
		close(respChan)
	}()

	// Accumulate the responses. The channel is closed when all sends are done
	for responses := range respChan {
		accResponses = append(accResponses, responses...)
	}
	return accResponses
}

func (srv *CentralSessionController) sendSingleSmPolicyUpdate(
	ctx context.Context,
	updateCtx *n7.SmPolicyUpdateReqCtx,
) []*protos.UsageMonitoringUpdateResponse {
	resp, err := srv.policyClient.PostSmPoliciesSmPolicyIdUpdateWithResponse(
		ctx, updateCtx.SmPolicyId, *updateCtx.ReqBody)
	if err == nil && resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("http error status-code=%d", resp.StatusCode())
	}
	metrics.ReportUpdateSmPolicy(err)
	if err != nil {
		glog.Errorf("SmPolicyUpdate request failed: %s policyId=%s", err, updateCtx.SmPolicyId)
		// Return failure usage monitoring response
		response := n7.GetUsageMonitoringUpdateResponseProto(updateCtx, false)
		return []*protos.UsageMonitoringUpdateResponse{response}
	}
	return n7.GetUsageMonitoringResponsesProto(updateCtx, resp.JSON200)
}
