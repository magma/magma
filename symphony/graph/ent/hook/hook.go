// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
)

// The ActionsRuleFunc type is an adapter to allow the use of ordinary
// function as ActionsRule mutator.
type ActionsRuleFunc func(context.Context, *ent.ActionsRuleMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ActionsRuleFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ActionsRuleMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ActionsRuleMutation", m)
	}
	return f(ctx, mv)
}

// The CheckListCategoryFunc type is an adapter to allow the use of ordinary
// function as CheckListCategory mutator.
type CheckListCategoryFunc func(context.Context, *ent.CheckListCategoryMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CheckListCategoryFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CheckListCategoryMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CheckListCategoryMutation", m)
	}
	return f(ctx, mv)
}

// The CheckListItemFunc type is an adapter to allow the use of ordinary
// function as CheckListItem mutator.
type CheckListItemFunc func(context.Context, *ent.CheckListItemMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CheckListItemFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CheckListItemMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CheckListItemMutation", m)
	}
	return f(ctx, mv)
}

// The CheckListItemDefinitionFunc type is an adapter to allow the use of ordinary
// function as CheckListItemDefinition mutator.
type CheckListItemDefinitionFunc func(context.Context, *ent.CheckListItemDefinitionMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CheckListItemDefinitionFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CheckListItemDefinitionMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CheckListItemDefinitionMutation", m)
	}
	return f(ctx, mv)
}

// The CommentFunc type is an adapter to allow the use of ordinary
// function as Comment mutator.
type CommentFunc func(context.Context, *ent.CommentMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CommentFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CommentMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CommentMutation", m)
	}
	return f(ctx, mv)
}

// The CustomerFunc type is an adapter to allow the use of ordinary
// function as Customer mutator.
type CustomerFunc func(context.Context, *ent.CustomerMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f CustomerFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.CustomerMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.CustomerMutation", m)
	}
	return f(ctx, mv)
}

// The EquipmentFunc type is an adapter to allow the use of ordinary
// function as Equipment mutator.
type EquipmentFunc func(context.Context, *ent.EquipmentMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f EquipmentFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.EquipmentMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.EquipmentMutation", m)
	}
	return f(ctx, mv)
}

// The EquipmentCategoryFunc type is an adapter to allow the use of ordinary
// function as EquipmentCategory mutator.
type EquipmentCategoryFunc func(context.Context, *ent.EquipmentCategoryMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f EquipmentCategoryFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.EquipmentCategoryMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.EquipmentCategoryMutation", m)
	}
	return f(ctx, mv)
}

// The EquipmentPortFunc type is an adapter to allow the use of ordinary
// function as EquipmentPort mutator.
type EquipmentPortFunc func(context.Context, *ent.EquipmentPortMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f EquipmentPortFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.EquipmentPortMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.EquipmentPortMutation", m)
	}
	return f(ctx, mv)
}

// The EquipmentPortDefinitionFunc type is an adapter to allow the use of ordinary
// function as EquipmentPortDefinition mutator.
type EquipmentPortDefinitionFunc func(context.Context, *ent.EquipmentPortDefinitionMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f EquipmentPortDefinitionFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.EquipmentPortDefinitionMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.EquipmentPortDefinitionMutation", m)
	}
	return f(ctx, mv)
}

// The EquipmentPortTypeFunc type is an adapter to allow the use of ordinary
// function as EquipmentPortType mutator.
type EquipmentPortTypeFunc func(context.Context, *ent.EquipmentPortTypeMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f EquipmentPortTypeFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.EquipmentPortTypeMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.EquipmentPortTypeMutation", m)
	}
	return f(ctx, mv)
}

// The EquipmentPositionFunc type is an adapter to allow the use of ordinary
// function as EquipmentPosition mutator.
type EquipmentPositionFunc func(context.Context, *ent.EquipmentPositionMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f EquipmentPositionFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.EquipmentPositionMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.EquipmentPositionMutation", m)
	}
	return f(ctx, mv)
}

// The EquipmentPositionDefinitionFunc type is an adapter to allow the use of ordinary
// function as EquipmentPositionDefinition mutator.
type EquipmentPositionDefinitionFunc func(context.Context, *ent.EquipmentPositionDefinitionMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f EquipmentPositionDefinitionFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.EquipmentPositionDefinitionMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.EquipmentPositionDefinitionMutation", m)
	}
	return f(ctx, mv)
}

// The EquipmentTypeFunc type is an adapter to allow the use of ordinary
// function as EquipmentType mutator.
type EquipmentTypeFunc func(context.Context, *ent.EquipmentTypeMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f EquipmentTypeFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.EquipmentTypeMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.EquipmentTypeMutation", m)
	}
	return f(ctx, mv)
}

