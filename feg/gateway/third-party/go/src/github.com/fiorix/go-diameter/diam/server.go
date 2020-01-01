// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter server, based on net/http.

package diam

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// The Handler interface allow arbitrary objects to be
// registered to serve particular messages like CER, DWR.
type Handler interface {
	// ServeDIAM should write messages to the Conn and then return.
	// Returning signals that the request is finished and that the
	// server can move on to the next request on the connection.
	ServeDIAM(Conn, *Message)
}

// Conn interface is used by a handler to send diameter messages.
type Conn interface {
	Write(b []byte) (int, error)                    // Writes a msg to the connection
	WriteStream(b []byte, stream uint) (int, error) // Writes a msg to the connection's stream
	Close()                                         // Close the connection
	LocalAddr() net.Addr                            // Returns the local IP
	RemoteAddr() net.Addr                           // Returns the remote IP
	TLS() *tls.ConnectionState                      // TLS or nil when not using TLS
	Dictionary() *dict.Parser                       // Dictionary parser of the connection
	Context() context.Context                       // Returns the internal context
	SetContext(ctx context.Context)                 // Stores a new context
	Connection() net.Conn                           // Returns network connection
}

// The CloseNotifier interface is implemented by Conns which
// allow detecting when the underlying connection has gone away.
//
// This mechanism can be used to detect if a peer has disconnected.
type CloseNotifier interface {
	// CloseNotify returns a channel that is closed
	// when the client connection has gone away.
	CloseNotify() <-chan struct{}
}

// A liveSwitchReader is a switchReader that's safe for concurrent
// reads and switches, if its mutex is held.
type liveSwitchReader struct {
	sync.Mutex
	r         io.Reader
	pr        *io.PipeReader
	pipeCopyF func()
}

func (sr *liveSwitchReader) Read(p []byte) (n int, err error) {
	sr.Lock()
	// Check if closeNotifier was created prior to this Read call & start it
	if sr.pr != nil && sr.pipeCopyF != nil {
		go sr.pipeCopyF()
		sr.r = sr.pr
		sr.pr = nil
		sr.pipeCopyF = nil
	}
	r := sr.r
	sr.Unlock()
	return r.Read(p)
}

// conn represents the server side of a diameter connection.
type conn struct {
	server   *Server              // the Server on which the connection arrived
	rwc      net.Conn             // i/o connection
	sr       liveSwitchReader     // reads from rwc
	buf      *bufio.ReadWriter    // buffered(sr, rwc)
	tlsState *tls.ConnectionState // or nil when not using TLS
	writer   *response            // the diam.Conn exposed to handlers

	mu           sync.Mutex // guards the following
	closeNotifyc chan struct{}
	clientGone   bool
}

func (c *conn) closeNotify() <-chan struct{} {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closeNotifyc == nil {
		c.closeNotifyc = make(chan struct{})

		if msc, isMulti := c.rwc.(MultistreamConn); isMulti {
			// MultistreamConn provides it's own error handler
			msc.SetErrorHandler(func(mc MultistreamConn, err error) {
				mc.Close()
				c.notifyClientGone()
			})
		} else {
			pr, pw := io.Pipe()
			c.sr.Lock()
			readSource := c.sr.r
			c.sr.pr = pr
			// Create closeNotifier pipe copy routine, but do not start it here
			// If we start it immediately, pipe Write can block indefinitely if we are already in
			// liveSwitchReader.Read() with original sr.r since Pipe.Write blocks in absence of corresponding
			// pipe reader
			// We should only swap the reader outside of r.Read call
			c.sr.pipeCopyF = func() {
				_, err := io.Copy(pw, readSource)
				if err == nil {
					err = io.EOF
				}
				pw.CloseWithError(err)
				c.notifyClientGone()
			}
			c.sr.Unlock()
		}
	}
	return c.closeNotifyc
}

func (c *conn) notifyClientGone() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closeNotifyc != nil && !c.clientGone {
		close(c.closeNotifyc) // unblock readers
		c.clientGone = true
	}
}

