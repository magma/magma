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

package integration

import (
	"magma/feg/cloud/go/protos"

	"github.com/go-openapi/swag"
)

func getPassAllRuleDefinition(ruleID, monitoringKey string, ratingGroup *uint32, precedence uint32) *protos.RuleDefinition {
	rule := &protos.RuleDefinition{
		RuleName:         ruleID,
		Precedence:       precedence,
		FlowDescriptions: []string{"permit out ip from any to any", "permit in ip from any to any"},
		MonitoringKey:    monitoringKey,
	}
	if ratingGroup != nil {
		rule.RatingGroup = swag.Uint32Value(ratingGroup)
	}
	return rule
}
