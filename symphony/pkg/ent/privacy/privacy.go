// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package privacy

import (
	"context"
	"errors"
	"fmt"

	"github.com/facebookincubator/symphony/pkg/ent"
)

var (
	// Allow may be returned by rules to indicate that the policy
	// evaluation should terminate with an allow decision.
	Allow = errors.New("ent/privacy: allow rule")

	// Deny may be returned by rules to indicate that the policy
	// evaluation should terminate with an deny decision.
	Deny = errors.New("ent/privacy: deny rule")

	// Skip may be returned by rules to indicate that the policy
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

type decisionCtxKey struct{}

// DecisionContext creates a decision context.
func DecisionContext(parent context.Context, decision error) context.Context {
	if decision == nil || errors.Is(decision, Skip) {
		return parent
	}
	return context.WithValue(parent, decisionCtxKey{}, decision)
}

func decisionFromContext(ctx context.Context) (error, bool) {
	decision, ok := ctx.Value(decisionCtxKey{}).(error)
	if ok && errors.Is(decision, Allow) {
		decision = nil
	}
	return decision, ok
}

type (
	// QueryPolicy combines multiple query rules into a single policy.
	QueryPolicy []QueryRule

	// QueryRule defines the interface deciding whether a
	// query is allowed and optionally modify it.
	QueryRule interface {
		EvalQuery(context.Context, ent.Query) error
	}
)

// EvalQuery evaluates a query against a query policy.
func (policy QueryPolicy) EvalQuery(ctx context.Context, q ent.Query) error {
	if decision, ok := decisionFromContext(ctx); ok {
		return decision
	}
	for _, rule := range policy {
		switch decision := rule.EvalQuery(ctx, q); {
		case decision == nil || errors.Is(decision, Skip):
		case errors.Is(decision, Allow):
			return nil
		default:
			return decision
		}
	}
	return nil
}

// QueryRuleFunc type is an adapter to allow the use of
// ordinary functions as query rules.
type QueryRuleFunc func(context.Context, ent.Query) error

// Eval returns f(ctx, q).
func (f QueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	return f(ctx, q)
}

type (
	// MutationPolicy combines multiple mutation rules into a single policy.
	MutationPolicy []MutationRule

	// MutationRule defines the interface deciding whether a
	// mutation is allowed and optionally modify it.
	MutationRule interface {
		EvalMutation(context.Context, ent.Mutation) error
	}
)

// EvalMutation evaluates a mutation against a mutation policy.
func (policy MutationPolicy) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if decision, ok := decisionFromContext(ctx); ok {
		return decision
	}
	for _, rule := range policy {
		switch decision := rule.EvalMutation(ctx, m); {
		case decision == nil || errors.Is(decision, Skip):
		case errors.Is(decision, Allow):
			return nil
		default:
			return decision
		}
	}
	return nil
}

// MutationRuleFunc type is an adapter to allow the use of
// ordinary functions as mutation rules.
type MutationRuleFunc func(context.Context, ent.Mutation) error

// EvalMutation returns f(ctx, m).
func (f MutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	return f(ctx, m)
}

// Policy groups query and mutation policies.
type Policy struct {
	Query    QueryPolicy
	Mutation MutationPolicy
}

// EvalQuery forwards evaluation to query policy.
func (policy Policy) EvalQuery(ctx context.Context, q ent.Query) error {
	return policy.Query.EvalQuery(ctx, q)
}

// EvalMutation forwards evaluation to mutation policy.
func (policy Policy) EvalMutation(ctx context.Context, m ent.Mutation) error {
	return policy.Mutation.EvalMutation(ctx, m)
}

// QueryMutationRule is the interface that groups query and mutation rules.
type QueryMutationRule interface {
	QueryRule
	MutationRule
}

// AlwaysAllowRule returns a rule that returns an allow decision.
func AlwaysAllowRule() QueryMutationRule {
	return fixedDecision{Allow}
}

// AlwaysDenyRule returns a rule that returns a deny decision.
func AlwaysDenyRule() QueryMutationRule {
	return fixedDecision{Deny}
}

type fixedDecision struct {
	decision error
}

func (f fixedDecision) EvalQuery(context.Context, ent.Query) error {
	return f.decision
}

func (f fixedDecision) EvalMutation(context.Context, ent.Mutation) error {
	return f.decision
}

type contextDecision struct {
	eval func(context.Context) error
}

