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

/*
policyd assignments servicer provides the gRPC interface for the REST and
services to interact with assignments from policy rules and subscribers.
*/
package servicers

import (
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	orcprotos "magma/orc8r/lib/go/protos"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PolicyAssignmentServer struct{}

func NewPolicyAssignmentServer() *PolicyAssignmentServer {
	return &PolicyAssignmentServer{}
}

func (srv *PolicyAssignmentServer) EnableStaticRules(ctx context.Context, req *protos.EnableStaticRuleRequest) (*orcprotos.Void, error) {
	networkID, err := getNetworkID(ctx)
	if err != nil {
		return nil, err
	}
	if !doesSubscriberAndRulesExist(networkID, req.Imsi, req.RuleIds, req.BaseNames) {
		return nil, status.Errorf(codes.InvalidArgument, "Either a subscriber or one more rules/basenames are not found")
	}
	var updates []configurator.EntityUpdateCriteria
	for _, ruleID := range req.RuleIds {
		updates = append(updates, getRuleUpdateForEnable(ruleID, req.Imsi))
	}
	for _, baseName := range req.BaseNames {
		updates = append(updates, getBaseNameUpdateForEnable(baseName, req.Imsi))
	}
	_, err = configurator.UpdateEntities(networkID, updates, serdes.Entity)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Failed to enable")
	}
	return &orcprotos.Void{}, nil
}

func (srv *PolicyAssignmentServer) DisableStaticRules(ctx context.Context, req *protos.DisableStaticRuleRequest) (*orcprotos.Void, error) {
	networkID, err := getNetworkID(ctx)
	if err != nil {
		return nil, err
	}
	if !doesSubscriberAndRulesExist(networkID, req.Imsi, req.RuleIds, req.BaseNames) {
		return nil, status.Errorf(codes.InvalidArgument, "Either a subscriber or one more rules/basenames are not found")
	}
	var updates []configurator.EntityUpdateCriteria
	for _, ruleID := range req.RuleIds {
		updates = append(updates, getRuleUpdateForDisableRule(ruleID, req.Imsi))
	}
	for _, baseName := range req.BaseNames {
		updates = append(updates, getBaseNameUpdateForDisable(baseName, req.Imsi))
	}
	_, err = configurator.UpdateEntities(networkID, updates, serdes.Entity)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Failed to disable")
	}
	return &orcprotos.Void{}, nil
}

func doesSubscriberAndRulesExist(networkID string, subscriberID string, ruleIDs []string, baseNames []string) bool {
	ids := []storage.TypeAndKey{{Type: lte.SubscriberEntityType, Key: subscriberID}}
	for _, ruleID := range ruleIDs {
		ids = append(ids, storage.TypeAndKey{Type: lte.PolicyRuleEntityType, Key: ruleID})
	}
	for _, baseName := range baseNames {
		ids = append(ids, storage.TypeAndKey{Type: lte.BaseNameEntityType, Key: baseName})
	}
	exists, err := configurator.DoEntitiesExist(networkID, ids)
	if err != nil {
		return false
	}
	return exists
}

func getRuleUpdateForEnable(ruleID string, subscriberID string) configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type:              lte.SubscriberEntityType,
		Key:               subscriberID,
		AssociationsToAdd: []storage.TypeAndKey{{Type: lte.PolicyRuleEntityType, Key: ruleID}},
	}
	return ret
}

func getRuleUpdateForDisableRule(ruleID string, subscriberID string) configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type:                 lte.SubscriberEntityType,
		Key:                  subscriberID,
		AssociationsToDelete: []storage.TypeAndKey{{Type: lte.PolicyRuleEntityType, Key: ruleID}},
	}
	return ret
}

func getBaseNameUpdateForEnable(baseName string, subscriberID string) configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type:              lte.SubscriberEntityType,
		Key:               subscriberID,
		AssociationsToAdd: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: baseName}},
	}
	return ret
}

func getBaseNameUpdateForDisable(baseName string, subscriberID string) configurator.EntityUpdateCriteria {
	ret := configurator.EntityUpdateCriteria{
		Type:                 lte.SubscriberEntityType,
		Key:                  subscriberID,
		AssociationsToDelete: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: baseName}},
	}
	return ret
}

func getNetworkID(ctx context.Context) (string, error) {
	id, err := orcprotos.GetGatewayIdentity(ctx)
	if err != nil {
		return "", err
	}
	return id.GetNetworkId(), nil
}
