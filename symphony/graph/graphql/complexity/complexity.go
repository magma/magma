// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package complexity

import (
	"math"
	"math/bits"

	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
)

// Infinite is the maximum possible complexity value.
const Infinite = 1<<(bits.UintSize-1) - 1

// New creates a graphql complexity root.
func New() (complexity generated.ComplexityRoot) {
	complexity.Location.Topology = func(childComplexity int, depth int) int {
		return childComplexity * int(math.Pow10(depth)) / 2
	}
	complexity.Query.CustomerSearch = SearchComplexity
	complexity.Query.Customers = PaginationComplexity
	complexity.Query.EquipmentPortDefinitions = PaginationComplexity
	complexity.Query.EquipmentPortTypes = PaginationComplexity
	complexity.Query.EquipmentPorts = PaginationComplexity
	complexity.Query.EquipmentSearch = func(childComplexity int, _ []*models.EquipmentFilterInput, limit *int) int {
		return SearchComplexity(childComplexity, limit)
	}
	complexity.Query.EquipmentTypes = PaginationComplexity
	complexity.Query.Equipments = PaginationComplexity
	complexity.Query.LinkSearch = func(childComplexity int, _ []*models.LinkFilterInput, limit *int) int {
		return SearchComplexity(childComplexity, limit)
	}
	complexity.Query.Links = PaginationComplexity
	complexity.Query.LocationSearch = func(childComplexity int, _ []*models.LocationFilterInput, limit *int) int {
		return SearchComplexity(childComplexity, limit)
	}
	complexity.Query.LocationTypes = PaginationComplexity
	complexity.Query.Locations = func(childComplexity int, _ *bool, _ []int, _ *string, _ *bool, after *ent.Cursor, first *int, before *ent.Cursor, last *int) int {
		return PaginationComplexity(childComplexity, after, first, before, last)
	}
	complexity.Query.PermissionsPolicies = PaginationComplexity
	complexity.Query.PermissionsPolicySearch = func(childComplexity int, _ []*models.PermissionsPolicyFilterInput, limit *int) int {
		return SearchComplexity(childComplexity, limit)
	}
	complexity.Query.PortSearch = func(childComplexity int, _ []*models.PortFilterInput, limit *int) int {
		return SearchComplexity(childComplexity, limit)
	}
	complexity.Query.ProjectSearch = func(childComplexity int, _ []*models.ProjectFilterInput, limit *int) int {
		return SearchComplexity(childComplexity, limit)
	}
	complexity.Query.ProjectTypes = PaginationComplexity
	complexity.Query.Projects = PaginationComplexity
	complexity.Query.SearchForNode = func(childComplexity int, _ string, after *ent.Cursor, first *int, before *ent.Cursor, last *int) int {
		return PaginationComplexity(childComplexity, after, first, before, last)
	}
	complexity.Query.ServiceSearch = func(childComplexity int, _ []*models.ServiceFilterInput, limit *int) int {
		return SearchComplexity(childComplexity, limit)
	}
	complexity.Query.ServiceTypes = PaginationComplexity
	complexity.Query.UserSearch = func(childComplexity int, _ []*models.UserFilterInput, limit *int) int {
		return SearchComplexity(childComplexity, limit)
	}
	complexity.Query.Users = PaginationComplexity
	complexity.Query.UsersGroups = PaginationComplexity
	complexity.Query.UsersGroupSearch = func(childComplexity int, _ []*models.UsersGroupFilterInput, limit *int) int {
		return SearchComplexity(childComplexity, limit)
	}
	complexity.Query.WorkOrderSearch = func(childComplexity int, _ []*models.WorkOrderFilterInput, limit *int) int {
		return SearchComplexity(childComplexity, limit)
	}
	complexity.Query.WorkOrderTypes = PaginationComplexity
	complexity.Query.WorkOrders = func(childComplexity int, after *ent.Cursor, first *int, before *ent.Cursor, last *int, _ *bool) int {
		return PaginationComplexity(childComplexity, after, first, before, last)
	}

	return complexity
}

// SearchComplexity returns the complexity function of searching queries.
func SearchComplexity(childComplexity int, limit *int) int {
	if limit != nil {
		return *limit * childComplexity
	}
	return Infinite
}

// PaginationComplexity returns the complexity function of paginating queries.
func PaginationComplexity(childComplexity int, _ *ent.Cursor, first *int, _ *ent.Cursor, last *int) int {
	switch {
	case first != nil:
		if last == nil || *first < *last {
			return *first * childComplexity
		}
		fallthrough
	case last != nil:
		return *last * childComplexity
	default:
		return Infinite
	}
}
