/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package servicers implements EAP-AKA GRPC service
package servicers

import (
	"log"
	"sync"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers/aka"
)

// Handler - is an AKA Subtype
type Handler func(srvr *EapAkaSrv, ctx *protos.Context, req eap.Packet) (eap.Packet, error)

var akaHandlers struct {
	rwl sync.RWMutex
	hm  map[aka.Subtype]Handler
}

func AddHandler(st aka.Subtype, h Handler) {
	if h == nil {
		return
	}
	akaHandlers.rwl.Lock()
	if akaHandlers.hm == nil {
		akaHandlers.hm = map[aka.Subtype]Handler{}
	}
	oldh, ok := akaHandlers.hm[st]
	if ok && oldh != nil {
		log.Printf("WARNING: EAP AKA Handler for subtype %d => %+v is already registered, will overwrite with %+v",
			st, oldh, h)
	}
	akaHandlers.hm[st] = h
	akaHandlers.rwl.Unlock()
}

func GetHandler(st aka.Subtype) Handler {
	akaHandlers.rwl.RLock()
	defer akaHandlers.rwl.RUnlock()
	res, ok := akaHandlers.hm[st]
	if ok {
		return res
	}
	return nil
}
