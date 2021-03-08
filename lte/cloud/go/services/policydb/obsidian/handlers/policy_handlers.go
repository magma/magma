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
	"net/http"
	"sort"
	"strings"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	baseNameParam   = "base_name"
	ruleIDParam     = "rule_id"
	qosProfileParam = "profile_id"
)

// Base names

func ListBaseNames(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	view := c.QueryParam("view")
	if strings.ToLower(view) == "full" {
		baseNames, _, err := configurator.LoadAllEntitiesOfType(
			networkID, lte.BaseNameEntityType,
			configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadAssocsToThis: true},
			serdes.Entity,
		)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		ret := map[string]*models.BaseNameRecord{}
		for _, bnEnt := range baseNames {
			ret[bnEnt.Key] = (&models.BaseNameRecord{}).FromEntity(bnEnt)
		}
		return c.JSON(http.StatusOK, ret)
	} else {
		names, err := configurator.ListEntityKeys(networkID, lte.BaseNameEntityType)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		sort.Strings(names)
		return c.JSON(http.StatusOK, names)
	}
}

func CreateBaseName(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	bnr := &models.BaseNameRecord{}
	if err := c.Bind(bnr); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	bnrEnt := bnr.ToEntity()

	// Verify that subscribers and policies exist
	parents := bnr.GetParentAssocs()
	doAssignedAssocsExist, _ := configurator.DoEntitiesExist(networkID, parents)
	if !doAssignedAssocsExist {
		return obsidian.HttpError(errors.New("failed to create base name: one or more subscribers or policies do not exist"), http.StatusInternalServerError)
	}

	// In one transaction
	// 	- create base name, with child assocs
	//	- update parent assocs: subscriber

	var writes []configurator.EntityWriteOperation
	writes = append(writes, bnr.ToEntity())
	for _, tk := range parents {
		if tk.Type == lte.SubscriberEntityType {
			w := configurator.EntityUpdateCriteria{
				Type:              lte.SubscriberEntityType,
				Key:               tk.Key,
				AssociationsToAdd: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: bnrEnt.Key}},
			}
			writes = append(writes, w)
		}
	}
	if err := configurator.WriteEntities(networkID, writes, serdes.Entity); err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to create base name"), http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, string(bnr.Name))
}

func GetBaseName(c echo.Context) error {
	networkID, baseName, nerr := getNetworkAndParam(c, baseNameParam)
	if nerr != nil {
		return nerr
	}

	ret, err := configurator.LoadEntity(
		networkID, lte.BaseNameEntityType, baseName,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadAssocsToThis: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, (&models.BaseNameRecord{}).FromEntity(ret))
}

func UpdateBaseName(c echo.Context) error {
	networkID, baseName, nerr := getNetworkAndParam(c, baseNameParam)
	if nerr != nil {
		return nerr
	}

	bnr := &models.BaseNameRecord{}
	if err := c.Bind(bnr); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if string(bnr.Name) != baseName {
		return obsidian.HttpError(errors.New("base name in body does not match URL param"), http.StatusBadRequest)
	}

	// 404 if the entity doesn't exist
	oldEnt, err := configurator.LoadEntity(
		networkID, lte.BaseNameEntityType, baseName,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadAssocsToThis: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(errors.Wrap(err, "failed to check if base name exists"), http.StatusInternalServerError)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// Verify that associated subscribers and policies exist
	parents := bnr.GetParentAssocs()
	assocsExist, _ := configurator.DoEntitiesExist(networkID, parents)
	if !assocsExist {
		return obsidian.HttpError(errors.New("failed to update base name: one or more subscribers or policies do not exist"), http.StatusInternalServerError)
	}

	// In one transaction
	// 	- modify base name, with child assocs
	//	- update parent assocs: subscriber

	var writes []configurator.EntityWriteOperation
	writes = append(writes, bnr.ToUpdateCriteria())

	remove, add := oldEnt.ParentAssociations.Difference(bnr.GetParentAssocs())
	for _, tk := range remove.Filter(lte.SubscriberEntityType) {
		w := configurator.EntityUpdateCriteria{
			Type:                 lte.SubscriberEntityType,
			Key:                  tk.Key,
			AssociationsToDelete: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: baseName}},
		}
		writes = append(writes, w)
	}
	for _, tk := range add.Filter(lte.SubscriberEntityType) {
		w := configurator.EntityUpdateCriteria{
			Type:              lte.SubscriberEntityType,
			Key:               tk.Key,
			AssociationsToAdd: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: baseName}},
		}
		writes = append(writes, w)
	}

	if err = configurator.WriteEntities(networkID, writes, serdes.Entity); err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to update base name"), http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteBaseName(c echo.Context) error {
	networkID, baseName, nerr := getNetworkAndParam(c, baseNameParam)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(networkID, lte.BaseNameEntityType, baseName)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

