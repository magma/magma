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

// Package test provides common definitions and function for eap related tests
package test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	eap_client "magma/feg/gateway/services/eap/client"
)

// EapClient is a test EAP Client interface
type EapClient interface {
	Handle(msg *protos.Eap) (*protos.Eap, error)
}

// Auth runs EAP-AKA auth sequence for a given IMSI & sends the result on 'done' chan if not nil
func Auth(t *testing.T, client EapClient, imsi string, iter int, done chan error) {
	var (
		err  error
		peap *protos.Eap
	)
	defer func() {
		if done != nil {
			done <- err
		}
		if err != nil {
			t.Fatal(err)
		}
	}()

	tst, found := Units[imsi]
	if !found {
		err = fmt.Errorf("Missing Test Data for IMSI: %s", imsi)
		return
	}

	for i := 0; i < iter; i++ {
		startTime := time.Now()
		eapCtx := &protos.Context{SessionId: eap.CreateSessionId()}
		peap, err = client.Handle(&protos.Eap{Payload: tst.EapIdentityResp, Ctx: eapCtx})
		if err != nil {
			err = fmt.Errorf("Error Handling Test EAP: %v", err)
			return
		}
		if !reflect.DeepEqual(peap.GetPayload(), tst.ExpectedChallengeReq) {
			err = fmt.Errorf("Unexpected identityResponse EAP\n\tReceived: %s\n\tExpected: %s",
				eap_client.BytesToStr(peap.GetPayload()), eap_client.BytesToStr(tst.ExpectedChallengeReq))
			return
		}
		time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond * 10))))
		eapCtx = peap.GetCtx()
		peap, err = client.Handle(&protos.Eap{Payload: tst.EapChallengeResp, Ctx: eapCtx})
		if err != nil {
			err = fmt.Errorf("Error Handling Test Challenge EAP: %v", err)
			return
		}
		successp := []byte{eap.SuccessCode, eap.Packet(tst.EapChallengeResp).Identifier(), 0, 4}
		if !reflect.DeepEqual(peap.GetPayload(), successp) {
			err = fmt.Errorf(
				"Unexpected Challenge Response EAP for Session: %s in %s\n\tReceived: %.3v\n\tExpected: %.3v",
				peap.GetCtx().GetSessionId(), time.Since(startTime), peap.GetPayload(), successp)
			return
		}
		// Check that we got expected MSISDN with the success EAP
		if peap.GetCtx().Msisdn != tst.MSISDN {
			err = fmt.Errorf("Unexpected MSISDN: %s, expected: %s", eapCtx.Msisdn, tst.MSISDN)
			return
		}
		time.Sleep(time.Duration(rand.Int63n(int64(time.Millisecond * 10))))
	}
}
