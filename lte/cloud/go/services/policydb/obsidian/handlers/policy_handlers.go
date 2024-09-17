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
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/obsidian"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/merrors"
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

	reqCtx := c.Request().Context()
	view := c.QueryParam("view")
	if strings.ToLower(view) == "full" {
		baseNames, _, err := configurator.LoadAllEntitiesOfType(
			reqCtx,
			networkID, lte.BaseNameEntityType,
			configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadAssocsToThis: true},
			serdes.Entity,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		ret := map[string]*models.BaseNameRecord{}
		for _, bnEnt := range baseNames {
			ret[bnEnt.Key] = (&models.BaseNameRecord{}).FromEntity(bnEnt)
		}
		return c.JSON(http.StatusOK, ret)
	} else {
		names, err := configurator.ListEntityKeys(reqCtx, networkID, lte.BaseNameEntityType)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	bnrEnt := bnr.ToEntity()
	reqCtx := c.Request().Context()

	// Verify that subscribers and policies exist
	parents := bnr.GetParentAssocs()
	doAssignedAssocsExist, _ := configurator.DoEntitiesExist(reqCtx, networkID, parents)
	if !doAssignedAssocsExist {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to create base name: one or more subscribers or policies do not exist"))
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
				AssociationsToAdd: storage.TKs{{Type: lte.BaseNameEntityType, Key: bnrEnt.Key}},
			}
			writes = append(writes, w)
		}
	}
	if err := configurator.WriteEntities(reqCtx, networkID, writes, serdes.Entity); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create base name: %v", err))
	}

	return c.JSON(http.StatusCreated, string(bnr.Name))
}

func GetBaseName(c echo.Context) error {
	networkID, baseName, nerr := getNetworkAndParam(c, baseNameParam)
	if nerr != nil {
		return nerr
	}

	ret, err := configurator.LoadEntity(
		c.Request().Context(),
		networkID, lte.BaseNameEntityType, baseName,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadAssocsToThis: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, (&models.BaseNameRecord{}).FromEntity(ret))
}

func UpdateBaseName(c echo.Context) error {
	networkID, baseName, nerr := getNetworkAndParam(c, baseNameParam)
	if nerr != nil {
		return nerr
	}
	reqCtx := c.Request().Context()

	bnr := &models.BaseNameRecord{}
	if err := c.Bind(bnr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if string(bnr.Name) != baseName {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("base name in body does not match URL param"))
	}

	// 404 if the entity doesn't exist
	oldEnt, err := configurator.LoadEntity(
		reqCtx,
		networkID, lte.BaseNameEntityType, baseName,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadAssocsToThis: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to check if base name exists: %v", err))
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Verify that associated subscribers and policies exist
	parents := bnr.GetParentAssocs()
	assocsExist, _ := configurator.DoEntitiesExist(reqCtx, networkID, parents)
	if !assocsExist {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to update base name: one or more subscribers or policies do not exist"))
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
			AssociationsToDelete: storage.TKs{{Type: lte.BaseNameEntityType, Key: baseName}},
		}
		writes = append(writes, w)
	}
	for _, tk := range add.Filter(lte.SubscriberEntityType) {
		w := configurator.EntityUpdateCriteria{
			Type:              lte.SubscriberEntityType,
			Key:               tk.Key,
			AssociationsToAdd: storage.TKs{{Type: lte.BaseNameEntityType, Key: baseName}},
		}
		writes = append(writes, w)
	}

	if err = configurator.WriteEntities(reqCtx, networkID, writes, serdes.Entity); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to update base name: %v", err))
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteBaseName(c echo.Context) error {
	networkID, baseName, nerr := getNetworkAndParam(c, baseNameParam)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(c.Request().Context(), networkID, lte.BaseNameEntityType, baseName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	reqCtx := c.Request().Context()
	if strings.ToLower(view) == "full" {
		rules, _, err := configurator.LoadAllEntitiesOfType(
			reqCtx,
			networkID, lte.PolicyRuleEntityType,
			configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true, LoadAssocsToThis: true},
			serdes.Entity,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		ret := map[string]*models.PolicyRule{}
		for _, ruleEnt := range rules {
			ret[ruleEnt.Key] = (&models.PolicyRule{}).FromEntity(ruleEnt)
		}
		return c.JSON(http.StatusOK, ret)
	} else {
		ruleIDs, err := configurator.ListEntityKeys(reqCtx, networkID, lte.PolicyRuleEntityType)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
	reqCtx := c.Request().Context()

	rule := &models.PolicyRule{}
	if err := c.Bind(rule); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := rule.ValidateModel(reqCtx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	updateToNewIPModel(rule.FlowList)

	// Verify that subscribers and policies exist
	var allAssocs storage.TKs
	allAssocs = append(allAssocs, rule.GetParentAssocs()...)
	allAssocs = append(allAssocs, rule.GetAssocs()...)
	assocsExist, _ := configurator.DoEntitiesExist(reqCtx, networkID, allAssocs)
	if !assocsExist {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to create policy: one or more subscribers or QoS profiles do not exist"))
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
			AssociationsToAdd: storage.TKs{{Type: lte.PolicyRuleEntityType, Key: createdEntity.Key}},
		}
		writes = append(writes, w)
	}

	if err := configurator.WriteEntities(reqCtx, networkID, writes, serdes.Entity); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to create policy: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}

func GetRule(c echo.Context) error {
	networkID, ruleID, nerr := getNetworkAndParam(c, ruleIDParam)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		c.Request().Context(),
		networkID, lte.PolicyRuleEntityType, ruleID,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true, LoadAssocsToThis: true},
		serdes.Entity,
	)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, (&models.PolicyRule{}).FromEntity(ent))
}

