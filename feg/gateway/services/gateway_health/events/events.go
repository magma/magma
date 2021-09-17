/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package events

import (
	"encoding/json"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/gateway/eventd"
	"magma/gateway/status"
	orcprotos "magma/orc8r/lib/go/protos"
)

const (
	healthStreamName               = "health"
	GatewayPromotionSucceededEvent = "gateway_promotion_succeeded"
	GatewayPromotionFailedEvent    = "gateway_promotion_failed"
	GatewayDemotionSucceededEvent  = "gateway_demotion_succeeded"
	GatewayDemotionFailedEvent     = "gateway_demotion_failed"
)

// GatewayHealthFailed event
type GatewayHealthFailed struct {
	FailureReason string `json:"failure_reason"`
}

// LogGatewayHealthSuccessEvent logs a successful promotion/demotion event.
func LogGatewayHealthSuccessEvent(eventType string, prevAction protos.HealthResponse_RequestedAction) {
	if !shouldLogEvent(eventType, prevAction) {
		return
	}
	go logEvent(eventType, "{}")
}

// LogGatewayHealthFailedEvent logs a failed promotion/demotion event with its
// associated error.
func LogGatewayHealthFailedEvent(eventType string, failureReason string, prevAction protos.HealthResponse_RequestedAction) {
	if !shouldLogEvent(eventType, prevAction) {
		return
	}
	failedEvent := GatewayHealthFailed{
		FailureReason: failureReason,
	}
	serializedHealthEvent, err := json.Marshal(failedEvent)
	if err != nil {
		glog.Errorf("Could not serialize %s event: %s", eventType, err)
		return
	}
	go logEvent(eventType, string(serializedHealthEvent))
}

func logEvent(eventType string, eventValue string) {
	hwid := status.GetHwId()
	event := &orcprotos.Event{
		StreamName: healthStreamName,
		EventType:  eventType,
		Tag:        hwid,
		Value:      eventValue,
	}
	err := eventd.V(eventd.DefaultVerbosity).Log(event)
	if err != nil {
		glog.Errorf("Sending %s event failed: %s", eventType, err)
	}
}

func shouldLogEvent(eventType string, prevAction protos.HealthResponse_RequestedAction) bool {
	// Avoid polluting event logs by only logging new gateway actions
	// (or repeated failures)
	switch prevAction {
	case protos.HealthResponse_SYSTEM_UP:
		if eventType == GatewayPromotionSucceededEvent {
			return false
		}
		return true
	case protos.HealthResponse_SYSTEM_DOWN:
		if eventType == GatewayDemotionSucceededEvent {
			return false
		}
		return true
	case protos.HealthResponse_NONE:
		return true
	default:
		return true
	}
}
