/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package eap (EAP Authenticator) provides interface to supported & registered EAP Authenticator Providers
//
//go:generate protoc --go_out=plugins=grpc,paths=source_relative:. protos/eap_auth.proto
//
package eap

import (
	"fmt"
	"io"
)

const (
	// EAP Message Payload Offsets
	EapMsgCode int = iota
	EapMsgIdentifier
	EapMsgLenHigh
	EapMsgLenLow
	EapMsgMethodType
	EapMsgData
	EapReserved1
	EapReserved2
	EapFirstAttribute
	EapFirstAttributeLen
)

const (
	// EapSubtype - pseudonym for EapMsgData
	EapSubtype   = EapMsgData
	EapHeaderLen = EapMsgMethodType
	// EapMaxLen maximum possible length of EAP Packet
	EapMaxLen uint = 1<<16 - 1

	RequestCode  = 1
	ResponseCode = 2
	SuccessCode  = 3
	FailureCode  = 4
)

// Packet represents EAP Packet
type Packet []byte

// NewPacket creates an EAP Packet with initialized header and appends provided data
// if additionalCapacity is specified, NewPacket reserves extra additionalCapacity bytes capacity in the
// returned packet byte slice
func NewPacket(code, identifier uint8, data []byte, additionalCapacity ...uint) Packet {
	l := len(data) + EapHeaderLen
	packetCap := l
	if len(additionalCapacity) > 0 && l < int(EapMaxLen) {
		ac := additionalCapacity[0]
		packetCap = l + int(ac)
		if packetCap > int(EapMaxLen) {
			packetCap = int(EapMaxLen)
		}
	}
	p := make([]byte, EapHeaderLen, packetCap)
	if l > EapHeaderLen {
		p = append(p, data...)
	}
	p[EapMsgCode], p[EapMsgIdentifier], p[EapMsgLenLow], p[EapMsgLenHigh] = code, identifier, uint8(l), uint8(l>>8)
	return p
}

// NewPreallocatedPacket creates an EAP Packet from/in passed data slice & initializes its header
func NewPreallocatedPacket(identifier uint8, data []byte) (Packet, error) {
	l := len(data)
	if l < EapHeaderLen {
		return nil, fmt.Errorf("Data is too short: %d, must be at least %d bytes", l, EapHeaderLen)
	}
	p := Packet(data)
	p[EapMsgIdentifier], p[EapMsgLenLow], p[EapMsgLenHigh] = identifier, uint8(l), uint8(l>>8)
	return p, nil
}

// Validate verifies that the packet is not nil & it's length is correct
func (p Packet) Validate() error {
	lp := len(p)
	if lp < EapHeaderLen {
		return io.ErrShortBuffer
	}
	if p.Len() > lp {
		return fmt.Errorf("Invalid Packet Length: header => %d, actual => %d", p.Len(), lp)
	}
	return nil
}

// Len returns EAP Packet length derived from its header (vs. len of []byte)
func (p Packet) Len() int {
	return (int(p[EapMsgLenHigh]) << 8) + int(p[EapMsgLenLow])
}

// Type returns EAP Method Type or 0 - reserved if not available
func (p Packet) Identifier() uint8 {
	return p[EapMsgIdentifier]
}

// Type returns EAP Method Type or 0 - reserved if not available
func (p Packet) Type() uint8 {
	if len(p) <= EapMsgMethodType {
		return 0
	}
	return p[EapMsgMethodType]
}