// The FileFunc type is an adapter to allow the use of ordinary
// function as File mutator.
type FileFunc func(context.Context, *ent.FileMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f FileFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.FileMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.FileMutation", m)
	}
	return f(ctx, mv)
}

// The FloorPlanFunc type is an adapter to allow the use of ordinary
// function as FloorPlan mutator.
type FloorPlanFunc func(context.Context, *ent.FloorPlanMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f FloorPlanFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.FloorPlanMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.FloorPlanMutation", m)
	}
	return f(ctx, mv)
}

// The FloorPlanReferencePointFunc type is an adapter to allow the use of ordinary
// function as FloorPlanReferencePoint mutator.
type FloorPlanReferencePointFunc func(context.Context, *ent.FloorPlanReferencePointMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f FloorPlanReferencePointFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.FloorPlanReferencePointMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.FloorPlanReferencePointMutation", m)
	}
	return f(ctx, mv)
}

// The FloorPlanScaleFunc type is an adapter to allow the use of ordinary
// function as FloorPlanScale mutator.
type FloorPlanScaleFunc func(context.Context, *ent.FloorPlanScaleMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f FloorPlanScaleFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.FloorPlanScaleMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.FloorPlanScaleMutation", m)
	}
	return f(ctx, mv)
}

// The HyperlinkFunc type is an adapter to allow the use of ordinary
// function as Hyperlink mutator.
type HyperlinkFunc func(context.Context, *ent.HyperlinkMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f HyperlinkFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.HyperlinkMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.HyperlinkMutation", m)
	}
	return f(ctx, mv)
}

// The LinkFunc type is an adapter to allow the use of ordinary
// function as Link mutator.
type LinkFunc func(context.Context, *ent.LinkMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f LinkFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.LinkMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.LinkMutation", m)
	}
	return f(ctx, mv)
}

// The LocationFunc type is an adapter to allow the use of ordinary
// function as Location mutator.
type LocationFunc func(context.Context, *ent.LocationMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f LocationFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.LocationMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.LocationMutation", m)
	}
	return f(ctx, mv)
}

// The LocationTypeFunc type is an adapter to allow the use of ordinary
// function as LocationType mutator.
type LocationTypeFunc func(context.Context, *ent.LocationTypeMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f LocationTypeFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.LocationTypeMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.LocationTypeMutation", m)
	}
	return f(ctx, mv)
}

// The ProjectFunc type is an adapter to allow the use of ordinary
// function as Project mutator.
type ProjectFunc func(context.Context, *ent.ProjectMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ProjectFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ProjectMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ProjectMutation", m)
	}
	return f(ctx, mv)
}

// The ProjectTypeFunc type is an adapter to allow the use of ordinary
// function as ProjectType mutator.
type ProjectTypeFunc func(context.Context, *ent.ProjectTypeMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ProjectTypeFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ProjectTypeMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ProjectTypeMutation", m)
	}
	return f(ctx, mv)
}

// The PropertyFunc type is an adapter to allow the use of ordinary
// function as Property mutator.
type PropertyFunc func(context.Context, *ent.PropertyMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f PropertyFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.PropertyMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.PropertyMutation", m)
	}
	return f(ctx, mv)
}

// The PropertyTypeFunc type is an adapter to allow the use of ordinary
// function as PropertyType mutator.
type PropertyTypeFunc func(context.Context, *ent.PropertyTypeMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f PropertyTypeFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.PropertyTypeMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.PropertyTypeMutation", m)
	}
	return f(ctx, mv)
}

// The ReportFilterFunc type is an adapter to allow the use of ordinary
// function as ReportFilter mutator.
type ReportFilterFunc func(context.Context, *ent.ReportFilterMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ReportFilterFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ReportFilterMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ReportFilterMutation", m)
	}
	return f(ctx, mv)
}

// The ServiceFunc type is an adapter to allow the use of ordinary
// function as Service mutator.
type ServiceFunc func(context.Context, *ent.ServiceMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ServiceFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ServiceMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ServiceMutation", m)
	}
	return f(ctx, mv)
}

// The ServiceEndpointFunc type is an adapter to allow the use of ordinary
// function as ServiceEndpoint mutator.
type ServiceEndpointFunc func(context.Context, *ent.ServiceEndpointMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ServiceEndpointFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ServiceEndpointMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ServiceEndpointMutation", m)
	}
	return f(ctx, mv)
}

