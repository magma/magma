// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent/predicate"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/comment"
	"github.com/facebookincubator/symphony/graph/ent/privacy"
)

// CommentReadPolicyRule grants read permission to comment based on policy.
func CommentReadPolicyRule() privacy.QueryRule {
	return privacy.CommentQueryRuleFunc(func(ctx context.Context, q *ent.CommentQuery) error {
		var predicates []predicate.Comment
		woPredicate := workOrderReadPredicate(ctx)
		if woPredicate != nil {
			predicates = append(predicates,
				comment.Or(
					comment.Not(comment.HasWorkOrder()),
					comment.HasWorkOrderWith(woPredicate)))
		}
		projectPredicate := projectReadPredicate(ctx)
		if projectPredicate != nil {
			predicates = append(predicates,
				comment.Or(
					comment.Not(comment.HasProject()),
					comment.HasProjectWith(projectPredicate)))
		}
		q.Where(predicates...)
		return privacy.Skip
	})
}

// CommentWritePolicyRule grants write permission to comment based on policy.
func CommentWritePolicyRule() privacy.MutationRule {
	return privacy.CommentMutationRuleFunc(func(ctx context.Context, m *ent.CommentMutation) error {
		commentID, exists := m.ID()
		if !exists {
			return privacy.Skip
		}
		comm, err := m.Client().Comment.Query().
			Where(comment.ID(commentID)).
			WithWorkOrder().
			WithProject().
			Only(ctx)

		if err != nil {
			if !ent.IsNotFound(err) {
				return privacy.Denyf("failed to fetch comment: %w", err)
			}
			return privacy.Skip
		}
		p := FromContext(ctx)
		switch {
		case comm.Edges.WorkOrder != nil:
			return allowOrSkipWorkOrder(ctx, p, comm.Edges.WorkOrder)
		case comm.Edges.Project != nil:
			projectTypeID, err := comm.Edges.Project.QueryType().OnlyID(ctx)
			if err != nil {
				return privacy.Denyf("failed to fetch project type id: %w", err)
			}
			return privacyDecision(checkWorkforce(p.WorkforcePolicy.Data.Update, nil, &projectTypeID))
		}
		return privacy.Skip
	})
}

// CommentCreatePolicyRule grants create permission to comment based on policy.
func CommentCreatePolicyRule() privacy.MutationRule {
	return privacy.CommentMutationRuleFunc(func(ctx context.Context, m *ent.CommentMutation) error {
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
		if projectID, exists := m.ProjectID(); exists {
			projectTypeID, err := m.Client().ProjectType.Query().
				Where(projecttype.HasProjectsWith(project.ID(projectID))).
				OnlyID(ctx)
			if err != nil {
				if ent.IsNotFound(err) {
					return privacy.Skip
				}
				return privacy.Denyf("failed to fetch project type id: %w", err)
			}
			return privacyDecision(checkWorkforce(p.WorkforcePolicy.Data.Update, nil, &projectTypeID))
		}
		return privacy.Skip
	})
}
