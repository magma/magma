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
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/wmnsk/go-gtp/gtpv2"
)

const (
	ANY_IP         = "0.0.0.0"
	GTPC_AUTO_PORT = 0 // if set to 0, port is set automatically
)

type Client struct {
	*gtpv2.Conn
	connType uint8
}

// NewRunningClient creates a GTP-C client. It also runs the GTP-C server waiting for incomming calls
// localIpAndPort is in form ip:port  (127.0.0.1:1)
// 	- In case localIpAndPort is empty it uses any IP and a random port
// 	- In case ip is not provided ( :port, or 0.0.0.0:port) it uses any interface
// 	- In case port is set to 0 it uses a random port ( 0.0.0.0:0, or 10.0.0.1:0)
// If you need to check server availability before any connection, use NewConnectedClient
func NewRunningClient(ctx context.Context, localIpAndPort string, connType uint8) (*Client, error) {
	if localIpAndPort == "" {
		localIpAndPort = fmt.Sprintf("%s:%d", ANY_IP, GTPC_AUTO_PORT)
	}
	splitted := strings.Split(localIpAndPort, ":")
	if len(splitted) != 2 {
		return nil, fmt.Errorf("LocalIpAndPort must be formatted as IP:Port, but %s was received", localIpAndPort)
	}
	ip := splitted[0]
	if ip == "" {
		ip = ANY_IP
	}
	port, err := strconv.Atoi(splitted[1])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse GTP port: %s", err)
	}
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		return nil, fmt.Errorf("Failed to parse IP address: %s from %s", ip, localIpAndPort)
	}

	localAddr := &net.UDPAddr{IP: ipAddr, Port: port, Zone: ""}
	c := newClient(localAddr, connType)
	c.enable(localAddr)
	err = c.run(ctx)
	if err != nil {
		return nil, err
	}
	c.WaitUntilClientIsReady(0)
	c.DisableValidation()
	return c, nil
}

// NewConnectedAutoClient creates a GTP client finding out automatically the local IP Address to
// be used to reach the remote IP.
// It checks if remote end is alive using echo.
// It runs the GTP-C server to serve incoming calls and responses.
func NewConnectedAutoClient(ctx context.Context, remoteIPAndPortStr string, connType uint8) (*Client, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", remoteIPAndPortStr)
	if err != nil {
		return nil, fmt.Errorf("could not resolve remote address %s: %s", remoteIPAndPortStr, err)
	}
	localAddrIp, err := GetOutboundIP(remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("could not find local address automatically:  %s", err)
	}
	localAddr := &net.UDPAddr{IP: localAddrIp, Port: GTPC_AUTO_PORT, Zone: ""}

	return NewConnectedClient(ctx, localAddr, remoteAddr, connType)
}

// NewConnectedClient creates a GTP-C client and checks with an echo if remote Addrs is
// available. It also runs the GTP-C server waiting for incoming calls
func NewConnectedClient(ctx context.Context, localAddr, remoteAddr *net.UDPAddr, connType uint8) (*Client, error) {
	var err error
	c := newClient(localAddr, connType)
	c.Conn, err = gtpv2.Dial(ctx, localAddr, remoteAddr, connType, 0)
	if err != nil {
		return nil, fmt.Errorf("could not connect to GTP-C %s server: %s", remoteAddr.String(), err)
	}
	c.DisableValidation()
	return c, nil
}

// NewClient creates basic configuration structure for a GTP-C client. It does
// not starts any connection or server.
func newClient(localAddr *net.UDPAddr, connType uint8) *Client {
	cli := &Client{
		connType: connType,
	}
	return cli
}

// Enable just creates the object connection enabling messages to be sent
func (c *Client) enable(localAddr *net.UDPAddr) {
	c.Conn = gtpv2.NewConn(localAddr, c.connType, 0)
}

// Run launches the actual GTP-C cluent which will be able to send and receive GTP-C messages
func (c *Client) run(ctx context.Context) error {
	if c.Conn == nil {
		return fmt.Errorf("nil conn object. You may need to Enable the client first")
	}
	go func() {
		if ctx == nil {
			ctx = context.Background()
		}
		if err := c.ListenAndServe(ctx); err != nil {
			glog.Errorf("error running gtp server: %s", err)
			return
		}
	}()
	//TODO: remove this wait once there is a way to check when the listener is ready
	return nil
}

func (c *Client) GetSessionAndCTeidByIMSI(imsi string) (*gtpv2.Session, uint32, error) {
	session, err := c.GetSessionByIMSI(imsi)
	if err != nil {
		glog.Errorf("Couldnt delete session. Couldnt find a session for IMSI %s:, %s", imsi, err)
		return nil, 0, err
	}
	teid, err := session.GetTEID(gtpv2.IFTypeS5S8PGWGTPC)
	if err != nil {
		glog.Errorf("Couldnt delete session. Couldnt find control TEID for IMSI %s:, %s", imsi, err)
		return nil, 0, err
	}
	return session, teid, nil
}

// Get preferred outbound ip of this machine
func GetOutboundIP(testIp *net.UDPAddr) (net.IP, error) {
	connection, err := net.Dial("udp", testIp.String())
	if err != nil {
		return nil, err
	}

	defer connection.Close()
	localAddr := connection.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

//TODO: remove this once we find a way to safely wait for initialization of the service
// WaitUntilClientIsReady is a hack to know when the client is ready and avoid null pointer issues using
// the GTP-C client too early. Since go-gtp doesn't offer any visibuility on the readines of the connection
// we use LocalAddrs as indicator
func (c *Client) WaitUntilClientIsReady(count int) {

	// TODO: only use those 3 waits for debugging
	time.Sleep(time.Millisecond * 20)
	time.Sleep(time.Millisecond * 20)
	time.Sleep(time.Millisecond * 20)

	defer func() {

		if count > 50 {
			time.Sleep(time.Millisecond * 20)
		}
		if count > 100 {
			glog.Errorf("Couldnt start GTP-Client")
			return
		}
		if r := recover(); r != nil {
			fmt.Print(".")
			c.WaitUntilClientIsReady(count + 1)
		}
	}()
	// this call will panic while client is starting
	addr := c.LocalAddr().String()
	fmt.Println()
	glog.V(2).Infof("Started GTP-C client in %s", addr)
}
