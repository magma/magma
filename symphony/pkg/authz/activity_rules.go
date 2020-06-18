// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

//TODO(T67933416): Return these rules
/*
// ActivityReadPolicyRule grants read permission to activity based on policy.
func ActivityReadPolicyRule() privacy.QueryRule {
	return privacy.ActivityQueryRuleFunc(func(ctx context.Context, q *ent.ActivityQuery) error {
		var predicates []predicate.Activity
		woPredicate := workOrderReadPredicate(ctx)
		if woPredicate != nil {
			predicates = append(predicates,
				activity.Or(
					activity.Not(activity.HasWorkOrder()),
					activity.HasWorkOrderWith(woPredicate)))
		}

		q.Where(predicates...)
		return privacy.Skip
	})
}

// ActivityWritePolicyRule grants write permission to activity based on policy.
func ActivityWritePolicyRule() privacy.MutationRule {
	return privacy.ActivityMutationRuleFunc(func(ctx context.Context, m *ent.ActivityMutation) error {
		activityID, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		wo, err := m.Client().Activity.Query().
			Where(activity.ID(activityID)).
			QueryWorkOrder().
			Only(ctx)

		if err != nil {
			if !ent.IsNotFound(err) {
				return privacy.Denyf("failed to fetch work order: %w", err)
			}
			return privacy.Skip
		}
		p := FromContext(ctx)
		if wo != nil {
			return allowOrSkipWorkOrder(ctx, p, wo)
		}
		return privacy.Skip
	})
}

// ActivityCreatePolicyRule grants create permission to activity based on policy.
func ActivityCreatePolicyRule() privacy.MutationRule {
	return privacy.ActivityMutationRuleFunc(func(ctx context.Context, m *ent.ActivityMutation) error {
		if !m.Op().Is(ent.OpCreate) {
			return privacy.Skip
		}
		p := FromContext(ctx)

		if workOrderID, exists := m.WorkOrderID(); exists {
			workOrder, err := m.Client().WorkOrder.Get(ctx, workOrderID)
			if err != nil {
				if !ent.IsNotFound(err) {
					return privacy.Denyf("failed to fetch work order: %w", err)
				}
				return privacy.Skip
			}
			return allowOrSkipWorkOrder(ctx, p, workOrder)
		}
		return privacy.Skip
	})
}
*/
