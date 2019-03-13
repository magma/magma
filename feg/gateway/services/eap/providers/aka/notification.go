/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSDstyle license found in the
LICENSE file in the root directory of this source tree.
*/

// package servce implements EAP-AKA GRPC service
package aka

import (
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/protos"
)

func NewAKANotificationReq(identifier uint8, code uint16) eap.Packet {
	return []byte{
		eap.RequestCode,
		identifier,
		0, 12, // EAP Len
		TYPE,
		byte(SubtypeNotification),
		0, 0,
		byte(AT_NOTIFICATION),
		1, // EAP AKA Attr Len
		uint8(code >> 8), uint8(code)}
}

func EapErrorResPacket(id uint8, code uint16, rpcCode codes.Code, f string, a ...interface{}) (eap.Packet, error) {
	return NewAKANotificationReq(id, code), Errorf(rpcCode, f, a...)
}

func EapErrorRes(id uint8, code uint16, rpcCode codes.Code, f string, a ...interface{}) (*protos.Eap, error) {
	return &protos.Eap{Payload: NewAKANotificationReq(id, code)}, Errorf(rpcCode, f, a...)
}

func Errorf(code codes.Code, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	log.Printf("AKA RPC [%s] %s", code, msg)
	return status.Errorf(code, msg)
}

func Error(code codes.Code, err error) error {
	log.Printf("AKA RPC [%s] %s", code, err)
	return status.Error(code, err.Error())
}
