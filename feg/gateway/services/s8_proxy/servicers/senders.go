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
	"github.com/golang/protobuf/proto"
	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/message"
	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/gtp/enriched_message"
	"time"

	"github.com/golang/glog"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
)

const (
	GtpTimeout = 5 * time.Second
)

// sendAndReceiveCreateSession creates a session in the gtp client, sends the create session request
// to PGW and waits for its answers.
// Returns a GRPC message translaged from the GTP-U create session response
func (s *S8Proxy) sendAndReceiveCreateSession(csReqIEs []*ie.IE, sessionTeids SessionFTeids) (*protos.CreateSessionResponsePgw, error) {
	glog.V(2).Infof("Send Create Session Request (gtp):\n%s",
		message.NewCreateSessionRequest(0, 0, csReqIEs...).String())

	session, seq, err := s.gtpClient.CreateSession(s.gtpClient.GetServerAddress(), csReqIEs...)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %s", err)
	}

	// add TEID to session and register session
	session.AddTEID(sessionTeids.uFTeid.MustInterfaceType(), sessionTeids.uFTeid.MustTEID())
	s.gtpClient.RegisterSession(sessionTeids.cFTeid.MustTEID(), session)

	grpcMessage, err := waitMessageAndExtractGrpc(session, seq)
	if err != nil {
		s.gtpClient.RemoveSession(session)
		return nil, fmt.Errorf("no response message to CreateSessionRequest: %s", err)
	}

	// check if message is proper
	csRes, ok := grpcMessage.(*protos.CreateSessionResponsePgw)
	if !ok {
		s.gtpClient.RemoveSession(session)
		return nil, fmt.Errorf("Wrong response type (no CreateSessionResponse), maybe received out of order response message: %s", err)
	}
	glog.V(2).Infof("Create Session Response (grpc):\n%s", csRes.String())
	return csRes, nil
}

// sendAndReceiveModifyBearer  sends modify bearer request GTP-U message to PGW and
// waits for its answers.
// Returns a GRPC message translaged from the GTP-U create session response
func (s *S8Proxy) sendAndReceiveModifyBearer(teid uint32, session *gtpv2.Session, mbReqIE []*ie.IE) (*protos.ModifyBearerResponsePgw, error) {
	glog.V(2).Infof("Send Modify Bearer Request (gtp):\n%s",
		message.NewModifyBearerRequest(teid, 0, mbReqIE...).String())

	seq, err := s.gtpClient.ModifyBearer(teid, session, mbReqIE...)
	if err != nil {
		return nil, err
	}
	grpcMessage, err := waitMessageAndExtractGrpc(session, seq)
	if err != nil {
		return nil, fmt.Errorf("no response message to ModifyBearerRequest: %s", err)
	}
	mbRes, ok := grpcMessage.(*protos.ModifyBearerResponsePgw)
	if !ok {
		return nil, fmt.Errorf("Wrong response type (no ModifyBearerResponse), maybe received out of order response message: %s", err)
	}
	glog.V(2).Infof("Modify Bearer Response (grpc):\n%s", mbRes.String())
	return mbRes, err
}

// sendAndReceiveDeleteSession  sends delete session request GTP-U message to PGW and
// waits for its answers.
// Returns a GRPC message translaged from the GTP-U create session response
func (s *S8Proxy) sendAndReceiveDeleteSession(teid uint32, session *gtpv2.Session) (*protos.DeleteSessionResponsePgw, error) {
	glog.V(2).Infof("Send Delete Session Request (gtp):\n%s",
		message.NewDeleteSessionRequest(teid, 0).String())

	seq, err := s.gtpClient.DeleteSession(teid, session)
	if err != nil {
		return nil, err
	}
	grpcMessage, err := waitMessageAndExtractGrpc(session, seq)
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

func (s *S8Proxy) sendAndReceiveEchoRequest() error {
	c := s.gtpClient.Conn
	_, err := c.EchoRequest(s.gtpClient.GetServerAddress())
	if err != nil {
		return err
	}
	return waitEchoResponse(s.echoChannel)
}

// waitMessageAndExtractGrpc blocks for GTP response with that specific sequence number
// It times out after GtpTimeout seconds
func waitMessageAndExtractGrpc(session *gtpv2.Session, sequence uint32) (proto.Message, error) {
	// Receive Create Session Response
	incomingMsg, err := session.WaitMessage(sequence, GtpTimeout)
	if err != nil {
		return nil, err
	}
	return enriched_message.ExtractGrpcMessageFromGtpMessage(incomingMsg)
}

func waitEchoResponse(ch chan error) error {
	select {
	case res := <-ch:
		return res
	case <-time.After(GtpTimeout):
		return fmt.Errorf("waitEchoResponse timeout")
	}

}
