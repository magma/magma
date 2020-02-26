// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

func (r mutationResolver) validatedPropertyInputsFromTemplate(
	ctx context.Context,
	input []*models.PropertyInput,
	tmplID string,
	entity models.PropertyEntity,
	skipMandatoryPropertiesCheck bool,
) ([]*models.PropertyInput, error) {
	var pTyps []*ent.PropertyType
	var erro error
	typeIDToInput := make(map[string]*models.PropertyInput)
	switch entity {
	case models.PropertyEntityWorkOrder:
		template, err := r.ClientFrom(ctx).WorkOrderType.Get(ctx, tmplID)
		if err != nil {
			return nil, fmt.Errorf("can't read work order type: %w", err)
		}
		pTyps, erro = template.QueryPropertyTypes().
			Where(propertytype.Deleted(false)).
			All(ctx)
	case models.PropertyEntityProject:
		template, err := r.ClientFrom(ctx).ProjectType.Get(ctx, tmplID)
		if err != nil {
			return nil, fmt.Errorf("can't read project type: %w", err)
		}
		pTyps, erro = template.QueryProperties().
			Where(propertytype.Deleted(false)).
			All(ctx)
	default:
		return nil, fmt.Errorf("can't query property types for %v", entity.String())
	}
	if erro != nil {
		return nil, erro
	}

	var validInput []*models.PropertyInput
	for _, pInput := range input {
		// verify it's in pTyps slice &&  not deleted
		candidate := findPropType(pTyps, pInput.PropertyTypeID)
		if candidate != nil {
			validInput = append(validInput, pInput)
			typeIDToInput[pInput.PropertyTypeID] = pInput
		} else {
			return nil, fmt.Errorf("invalid property type (id=%v), either deleted or belongs to other template", pInput.PropertyTypeID)
		}
	}
	for _, propTyp := range pTyps {
		if _, ok := typeIDToInput[propTyp.ID]; !ok {
			// propTyp not in inputs
			if !skipMandatoryPropertiesCheck && propTyp.Mandatory {
				return nil, fmt.Errorf("property type %v is mandatory and must be specified", propTyp.Name)
			}
			validInput = append(validInput, &models.PropertyInput{
				PropertyTypeID:     propTyp.ID,
				StringValue:        &propTyp.StringVal,
				IntValue:           &propTyp.IntVal,
				BooleanValue:       &propTyp.BoolVal,
				FloatValue:         &propTyp.FloatVal,
				LatitudeValue:      &propTyp.LatitudeVal,
				LongitudeValue:     &propTyp.LongitudeVal,
				RangeFromValue:     &propTyp.RangeFromVal,
				RangeToValue:       &propTyp.RangeToVal,
				IsInstanceProperty: &propTyp.IsInstanceProperty,
				IsEditable:         &propTyp.Editable,
			})
		}
	}
	return validInput, nil
}

func findPropType(allTypes []*ent.PropertyType, id string) *ent.PropertyType {
	for _, typ := range allTypes {
		if typ.ID == id {
			return typ
		}
	}
	return nil
}
