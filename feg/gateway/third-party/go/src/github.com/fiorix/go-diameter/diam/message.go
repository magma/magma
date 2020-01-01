// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// MessageBufferLength is the default buffer length for Diameter messages.
var MessageBufferLength = 1 << 10

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Message represents a Diameter message.
type Message struct {
	Header *Header
	AVP    []*AVP // AVPs in this message.

	// dictionary parser object used to encode and decode AVPs.
	dictionary *dict.Parser
	stream     uint // the stream this message was received on (if any)
}

var readerBufferPool sync.Pool

func newReaderBuffer() *bytes.Buffer {
	if v := readerBufferPool.Get(); v != nil {
		return v.(*bytes.Buffer)
	}
	return bytes.NewBuffer(make([]byte, MessageBufferLength))
}

func putReaderBuffer(b *bytes.Buffer) {
	if cap(b.Bytes()) == MessageBufferLength {
		b.Reset()
		readerBufferPool.Put(b)
	}
}

func readerBufferSlice(buf *bytes.Buffer, l int) []byte {
	b := buf.Bytes()
	if l <= MessageBufferLength && cap(b) >= MessageBufferLength {
		return b[:l]
	}
	return make([]byte, l)
}

// ReadMessage reads a binary stream from the reader and uses the given
// dictionary to parse it.
func ReadMessage(reader io.Reader, dictionary *dict.Parser) (*Message, error) {
	buf := newReaderBuffer()
	defer putReaderBuffer(buf)
	m := &Message{dictionary: dictionary}
	cmd, stream, err := m.readHeader(reader, buf)
	if err != nil {
		return nil, err
	}
	m.stream = stream
	if err = m.readBody(reader, buf, cmd, stream); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageStream returns the stream #, the message was received on (when applicable)
func (m *Message) MessageStream() uint {
	return m.stream
}

func (m *Message) readHeader(r io.Reader, buf *bytes.Buffer) (cmd *dict.Command, stream uint, err error) {
	b := buf.Bytes()[:HeaderLength]
	msr, isMulti := r.(MultistreamReader)
	if isMulti {
		_, stream, err = msr.ReadAtLeast(b, HeaderLength, InvalidStreamID)
		if err == nil {
			msr.SetCurrentStream(stream)
		}
	} else {
		_, err = io.ReadFull(r, b)
	}
	if err != nil {
		return nil, stream, err
	}
	m.Header, err = DecodeHeader(b)
	if err != nil {
		return nil, stream, err
	}
	cmd, err = m.Dictionary().FindCommand(
		m.Header.ApplicationID,
		m.Header.CommandCode,
	)
	if err != nil {
		return nil, stream, err
	}
	return cmd, stream, nil
}

func (m *Message) readBody(r io.Reader, buf *bytes.Buffer, cmd *dict.Command, stream uint) error {
	var err error
	var n int
	b := readerBufferSlice(buf, int(m.Header.MessageLength-HeaderLength))
	msr, isMulti := r.(MultistreamReader)
	if isMulti {
		n, _, err = msr.ReadAtLeast(b, len(b), stream)
	} else {
		n, err = io.ReadFull(r, b)
	}
	if err != nil {
		return fmt.Errorf("readBody Error: %v, %d bytes read", err, n)
	}
	n = m.maxAVPsFor(cmd)
	if n == 0 {
		// TODO: fail to load the dictionary instead.
		return fmt.Errorf(
			"Command %s (%d) has no AVPs defined in the dictionary.",
			cmd.Name, cmd.Code)
	}
	// Pre-allocate max # of AVPs for this message.
	m.AVP = make([]*AVP, 0, n)
	if err = m.decodeAVPs(b); err != nil {
		return err
	}
	return nil
}

func (m *Message) maxAVPsFor(cmd *dict.Command) int {
	if m.Header.CommandFlags&RequestFlag == RequestFlag {
		return len(cmd.Request.Rule)
	}
	return len(cmd.Answer.Rule)
}

func (m *Message) decodeAVPs(b []byte) error {
	var a *AVP
	var err error
	for n := 0; n < len(b); {
		a, err = DecodeAVP(b[n:], m.Header.ApplicationID, m.Dictionary())
		if err != nil {
			return fmt.Errorf("Failed to decode AVP: %s", err)
		}
		m.AVP = append(m.AVP, a)
		n += a.Len()
	}
	return nil
}

// NewMessage creates and initializes a Message.
func NewMessage(cmd uint32, flags uint8, appid, hopbyhop, endtoend uint32, dictionary *dict.Parser) *Message {
	if hopbyhop == 0 {
		hopbyhop = rand.Uint32()
	}
	if endtoend == 0 {
		endtoend = rand.Uint32()
	}
	return &Message{
		Header: &Header{
			Version:       1,
			MessageLength: HeaderLength,
			CommandFlags:  flags,
			CommandCode:   cmd,
			ApplicationID: appid,
			HopByHopID:    hopbyhop,
			EndToEndID:    endtoend,
		},
		dictionary: dictionary,
		stream:     InvalidStreamID,
	}
}

// NewRequest creates a new Message with the Request bit set.
func NewRequest(cmd uint32, appid uint32, dictionary *dict.Parser) *Message {
	return NewMessage(cmd, RequestFlag, appid, 0, 0, dictionary)
}

// Dictionary returns the dictionary parser object associated with this
// message. This dictionary is used to encode and decode the message.
// If no dictionary is associated then it returns the default dictionary.
func (m *Message) Dictionary() *dict.Parser {
	if m.dictionary == nil {
		return dict.Default
	}
	return m.dictionary
}

// NewAVP creates and initializes a new AVP and adds it to the Message.
// It is not safe for concurrent calls.
func (m *Message) NewAVP(code interface{}, flags uint8, vendor uint32, data datatype.Type) (*AVP, error) {
	var a *AVP
	switch code.(type) {
	case int:
		a = NewAVP(uint32(code.(int)), flags, vendor, data)
	case uint32:
		a = NewAVP(code.(uint32), flags, vendor, data)
	case string:
		dictAVP, err := m.Dictionary().FindAVPWithVendor(
			m.Header.ApplicationID,
			code.(string),
			vendor,
		)
		if err != nil {
			return nil, err
		}
		a = NewAVP(dictAVP.Code, flags, vendor, data)
	}
	m.AVP = append(m.AVP, a)
	m.Header.MessageLength += uint32(a.Len())
	return a, nil
}

// AddAVP adds the AVP to the Message. It is not safe for concurrent calls.
func (m *Message) AddAVP(a *AVP) {
	m.AVP = append(m.AVP, a)
	m.Header.MessageLength += uint32(a.Len())
}

// InsertAVP inserts the AVP to the Message as the first AVP. It is not
// safe for concurrent calls.
func (m *Message) InsertAVP(a *AVP) {
	m.AVP = append([]*AVP{a}, m.AVP...)
	m.Header.MessageLength += uint32(a.Len())
}

var writerBufferPool sync.Pool

func newWriterBuffer(min int) *bytes.Buffer {
	if min > MessageBufferLength {
		return bytes.NewBuffer(make([]byte, min))
	}
	if v := writerBufferPool.Get(); v != nil {
		return v.(*bytes.Buffer)
	}
	return bytes.NewBuffer(make([]byte, MessageBufferLength))
}

func putWriterBuffer(b *bytes.Buffer) {
	b.Reset()
	if cap(b.Bytes()) == MessageBufferLength {
		writerBufferPool.Put(b)
	}
}

// WriteTo serializes the Message and writes into the writer.
func (m *Message) WriteTo(writer io.Writer) (int64, error) {
	n, err := m.WriteToStream(writer, m.stream)
	return int64(n), err
}

// WriteToStream serializes the Message and writes into the writer with given retries if needed
func (m *Message) WriteToWithRetry(writer io.Writer, retries uint) (int64, error) {
	n, err := m.WriteToStreamWithRetry(writer, m.stream, retries)
	return int64(n), err
}

// WriteToStream serializes the Message and writes into the writer
// If writer implements MultistreamWriter, writes the message into specified stream
func (m *Message) WriteToStream(writer io.Writer, stream uint) (n int, err error) {
	return m.WriteToStreamWithRetry(writer, stream, 0)
}

// WriteToStreamWithRetry serializes the Message and writes into the writer with specified number of retries
// if needed
// If writer implements MultistreamWriter, writes the message into specified stream
func (m *Message) WriteToStreamWithRetry(writer io.Writer, stream, retries uint) (n int, err error) {
	l := m.Len()
	buf := newWriterBuffer(l)
	defer putWriterBuffer(buf)
	b := buf.Bytes()[0:l]
	if err := m.SerializeTo(b); err != nil {
		return 0, err
	}
	switch w := writer.(type) {
	case MultistreamWriter:
		return writeStreamRetry(w, b, stream, retries)
	default:
		return writeRetry(writer, b, retries)
	}
}

func writeRetry(w io.Writer, b []byte, retries uint) (n int, err error) {
	var wn int
	for {
		wn, err = w.Write(b)
		n += wn
		if err == nil || retries == 0 {
			return
		}
		if nerr, isNetErr := err.(net.Error); !(isNetErr && nerr.Temporary()) {
			return
		}
		if wn > 0 {
			b = b[wn:]
		}
		retries--
	}
}

func writeStreamRetry(w MultistreamWriter, b []byte, stream, retries uint) (n int, err error) {
	var wn int
	for {
		wn, err = w.WriteStream(b, stream)
		n += wn
		if err == nil || retries == 0 {
			return
		}
		if nerr, isNetErr := err.(net.Error); !(isNetErr && nerr.Temporary()) {
			return
		}
		if wn > 0 {
			b = b[wn:]
		}
		retries--
	}
}

// Serialize returns the serialized bytes of the Message.
func (m *Message) Serialize() ([]byte, error) {
	b := make([]byte, m.Len())
	if err := m.SerializeTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// SerializeTo writes the serialized bytes of the Message into b.
func (m *Message) SerializeTo(b []byte) (err error) {
	m.Header.SerializeTo(b[0:HeaderLength])
	offset := HeaderLength
	for _, avp := range m.AVP {
		if err = avp.SerializeTo(b[offset:]); err != nil {
			return err
		}
		offset += avp.Len()
	}
	return nil
}

// Len returns the length of the Message in bytes.
func (m *Message) Len() int {
	l := HeaderLength
	for _, avp := range m.AVP {
		l += avp.Len()
	}
	return l
}

func findFromAVP(avps []*AVP, code uint32, findMultiple bool) ([]*AVP, error) {
	var avpResult []*AVP
	for _, a := range avps {

		if a.Code == code {
			avpResult = append(avpResult, a)
			if !findMultiple {
				return avpResult, nil
			}
		}

		if a.Data.Type() == GroupedAVPType {
			groupedAVP := a.Data
			result, err := findFromAVP(groupedAVP.(*GroupedAVP).AVP, code, findMultiple)
			if err == nil {
				avpResult = append(avpResult, result...)
				if !findMultiple {
					return avpResult, nil
				}
			}
		}
	}

	if len(avpResult) == 0 {
		return nil, errors.New("AVP not found")
	}

	return avpResult, nil
}

// Strict path search (eg: in case of groups)
// Can be also used to search AVPs as findFromAVP
func avpsWithPath(avps []*AVP, path []uint32) []*AVP {
	if len(path) == 0 {
		return avps
	}
	var avsOnPath []*AVP
	for _, avp := range avps {
		if avp.Code != path[0] {
			continue
		}
		if len(path) == 1 { // Reached end
			avsOnPath = append(avsOnPath, avp)
			continue
		}
		if avp.Data.Type() != GroupedAVPType {
			continue
		}
		avpsOnSubpath := avpsWithPath(avp.Data.(*GroupedAVP).AVP, path[1:])
		if len(avpsOnSubpath) != 0 {
			avsOnPath = append(avsOnPath, avpsOnSubpath...)
		}
	}
	return avsOnPath
}

// FindAVPs searches the Message for all avps that match the search criteria.
// The code can be either the AVP code (int, uint32) or name (string).
//
// Example:
//
//	avps, err := m.FindAVPs(264)
//	avps, err := m.FindAVPs(avp.OriginHost)
//	avps, err := m.FindAVPs("Origin-Host")
//
func (m *Message) FindAVPs(code interface{}, vendorID uint32) ([]*AVP, error) {
	dictAVP, err := m.Dictionary().FindAVPWithVendor(m.Header.ApplicationID, code, vendorID)

	if err != nil {
		return nil, err
	}

	return findFromAVP(m.AVP, dictAVP.Code, true)
}

// FindAVP searches the Message for a specific AVP.
// The code can be either the AVP code (int, uint32) or name (string).
//
// Example:
//
//	avp, err := m.FindAVP(264)
//	avp, err := m.FindAVP(avp.OriginHost)
//	avp, err := m.FindAVP("Origin-Host")
//
func (m *Message) FindAVP(code interface{}, vendorID uint32) (*AVP, error) {
	dictAVP, err := m.Dictionary().FindAVPWithVendor(m.Header.ApplicationID, code, vendorID)

	if err != nil {
		return nil, err
	}

	result, err := findFromAVP(m.AVP, dictAVP.Code, false)

	if err == nil {
		return result[0], err
	}
	return nil, err
}

// FindAVPsWithPath searches the Message for AVPs on specific path.
// Used for example on group hierarchies.
// The path elements can be either AVP code (int, uint32), name (string) or combination of them.
//
// Example:
//
//	avp, err := m.FindAVPsWithPath([]interface{}{264})
//	avp, err := m.FindAVPsWithPath([]interface{}{avp.OriginHost})
//	avp, err := m.FindAVPsWithPath([]interface{}{"Origin-Host"})
//
func (m *Message) FindAVPsWithPath(path []interface{}, vendorID uint32) ([]*AVP, error) {
	pathCodes := make([]uint32, len(path))
	for i, pathCode := range path {
		dictAVP, err := m.Dictionary().FindAVPWithVendor(m.Header.ApplicationID, pathCode, vendorID)
		if err != nil {
			return nil, err
		}
		pathCodes[i] = dictAVP.Code
	}
	return avpsWithPath(m.AVP, pathCodes), nil
}

// Answer creates an answer for the current Message
// with optinal ResultCode AVP
func (m *Message) Answer(resultCode uint32) *Message {
	nm := NewMessage(
		m.Header.CommandCode,
		m.Header.CommandFlags&^RequestFlag, // Reset the Request bit.
		m.Header.ApplicationID,
		m.Header.HopByHopID,
		m.Header.EndToEndID,
		m.Dictionary(),
	)
	if resultCode != 0 {
		nm.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(resultCode))
	}
	nm.stream = m.stream
	return nm
}

func (m *Message) String() string {
	var b bytes.Buffer
	var typ string
	if m.Header.CommandFlags&RequestFlag == RequestFlag {
		typ = "Request"
	} else {
		typ = "Answer"
	}
	if dictCMD, err := m.Dictionary().FindCommand(
		m.Header.ApplicationID,
		m.Header.CommandCode,
	); err != nil {
		fmt.Fprintf(&b, "Unknown-%s\n%s\n", typ, m.Header)
	} else {
		fmt.Fprintf(&b, "%s-%s (%s%c)\n%s\n",
			dictCMD.Name,
			typ,
			dictCMD.Short,
			typ[0],
			m.Header,
		)
	}
	for _, a := range m.AVP {
		if dictAVP, err := m.Dictionary().FindAVPWithVendor(
			m.Header.ApplicationID,
			a.Code,
			a.VendorID,
		); err != nil {
			fmt.Fprintf(&b, "\tUnknown %s (%s)\n", a, err)
		} else if a.Data.Type() == GroupedAVPType {
			fmt.Fprintf(&b, "\t%s %s\n", dictAVP.Name, printGrouped("\t", m, a, 1))
		} else {
			fmt.Fprintf(&b, "\t%s %s\n", dictAVP.Name, a)
		}
	}
	return b.String()
}

func printGrouped(prefix string, m *Message, a *AVP, indent int) string {
	var b bytes.Buffer
	fmt.Fprintf(&b, "{Code:%d,Flags:0x%x,Length:%d,VendorId:%d,Value:Grouped{\n",
		a.Code,
		a.Flags,
		a.Len(),
		a.VendorID,
	)
	for _, ga := range a.Data.(*GroupedAVP).AVP {
		if dictAVP, err := m.Dictionary().FindAVPWithVendor(
			m.Header.ApplicationID,
			ga.Code,
			ga.VendorID,
		); err != nil {
			if dictAVP != nil {
				fmt.Fprintf(&b, "%s\t%s %s (%s),\n", prefix, dictAVP.Name, ga, err)
			} else {
				fmt.Fprintf(&b, "%s\tUnknown %s (%s),\n", prefix, ga, err)
			}
		} else {
			if ga.Data.Type() == GroupedAVPType {
				indent++
				tabs := indentTabs(indent)
				fmt.Fprintf(&b, "%s%s %s\n", tabs, dictAVP.Name, printGrouped(tabs, m, ga, indent))
			} else {
				fmt.Fprintf(&b, "%s\t%s %s,\n", prefix, dictAVP.Name, ga)
			}
		}
	}
	fmt.Fprintf(&b, "%s}}", prefix)
	return b.String()
}

func indentTabs(n int) string {
	var s string
	for i := 0; i < n; i++ {
		s += "\t"
	}
	return s
}
