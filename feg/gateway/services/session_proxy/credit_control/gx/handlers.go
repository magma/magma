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

package gx

import (
	"context"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/golang/glog"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/metrics"
	"magma/feg/gateway/services/session_proxy/relay"
	"magma/gateway/service_registry"
	"magma/lte/cloud/go/protos"
)

// ccaHandler parses a CCADiameterMessage received over Gx and returns the
// `KeyAndAnswer` packed inside the CCA message.
func ccaHandler(message *diam.Message) diameter.KeyAndAnswer {
	var cca CCADiameterMessage
	glog.V(2).Infof("Received Gx CCA message:\n%s\n", message)
	if err := message.Unmarshal(&cca); err != nil {
		metrics.GxUnparseableMsg.Inc()
		glog.Errorf("Received unparseable CCA over Gx: %s", err)
		return diameter.KeyAndAnswer{}
	}
	sid := diameter.DecodeSessionID(cca.SessionID)
	return diameter.KeyAndAnswer{
		Key: credit_control.GetRequestKey(credit_control.Gx, sid, cca.RequestNumber),
		Answer: &CreditControlAnswer{
			ResultCode:             cca.ResultCode,
			ExperimentalResultCode: cca.ExperimentalResult.ExperimentalResultCode,
			SessionID:              sid,
			OriginHost:             cca.OriginHost,
			RequestNumber:          cca.RequestNumber,
			RuleInstallAVP:         cca.RuleInstalls,
			RuleRemoveAVP:          cca.RuleRemovals,
			UsageMonitors:          cca.UsageMonitors[:],
			EventTriggers:          cca.EventTriggers,
			RevalidationTime:       cca.RevalidationTime,
			Qos:                    cca.Qos,
		},
	}
}

type PolicyReAuthHandler func(request *PolicyReAuthRequest) *PolicyReAuthAnswer

// Factory function for a RAR message handler which relays to the corresponding
// gateway.
func GetGxReAuthHandler(cloudRegistry service_registry.GatewayRegistry, policyDBClient policydb.PolicyDBClient) PolicyReAuthHandler {
	return func(request *PolicyReAuthRequest) *PolicyReAuthAnswer {
		sid := diameter.DecodeSessionID(request.SessionID)
		imsi, err := protos.GetIMSIwithPrefixFromSessionId(sid)
		if err != nil {
			glog.Errorf("Error retrieving IMSI from session ID %s: %s", request.SessionID, err)
			return &PolicyReAuthAnswer{
				SessionID:  request.SessionID,
				ResultCode: diam.UnknownSessionID,
			}
		}

		client, err := relay.GetSessionProxyResponderClient(cloudRegistry)
		if err != nil {
			glog.Error(err)
			return &PolicyReAuthAnswer{
				SessionID:  request.SessionID,
				ResultCode: diam.UnableToDeliver,
			}
		}
		defer client.Close()

		gwReq := request.ToProto(imsi, sid, policyDBClient)
		ans, err := client.PolicyReAuth(context.Background(), gwReq)
		if err != nil {
			glog.Errorf("Error relaying Gx reauth request to gateway: %s", err)
			return &PolicyReAuthAnswer{
				SessionID:  request.SessionID,
				ResultCode: diam.UnableToDeliver,
			}
		}
		return (&PolicyReAuthAnswer{}).FromProto(request.SessionID, ans)
	}
}
