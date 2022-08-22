/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package store_test

import (
	"math/rand"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"magma/feg/gateway/services/aaa"
	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/aaa/store"
)

var (
	sharedSid     = aaa.CreateSessionId()
	sharedImsi    = "0123456789123456"
	sharedSession aaa.Session
)

func TestInMemSessionTable(t *testing.T) {
	const routines = 30

	st := store.NewMemorySessionTable()
	var err error
	sharedSession, err = st.AddSession(&protos.Context{SessionId: sharedSid, Imsi: sharedImsi}, time.Minute*10, nil)
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

type callbackDone int32

func (done *callbackDone) timeoutCallback(s aaa.Session) error {
	if done != nil {
		atomic.StoreInt32((*int32)(done), 1)
	}
	return nil
}

func runTest(t *testing.T, st aaa.SessionTable, c chan struct{}) {
	defer func() { c <- struct{}{} }()

	shared, err := st.AddSession(
		&protos.Context{SessionId: sharedSid,
			Imsi: strconv.FormatUint(rand.Uint64(), 10)[:15]},
		time.Minute*10, nil)
	assert.Error(t, err)
	assert.Equal(t, sharedSid, st.FindSession(sharedImsi))

	shared.Lock()
	assert.Equal(t, shared, sharedSession)
	assert.Equal(t, sharedImsi, shared.GetCtx().GetImsi())
	shared.Unlock()

	sid := aaa.CreateSessionId()
	imsi := strconv.FormatUint(rand.Uint64(), 10)[:15]
	pc := &protos.Context{SessionId: sid, Imsi: imsi}

	// Test Crete session
	var done callbackDone
	s, err := st.AddSession(pc, time.Millisecond*40, (&done).timeoutCallback)
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
	checkSid := st.FindSession(imsi)
	assert.Equal(t, sid, checkSid)
	s1.Lock()

	// Test timeout cleanup
	time.Sleep(time.Millisecond * 300)
	s1.Unlock()

	assert.NotEqual(t, 0, atomic.LoadInt32((*int32)(&done)))
	s2 := st.GetSession(sid)
	assert.Nil(t, s2)
	checkSid = st.FindSession(imsi)
	assert.Equal(t, "", checkSid)

	// Test Remove session
	s, err = st.AddSession(pc, time.Minute, nil) // don't expire
	assert.NoError(t, err)
	assert.NotNil(t, s)
	s1 = st.GetSession(sid)
	assert.Equal(t, s, s1)
	checkSid = st.FindSession(imsi)
	assert.Equal(t, sid, checkSid)
	s2 = st.RemoveSession(sid)
	assert.Equal(t, s1, s2)
	checkSid = st.FindSession(imsi)
	assert.Equal(t, "", checkSid)

	// Test SetTimeout
	s, err = st.AddSession(pc, time.Minute, nil)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	s1 = st.GetSession(sid)
	assert.Equal(t, s, s1)
	atomic.StoreInt32((*int32)(&done), 0)
	success := st.SetTimeout(sid, time.Millisecond*5, (&done).timeoutCallback)
	assert.True(t, success)
	time.Sleep(time.Millisecond * 300)

	assert.NotEqual(t, 0, atomic.LoadInt32((*int32)(&done)))
	s2 = st.GetSession(sid)
	assert.Nil(t, s2)

	success = st.SetTimeout(sid, time.Millisecond*10, nil)
	assert.False(t, success)
}
