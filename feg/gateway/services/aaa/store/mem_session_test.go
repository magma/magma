/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package store_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"magma/feg/gateway/services/aaa"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/aaa/store"
)

var (
	sharedSid     = aaa.CreateSessionId()
	sharedSession aaa.Session
)

func TestInMemSessionTable(t *testing.T) {
	const routines = 30

	st := store.NewMemorySessionTable()
	var err error
	sharedSession, err = st.AddSession(&protos.Context{SessionId: sharedSid}, time.Minute*10, nil)
	assert.NoError(t, err)
	assert.NotNil(t, sharedSession)

	c := make(chan struct{})
	i := 0
	for ; i < routines; i++ {
		go runTest(t, st, c)
	}
	t.Logf("Started %d test routines\n", i)
	for i = 0; i < routines; i++ {
		<-c
	}
}

type callbackDone bool

func (done *callbackDone) timeoutCallback(s aaa.Session) error {
	if done != nil {
		*done = true
	}
	return nil
}

func runTest(t *testing.T, st aaa.SessionTable, c chan struct{}) {
	defer func() { c <- struct{}{} }()

	shared, err := st.AddSession(&protos.Context{SessionId: sharedSid}, time.Minute*10, nil)
	assert.Error(t, err)
	assert.Equal(t, shared, sharedSession)

	shared.Lock()

	sid := aaa.CreateSessionId()
	pc := &protos.Context{SessionId: sid}

	shared.Unlock()

	// Test Crete session
	var done callbackDone
	s, err := st.AddSession(pc, time.Millisecond*10, (&done).timeoutCallback)
	assert.NoError(t, err)
	assert.NotNil(t, s)

	shared = st.GetSession(sharedSid)
	assert.Equal(t, shared, sharedSession)
	shared.Lock()
	shared.GetCtx().Identity = time.Now().String()
	shared.Unlock()

	// Test Find session
	s1 := st.GetSession(sid)
	assert.Equal(t, s, s1)
	s1.Lock()

	// Test timeout cleanup
	time.Sleep(time.Millisecond * 300)
	s1.Unlock()

	assert.True(t, bool(done))
	s2 := st.GetSession(sid)
	assert.Nil(t, s2)

	// Test Remove session
	s, err = st.AddSession(pc, time.Minute, nil) // don't expire
	assert.NoError(t, err)
	assert.NotNil(t, s)
	s1 = st.GetSession(sid)
	assert.Equal(t, s, s1)
	s2 = st.RemoveSession(sid)
	assert.Equal(t, s1, s2)

	// Test SetTimeout
	s, err = st.AddSession(pc, time.Minute, nil)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	s1 = st.GetSession(sid)
	assert.Equal(t, s, s1)
	done = false
	success := st.SetTimeout(sid, time.Millisecond*5, (&done).timeoutCallback)
	assert.True(t, success)
	time.Sleep(time.Millisecond * 300)

	assert.True(t, bool(done))
	s2 = st.GetSession(sid)
	assert.Nil(t, s2)

	success = st.SetTimeout(sid, time.Millisecond*10, nil)
	assert.False(t, success)
}
