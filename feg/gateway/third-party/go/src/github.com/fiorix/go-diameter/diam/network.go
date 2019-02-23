package diam

import (
	"io"
	"net"
	"time"

	"github.com/ishidawataru/sctp"
)

// InvalidStreamID - const used for unspecified/uninitialized stream #
const InvalidStreamID = ^uint(0)

// MultistreamReader - reader interface for multi-stream protocols
type MultistreamReader interface {
	io.Reader
	// ReadAny reads data from any connection's stream (if available).
	// Returns number of bytes read and the stream number
	ReadAny(b []byte) (n int, stream uint, err error)
	// ReadStream reads data from the specified connection's stream
	ReadStream(b []byte, stream uint) (n int, err error)
	// ReadAtLeast reads into b from the stream until it has read at least min bytes.
	// It returns the number of bytes copied and an error if fewer bytes were read.
	// If min is greater than the length of buf, ReadAtLeast returns ErrShortBuffer.
	ReadAtLeast(b []byte, min int, strm uint) (n int, stream uint, err error)
	// CurrentStream returns the last stream read by Read adaptor
	CurrentStream() uint
	// ResetCurrentStream resets current read stream so the next Read adaptor call will read from any stream
	ResetCurrentStream()
	// SetCurrentStream sets current read stream so the next Read adaptor call will be forced to read from this stream
	SetCurrentStream(uint) uint
}

// MultistreamWriter - writer interface for multi-stream protocols
type MultistreamWriter interface {
	io.Writer
	// WriteStream writes data to the connection's stream
	WriteStream(b []byte, stream uint) (n int, err error)
	// CurrentWriterStream returns the stream that the next call to Write adaptor will be used for writing
	CurrentWriterStream() uint
	// ResetWriterStream resets current read stream so the next Write adaptor call
	// will use either the current Read stream (if set) or the protocol specific default stream
	ResetWriterStream()
	// SetWriterStream sets current write stream so the next Write adaptor call will be forced to read from this stream
	SetWriterStream(uint) uint
}

// MutistreamConnErrorHandler is a handler for Multi stream connection to be called on read errors if needed
type MutistreamConnErrorHandler func(MultistreamConn, error)

// MultistreamConn provides interface for multi-streamed association
// Alone with stream specific read/write functionality, MultistreamConn also provides
// Read & Write "adaptor" methods which comply with io.Reader, io.Writer interfaces and ensure read/write stream
// continuity for a single Read/Write
// Implements net.Conn, MultistreamReader & MultistreamWriter
type MultistreamConn interface {
	net.Conn
	//
	// MultistreamReader interface
	//
	// ReadAny reads data from any connection's stream (if available).
	// Returns number of bytes read and the stream number
	ReadAny(b []byte) (n int, stream uint, err error)
	// ReadStream reads data from the specified connection's stream
	ReadStream(b []byte, stream uint) (n int, err error)
	// ReadAtLeast reads into b from the stream until it has read at least min bytes.
	// It returns the number of bytes copied and an error if fewer bytes were read.
	// If min is greater than the length of buf, ReadAtLeast returns ErrShortBuffer.
	ReadAtLeast(b []byte, min int, strm uint) (n int, stream uint, err error)
	// CurrentStream returns the last stream read by Read adaptor
	CurrentStream() uint
	// ResetCurrentStream resets current read stream so the next Read adaptor call will read from any stream
	ResetCurrentStream()
	// SetCurrentStream sets current read stream so the next Read adaptor call will be forced to read from this stream
	SetCurrentStream(uint) uint
	//
	// MultistreamWriter interface
	//
	// WriteStream writes data to the connection's stream
	WriteStream(b []byte, stream uint) (n int, err error)
	// CurrentWriterStream returns the stream that the next call to Write adaptor will be used for writing
	CurrentWriterStream() uint
	// ResetWriterStream resets current write stream so the next Write adaptor call
	// will use either the current Read stream (if set) or the protocol specific default stream to write to
	ResetWriterStream()
	// SetWriterStream sets current write stream so the next Write adaptor call will be forced to write to this stream
	SetWriterStream(uint) uint
	//
	// SetErrorHandler sets reader error notification handler, it'll be called on any read IO error for the connection
	SetErrorHandler(MutistreamConnErrorHandler)
}

// Dialer interface, see https://golang.org/pkg/net/#Dialer.Dial
type Dialer interface {
	// Dial connects to the address on the named network
	Dial(network, address string) (net.Conn, error)
}

func getDialer(network string, timeout time.Duration, laddr net.Addr) Dialer {
	switch network {
	case "sctp", "sctp4", "sctp6":
		la, _ := laddr.(*sctp.SCTPAddr)
		return sctpSingleStreamDialer{LocalAddr: la}
	default:
		return &net.Dialer{Timeout: timeout, LocalAddr: laddr}
	}
}

// getMultistreamDialer returns Dailer with multistreaming support when appropriate for the network/protocol
func getMultistreamDialer(network string, timeout time.Duration, laddr net.Addr) Dialer {
	switch network {
	case "sctp", "sctp4", "sctp6":
		la, _ := laddr.(*sctp.SCTPAddr)
		return sctpDialer{LocalAddr: la}
	default:
		return &net.Dialer{Timeout: timeout, LocalAddr: laddr}
	}
}

func resolveAddress(network, addr string) (net.Addr, error) {
	switch network {
	case "sctp", "sctp4", "sctp6":
		return sctp.ResolveSCTPAddr(network, addr)
	case "":
		network = "tcp"
		fallthrough
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, addr)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}

func listenSCTP(network, address string) (*sctp.SCTPListener, error) {
	sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
	if err != nil {
		return nil, err
	}
	return sctp.ListenSCTPExt(
		network,
		sctpAddr,
		sctp.InitMsg{
			NumOstreams:  MaxOutboundSCTPStreams,
			MaxInstreams: MaxInboundSCTPStreams})
}

// Listen announces on the local network address
func Listen(network, address string) (net.Listener, error) {
	switch network {
	case "sctp", "sctp4", "sctp6":
		return listenSCTP(network, address)
	default:
		return net.Listen(network, address)
	}
}

// MultistreamListen returns Listener with multistreaming support when appropriate for the network/protocol
func MultistreamListen(network, address string) (net.Listener, error) {
	switch network {
	case "sctp", "sctp4", "sctp6":
		lis, err := listenSCTP(network, address)
		return sctpListener{lis}, err
	default:
		return net.Listen(network, address)
	}
}
