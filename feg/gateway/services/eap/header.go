/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package eap (EAP Authenticator) provides interface to supported & registered EAP Authenticator Providers
//
//go:generate protoc -I. -I ../aaa/protos --go_out=plugins=grpc,paths=source_relative:. protos/eap_auth.proto
//
package eap

import "magma/feg/gateway/services/aaa/protos"

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

	UndefinedCode = 0
	RequestCode   = 1
	ResponseCode  = 2
	SuccessCode   = 3
	FailureCode   = 4
)

const (
	// EAP Related Consts
	MethodIdentity = uint8(protos.EapType_Identity)
	MethodNak      = uint8(protos.EapType_Legacy_Nak)
	CodeResponse   = uint8(protos.EapCode_Response)
)