// ContextQueryMutationRule creates a query/mutation rule from a context eval func.
func ContextQueryMutationRule(eval func(context.Context) error) QueryMutationRule {
	return contextDecision{eval}
}

func (c contextDecision) EvalQuery(ctx context.Context, _ ent.Query) error {
	return c.eval(ctx)
}

func (c contextDecision) EvalMutation(ctx context.Context, _ ent.Mutation) error {
	return c.eval(ctx)
}

// OnMutationOperation evaluates the given rule only on a given mutation operation.
func OnMutationOperation(rule MutationRule, op ent.Op) MutationRule {
	return MutationRuleFunc(func(ctx context.Context, m ent.Mutation) error {
		if m.Op().Is(op) {
			return rule.EvalMutation(ctx, m)
		}
		return Skip
	})
}

// DenyMutationOperationRule returns a rule denying specified mutation operation.
func DenyMutationOperationRule(op ent.Op) MutationRule {
	rule := MutationRuleFunc(func(_ context.Context, m ent.Mutation) error {
		return Denyf("ent/privacy: operation %s is not allowed", m.Op())
	})
	return OnMutationOperation(rule, op)
}

// The ActionsRuleQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type ActionsRuleQueryRuleFunc func(context.Context, *ent.ActionsRuleQuery) error

// EvalQuery return f(ctx, q).
func (f ActionsRuleQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.ActionsRuleQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.ActionsRuleQuery", q)
}

// The ActionsRuleMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type ActionsRuleMutationRuleFunc func(context.Context, *ent.ActionsRuleMutation) error

// EvalMutation calls f(ctx, m).
func (f ActionsRuleMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ActionsRuleMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ActionsRuleMutation", m)
}

// The ActivityQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type ActivityQueryRuleFunc func(context.Context, *ent.ActivityQuery) error

// EvalQuery return f(ctx, q).
func (f ActivityQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.ActivityQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.ActivityQuery", q)
}

// The ActivityMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type ActivityMutationRuleFunc func(context.Context, *ent.ActivityMutation) error

// EvalMutation calls f(ctx, m).
func (f ActivityMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ActivityMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ActivityMutation", m)
}

// The CheckListCategoryQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type CheckListCategoryQueryRuleFunc func(context.Context, *ent.CheckListCategoryQuery) error

// EvalQuery return f(ctx, q).
func (f CheckListCategoryQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.CheckListCategoryQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.CheckListCategoryQuery", q)
}

// The CheckListCategoryMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type CheckListCategoryMutationRuleFunc func(context.Context, *ent.CheckListCategoryMutation) error

// EvalMutation calls f(ctx, m).
func (f CheckListCategoryMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CheckListCategoryMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CheckListCategoryMutation", m)
}

// The CheckListCategoryDefinitionQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type CheckListCategoryDefinitionQueryRuleFunc func(context.Context, *ent.CheckListCategoryDefinitionQuery) error

// EvalQuery return f(ctx, q).
func (f CheckListCategoryDefinitionQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.CheckListCategoryDefinitionQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.CheckListCategoryDefinitionQuery", q)
}

// The CheckListCategoryDefinitionMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type CheckListCategoryDefinitionMutationRuleFunc func(context.Context, *ent.CheckListCategoryDefinitionMutation) error

// EvalMutation calls f(ctx, m).
func (f CheckListCategoryDefinitionMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CheckListCategoryDefinitionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CheckListCategoryDefinitionMutation", m)
}

// The CheckListItemQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type CheckListItemQueryRuleFunc func(context.Context, *ent.CheckListItemQuery) error

// EvalQuery return f(ctx, q).
func (f CheckListItemQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.CheckListItemQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.CheckListItemQuery", q)
}

// The CheckListItemMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type CheckListItemMutationRuleFunc func(context.Context, *ent.CheckListItemMutation) error

// EvalMutation calls f(ctx, m).
func (f CheckListItemMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CheckListItemMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CheckListItemMutation", m)
}

// The CheckListItemDefinitionQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type CheckListItemDefinitionQueryRuleFunc func(context.Context, *ent.CheckListItemDefinitionQuery) error

// EvalQuery return f(ctx, q).
func (f CheckListItemDefinitionQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.CheckListItemDefinitionQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.CheckListItemDefinitionQuery", q)
}

