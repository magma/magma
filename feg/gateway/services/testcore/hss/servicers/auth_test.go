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

package servicers_test

import (
	"testing"

	"github.com/emakeev/milenage"
	"github.com/stretchr/testify/assert"

	"magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/servicers/test_utils"
	"magma/lte/cloud/go/protos"
)

var (
	defaultPlmn       = []byte("\x02\xf8\x59")
	defaultLteAuthOp  = []byte("\xcd\xc2\x02\xd5\x12> \xf6+mgj\xc7,\xb3\x18")
	defaultLteAuthAmf = []byte("\x80\x00")
	defaultAuthSqnInd = uint64(0)
)

func TestSeqToSqn(t *testing.T) {
	assert.Equal(t, uint64(0x1FE000), servicers.SeqToSqn(0xFF00, 0))
	assert.Equal(t, uint64(0xFFFFFFFFFA00), servicers.SeqToSqn(0xFFFFFFFFFFD0, 0))
	assert.Equal(t, uint64(0x142), servicers.SeqToSqn(0xA, 2))
	assert.Equal(t, uint64(0xFFFFFFFFF805), servicers.SeqToSqn(0xFFFFFFFFFFC0, 5))
}

func TestSplitSqn(t *testing.T) {
	sqn, ind := servicers.SplitSqn(0x1FE001)
	assert.Equal(t, uint64(0xFF00), sqn)
	assert.Equal(t, uint64(0x1), ind)

	sqn, ind = servicers.SplitSqn(0xFFFFFFFFFA1F)
	assert.Equal(t, uint64(0x7FFFFFFFFD0), sqn)
	assert.Equal(t, uint64(0x1F), ind)
}

func TestGetOrGenerateOpc(t *testing.T) {
	lte := &protos.LTESubscription{AuthOpc: []byte("\xcdc\xcbq\x95J\x9fNH\xa5\x99N7\xa0+\xaf")}
	opc, err := servicers.GetOrGenerateOpc(lte, defaultLteAuthOp)
	assert.NoError(t, err)
	assert.Equal(t, lte.AuthOpc, opc)

	lte = &protos.LTESubscription{AuthKey: []byte("\x46\x5b\x5c\xe8\xb1\x99\xb4\x9f\xaa\x5f\x0a\x2e\xe2\x38\xa6\xbc")}
	opc, err = servicers.GetOrGenerateOpc(lte, defaultLteAuthOp)
	assert.NoError(t, err)
	expectedOpc, err := milenage.GenerateOpc(lte.AuthKey, defaultLteAuthOp)
	assert.NoError(t, err)
	assert.Equal(t, expectedOpc[:], opc)
}

func TestGenerateLteAuthVector_MissingLTE(t *testing.T) {
	mcipher, err := milenage.NewCipher(defaultLteAuthAmf)
	assert.NoError(t, err)

	subscriber := &protos.SubscriberData{State: &protos.SubscriberState{}}
	_, _, err = servicers.GenerateLteAuthVector(mcipher, subscriber, defaultPlmn, defaultLteAuthOp, defaultAuthSqnInd)
	assert.Exactly(t, servicers.NewAuthRejectedError("Subscriber data missing LTE subscription"), err)
}

func TestGenerateLteAuthVector_MissingSubscriberState(t *testing.T) {
	mcipher, err := milenage.NewCipher(defaultLteAuthAmf)
	assert.NoError(t, err)

	subscriber := &protos.SubscriberData{
		Lte: &protos.LTESubscription{
			State:    protos.LTESubscription_ACTIVE,
			AuthAlgo: protos.LTESubscription_MILENAGE,
		},
	}
	_, _, err = servicers.GenerateLteAuthVector(mcipher, subscriber, defaultPlmn, defaultLteAuthOp, defaultAuthSqnInd)
	assert.Exactly(t, servicers.NewAuthRejectedError("Subscriber data missing subscriber state"), err)
}

func TestGenerateLteAuthVector_InactiveLTESubscription(t *testing.T) {
	mcipher, err := milenage.NewCipher(defaultLteAuthAmf)
	assert.NoError(t, err)

	subscriber := &protos.SubscriberData{
		Lte: &protos.LTESubscription{
			State:    protos.LTESubscription_INACTIVE,
			AuthAlgo: protos.LTESubscription_MILENAGE,
		},
		State: &protos.SubscriberState{},
	}
	_, _, err = servicers.GenerateLteAuthVector(mcipher, subscriber, defaultPlmn, defaultLteAuthOp, defaultAuthSqnInd)
	assert.Exactly(t, servicers.NewAuthRejectedError("LTE Service not active"), err)
}

func TestGenerateLteAuthVector_UnknownLTEAuthAlgo(t *testing.T) {
	mcipher, err := milenage.NewCipher(defaultLteAuthAmf)
	assert.NoError(t, err)

	subscriber := &protos.SubscriberData{
		Lte: &protos.LTESubscription{
			State:    protos.LTESubscription_ACTIVE,
			AuthAlgo: 10,
		},
		State: &protos.SubscriberState{},
	}
	_, _, err = servicers.GenerateLteAuthVector(mcipher, subscriber, defaultPlmn, defaultLteAuthOp, defaultAuthSqnInd)
	assert.Exactly(t, servicers.NewAuthRejectedError("Unsupported milenage algorithm: 10"), err)
}

