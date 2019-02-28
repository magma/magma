/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package eapauth (EAP Authenticator) provides interface to supported & registered EAP Authenticator Providers
package eapauth

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
)

// Packet represents EAP Packet
type Packet []byte

// NewEapPacket creates an EAP Packet with initialized header and appends provided data
// if additionalCapacity is specified, NewEapPacket reserves extra additionalCapacity bytes capacity in the
// returned packet byte slice
func NewEapPacket(code, identifier uint8, data []byte, additionalCapacity ...uint) Packet {
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
