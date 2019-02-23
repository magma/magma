/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package test

import (
	"testing"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/storage"
	"magma/lte/cloud/go/protos"

	"github.com/stretchr/testify/assert"
)

const (
	defaultServerProtocol = "tcp"
	defaultServerAddr     = "127.0.0.1:0"
	defaultServerHost     = "magma.com"
	defaultServerRealm    = "magma.com"
	defaultMaxUlBitRate   = uint64(100000000)
	defaultMaxDlBitRate   = uint64(200000000)
)

var (
	defaultLteAuthAmf = []byte("\x80\x00")
	defaultLteAuthOp  = []byte("\xcd\xc2\x02\xd5\x12> \xf6+mgj\xc7,\xb3\x18")
)

// newTestHomeSubscriberServer creates a HSS with test users so its functionality
// can be tested.
func newTestHomeSubscriberServer(t *testing.T) *servicers.HomeSubscriberServer {
	store := storage.NewMemorySubscriberStore()

	sub := &protos.SubscriberData{
		Sid: &protos.SubscriberID{Id: "sub1"},
		Lte: &protos.LTESubscription{
			State:    protos.LTESubscription_ACTIVE,
			AuthAlgo: protos.LTESubscription_MILENAGE,
			AuthKey:  []byte("\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb"),
			AuthOpc:  []byte("\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]"),
		},
		State: &protos.SubscriberState{
			LteAuthNextSeq: 7350,
		},
	}
	err := store.AddSubscriber(sub)
	assert.NoError(t, err)

	config := &mconfig.HSSConfig{
		Server: &mconfig.DiamServerConfig{
			Protocol:  defaultServerProtocol,
			Address:   defaultServerAddr,
			DestHost:  defaultServerHost,
			DestRealm: defaultServerRealm,
		},
		LteAuthAmf: defaultLteAuthAmf,
		LteAuthOp:  defaultLteAuthOp,
		DefaultSubProfile: &mconfig.HSSConfig_SubscriptionProfile{
			MaxUlBitRate: defaultMaxUlBitRate,
			MaxDlBitRate: defaultMaxDlBitRate,
		},
		SubProfiles: make(map[string]*mconfig.HSSConfig_SubscriptionProfile),
	}
	server, err := servicers.NewHomeSubscriberServer(store, config)
	assert.NoError(t, err)
	return server
}
