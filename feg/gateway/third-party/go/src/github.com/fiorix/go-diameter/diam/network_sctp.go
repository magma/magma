package diam

import (
	"bytes"
	"container/heap"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/ishidawataru/sctp"
)

const (
	// MaxInboundSCTPStreams - max inbound streams default for new connections.
	// see https://tools.ietf.org/html/rfc4960#page-25
	MaxInboundSCTPStreams = 16

	// MaxOutboundSCTPStreams - max outbound streams default for new connections.
	MaxOutboundSCTPStreams = MaxInboundSCTPStreams

	// DiameterPPID - SCTP Payload Protocol Identifier for Diameter
	// see: https://tools.ietf.org/html/rfc4960#section-14.4 and https://tools.ietf.org/html/rfc6733#page-24
	DiameterPPID uint32 = 46
)

type sctpDialer struct {
	LocalAddr *sctp.SCTPAddr
}

// sctpSingleStreamDialer - SCTP Dialer for stream unaware applications.
type sctpSingleStreamDialer sctpDialer

type sctpListener struct {
	*sctp.SCTPListener
}

type streamBuffer struct {
	*bytes.Buffer
	stream uint
	idx    int
}

type streams struct {
	streamMap  map[uint]*streamBuffer
	streamHeap []*streamBuffer
}

// Go Heap interface implementation
func (pq *streams) Len() int {
	return len(pq.streamHeap)
}

func (pq *streams) Less(i, j int) bool {
	return pq.streamHeap[i].Len() > pq.streamHeap[j].Len()
}

func (pq *streams) Swap(i, j int) {
	pq.streamHeap[i], pq.streamHeap[j] = pq.streamHeap[j], pq.streamHeap[i]
	pq.streamHeap[i].idx = i
	pq.streamHeap[j].idx = j
}

func (pq *streams) Push(x interface{}) {
	sb := x.(*streamBuffer)
	if _, ok := pq.streamMap[sb.stream]; ok {
		panic("pushing an existing stream " + strconv.Itoa(int(sb.stream)))
	}
	sb.idx = len(pq.streamHeap)
	pq.streamHeap = append(pq.streamHeap, sb)
	if pq.streamMap == nil {
		pq.streamMap = map[uint]*streamBuffer{sb.stream: sb}
	} else {
		pq.streamMap[sb.stream] = sb
	}
}

func (pq *streams) Pop() interface{} {
	n := len(pq.streamHeap)
	sb := pq.streamHeap[n-1]
	delete(pq.streamMap, sb.stream)
	pq.streamHeap = pq.streamHeap[0 : n-1]
	sb.idx = -1
	return sb
}

// SCTPConn - MutistreamConn implementation for SCTP.
type SCTPConn struct {
	*sctp.SCTPConn

	streamBuffMu sync.Mutex
	s            *streams
	mu           sync.RWMutex
	currStream   uint
	wmu          sync.RWMutex
	writerStream uint
	errorHandler MutistreamConnErrorHandler
}

// NewSCTPConn - creates new MultistreamConn (diam.SCTPConn) from provided sctp.SCTPConn.
func NewSCTPConn(sctpConn *sctp.SCTPConn) MultistreamConn {
	if sctpConn == nil {
		return nil
	}
	sctpConn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
	return &SCTPConn{SCTPConn: sctpConn, s: &streams{}, currStream: InvalidStreamID, writerStream: InvalidStreamID}
}

// ReadAny reads data from any association's stream (if available).
// Returns number of bytes read and the stream number.
func (msc *SCTPConn) ReadAny(b []byte) (n int, stream uint, err error) {
	// First see if we can consume an existing, previously received buffer
	msc.streamBuffMu.Lock()
	// Consume the longest buffer first
	if msc.s.Len() > 0 && msc.s.streamHeap[0].Len() > 0 {
		sb := msc.s.streamHeap[0]
		n, err = sb.Read(b)
		stream = sb.stream
		heap.Fix(msc.s, 0)
		msc.streamBuffMu.Unlock()
		return
	}
	msc.streamBuffMu.Unlock()
	var info *sctp.SndRcvInfo
	n, info, err = msc.SCTPRead(b)
	if n < 0 {
		n = 0
	}
	if info != nil {
		stream = uint(info.Stream)
	} else if n > 0 { // reset current stream only if there was some data received
		stream = InvalidStreamID
	}
	// shortcut for non empty stream buffer
	n, err = msc.verifyStreamBuff(b, n, stream, err)
	if err != nil {
		hptr := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&msc.errorHandler)))
		if hptr != nil {
			(*(*MutistreamConnErrorHandler)(unsafe.Pointer(&hptr)))(msc, err)
		}
	}
	return
}