func UpdateRule(c echo.Context) error {
	networkID, ruleID, nerr := getNetworkAndParam(c, ruleIDParam)
	if nerr != nil {
		return nerr
	}
	reqCtx := c.Request().Context()

	rule := &models.PolicyRule{}
	if err := c.Bind(rule); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := rule.ValidateModel(reqCtx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if ruleID != string(*rule.ID) {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("rule ID in body does not match URL param"))
	}

	updateToNewIPModel(rule.FlowList)

	// 404 if the rule doesn't exist
	oldEnt, err := configurator.LoadEntity(
		reqCtx,
		networkID, lte.PolicyRuleEntityType, ruleID,
		configurator.EntityLoadCriteria{LoadAssocsToThis: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to check if policy exists: %v", err))
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// Verify subscribers and policies exist
	var allAssocs storage.TKs
	allAssocs = append(allAssocs, rule.GetParentAssocs()...)
	allAssocs = append(allAssocs, rule.GetAssocs()...)
	assocsExist, _ := configurator.DoEntitiesExist(reqCtx, networkID, allAssocs)
	if !assocsExist {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("failed to create policy: one or more subscribers or QoS profiles do not exist"))
	}

	// In one transaction
	// 	- modify policy rule
	// 	- update parent assocs: subscriber
	//	- update child assocs: policy_qos_profile

	var writes []configurator.EntityWriteOperation

	removedAssocs, addedAssocs := oldEnt.Associations.Difference(rule.GetAssocs())
	writes = append(writes, rule.ToEntityUpdateCriteria(
		addedAssocs.Filter(lte.PolicyQoSProfileEntityType),
		removedAssocs.Filter(lte.PolicyQoSProfileEntityType),
	))

	remove, add := oldEnt.ParentAssociations.Difference(rule.GetParentAssocs())
	for _, tk := range remove.Filter(lte.SubscriberEntityType) {
		w := configurator.EntityUpdateCriteria{
			Type:                 lte.SubscriberEntityType,
			Key:                  tk.Key,
			AssociationsToDelete: storage.TKs{{Type: lte.PolicyRuleEntityType, Key: ruleID}},
		}
		writes = append(writes, w)
	}
	for _, tk := range add.Filter(lte.SubscriberEntityType) {
		w := configurator.EntityUpdateCriteria{
			Type:              lte.SubscriberEntityType,
			Key:               tk.Key,
			AssociationsToAdd: storage.TKs{{Type: lte.PolicyRuleEntityType, Key: ruleID}},
		}
		writes = append(writes, w)
	}

	if err = configurator.WriteEntities(reqCtx, networkID, writes, serdes.Entity); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("failed to update policy rule: %v", err))
	}

	return c.NoContent(http.StatusNoContent)
}

func DeleteRule(c echo.Context) error {
	networkID, ruleID, nerr := getNetworkAndParam(c, ruleIDParam)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(c.Request().Context(), networkID, lte.PolicyRuleEntityType, ruleID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
		c.Request().Context(),
		networkID, lte.PolicyQoSProfileEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
	reqCtx := c.Request().Context()

	profile := &models.PolicyQosProfile{}
	if err := c.Bind(profile); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := profile.ValidateModel(reqCtx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	exists, err := configurator.DoesEntityExist(reqCtx, networkID, lte.PolicyQoSProfileEntityType, profile.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	if exists {
		return echo.ErrBadRequest
	}

	_, err = configurator.CreateEntity(reqCtx, networkID, profile.ToEntity(), serdes.Entity)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusCreated)
}

func deleteQoSProfile(c echo.Context) error {
	networkID, profileID, nerr := getNetworkAndParam(c, qosProfileParam)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(c.Request().Context(), networkID, lte.PolicyQoSProfileEntityType, profileID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