// Create new connection from rwc.
func (srv *Server) newConn(rwc net.Conn) (c *conn, err error) {
	msc, isMulti := rwc.(MultistreamConn)
	if isMulti {
		c = &conn{
			server: srv,
			rwc:    msc,
		}
	} else {
		c = &conn{
			server: srv,
			rwc:    rwc,
			sr:     liveSwitchReader{r: rwc},
		}
		c.buf = bufio.NewReadWriter(bufio.NewReader(&c.sr), bufio.NewWriter(rwc))
	}
	c.writer = &response{conn: c}
	return c, nil
}

// Read next message from connection.
func (c *conn) readMessage() (m *Message, err error) {
	if c.server.ReadTimeout > 0 {
		c.rwc.SetReadDeadline(time.Now().Add(c.server.ReadTimeout))
	}
	if msc, isMulti := c.rwc.(MultistreamConn); isMulti {
		// If it's a multi-stream association - reset the stream to "undefined" prior to reading next message
		msc.ResetCurrentStream()
		m, err = ReadMessage(msc, c.dictionary()) // MultistreamConn has it's own buffering
	} else {
		m, err = ReadMessage(c.buf.Reader, c.dictionary())
	}
	if err != nil {
		return nil, err
	}
	return m, nil
}

// Serve a new connection.
func (c *conn) serve() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("diam: panic serving %v: %v\n%s",
				c.rwc.RemoteAddr().String(), err, buf)
		}
		c.rwc.Close()
	}()
	if tlsConn, ok := c.rwc.(*tls.Conn); ok {
		if err := tlsConn.Handshake(); err != nil {
			return
		}
		c.tlsState = &tls.ConnectionState{}
		*c.tlsState = tlsConn.ConnectionState()
	}
	for {
		m, err := c.readMessage()
		if err != nil {
			c.rwc.Close()
			// Report errors to the channel, except EOF.
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				h := c.server.Handler
				if h == nil {
					h = DefaultServeMux
				}
				if er, ok := h.(ErrorReporter); ok {
					er.Error(&ErrorReport{c.writer, m, err})
				}
			}
			break
		}
		// Handle messages in this goroutine.
		serverHandler{c.server}.ServeDIAM(c.writer, m)
	}
}

// dictionary returns the dictionary parser associated to the Server instance
// or dict.Default.
func (c *conn) dictionary() *dict.Parser {
	if c.server.Dict == nil {
		return dict.Default
	}
	return c.server.Dict
}

// A response represents the server side of a diameter response.
// It implements the Conn and CloseNotifier interfaces.
type response struct {
	mu   sync.Mutex      // guards conn and Write
	conn *conn           // socket, reader and writer
	xmu  sync.Mutex      // guards ctx
	ctx  context.Context // context for this Conn
}

// Write writes the message m to the connection.
func (w *response) Write(b []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.conn.server.WriteTimeout > 0 {
		w.conn.rwc.SetWriteDeadline(time.Now().Add(w.conn.server.WriteTimeout))
	}
	msc, isMulti := w.conn.rwc.(MultistreamConn) // Note - SetWriteDeadline is not currently supported for SCTP
	if isMulti {                                 // don't use buffered writer for muti-streamming writes it'll mix up streams
		return msc.Write(b)
	}
	n, err := w.conn.buf.Writer.Write(b)
	if err != nil {
		return 0, err
	}
	if err = w.conn.buf.Writer.Flush(); err != nil {
		return 0, err
	}
	return n, nil
}

// WriteStream of MultistreamWriter interface
func (w *response) WriteStream(b []byte, stream uint) (int, error) {
	// TODO - SetWriteDeadline is not currently supported
	if msc, isMulti := w.conn.rwc.(MultistreamConn); isMulti {
		// don't use buffered writer for muti-streamming writes it'll mix up streams
		return msc.WriteStream(b, stream)
	}
	return w.Write(b)
}

// CurrentWriterStream of MultistreamWriter interface
func (w *response) CurrentWriterStream() uint {
	if msc, isMulti := w.conn.rwc.(MultistreamConn); isMulti {
		return msc.CurrentWriterStream()
	}
	return 0
}

// ResetWriterStream of MultistreamWriter interface
func (w *response) ResetWriterStream() {
	if msc, isMulti := w.conn.rwc.(MultistreamConn); isMulti {
		msc.CurrentWriterStream()
	}
}

