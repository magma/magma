// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package privacy

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent"
)

var (
	// Allow may be returned by read/write rules to indicate that the policy
	// evaluation should terminate with an allow decision.
	Allow = errors.New("ent/privacy: allow rule")

	// Deny may be returned by read/write rules to indicate that the policy
	// evaluation should terminate with an deny decision.
	Deny = errors.New("ent/privacy: deny rule")

	// Skip may be returned by read/write rules to indicate that the policy
	// evaluation should continue to the next rule.
	Skip = errors.New("ent/privacy: skip rule")
)

// Allowf returns an formatted wrapped Allow decision.
func Allowf(format string, a ...interface{}) error {
	return fmt.Errorf(format+": %w", append(a, Allow)...)
}

// Denyf returns an formatted wrapped Deny decision.
func Denyf(format string, a ...interface{}) error {
	return fmt.Errorf(format+": %w", append(a, Deny)...)
}

// Skipf returns an formatted wrapped Skip decision.
func Skipf(format string, a ...interface{}) error {
	return fmt.Errorf(format+": %w", append(a, Skip)...)
}

type (
	// ReadPolicy combines multiple read rules into a single policy.
	ReadPolicy []ReadRule

	// ReadRule defines the interface deciding whether a read is allowed.
	ReadRule interface {
		EvalRead(context.Context, ent.Value) error
	}
)

// EvalRead evaluates a load against a read policy.
func (policy ReadPolicy) EvalRead(ctx context.Context, v ent.Value) error {
	for _, rule := range policy {
		switch err := rule.EvalRead(ctx, v); {
		case err == nil || errors.Is(err, Skip):
		case errors.Is(err, Allow):
			return nil
		default:
			return err
		}
	}
	return nil
}

// ReadRuleFunc type is an adapter to allow the use of
// ordinary functions as read rules.
type ReadRuleFunc func(context.Context, ent.Value) error

// Eval calls f(ctx, v).
func (f ReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	return f(ctx, v)
}

type (
	// WritePolicy combines multiple write rules into a single policy.
	WritePolicy []WriteRule

	// WriteRule defines the interface deciding whether a write is allowed.
	WriteRule interface {
		EvalWrite(context.Context, ent.Mutation) error
	}
)

// EvalWrite evaluates a mutation against a write policy.
func (policy WritePolicy) EvalWrite(ctx context.Context, m ent.Mutation) error {
	for _, rule := range policy {
		switch err := rule.EvalWrite(ctx, m); {
		case err == nil || errors.Is(err, Skip):
		case errors.Is(err, Allow):
			return nil
		default:
			return err
		}
	}
	return nil
}

// WriteRuleFunc type is an adapter to allow the use of
// ordinary functions as write rules.
type WriteRuleFunc func(context.Context, ent.Mutation) error

// Eval calls f(ctx, m).
func (f WriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	return f(ctx, m)
}

// Policy groups read and write policies.
type Policy struct {
	Read  ReadPolicy
	Write WritePolicy
}

// EvalRead forwards evaluation to read policy.
func (policy Policy) EvalRead(ctx context.Context, v ent.Value) error {
	return policy.Read.EvalRead(ctx, v)
}

// EvalWrite forwards evaluation to write policy.
func (policy Policy) EvalWrite(ctx context.Context, m ent.Mutation) error {
	return policy.Write.EvalWrite(ctx, m)
}

// ReadWriteRule is the interface that groups read and write rules.
type ReadWriteRule interface {
	ReadRule
	WriteRule
}

// AlwaysAllowRule returns a read/write rule that returns an allow decision.
func AlwaysAllowRule() ReadWriteRule {
	return fixedDecisionRule{Allow}
}

// AlwaysDenyRule returns a read/write rule that returns a deny decision.
func AlwaysDenyRule() ReadWriteRule {
	return fixedDecisionRule{Deny}
}

type fixedDecisionRule struct{ err error }

func (f fixedDecisionRule) EvalRead(context.Context, ent.Value) error     { return f.err }
func (f fixedDecisionRule) EvalWrite(context.Context, ent.Mutation) error { return f.err }

// The ActionsRuleReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type ActionsRuleReadRuleFunc func(context.Context, *ent.ActionsRule) error

// EvalRead calls f(ctx, v).
func (f ActionsRuleReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.ActionsRule); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.ActionsRule", v)
}

