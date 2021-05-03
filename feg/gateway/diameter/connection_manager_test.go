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

// Package diameter_test tests diameter calls within the magma setting
package diameter

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
)

// TestGyClient tests CCR init, update, and terminate messages using a fake
// server
func TestConnectionManager(t *testing.T) {
	const (
		ServerHost  = datatype.DiameterIdentity("test.test.com")
		ServerRealm = datatype.DiameterIdentity("test.com")
		ClientHost  = datatype.DiameterIdentity("test.magma.com")
		ClientRealm = datatype.DiameterIdentity("magma.com")
	)

	var (
		serverConfig = &DiameterServerConfig{DiameterServerConnConfig: DiameterServerConnConfig{
			Addr:     "127.0.0.1:0", // Addr will be updated by startTestServer to reflect assigned port
			Protocol: "tcp"},
		}
		serverConnMan = NewConnectionManager()
		mux           = sm.New(&sm.Settings{
			OriginHost:  ClientHost,
			OriginRealm: ClientRealm,
			VendorID:    datatype.Unsigned32(Vendor3GPP),
			ProductName: "connection manager",
		})
		cli = &sm.Client{
			Dict:               dict.Default,
			Handler:            mux,
			MaxRetransmits:     3,
			RetransmitInterval: time.Second,
			EnableWatchdog:     true,
			WatchdogInterval:   5 * time.Second,
			AuthApplicationID: []*diam.AVP{
				diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.CHARGING_CONTROL_APP_ID)),
			},
		}
		serverSmCli = &sm.Client{
			Dict:               dict.Default,
			Handler:            mux,
			MaxRetransmits:     3,
			RetransmitInterval: time.Second,
			EnableWatchdog:     true,
			WatchdogInterval:   5 * time.Second,
			AuthApplicationID: []*diam.AVP{
				diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.CHARGING_CONTROL_APP_ID)),
			},
		}
		testServerMux = sm.New(&sm.Settings{
			OriginHost:  ServerHost,
			OriginRealm: ServerRealm,
			VendorID:    datatype.Unsigned32(Vendor3GPP),
			ProductName: "hello",
		})

		serverStarted = make(chan struct{})

		newMessage = func() *diam.Message {
			m := diam.NewRequest(diam.CreditControl, diam.CHARGING_CONTROL_APP_ID, nil)
			m.NewAVP(avp.OriginHost, avp.Mbit, 0, ClientHost)
			m.NewAVP(avp.OriginRealm, avp.Mbit, 0, ClientRealm)
			return m
		}

		startTestServer = func(server *DiameterServerConfig, started chan struct{}) {
			srv := &diam.Server{Network: server.Protocol, Addr: server.Addr, Handler: testServerMux, Dict: nil}
			l, err := diam.MultistreamListen(server.Protocol, server.Addr)
			if err != nil {
				log.Fatalf("Could not create server socket on: %s, %v", server.Addr, err)
			}
			server.Addr = l.Addr().String()
			started <- struct{}{}
			err = srv.Serve(l)
			if err != nil {
				log.Fatalf("Could not start server, %s", err.Error())
				return
			}
		}

		verifyRequest = func(conn diam.Conn, message *diam.Message, server *DiameterServerConfig) error {
			host, _ := message.FindAVP(avp.DestinationHost, 0)
			realm, _ := message.FindAVP(avp.DestinationRealm, 0)
			if host == nil {
				return fmt.Errorf("Missing Destination-Host AVP; realm: %+v", realm)
			}
			if realm == nil {
				return fmt.Errorf("Missing Destination-Realm AVP; host: %+v", host)
			}
			return serverConnMan.AddExistingConnection(conn, serverSmCli, server)
		}
	)

	go startTestServer(serverConfig, serverStarted)
	<-serverStarted
	time.Sleep(time.Millisecond * 10)

	connMan := NewConnectionManager()
	conn, _ := connMan.GetConnection(cli, serverConfig)

	// assert multiple calls to GetConnection returns same object
	conn2, _ := connMan.GetConnection(cli, serverConfig)
	assert.Equal(t, conn, conn2)
	var err error

	// basic test

	handlerChan := make(chan error)
	testServerMux.HandleIdx(
		diam.CommandIndex{AppID: diam.CHARGING_CONTROL_APP_ID, Code: diam.CreditControl, Request: true},
		diam.HandlerFunc(func(conn diam.Conn, m *diam.Message) {
			handlerChan <- verifyRequest(conn, m, serverConfig)
		}))
	err = conn.SendRequest(newMessage(), 0)
	assert.NoError(t, err)
	select {
	case e := <-handlerChan:
		assert.NoError(t, e)
	case <-time.After(time.Second):
		t.Fatal("SendRequest timeout")
	}
	serverConn, ok := serverConnMan.connMap[serverConfig.DiameterServerConnConfig]
	assert.True(t, ok)
	assert.Equal(t, serverConn.conn.LocalAddr(), conn.conn.RemoteAddr())
	assert.Equal(t, serverConn.conn.RemoteAddr(), conn.conn.LocalAddr())

	// Test retries
	// On first call to Write, return error, next time, return nil
	c, _, err := conn.getDiamConnection()
	assert.NoError(t, err)
	c.Close() // should error out after that
	err = conn.SendRequest(newMessage(), 0)
	assert.Error(t, err)
	// destroyConnection was called - next send will succeed
	err = conn.SendRequest(newMessage(), 0)
	assert.NoError(t, err)
	select {
	case e := <-handlerChan:
		assert.NoError(t, e)
	case <-time.After(time.Second):
		t.Fatal("SendRequest2 timeout")
	}
	// Now, do it all in one send with retries
	c, _, err = conn.getDiamConnection()
	require.NoError(t, err)
	c.Close()
	err = conn.SendRequest(newMessage(), 1)
	assert.NoError(t, err)
	select {
	case e := <-handlerChan:
		assert.NoError(t, e)
	case <-time.After(time.Second):
		t.Fatal("SendRequest3 timeout")
	}

	connMan.DisableFor(time.Millisecond * 30)
	_, err = connMan.GetConnection(cli, serverConfig)
	assert.Error(t, err)
	time.Sleep(time.Millisecond * 100)
	_, err = connMan.GetConnection(cli, serverConfig)
	assert.NoError(t, err)
}