// SetWriterStream of MultistreamWriter interface
func (w *response) SetWriterStream(stream uint) uint {
	if msc, isMulti := w.conn.rwc.(MultistreamConn); isMulti {
		return msc.SetWriterStream(stream)
	}
	return 0
}

// Close closes the connection.
func (w *response) Close() {
	w.conn.rwc.Close()
}

// LocalAddr returns the local address of the connection.
func (w *response) LocalAddr() net.Addr {
	return w.conn.rwc.LocalAddr()
}

// RemoteAddr returns the peer address of the connection.
func (w *response) RemoteAddr() net.Addr {
	return w.conn.rwc.RemoteAddr()
}

// TLS returns the TLS connection state, or nil.
func (w *response) TLS() *tls.ConnectionState {
	return w.conn.tlsState
}

// Dictionary returns the dictionary parser associated to this connection.
// If none was provided then it returns the default dictionary.
func (w *response) Dictionary() *dict.Parser {
	return w.conn.dictionary()
}

// CloseNotify implements the CloseNotifier interface.
func (w *response) CloseNotify() <-chan struct{} {
	return w.conn.closeNotify()
}

// Context returns the internal context or a new context.Background.
func (w *response) Context() context.Context {
	w.xmu.Lock()
	defer w.xmu.Unlock()
	if w.ctx == nil {
		w.ctx = context.Background()
	}
	return w.ctx
}

// SetContext replaces the internal context with the given one.
func (w *response) SetContext(ctx context.Context) {
	w.xmu.Lock()
	w.ctx = ctx
	w.xmu.Unlock()
}

func (w *response) Connection() net.Conn {
	return w.conn.rwc
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as diameter handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(Conn, *Message)

// ServeDIAM calls f(c, m).
func (f HandlerFunc) ServeDIAM(c Conn, m *Message) {
	f(c, m)
}

// The ErrorReporter interface is implemented by Handlers that
// allow reading errors from the underlying connection, like
// parsing diameter messages or connection errors.
type ErrorReporter interface {
	// Error writes an error to the reporter.
	Error(err *ErrorReport)

	// ErrorReports returns a channel that receives
	// errors from the connection.
	ErrorReports() <-chan *ErrorReport
}

// ErrorReport is sent out of the server in case it fails to
// read messages due to a bad dictionary or network errors.
type ErrorReport struct {
	Conn    Conn     // Peer that caused the error
	Message *Message // Message that caused the error
	Error   error    // Error message
}

// String returns an error message. It does not render the Message field.
func (er *ErrorReport) String() string {
	if er.Conn == nil {
		return fmt.Sprintf("diameter error: %s", er.Error)
	}
	return fmt.Sprintf("diameter error on %s: %s", er.Conn.RemoteAddr(), er.Error)
}

// ServeMux is a diameter message multiplexer. It matches the
// command from the incoming message against a list of
// registered commands and calls the handler.
type ServeMux struct {
	e      chan *ErrorReport
	mu     sync.RWMutex // Guards m.
	m      map[string]muxEntry
	idxMap map[CommandIndex]muxEntry
}

type muxEntry struct {
	h      Handler
	cmd    string
	cmdIdx CommandIndex
}

type CommandIndex struct {
	AppID   uint32
	Code    uint32
	Request bool
}

var ALL_CMD_INDEX = CommandIndex{^uint32(0), ^uint32(0), false}

// NewServeMux allocates and returns a new ServeMux.
func NewServeMux() *ServeMux {
	return &ServeMux{
		e:      make(chan *ErrorReport, 1),
		m:      make(map[string]muxEntry),
		idxMap: make(map[CommandIndex]muxEntry),
	}
}

// DefaultServeMux is the default ServeMux used by Serve.
var DefaultServeMux = NewServeMux()

// Error implements the ErrorReporter interface.
func (mux *ServeMux) Error(err *ErrorReport) {
	select {
	case mux.e <- err:
	default:
	}
}

// ErrorReports implement the ErrorReporter interface.
func (mux *ServeMux) ErrorReports() <-chan *ErrorReport {
	return mux.e
}

// ServeDIAM dispatches the request to the handler that match the code
// in the incoming message. If the special "ALL" handler is registered
// it is used as a catch-all. Otherwise an ErrorReport is sent out.
func (mux *ServeMux) ServeDIAM(c Conn, m *Message) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	dcmd, err := m.Dictionary().FindCommand(
		m.Header.ApplicationID,
		m.Header.CommandCode)

	if err != nil {
		// Try the catch-all.
		mux.serveIdx(ALL_CMD_INDEX, c, m)
		return
	}

	idx := CommandIndex{
		m.Header.ApplicationID,
		m.Header.CommandCode,
		m.Header.CommandFlags&RequestFlag == RequestFlag}
	_, ok := mux.idxMap[idx]
	if ok {
		mux.serveIdx(idx, c, m)
		return
	}

	var cmd string
	if m.Header.CommandFlags&RequestFlag == RequestFlag {
		cmd = dcmd.Short + "R"
	} else {
		cmd = dcmd.Short + "A"
	}
	mux.serve(cmd, c, m)
}

