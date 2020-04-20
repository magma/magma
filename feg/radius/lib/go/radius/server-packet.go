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

package radius

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"sync/atomic"
)

type packetResponseWriter struct {
	// listener that received the packet
	conn                 net.PacketConn
	addr                 net.Addr
	requestAuthenticator [16]byte
	secret               []byte
}

// MessageAuthenticatorAttrLength the length, in bytes, of the Message-Authenticator
// attribute, including attribute Type and Length fields
const MessageAuthenticatorAttrLength uint16 = 18

func (r *packetResponseWriter) Write(packet *Packet) error {
	encoded, err := packet.Encode()
	if err != nil {
		return err
	}

	// Add Message-Authenticator if needed
	// TODO: Cannot reference rfc2869 package and use rfc2869.EAPMessage_Type,
	// because this creates a circular dependecy.
	_, hasEapMessage := packet.Lookup(Type(79))
	if hasEapMessage && packet.Code.ImpliesMessageAuthenticatorNeeded() {
		encoded = r.addMessageAuthenticator(encoded)
	}

	if _, err := r.conn.WriteTo(encoded, r.addr); err != nil {
		return err
	}
	return nil
}

// ImpliesMessageAuthenticatorNeeded indicates if the RadiusCode implies
// Message Authenticator is needed, as per rfc3579 section 3.2 and
// rfc2869 section 5.14
func (c Code) ImpliesMessageAuthenticatorNeeded() bool {
	return c == CodeAccessAccept || c == CodeAccessReject || c == CodeAccessChallenge
}

func (r *packetResponseWriter) addMessageAuthenticator(encoded []byte) []byte {
	// Fix the size
	size := binary.BigEndian.Uint16(encoded[2:4]) + MessageAuthenticatorAttrLength
	binary.BigEndian.PutUint16(encoded[2:4], uint16(size))

	// Add Message Authenticator Attribute, 0 padded, and flatten
	zeroedOutMsgAuthenticator := [16]byte{}
	allBytes := [][]byte{
		encoded[:4],
		r.requestAuthenticator[:],
		encoded[20:],
		[]byte{80, 18},
		zeroedOutMsgAuthenticator[:],
	}
	var radiusMsg []byte
	for _, b := range allBytes {
		radiusMsg = append(radiusMsg[:], b[:]...)
	}

	// Calculate Message Authenticator & Overwrite
	hash := hmac.New(md5.New, r.secret)
	hash.Write(radiusMsg)
	encoded = hash.Sum(radiusMsg[:len(radiusMsg)-16])

	// Re-calc the Response Authenticator
	resAuth := md5.New()
	resAuth.Write(encoded[:4])
	resAuth.Write(r.requestAuthenticator[:])
	resAuth.Write(encoded[20:])
	resAuth.Write(r.secret)
	resAuth.Sum(encoded[4:4:20])

	return encoded
}

// PacketServer listens for RADIUS requests on a packet-based protocols (e.g.
// UDP).
type PacketServer struct {
	// The address on which the server listens. Defaults to :1812.
	Addr string
	// The network on which the server listens. Defaults to udp.
	Network      string
	SecretSource SecretSource
	Handler      Handler

	// Skip incoming packet authenticity validation.
	// This should only be set to true for debugging purposes.
	InsecureSkipVerify bool

	// Channel to indicate when server is listenning and ready to serve requests
	Ready chan bool

	mu           sync.Mutex
	shuttingDown bool
	ctx          context.Context
	ctxDone      context.CancelFunc
	running      chan struct{}
	listeners    map[net.PacketConn]int
	activeCount  int32
}

// TODO: logger on PacketServer

