/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package protos includes generated GRPC sources as well as corresponding helper functions
package protos

import (
	"strconv"
	"sync"
)

const sidLocksSize = 64

var subscriberLocks [sidLocksSize]sync.Mutex

func (sid *SubscriberID) hash() uint64 {
	h, _ := strconv.ParseUint(sid.GetId(), 10, 64)
	h += uint64(sid.GetType()) << 46
	return h
}

// Lock - lockable interface implementation
func (subscr *SubscriberData) Lock() {
	subscriberLocks[subscr.GetSid().hash()%sidLocksSize].Lock()
}

// Unlock - lockable interface implementation
func (subscr *SubscriberData) Unlock() {
	subscriberLocks[subscr.GetSid().hash()%sidLocksSize].Unlock()
}