// The CheckListItemDefinitionMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type CheckListItemDefinitionMutationRuleFunc func(context.Context, *ent.CheckListItemDefinitionMutation) error

// EvalMutation calls f(ctx, m).
func (f CheckListItemDefinitionMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CheckListItemDefinitionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CheckListItemDefinitionMutation", m)
}

// The CommentQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type CommentQueryRuleFunc func(context.Context, *ent.CommentQuery) error

// EvalQuery return f(ctx, q).
func (f CommentQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.CommentQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.CommentQuery", q)
}

// The CommentMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type CommentMutationRuleFunc func(context.Context, *ent.CommentMutation) error

// EvalMutation calls f(ctx, m).
func (f CommentMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CommentMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CommentMutation", m)
}

// The CustomerQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type CustomerQueryRuleFunc func(context.Context, *ent.CustomerQuery) error

// EvalQuery return f(ctx, q).
func (f CustomerQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.CustomerQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.CustomerQuery", q)
}

// The CustomerMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type CustomerMutationRuleFunc func(context.Context, *ent.CustomerMutation) error

// EvalMutation calls f(ctx, m).
func (f CustomerMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CustomerMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CustomerMutation", m)
}

// The EquipmentQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type EquipmentQueryRuleFunc func(context.Context, *ent.EquipmentQuery) error

// EvalQuery return f(ctx, q).
func (f EquipmentQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.EquipmentQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.EquipmentQuery", q)
}

// The EquipmentMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type EquipmentMutationRuleFunc func(context.Context, *ent.EquipmentMutation) error

// EvalMutation calls f(ctx, m).
func (f EquipmentMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentMutation", m)
}

// The EquipmentCategoryQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type EquipmentCategoryQueryRuleFunc func(context.Context, *ent.EquipmentCategoryQuery) error

// EvalQuery return f(ctx, q).
func (f EquipmentCategoryQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.EquipmentCategoryQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.EquipmentCategoryQuery", q)
}

// The EquipmentCategoryMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type EquipmentCategoryMutationRuleFunc func(context.Context, *ent.EquipmentCategoryMutation) error

// EvalMutation calls f(ctx, m).
func (f EquipmentCategoryMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentCategoryMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentCategoryMutation", m)
}

// The EquipmentPortQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type EquipmentPortQueryRuleFunc func(context.Context, *ent.EquipmentPortQuery) error

// EvalQuery return f(ctx, q).
func (f EquipmentPortQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.EquipmentPortQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.EquipmentPortQuery", q)
}

// The EquipmentPortMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type EquipmentPortMutationRuleFunc func(context.Context, *ent.EquipmentPortMutation) error

// EvalMutation calls f(ctx, m).
func (f EquipmentPortMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentPortMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentPortMutation", m)
}

// The EquipmentPortDefinitionQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type EquipmentPortDefinitionQueryRuleFunc func(context.Context, *ent.EquipmentPortDefinitionQuery) error

// EvalQuery return f(ctx, q).
func (f EquipmentPortDefinitionQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.EquipmentPortDefinitionQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.EquipmentPortDefinitionQuery", q)
}

// The EquipmentPortDefinitionMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type EquipmentPortDefinitionMutationRuleFunc func(context.Context, *ent.EquipmentPortDefinitionMutation) error

// EvalMutation calls f(ctx, m).
func (f EquipmentPortDefinitionMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentPortDefinitionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentPortDefinitionMutation", m)
}

// The EquipmentPortTypeQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type EquipmentPortTypeQueryRuleFunc func(context.Context, *ent.EquipmentPortTypeQuery) error

// EvalQuery return f(ctx, q).
func (f EquipmentPortTypeQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.EquipmentPortTypeQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.EquipmentPortTypeQuery", q)
}

// The EquipmentPortTypeMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type EquipmentPortTypeMutationRuleFunc func(context.Context, *ent.EquipmentPortTypeMutation) error

// EvalMutation calls f(ctx, m).
func (f EquipmentPortTypeMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentPortTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentPortTypeMutation", m)
}

// The EquipmentPositionQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type EquipmentPositionQueryRuleFunc func(context.Context, *ent.EquipmentPositionQuery) error

// EvalQuery return f(ctx, q).
func (f EquipmentPositionQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.EquipmentPositionQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.EquipmentPositionQuery", q)
}

