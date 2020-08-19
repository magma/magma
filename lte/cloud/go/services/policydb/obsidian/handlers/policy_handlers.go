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
		baseNames, err := configurator.LoadAllEntitiesInNetwork(networkID, lte.BaseNameEntityType, configurator.EntityLoadCriteria{LoadAssocsToThis: true})
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
	bnr := new(models.BaseNameRecord)
	if err := c.Bind(bnr); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// Verify that subscribers and policies exist
	allAssocs := bnr.GetParentAssociations()
	doAssignedAssocsExist, _ := configurator.DoEntitiesExist(networkID, allAssocs)
	if !doAssignedAssocsExist {
		return obsidian.HttpError(errors.New("failed to create base name: one or more subscribers or policies do not exist"), http.StatusInternalServerError)
	}

	// In one transaction, create the base name and associate subscribers
	// and policies to it. Succeeds or fails in its entirety.
	// Create entity
	createdEntity := bnr.ToEntity()
	writes := []configurator.EntityWriteOperation{}
	writes = append(writes, createdEntity)
	// Update entity operations for subscribers and policies to point
	for _, tk := range allAssocs {
		if tk.Type == lte.PolicyRuleEntityType {
			writes = append(writes, configurator.EntityUpdateCriteria{
				Type:              lte.PolicyRuleEntityType,
				Key:               tk.Key,
				AssociationsToAdd: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: createdEntity.Key}},
			})
		} else if tk.Type == lte.SubscriberEntityType {
			writes = append(writes, configurator.EntityUpdateCriteria{
				Type:              lte.SubscriberEntityType,
				Key:               tk.Key,
				AssociationsToAdd: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: createdEntity.Key}},
			})
		}
	}
	if err := configurator.WriteEntities(networkID, writes...); err != nil {
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
		networkID,
		lte.BaseNameEntityType,
		baseName,
		configurator.EntityLoadCriteria{LoadAssocsToThis: true},
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
	prevBaseNameEnt, err := configurator.LoadEntity(
		networkID,
		lte.BaseNameEntityType,
		baseName,
		configurator.EntityLoadCriteria{LoadAssocsToThis: true},
	)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(errors.Wrap(err, "failed to check if base name exists"), http.StatusInternalServerError)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	prevBaseName := (&models.BaseNameRecord{}).FromEntity(prevBaseNameEnt)

	// Verify that associated subscribers and policies exist
	allAssocs := bnr.GetParentAssociations()
	doAssignedAssocsExist, _ := configurator.DoEntitiesExist(networkID, allAssocs)
	if !doAssignedAssocsExist {
		return obsidian.HttpError(errors.New("failed to update base name: one or more subscribers or policies do not exist"), http.StatusInternalServerError)
	}

	// In one transaction, modify the base name, and change associations
	// between subscribers/policies and the base name
	// Succeeds or fails in its entirety.
	writes := []configurator.EntityWriteOperation{}
	prevAssocs := prevBaseName.GetParentAssociations()
	assocsToRemove := getTypeAndKeyDiff(prevAssocs, allAssocs)
	for _, tk := range assocsToRemove {
		if tk.Type == lte.PolicyRuleEntityType {
			writes = append(writes, configurator.EntityUpdateCriteria{
				Type:                 lte.PolicyRuleEntityType,
				Key:                  tk.Key,
				AssociationsToDelete: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: baseName}},
			})
		} else if tk.Type == lte.SubscriberEntityType {
			writes = append(writes, configurator.EntityUpdateCriteria{
				Type:                 lte.SubscriberEntityType,
				Key:                  tk.Key,
				AssociationsToDelete: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: baseName}},
			})
		}
	}
	assocsToAdd := getTypeAndKeyDiff(allAssocs, prevAssocs)
	for _, tk := range assocsToAdd {
		if tk.Type == lte.PolicyRuleEntityType {
			writes = append(writes, configurator.EntityUpdateCriteria{
				Type:              lte.PolicyRuleEntityType,
				Key:               tk.Key,
				AssociationsToAdd: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: baseName}},
			})
		} else if tk.Type == lte.SubscriberEntityType {
			writes = append(writes, configurator.EntityUpdateCriteria{
				Type:              lte.SubscriberEntityType,
				Key:               tk.Key,
				AssociationsToAdd: []storage.TypeAndKey{{Type: lte.BaseNameEntityType, Key: baseName}},
			})
		}
	}
	if err = configurator.WriteEntities(networkID, writes...); err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to update base name"), http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func getTypeAndKeyDiff(a []storage.TypeAndKey, b []storage.TypeAndKey) []storage.TypeAndKey {
	aLessB := map[string]storage.TypeAndKey{}
	for _, tk := range a {
		aLessB[tk.String()] = tk
	}
	for _, tk := range b {
		delete(aLessB, tk.String())
	}
	ret := []storage.TypeAndKey{}
	for _, v := range aLessB {
		ret = append(ret, v)
	}
	return ret
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
		rules, err := configurator.LoadAllEntitiesInNetwork(
			networkID, lte.PolicyRuleEntityType,
			configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
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

	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := rule.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// Verify that subscribers and policies exist
	allAssocs := rule.GetParentAssociations()
	doAssignedAssocsExist, _ := configurator.DoEntitiesExist(networkID, allAssocs)
	if !doAssignedAssocsExist {
		return obsidian.HttpError(errors.New("failed to create policy: one or more subscribers do not exist"), http.StatusInternalServerError)
	}

	// In one transaction, create the policy rule and associate subscribers
	// to it. Succeeds or fails in its entirety.
	// Create entity
	createdEntity := rule.ToEntity()
	writes := []configurator.EntityWriteOperation{}
	writes = append(writes, createdEntity)
	// Update entity operations for subscribers and policies to point
	for _, tk := range allAssocs {
		writes = append(writes, configurator.EntityUpdateCriteria{
			Type:              lte.SubscriberEntityType,
			Key:               tk.Key,
			AssociationsToAdd: []storage.TypeAndKey{{Type: lte.PolicyRuleEntityType, Key: createdEntity.Key}},
		})
	}
	if err := configurator.WriteEntities(networkID, writes...); err != nil {
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
		networkID,
		lte.PolicyRuleEntityType,
		ruleID,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsToThis: true},
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

	rule := new(models.PolicyRule)
	if err := c.Bind(rule); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := rule.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if ruleID != string(rule.ID) {
		return obsidian.HttpError(errors.New("rule ID in body does not match URL param"), http.StatusBadRequest)
	}

	// 404 if the rule doesn't exist
	prevPolicyEnt, err := configurator.LoadEntity(
		networkID,
		lte.PolicyRuleEntityType,
		ruleID,
		configurator.EntityLoadCriteria{LoadAssocsToThis: true},
	)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(errors.Wrap(err, "Failed to check if policy exists"), http.StatusInternalServerError)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	prevPolicy := (&models.PolicyRule{}).FromEntity(prevPolicyEnt)

	// Verify that associated subscribers and policies exist
	allAssocs := rule.GetParentAssociations()
	doAssignedAssocsExist, _ := configurator.DoEntitiesExist(networkID, allAssocs)
	if !doAssignedAssocsExist {
		return obsidian.HttpError(errors.New("failed to update policy rule: one or more subscribers do not exist"), http.StatusInternalServerError)
	}

	// In one transaction, modify the policy rule, and change associations
	// from subscribers.
	// Succeeds or fails in its entirety.
	writes := []configurator.EntityWriteOperation{}
	writes = append(writes, rule.ToEntityUpdateCriteria())
	prevAssocs := prevPolicy.GetParentAssociations()
	assocsToRemove := getTypeAndKeyDiff(prevAssocs, allAssocs)
	for _, tk := range assocsToRemove {
		if tk.Type == lte.SubscriberEntityType {
			writes = append(writes, configurator.EntityUpdateCriteria{
				Type:                 lte.SubscriberEntityType,
				Key:                  tk.Key,
				AssociationsToDelete: []storage.TypeAndKey{{Type: lte.PolicyRuleEntityType, Key: ruleID}},
			})
		}
	}
	assocsToAdd := getTypeAndKeyDiff(allAssocs, prevAssocs)
	for _, tk := range assocsToAdd {
		if tk.Type == lte.SubscriberEntityType {
			writes = append(writes, configurator.EntityUpdateCriteria{
				Type:              lte.SubscriberEntityType,
				Key:               tk.Key,
				AssociationsToAdd: []storage.TypeAndKey{{Type: lte.PolicyRuleEntityType, Key: ruleID}},
			})
		}
	}
	if err = configurator.WriteEntities(networkID, writes...); err != nil {
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

func ListQoSProfiles(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ids, err := configurator.ListEntityKeys(networkID, lte.PolicyQoSProfileEntityType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	sort.Strings(ids)
	return c.JSON(http.StatusOK, ids)
}

func CreateQoSProfile(c echo.Context) error {
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

	exists, err := configurator.DoesEntityExist(networkID, lte.PolicyQoSProfileEntityType, *profile.ID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if exists {
		return echo.ErrNotFound
	}

	_, err = configurator.CreateEntity(networkID, profile.ToEntity())
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, *profile.ID)
}

func DeleteQoSProfile(c echo.Context) error {
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
