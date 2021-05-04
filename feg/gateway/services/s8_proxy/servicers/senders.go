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

// senders contains the calls that will be run when a GTP command is sent.
// Those functions will also return the result of the call

package servicers

import (
	"fmt"
	"net"

	"magma/feg/cloud/go/protos"

	"github.com/golang/glog"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

// sendAndReceiveCreateSession creates a session in the gtp client, sends the create session request
// to PGW and waits for its answers.
// Returns a GRPC message translated from the GTP-C create session response
func (s *S8Proxy) sendAndReceiveCreateSession(
	csReq *protos.CreateSessionRequestPgw,
	cPgwUDPAddr *net.UDPAddr,
	csReqMsg message.Message) (*protos.CreateSessionResponsePgw, error) {
	glog.V(2).Infof("Send Create Session Request (grpc) to %s:\n%s",
		cPgwUDPAddr.String(), csReq.String())
	glog.V(2).Infof("Send Create Session Request (gtp) to %s:\n%s",
		cPgwUDPAddr.String(), csReqMsg.(*message.CreateSessionRequest).String())
	//glog.V(4).Infof("Send Create Session Request (gtp) to %s:\n%s",
	//	cPgwUDPAddr.String(), message.Prettify(csReqMsg))

	grpcMessage, err := s.gtpClient.SendMessageAndExtractGrpc(csReq.Imsi, csReq.CAgwTeid, cPgwUDPAddr, csReqMsg)
	if err != nil {
		return nil, fmt.Errorf("no response message to CreateSessionRequest: %s", err)
	}
	// check if message is proper
	csRes, ok := grpcMessage.(*protos.CreateSessionResponsePgw)
	if !ok {
		s.gtpClient.RemoveSessionByIMSI(csReq.Imsi)
		return nil, fmt.Errorf("Wrong response type (no CreateSessionResponse), maybe received out of order response message: %s", err)
	}
	glog.V(2).Infof("Create Session Response (grpc):\n%s", csRes.String())
	return csRes, nil
}

// sendAndReceiveDeleteSession  sends delete session request GTP-C message to PGW and
// waits for its answers.
// Returns a GRPC message translated from the GTP-C delete session response
func (s *S8Proxy) sendAndReceiveDeleteSession(req *protos.DeleteSessionRequestPgw,
	cPgwUDPAddr *net.UDPAddr,
	dsReqMsg message.Message) (*protos.DeleteSessionResponsePgw, error) {
	glog.V(2).Infof("Send Delete Session Request (grpc) to %s:\n%s", cPgwUDPAddr,
		dsReqMsg)
	glog.V(2).Infof("Send Delete Session Request (gtp) to %s:\n%s",
		cPgwUDPAddr.String(), dsReqMsg.(*message.DeleteSessionRequest).String())
	//glog.V(4).Infof("Send Delete Session Request (gtp) to %s:\n%s",
	//	cPgwUDPAddr.String(), message.Prettify(dsReqMsg))
	grpcMessage, err := s.gtpClient.SendMessageAndExtractGrpc(req.Imsi, req.CAgwTeid, cPgwUDPAddr, dsReqMsg)
	if err != nil {
		return nil, fmt.Errorf("no response message to DeleteSessionRequest: %s", err)
	}
	dsRes, ok := grpcMessage.(*protos.DeleteSessionResponsePgw)
	if !ok {
		return nil, fmt.Errorf("Wrong response type (no DeleteSessionResponse), maybe received out of order response message: %s", err)
	}
	glog.V(2).Infof("Delete Session Response (grpc):\n%s", dsRes.String())
	return dsRes, err
}