func (mux *ServeMux) serveIdx(cmd CommandIndex, c Conn, m *Message) {
	entry, ok := mux.idxMap[cmd]
	if ok {
		entry.h.ServeDIAM(c, m)
		return
	}
	// Try catch-all.
	entry, ok = mux.idxMap[ALL_CMD_INDEX]
	if ok {
		entry.h.ServeDIAM(c, m)
		return
	}
	mux.Error(&ErrorReport{
		Conn:    c,
		Message: m,
		Error:   fmt.Errorf("unhandled message for index: %+v", cmd),
	})
}

func (mux *ServeMux) serve(cmd string, c Conn, m *Message) {
	entry, ok := mux.m[cmd]
	if ok {
		entry.h.ServeDIAM(c, m)
		return
	}
	// Try catch-all.
	entry, ok = mux.idxMap[ALL_CMD_INDEX]
	if ok {
		entry.h.ServeDIAM(c, m)
		return
	}
	mux.Error(&ErrorReport{
		Conn:    c,
		Message: m,
		Error:   fmt.Errorf("unhandled message for '%s'", cmd),
	})
}

// Handle registers the handler for the given code.
// If a handler already exists for code, Handle panics.
func (mux *ServeMux) Handle(shortCmd string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	if handler == nil {
		panic("DIAM: nil handler")
	}
	if shortCmd == "ALL" {
		mux.idxMap[ALL_CMD_INDEX] = muxEntry{h: handler, cmd: shortCmd}
		return
	}
	mux.m[shortCmd] = muxEntry{h: handler, cmd: shortCmd}
}

// Handle registers the handler for the given code.
// If a handler already exists for code, Handle panics.
func (mux *ServeMux) HandleIdx(cmd CommandIndex, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	if handler == nil {
		panic("DIAM: nil handler")
	}
	mux.idxMap[cmd] = muxEntry{h: handler, cmdIdx: cmd}
}

// HandleFunc registers the handler function for the given command.
// Special cmd "ALL" may be used as a catch all.
func (mux *ServeMux) HandleFunc(cmd string, handler func(Conn, *Message)) {
	mux.Handle(cmd, HandlerFunc(handler))
}

// Handle registers the handler object for the given command
// in the DefaultServeMux.
func Handle(cmd string, handler Handler) {
	DefaultServeMux.Handle(cmd, handler)
}

// HandleFunc registers the handler function for the given command
// in the DefaultServeMux.
func HandleFunc(cmd string, handler func(Conn, *Message)) {
	DefaultServeMux.HandleFunc(cmd, handler)
}

// ErrorReports returns the ErrorReport channel of the DefaultServeMux.
func ErrorReports() <-chan *ErrorReport {
	return DefaultServeMux.ErrorReports()
}

// Serve accepts incoming diameter connections on the listener l,
// creating a new service goroutine for each.  The service goroutines
// read messages and then call handler to reply to them.
// Handler is typically nil, in which case the DefaultServeMux is used.
func Serve(l net.Listener, handler Handler) error {
	srv := &Server{Handler: handler}
	return srv.Serve(l)
}

