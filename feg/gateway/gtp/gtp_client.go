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
	DefaultGtpTimeout     = 3 * time.Second
	SGWControlPlaneIfType = gtpv2.IFTypeS5S8SGWGTPC

	ANY_IP         = "0.0.0.0"
	GTPC_AUTO_PORT = 0 // if set to 0, port is set automatically
)

type Client struct {
	*gtpv2.Conn
	connType   uint8
	GtpTimeout time.Duration
}

// NewRunningClient creates a GTP-C client. It also runs the GTP-C server waiting for incomming calls
// localIpAndPort is in form ip:port  (127.0.0.1:1)
// 	- In case localIpAndPort is empty it uses any IP and a random port
// 	- In case ip is not provided ( :port, or 0.0.0.0:port) it uses any interface
// 	- In case port is set to 0 it uses a random port ( 0.0.0.0:0, or 10.0.0.1:0)
// If you need to check server availability before any connection, use NewConnectedClient
func NewRunningClient(ctx context.Context, localIpAndPort string, connType uint8, gtpTimeout time.Duration) (*Client, error) {
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
	c := newClient(localAddr, connType, gtpTimeout)
	err = c.run(ctx, localAddr)
	if err != nil {
		return nil, err
	}

	// We need to disable GTP validation in order to support stateless operation
	// Otherwise if we receive a message which doesnt have an actie session, the message
	// will be discarded.
	c.DisableValidation()
	return c, nil
}

// NewConnectedAutoClient creates a GTP client finding out automatically the local IP Address to
// be used to reach the remote IP.
// It checks if remote end is alive using echo.
// It runs the GTP-C server to serve incoming calls and responses.
func NewConnectedAutoClient(ctx context.Context, remoteIPAndPortStr string, connType uint8, gtpTimeout time.Duration) (*Client, error) {
	remoteAddr, err := net.ResolveUDPAddr("udp", remoteIPAndPortStr)
	if err != nil {
		return nil, fmt.Errorf("could not resolve remote address %s: %s", remoteIPAndPortStr, err)
	}
	localAddrIp, err := GetLocalOutboundIP(remoteAddr)
	if err != nil {
		return nil, fmt.Errorf("could not find local address automatically:  %s", err)
	}
	localAddr := &net.UDPAddr{IP: localAddrIp, Port: GTPC_AUTO_PORT, Zone: ""}

	return NewConnectedClient(ctx, localAddr, remoteAddr, connType, gtpTimeout)
}

// NewConnectedClient creates a GTP-C client and checks with an echo if remote Addrs is
// available. It also runs the GTP-C server waiting for incoming calls
func NewConnectedClient(ctx context.Context, localAddr, remoteAddr *net.UDPAddr, connType uint8, gtpTimeout time.Duration) (*Client, error) {
	var err error
	c := newClient(localAddr, connType, gtpTimeout)
	c.Conn, err = gtpv2.Dial(ctx, localAddr, remoteAddr, connType, 0)
	if err != nil {
		return nil, fmt.Errorf("could not connect to GTP-C %s server: %s", remoteAddr.String(), err)
	}
	c.DisableValidation()
	return c, nil
}

// NewClient creates basic configuration structure for a GTP-C client. It does
// not starts any connection or server.
func newClient(localAddr *net.UDPAddr, connType uint8, gtpTimeout time.Duration) *Client {
	cli := &Client{
		connType:   connType,
		GtpTimeout: configOrDefaultTimeout(gtpTimeout),
	}
	return cli
}

// run starts the listener and launches the actual GTP-C routine
func (c *Client) run(ctx context.Context, localAddr *net.UDPAddr) error {
	c.Conn = gtpv2.NewConn(localAddr, c.connType, 0)
	if err := c.Conn.Listen(ctx); err != nil {
		return fmt.Errorf("error enabling GTP client: %s", err)
	}
	go func() {
		if ctx == nil {
			ctx = context.Background()
		}
		if err := c.Serve(ctx); err != nil {
			glog.Errorf("error running GTP client: %s", err)
			return
		}
	}()
	//TODO: remove this wait once there is a way to check when the listener is ready
	return nil
}

// Get preferred outbound ip of this machine
func GetLocalOutboundIP(testIp *net.UDPAddr) (net.IP, error) {
	connection, err := net.Dial("udp", testIp.String())
	if err != nil {
		return nil, err
	}

	defer connection.Close()
	localAddr := connection.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// configOrDefaultTimeout sets a default timeout if config timeout is 0
func configOrDefaultTimeout(configTimeout time.Duration) time.Duration {
	if configTimeout == 0 {
		return DefaultGtpTimeout
	}
	return configTimeout
}
