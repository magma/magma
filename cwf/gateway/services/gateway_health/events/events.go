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
	"fmt"

	"magma/gateway/eventd"
	"magma/gateway/status"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
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
func LogGatewayHealthSuccessEvent(eventType string) {
	go logEvent(eventType, "{}")
}

// LogGatewayHealthFailedEvent logs a failed promotion/demotion event with its
// associated error.
func LogGatewayHealthFailedEvent(eventType string, failureReason string) {
	failedEvent := GatewayHealthFailed{
		FailureReason: failureReason,
	}
	serializedHealthEvent, err := json.Marshal(failedEvent)
	if err != nil {
		glog.Errorf("Could not serialize %s event: %s", eventType, err)
		return
	}
	go logEvent(eventType, fmt.Sprintf("%s", serializedHealthEvent))
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