// A Server defines parameters for running a diameter server.
type Server struct {
	Network      string        // network of the address - empty string defaults to tcp
	Addr         string        // address to listen on, ":3868" if empty
	Handler      Handler       // handler to invoke, DefaultServeMux if nil
	Dict         *dict.Parser  // diameter dictionaries for this server
	ReadTimeout  time.Duration // maximum duration before timing out read of the request
	WriteTimeout time.Duration // maximum duration before timing out write of the response
	TLSConfig    *tls.Config   // optional TLS config, used by ListenAndServeTLS
	LocalAddr    net.Addr      // optional Local Address to bind dailer's (Dail...) socket to
}

// serverHandler delegates to either the server's Handler or DefaultServeMux.
type serverHandler struct {
	srv *Server
}

func (sh serverHandler) ServeDIAM(w Conn, m *Message) {
	handler := sh.srv.Handler
	if handler == nil {
		handler = DefaultServeMux
	}
	handler.ServeDIAM(w, m)
}

// ListenAndServe listens on the network address srv.Addr and then
// calls Serve to handle requests on incoming connections.  If
//
// If srv.Network is blank, "tcp" is used
// If srv.Addr is blank, ":3868" is used.
func (srv *Server) ListenAndServe() error {
	network := srv.Network
	if len(network) == 0 {
		network = "tcp"
	}
	addr := srv.Addr
	if len(addr) == 0 {
		addr = ":3868"
	}
	l, e := MultistreamListen(network, addr)
	if e != nil {
		return e
	}
	return srv.Serve(l)
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("diam: accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			network := "<nil>"
			address := network
			addr := l.Addr()
			if addr != nil {
				network = addr.Network()
				address = addr.String()
			}
			log.Printf("diam: accept error: %v for %s %s", e, network, address)
			return e
		}
		tempDelay = 0
		if c, err := srv.newConn(rw); err != nil {
			log.Printf("srv.newConn error: %v", err)
			continue
		} else {
			go c.serve()
		}
	}
}

// ListenAndServeNetwork listens on the network & addr
// and then calls Serve with handler to handle requests
// on incoming connections.
//
// If handler is nil, DefaultServeMux is used.
//
// If dict is nil, dict.Default is used.
func ListenAndServeNetwork(network, addr string, handler Handler, dp *dict.Parser) error {
	server := &Server{Network: network, Addr: addr, Handler: handler, Dict: dp}
	return server.ListenAndServe()
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests
// on incoming connections.
//
// If handler is nil, DefaultServeMux is used.
//
// If dict is nil, dict.Default is used.
func ListenAndServe(addr string, handler Handler, dp *dict.Parser) error {
	return ListenAndServeNetwork("tcp", addr, handler, dp)
}

// ListenAndServeTLS listens on the network address srv.Addr and
// then calls Serve to handle requests on incoming TLS connections.
//
// Filenames containing a certificate and matching private key for
// the server must be provided. If the certificate is signed by a
// certificate authority, the certFile should be the concatenation
// of the server's certificate followed by the CA's certificate.
//
// If srv.Network is blank, "tcp" is used
// If srv.Addr is blank, ":3868" is used.
func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	network := srv.Network
	if len(network) == 0 {
		network = "tcp"
	}
	addr := srv.Addr
	if len(addr) == 0 {
		addr = ":3868"
	}
	var config *tls.Config
	if srv.TLSConfig == nil {
		config = new(tls.Config)
	} else {
		config = TLSConfigClone(srv.TLSConfig)
	}
	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}
	conn, err := Listen(network, addr)
	if err != nil {
		return err
	}
	tlsListener := tls.NewListener(conn, config)
	return srv.Serve(tlsListener)
}

// ListenAndServeNetworkTLS acts identically to ListenAndServeNetwork, except that it
// expects SSL connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate followed by the CA's certificate.
//
// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
func ListenAndServeNetworkTLS(network, addr string, certFile string, keyFile string, handler Handler, dp *dict.Parser) error {
	server := &Server{Network: network, Addr: addr, Handler: handler, Dict: dp}
	return server.ListenAndServeTLS(certFile, keyFile)
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects SSL connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate followed by the CA's certificate.
//
// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler, dp *dict.Parser) error {
	return ListenAndServeNetworkTLS("tcp", addr, certFile, keyFile, handler, dp)
}
