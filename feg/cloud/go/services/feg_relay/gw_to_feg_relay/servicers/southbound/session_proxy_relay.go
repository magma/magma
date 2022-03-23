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
	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
)

// SessionProxyServer implementation
//
// CreateSession Notifies OCS/PCRF of new session and return rules associated with subscriber
// along with credits for each rule
func (s *RelayRouter) CreateSession(
	ctx context.Context, r *protos.CreateSessionRequest) (*protos.CreateSessionResponse, error) {

	// CreateSession's SID is just IMSI with "IMSI" prefix which findServingNHFeg() should remove
	client, ctx, cancel, err := s.getSessionControllerProxyClient(ctx, getSessionControllerProxyType(r.GetCommonContext().GetRatType()), r.GetCommonContext().GetSid().GetId())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.CreateSession(ctx, r)
}

// UpdateSession Updates OCS/PCRF with used credit and terminations from gateway
func (s *RelayRouter) UpdateSession(
	ctx context.Context, r *protos.UpdateSessionRequest) (*protos.UpdateSessionResponse, error) {

	// Group requests by [PLMNID][GwServiceType]
	// - PLMNID: Use the longest possible PLMN ID (6), in the worst case we'll fragment requests more then needed,
	//   but the routing will still be logically correct
	// - GwServiceType: Each request will be sent to either Session_Proxy or for N7_N40_proxy
	reqMap := map[string]map[gateway_registry.GwServiceType]*protos.UpdateSessionRequest{}
	for _, u := range r.GetUpdates() {
		if u == nil {
			continue
		}
		plmnid := getPlmnId6(u.GetCommonContext().GetSid().GetId())

		// Group requests by PLMNID
		pMap, ok := reqMap[plmnid]
		if !ok || pMap == nil {
			// Group requests by GwServiceType
			addReqMapEntry(reqMap, plmnid)
		}

		// Append the Updates to the Session_Proxy or N7_N40_Proxy map request,
		// based on the RAT Type
		proxyType := getSessionControllerProxyType(u.GetCommonContext().GetRatType())
		reqMap[plmnid][proxyType].Updates = append(reqMap[plmnid][proxyType].Updates, u)
	}
	for _, m := range r.GetUsageMonitors() {
		if m == nil {
			continue
		}
		plmnid := getPlmnId6(m.GetSid())
		pMap, ok := reqMap[plmnid]
		if !ok || pMap == nil {
			addReqMapEntry(reqMap, plmnid)
		}

		// Append the UsageMonitors to the Session_Proxy or N7_N40_Proxy map request,
		// based on the RAT Type
		proxyType := getSessionControllerProxyType(m.GetRatType())
		reqMap[plmnid][proxyType].UsageMonitors = append(reqMap[plmnid][proxyType].UsageMonitors, m)
	}

	// Each request in reqMap is segregated into two requests and
	// relays each request to respective proxies (session_proxy and n7n40_proxy).
	// Hence length of resultChan is len(reqMap) * 2

	resultChan := make(chan *protos.UpdateSessionResponse, len(reqMap)*2)
	var numOfProxyRreq int

	// send a separate Update request for each unique PLMN ID
	for plmnid, proxyMap := range reqMap {
		for gwProxyType, req := range proxyMap {
			if (req == nil) || ((len(req.Updates) == 0) && (len(req.UsageMonitors) == 0)) {
				continue
			}
			numOfProxyRreq++

			go func(plmnid string, req *protos.UpdateSessionRequest, gwProxyType gateway_registry.GwServiceType) {
				client, ctx, cancel, err := s.getSessionControllerProxyClient(ctx, gwProxyType, plmnid)
				if err != nil {
					glog.Errorf("failed connect to %s for PLMNID '%s': %v", gwProxyType, plmnid, err)
					resultChan <- genUpdateErrorResp(req)
					return
				}
				defer cancel()
				resp, err := client.UpdateSession(ctx, req)
				if err != nil {
					glog.Errorf("failed %s Update for PLMNID '%s': %v", gwProxyType, plmnid, err)
					resultChan <- genUpdateErrorResp(req)
					return
				}
				resultChan <- resp
			}(plmnid, req, gwProxyType)
		}
	}

	// Combine received responses
	resp := &protos.UpdateSessionResponse{
		Responses:             []*protos.CreditUpdateResponse{},
		UsageMonitorResponses: []*protos.UsageMonitoringUpdateResponse{},
	}
	for i := numOfProxyRreq; i > 0; i-- {
		nhResp := <-resultChan
		resp.Responses = append(resp.Responses, nhResp.GetResponses()...)
		resp.UsageMonitorResponses = append(resp.UsageMonitorResponses, nhResp.GetUsageMonitorResponses()...)
	}
	// leave resultChan open, this thread is a 'reader'
	return resp, nil
}

// TerminateSession Terminates session in OCS/PCRF for a subscriber
func (s *RelayRouter) TerminateSession(
	ctx context.Context, r *protos.SessionTerminateRequest) (*protos.SessionTerminateResponse, error) {

	// TerminateSession's SID is just IMSI with "IMSI" prefix which findServingNHFeg() should remove
	client, ctx, cancel, err := s.getSessionControllerProxyClient(ctx, getSessionControllerProxyType(r.GetCommonContext().GetRatType()), r.GetCommonContext().GetSid().GetId())
	if err != nil {
		return nil, err
	}
	defer cancel()
	return client.TerminateSession(ctx, r)
}

// getSessionControllerProxyClient Get Proxy Client based on GwServiceType
func (s *RelayRouter) getSessionControllerProxyClient(c context.Context, gwProxyType gateway_registry.GwServiceType, imsi string) (protos.CentralSessionControllerClient, context.Context, context.CancelFunc, error) {
	conn, ctx, cancel, err := s.GetFegServiceConnection(c, imsi, gwProxyType)
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

// getSessionControllerProxyType Service selection based on the RAT Type
func getSessionControllerProxyType(r protos.RATType) gateway_registry.GwServiceType {
	switch r {
	case protos.RATType_TGPP_NR:
		return FegN7N40Proxy
	default:
		return FegSessionProxy
	}
}

// addReqMapEntry Defines the map for Session_Proxy and N7_N40_Proxy per PLMN
func addReqMapEntry(reqMap map[string]map[gateway_registry.GwServiceType]*protos.UpdateSessionRequest,
	plmnid string) {

	reqMap[plmnid] = map[gateway_registry.GwServiceType]*protos.UpdateSessionRequest{}
	reqMap[plmnid][FegSessionProxy] = &protos.UpdateSessionRequest{
		Updates:       []*protos.CreditUsageUpdate{},
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{},
	}
	reqMap[plmnid][FegN7N40Proxy] = &protos.UpdateSessionRequest{
		Updates:       []*protos.CreditUsageUpdate{},
		UsageMonitors: []*protos.UsageMonitoringUpdateRequest{},
	}
}