func TestEncodeDecodeSID(t *testing.T) {
	// Magma SIDs
	assert.Equal(t,
		EncodeSessionID("gx.magma.com", "IMSI123456789012345-987654321"),
		"gx.magma.com;9876;54321;IMSI123456789012345")
	assert.Equal(t,
		DecodeSessionID("gx.magma.com;9876;54321;IMSI123456789012345"),
		"IMSI123456789012345-987654321")

	assert.Equal(t,
		EncodeSessionID("gx.magma.com", "IMSI123456789012345-98"),
		"gx.magma.com;9;8;IMSI123456789012345")
	assert.Equal(t,
		DecodeSessionID("gx.magma.com;9;8;IMSI123456789012345"),
		"IMSI123456789012345-98")

	assert.Equal(t,
		EncodeSessionID("gx.magma.com", "IMSI123456789012345-1"),
		"gx.magma.com;;1;IMSI123456789012345")
	assert.Equal(t,
		DecodeSessionID("gx.magma.com;;1;IMSI123456789012345"),
		"IMSI123456789012345-1")

	// With Bearer ID
	assert.Equal(t,
		EncodeSessionID("gx.magma.com", "IMSI123456789012345_7-987654321"),
		"gx.magma.com;9876;54321;IMSI123456789012345_7")
	assert.Equal(t,
		DecodeSessionID("gx.magma.com;9876;54321;IMSI123456789012345_7"),
		"IMSI123456789012345_7-987654321")

	// Non magma SIDs
	assert.Equal(t,
		EncodeSessionID("gx.magma.com", "IMSI123456789012345"),
		"IMSI123456789012345")
	assert.Equal(t,
		EncodeSessionID("gx.magma.com", "123456789012345-987654321"),
		"123456789012345-987654321")
	assert.Equal(t, DecodeSessionID("IMSI123456789012345"), "IMSI123456789012345")
	assert.Equal(t,
		DecodeSessionID("gx.magma.com;987654321;987654321;123456789012345"),
		"gx.magma.com;987654321;987654321;123456789012345")

	// Parsing
	r := [5]string{}
	r[0], r[1], r[2], r[3], r[4] =
		ParseDiamSessionID("gx.magma.com;987654321;987654322;IMSI123456789012345_123")
	assert.Equal(t, r, [5]string{"gx.magma.com", "987654321", "987654322", "123456789012345", "123"})
	r[0], r[1], r[2], r[3], r[4] =
		ParseDiamSessionID("gx.magma.com;987654321;987654322;IMSI123456789012345")
	assert.Equal(t, r, [5]string{"gx.magma.com", "987654321", "987654322", "123456789012345", ""})
	r[0], r[1], r[2], r[3], r[4] = ParseDiamSessionID("blablabla")
	assert.Equal(t, r, [5]string{"blablabla", "", "", "", ""})
}
