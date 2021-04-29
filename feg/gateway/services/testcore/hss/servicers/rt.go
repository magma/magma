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

package servicers

import (
	"fmt"
	"strings"
	"time"

	"magma/feg/gateway/diameter"
	swx "magma/feg/gateway/services/swx_proxy/servicers"
	"magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//Permanently terminate the non-3gpp subscription
const PermanentTermination = 0

func (srv *HomeSubscriberServer) TerminateRegistration(sub *protos.SubscriberData) error {
	if sub.GetState().GetTgppAaaServerName() == "" {
		return fmt.Errorf("No AAA server found for subscriber: %s. Cannot send RTR", sub.GetSid().GetId())
	}
	aaaServerCfg, err := srv.genAAAServerConfig(sub.GetState().GetTgppAaaServerName())
	if err != nil {
		return fmt.Errorf("TerminateRegistration error: %s", err)
	}
	sid := (&diameter.DiameterClientConfig{}).GenSessionID("swx")

	ch := make(chan interface{})
	srv.requestTracker.RegisterRequest(sid, ch)
	// if request hasn't been removed by end of transaction, remove it
	defer srv.requestTracker.DeregisterRequest(sid)

	rtrMsg := srv.createRTR(sid, sub.GetSid().GetId())
	err = srv.sendDiameterMsg(rtrMsg, aaaServerCfg, maxDiamRetries)
	if err != nil {
		return err
	}
	select {
	case resp, open := <-ch:
		if !open {
			err = status.Errorf(codes.Aborted, "RTA for Session ID: %s is canceled", sid)
			glog.Error(err)
			return err
		}
		rta, ok := resp.(*swx.RTA)
		if !ok {
			err = status.Errorf(codes.Internal, "Invalid Response Type: %T, RTA expected.", resp)
			glog.Error(err)
			return err
		}
		if err = diameter.TranslateDiamResultCode(rta.ResultCode); err != nil {
			return err
		}
		// If there is no base diameter error, check that there is no experimental error either
		if err = diameter.TranslateDiamResultCode(rta.ExperimentalResult.ExperimentalResultCode); err != nil {
			return err
		}
		return srv.deregisterSubscriber(sub)

	case <-time.After(time.Second * timeoutSeconds):
		err = status.Errorf(codes.DeadlineExceeded, "RTA Timed Out for Session ID: %s", sid)
		glog.Error(err)
		return err
	}
}

func (srv *HomeSubscriberServer) sendDiameterMsg(msg *diam.Message, aaaCfg *diameter.DiameterServerConfig, retryCount uint) error {
	conn, err := srv.connMan.GetConnection(srv.smClient, aaaCfg)
	if err != nil {
		return err
	}
	err = conn.SendRequest(msg, retryCount)
	if err != nil {
		err = status.Errorf(codes.DataLoss, err.Error())
	}
	return err
}

// createRTR creates a Registration Termination Request with provided SessionID (sid)
// and userName to be sent over diameter to AAA Server
func (srv *HomeSubscriberServer) createRTR(sessionID string, username string) *diam.Message {
	msg := diameter.NewProxiableRequest(diam.RegistrationTermination, diam.TGPP_SWX_APP_ID, dict.Default)
	msg.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sessionID))
	msg.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_SWX_APP_ID)),
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
		},
	})

	// Set origin host and realm to server's host and realm since RTR is sent from HSS
	msg.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(srv.Config.Server.DestHost))
	msg.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(srv.Config.Server.DestRealm))
	msg.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(username))
	msg.NewAVP(avp.DeregistrationReason, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ReasonCode, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(PermanentTermination)),
			diam.NewAVP(avp.ReasonInfo, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("Manually initiated termination")),
		},
	})
	msg.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	return msg
}

func (srv *HomeSubscriberServer) deregisterSubscriber(subscriber *protos.SubscriberData) error {
	subscriber.State.TgppAaaServerRegistered = false
	subscriber.State.TgppAaaServerName = ""
	return srv.store.UpdateSubscriber(subscriber)
}

func (srv *HomeSubscriberServer) genAAAServerConfig(serverName string) (*diameter.DiameterServerConfig, error) {
	var destRealm string
	splitServerName := strings.Split(serverName, ".")
	if len(splitServerName) < 2 {
		destRealm = serverName
	} else {
		destRealm = strings.Join(splitServerName[1:], ".")
	}
	addr, ok := srv.clientMapping[serverName]
	if !ok {
		return nil, fmt.Errorf("could not find IP address for AAA server: %s", serverName)
	}
	return &diameter.DiameterServerConfig{
		DestHost:  serverName,
		DestRealm: destRealm,
		DiameterServerConnConfig: diameter.DiameterServerConnConfig{
			Addr:      addr,
			Protocol:  srv.Config.GetServer().GetProtocol(),
			LocalAddr: srv.Config.GetServer().GetAddress(),
		},
	}, nil
}

func handleRTA(srv *HomeSubscriberServer) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var rta swx.RTA
		err := m.Unmarshal(&rta)
		if err != nil {
			glog.Errorf("RTA Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		ch := srv.requestTracker.DeregisterRequest(rta.SessionID)
		if ch != nil {
			ch <- &rta
		} else {
			glog.Errorf("RTA SessionID %s not found. Message: %s, Remote: %s", rta.SessionID, m, c.RemoteAddr())
		}
	}
}