// SetErrorHandler sets reader error notification handler,
// it'll be called on any read IO error for the connection.
func (msc *SCTPConn) SetErrorHandler(h MutistreamConnErrorHandler) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&msc.errorHandler)), *(*unsafe.Pointer)(unsafe.Pointer(&h)))
}

// ReadStream reads data from the specified association's stream.
func (msc *SCTPConn) ReadStream(b []byte, stream uint) (n int, err error) {
	var info *sctp.SndRcvInfo
	var currStream uint
	msc.streamBuffMu.Lock()
	for {
		// First see if we can consume an existing, previously received stream buffer
		if sb, ok := msc.s.streamMap[stream]; ok && sb.Len() > 0 {
			n, err = sb.Read(b)
			heap.Fix(msc.s, sb.idx)
			msc.streamBuffMu.Unlock()
			return
		}

		msc.streamBuffMu.Unlock()
		n, info, err = msc.SCTPRead(b)
		if n <= 0 {
			return 0, err
		}

		if info == nil {
			// info == nil => no stream info, the socket was not initialized properly, assign to InvalidStreamID stream
			currStream = InvalidStreamID
		} else {
			currStream = uint(info.Stream)
			// shortcut for non empty stream buffer
			if currStream == stream {
				// We got data for the requested stream, but we still need to make sure that the stream buffer didn't
				// get any data buffered into while we were waiting on SCTPRead outside of lock
				return msc.verifyStreamBuff(b, n, stream, err)
			}
		}
		rb := b[0:n]
		msc.streamBuffMu.Lock()
		msc.bufferStreamData(rb, currStream)
	}
}

