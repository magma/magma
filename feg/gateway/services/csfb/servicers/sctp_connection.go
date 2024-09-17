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

package servicers

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"unsafe"

	"github.com/golang/glog"
	"github.com/ishidawataru/sctp"
)

const (
	VLRAddrEnv          = "VLR_ADDR"
	DefaultVLRIPAddress = "127.0.0.1"
	DefaultVLRPort      = 1357
	LocalAddrEnv        = "SGS_LOCAL_ADDR"
	LocalIPAddress      = "127.0.0.1"
	LocalPort           = 0
)

type SCTPClientConnection struct {
	sendConn     *sctp.SCTPConn
	vlrSCTPAddr  *sctp.SCTPAddr
	localSGsAddr *sctp.SCTPAddr
}

func NewSCTPClientConnection(vlrSCTPAddr *sctp.SCTPAddr, localSGsAddr *sctp.SCTPAddr) (*SCTPClientConnection, error) {
	return &SCTPClientConnection{
		vlrSCTPAddr:  vlrSCTPAddr,
		localSGsAddr: localSGsAddr, // nil when it's not specified
	}, nil
}

func (conn *SCTPClientConnection) EstablishConn() error {
	glog.V(2).Infof("Establishing SCTP connection with %s", conn.vlrSCTPAddr)
	sendConn, err := sctp.DialSCTP(
		"sctp",
		conn.localSGsAddr,
		conn.vlrSCTPAddr,
	)
	if err != nil {
		return err
	}
	glog.V(2).Info("SCTP connection with VLR established")
	conn.sendConn = sendConn

	return nil
}

func (conn *SCTPClientConnection) CloseConn() error {
	glog.V(2).Info("Closing SCTP connection with VLR")
	if conn.sendConn == nil {
		return errors.New("connection to VLR not established")
	}

	err := conn.sendConn.Close()
	conn.sendConn = nil
	if err != nil {
		return err
	}
	glog.V(2).Info("SCTP connection with VLR closed successfully")

	return nil
}

func (conn *SCTPClientConnection) Send(message []byte) error {
	glog.V(2).Info("Sending message to VLR through SCTP")
	if conn.sendConn == nil {
		return errors.New("connection to VLR not established")
	}

	ppid := 0

	info := &sctp.SndRcvInfo{
		Stream: uint16(ppid),
		PPID:   uint32(ppid),
	}

	conn.sendConn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	_, err := conn.sendConn.SCTPWrite(message, info)
	if err != nil {
		return err
	}
	glog.V(2).Info("Message sent successfully to VLR")

	return nil
}

func (conn *SCTPClientConnection) Receive() ([]byte, error) {
	if conn.sendConn == nil {
		return []byte{}, errors.New("connection to VLR not established")
	}

	buf := make([]byte, 254)
	n, _, err := conn.sendConn.SCTPRead(buf)
	if err != nil {
		return []byte{}, err
	}

	return buf[:n], err
}

type SCTPServerConnection struct {
	rcvListener     *sctp.SCTPListener
	sendConn        *sctp.SCTPConn
	infoWrappedConn *sctp.SCTPSndRcvInfoWrappedConn
}

func NewSCTPServerConnection() (*SCTPServerConnection, error) {
	return &SCTPServerConnection{
		rcvListener: nil,
		sendConn:    nil,
	}, nil
}

func (conn *SCTPServerConnection) StartListener(ipAddr string, port PortNumber) (PortNumber, error) {
	ln, err := sctp.ListenSCTP("sctp", ConstructSCTPAddr(ipAddr, port))
	if err != nil {
		return -1, err
	}

	conn.rcvListener = ln
	address := url.URL{Host: ln.Addr().String()}
	portNumber, err := strconv.Atoi(address.Port())
	if err != nil {
		return -1, err
	}

	return portNumber, nil
}

func (conn *SCTPServerConnection) CloseListener() error {
	err := conn.rcvListener.Close()
	if err != nil {
		return err
	}

	return nil
}

func (conn *SCTPServerConnection) ConnectionEstablished() bool {
	return conn.sendConn != nil
}

func (conn *SCTPServerConnection) AcceptConn() error {
	netConn, err := conn.rcvListener.Accept()
	if err != nil {
		return err
	}
	conn.infoWrappedConn = sctp.NewSCTPSndRcvInfoWrappedConn(netConn.(*sctp.SCTPConn))
	conn.sendConn = netConn.(*sctp.SCTPConn)
	return nil
}

func (conn *SCTPServerConnection) CloseConn() error {
	if conn.sendConn == nil {
		return errors.New("connection to client not established")
	}
	err := conn.sendConn.Close()
	conn.sendConn = nil
	conn.infoWrappedConn = nil
	if err != nil {
		return err
	}
	return nil
}

func (conn *SCTPServerConnection) ReceiveThroughListener() ([]byte, error) {
	if conn.infoWrappedConn == nil {
		return []byte{}, errors.New("connection to client not established")
	}

	buf := make([]byte, 254)
	n, err := conn.infoWrappedConn.Read(buf)

	if err != nil {
		return []byte{}, err
	}

	return buf[unsafe.Sizeof(sctp.SndRcvInfo{}):n], err
}

func (conn *SCTPServerConnection) SendFromServer(msg []byte) error {
	if conn.sendConn == nil {
		return errors.New("connection to client not established")
	}

	err := conn.sendConn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	if err != nil {
		return err
	}

	ppid := 0
	info := &sctp.SndRcvInfo{
		Stream: uint16(ppid),
		PPID:   uint32(ppid),
	}

	_, err = conn.sendConn.SCTPWrite(msg, info)
	if err != nil {
		return err
	}

	return nil
}

func (conn *SCTPClientConnection) GetVlrIPandPort() ([]string, int) {
	return conn.getIPandPort(conn.vlrSCTPAddr)
}

func (conn *SCTPClientConnection) GetLocalIPandPort() ([]string, int) {
	return conn.getIPandPort(conn.localSGsAddr)
}

func (conn *SCTPClientConnection) getIPandPort(sctpAddr *sctp.SCTPAddr) ([]string, int) {
	ipsString := make([]string, 1)
	for _, IPAddr := range sctpAddr.IPAddrs {
		ipsString = append(ipsString, IPAddr.String())
	}
	return ipsString, sctpAddr.Port
}

func convertIPAddressFromStrip(addresString string) (*sctp.SCTPAddr, error) {
	ips, port, err := SplitIP(addresString)
	if err != nil {
		return &sctp.SCTPAddr{}, err
	}
	sctpAddr := ConstructSCTPAddr(ips, port)
	return sctpAddr, nil
}

func SplitIP(IPsandPort string) (ipStr string, portInt int, err error) {
	portStr := ""
	ipStr, portStr, err = net.SplitHostPort(IPsandPort)
	portInt = 0
	if err != nil {
		glog.Errorf("Couldn't parse the whole string of IPs %s", IPsandPort)
		return
	}
	portInt, err = strconv.Atoi(portStr)
	if err != nil {
		glog.Errorf("Couldn't parse port %s", IPsandPort)
		return
	}
	return
}

// Suports strings with multiple IP comma separated like "ip" or "ip1,ip2,ip3"
func ConstructSCTPAddr(ip string, port int) *sctp.SCTPAddr {
	ips := []net.IPAddr{}
	for _, i := range strings.Split(ip, ",") {
		if a, err := net.ResolveIPAddr("ip", i); err == nil {
			ips = append(ips, *a)
		}
	}
	return &sctp.SCTPAddr{
		IPAddrs: ips,
		Port:    port,
	}
}