// Rules

func ListRules(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	view := c.QueryParam("view")
	if strings.ToLower(view) == "full" {
		rules, _, err := configurator.LoadAllEntitiesOfType(
			networkID, lte.PolicyRuleEntityType,
			configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true, LoadAssocsToThis: true},
			serdes.Entity,
		)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		ret := map[string]*models.PolicyRule{}
		for _, ruleEnt := range rules {
			ret[ruleEnt.Key] = (&models.PolicyRule{}).FromEntity(ruleEnt)
		}
		return c.JSON(http.StatusOK, ret)
	} else {
		ruleIDs, err := configurator.ListEntityKeys(networkID, lte.PolicyRuleEntityType)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		sort.Strings(ruleIDs)
		return c.JSON(http.StatusOK, ruleIDs)
	}
}

func CreateRule(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	rule := &models.PolicyRule{}
	if err := c.Bind(rule); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := rule.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	updateToNewIPModel(rule.FlowList)

	// Verify that subscribers and policies exist
	var allAssocs storage.TKs
	allAssocs = append(allAssocs, rule.GetParentAssocs()...)
	allAssocs = append(allAssocs, rule.GetAssocs()...)
	assocsExist, _ := configurator.DoEntitiesExist(networkID, allAssocs)
	if !assocsExist {
		return obsidian.HttpError(errors.New("failed to create policy: one or more subscribers or QoS profiles do not exist"), http.StatusInternalServerError)
	}

	// In one transaction, create the policy rule and associate subscribers
	// to it. Succeeds or fails in its entirety.
	// Create entity
	createdEntity := rule.ToEntity()
	var writes []configurator.EntityWriteOperation
	writes = append(writes, createdEntity)

	for _, tk := range rule.GetParentAssocs().Filter(lte.SubscriberEntityType) {
		w := configurator.EntityUpdateCriteria{
			Type:              lte.SubscriberEntityType,
			Key:               tk.Key,
			AssociationsToAdd: []storage.TypeAndKey{{Type: lte.PolicyRuleEntityType, Key: createdEntity.Key}},
		}
		writes = append(writes, w)
	}

	if err := configurator.WriteEntities(networkID, writes, serdes.Entity); err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to create policy"), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func GetRule(c echo.Context) error {
	networkID, ruleID, nerr := getNetworkAndParam(c, ruleIDParam)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		networkID, lte.PolicyRuleEntityType, ruleID,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true, LoadAssocsToThis: true},
		serdes.Entity,
	)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, (&models.PolicyRule{}).FromEntity(ent))
}

