package authz

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent/comment"

	"github.com/facebookincubator/symphony/graph/ent/predicate"

	"github.com/facebookincubator/symphony/graph/ent"
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
