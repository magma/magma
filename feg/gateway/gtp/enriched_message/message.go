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

package enriched_message

import (
	"fmt"

	proto "github.com/golang/protobuf/proto"
	"github.com/wmnsk/go-gtp/gtpv2/message"
)

// MessageWithGrpc wraps Message interface so we can use it as a Message.
// grpcMessage field to store the GRPC version of the Message
// err stores any possible error that were associated with parsing message into grpcMessage
type MessageWithGrpc struct {
	message.Message               // GTP
	grpcMessage     proto.Message // GRPC
	err             error
}

// NewMessageWithGrpc returns a full valid MessageWithGrpc which include all parameters
func NewMessageWithGrpc(gtpMessage message.Message, grpcMessage proto.Message, err error) *MessageWithGrpc {
	return &MessageWithGrpc{
		Message:     gtpMessage,
		grpcMessage: grpcMessage,
		err:         err,
	}
}

func (m MessageWithGrpc) GetGrpcMessage() proto.Message {
	return m.grpcMessage
}

func ExtractGrpcMessageFromGtpMessage(incomingMsg message.Message) (proto.Message, error) {
	// check if it is NewMessageWithGrpc
	var withGrpc *MessageWithGrpc
	switch m := incomingMsg.(type) {
	case *MessageWithGrpc:
		withGrpc = m
	default:
		return nil, fmt.Errorf("incoming message it is not MessageWithGrpc type %+v", incomingMsg)
	}
	if withGrpc.err != nil {
		return nil, withGrpc.err
	}
	grpcMessage := withGrpc.GetGrpcMessage()
	return grpcMessage, nil
}