// The ActionsRuleWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type ActionsRuleWriteRuleFunc func(context.Context, *ent.ActionsRuleMutation) error

// EvalWrite calls f(ctx, m).
func (f ActionsRuleWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ActionsRuleMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ActionsRuleMutation", m)
}

// The CheckListCategoryReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type CheckListCategoryReadRuleFunc func(context.Context, *ent.CheckListCategory) error

// EvalRead calls f(ctx, v).
func (f CheckListCategoryReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.CheckListCategory); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.CheckListCategory", v)
}

// The CheckListCategoryWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type CheckListCategoryWriteRuleFunc func(context.Context, *ent.CheckListCategoryMutation) error

// EvalWrite calls f(ctx, m).
func (f CheckListCategoryWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CheckListCategoryMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CheckListCategoryMutation", m)
}

// The CheckListItemReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type CheckListItemReadRuleFunc func(context.Context, *ent.CheckListItem) error

// EvalRead calls f(ctx, v).
func (f CheckListItemReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.CheckListItem); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.CheckListItem", v)
}

// The CheckListItemWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type CheckListItemWriteRuleFunc func(context.Context, *ent.CheckListItemMutation) error

// EvalWrite calls f(ctx, m).
func (f CheckListItemWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CheckListItemMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CheckListItemMutation", m)
}

// The CheckListItemDefinitionReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type CheckListItemDefinitionReadRuleFunc func(context.Context, *ent.CheckListItemDefinition) error

// EvalRead calls f(ctx, v).
func (f CheckListItemDefinitionReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.CheckListItemDefinition); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.CheckListItemDefinition", v)
}

// The CheckListItemDefinitionWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type CheckListItemDefinitionWriteRuleFunc func(context.Context, *ent.CheckListItemDefinitionMutation) error

// EvalWrite calls f(ctx, m).
func (f CheckListItemDefinitionWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CheckListItemDefinitionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CheckListItemDefinitionMutation", m)
}

// The CommentReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type CommentReadRuleFunc func(context.Context, *ent.Comment) error

// EvalRead calls f(ctx, v).
func (f CommentReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Comment); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Comment", v)
}

// The CommentWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type CommentWriteRuleFunc func(context.Context, *ent.CommentMutation) error

// EvalWrite calls f(ctx, m).
func (f CommentWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CommentMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CommentMutation", m)
}

// The CustomerReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type CustomerReadRuleFunc func(context.Context, *ent.Customer) error

// EvalRead calls f(ctx, v).
func (f CustomerReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Customer); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Customer", v)
}

// The CustomerWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type CustomerWriteRuleFunc func(context.Context, *ent.CustomerMutation) error

// EvalWrite calls f(ctx, m).
func (f CustomerWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CustomerMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CustomerMutation", m)
}

// The EquipmentReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type EquipmentReadRuleFunc func(context.Context, *ent.Equipment) error

// EvalRead calls f(ctx, v).
func (f EquipmentReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Equipment); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Equipment", v)
}

// The EquipmentWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type EquipmentWriteRuleFunc func(context.Context, *ent.EquipmentMutation) error

// EvalWrite calls f(ctx, m).
func (f EquipmentWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentMutation", m)
}

// The EquipmentCategoryReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type EquipmentCategoryReadRuleFunc func(context.Context, *ent.EquipmentCategory) error

// EvalRead calls f(ctx, v).
func (f EquipmentCategoryReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.EquipmentCategory); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.EquipmentCategory", v)
}

// The EquipmentCategoryWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type EquipmentCategoryWriteRuleFunc func(context.Context, *ent.EquipmentCategoryMutation) error

// EvalWrite calls f(ctx, m).
func (f EquipmentCategoryWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentCategoryMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentCategoryMutation", m)
}

// The EquipmentPortReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type EquipmentPortReadRuleFunc func(context.Context, *ent.EquipmentPort) error

// EvalRead calls f(ctx, v).
func (f EquipmentPortReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.EquipmentPort); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.EquipmentPort", v)
}

// The EquipmentPortWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type EquipmentPortWriteRuleFunc func(context.Context, *ent.EquipmentPortMutation) error

// EvalWrite calls f(ctx, m).
func (f EquipmentPortWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentPortMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentPortMutation", m)
}

// The EquipmentPortDefinitionReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type EquipmentPortDefinitionReadRuleFunc func(context.Context, *ent.EquipmentPortDefinition) error

