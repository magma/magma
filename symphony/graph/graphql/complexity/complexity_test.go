// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complexity_test

import (
	"reflect"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/graphql/complexity"
	"github.com/scylladb/go-set/strset"
	"github.com/stretchr/testify/assert"
)

func TestComplexityRoot(t *testing.T) {
	exempted := strset.New(
		"Equipment.Ports",
		"Location.DistanceKm",
		"LocationType.Locations",
	)
	root := reflect.ValueOf(complexity.New())
	trivial := reflect.TypeOf(func(int) (_ int) { return })
	for i, n := 0, root.NumField(); i < n; i++ {
		name := root.Type().Field(i).Name
		if name == "Query" || name == "Mutation" {
			continue
		}
		outer := root.Field(i)
		for i, n := 0, outer.NumField(); i < n; i++ {
			inner := outer.Field(i)
			if !inner.IsNil() {
				continue
			}
			name := name + "." + outer.Type().Field(i).Name
			if exempted.Has(name) {
				exempted.Remove(name)
				continue
			}
			assert.Truef(t, trivial.AssignableTo(inner.Type()),
				`field %q has non trivial complexity func `+
					"either define one or add it to exempted set", name,
			)
		}
	}
	assert.Truef(t, exempted.IsEmpty(),
		"exempted set contains unchecked items: %s", exempted,
	)
}
func TestQueryComplexity(t *testing.T) {
	exempted := strset.New(
		"ActionsRules", "ActionsTriggers", "LatestPythonPackage",
		"Me", "NearestSites", "Node", "PossibleProperties",
		"PythonPackages", "ReportFilters", "Surveys", "User", "Vertex",
	)
	query := reflect.ValueOf(complexity.New().Query)
	for i, n := 0, query.NumField(); i < n; i++ {
		name := query.Type().Field(i).Name
		if exempted.Has(name) {
			exempted.Remove(name)
			continue
		}
		assert.Falsef(t, query.Field(i).IsNil(),
			"query field %q has no complexity func, "+
				"either define one or add it to exempted set", name,
		)
	}
	assert.Truef(t, exempted.IsEmpty(),
		"exempted set contains unchecked items: %s", exempted,
	)
}

func TestPaginationComplexity(t *testing.T) {
	tests := []struct {
		name            string
		childComplexity int
		first, last     *int
		want            int
	}{
		{
			name:            "Forwards",
			childComplexity: 10,
			first:           pointer.ToInt(10),
			want:            100,
		},
		{
			name:            "Backwards",
			childComplexity: 5,
			last:            pointer.ToInt(100),
			want:            500,
		},
		{
			name:            "ForwardsBackwards",
			childComplexity: 2,
			first:           pointer.ToInt(4),
			last:            pointer.ToInt(3),
			want:            6,
		},
		{
			name: "NoLimit",
			want: complexity.Infinite,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := complexity.PaginationComplexity(
				tc.childComplexity, nil, tc.first, nil, tc.last,
			)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestSearchComplexity(t *testing.T) {
	tests := []struct {
		name            string
		childComplexity int
		limit           *int
		want            int
	}{
		{
			name:            "Limited",
			childComplexity: 10,
			limit:           pointer.ToInt(5),
			want:            50,
		},
		{
			name: "NoLimit",
			want: complexity.Infinite,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := complexity.SearchComplexity(tc.childComplexity, tc.limit)
			assert.Equal(t, tc.want, got)
		})
	}
}
