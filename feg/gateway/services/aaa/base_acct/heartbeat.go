/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package base_acct provides a client API for interacting with the
// base_acct cloud service
package base_acct

import (
	"context"
	"sync"
	"time"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
)

const HeartbeatDuration = time.Minute * 5 // TBD: make the duration configurable

var (
	initApn         sync.Once
	apn             string
	heartbeatTicker *time.Ticker
)

func StartBaseAccountingHeartbeat() {
	heartbeatTicker = time.NewTicker(HeartbeatDuration)
	go func() {
		for range heartbeatTicker.C {
			client, err := getBaseAcctClient()
			if err != nil {
				glog.Errorf("base acct heartbeat error: %v", err)
			} else {
				_, err := client.Update(
					context.Background(),
					&protos.AcctUpdateReq{Session: &protos.AcctSession{ServingApn: apn}})
				if err != nil {
					glog.Warningf("APN '%s' heartbeat failure: %v", apn, err)
				}
			}
		}
	}()
}

func initApnFromSession(session *protos.AcctSession) {
	if session != nil {
		apnStr := session.GetServingApn()
		if len(apnStr) > 0 {
			initApn.Do(func() { apn = apnStr })
		}
	}
}
