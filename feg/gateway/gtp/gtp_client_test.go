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

package gtp

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wmnsk/go-gtp/gtpv2"
	"github.com/wmnsk/go-gtp/gtpv2/ie"
)

const (
	sgwAddrs  = ":0"
	pgwAddrs  = "127.0.0.1:0" //port 0 means go will choose. Selected port will be injected on getDefaultConfig
	IMSI1     = "123456789012345"
	bearerId1 = uint8(5)
	bearerId2 = uint8(6)

	qci1 = uint8(5)
	qci2 = uint8(6)
)

func TestEcho(t *testing.T) {
	// run GTP server (PGW)
	pgwCli := startGTPServer(t)
	actualServerIPAndPort := pgwCli.LocalAddr().String()

	// run GTP client (SGW) and send echo message to check if server is available
	_, err := NewConnectedAutoClient(context.Background(), actualServerIPAndPort, gtpv2.IFTypeS5S8SGWGTPC)

	// if no error service was started and echo was received properly
	assert.NoError(t, err)
}

// TestReCreateSession tests the creation of new session over an existing session with the same IMSI
func TestGtpClient(t *testing.T) {
	// run GTP server (PGW)
	pgwCli := startGTPServer(t)
	actualServerIPAndPort := pgwCli.LocalAddr().String()

	// run GTP client (SGW) but do not attach to any server
	cli, err := NewRunningClient(context.Background(), sgwAddrs, gtpv2.IFTypeS5S8SGWGTPC)
	assert.NoError(t, err)

	remoteAddr, err := net.ResolveUDPAddr("udp", actualServerIPAndPort)
	assert.NoError(t, err)

	// find out the local interface to be used (because it is not specified for testing in sgwAddrs)
	localIP, err := GetOutboundIP(remoteAddr)
	assert.NoError(t, err)

	// test CreateSession
	csr := getCreateSessionRequest(t, cli, localIP, actualServerIPAndPort, bearerId1, qci1)
	expectedSession, _, err := cli.CreateSession(remoteAddr, csr...)
	assert.NoError(t, err)

	// tesat GetSessionAndCTeidByIMSI
	session, cteid, err := cli.GetSessionAndCTeidByIMSI(IMSI1)
	assert.NoError(t, err)
	assert.Equal(t, expectedSession, session)
	assert.Equal(t, uint32(0x0), cteid)
	assert.Equal(t, bearerId1, session.GetDefaultBearer().EBI)
	assert.Equal(t, qci1, session.GetDefaultBearer().QCI)

	// create same session with differnt QCI (old session should be removed)
	csr = getCreateSessionRequest(t, cli, localIP, actualServerIPAndPort, bearerId1, qci2)
	session, _, err = cli.CreateSession(remoteAddr, csr...)
	assert.NoError(t, err)
	cteid, err = session.GetTEID(gtpv2.IFTypeS5S8PGWGTPC)
	assert.Equal(t, uint32(0x0), cteid)
	assert.Equal(t, bearerId1, session.GetDefaultBearer().EBI)
	assert.Equal(t, qci2, session.GetDefaultBearer().QCI)
}

func startGTPServer(t *testing.T) *Client {
	pgwConn, err := NewRunningClient(context.Background(), pgwAddrs, gtpv2.IFTypeS5S8PGWGTPC)
	assert.NoError(t, err)
	return pgwConn
}

func getCreateSessionRequest(t *testing.T, cli *Client, laddrs net.IP, raddrs string, bearerId, qci uint8) []*ie.IE {
	// sgw will chose random teid on sgw side
	cSgwFTeid := cli.NewSenderFTEID(laddrs.String(), "")
	// pgw TEID will be 0. PGW will select one for sgw in the response
	cPgwFTeid := ie.NewFullyQualifiedTEID(gtpv2.IFTypeS5S8PGWGTPC, 0, raddrs, "").WithInstance(1)

	return []*ie.IE{
		ie.NewIMSI(IMSI1),
		ie.NewMSISDN("8130900000005"),
		ie.NewMobileEquipmentIdentity("123456780000015"),
		ie.NewUserLocationInformation(
			0, 0, 0, 1, 1, 0, 0, 0,
			"123", "456", 0, 0, 0, 0, 1, 1, 0, 0,
		),
		ie.NewRATType(gtpv2.RATTypeEUTRAN),
		ie.NewIndicationFromOctets(0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00),
		cSgwFTeid,
		cPgwFTeid,
		ie.NewAccessPointName("internet"),
		ie.NewSelectionMode(gtpv2.SelectionModeMSorNetworkProvidedAPNSubscribedVerified),
		ie.NewPDNType(gtpv2.PDNTypeIPv4),
		ie.NewAPNRestriction(gtpv2.APNRestrictionNoExistingContextsorRestriction),
		ie.NewAggregateMaximumBitRate(0, 0),
		ie.NewBearerContext(
			ie.NewEPSBearerID(bearerId),
			ie.NewBearerQoS(1, 2, 3, qci, 40000, 4001, 4000, 4001),
		),
	}
}
