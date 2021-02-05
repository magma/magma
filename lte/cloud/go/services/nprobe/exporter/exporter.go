package exporter

import (
	"crypto/tls"
	"errors"
	"sync"

	"github.com/gogf/gf/net/gtcp"
	"github.com/golang/glog"
)

// RecordExporter sends records to a remote host over tcp/tls
type RecordExporter struct {
	tlsConfig  *tls.Config
	conn       *gtcp.Conn
	remoteAddr string
	mutex      sync.Mutex
}

// NewTlsConfig creates a new TLS config from the client certificates
func NewTlsConfig(crtFile, keyFile, rootCAFile string, skipVerify bool) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		return nil, err
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: skipVerify}, nil
}

// NewRecordExporter creates a new tls exporter and attempt to establish a connection at start
func NewRecordExporter(remoteAddr string, tlsConfig *tls.Config) (*RecordExporter, error) {
	client := &RecordExporter{
		tlsConfig:  tlsConfig,
		remoteAddr: remoteAddr,
	}
	var err error
	go func() { // init connection in a goroutine, it can block for long time
		conn, err := client.getTlsConnection() // attempt to establish connection at start
		if err != nil {
			glog.Errorf(
				"Failed to establish new TLS connection from to '%s'; error: %v, will retry later.",
				remoteAddr, err)
		}
		client.conn = conn
	}()
	return client, err
}

// SendRecord writes data to remote address with a retry counter
func (c *RecordExporter) SendRecord(record []byte, retryCount uint32) error {
	return c.sendMessageWithRetries(record, retryCount)
}

func (c *RecordExporter) sendMessageWithRetries(message []byte, retryCount uint32) error {
	var err error
	timesToSend := retryCount + 1
	for ; timesToSend > 0; timesToSend-- {
		err = c.sendMessage(message)
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
func (c *RecordExporter) sendMessage(message []byte) error {
	var err error
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

// getDiamConnection returns the existing connection or
// dials and initializes a connection if it doesn't exist
func (c *RecordExporter) getTlsConnection() (*gtcp.Conn, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.conn != nil {
		return c.conn, nil
	}
	var err error
	if len(c.remoteAddr) == 0 {
		return nil, errors.New("Invalid remote address")
	}
	if c.tlsConfig == nil {
		return nil, errors.New("Invalid tls config")
	}

	conn, err := gtcp.NewConnTLS(c.remoteAddr, c.tlsConfig)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	return conn, nil
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
	if conn == c.conn {
		c.conn = nil
	}
	c.mutex.Unlock()
	conn.Close()
}

// cleanupConnection is similar to destroyConnection, but it closes and
// cleans up connection unconditionally
func (c *RecordExporter) cleanupConnection() {
	c.mutex.Lock()
	conn := c.conn
	if conn != nil {
		c.conn = nil
		c.mutex.Unlock()
		conn.Close()
	} else {
		c.mutex.Unlock()
	}
}
