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

func (s *s8Proxy) sendAndReceiveCreateSession(csReqIE []*ie.IE, sessionTeids SessionFTeids) (*protos.CreateSessionResponsePgw, error) {
	// Send Create Session Req
	session, seq, err := s.gtpClient.CreateSession(s.gtpClient.GetServerAddress(), csReqIE...)
	glog.V(2).Infof("Send Create Session Request:\n%s",
		message.NewCreateSessionRequest(0, 0, csReqIE...).String())
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %s", err)
	}

	// add TEID to session and register session
	session.AddTEID(sessionTeids.uFTeid.MustInterfaceType(), sessionTeids.uFTeid.MustTEID())
	s.gtpClient.RegisterSession(sessionTeids.cFTeid.MustTEID(), session)

	grpcMessage, err := waitMessageAndExtractGrpc(session, seq)
	if err != nil {
		//TODO: remove session properly
		s.gtpClient.RemoveSession(session)
		return nil, fmt.Errorf("no response message: %s", err)
	}

	// check if message is proper
	csRes, ok := grpcMessage.(*protos.CreateSessionResponsePgw)
	if !ok {
		//TODO handle  error case (remove session properly)
		s.gtpClient.RemoveSession(session)
		return nil, fmt.Errorf("Wrong response type, maybe received out of order response message: %s", err)
	}
	return csRes, nil
}

func waitMessageAndExtractGrpc(session *gtpv2.Session, sequence uint32) (proto.Message, error) {
	// Receive Create Session Response
	incomingMsg, err := session.WaitMessage(sequence, GtpTimeout)
	if err != nil {
		return nil, err
	}
	grpcMessage, err := enriched_message.ExtractGrpcMessageFromGtpMessage(incomingMsg)
	if err != nil {
		return nil, err
	}
	return grpcMessage, nil
}
