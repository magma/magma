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

package credit_control

import (
	"context"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"

	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/session_proxy/relay"
	"magma/gateway/service_registry"
	"magma/lte/cloud/go/protos"
)

type asrHandler struct {
	diamClient *diameter.Client
	registry   service_registry.GatewayRegistry
}

func NewASRHandler(diamClient *diameter.Client, cloudRegistry service_registry.GatewayRegistry) diam.Handler {
	return &asrHandler{diamClient: diamClient, registry: cloudRegistry}
}

func (h *asrHandler) ServeDIAM(conn diam.Conn, m *diam.Message) {
	if h == nil || h.diamClient == nil {
		glog.Errorf("Invalid ASR Handler")
	}
	asr := &diameter.ASR{}
	if err := m.Unmarshal(asr); err != nil {
		glog.Errorf("Received unparseable ASR %s\n%s", m, err)
		return
	}
	go func() {
		var err error
		imsi := string(asr.UserName)
		if len(imsi) == 0 {
			imsi, err = diameter.ExtractImsiFromSessionID(asr.SessionID)
			if err != nil {
				glog.Errorf("Error retreiving IMSI from Session ID %s: %s", asr.SessionID, err)
				h.sendASA(conn, m, asr.SessionID, diam.UnknownSessionID)
				return
			}
		}
		client, err := relay.GetAbortSessionResponderClient(h.registry)
		if err != nil {
			glog.Error(err)
			h.sendASA(conn, m, asr.SessionID, diam.UnableToDeliver)
			return
		}
		defer client.Close()

		res, err := client.AbortSession(context.Background(), &protos.AbortSessionRequest{
			UserName:  imsi,
			SessionId: diameter.DecodeSessionID(asr.SessionID),
		})
		if err != nil {
			glog.Errorf("Error relaying ASR to gateway: %s", err)
			h.sendASA(conn, m, asr.SessionID, diam.UnableToDeliver)
			return
		}
		var resCode uint32
		switch res.GetCode() {
		case protos.AbortSessionResult_GATEWAY_NOT_FOUND:
			glog.Errorf("Failed ASR to gateway: %s", res.GetErrorMessage())
			resCode = diam.UnableToDeliver
		case protos.AbortSessionResult_SESSION_NOT_FOUND:
			glog.Errorf("Unknown Session in ASR: %s", res.GetErrorMessage())
			resCode = diam.UnknownSessionID
		case protos.AbortSessionResult_USER_NOT_FOUND:
			glog.Errorf("Unknown User in ASR: %s", res.GetErrorMessage())
			resCode = diam.UnknownUser
		case protos.AbortSessionResult_SESSION_REMOVED:
			resCode = diam.Success
		default:
			if len(res.GetErrorMessage()) > 0 {
				glog.Errorf("Limited ASR Success: %s", res.GetErrorMessage())
			}
			resCode = diam.LimitedSuccess
		}
		h.sendASA(conn, m, asr.SessionID, resCode)
	}()
}

func (h *asrHandler) sendASA(conn diam.Conn, m *diam.Message, sid string, code uint32) {
	asaMsg := m.Answer(code)
	asaMsg.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid)))
	asaMsg = h.diamClient.AddOriginAVPsToMessage(asaMsg)
	_, err := asaMsg.WriteToWithRetry(conn, h.diamClient.Retries())
	if err != nil {
		glog.Errorf(
			"ASA Write Failed for %s->%s, SessionID: %s - %v",
			conn.LocalAddr(), conn.RemoteAddr(), sid, err)
		conn.Close() // close connection on error
	}
}