// EvalRead calls f(ctx, v).
func (f EquipmentPortDefinitionReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.EquipmentPortDefinition); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.EquipmentPortDefinition", v)
}

// The EquipmentPortDefinitionWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type EquipmentPortDefinitionWriteRuleFunc func(context.Context, *ent.EquipmentPortDefinitionMutation) error

// EvalWrite calls f(ctx, m).
func (f EquipmentPortDefinitionWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentPortDefinitionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentPortDefinitionMutation", m)
}

// The EquipmentPortTypeReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type EquipmentPortTypeReadRuleFunc func(context.Context, *ent.EquipmentPortType) error

// EvalRead calls f(ctx, v).
func (f EquipmentPortTypeReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.EquipmentPortType); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.EquipmentPortType", v)
}

// The EquipmentPortTypeWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type EquipmentPortTypeWriteRuleFunc func(context.Context, *ent.EquipmentPortTypeMutation) error

// EvalWrite calls f(ctx, m).
func (f EquipmentPortTypeWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentPortTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentPortTypeMutation", m)
}

// The EquipmentPositionReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type EquipmentPositionReadRuleFunc func(context.Context, *ent.EquipmentPosition) error

// EvalRead calls f(ctx, v).
func (f EquipmentPositionReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.EquipmentPosition); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.EquipmentPosition", v)
}

// The EquipmentPositionWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type EquipmentPositionWriteRuleFunc func(context.Context, *ent.EquipmentPositionMutation) error

// EvalWrite calls f(ctx, m).
func (f EquipmentPositionWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentPositionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentPositionMutation", m)
}

// The EquipmentPositionDefinitionReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type EquipmentPositionDefinitionReadRuleFunc func(context.Context, *ent.EquipmentPositionDefinition) error

// EvalRead calls f(ctx, v).
func (f EquipmentPositionDefinitionReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.EquipmentPositionDefinition); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.EquipmentPositionDefinition", v)
}

// The EquipmentPositionDefinitionWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type EquipmentPositionDefinitionWriteRuleFunc func(context.Context, *ent.EquipmentPositionDefinitionMutation) error

// EvalWrite calls f(ctx, m).
func (f EquipmentPositionDefinitionWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentPositionDefinitionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentPositionDefinitionMutation", m)
}

// The EquipmentTypeReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type EquipmentTypeReadRuleFunc func(context.Context, *ent.EquipmentType) error

// EvalRead calls f(ctx, v).
func (f EquipmentTypeReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.EquipmentType); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.EquipmentType", v)
}

// The EquipmentTypeWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type EquipmentTypeWriteRuleFunc func(context.Context, *ent.EquipmentTypeMutation) error

// EvalWrite calls f(ctx, m).
func (f EquipmentTypeWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentTypeMutation", m)
}

// The FileReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type FileReadRuleFunc func(context.Context, *ent.File) error

// EvalRead calls f(ctx, v).
func (f FileReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.File); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.File", v)
}

// The FileWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type FileWriteRuleFunc func(context.Context, *ent.FileMutation) error

// EvalWrite calls f(ctx, m).
func (f FileWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.FileMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.FileMutation", m)
}

// The FloorPlanReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type FloorPlanReadRuleFunc func(context.Context, *ent.FloorPlan) error

// EvalRead calls f(ctx, v).
func (f FloorPlanReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.FloorPlan); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.FloorPlan", v)
}

// The FloorPlanWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type FloorPlanWriteRuleFunc func(context.Context, *ent.FloorPlanMutation) error

// EvalWrite calls f(ctx, m).
func (f FloorPlanWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.FloorPlanMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.FloorPlanMutation", m)
}

// The FloorPlanReferencePointReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type FloorPlanReferencePointReadRuleFunc func(context.Context, *ent.FloorPlanReferencePoint) error

// EvalRead calls f(ctx, v).
func (f FloorPlanReferencePointReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.FloorPlanReferencePoint); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.FloorPlanReferencePoint", v)
}

// The FloorPlanReferencePointWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type FloorPlanReferencePointWriteRuleFunc func(context.Context, *ent.FloorPlanReferencePointMutation) error

// EvalWrite calls f(ctx, m).
func (f FloorPlanReferencePointWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.FloorPlanReferencePointMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.FloorPlanReferencePointMutation", m)
}

