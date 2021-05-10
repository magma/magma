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

// gtp module wraps gp-gtp v2 client providing the client with some custom functions
// to ease its instantiation

package gtp

import (
	"fmt"
	"net"

	"magma/feg/gateway/gtp/enriched_message"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

type GtpParserFunc = func(msg message.Message) (csRes proto.Message, err error)

// SendMessageAndExtractGrpc Send the message and blocks for GTP response with that specific sequence number
// It times out after GtpTimeout seconds
func (c *Client) SendMessageAndExtractGrpc(imsi string, srcTEID uint32, peerAddr net.Addr, msg message.Message) (
	proto.Message, error) {
	// Receive Create Session Response
	session := c.getSessionOrCreateNew(imsi, srcTEID, peerAddr)
	// session will be removed once we are done. No need to keep it
	defer c.RemoveSession(session)

	sequence, err := c.SendMessageTo(msg, session.PeerAddr())
	if err != nil {
		return nil, err
	}

	incomingMsg, err := session.WaitMessage(sequence, c.GtpTimeout)
	if err != nil {
		return nil, err
	}
	grpcMsg, err := enriched_message.ExtractGrpcMessageFromGtpMessage(incomingMsg)
	if err != nil {
		return nil, fmt.Errorf("GTP server return an error: %s", err)
	}
	return grpcMsg, nil
}

// getSessionOrCreateNew it used on sessionless mode before sending a message. That function will retrieve a session
// for that imsi if exist, otherwsie will create a new one. This way, we can send any kind of message no matter if
// a previous session existed
func (c *Client) getSessionOrCreateNew(imsi string, srcTEID uint32, peerAddr net.Addr) *gtpv2.Session {
	session, err := c.GetSessionByIMSI(imsi)
	if err == nil {
		return session
	}
	session = gtpv2.NewSession(peerAddr, &gtpv2.Subscriber{Location: &gtpv2.Location{}})
	session.IMSI = imsi
	session.Activate()
	c.RegisterSession(srcTEID, session)
	return session
}

func (c *Client) RemoveSessionByIMSI(imsi string) {
	session, err := c.GetSessionByIMSI(imsi)
	if err != nil {
		glog.Warningf("Couldn't remove sessiong for imsi %s because it was not found", imsi)
		return
	}
	c.RemoveSession(session)
}

// PassMessage will send and enriched_message to the session with that specific teid
// If there were an error during parsing, enriched_message will contain that error
// If session can not be found, then the caller will never receive an answer and will time out
func (c *Client) PassMessage(teid uint32, senderAddr net.Addr,
	gtpMessage message.Message, grpcMessage proto.Message, incomingError error) error {
	session, err := c.GetSessionByTEID(teid, senderAddr)
	if err != nil {
		return err
	}
	enrichedMsg := enriched_message.NewMessageWithGrpc(gtpMessage, grpcMessage, incomingError)
	// pass message to same session
	if err = gtpv2.PassMessageTo(session, enrichedMsg, c.GtpTimeout); err != nil {
		return err
	}
	return nil
}
