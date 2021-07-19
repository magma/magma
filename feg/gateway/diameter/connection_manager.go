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
	"fmt"
	"sync"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
	"github.com/golang/glog"
)

// ConnectionManager holds a map of connections keyed by the server ip/protocol
// pair
type ConnectionManager struct {
	connMap  map[DiameterServerConnConfig]*Connection // map of DiameterServerConfig -> *lockedConnection
	disabled bool                                     // true is new connection creation is disabled
	rwl      sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{connMap: map[DiameterServerConnConfig]*Connection{}}
}

// GetConnection either gets the existing connection or creates a new one
// if it doesn't exist. This is threadsafe
func (cm *ConnectionManager) GetConnection(client *sm.Client, server *DiameterServerConfig) (*Connection, error) {
	cm.rwl.RLock()
	conn, ok := cm.connMap[server.DiameterServerConnConfig]
	if ok && conn != nil {
		cm.rwl.RUnlock()
		return conn, nil
	}
	cm.rwl.RUnlock()

	glog.V(2).Infof("ConnectionManager: no cached connection for %+v", server)

	cm.rwl.Lock()
	defer cm.rwl.Unlock()
	if cm.disabled {
		err := fmt.Errorf("cannot create connection for %+v, ConnectionManager is disabled", server)
		glog.Error(err)
		return nil, err
	}
	conn, ok = cm.connMap[server.DiameterServerConnConfig]
	if ok && conn != nil { // check again, another thread may have added a connection between RUnlock() & Lock()
		return conn, nil
	}
	conn = newConnection(client, server)
	cm.connMap[server.DiameterServerConnConfig] = conn
	glog.V(2).Infof("ConnectionManager: created connection for %+v", server)

	return conn, nil
}

// AddExistingConnection adds an already in-use connection to the connection manager.
// This is used for servers that need to maintain a connection mapping to clients
// If a connection already exists for the provided server config, update the
// connection manager with the new connection. This is threadsafe
func (cm *ConnectionManager) AddExistingConnection(conn diam.Conn, client *sm.Client, server *DiameterServerConfig) error {
	glog.V(2).Infof("ConnectionManager: adding existing connection for %+v", server)
	cm.rwl.Lock()
	defer cm.rwl.Unlock()
	meta, ok := smpeer.FromContext(conn.Context())
	if !ok {
		err := fmt.Errorf(
			"ConnectionManager: cannot add existing connection for %+v, failed to fetch connection context",
			server)
		glog.Error(err)
		return err
	}
	diameterConnection := &Connection{
		server:   server,
		client:   client,
		conn:     conn,
		metadata: meta,
	}
	cm.connMap[server.DiameterServerConnConfig] = diameterConnection
	return nil
}

// CleanupAllConnections does exactly that
func (cm *ConnectionManager) CleanupAllConnections() {
	cm.rwl.Lock()
	defer cm.rwl.Unlock()
	cm.cleanupConnections()
}

// DisableFor - cleans up all existing connections & disables new connection creations for given duration
func (cm *ConnectionManager) DisableFor(period time.Duration) {
	glog.V(2).Infof("ConnectionManager: disabling connections for %s", period.String())
	cm.rwl.Lock()
	defer cm.rwl.Unlock()
	cm.disabled = true
	cm.cleanupConnections()
	time.AfterFunc(period, func() { cm.Enable() })
}

// Enable - enables new connection creations
func (cm *ConnectionManager) Enable() {
	glog.V(2).Info("ConnectionManager: enabling connections")
	cm.rwl.Lock()
	defer cm.rwl.Unlock()
	cm.disabled = false
}

// Find returns connection object matching given diameter connection or nil if not found
func (cm *ConnectionManager) Find(dc diam.Conn) *Connection {
	if cm == nil || dc == nil {
		return nil
	}
	cm.rwl.RLock()
	defer cm.rwl.RUnlock()
	for _, conn := range cm.connMap {
		if conn != nil {
			if conn.conn == dc {
				return conn
			}
			glog.V(2).Infof("connection mismatch: %+v (%T) != %+v (%T)", conn.conn, conn.conn, dc, dc)
		}
	}
	glog.V(2).Infof("failed to find a match for diam.Conn: '%s' %s->%s",
		dc.RemoteAddr().Network(), dc.LocalAddr().String(), dc.RemoteAddr().String())
	return nil
}

// FindConnection returns an existing connection or returns nil if not found
func (cm *ConnectionManager) FindConnection(server *DiameterServerConfig) *Connection {
	cm.rwl.RLock()
	defer cm.rwl.RUnlock()
	conn, ok := cm.connMap[server.DiameterServerConnConfig]
	if ok {
		return conn
	}
	return nil
}

func (cm *ConnectionManager) cleanupConnections() {
	glog.V(2).Info("ConnectionManager: removing all existing connections")
	for _, c := range cm.connMap {
		if c != nil {
			c.cleanupConnection()
		}
	}
	cm.connMap = map[DiameterServerConnConfig]*Connection{}
}
