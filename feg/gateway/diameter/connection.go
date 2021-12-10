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

package diameter

import (
	"errors"
	"net"
	"strings"
	"sync"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
	"github.com/golang/glog"
	"github.com/ishidawataru/sctp"
)

type messageTypeEnum uint8

const (
	requestMessage messageTypeEnum = 1
	answerMessage  messageTypeEnum = 2
)

// Connection is representing a diameter connection that you can
// send messages to and get metadata from (for building AVPs)
type Connection struct {
	conn     diam.Conn
	metadata *smpeer.Metadata
	server   *DiameterServerConfig
	client   *sm.Client
	mutex    sync.Mutex
}

func newConnection(client *sm.Client, server *DiameterServerConfig) *Connection {
	conn := &Connection{
		server: server,
		client: client,
	}
	go func() { // init connection in a goroutine, it can block for long time
		_, _, err := conn.getDiamConnection() // attempt to establish connection at start
		if err != nil {
			glog.Errorf(
				"Failed to establish new %s diameter connection from '%s' to '%s'; error: %v, will retry later.",
				server.Protocol, server.LocalAddr, server.Addr, err)
		}
	}()
	return conn
}

func (c *Connection) SendAnswer(message *diam.Message, retryCount uint) error {
	return c.sendMessageWithRetries(message, answerMessage, retryCount, nil)
}

func (c *Connection) SendRequest(message *diam.Message, retryCount uint) error {
	return c.sendMessageWithRetries(message, requestMessage, retryCount, nil)
}

func (c *Connection) SendRequestToServer(message *diam.Message, retryCount uint, server *DiameterServerConfig) error {
	return c.sendMessageWithRetries(message, requestMessage, retryCount, server)
}

func (c *Connection) sendMessageWithRetries(
	message *diam.Message, messageType messageTypeEnum, retryCount uint, server *DiameterServerConfig) error {

	var err error
	timesToSend := retryCount + 1
	for ; timesToSend > 0; timesToSend-- {
		err = c.sendMessage(message, messageType, server)
		// send succeeded
		if err == nil {
			break
		}
	}
	return err
}

// sendMessage sends a single message on the connection. If the connection is
// not established, this establishes it. If the message sending fails, the
// connection is closed
func (c *Connection) sendMessage(
	message *diam.Message, messageType messageTypeEnum, server *DiameterServerConfig) error {

	var err error
	conn, metadata, err := c.getDiamConnection()
	if err != nil {
		return err
	}

	if messageType == requestMessage {
		srvr := c.server
		if server != nil {
			srvr = server
		}
		// add destination for diameter requests only
		message, err = addDestinationToMessage(message, metadata, srvr)
		if err != nil {
			return err
		}
	}

	// It's possible that the connection is closed here in contention for the
	// connection. This is handled as an error and the sender can retry
	_, err = message.WriteTo(conn)
	if err != nil {
		// write failed, close and cleanup connection
		c.destroyConnection(conn)
	}
	return err
}

// getDiamConnection returns the existing connection and its metadata or
// dials and initializes a connection if it doesn't exist and returns it and its metadata,
func (c *Connection) getDiamConnection() (diam.Conn, *smpeer.Metadata, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.conn != nil {
		return c.conn, c.metadata, nil
	}
	var (
		localAddr net.Addr
		err       error
	)
	if len(c.server.LocalAddr) > 0 {
		if len(c.server.Protocol) == 0 || strings.HasPrefix(c.server.Protocol, "tcp") {
			localAddr, err = net.ResolveTCPAddr(c.server.Protocol, c.server.LocalAddr)
		} else if strings.HasPrefix(c.server.Protocol, "sctp") {
			localAddr, err = sctp.ResolveSCTPAddr(c.server.Protocol, c.server.LocalAddr)
		}
		if err != nil {
			return nil, nil, errors.New(
				"Invalid " + c.server.Protocol + " local address '" + c.server.LocalAddr + "':" + err.Error())
		}
	}
	conn, err := c.client.DialExt(c.server.Protocol, c.server.Addr, 0, localAddr)
	if err != nil {
		return nil, nil, err
	}
	metadata, ok := smpeer.FromContext(conn.Context())
	if !ok {
		conn.Close()
		return nil, nil, errors.New("Could not obtain metadata from connection")
	}
	c.conn, c.metadata = conn, metadata
	return conn, metadata, nil
}

// destroyConnection closes a bad connection. If the connection
// passed is the same as the one stored in the locked connection, it is nullified.
// If the passed diam connection is not the same, this probably means another go routine
// already created a new connection - just try to close it and return.
func (c *Connection) destroyConnection(conn diam.Conn) {
	if conn == nil {
		return
	}
	c.mutex.Lock()
	match := conn == c.conn
	if match {
		c.conn = nil
		c.metadata = nil
	}
	c.mutex.Unlock()

	if glog.V(2) {
		localAddress := "nil"
		remoteAddress := "nil"
		if conn.LocalAddr() != nil {
			localAddress = conn.LocalAddr().String()
		}
		if conn.LocalAddr() != nil {
			remoteAddress = conn.RemoteAddr().String()
		}
		if match {
			glog.Infof("destroyed %s->%s connection", localAddress, remoteAddress)
		} else {
			glog.Infof(
				"cannot destroy mismatched %s->%s connection", localAddress, remoteAddress)
		}
	}
	conn.Close()
}

// cleanupConnection is similar to destroyConnection, but it closes & cleans up connection unconditionally
func (c *Connection) cleanupConnection() {
	c.mutex.Lock()
	conn := c.conn
	if conn != nil {
		c.conn = nil
		c.metadata = nil
		c.mutex.Unlock()
		if glog.V(2) {
			glog.Infof("cleaned up %s->%s connection", conn.LocalAddr().String(), conn.RemoteAddr().String())
		}
		conn.Close()
	} else {
		c.mutex.Unlock()
	}
}

// addDestinationToMessage adds the destination host/realm AVPs to the message
// unless they are already in the message
func addDestinationToMessage(
	message *diam.Message, metadata *smpeer.Metadata, server *DiameterServerConfig) (*diam.Message, error) {

	if metadata == nil && (server == nil || len(server.DestHost) == 0 || len(server.DestRealm) == 0) {
		return message, errors.New("Could not add destination to message, empty metadata & invalid server realm")
	}
	var destHost datatype.DiameterIdentity
	if len(server.DestHost) > 0 {
		destHost = datatype.DiameterIdentity(server.DestHost)
	} else {
		destHost = metadata.OriginHost
	}
	var destRealm datatype.DiameterIdentity
	if len(server.DestRealm) > 0 {
		destRealm = datatype.DiameterIdentity(server.DestRealm)
	} else {
		destRealm = metadata.OriginRealm
	}
	realmAVP, err := message.FindAVP(avp.DestinationRealm, 0)
	if err != nil {
		message.NewAVP(avp.DestinationRealm, avp.Mbit, 0, destRealm)
	} else if realmAVP != nil {
		// apply new realm
		realmAVP.Data = destRealm
	}
	if server.DisableDestHost {
		return message, nil
	}
	hostAVP, err := message.FindAVP(avp.DestinationHost, 0)
	if err != nil {
		message.NewAVP(avp.DestinationHost, avp.Mbit, 0, destHost)
	} else if hostAVP != nil {
		if server.OverwriteDestHost {
			// apply new host
			hostAVP.Data = destHost
		}
	}
	return message, nil
}
