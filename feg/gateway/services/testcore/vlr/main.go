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

package main

import (
	"flag"
	"log"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	csfbServicers "magma/feg/gateway/services/csfb/servicers"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/decode/message"
	"magma/feg/gateway/services/testcore/vlr/servicers"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func main() {
	srv, err := service.NewOrchestratorService(registry.ModuleName, registry.MOCK_VLR)
	if err != nil {
		log.Fatalf("Error creating mock VLR service: %s", err)
	}

	conn, err := csfbServicers.NewSCTPServerConnection()
	if err != nil {
		log.Fatalf("Failed to create SCTP connection for mocked VLR server: %s", err)
	}

	servicer, err := servicers.NewVLRServer(conn)
	if err != nil {
		log.Fatalf("Failed to create mocked VLR server: %s", err)
	}

	protos.RegisterMockCoreConfiguratorServer(srv.GrpcServer, servicer)

	_, err = servicer.Conn.StartListener(
		csfbServicers.LocalIPAddress,
		csfbServicers.LocalPort,
	)
	defer servicer.Conn.CloseListener()
	if err != nil {
		log.Fatalf("Failed to start SCTP listener on mocked VLR server: %s", err)
	}

	go func() {
		for {
			if !servicer.Conn.ConnectionEstablished() {
				err = servicer.Conn.AcceptConn()
				if err != nil {
					glog.Errorf("Failed to accept connection from a client: %s", err)
					continue
				}
			}
			// receive messages
			glog.V(2).Info("Wait for incoming messages.")
			receivedMsg, err := servicer.Conn.ReceiveThroughListener()
			if err != nil {
				glog.Errorf("Failed to receive message: %s", err)
				servicer.Conn.CloseConn()
				continue
			}
			msgType, _, err := message.SGsMessageDecoder(receivedMsg)
			if err != nil {
				glog.Errorf("Failed to decode message: %s", err)
				// pop one element of the queue
				servicer.GetNextRequestReplyPair()
				continue
			}
			glog.V(2).Infof("Received %s", decode.MsgTypeNameByCode[msgType])
			err = servicer.ReplyClient(receivedMsg)
			if err != nil {
				glog.Errorf("Failed to reply to the client: %s", err)
			}
		}
	}()

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running mock VLR service: %s", err)
	}
}
