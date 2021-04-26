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

package policydb

import (
	"magma/feg/gateway/object_store"
	"magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

type storedObjectListener struct {
	streamMap object_store.ObjectMap
	name      string
	protoBuf  proto.Message
}

type BaseNameStreamListener struct {
	storedObjectListener
}

type PolicyDBStreamListener struct {
	storedObjectListener
}

type OmnipresentRulesStreamListener struct {
	storedObjectListener
}

func NewBaseNameStreamListener(streamMap object_store.ObjectMap) *BaseNameStreamListener {
	return &BaseNameStreamListener{storedObjectListener: storedObjectListener{streamMap: streamMap, name: "base_names", protoBuf: &protos.ChargingRuleNameSet{}}}
}

func NewPolicyDBStreamListener(streamMap object_store.ObjectMap) *PolicyDBStreamListener {
	return &PolicyDBStreamListener{storedObjectListener: storedObjectListener{streamMap: streamMap, name: "policydb", protoBuf: &protos.PolicyRule{}}}
}

func NewOmnipresentRulesListener(streamMap object_store.ObjectMap) *OmnipresentRulesStreamListener {
	return &OmnipresentRulesStreamListener{storedObjectListener: storedObjectListener{streamMap: streamMap, name: "network_wide_rules", protoBuf: &protos.AssignedPolicies{}}}
}

// Gateway Streamer Listener Interface Implementation
func (listener *storedObjectListener) GetName() string {
	return listener.name
}

func (listener *storedObjectListener) ReportError(e error) error {
	glog.Errorf("GxClient %s streaming error: %v", listener.name, e)
	return nil // continue listener
}

func (listener *storedObjectListener) Update(ub *orcprotos.DataUpdateBatch) bool {
	if !ub.GetResync() {
		return true
	}

	currMap, err := listener.streamMap.GetAll()
	if err != nil {
		glog.Errorf("Streamer error getting current %s: %v", listener.name, err)
		return true
	}

	for _, u := range ub.GetUpdates() {
		messageSet := proto.Clone(listener.protoBuf)
		if err := proto.Unmarshal(u.GetValue(), messageSet); err != nil {
			glog.Errorf("Streamer Unmarshal Error: %v for %s '%s'", err, listener.name, u.GetKey())
			continue
		}
		if err := listener.streamMap.Set(u.GetKey(), messageSet); err != nil {
			glog.Errorf("Streamer store Error: %v for %s '%s'", err, listener.name, u.GetKey())
		}
		delete(currMap, u.GetKey())
	}
	for key := range currMap {
		// leftovers
		if err := listener.streamMap.Delete(key); err != nil {
			glog.Errorf("Streamer deletion Error: %v for %s '%s'", err, listener.name, key)
		}
	}
	return true
}

func (listener *storedObjectListener) GetExtraArgs() *any.Any {
	return nil
}
