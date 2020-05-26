package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent/activity"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

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