// ReadAtLeast reads into b from the stream until it has read at least min bytes.
// It returns the number of bytes copied and an error if fewer bytes were read.
// If min is greater than the length of buf, ReadAtLeast returns ErrShortBuffer.
func (msc *SCTPConn) ReadAtLeast(buf []byte, min int, strm uint) (n int, stream uint, err error) {
	if len(buf) < min {
		return 0, InvalidStreamID, io.ErrShortBuffer
	}
	var nn int
	if strm == InvalidStreamID {
		nn, stream, err = msc.ReadAny(buf)
		n += nn
	} else {
		stream = strm
	}
	for n < min && err == nil {
		nn, err = msc.ReadStream(buf[n:], stream)
		n += nn
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

// verifyStreamBuff checks is there is a ready buffered data for the stream already,
// pipes b through the same buffer if there is and re-reads b from the beginning of the buffer.
func (msc *SCTPConn) verifyStreamBuff(b []byte, n int, stream uint, currErr error) (int, error) {
	msc.streamBuffMu.Lock()
	defer msc.streamBuffMu.Unlock()
	if sb, ok := msc.s.streamMap[stream]; ok && sb.Len() > 0 {
		sb.Write(b[0:n])
		heap.Fix(msc.s, sb.idx)
		return sb.Read(b)
	}
	return n, currErr
}

// bufferStreamData buffers b into the corresponding stream buffer.
func (msc *SCTPConn) bufferStreamData(b []byte, stream uint) {
	sb, ok := msc.s.streamMap[stream]
	if ok {
		sb.Write(b)
		heap.Fix(msc.s, sb.idx)
	} else {
		sb = &streamBuffer{Buffer: new(bytes.Buffer), stream: stream}
		sb.Write(b)
		heap.Push(msc.s, sb)
	}
}

// WriteStream writes data to the association's stream.
func (msc *SCTPConn) WriteStream(b []byte, stream uint) (int, error) {
	info := &sctp.SndRcvInfo{PPID: DiameterPPID}
	if stream != InvalidStreamID {
		info.Stream = uint16(stream)
	}
	return msc.SCTPWrite(b, info)
}

// CurrentStream returns the last stream read by Read 'adaptor'.
func (msc *SCTPConn) CurrentStream() uint {
	msc.mu.RLock()
	defer msc.mu.RUnlock()
	return msc.currStream
}

// ResetCurrentStream resets current read stream so the next Read adaptor call will read from any stream.
func (msc *SCTPConn) ResetCurrentStream() {
	msc.mu.Lock()
	msc.currStream = InvalidStreamID
	msc.mu.Unlock()
}

// SetCurrentStream sets current read stream so the next Read adaptor call will be forced to read from it.
func (msc *SCTPConn) SetCurrentStream(stream uint) uint {
	msc.mu.Lock()
	stream, msc.currStream = msc.currStream, stream
	msc.mu.Unlock()
	return stream
}

// CurrentWriterStream returns the stream that the next call to Write adaptor will be used for writing.
func (msc *SCTPConn) CurrentWriterStream() uint {
	msc.wmu.RLock()
	defer msc.wmu.RUnlock()
	return msc.writerStream
}

// ResetWriterStream resets current write stream so the next Write adaptor call
// will use either the current Read stream (if set) or the protocol specific
// default stream to write to.
func (msc *SCTPConn) ResetWriterStream() {
	msc.wmu.Lock()
	msc.writerStream = InvalidStreamID
	msc.wmu.Unlock()
}

// SetWriterStream resets current write stream so the next Write adaptor call
// will use either the current Write stream (if set), current Read stream (if set)
// or the protocol specific default stream.
func (msc *SCTPConn) SetWriterStream(stream uint) uint {
	msc.wmu.Lock()
	stream, msc.writerStream = msc.writerStream, stream
	msc.wmu.Unlock()
	return stream
}

// Read 'adaptor' reads data from the specified association while maintaining stream continuity.
// Read implements io.Reader interface on multi-stream protocols for stream unaware (legacy) applications.
// Read will guarantee that every read into a single buffer 'b' is sourced from a single stream and that
// multiple consecutive & concurrent reads will continue from the same stream between calls to ResetCurrentStream or
// SetCurrentStream.
// Calls to Read & ResetCurrentStream/SetCurrentStream should be synchronized
func (msc *SCTPConn) Read(b []byte) (n int, err error) {
	var strm uint
	msc.mu.RLock() // block changes to currStream during a single read

	for {
		if msc.currStream == InvalidStreamID { // first read after reset]
			n, strm, err = msc.ReadAny(b)
			if err == nil && msc.currStream != strm {
				msc.mu.RUnlock()
				msc.mu.Lock()
				if msc.currStream == InvalidStreamID { // msc.currStream may have changed, check again
					msc.currStream = strm
				} else if msc.currStream != strm {
					// Worst case scenario - concurrent read completed first and received from a different stream
					msc.bufferStreamData(b[0:n], strm)
					msc.mu.Unlock()
					msc.mu.RLock()
					continue
				}
				msc.mu.Unlock()
			} else {
				msc.mu.RUnlock()
			}
			return
		}
		// Stream # is already set, keep reading from it until next reset
		strm = msc.currStream
		n, err = msc.ReadStream(b, strm)
		msc.mu.RUnlock()
		return
	}
}

// Write writes data to the association while maintaining stream continuity.
// The write stream will be selected in the following order:
//   1) If current write stream is set (is not InvalidStreamID), it'll be used for writing.
//   2) If current read stream is set, it'll be used for writing.
//   3) If neither current write nor current read streams are set, write will use default
//      protocol stream (0 for the current SCTP implementation).
func (msc *SCTPConn) Write(b []byte) (int, error) {
	msc.wmu.RLock() // block changes to msc.writerStream during write
	defer msc.wmu.RUnlock()

	stream := msc.writerStream
	// If writer stream is set by a user, stick with it
	if stream == InvalidStreamID {
		stream = msc.CurrentStream()
	}
	info := &sctp.SndRcvInfo{PPID: DiameterPPID}

	// If writer stream is not set, stick with the reader stream #
	if stream != InvalidStreamID {
		info.Stream = uint16(stream)
	}
	return msc.SCTPWrite(b, info)
}

// Dial connects to the address on the named SCTP network.
func (d sctpDialer) Dial(network, address string) (net.Conn, error) {
	sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
	if err != nil {
		return nil, err
	}

	conn, err := sctp.DialSCTPExt(
		network,
		d.LocalAddr,
		sctpAddr,
		sctp.InitMsg{
			NumOstreams:  MaxOutboundSCTPStreams,
			MaxInstreams: MaxInboundSCTPStreams})
	return NewSCTPConn(conn), err
}

// Dial - SCTP dial for stream unaware apps.
func (d sctpSingleStreamDialer) Dial(network, address string) (net.Conn, error) {
	sctpAddr, err := sctp.ResolveSCTPAddr(network, address)
	if err != nil {
		return nil, err
	}

	return sctp.DialSCTPExt(
		network,
		d.LocalAddr,
		sctpAddr,
		sctp.InitMsg{
			NumOstreams:  MaxOutboundSCTPStreams,
			MaxInstreams: MaxInboundSCTPStreams})
}

// Accept implements the Accept method in the listener interface for sctpListener (see: MultistreamListen).
func (l sctpListener) Accept() (net.Conn, error) {
	conn, err := l.AcceptSCTP()
	return NewSCTPConn(conn), err
}