// The EquipmentPositionMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type EquipmentPositionMutationRuleFunc func(context.Context, *ent.EquipmentPositionMutation) error

// EvalMutation calls f(ctx, m).
func (f EquipmentPositionMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentPositionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentPositionMutation", m)
}

// The EquipmentPositionDefinitionQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type EquipmentPositionDefinitionQueryRuleFunc func(context.Context, *ent.EquipmentPositionDefinitionQuery) error

// EvalQuery return f(ctx, q).
func (f EquipmentPositionDefinitionQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.EquipmentPositionDefinitionQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.EquipmentPositionDefinitionQuery", q)
}

// The EquipmentPositionDefinitionMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type EquipmentPositionDefinitionMutationRuleFunc func(context.Context, *ent.EquipmentPositionDefinitionMutation) error

// EvalMutation calls f(ctx, m).
func (f EquipmentPositionDefinitionMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentPositionDefinitionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentPositionDefinitionMutation", m)
}

// The EquipmentTypeQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type EquipmentTypeQueryRuleFunc func(context.Context, *ent.EquipmentTypeQuery) error

// EvalQuery return f(ctx, q).
func (f EquipmentTypeQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.EquipmentTypeQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.EquipmentTypeQuery", q)
}

// The EquipmentTypeMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type EquipmentTypeMutationRuleFunc func(context.Context, *ent.EquipmentTypeMutation) error

// EvalMutation calls f(ctx, m).
func (f EquipmentTypeMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.EquipmentTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.EquipmentTypeMutation", m)
}

// The FileQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type FileQueryRuleFunc func(context.Context, *ent.FileQuery) error

// EvalQuery return f(ctx, q).
func (f FileQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.FileQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.FileQuery", q)
}

// The FileMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type FileMutationRuleFunc func(context.Context, *ent.FileMutation) error

// EvalMutation calls f(ctx, m).
func (f FileMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.FileMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.FileMutation", m)
}

// The FloorPlanQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type FloorPlanQueryRuleFunc func(context.Context, *ent.FloorPlanQuery) error

// EvalQuery return f(ctx, q).
func (f FloorPlanQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.FloorPlanQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.FloorPlanQuery", q)
}

// The FloorPlanMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type FloorPlanMutationRuleFunc func(context.Context, *ent.FloorPlanMutation) error

// EvalMutation calls f(ctx, m).
func (f FloorPlanMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.FloorPlanMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.FloorPlanMutation", m)
}

// The FloorPlanReferencePointQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type FloorPlanReferencePointQueryRuleFunc func(context.Context, *ent.FloorPlanReferencePointQuery) error

// EvalQuery return f(ctx, q).
func (f FloorPlanReferencePointQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.FloorPlanReferencePointQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.FloorPlanReferencePointQuery", q)
}

// The FloorPlanReferencePointMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type FloorPlanReferencePointMutationRuleFunc func(context.Context, *ent.FloorPlanReferencePointMutation) error

// EvalMutation calls f(ctx, m).
func (f FloorPlanReferencePointMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.FloorPlanReferencePointMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.FloorPlanReferencePointMutation", m)
}

// The FloorPlanScaleQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type FloorPlanScaleQueryRuleFunc func(context.Context, *ent.FloorPlanScaleQuery) error

// EvalQuery return f(ctx, q).
func (f FloorPlanScaleQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.FloorPlanScaleQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.FloorPlanScaleQuery", q)
}

// The FloorPlanScaleMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type FloorPlanScaleMutationRuleFunc func(context.Context, *ent.FloorPlanScaleMutation) error

// EvalMutation calls f(ctx, m).
func (f FloorPlanScaleMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.FloorPlanScaleMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.FloorPlanScaleMutation", m)
}

// The HyperlinkQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type HyperlinkQueryRuleFunc func(context.Context, *ent.HyperlinkQuery) error

// EvalQuery return f(ctx, q).
func (f HyperlinkQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.HyperlinkQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.HyperlinkQuery", q)
}

// The HyperlinkMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type HyperlinkMutationRuleFunc func(context.Context, *ent.HyperlinkMutation) error

// EvalMutation calls f(ctx, m).
func (f HyperlinkMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.HyperlinkMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.HyperlinkMutation", m)
}

// The LinkQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type LinkQueryRuleFunc func(context.Context, *ent.LinkQuery) error

