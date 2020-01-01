/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package diameter

import (
	"errors"
	"sync"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"github.com/fiorix/go-diameter/v4/diam/sm/smpeer"
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
	cm.rwl.Lock()
	defer cm.rwl.Unlock()
	if cm.disabled {
		return nil, errors.New("ConnectionManager: Connection Creation Is Disabled")
	}
	conn, ok = cm.connMap[server.DiameterServerConnConfig]
	if ok && conn != nil { // check again, another thread may have added a connection between RUnlock() & Lock()
		return conn, nil
	}
	conn = newConnection(client, server)
	cm.connMap[server.DiameterServerConnConfig] = conn
	return conn, nil
}

// AddExistingConnection adds an already in-use connection to the connection manager.
// This is used for servers that need to maintain a connection mapping to clients
// If a connection already exists for the provided server config, update the
// connection manager with the new connection. This is threadsafe
func (cm *ConnectionManager) AddExistingConnection(conn diam.Conn, client *sm.Client, server *DiameterServerConfig) error {
	cm.rwl.Lock()
	defer cm.rwl.Unlock()
	meta, ok := smpeer.FromContext(conn.Context())
	if !ok {
		return errors.New("ConnectionManager: Cannot add existing connection - could not fetch connection context")
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
	cm.rwl.Lock()
	defer cm.rwl.Unlock()
	cm.disabled = true
	cm.cleanupConnections()
	time.AfterFunc(period, func() { cm.Enable() })
}

// Enable - enables new connection creations
func (cm *ConnectionManager) Enable() {
	cm.rwl.Lock()
	defer cm.rwl.Unlock()
	cm.disabled = false
}

func (cm *ConnectionManager) cleanupConnections() {
	for _, c := range cm.connMap {
		if c != nil {
			c.cleanupConnection()
		}
	}
	cm.connMap = map[DiameterServerConnConfig]*Connection{}
}
