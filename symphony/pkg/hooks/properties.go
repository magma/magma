// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hooks

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/hook"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
	"github.com/facebookincubator/symphony/pkg/ent/workordertemplate"
	"github.com/facebookincubator/symphony/pkg/viewer"
)

const (
	NodeTypeLocation  = "location"
	NodeTypeEquipment = "equipment"
	NodeTypeService   = "service"
	NodeTypeWorkOrder = "work_order"
	NodeTypeUser      = "user"
)

// errWorkOrderTemplateNotFound error is returned when work order template not found.
var errWorkOrderTemplateNotFound = errors.New("work order template not found")

func isEmptyNodeProp(ctx context.Context, property *ent.Property, nodeType string) (bool, error) {
	var (
		exists bool
		err    error
	)
	if property == nil {
		return true, nil
	}
	switch nodeType {
	case NodeTypeLocation:
		exists, err = property.QueryLocationValue().Exist(ctx)
	case NodeTypeEquipment:
		exists, err = property.QueryEquipmentValue().Exist(ctx)
	case NodeTypeService:
		exists, err = property.QueryServiceValue().Exist(ctx)
	case NodeTypeWorkOrder:
		exists, err = property.QueryWorkOrderValue().Exist(ctx)
	case NodeTypeUser:
		exists, err = property.QueryUserValue().Exist(ctx)
	default:
		return false, fmt.Errorf("unknown node type: %s", nodeType)
	}
	if err != nil {
		if !ent.IsNotFound(err) {
			return true, fmt.Errorf("failed to validate node exists: %w", err)
		}
		return true, nil
	}
	return !exists, nil
}

func isEmptyPrimitiveProp(typ propertytype.Type, strVal *string, intVal *int, boolVal *bool,
	floatVal, lat, long, rangeTo, rangeFrom *float64) (bool, error) {
	switch typ {
	case propertytype.TypeDate,
		propertytype.TypeEmail,
		propertytype.TypeString,
		propertytype.TypeDatetimeLocal,
		propertytype.TypeEnum:
		return pointer.GetString(strVal) == "", nil
	case propertytype.TypeInt:
		return intVal == nil, nil
	case propertytype.TypeGpsLocation:
		if lat == nil || long == nil {
			return true, nil
		}
		return *lat == 0 && *long == 0, nil
	case propertytype.TypeRange:
		if rangeTo == nil || rangeFrom == nil {
			return true, nil
		}
		return *rangeTo == 0 && *rangeFrom == 0, nil
	case propertytype.TypeBool:
		return boolVal == nil, nil
	case propertytype.TypeFloat:
		return floatVal == nil, nil
	default:
		return false, fmt.Errorf("unknown type: %s", typ)
	}
}

func isEmptyProp(ctx context.Context, propertyType *ent.PropertyType, property *ent.Property) (bool, error) {
	if propertyType.Type == propertytype.TypeNode {
		return isEmptyNodeProp(ctx, property, propertyType.NodeType)
	}
	return isEmptyPrimitiveProp(
		propertyType.Type, property.StringVal, property.IntVal, property.BoolVal,
		property.FloatVal, property.LatitudeVal, property.LongitudeVal, property.RangeToVal, property.RangeFromVal)
}

func isEmptyPropType(propertyType *ent.PropertyType) (bool, error) {
	if propertyType.Type == propertytype.TypeEnum || propertyType.Type == propertytype.TypeNode {
		return true, nil
	}
	return isEmptyPrimitiveProp(
		propertyType.Type, propertyType.StringVal, propertyType.IntVal, propertyType.BoolVal,
		propertyType.FloatVal, propertyType.LatitudeVal, propertyType.LongitudeVal, propertyType.RangeToVal, propertyType.RangeFromVal)
}

func getWorkOrderTemplateID(ctx context.Context, m *ent.WorkOrderMutation) (int, error) {
	client := m.Client()
	if m.Op().Is(ent.OpCreate) {
		templateID, exists := m.TemplateID()
		if !exists {
			return templateID, errWorkOrderTemplateNotFound
		}
		return templateID, nil
	}
	id, exists := m.ID()
	if !exists {
		return id, errWorkOrderTemplateNotFound
	}
	templateID, err := client.WorkOrder.Query().
		Where(workorder.ID(id)).
		QueryTemplate().
		OnlyID(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return templateID, errWorkOrderTemplateNotFound
		}
		return templateID, err
	}
	return templateID, nil
}