// The FloorPlanScaleReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type FloorPlanScaleReadRuleFunc func(context.Context, *ent.FloorPlanScale) error

// EvalRead calls f(ctx, v).
func (f FloorPlanScaleReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.FloorPlanScale); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.FloorPlanScale", v)
}

// The FloorPlanScaleWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type FloorPlanScaleWriteRuleFunc func(context.Context, *ent.FloorPlanScaleMutation) error

// EvalWrite calls f(ctx, m).
func (f FloorPlanScaleWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.FloorPlanScaleMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.FloorPlanScaleMutation", m)
}

// The HyperlinkReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type HyperlinkReadRuleFunc func(context.Context, *ent.Hyperlink) error

// EvalRead calls f(ctx, v).
func (f HyperlinkReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Hyperlink); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Hyperlink", v)
}

// The HyperlinkWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type HyperlinkWriteRuleFunc func(context.Context, *ent.HyperlinkMutation) error

// EvalWrite calls f(ctx, m).
func (f HyperlinkWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.HyperlinkMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.HyperlinkMutation", m)
}

// The LinkReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type LinkReadRuleFunc func(context.Context, *ent.Link) error

// EvalRead calls f(ctx, v).
func (f LinkReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Link); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Link", v)
}

// The LinkWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type LinkWriteRuleFunc func(context.Context, *ent.LinkMutation) error

// EvalWrite calls f(ctx, m).
func (f LinkWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.LinkMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.LinkMutation", m)
}

// The LocationReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type LocationReadRuleFunc func(context.Context, *ent.Location) error

// EvalRead calls f(ctx, v).
func (f LocationReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Location); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Location", v)
}

// The LocationWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type LocationWriteRuleFunc func(context.Context, *ent.LocationMutation) error

// EvalWrite calls f(ctx, m).
func (f LocationWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.LocationMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.LocationMutation", m)
}

// The LocationTypeReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type LocationTypeReadRuleFunc func(context.Context, *ent.LocationType) error

// EvalRead calls f(ctx, v).
func (f LocationTypeReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.LocationType); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.LocationType", v)
}

// The LocationTypeWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type LocationTypeWriteRuleFunc func(context.Context, *ent.LocationTypeMutation) error

// EvalWrite calls f(ctx, m).
func (f LocationTypeWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.LocationTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.LocationTypeMutation", m)
}

// The ProjectReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type ProjectReadRuleFunc func(context.Context, *ent.Project) error

// EvalRead calls f(ctx, v).
func (f ProjectReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Project); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Project", v)
}

// The ProjectWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type ProjectWriteRuleFunc func(context.Context, *ent.ProjectMutation) error

// EvalWrite calls f(ctx, m).
func (f ProjectWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ProjectMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ProjectMutation", m)
}

// The ProjectTypeReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type ProjectTypeReadRuleFunc func(context.Context, *ent.ProjectType) error

// EvalRead calls f(ctx, v).
func (f ProjectTypeReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.ProjectType); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.ProjectType", v)
}

// The ProjectTypeWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type ProjectTypeWriteRuleFunc func(context.Context, *ent.ProjectTypeMutation) error

// EvalWrite calls f(ctx, m).
func (f ProjectTypeWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ProjectTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ProjectTypeMutation", m)
}

// The PropertyReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type PropertyReadRuleFunc func(context.Context, *ent.Property) error

// EvalRead calls f(ctx, v).
func (f PropertyReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Property); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Property", v)
}

// The PropertyWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type PropertyWriteRuleFunc func(context.Context, *ent.PropertyMutation) error

// EvalWrite calls f(ctx, m).
func (f PropertyWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.PropertyMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.PropertyMutation", m)
}

// The PropertyTypeReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type PropertyTypeReadRuleFunc func(context.Context, *ent.PropertyType) error

// EvalRead calls f(ctx, v).
func (f PropertyTypeReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.PropertyType); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.PropertyType", v)
}

// The PropertyTypeWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type PropertyTypeWriteRuleFunc func(context.Context, *ent.PropertyTypeMutation) error

// EvalWrite calls f(ctx, m).
func (f PropertyTypeWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.PropertyTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.PropertyTypeMutation", m)
}

// The ReportFilterReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type ReportFilterReadRuleFunc func(context.Context, *ent.ReportFilter) error