// EvalQuery return f(ctx, q).
func (f LinkQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.LinkQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.LinkQuery", q)
}

// The LinkMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type LinkMutationRuleFunc func(context.Context, *ent.LinkMutation) error

// EvalMutation calls f(ctx, m).
func (f LinkMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.LinkMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.LinkMutation", m)
}

// The LocationQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type LocationQueryRuleFunc func(context.Context, *ent.LocationQuery) error

// EvalQuery return f(ctx, q).
func (f LocationQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.LocationQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.LocationQuery", q)
}

// The LocationMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type LocationMutationRuleFunc func(context.Context, *ent.LocationMutation) error

// EvalMutation calls f(ctx, m).
func (f LocationMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.LocationMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.LocationMutation", m)
}

// The LocationTypeQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type LocationTypeQueryRuleFunc func(context.Context, *ent.LocationTypeQuery) error

// EvalQuery return f(ctx, q).
func (f LocationTypeQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.LocationTypeQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.LocationTypeQuery", q)
}

// The LocationTypeMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type LocationTypeMutationRuleFunc func(context.Context, *ent.LocationTypeMutation) error

// EvalMutation calls f(ctx, m).
func (f LocationTypeMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.LocationTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.LocationTypeMutation", m)
}

// The PermissionsPolicyQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type PermissionsPolicyQueryRuleFunc func(context.Context, *ent.PermissionsPolicyQuery) error

// EvalQuery return f(ctx, q).
func (f PermissionsPolicyQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.PermissionsPolicyQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.PermissionsPolicyQuery", q)
}

// The PermissionsPolicyMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type PermissionsPolicyMutationRuleFunc func(context.Context, *ent.PermissionsPolicyMutation) error

// EvalMutation calls f(ctx, m).
func (f PermissionsPolicyMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.PermissionsPolicyMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.PermissionsPolicyMutation", m)
}

// The ProjectQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type ProjectQueryRuleFunc func(context.Context, *ent.ProjectQuery) error

// EvalQuery return f(ctx, q).
func (f ProjectQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.ProjectQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.ProjectQuery", q)
}

// The ProjectMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type ProjectMutationRuleFunc func(context.Context, *ent.ProjectMutation) error

// EvalMutation calls f(ctx, m).
func (f ProjectMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ProjectMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ProjectMutation", m)
}

// The ProjectTypeQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type ProjectTypeQueryRuleFunc func(context.Context, *ent.ProjectTypeQuery) error

// EvalQuery return f(ctx, q).
func (f ProjectTypeQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.ProjectTypeQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.ProjectTypeQuery", q)
}

// The ProjectTypeMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type ProjectTypeMutationRuleFunc func(context.Context, *ent.ProjectTypeMutation) error

// EvalMutation calls f(ctx, m).
func (f ProjectTypeMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ProjectTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ProjectTypeMutation", m)
}

// The PropertyQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type PropertyQueryRuleFunc func(context.Context, *ent.PropertyQuery) error

// EvalQuery return f(ctx, q).
func (f PropertyQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.PropertyQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.PropertyQuery", q)
}

// The PropertyMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type PropertyMutationRuleFunc func(context.Context, *ent.PropertyMutation) error

// EvalMutation calls f(ctx, m).
func (f PropertyMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.PropertyMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.PropertyMutation", m)
}

// The PropertyTypeQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type PropertyTypeQueryRuleFunc func(context.Context, *ent.PropertyTypeQuery) error

// EvalQuery return f(ctx, q).
func (f PropertyTypeQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.PropertyTypeQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.PropertyTypeQuery", q)
}

// The PropertyTypeMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type PropertyTypeMutationRuleFunc func(context.Context, *ent.PropertyTypeMutation) error

// EvalMutation calls f(ctx, m).
func (f PropertyTypeMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.PropertyTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.PropertyTypeMutation", m)
}

// The ReportFilterQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type ReportFilterQueryRuleFunc func(context.Context, *ent.ReportFilterQuery) error

// EvalQuery return f(ctx, q).
func (f ReportFilterQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.ReportFilterQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.ReportFilterQuery", q)
}

// The ReportFilterMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type ReportFilterMutationRuleFunc func(context.Context, *ent.ReportFilterMutation) error

