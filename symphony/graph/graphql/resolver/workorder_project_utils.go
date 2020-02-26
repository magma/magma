// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
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
	typeIDToInput := make(map[string]*models.PropertyInput, len(input))
	switch entity {
	case models.PropertyEntityWorkOrder:
		template, err := r.ClientFrom(ctx).WorkOrderType.Get(ctx, tmplID)
		if err != nil {
			return nil, fmt.Errorf("can't read work order type: %w", err)
		}
		pTyps, erro = template.QueryPropertyTypes().All(ctx)
	case models.PropertyEntityProject:
		template, err := r.ClientFrom(ctx).ProjectType.Get(ctx, tmplID)
		if err != nil {
			return nil, fmt.Errorf("can't read project type: %w", err)
		}
		pTyps, erro = template.QueryProperties().All(ctx)
	default:
		return nil, fmt.Errorf("can't query property types for %v", entity.String())
	}
	if erro != nil {
		return nil, erro
	}
	for _, pInput := range input {
		typeIDToInput[pInput.PropertyTypeID] = pInput
	}
	for _, propTyp := range pTyps {
		if propTyp.Deleted {
			continue
		}
		if _, ok := typeIDToInput[propTyp.ID]; !ok {
			// propTyp not in inputs
			if !skipMandatoryPropertiesCheck && propTyp.Mandatory {
				return nil, fmt.Errorf("property type %v is mandatory and must be specified", propTyp.Name)
			}
			input = append(input, &models.PropertyInput{
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
	return input, nil
}
