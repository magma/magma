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
	"fmt"

	"github.com/golang/glog"
	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
	"github.com/wmnsk/go-gtp/gtpv2/message"

	"magma/feg/cloud/go/protos"
)

// parseCreateSessionResponse parses a gtp message into a CreateSessionResponsePgw. In case
// the message is proper it also returns the session. In case it there is an error it returns
// the cause of error
func parseCreateSessionResponse(msg message.Message) (csRes *protos.CreateSessionResponsePgw, err error) {
	csResGtp := msg.(*message.CreateSessionResponse)
	glog.V(2).Infof("Received Create Session Response (gtp):\n%s", csResGtp.String())

	csRes = &protos.CreateSessionResponsePgw{}
	// check Cause value first.
	if causeIE := csResGtp.Cause; causeIE != nil {
		cause, err2 := causeIE.Cause()
		if err2 != nil {
			err = fmt.Errorf("Couldn't check cause of csRes: %s", err2)
			return
		}
		if cause != gtpv2.CauseRequestAccepted {
			err = &gtpv2.CauseNotOKError{
				MsgType: csResGtp.MessageTypeName(),
				Cause:   cause,
				Msg:     fmt.Sprintf("Not accepted"),
			}
			return
		}
	} else {
		err = &gtpv2.RequiredIEMissingError{Type: ie.Cause}
		return
	}

	// get values sent by pgw
	if paaIE := csResGtp.PAA; paaIE != nil {
		ip, err2 := paaIE.IPAddress()
		if err2 != nil {
			err = err2
			return
		}
		csRes.SubscriberIp = ip
	} else {
		err = &gtpv2.RequiredIEMissingError{Type: ie.PDNAddressAllocation}
		return
	}

	// Pgw control plane fteid
	if pgwCFteidIE := csResGtp.PGWS5S8FTEIDC; pgwCFteidIE != nil {
		pgwCFteid, _, err2 := handleFTEID(pgwCFteidIE)
		if err2 != nil {
			err = fmt.Errorf("Couldn't get PGW control plane FTEID: %s ", err2)
			return
		}
		//session.AddTEID(interfaceType, pgwCFteid.GetTeid())
		csRes.CPgwFteid = pgwCFteid
	} else {
		err = &gtpv2.RequiredIEMissingError{Type: ie.FullyQualifiedTEID}
		return
	}

	// TODO: handle more than one bearer
	if brCtxIE := csResGtp.BearerContextsCreated; brCtxIE != nil {
		bearerCtx := &protos.BearerContext{}
		for _, childIE := range brCtxIE.ChildIEs {
			switch childIE.Type {
			case ie.Cause:
				cause, err2 := childIE.Cause()
				if err2 != nil {
					err = err2
					return
				}
				if cause != gtpv2.CauseRequestAccepted {
					err = &gtpv2.CauseNotOKError{
						MsgType: csResGtp.MessageTypeName(),
						Cause:   cause,
						//Msg:     fmt.Sprintf("subscriber: %s", session.IMSI),
					}
					return
				}
			case ie.EPSBearerID:
				ebi, err2 := childIE.EPSBearerID()
				if err2 != nil {
					err = err2
					return
				}
				bearerCtx.Id = uint32(ebi)
			case ie.FullyQualifiedTEID:
				uFteid, _, err2 := handleFTEID(childIE)
				if err2 != nil {
					err = err2
					return
				}
				bearerCtx.UserPlaneFteid = uFteid
			case ie.ChargingID:
				bearerCtx.ChargingId, err = childIE.ChargingID()
				if err != nil {
					return
				}
			}
		}
		csRes.BearerContext = bearerCtx
	} else {
		err = &gtpv2.RequiredIEMissingError{Type: ie.BearerContext}
		return
	}
	return csRes, nil
}

// parseDelteSessionResponse parses a gtp message into a DeleteSessionResponsePgw. In case
// the message is proper it also returns the session. In case it there is an error it returns
// the cause of error
func parseDelteSessionResponse(msg message.Message) (
	dsRes *protos.DeleteSessionResponsePgw, err error) {
	cdResGtp := msg.(*message.DeleteSessionResponse)
	glog.V(2).Infof("Received Delete Session Response (gtp):\n%s", cdResGtp.String())

	dsRes = &protos.DeleteSessionResponsePgw{}
	// check Cause value first.
	if causeIE := cdResGtp.Cause; causeIE != nil {
		cause, err2 := causeIE.Cause()
		if err2 != nil {
			err = fmt.Errorf("Couldn't check cause of delete session response: %s", err2)
			return
		}
		if cause != gtpv2.CauseRequestAccepted {
			err = &gtpv2.CauseNotOKError{
				MsgType: cdResGtp.MessageTypeName(),
				Cause:   cause,
				Msg:     fmt.Sprintf("Delete Session Response not accepted"),
			}
			return
		}
	} else {
		err = &gtpv2.RequiredIEMissingError{
			Type: ie.Cause,
		}
		return
	}
	return dsRes, nil
}

// handleFTEID converts FTEID IE format into Proto format returning also the type of interface
func handleFTEID(fteidIE *ie.IE) (*protos.Fteid, uint8, error) {
	interfaceType, err := fteidIE.InterfaceType()
	if err != nil {
		return nil, interfaceType, err
	}
	teid, err := fteidIE.TEID()
	if err != nil {
		return nil, interfaceType, err
	}

	fteid := &protos.Fteid{Teid: teid}

	if !fteidIE.HasIPv4() && !fteidIE.HasIPv6() {
		return nil, interfaceType, fmt.Errorf("Error: fteid %+v has no ips", fteidIE.String())
	}
	if fteidIE.HasIPv4() {
		ipv4, err := fteidIE.IPv4()
		if err != nil {
			return nil, interfaceType, err
		}
		fteid.Ipv4Address = ipv4.String()
	}
	if fteidIE.HasIPv6() {
		ipv6, err := fteidIE.IPv6()
		if err != nil {
			return nil, interfaceType, err
		}
		fteid.Ipv6Address = ipv6.String()
	}
	return fteid, interfaceType, nil
}
