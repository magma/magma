// package servce implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/examples/s6a_proxy/protos"
	"google.golang.org/grpc/codes"
)

// genSID generates a unique string to be used as diameter message Session ID
// the generated session ID consists of a const prefix 's6a_proxy;', unix timestamp & rand
func genSID() string {
	return fmt.Sprintf("s6a_proxy;%d_%X", uint(time.Now().Unix()), rand.Uint32())
}

// updateSession resets chan in the session map & closes and returns previous chan associated with the sid if any
func (s *s6aProxy) updateSession(sid string, ch chan interface{}) chan interface{} {
	s.sessionsMu.Lock()
	oldCh, ok := s.sessions[sid]
	s.sessions[sid] = ch
	s.sessionsMu.Unlock()
	if ok {
		if oldCh != nil {
			close(oldCh)
		}
		return oldCh
	}
	return nil
}

// updateSession deletes chanel entry in the session map & closes and returns previous chan associated with the sid
func (s *s6aProxy) cleanupSession(sid string) chan interface{} {
	s.sessionsMu.Lock()
	oldCh, ok := s.sessions[sid]
	if ok {
		delete(s.sessions, sid)
		s.sessionsMu.Unlock()
		if oldCh != nil {
			close(oldCh)
		}
		return oldCh
	}
	s.sessionsMu.Unlock()
	return nil
}

// acquireConnection takes Reader lock for s and ensures that connection is established,
// It leaves s in the "locked" state and must be followed by call to releaseConnection()
func (s *s6aProxy) acquireConnection() (diam.Conn, error) {
	var err error
	s.mu.RLock()
	if s.conn == nil {
		s.mu.RUnlock()
		s.mu.Lock()
		if s.conn == nil {
			s.conn, err = s.smClient.DialNetwork(s.cfg.Protocol, s.cfg.HssAddr)
			if err == nil && s.conn != nil {
				closeChan := s.conn.(diam.CloseNotifier).CloseNotify()
				go func() {
					<-closeChan
					s.mu.Lock()
					if s.conn != nil {
						s.conn.Close()
					}
					s.conn = nil
					s.mu.Unlock()
				}()
			}
		}
		s.mu.Unlock()
		s.mu.RLock()
		// There is a chance, somebody closes the conn between Unlock() & RLock() due to a R/W error
		if s.conn == nil && err == nil {
			err = fmt.Errorf("Connection was closed unexpectedly")
		}
	}
	return s.conn, err
}

func (s *s6aProxy) releaseConnection() {
	s.mu.RUnlock()
}

func (s *s6aProxy) cleanupConn(c diam.Conn) {
	s.mu.Lock()
	if c == s.conn {
		s.conn = nil
	}
	s.mu.Unlock()
	c.Close()
}

// TranslateBaseDiamResultCode maps Base Diameter Result Code to GRPC Status Error and returns it,
// Diam success codes will result in nil error returned
func TranslateBaseDiamResultCode(diamResult uint32) error {
	if diamResult == uint32(protos.ErrorCode_UNDEFINED) { // diamResult was not set (default will be 0)
		return nil
	}
	// diam result code is 2xxx
	if diamResult >= uint32(protos.ErrorCode_SUCCESS) && diamResult < uint32(protos.ErrorCode_COMMAND_UNSUPORTED) {
		return nil
	}
	errName, ok := protos.ErrorCode_name[int32(diamResult)]
	if !ok {
		errName = "BASE_DIAMETER"
	}
	return Errorf(codes.Code(diamResult), "Diameter Error: %d (%s)", diamResult, errName)
}