func UpdateRule(c echo.Context) error {
	networkID, ruleID, nerr := getNetworkAndParam(c, ruleIDParam)
	if nerr != nil {
		return nerr
	}

	rule := &models.PolicyRule{}
	if err := c.Bind(rule); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := rule.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if ruleID != string(rule.ID) {
		return obsidian.HttpError(errors.New("rule ID in body does not match URL param"), http.StatusBadRequest)
	}

	updateToNewIPModel(rule.FlowList)

	// 404 if the rule doesn't exist
	oldEnt, err := configurator.LoadEntity(
		networkID, lte.PolicyRuleEntityType, ruleID,
		configurator.EntityLoadCriteria{LoadAssocsToThis: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(errors.Wrap(err, "Failed to check if policy exists"), http.StatusInternalServerError)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	// Verify subscribers and policies exist
	var allAssocs storage.TKs
	allAssocs = append(allAssocs, rule.GetParentAssocs()...)
	allAssocs = append(allAssocs, rule.GetAssocs()...)
	assocsExist, _ := configurator.DoEntitiesExist(networkID, allAssocs)
	if !assocsExist {
		return obsidian.HttpError(errors.New("failed to create policy: one or more subscribers or QoS profiles do not exist"), http.StatusInternalServerError)
	}

	// In one transaction
	// 	- modify policy rule
	// 	- update parent assocs: subscriber
	//	- update child assocs: policy_qos_profile

	var writes []configurator.EntityWriteOperation
	writes = append(writes, rule.ToEntityUpdateCriteria())

	remove, add := oldEnt.ParentAssociations.Difference(rule.GetParentAssocs())
	for _, tk := range remove.Filter(lte.SubscriberEntityType) {
		w := configurator.EntityUpdateCriteria{
			Type:                 lte.SubscriberEntityType,
			Key:                  tk.Key,
			AssociationsToDelete: []storage.TypeAndKey{{Type: lte.PolicyRuleEntityType, Key: ruleID}},
		}
		writes = append(writes, w)
	}
	for _, tk := range add.Filter(lte.SubscriberEntityType) {
		w := configurator.EntityUpdateCriteria{
			Type:              lte.SubscriberEntityType,
			Key:               tk.Key,
			AssociationsToAdd: []storage.TypeAndKey{{Type: lte.PolicyRuleEntityType, Key: ruleID}},
		}
		writes = append(writes, w)
	}

	if err = configurator.WriteEntities(networkID, writes, serdes.Entity); err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to update policy rule"), http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteRule(c echo.Context) error {
	networkID, ruleID, nerr := getNetworkAndParam(c, ruleIDParam)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(networkID, lte.PolicyRuleEntityType, ruleID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

// QoS profiles

func getQoSProfiles(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	profiles, _, err := configurator.LoadAllEntitiesOfType(
		networkID, lte.PolicyQoSProfileEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := map[string]*models.PolicyQosProfile{}
	for _, ent := range profiles {
		ret[ent.Key] = ent.Config.(*models.PolicyQosProfile)
	}

	return c.JSON(http.StatusOK, ret)
}

func createQoSProfile(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	profile := &models.PolicyQosProfile{}
	if err := c.Bind(profile); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := profile.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	exists, err := configurator.DoesEntityExist(networkID, lte.PolicyQoSProfileEntityType, profile.ID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if exists {
		return echo.ErrBadRequest
	}

	_, err = configurator.CreateEntity(networkID, profile.ToEntity(), serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func deleteQoSProfile(c echo.Context) error {
	networkID, profileID, nerr := getNetworkAndParam(c, qosProfileParam)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(networkID, lte.PolicyQoSProfileEntityType, profileID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func getNetworkAndParam(c echo.Context, paramName string) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", paramName)
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}

func updateToNewIPModel(flowList []*models.FlowDescription) {
	for _, flow_desc := range flowList {
		if flow_desc.Match.IPV4Src != "" {
			flow_desc.Match.IPSrc = &models.IPAddress{
				Version: models.IPAddressVersionIPV4,
				Address: flow_desc.Match.IPV4Src,
			}
			flow_desc.Match.IPV4Src = ""
		}
		if flow_desc.Match.IPV4Dst != "" {
			flow_desc.Match.IPDst = &models.IPAddress{
				Version: models.IPAddressVersionIPV4,
				Address: flow_desc.Match.IPV4Dst,
			}
			flow_desc.Match.IPV4Dst = ""
		}
	}
}