// EvalMutation calls f(ctx, m).
func (f ReportFilterMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ReportFilterMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ReportFilterMutation", m)
}

// The ServiceQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type ServiceQueryRuleFunc func(context.Context, *ent.ServiceQuery) error

// EvalQuery return f(ctx, q).
func (f ServiceQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.ServiceQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.ServiceQuery", q)
}

// The ServiceMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type ServiceMutationRuleFunc func(context.Context, *ent.ServiceMutation) error

// EvalMutation calls f(ctx, m).
func (f ServiceMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ServiceMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ServiceMutation", m)
}

// The ServiceEndpointQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type ServiceEndpointQueryRuleFunc func(context.Context, *ent.ServiceEndpointQuery) error

// EvalQuery return f(ctx, q).
func (f ServiceEndpointQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.ServiceEndpointQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.ServiceEndpointQuery", q)
}

// The ServiceEndpointMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type ServiceEndpointMutationRuleFunc func(context.Context, *ent.ServiceEndpointMutation) error

// EvalMutation calls f(ctx, m).
func (f ServiceEndpointMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ServiceEndpointMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ServiceEndpointMutation", m)
}

// The ServiceEndpointDefinitionQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type ServiceEndpointDefinitionQueryRuleFunc func(context.Context, *ent.ServiceEndpointDefinitionQuery) error

// EvalQuery return f(ctx, q).
func (f ServiceEndpointDefinitionQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.ServiceEndpointDefinitionQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.ServiceEndpointDefinitionQuery", q)
}

// The ServiceEndpointDefinitionMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type ServiceEndpointDefinitionMutationRuleFunc func(context.Context, *ent.ServiceEndpointDefinitionMutation) error

// EvalMutation calls f(ctx, m).
func (f ServiceEndpointDefinitionMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ServiceEndpointDefinitionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ServiceEndpointDefinitionMutation", m)
}

// The ServiceTypeQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type ServiceTypeQueryRuleFunc func(context.Context, *ent.ServiceTypeQuery) error

// EvalQuery return f(ctx, q).
func (f ServiceTypeQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.ServiceTypeQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.ServiceTypeQuery", q)
}

// The ServiceTypeMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type ServiceTypeMutationRuleFunc func(context.Context, *ent.ServiceTypeMutation) error

// EvalMutation calls f(ctx, m).
func (f ServiceTypeMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.ServiceTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.ServiceTypeMutation", m)
}

// The SurveyQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type SurveyQueryRuleFunc func(context.Context, *ent.SurveyQuery) error

// EvalQuery return f(ctx, q).
func (f SurveyQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.SurveyQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.SurveyQuery", q)
}

// The SurveyMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type SurveyMutationRuleFunc func(context.Context, *ent.SurveyMutation) error

// EvalMutation calls f(ctx, m).
func (f SurveyMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyMutation", m)
}

// The SurveyCellScanQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type SurveyCellScanQueryRuleFunc func(context.Context, *ent.SurveyCellScanQuery) error

// EvalQuery return f(ctx, q).
func (f SurveyCellScanQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.SurveyCellScanQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.SurveyCellScanQuery", q)
}

// The SurveyCellScanMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type SurveyCellScanMutationRuleFunc func(context.Context, *ent.SurveyCellScanMutation) error

// EvalMutation calls f(ctx, m).
func (f SurveyCellScanMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyCellScanMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyCellScanMutation", m)
}

// The SurveyQuestionQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type SurveyQuestionQueryRuleFunc func(context.Context, *ent.SurveyQuestionQuery) error

// EvalQuery return f(ctx, q).
func (f SurveyQuestionQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.SurveyQuestionQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.SurveyQuestionQuery", q)
}

// The SurveyQuestionMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type SurveyQuestionMutationRuleFunc func(context.Context, *ent.SurveyQuestionMutation) error

// EvalMutation calls f(ctx, m).
func (f SurveyQuestionMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyQuestionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyQuestionMutation", m)
}

// The SurveyTemplateCategoryQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type SurveyTemplateCategoryQueryRuleFunc func(context.Context, *ent.SurveyTemplateCategoryQuery) error

// EvalQuery return f(ctx, q).
func (f SurveyTemplateCategoryQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.SurveyTemplateCategoryQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.SurveyTemplateCategoryQuery", q)
}

// The SurveyTemplateCategoryMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type SurveyTemplateCategoryMutationRuleFunc func(context.Context, *ent.SurveyTemplateCategoryMutation) error

// EvalMutation calls f(ctx, m).
func (f SurveyTemplateCategoryMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyTemplateCategoryMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyTemplateCategoryMutation", m)
}

// The SurveyTemplateQuestionQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type SurveyTemplateQuestionQueryRuleFunc func(context.Context, *ent.SurveyTemplateQuestionQuery) error

// EvalQuery return f(ctx, q).
func (f SurveyTemplateQuestionQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.SurveyTemplateQuestionQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.SurveyTemplateQuestionQuery", q)
}

// The SurveyTemplateQuestionMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type SurveyTemplateQuestionMutationRuleFunc func(context.Context, *ent.SurveyTemplateQuestionMutation) error

// EvalMutation calls f(ctx, m).
func (f SurveyTemplateQuestionMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyTemplateQuestionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyTemplateQuestionMutation", m)
}

// The SurveyWiFiScanQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type SurveyWiFiScanQueryRuleFunc func(context.Context, *ent.SurveyWiFiScanQuery) error

// EvalQuery return f(ctx, q).
func (f SurveyWiFiScanQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.SurveyWiFiScanQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.SurveyWiFiScanQuery", q)
}

// The SurveyWiFiScanMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type SurveyWiFiScanMutationRuleFunc func(context.Context, *ent.SurveyWiFiScanMutation) error

// EvalMutation calls f(ctx, m).
func (f SurveyWiFiScanMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SurveyWiFiScanMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SurveyWiFiScanMutation", m)
}

// The UserQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type UserQueryRuleFunc func(context.Context, *ent.UserQuery) error

// EvalQuery return f(ctx, q).
func (f UserQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.UserQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.UserQuery", q)
}

// The UserMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type UserMutationRuleFunc func(context.Context, *ent.UserMutation) error

// EvalMutation calls f(ctx, m).
func (f UserMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.UserMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.UserMutation", m)
}

// The UsersGroupQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type UsersGroupQueryRuleFunc func(context.Context, *ent.UsersGroupQuery) error

// EvalQuery return f(ctx, q).
func (f UsersGroupQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.UsersGroupQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.UsersGroupQuery", q)
}

// The UsersGroupMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type UsersGroupMutationRuleFunc func(context.Context, *ent.UsersGroupMutation) error

// EvalMutation calls f(ctx, m).
func (f UsersGroupMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.UsersGroupMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.UsersGroupMutation", m)
}

// The WorkOrderQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type WorkOrderQueryRuleFunc func(context.Context, *ent.WorkOrderQuery) error

// EvalQuery return f(ctx, q).
func (f WorkOrderQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.WorkOrderQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.WorkOrderQuery", q)
}

// The WorkOrderMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type WorkOrderMutationRuleFunc func(context.Context, *ent.WorkOrderMutation) error

// EvalMutation calls f(ctx, m).
func (f WorkOrderMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.WorkOrderMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.WorkOrderMutation", m)
}

// The WorkOrderDefinitionQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type WorkOrderDefinitionQueryRuleFunc func(context.Context, *ent.WorkOrderDefinitionQuery) error

// EvalQuery return f(ctx, q).
func (f WorkOrderDefinitionQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.WorkOrderDefinitionQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.WorkOrderDefinitionQuery", q)
}

// The WorkOrderDefinitionMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type WorkOrderDefinitionMutationRuleFunc func(context.Context, *ent.WorkOrderDefinitionMutation) error

// EvalMutation calls f(ctx, m).
func (f WorkOrderDefinitionMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.WorkOrderDefinitionMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.WorkOrderDefinitionMutation", m)
}

// The WorkOrderTypeQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type WorkOrderTypeQueryRuleFunc func(context.Context, *ent.WorkOrderTypeQuery) error

// EvalQuery return f(ctx, q).
func (f WorkOrderTypeQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.WorkOrderTypeQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.WorkOrderTypeQuery", q)
}

// The WorkOrderTypeMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type WorkOrderTypeMutationRuleFunc func(context.Context, *ent.WorkOrderTypeMutation) error

// EvalMutation calls f(ctx, m).
func (f WorkOrderTypeMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.WorkOrderTypeMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.WorkOrderTypeMutation", m)
}
