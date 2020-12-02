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

package handlers

import (
	"magma/lte/cloud/go/serdes"
	lte_handlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	policydb_models "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
)

const (
	qosProfileRootPath   = lte_handlers.ManageNetworkPath + obsidian.UrlSep + "policy_qos_profiles"
	qosProfileManagePath = qosProfileRootPath + obsidian.UrlSep + ":profile_id"

	policiesRootPath         = handlers.ManageNetworkPath + obsidian.UrlSep + "policies"
	policyRuleRootPath       = policiesRootPath + obsidian.UrlSep + "rules"
	policyRuleManagePath     = policyRuleRootPath + obsidian.UrlSep + ":rule_id"
	policyBaseNameRootPath   = policiesRootPath + obsidian.UrlSep + "base_names"
	policyBaseNameManagePath = policyBaseNameRootPath + obsidian.UrlSep + ":base_name"

	ratingGroupsRootPath   = handlers.ManageNetworkPath + obsidian.UrlSep + "rating_groups"
	ratingGroupsManagePath = ratingGroupsRootPath + obsidian.UrlSep + ":rating_group_id"
)

func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: policyBaseNameRootPath, Methods: obsidian.GET, HandlerFunc: ListBaseNames},
		{Path: policyBaseNameRootPath, Methods: obsidian.POST, HandlerFunc: CreateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.GET, HandlerFunc: GetBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateBaseName},
		{Path: policyBaseNameManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteBaseName},

		{Path: qosProfileRootPath, Methods: obsidian.GET, HandlerFunc: getQoSProfiles},
		{Path: qosProfileRootPath, Methods: obsidian.POST, HandlerFunc: createQoSProfile},
		{Path: qosProfileManagePath, Methods: obsidian.DELETE, HandlerFunc: deleteQoSProfile},

		{Path: policyRuleRootPath, Methods: obsidian.GET, HandlerFunc: ListRules},
		{Path: policyRuleRootPath, Methods: obsidian.POST, HandlerFunc: CreateRule},
		{Path: policyRuleManagePath, Methods: obsidian.GET, HandlerFunc: GetRule},
		{Path: policyRuleManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateRule},
		{Path: policyRuleManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteRule},

		{Path: ratingGroupsRootPath, Methods: obsidian.GET, HandlerFunc: ListRatingGroups},
		{Path: ratingGroupsRootPath, Methods: obsidian.POST, HandlerFunc: CreateRatingGroup},
		{Path: ratingGroupsManagePath, Methods: obsidian.GET, HandlerFunc: GetRatingGroup},
		{Path: ratingGroupsManagePath, Methods: obsidian.PUT, HandlerFunc: UpdateRatingGroup},
		{Path: ratingGroupsManagePath, Methods: obsidian.DELETE, HandlerFunc: DeleteRatingGroup},
	}

	ret = append(ret, handlers.GetPartialEntityHandlers(qosProfileManagePath, "profile_id", &policydb_models.PolicyQosProfile{}, serdes.Entity)...)

	return ret
}
