/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package restrictor

import (
	"fmt"

	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/promql"
)

// QueryRestrictor provides functionality to add restrictor labels to a
// Prometheus query
type QueryRestrictor struct {
	restrictors map[string]string
}

// NewQueryRestrictor returns a new QueryRestrictor with the given labels
func NewQueryRestrictor(restrictors map[string]string) *QueryRestrictor {
	return &QueryRestrictor{
		restrictors: restrictors,
	}
}

// RestrictQuery appends a label selector to each metric in a given query so
// that only metrics with those labels are returned from the query.
func (q *QueryRestrictor) RestrictQuery(query string) (string, error) {
	if query == "" {
		return "", fmt.Errorf("empty query string")
	}

	promQuery, err := promql.ParseExpr(query)
	if err != nil {
		return "", fmt.Errorf("error parsing query: %v", err)
	}
	promql.Inspect(promQuery, q.addRestrictorLabels())
	return promQuery.String(), nil
}

func (q *QueryRestrictor) addRestrictorLabels() func(n promql.Node, path []promql.Node) error {
	return func(n promql.Node, path []promql.Node) error {
		if n == nil {
			return nil
		}
		for labelName, labelValue := range q.restrictors {
			injectedLabelMatcher, err := labels.NewMatcher(labels.MatchEqual, labelName, labelValue)
			if err != nil {
				return fmt.Errorf("error creating labelMatcher: %v", err)
			}
			switch n := n.(type) {
			case *promql.VectorSelector:
				n.LabelMatchers = appendOrReplaceMatcher(n.LabelMatchers, *injectedLabelMatcher)
			case *promql.MatrixSelector:
				n.LabelMatchers = appendOrReplaceMatcher(n.LabelMatchers, *injectedLabelMatcher)
			}
		}
		return nil
	}
}

func appendOrReplaceMatcher(matchers []*labels.Matcher, newMatcher labels.Matcher) []*labels.Matcher {
	if getMatcherIndex(matchers, newMatcher.Name) >= 0 {
		return replaceLabelValue(matchers, newMatcher.Name, newMatcher.Value)
	} else {
		return append(matchers, &newMatcher)
	}
}

func getMatcherIndex(matchers []*labels.Matcher, name string) int {
	for idx, match := range matchers {
		if match.Name == name {
			return idx
		}
	}
	return -1
}

func replaceLabelValue(matchers []*labels.Matcher, name, value string) []*labels.Matcher {
	idx := getMatcherIndex(matchers, name)
	if idx >= -1 {
		matchers[idx].Value = value
	}
	return matchers
}
