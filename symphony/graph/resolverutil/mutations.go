// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolverutil

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/pkg/errors"
)

func ValidateAndGetPositionIfExists(ctx context.Context, client *ent.Client, parentEquipmentID, positionDefinitionID *string, mustBeEmpty bool) (*ent.EquipmentPosition, error) {
	if parentEquipmentID == nil || positionDefinitionID == nil {
		if parentEquipmentID == nil && positionDefinitionID == nil {
			return nil, nil
		}
		return nil, errors.New("both position definition and parent equipment must not be nil")
	}
	ep, err := client.Equipment.Query().
		Where(equipment.ID(*parentEquipmentID)).
		QueryPositions().
		Where(equipmentposition.HasDefinitionWith(
			equipmentpositiondefinition.ID(*positionDefinitionID),
		)).
		Only(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, errors.Wrapf(err, "querying equipment: definition=%q, parent=%q", *positionDefinitionID, *parentEquipmentID)
	}

	if ep != nil {
		if mustBeEmpty {
			hasAttachment, err := ep.QueryAttachment().Exist(ctx)
			if err != nil {
				return nil, err
			}
			if hasAttachment {
				return nil, errors.Wrapf(err, "position already has attachment, position: %q", ep.ID)
			}
		}
		return ep, nil
	}
	return nil, nil
}

func GetOrCreatePosition(ctx context.Context, client *ent.Client, parentEquipmentID, positionDefinitionID *string, mustBeEmpty bool) (*ent.EquipmentPosition, error) {
	if parentEquipmentID == nil && positionDefinitionID == nil {
		return nil, nil
	}
	ep, err := ValidateAndGetPositionIfExists(ctx, client, parentEquipmentID, positionDefinitionID, mustBeEmpty)
	if err != nil {
		return nil, errors.Wrapf(err, "error validating before creating position")
	}
	if ep != nil {
		return ep, nil
	}
	if ep, err = client.EquipmentPosition.Create().
		SetDefinitionID(*positionDefinitionID).
		SetParentID(*parentEquipmentID).
		Save(ctx); err != nil {
		return nil, errors.Wrap(err, "creating equipment position")
	}
	return ep, nil
}

func GetDifferenceBetweenSlices(left, right []string) ([]string, []string) {
	var (
		added, deleted []string
		seen           = map[string]bool{}
	)
	for _, str := range left {
		seen[str] = false
	}
	for _, str := range right {
		if _, ok := seen[str]; ok {
			seen[str] = true
		} else {
			added = append(added, str)
		}
	}
	for str, val := range seen {
		if !val {
			deleted = append(deleted, str)
		}
	}
	return added, deleted
}
