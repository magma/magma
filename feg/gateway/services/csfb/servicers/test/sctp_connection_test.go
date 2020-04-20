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

package test

import (
	"testing"
	"time"
	"unsafe"

	"magma/feg/gateway/services/csfb/servicers"
	"magma/feg/gateway/services/csfb/test_init"

	"github.com/ishidawataru/sctp"
	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	// Mock the server which is receiving message
	ln, port := test_init.GetMockVLRListenerAndPort(t)
	defer ln.Close()

	msgReceivedFlag := make(chan bool)
	go func() {
		netConn, err := ln.Accept()
		assert.NoError(t, err)
		defer netConn.Close()

		wconn := sctp.NewSCTPSndRcvInfoWrappedConn(netConn.(*sctp.SCTPConn))
		defer wconn.Close()

		buf := make([]byte, 254)
		n, err := wconn.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, buf[unsafe.Sizeof(sctp.SndRcvInfo{}):n], []byte("hello"))

		msgReceivedFlag <- true
	}()

	// Send message to the mocked server
	vlrSCTPAddr := servicers.ConstructSCTPAddr(
		servicers.DefaultVLRIPAddress,
		port,
	)
	vlrConn, err := servicers.NewSCTPClientConnection(vlrSCTPAddr, nil)
	assert.NoError(t, err)

	err = vlrConn.EstablishConn()
	defer vlrConn.CloseConn()
	assert.NoError(t, err)

	err = vlrConn.Send([]byte("hello"))
	assert.NoError(t, err)

	<-msgReceivedFlag
}

func TestReceive(t *testing.T) {
	// Mock the server which is sending message
	ln, port := test_init.GetMockVLRListenerAndPort(t)
	defer ln.Close()

	goRoutineEndFlag := make(chan bool)
	go func() {
		netConn, err := ln.Accept()
		assert.NoError(t, err)
		defer netConn.Close()

		sendConn := netConn.(*sctp.SCTPConn)
		defer sendConn.Close()

		ppid := 0
		info := &sctp.SndRcvInfo{
			Stream: uint16(ppid),
			PPID:   uint32(ppid),
		}

		sendConn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
		_, err = sendConn.SCTPWrite([]byte("hello hello"), info)
		assert.NoError(t, err)

		goRoutineEndFlag <- true
	}()

	// Receive message from the mocked server
	vlrSCTPAddr := servicers.ConstructSCTPAddr(
		servicers.DefaultVLRIPAddress,
		port,
	)
	vlrConn, err := servicers.NewSCTPClientConnection(vlrSCTPAddr, nil)
	assert.NoError(t, err)

	err = vlrConn.EstablishConn()
	defer vlrConn.CloseConn()
	assert.NoError(t, err)

	msg, err := vlrConn.Receive()
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello hello"), msg)

	<-goRoutineEndFlag
}

func TestServerReceiveAndReply(t *testing.T) {
	srv, err := servicers.NewSCTPServerConnection()
	assert.NoError(t, err)

	portNumber, err := srv.StartListener(
		servicers.LocalIPAddress,
		servicers.LocalPort,
	)
	defer srv.CloseListener()
	assert.NoError(t, err)

	// Client sends and receives message
	clientReadyToReceive := make(chan bool)
	go func() {
		vlrSCTPAddr := servicers.ConstructSCTPAddr(
			servicers.DefaultVLRIPAddress,
			portNumber,
		)
		vlrConn, err := servicers.NewSCTPClientConnection(vlrSCTPAddr, nil)
		assert.NoError(t, err)

		err = vlrConn.EstablishConn()
		defer vlrConn.CloseConn()
		assert.NoError(t, err)

		err = vlrConn.Send([]byte("12345"))
		assert.NoError(t, err)

		clientReadyToReceive <- true
		serverMsg, err := vlrConn.Receive()
		assert.NoError(t, err)
		assert.Equal(t, []byte("54321"), serverMsg)

		err = vlrConn.Send([]byte("123456"))
		assert.NoError(t, err)
	}()

	// Server accepts the connection from client
	err = srv.AcceptConn()
	assert.NoError(t, err)

	// Server receives
	clientMsg, err := srv.ReceiveThroughListener()
	assert.NoError(t, err)
	assert.Equal(t, []byte("12345"), clientMsg)

	// Use channel variable and Sleep() to ensure that the client
	// is ready to receive (i.e. blocked at Receive())
	<-clientReadyToReceive
	time.Sleep(time.Second * 1)

	// Server replies
	err = srv.SendFromServer([]byte("54321"))
	assert.NoError(t, err)

	// Server receives again
	clientMsg, err = srv.ReceiveThroughListener()
	assert.NoError(t, err)
	assert.Equal(t, []byte("123456"), clientMsg)

	err = srv.CloseListener()
	assert.NoError(t, err)

	err = srv.CloseConn()
	assert.NoError(t, err)
}
