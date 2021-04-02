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

package servicers

import (
	"context"

	"github.com/golang/glog"

	"magma/lte/cloud/go/protos"
)

// SessionProxyServer implementation
//
// Notify OCS/PCRF of new session and return rules associated with subscriber
// along with credits for each rule
func (s *RelayRouter) CreateSession(
	ctx context.Context, r *protos.CreateSessionRequest) (*protos.CreateSessionResponse, error) {

	// CreateSession's SID is just IMSI with "IMSI" prefix which findServingNHFeg() should remove
	client, ctx, cancel, err := s.getSessionProxyClient(ctx, r.GetCommonContext().GetSid().GetId())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.CreateSession(ctx, r)
}

// Updates OCS/PCRF with used credit and terminations from gateway
func (s *RelayRouter) UpdateSession(
	ctx context.Context, r *protos.UpdateSessionRequest) (*protos.UpdateSessionResponse, error) {

	// Group requests by PLMNID
	// Use the longest possible PLMN ID (6), in the worst case we'll fragment requests more then needed, but the routing
	// will still be logically correct
	reqMap := map[string]*protos.UpdateSessionRequest{}
	for _, u := range r.GetUpdates() {
		if u == nil {
			continue
		}
		plmnid := getPlmnId6(u.GetCommonContext().GetSid().GetId())
		req, ok := reqMap[plmnid]
		if !ok || req == nil {
			req = &protos.UpdateSessionRequest{Updates: []*protos.CreditUsageUpdate{u}}
			reqMap[plmnid] = req
		} else {
			req.Updates = append(req.Updates, u)
		}
	}
	for _, m := range r.GetUsageMonitors() {
		if m == nil {
			continue
		}
		plmnid := getPlmnId6(m.GetSid())
		req, ok := reqMap[plmnid]
		if !ok || req == nil {
			req = &protos.UpdateSessionRequest{UsageMonitors: []*protos.UsageMonitoringUpdateRequest{m}}
			reqMap[plmnid] = req
		} else {
			req.UsageMonitors = append(req.UsageMonitors, m)
		}
	}
	resultChan := make(chan *protos.UpdateSessionResponse, len(reqMap))

	// send a separate Update request for each unique PLMN ID
	for plmnid, req := range reqMap {
		go func(plmnid string, req *protos.UpdateSessionRequest) {
			client, ctx, cancel, err := s.getSessionProxyClient(ctx, plmnid)
			if err != nil {
				glog.Errorf("failed connect to Session Proxy for PLMNID '%s': %v", plmnid, err)
				resultChan <- genUpdateErrorResp(req)
				return
			}
			defer cancel()
			resp, err := client.UpdateSession(ctx, req)
			if err != nil {
				glog.Errorf("failed Session Proxy Update for PLMNID '%s': %v", plmnid, err)
				resultChan <- genUpdateErrorResp(req)
				return
			}
			resultChan <- resp
		}(plmnid, req)
	}
	// Combine received responses
	resp := &protos.UpdateSessionResponse{
		Responses:             []*protos.CreditUpdateResponse{},
		UsageMonitorResponses: []*protos.UsageMonitoringUpdateResponse{},
	}
	for i := len(reqMap); i > 0; i-- {
		nhResp := <-resultChan
		resp.Responses = append(resp.Responses, nhResp.GetResponses()...)
		resp.UsageMonitorResponses = append(resp.UsageMonitorResponses, nhResp.GetUsageMonitorResponses()...)
	}
	// leave resultChan open, this thread is a 'reader'
	return resp, nil
}

// Terminates session in OCS/PCRF for a subscriber
func (s *RelayRouter) TerminateSession(
	ctx context.Context, r *protos.SessionTerminateRequest) (*protos.SessionTerminateResponse, error) {

	// TerminateSession's SID is just IMSI with "IMSI" prefix which findServingNHFeg() should remove
	client, ctx, cancel, err := s.getSessionProxyClient(ctx, r.GetCommonContext().GetSid().GetId())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.TerminateSession(ctx, r)
}

func (s *RelayRouter) getSessionProxyClient(
	c context.Context, imsi string) (protos.CentralSessionControllerClient, context.Context, context.CancelFunc, error) {

	conn, ctx, cancel, err := s.GetFegServiceConnection(c, imsi, FegSessionProxy)
	if err != nil {
		return nil, nil, nil, err
	}
	return protos.NewCentralSessionControllerClient(conn), ctx, cancel, nil
}

func genUpdateErrorResp(req *protos.UpdateSessionRequest) *protos.UpdateSessionResponse {
	resp := &protos.UpdateSessionResponse{
		Responses:             []*protos.CreditUpdateResponse{},
		UsageMonitorResponses: []*protos.UsageMonitoringUpdateResponse{},
	}
	for _, u := range req.GetUpdates() {
		resp.Responses = append(resp.Responses, &protos.CreditUpdateResponse{
			Success:    false,
			Sid:        u.GetCommonContext().GetSid().GetId(),
			SessionId:  u.GetSessionId(),
			ResultCode: DiamUnableToDeliverErr,
		})
	}
	for _, m := range req.GetUsageMonitors() {
		resp.UsageMonitorResponses = append(resp.UsageMonitorResponses, &protos.UsageMonitoringUpdateResponse{
			Success:    false,
			Sid:        m.GetSid(),
			SessionId:  m.GetSessionId(),
			ResultCode: DiamUnableToDeliverErr,
		})
	}
	return resp
}