// EvalRead calls f(ctx, v).
func (f ReportFilterReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.ReportFilter); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.ReportFilter", v)
}

// The ReportFilterWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type ReportFilterWriteRuleFunc func(context.Context, *ent.ReportFilterMutation) error

// EvalWrite calls f(ctx, m).
func (f ReportFilterWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ReportFilterMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ReportFilterMutation", m)
}

// The ServiceReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type ServiceReadRuleFunc func(context.Context, *ent.Service) error

// EvalRead calls f(ctx, v).
func (f ServiceReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Service); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Service", v)
}

// The ServiceWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type ServiceWriteRuleFunc func(context.Context, *ent.ServiceMutation) error

// EvalWrite calls f(ctx, m).
func (f ServiceWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ServiceMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ServiceMutation", m)
}

// The ServiceEndpointReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type ServiceEndpointReadRuleFunc func(context.Context, *ent.ServiceEndpoint) error

// EvalRead calls f(ctx, v).
func (f ServiceEndpointReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.ServiceEndpoint); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.ServiceEndpoint", v)
}

// The ServiceEndpointWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type ServiceEndpointWriteRuleFunc func(context.Context, *ent.ServiceEndpointMutation) error

// EvalWrite calls f(ctx, m).
func (f ServiceEndpointWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ServiceEndpointMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ServiceEndpointMutation", m)
}

// The ServiceTypeReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type ServiceTypeReadRuleFunc func(context.Context, *ent.ServiceType) error

// EvalRead calls f(ctx, v).
func (f ServiceTypeReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.ServiceType); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.ServiceType", v)
}

// The ServiceTypeWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type ServiceTypeWriteRuleFunc func(context.Context, *ent.ServiceTypeMutation) error

// EvalWrite calls f(ctx, m).
func (f ServiceTypeWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ServiceTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ServiceTypeMutation", m)
}

// The SurveyReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type SurveyReadRuleFunc func(context.Context, *ent.Survey) error

// EvalRead calls f(ctx, v).
func (f SurveyReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Survey); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Survey", v)
}

// The SurveyWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type SurveyWriteRuleFunc func(context.Context, *ent.SurveyMutation) error

// EvalWrite calls f(ctx, m).
func (f SurveyWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyMutation", m)
}

// The SurveyCellScanReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type SurveyCellScanReadRuleFunc func(context.Context, *ent.SurveyCellScan) error

// EvalRead calls f(ctx, v).
func (f SurveyCellScanReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.SurveyCellScan); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.SurveyCellScan", v)
}

// The SurveyCellScanWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type SurveyCellScanWriteRuleFunc func(context.Context, *ent.SurveyCellScanMutation) error

// EvalWrite calls f(ctx, m).
func (f SurveyCellScanWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyCellScanMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyCellScanMutation", m)
}

// The SurveyQuestionReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type SurveyQuestionReadRuleFunc func(context.Context, *ent.SurveyQuestion) error

// EvalRead calls f(ctx, v).
func (f SurveyQuestionReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.SurveyQuestion); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.SurveyQuestion", v)
}

// The SurveyQuestionWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type SurveyQuestionWriteRuleFunc func(context.Context, *ent.SurveyQuestionMutation) error

// EvalWrite calls f(ctx, m).
func (f SurveyQuestionWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyQuestionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyQuestionMutation", m)
}

// The SurveyTemplateCategoryReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type SurveyTemplateCategoryReadRuleFunc func(context.Context, *ent.SurveyTemplateCategory) error

// EvalRead calls f(ctx, v).
func (f SurveyTemplateCategoryReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.SurveyTemplateCategory); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.SurveyTemplateCategory", v)
}

// The SurveyTemplateCategoryWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type SurveyTemplateCategoryWriteRuleFunc func(context.Context, *ent.SurveyTemplateCategoryMutation) error

// EvalWrite calls f(ctx, m).
func (f SurveyTemplateCategoryWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyTemplateCategoryMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyTemplateCategoryMutation", m)
}

// The SurveyTemplateQuestionReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type SurveyTemplateQuestionReadRuleFunc func(context.Context, *ent.SurveyTemplateQuestion) error

// EvalRead calls f(ctx, v).
func (f SurveyTemplateQuestionReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.SurveyTemplateQuestion); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.SurveyTemplateQuestion", v)
}

// The SurveyTemplateQuestionWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type SurveyTemplateQuestionWriteRuleFunc func(context.Context, *ent.SurveyTemplateQuestionMutation) error

// EvalWrite calls f(ctx, m).
func (f SurveyTemplateQuestionWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyTemplateQuestionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyTemplateQuestionMutation", m)
}

// The SurveyWiFiScanReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type SurveyWiFiScanReadRuleFunc func(context.Context, *ent.SurveyWiFiScan) error

// EvalRead calls f(ctx, v).
func (f SurveyWiFiScanReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.SurveyWiFiScan); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.SurveyWiFiScan", v)
}

// The SurveyWiFiScanWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type SurveyWiFiScanWriteRuleFunc func(context.Context, *ent.SurveyWiFiScanMutation) error

// EvalWrite calls f(ctx, m).
func (f SurveyWiFiScanWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyWiFiScanMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyWiFiScanMutation", m)
}

// The TechnicianReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type TechnicianReadRuleFunc func(context.Context, *ent.Technician) error

// EvalRead calls f(ctx, v).
func (f TechnicianReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.Technician); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.Technician", v)
}

// The TechnicianWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type TechnicianWriteRuleFunc func(context.Context, *ent.TechnicianMutation) error

// EvalWrite calls f(ctx, m).
func (f TechnicianWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.TechnicianMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.TechnicianMutation", m)
}

// The UserReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type UserReadRuleFunc func(context.Context, *ent.User) error

// EvalRead calls f(ctx, v).
func (f UserReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.User); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.User", v)
}

// The UserWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type UserWriteRuleFunc func(context.Context, *ent.UserMutation) error

// EvalWrite calls f(ctx, m).
func (f UserWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.UserMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.UserMutation", m)
}

// The UsersGroupReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type UsersGroupReadRuleFunc func(context.Context, *ent.UsersGroup) error

// EvalRead calls f(ctx, v).
func (f UsersGroupReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.UsersGroup); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.UsersGroup", v)
}

// The UsersGroupWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type UsersGroupWriteRuleFunc func(context.Context, *ent.UsersGroupMutation) error

// EvalWrite calls f(ctx, m).
func (f UsersGroupWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.UsersGroupMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.UsersGroupMutation", m)
}

// The WorkOrderReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type WorkOrderReadRuleFunc func(context.Context, *ent.WorkOrder) error

// EvalRead calls f(ctx, v).
func (f WorkOrderReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.WorkOrder); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.WorkOrder", v)
}

// The WorkOrderWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type WorkOrderWriteRuleFunc func(context.Context, *ent.WorkOrderMutation) error

// EvalWrite calls f(ctx, m).
func (f WorkOrderWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.WorkOrderMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.WorkOrderMutation", m)
}

// The WorkOrderDefinitionReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type WorkOrderDefinitionReadRuleFunc func(context.Context, *ent.WorkOrderDefinition) error

// EvalRead calls f(ctx, v).
func (f WorkOrderDefinitionReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.WorkOrderDefinition); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.WorkOrderDefinition", v)
}

// The WorkOrderDefinitionWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type WorkOrderDefinitionWriteRuleFunc func(context.Context, *ent.WorkOrderDefinitionMutation) error

// EvalWrite calls f(ctx, m).
func (f WorkOrderDefinitionWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.WorkOrderDefinitionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.WorkOrderDefinitionMutation", m)
}

// The WorkOrderTypeReadRuleFunc type is an adapter to allow the use of ordinary
// functions as a read rule.
type WorkOrderTypeReadRuleFunc func(context.Context, *ent.WorkOrderType) error

// EvalRead calls f(ctx, v).
func (f WorkOrderTypeReadRuleFunc) EvalRead(ctx context.Context, v ent.Value) error {
	if v, ok := v.(*ent.WorkOrderType); ok {
		return f(ctx, v)
	}
	return Denyf("ent/privacy: unexpected value type %T, expect *ent.WorkOrderType", v)
}

// The WorkOrderTypeWriteRuleFunc type is an adapter to allow the use of ordinary
// functions as a write rule.
type WorkOrderTypeWriteRuleFunc func(context.Context, *ent.WorkOrderTypeMutation) error

// EvalWrite calls f(ctx, m).
func (f WorkOrderTypeWriteRuleFunc) EvalWrite(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.WorkOrderTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.WorkOrderTypeMutation", m)
}
