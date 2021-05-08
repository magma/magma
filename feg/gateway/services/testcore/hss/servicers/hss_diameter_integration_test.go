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

package servicers_test

import (
	"fmt"
	"testing"
	"time"

	"magma/feg/gateway/diameter"
	s6a "magma/feg/gateway/services/s6a_proxy/servicers"
	swx "magma/feg/gateway/services/swx_proxy/servicers"
	hss "magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/servicers/test_utils"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/stretchr/testify/assert"
)

func TestHomeSubscriberServer_handleAIR(t *testing.T) {
	// Create a test client interface which expects a valid AIA response.
	clientHandler := func(conn diam.Conn, msg *diam.Message) {
		var aia s6a.AIA

		// Check that the AIA is a success and has the expected data.
		err := msg.Unmarshal(&aia)
		assert.NoError(t, err)
		assert.Equal(t, "magma;123_1234", aia.SessionID)
		assert.Equal(t, diam.Success, int(aia.ResultCode))
		assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginHost)
		assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginRealm)
	}

	air := createAIR("sub1")
	testDiameterMessage(t, clientHandler, air)
}

func TestHomeSubscriberServer_handleULA(t *testing.T) {
	// Create a test client interface which expects a valid ULA response.
	clientHandler := func(conn diam.Conn, msg *diam.Message) {
		var ula s6a.ULA

		// Check that the ULA is a success and has the expected data.
		err := msg.Unmarshal(&ula)
		assert.NoError(t, err)
		assert.Equal(t, "magma;123_1234", ula.SessionID)
		assert.Equal(t, diam.Success, int(ula.ResultCode))
		assert.Equal(t, datatype.DiameterIdentity("magma.com"), ula.OriginHost)
		assert.Equal(t, datatype.DiameterIdentity("magma.com"), ula.OriginRealm)
	}

	ulr := createULR("sub1")
	testDiameterMessage(t, clientHandler, ulr)
}

func TestHomeSubscriberServer_handleMAR(t *testing.T) {
	// Create a test client interface which expects a valid MAA response.
	clientHandler := func(conn diam.Conn, msg *diam.Message) {
		var maa swx.MAA

		// Check that the MAA is a success and has the expected data.
		err := msg.Unmarshal(&maa)
		assert.NoError(t, err)
		assert.Equal(t, "magma;123_1234", maa.SessionID)
		assert.Equal(t, diam.Success, int(maa.ResultCode))
		assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginHost)
		assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginRealm)
		checkSIPAuthVectors(t, maa, 1)
	}

	mar := createMARWithSingleAuthItem("sub1")
	testDiameterMessage(t, clientHandler, mar)
}

func TestHomeSubscriberServer_handleSAR(t *testing.T) {
	// Create a test client interface which expects a valid SAA response.
	clientHandler := func(conn diam.Conn, msg *diam.Message) {
		checkSAASuccess(t, msg)
	}

	sar := createSAR("sub1", swx.ServerAssignmentType_REGISTRATION)
	testDiameterMessage(t, clientHandler, sar)
}

// testDiameterMessage sends a message to a test diameter server and provides the
// response in a callback function.
// Inputs: t - test interface
//         clientHandler - receives responses from the server
//         msg - the message to send to the test server
func testDiameterMessage(t *testing.T, clientHandler diam.HandlerFunc, msg *diam.Message) {
	// Wrap the test client interface so we can signal that a message has been
	// received.
	signal := make(chan struct{})
	handler := func(conn diam.Conn, msg *diam.Message) {
		clientHandler(conn, msg)
		close(signal)
	}

	// Create a test client-server connection.
	conn, err := getConnectionToTestHSS(t, handler)
	assert.NoError(t, err)
	defer conn.Close()

	// Send the message.
	_, err = msg.WriteTo(conn)
	assert.NoError(t, err)

	// Wait until the client receives a message or we time out.
	select {
	case <-signal:
		// Received a message.
	case <-time.After(time.Second):
		assert.Fail(t, "service timed out before receiving a response")
	}
}

// getTestHSSDiameterServer returns a test home subscriber server with a
// running diameter server listening for new connections.
func getTestHSSDiameterServer(t *testing.T) *hss.HomeSubscriberServer {
	// Start s6a diameter server
	result := make(chan error)
	hss := test_utils.NewTestHomeSubscriberServer(t)
	started := make(chan string)
	go func() {
		err := hss.Start(started)
		if err != nil {
			fmt.Printf("getConnectionToTestHSS Error: %v for address: %s\n", err, hss.Config.Server.Address)
			result <- err
		}
	}()

	// Wait for the server to start up.
	select {
	case err := <-result:
		assert.Fail(t, "%v", err)
		return nil
	case localAddr := <-started:
		// Overwrite server address with actual local address in case, the test HSS was started
		// without a configured port # (for example: 127.0.0.1:0)
		hss.Config.Server.Address = localAddr
	}
	return hss
}

// getConnectionToTestHSS starts a new Test Home Subscriber Server on given network & address
// Inputs: The client handler function receives messages from the server
// Outputs: a diameter connection to the server or an error
func getConnectionToTestHSS(t *testing.T, clientHandler diam.HandlerFunc) (diam.Conn, error) {
	hss := getTestHSSDiameterServer(t)

	// Create a client to receive the server's messages.
	clientMux := sm.New(&sm.Settings{
		OriginHost:       "magma.com",
		OriginRealm:      "magma.com",
		VendorID:         datatype.Unsigned32(diameter.Vendor3GPP),
		ProductName:      "magma",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
	})
	clientMux.Handle("ALL", clientHandler) // Catch all.

	// Create a connection to the server.
	client := &sm.Client{
		Handler: clientMux,
		SupportedVendorID: []*diam.AVP{
			diam.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
		},
		VendorSpecificApplicationID: []*diam.AVP{
			diam.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_S6A_APP_ID)),
					diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
				},
			}),
		},
	}
	serverCfg := hss.Config.Server
	return client.DialNetwork(serverCfg.Protocol, serverCfg.Address)
}