// The ServiceTypeFunc type is an adapter to allow the use of ordinary
// function as ServiceType mutator.
type ServiceTypeFunc func(context.Context, *ent.ServiceTypeMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ServiceTypeFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.ServiceTypeMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ServiceTypeMutation", m)
	}
	return f(ctx, mv)
}

// The SurveyFunc type is an adapter to allow the use of ordinary
// function as Survey mutator.
type SurveyFunc func(context.Context, *ent.SurveyMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f SurveyFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.SurveyMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.SurveyMutation", m)
	}
	return f(ctx, mv)
}

// The SurveyCellScanFunc type is an adapter to allow the use of ordinary
// function as SurveyCellScan mutator.
type SurveyCellScanFunc func(context.Context, *ent.SurveyCellScanMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f SurveyCellScanFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.SurveyCellScanMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.SurveyCellScanMutation", m)
	}
	return f(ctx, mv)
}

// The SurveyQuestionFunc type is an adapter to allow the use of ordinary
// function as SurveyQuestion mutator.
type SurveyQuestionFunc func(context.Context, *ent.SurveyQuestionMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f SurveyQuestionFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.SurveyQuestionMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.SurveyQuestionMutation", m)
	}
	return f(ctx, mv)
}

// The SurveyTemplateCategoryFunc type is an adapter to allow the use of ordinary
// function as SurveyTemplateCategory mutator.
type SurveyTemplateCategoryFunc func(context.Context, *ent.SurveyTemplateCategoryMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f SurveyTemplateCategoryFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.SurveyTemplateCategoryMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.SurveyTemplateCategoryMutation", m)
	}
	return f(ctx, mv)
}

// The SurveyTemplateQuestionFunc type is an adapter to allow the use of ordinary
// function as SurveyTemplateQuestion mutator.
type SurveyTemplateQuestionFunc func(context.Context, *ent.SurveyTemplateQuestionMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f SurveyTemplateQuestionFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.SurveyTemplateQuestionMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.SurveyTemplateQuestionMutation", m)
	}
	return f(ctx, mv)
}

// The SurveyWiFiScanFunc type is an adapter to allow the use of ordinary
// function as SurveyWiFiScan mutator.
type SurveyWiFiScanFunc func(context.Context, *ent.SurveyWiFiScanMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f SurveyWiFiScanFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.SurveyWiFiScanMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.SurveyWiFiScanMutation", m)
	}
	return f(ctx, mv)
}

// The TechnicianFunc type is an adapter to allow the use of ordinary
// function as Technician mutator.
type TechnicianFunc func(context.Context, *ent.TechnicianMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f TechnicianFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.TechnicianMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.TechnicianMutation", m)
	}
	return f(ctx, mv)
}

// The UserFunc type is an adapter to allow the use of ordinary
// function as User mutator.
type UserFunc func(context.Context, *ent.UserMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f UserFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.UserMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.UserMutation", m)
	}
	return f(ctx, mv)
}

// The WorkOrderFunc type is an adapter to allow the use of ordinary
// function as WorkOrder mutator.
type WorkOrderFunc func(context.Context, *ent.WorkOrderMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f WorkOrderFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.WorkOrderMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.WorkOrderMutation", m)
	}
	return f(ctx, mv)
}

// The WorkOrderDefinitionFunc type is an adapter to allow the use of ordinary
// function as WorkOrderDefinition mutator.
type WorkOrderDefinitionFunc func(context.Context, *ent.WorkOrderDefinitionMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f WorkOrderDefinitionFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.WorkOrderDefinitionMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.WorkOrderDefinitionMutation", m)
	}
	return f(ctx, mv)
}

// The WorkOrderTypeFunc type is an adapter to allow the use of ordinary
// function as WorkOrderType mutator.
type WorkOrderTypeFunc func(context.Context, *ent.WorkOrderTypeMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f WorkOrderTypeFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	mv, ok := m.(*ent.WorkOrderTypeMutation)
	if !ok {
		return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.WorkOrderTypeMutation", m)
	}
	return f(ctx, mv)
}

// On executes the given hook only of the given operation.
//
//	hook.On(Log, ent.Delete|ent.Create)
//
func On(hk ent.Hook, op ent.Op) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if m.Op().Is(op) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []ent.Hook {
//		return []ent.Hook{
//			Reject(ent.Delete|ent.Update),
//		}
//	}
//
func Reject(op ent.Op) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if m.Op().Is(op) {
				return nil, fmt.Errorf("%s operation is not allowed", m.Op())
			}
			return next.Mutate(ctx, m)
		})
	}
}