func isWorkOrderClosed(ctx context.Context, client *ent.Client, templateID int) (bool, error) {
	workOrder, err := client.WorkOrder.Query().
		Where(workorder.HasTemplateWith(workordertemplate.ID(templateID))).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to query work order: %w", err)
	}
	return workOrder.Status == workorder.StatusDONE, nil
}

func mandatoryPropertyOnClose(ctx context.Context, client *ent.Client, templateID int) error {
	propertyTypes, err := client.WorkOrderTemplate.Query().
		Where(workordertemplate.ID(templateID)).
		QueryPropertyTypes().
		Where(propertytype.And(
			propertytype.Mandatory(true),
			propertytype.IsInstanceProperty(true),
		)).
		WithProperties().
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed to query work order mandatory property types: %w", err)
	}
	if len(propertyTypes) == 0 {
		return nil
	}
	for _, propType := range propertyTypes {
		properties := propType.Edges.Properties
		var (
			emptyProp bool
			err       error
		)
		switch {
		case len(properties) >= 2:
			return fmt.Errorf("multiple properties of type %d were found", propType.ID)
		case len(properties) == 0:
			emptyProp, err = isEmptyPropType(propType)
		case len(properties) == 1:
			emptyProp, err = isEmptyProp(ctx, propType, properties[0])
		}
		if err != nil {
			return err
		}
		if emptyProp {
			return fmt.Errorf("property of type %d is empty", propType.ID)
		}
	}
	return nil
}

func isClosingWorkOrder(ctx context.Context, m *ent.WorkOrderMutation) (bool, error) {
	newStatus, exists := m.Status()
	if !exists {
		return false, nil
	}
	if m.Op().Is(ent.OpCreate) {
		return newStatus == workorder.StatusDONE, nil
	}
	oldStatus, err := m.OldStatus(ctx)
	if err != nil {
		return false, err
	}
	return oldStatus != newStatus && newStatus == workorder.StatusDONE, nil
}

func rollbackAndErr(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); err != nil {
		err = fmt.Errorf("%v: %v", err, rerr)
	}
	return err
}

func WorkOrderMandatoryPropertyOnClose() ent.Hook {
	hk := func(next ent.Mutator) ent.Mutator {
		return hook.WorkOrderFunc(func(ctx context.Context, m *ent.WorkOrderMutation) (ent.Value, error) {
			if !viewer.FromContext(ctx).Features().Enabled(viewer.FeatureMandatoryPropertiesOnWorkOrderClose) {
				return next.Mutate(ctx, m)
			}
			closing, err := isClosingWorkOrder(ctx, m)
			if err != nil {
				return nil, fmt.Errorf("failed to check if work order is closed: %w", err)
			}
			if !closing {
				return next.Mutate(ctx, m)
			}
			templateID, err := getWorkOrderTemplateID(ctx, m)
			if err != nil {
				if err == errWorkOrderTemplateNotFound {
					return next.Mutate(ctx, m)
				}
				return nil, fmt.Errorf("failed to get work order template: %w", err)
			}
			tx := ent.TxFromContext(ctx)
			if tx == nil {
				if err := mandatoryPropertyOnClose(ctx, m.Client(), templateID); err != nil {
					return nil, fmt.Errorf("mandatory properties in work order are missing: %w", err)
				}
			} else {
				tx.OnCommit(func(next ent.Committer) ent.Committer {
					return ent.CommitFunc(func(ctx context.Context, tx *ent.Tx) error {
						closed, err := isWorkOrderClosed(ctx, tx.Client(), templateID)
						if err != nil {
							return rollbackAndErr(tx, err)
						}
						if !closed {
							return next.Commit(ctx, tx)
						}
						if err := mandatoryPropertyOnClose(ctx, tx.Client(), templateID); err != nil {
							return rollbackAndErr(
								tx, fmt.Errorf("mandatory properties in work order are missing: %w", err))
						}
						return next.Commit(ctx, tx)
					})
				})
			}
			return next.Mutate(ctx, m)
		})
	}
	return hook.On(hk, ent.OpUpdateOne|ent.OpCreate)
}