func TestGenerateLteAuthVector_Success(t *testing.T) {
	rand := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f")
	mcipher, err := milenage.NewMockCipher([]byte("\x80\x00"), rand)
	assert.NoError(t, err)

	subscriber := &protos.SubscriberData{
		Sid: &protos.SubscriberID{Id: "sub1"},
		Lte: &protos.LTESubscription{
			State:    protos.LTESubscription_ACTIVE,
			AuthAlgo: protos.LTESubscription_MILENAGE,
			AuthKey:  []byte("\x8b\xafG?/\x8f\xd0\x94\x87\xcc\xcb\xd7\t|hb"),
			AuthOpc:  []byte("\x8e'\xb6\xaf\x0ei.u\x0f2fz;\x14`]"),
		},
		State: &protos.SubscriberState{LteAuthNextSeq: 229},
	}
	vector, lteAuthNextSeq, err := servicers.GenerateLteAuthVector(mcipher, subscriber, defaultPlmn, defaultLteAuthOp, 23)
	assert.NoError(t, err)
	assert.Equal(t, uint64(230), lteAuthNextSeq)

	assert.Equal(t, rand, vector.Rand[:])
	assert.Equal(t, []byte("\x2d\xaf\x87\x3d\x73\xf3\x10\xc6"), vector.Xres[:])
	assert.Equal(t, []byte("o\xbf\xa3\x80\x1fW\x80\x00{\xdeY\x88n\x96\xe4\xfe"), vector.Autn[:])
	assert.Equal(t, []byte("\x87H\xc1\xc0\xa2\x82o\xa4\x05\xb1\xe2~\xa1\x04CJ\xe5V\xc7e\xe8\xf0a\xeb\xdb\x8a\xe2\x86\xc4F\x16\xc2"), vector.Kasme[:])
}

func TestResyncLteAuthSeq(t *testing.T) {
	subscriber := test_utils.GetTestSubscribers()[0]
	lteAuthNextSeq, err := servicers.ResyncLteAuthSeq(subscriber, nil, defaultLteAuthOp)
	assert.NoError(t, err)
	assert.Equal(t, lteAuthNextSeq, subscriber.GetState().GetLteAuthNextSeq())

	lteAuthNextSeq, err = servicers.ResyncLteAuthSeq(subscriber, make([]byte, 30), defaultLteAuthOp)
	assert.NoError(t, err)
	assert.Equal(t, lteAuthNextSeq, subscriber.GetState().GetLteAuthNextSeq())

	resyncInfo := make([]byte, 50)
	resyncInfo[25] = 1
	_, err = servicers.ResyncLteAuthSeq(subscriber, resyncInfo, defaultLteAuthOp)
	assert.Exactly(t, servicers.NewAuthRejectedError("resync info incorrect length. expected 30 bytes, but got 50 bytes"), err)

	resyncInfo = make([]byte, 30)
	resyncInfo[0] = 0xFF
	_, err = servicers.ResyncLteAuthSeq(subscriber, resyncInfo, defaultLteAuthOp)
	assert.Exactly(t, servicers.NewAuthRejectedError("Invalid resync authentication code"), err)

	macS := []byte{132, 178, 239, 23, 199, 61, 138, 176}
	copy(resyncInfo[22:], macS)
	lteAuthNextSeq, err = servicers.ResyncLteAuthSeq(subscriber, resyncInfo, defaultLteAuthOp)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0x4204c05f18b), lteAuthNextSeq)
}

func TestGetNextLteAuthSqnAfterResync(t *testing.T) {
	state := &protos.SubscriberState{LteAuthNextSeq: 1 << 30}
	_, err := servicers.GetNextLteAuthSqnAfterResync(state, servicers.SeqToSqn(1<<30-1<<10, 2))
	assert.Exactly(t, servicers.NewAuthRejectedError("Re-sync delta in range but UE rejected auth: 1023"), err)

	lteAuthNextSeq, err := servicers.GetNextLteAuthSqnAfterResync(state, servicers.SeqToSqn(1<<30-1, 3))
	assert.NoError(t, err)
	assert.Equal(t, uint64(1<<30), lteAuthNextSeq)

	_, err = servicers.GetNextLteAuthSqnAfterResync(nil, 0)
	assert.Exactly(t, servicers.NewAuthDataUnavailableError("subscriber state was nil"), err)
}

func TestValidateLteSubscription(t *testing.T) {
	err := servicers.ValidateLteSubscription(nil)
	assert.EqualError(t, err, "Subscriber data missing LTE subscription")

	lte := &protos.LTESubscription{
		State:    protos.LTESubscription_INACTIVE,
		AuthAlgo: protos.LTESubscription_MILENAGE,
	}
	err = servicers.ValidateLteSubscription(lte)
	assert.EqualError(t, err, "LTE Service not active")

	lte = &protos.LTESubscription{
		State:    protos.LTESubscription_ACTIVE,
		AuthAlgo: 50,
	}
	err = servicers.ValidateLteSubscription(lte)
	assert.EqualError(t, err, "Unsupported milenage algorithm: 50")

	lte = &protos.LTESubscription{
		State:    protos.LTESubscription_ACTIVE,
		AuthAlgo: protos.LTESubscription_MILENAGE,
	}
	err = servicers.ValidateLteSubscription(lte)
	assert.NoError(t, err)
}

func TestIsAllZero(t *testing.T) {
	assert.Equal(t, true, servicers.IsAllZero(nil))
	assert.Equal(t, true, servicers.IsAllZero(make([]byte, 50)))

	bytes := make([]byte, 30)
	bytes[25] = 1
	assert.Equal(t, false, servicers.IsAllZero(bytes))
}