// Serve accepts incoming connections on conn.
func (s *PacketServer) Serve(conn net.PacketConn) error {
	if s.Handler == nil {
		return errors.New("radius: nil Handler")
	}
	if s.SecretSource == nil {
		return errors.New("radius: nil SecretSource")
	}

	s.mu.Lock()
	if s.shuttingDown {
		s.mu.Unlock()
		return ErrServerShutdown
	}
	var ctx context.Context
	if s.ctx == nil {
		s.ctx, s.ctxDone = context.WithCancel(context.Background())
		ctx = s.ctx
	}
	if s.running == nil {
		s.running = make(chan struct{})
	}
	if s.listeners == nil {
		s.listeners = make(map[net.PacketConn]int)
	}
	s.listeners[conn]++
	s.mu.Unlock()

	type activeKey struct {
		IP         string
		Identifier byte
	}

	var (
		activeLock sync.Mutex
		active     = map[activeKey]struct{}{}
	)

	atomic.AddInt32(&s.activeCount, 1)
	defer func() {
		s.mu.Lock()
		s.listeners[conn]--
		if s.listeners[conn] == 0 {
			delete(s.listeners, conn)
		}
		s.mu.Unlock()

		if atomic.AddInt32(&s.activeCount, -1) == 0 {
			s.mu.Lock()
			s.shuttingDown = false
			close(s.running)
			s.running = nil
			s.ctx = nil
			s.mu.Unlock()
		}
	}()

	var buff [MaxPacketLength]byte
	for {
		n, remoteAddr, err := conn.ReadFrom(buff[:])
		if err != nil {
			s.mu.Lock()
			if s.shuttingDown {
				s.mu.Unlock()
				return nil
			}
			s.mu.Unlock()

			if ne, ok := err.(net.Error); ok && !ne.Temporary() {
				return err
			}
			// TODO: log error?
			continue
		}

		buffCopy := append([]byte(nil), buff[:n]...)

		atomic.AddInt32(&s.activeCount, 1)
		go func(buff []byte, remoteAddr net.Addr) {
			secret, err := s.SecretSource.RADIUSSecret(ctx, remoteAddr)
			if err != nil {
				// TODO: log only if server is not shutting down?
				return
			}
			if len(secret) == 0 {
				return
			}

			if !s.InsecureSkipVerify && !IsAuthenticRequest(buff, secret) {
				// TODO: log?
				return
			}

			packet, err := Parse(buff, secret)
			if err != nil {
				// TODO: error logger
				return
			}

			key := activeKey{
				IP:         remoteAddr.String(),
				Identifier: packet.Identifier,
			}
			activeLock.Lock()
			if _, ok := active[key]; ok {
				activeLock.Unlock()
				return
			}
			active[key] = struct{}{}
			activeLock.Unlock()

			response := packetResponseWriter{
				conn:                 conn,
				addr:                 remoteAddr,
				requestAuthenticator: packet.Authenticator,
				secret:               secret,
			}

			defer func() {
				activeLock.Lock()
				delete(active, key)
				activeLock.Unlock()

				if atomic.AddInt32(&s.activeCount, -1) == 0 {
					s.mu.Lock()
					s.shuttingDown = false
					close(s.running)
					s.running = nil
					s.ctx = nil
					s.mu.Unlock()
				}
			}()

			request := Request{
				LocalAddr:  conn.LocalAddr(),
				RemoteAddr: remoteAddr,
				Packet:     packet,
				ctx:        ctx,
			}

			s.Handler.ServeRADIUS(&response, &request)
		}(buffCopy, remoteAddr)
	}
}

// ListenAndServe starts a RADIUS server on the address given in s.
func (s *PacketServer) ListenAndServe() error {
	if s.Handler == nil {
		return errors.New("radius: nil Handler")
	}
	if s.SecretSource == nil {
		return errors.New("radius: nil SecretSource")
	}

	addrStr := ":1812"
	if s.Addr != "" {
		addrStr = s.Addr
	}

	network := "udp"
	if s.Network != "" {
		network = s.Network
	}
	pc, err := net.ListenPacket(network, addrStr)
	if err != nil {
		if s.Ready != nil {
			s.Ready <- false
		}
		return err
	}
	defer pc.Close()

	// Signal server is ready & serving requests
	if s.Ready != nil {
		s.Ready <- true
	}
	return s.Serve(pc)
}

// Shutdown gracefully stops the server. It first closes all listeners (which
// stops accepting new packets) and then waits for running handlers to complete.
//
// Shutdown returns after all handlers have completed, or when ctx is canceled.
// The PacketServer is ready for re-use once the function returns nil.
func (s *PacketServer) Shutdown(ctx context.Context) error {
	s.mu.Lock()

	if len(s.listeners) == 0 {
		s.mu.Unlock()
		return nil
	}

	if !s.shuttingDown {
		s.shuttingDown = true
		s.ctxDone()
		for listener := range s.listeners {
			listener.Close()
		}
	}

	ch := s.running
	s.mu.Unlock()
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
