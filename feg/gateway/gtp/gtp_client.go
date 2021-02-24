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
	"time"

	"github.com/golang/glog"
	"github.com/wmnsk/go-gtp/gtpv2"
)

const (
	GTPC_PORT = 0 // if set to 0, port is set automatically
)

type Client struct {
	*gtpv2.Conn
	connType   uint8
	localAddr  *net.UDPAddr
	remoteAddr *net.UDPAddr
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
	localAddr := &net.UDPAddr{IP: localAddrIp, Port: GTPC_PORT, Zone: ""}

	return NewConnectedClient(ctx, localAddr, remoteAddr, connType)
}

// NewConnectedClient creates a GTP-C client and checks with an echo if remote Addrs is
// available. It also runs the GTP-C server waiting for incoming calls
func NewConnectedClient(ctx context.Context, localAddr, remoteAddr *net.UDPAddr, connType uint8) (*Client, error) {
	var err error
	c := NewClient(localAddr, remoteAddr, connType)
	c.Conn, err = gtpv2.Dial(ctx, localAddr, remoteAddr, connType, 0)
	if err != nil {
		return nil, fmt.Errorf("could not connect to GTP-C %s server: %s", remoteAddr.String(), err)
	}
	return c, nil
}

// NewRunningAutoClient creates a GTP-C client inding out automatically the local IP Address to
// be used to reach the remote IP.
// It DOES NOT send initial echo to check if the server is alive
// It runs the GTP-C server to serve incoming calls and responses.
func NewRunningAutoClient(ctx context.Context, remoteIPAndPortStr string, connType uint8) (*Client, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", remoteIPAndPortStr)
	if err != nil {
		return nil, fmt.Errorf("could not resolve remote address %s: %s", remoteIPAndPortStr, err)
	}
	localAddrIp, err := GetOutboundIP(remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("could not find local address automatically:  %s", err)
	}
	localAddr := &net.UDPAddr{IP: localAddrIp, Port: GTPC_PORT, Zone: ""}
	return NewRunningClient(ctx, localAddr, remoteAddr, connType)
}

// NewRunningClient creates a GTP-C client. It also runs the GTP-C server waiting for incomming calls
// If you need to check raddrs availability, use NewConnectedClient
func NewRunningClient(ctx context.Context, localAddr, remoteAddr *net.UDPAddr, connType uint8) (*Client, error) {
	c := NewClient(localAddr, remoteAddr, connType)
	c.Enable()
	err := c.Run(ctx)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// NewClient creates basic configuration structure for a GTP-C client. It does
// not starts any connection or server.
func NewClient(localAddr, remoteAddr *net.UDPAddr, connType uint8) *Client {
	return &Client{
		connType:   connType,
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
	}
}

// Enable just creates the object connection enabling messages to be sent
func (c *Client) Enable() {
	c.Conn = gtpv2.NewConn(c.localAddr, c.connType, 0)
}

// Run launches the actual GTP-C cluent which will be able to send and receive GTP-C messages
func (c *Client) Run(ctx context.Context) error {
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
	return nil
}

func (c *Client) GetServerAddress() *net.UDPAddr {
	return c.remoteAddr
}

func (c *Client) GetLocalAddress() *net.UDPAddr {
	return c.localAddr
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
