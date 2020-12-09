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
	"fmt"
	"net"

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

// NewConnectedAutoClient creates a GTP-C client finding out automatically the local Address that needs to
// be used to reach the remote IP. It checks remote address using echo. And It also runs the GTP-C
// server waiting for incoming calls.
func NewConnectedAutoClient(ctx context.Context, remoteAddrStr string, connType uint8) (*Client, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", remoteAddrStr)
	if err != nil {
		return nil, fmt.Errorf("could not resolve remote address: %s", err)
	}
	localAddrIp, err := GetOutboundIP(remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("Could not find automatically local address:  %s", err)
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
		return nil, fmt.Errorf("Could not connect to GTP-C %s server:  %s", remoteAddr.String(), err)
	}
	return c, nil
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

func (c *Client) Enable() {
	c.Conn = gtpv2.NewConn(c.localAddr, c.connType, 0)
}

func (c *Client) Run(ctx context.Context) error {
	if c.Conn == nil {
		return fmt.Errorf("nil conn object. You may need to Enable the client first")
	}
	go func() {
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
	//Conn, err := net.Dial("udp", "8.8.8.8:80")
	connection, err := net.Dial("udp", testIp.String())
	if err != nil {
		return nil, fmt.Errorf("Couldn't local io running in this server: %s", err)
	}
	defer connection.Close()
	localAddr := connection.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
