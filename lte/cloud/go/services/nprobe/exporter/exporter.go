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

package exporter

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"sync"

	"github.com/gogf/gf/net/gtcp"
	"github.com/golang/glog"

	"magma/lte/cloud/go/services/nprobe/obsidian/models"
)

// RecordExporter sends records to a remote host over tcp/tls
type RecordExporter struct {
	tlsConfig  *tls.Config
	conn       *gtcp.Conn
	remoteAddr string
	mutex      sync.Mutex
}

// newTlsConfig creates a new TLS config from the client certificates
func newTlsConfig(crtFile, keyFile []byte, skipVerify bool) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(crtFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: skipVerify,
	}, nil
}

// NewRecordExporter creates a new tls exporter and attempt to establish a connection at start
func NewRecordExporter(destination *models.NetworkProbeDestination) (*RecordExporter, error) {
	details := destination.DestinationDetails
	if details.Certificate == nil || details.PrivateKey == nil {
		return nil, fmt.Errorf("Certificates not provided for destination %v", destination.DestinationID)
	}

	certPem, err := base64.StdEncoding.DecodeString(details.Certificate.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to parse certificate for destination %v", destination.DestinationID)
	}

	keyPem, err := base64.StdEncoding.DecodeString(details.PrivateKey.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to parse private key for destination %v", destination.DestinationID)
	}

	tlsConfig, err := newTlsConfig(certPem, keyPem, details.SkipVerifyServer)
	if err != nil {
		return nil, fmt.Errorf("Failed to to create new TlsConfig %v", err)
	}

	client := &RecordExporter{
		tlsConfig:  tlsConfig,
		remoteAddr: details.DeliveryAddress,
	}
	conn, err := client.getTlsConnection() // attempt to establish connection at start
	if err != nil {
		glog.Errorf(
			"Failed to establish new TLS connection from to '%s'; error: %v, will retry later.",
			details.DeliveryAddress, err)
		return nil, err
	}
	client.conn = conn
	return client, nil
}

// SendMessageWithRetries writes data to remote address with a retry counter
func (c *RecordExporter) SendMessageWithRetries(message []byte, retryCount uint32) error {
	var err error
	for i := 0; i < int(retryCount); i++ {
		err = c.sendMessage(message)
		// send succeeded
		if err == nil {
			return nil
		}
	}
	return err
}

// sendMessage sends a single message on the connection. If the connection is
// not established, this establishes it. If the message sending fails, the
// connection is closed
func (c *RecordExporter) sendMessage(message []byte) error {
	conn, err := c.getTlsConnection()
	if err != nil {
		return err
	}

	// It's possible that the connection is closed here in contention for the
	// connection. This is handled as an error and the sending can retry
	err = c.conn.Send(message)
	if err != nil {
		// write failed, close and cleanup connection
		c.destroyConnection(conn)
	}
	return err
}

// getTlsConnection returns the existing connection or
// dials and initializes a connection if it doesn't exist
func (c *RecordExporter) getTlsConnection() (*gtcp.Conn, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.conn != nil {
		return c.conn, nil
	}
	if len(c.remoteAddr) == 0 {
		return nil, errors.New("Invalid remote address")
	}

	conn, err := gtcp.NewConnTLS(c.remoteAddr, c.tlsConfig)
	c.conn = conn
	return c.conn, err
}

// destroyConnection closes a bad connection. If the connection
// passed is the same as the one stored in the locked connection, it is nullified.
// If the passed connection is not the same, this probably means another go routine
// already created a new connection - just try to close it and return.
func (c *RecordExporter) destroyConnection(conn *gtcp.Conn) {
	if conn == nil {
		return
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if conn == c.conn {
		c.conn = nil
	}
	conn.Close()
}
