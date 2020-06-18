// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/facebookincubator/symphony/pkg/ent/actionsrule"
	"github.com/facebookincubator/symphony/pkg/ent/activity"
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategory"
	"github.com/facebookincubator/symphony/pkg/ent/checklistcategorydefinition"
	"github.com/facebookincubator/symphony/pkg/ent/checklistitem"
	"github.com/facebookincubator/symphony/pkg/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/comment"
	"github.com/facebookincubator/symphony/pkg/ent/customer"
	"github.com/facebookincubator/symphony/pkg/ent/equipment"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentport"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentposition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/pkg/ent/equipmenttype"
	"github.com/facebookincubator/symphony/pkg/ent/file"
	"github.com/facebookincubator/symphony/pkg/ent/floorplan"
	"github.com/facebookincubator/symphony/pkg/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/pkg/ent/floorplanscale"
	"github.com/facebookincubator/symphony/pkg/ent/hyperlink"
	"github.com/facebookincubator/symphony/pkg/ent/link"
	"github.com/facebookincubator/symphony/pkg/ent/location"
	"github.com/facebookincubator/symphony/pkg/ent/locationtype"
	"github.com/facebookincubator/symphony/pkg/ent/permissionspolicy"
	"github.com/facebookincubator/symphony/pkg/ent/project"
	"github.com/facebookincubator/symphony/pkg/ent/projecttype"
	"github.com/facebookincubator/symphony/pkg/ent/property"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/ent/reportfilter"
	"github.com/facebookincubator/symphony/pkg/ent/service"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/pkg/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/servicetype"
	"github.com/facebookincubator/symphony/pkg/ent/survey"
	"github.com/facebookincubator/symphony/pkg/ent/surveycellscan"
	"github.com/facebookincubator/symphony/pkg/ent/surveyquestion"
	"github.com/facebookincubator/symphony/pkg/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/pkg/ent/surveytemplatequestion"
	"github.com/facebookincubator/symphony/pkg/ent/surveywifiscan"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/ent/usersgroup"
	"github.com/facebookincubator/symphony/pkg/ent/workorder"
	"github.com/facebookincubator/symphony/pkg/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/pkg/ent/workordertype"
	"github.com/ugorji/go/codec"
)

// PageInfo of a connection type.
type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	HasPreviousPage bool    `json:"hasPreviousPage"`
	StartCursor     *Cursor `json:"startCursor"`
	EndCursor       *Cursor `json:"endCursor"`
}

// Cursor of an edge type.
type Cursor struct {
	ID int
}

// ErrInvalidPagination error is returned when paginating with invalid parameters.
var ErrInvalidPagination = errors.New("ent: invalid pagination parameters")

var quote = []byte(`"`)

// MarshalGQL implements graphql.Marshaler interface.
func (c Cursor) MarshalGQL(w io.Writer) {
	w.Write(quote)
	defer w.Write(quote)
	wc := base64.NewEncoder(base64.StdEncoding, w)
	defer wc.Close()
	_ = codec.NewEncoder(wc, &codec.MsgpackHandle{}).Encode(c)
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (c *Cursor) UnmarshalGQL(v interface{}) error {
	s, ok := v.(string)
	if !ok {
		return fmt.Errorf("%T is not a string", v)
	}
	if err := codec.NewDecoder(
		base64.NewDecoder(
			base64.StdEncoding,
			strings.NewReader(s),
		),
		&codec.MsgpackHandle{},
	).Decode(c); err != nil {
		return fmt.Errorf("decode cursor: %w", err)
	}
	return nil
}

// ActionsRuleEdge is the edge representation of ActionsRule.
type ActionsRuleEdge struct {
	Node   *ActionsRule `json:"node"`
	Cursor Cursor       `json:"cursor"`
}

// ActionsRuleConnection is the connection containing edges to ActionsRule.
type ActionsRuleConnection struct {
	Edges    []*ActionsRuleEdge `json:"edges"`
	PageInfo PageInfo           `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to ActionsRule.
func (ar *ActionsRuleQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*ActionsRuleConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &ActionsRuleConnection{
				Edges: []*ActionsRuleEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &ActionsRuleConnection{
				Edges: []*ActionsRuleEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		ar = ar.Where(actionsrule.IDGT(after.ID))
	}
	if before != nil {
		ar = ar.Where(actionsrule.IDLT(before.ID))
	}
	if first != nil {
		ar = ar.Order(Asc(actionsrule.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		ar = ar.Order(Desc(actionsrule.FieldID)).Limit(*last + 1)
	}
	ar = ar.collectConnectionFields(ctx)

	nodes, err := ar.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &ActionsRuleConnection{
			Edges: []*ActionsRuleEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn ActionsRuleConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*ActionsRuleEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &ActionsRuleEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (ar *ActionsRuleQuery) collectConnectionFields(ctx context.Context) *ActionsRuleQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		ar = ar.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return ar
}

// ActivityEdge is the edge representation of Activity.
type ActivityEdge struct {
	Node   *Activity `json:"node"`
	Cursor Cursor    `json:"cursor"`
}

// ActivityConnection is the connection containing edges to Activity.
type ActivityConnection struct {
	Edges    []*ActivityEdge `json:"edges"`
	PageInfo PageInfo        `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Activity.
func (a *ActivityQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*ActivityConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &ActivityConnection{
				Edges: []*ActivityEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &ActivityConnection{
				Edges: []*ActivityEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		a = a.Where(activity.IDGT(after.ID))
	}
	if before != nil {
		a = a.Where(activity.IDLT(before.ID))
	}
	if first != nil {
		a = a.Order(Asc(activity.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		a = a.Order(Desc(activity.FieldID)).Limit(*last + 1)
	}
	a = a.collectConnectionFields(ctx)

	nodes, err := a.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &ActivityConnection{
			Edges: []*ActivityEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn ActivityConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*ActivityEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &ActivityEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (a *ActivityQuery) collectConnectionFields(ctx context.Context) *ActivityQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		a = a.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return a
}

// CheckListCategoryEdge is the edge representation of CheckListCategory.
type CheckListCategoryEdge struct {
	Node   *CheckListCategory `json:"node"`
	Cursor Cursor             `json:"cursor"`
}

// CheckListCategoryConnection is the connection containing edges to CheckListCategory.
type CheckListCategoryConnection struct {
	Edges    []*CheckListCategoryEdge `json:"edges"`
	PageInfo PageInfo                 `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to CheckListCategory.
func (clc *CheckListCategoryQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*CheckListCategoryConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &CheckListCategoryConnection{
				Edges: []*CheckListCategoryEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &CheckListCategoryConnection{
				Edges: []*CheckListCategoryEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		clc = clc.Where(checklistcategory.IDGT(after.ID))
	}
	if before != nil {
		clc = clc.Where(checklistcategory.IDLT(before.ID))
	}
	if first != nil {
		clc = clc.Order(Asc(checklistcategory.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		clc = clc.Order(Desc(checklistcategory.FieldID)).Limit(*last + 1)
	}
	clc = clc.collectConnectionFields(ctx)

	nodes, err := clc.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &CheckListCategoryConnection{
			Edges: []*CheckListCategoryEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn CheckListCategoryConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*CheckListCategoryEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &CheckListCategoryEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (clc *CheckListCategoryQuery) collectConnectionFields(ctx context.Context) *CheckListCategoryQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		clc = clc.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return clc
}

// CheckListCategoryDefinitionEdge is the edge representation of CheckListCategoryDefinition.
type CheckListCategoryDefinitionEdge struct {
	Node   *CheckListCategoryDefinition `json:"node"`
	Cursor Cursor                       `json:"cursor"`
}

// CheckListCategoryDefinitionConnection is the connection containing edges to CheckListCategoryDefinition.
type CheckListCategoryDefinitionConnection struct {
	Edges    []*CheckListCategoryDefinitionEdge `json:"edges"`
	PageInfo PageInfo                           `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to CheckListCategoryDefinition.
func (clcd *CheckListCategoryDefinitionQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*CheckListCategoryDefinitionConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &CheckListCategoryDefinitionConnection{
				Edges: []*CheckListCategoryDefinitionEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &CheckListCategoryDefinitionConnection{
				Edges: []*CheckListCategoryDefinitionEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		clcd = clcd.Where(checklistcategorydefinition.IDGT(after.ID))
	}
	if before != nil {
		clcd = clcd.Where(checklistcategorydefinition.IDLT(before.ID))
	}
	if first != nil {
		clcd = clcd.Order(Asc(checklistcategorydefinition.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		clcd = clcd.Order(Desc(checklistcategorydefinition.FieldID)).Limit(*last + 1)
	}
	clcd = clcd.collectConnectionFields(ctx)

	nodes, err := clcd.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &CheckListCategoryDefinitionConnection{
			Edges: []*CheckListCategoryDefinitionEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn CheckListCategoryDefinitionConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*CheckListCategoryDefinitionEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &CheckListCategoryDefinitionEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (clcd *CheckListCategoryDefinitionQuery) collectConnectionFields(ctx context.Context) *CheckListCategoryDefinitionQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		clcd = clcd.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return clcd
}

// CheckListItemEdge is the edge representation of CheckListItem.
type CheckListItemEdge struct {
	Node   *CheckListItem `json:"node"`
	Cursor Cursor         `json:"cursor"`
}

// CheckListItemConnection is the connection containing edges to CheckListItem.
type CheckListItemConnection struct {
	Edges    []*CheckListItemEdge `json:"edges"`
	PageInfo PageInfo             `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to CheckListItem.
func (cli *CheckListItemQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*CheckListItemConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &CheckListItemConnection{
				Edges: []*CheckListItemEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &CheckListItemConnection{
				Edges: []*CheckListItemEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		cli = cli.Where(checklistitem.IDGT(after.ID))
	}
	if before != nil {
		cli = cli.Where(checklistitem.IDLT(before.ID))
	}
	if first != nil {
		cli = cli.Order(Asc(checklistitem.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		cli = cli.Order(Desc(checklistitem.FieldID)).Limit(*last + 1)
	}
	cli = cli.collectConnectionFields(ctx)

	nodes, err := cli.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &CheckListItemConnection{
			Edges: []*CheckListItemEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn CheckListItemConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*CheckListItemEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &CheckListItemEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (cli *CheckListItemQuery) collectConnectionFields(ctx context.Context) *CheckListItemQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		cli = cli.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return cli
}

// CheckListItemDefinitionEdge is the edge representation of CheckListItemDefinition.
type CheckListItemDefinitionEdge struct {
	Node   *CheckListItemDefinition `json:"node"`
	Cursor Cursor                   `json:"cursor"`
}

// CheckListItemDefinitionConnection is the connection containing edges to CheckListItemDefinition.
type CheckListItemDefinitionConnection struct {
	Edges    []*CheckListItemDefinitionEdge `json:"edges"`
	PageInfo PageInfo                       `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to CheckListItemDefinition.
func (clid *CheckListItemDefinitionQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*CheckListItemDefinitionConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &CheckListItemDefinitionConnection{
				Edges: []*CheckListItemDefinitionEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &CheckListItemDefinitionConnection{
				Edges: []*CheckListItemDefinitionEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		clid = clid.Where(checklistitemdefinition.IDGT(after.ID))
	}
	if before != nil {
		clid = clid.Where(checklistitemdefinition.IDLT(before.ID))
	}
	if first != nil {
		clid = clid.Order(Asc(checklistitemdefinition.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		clid = clid.Order(Desc(checklistitemdefinition.FieldID)).Limit(*last + 1)
	}
	clid = clid.collectConnectionFields(ctx)

	nodes, err := clid.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &CheckListItemDefinitionConnection{
			Edges: []*CheckListItemDefinitionEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn CheckListItemDefinitionConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*CheckListItemDefinitionEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &CheckListItemDefinitionEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (clid *CheckListItemDefinitionQuery) collectConnectionFields(ctx context.Context) *CheckListItemDefinitionQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		clid = clid.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return clid
}

// CommentEdge is the edge representation of Comment.
type CommentEdge struct {
	Node   *Comment `json:"node"`
	Cursor Cursor   `json:"cursor"`
}

// CommentConnection is the connection containing edges to Comment.
type CommentConnection struct {
	Edges    []*CommentEdge `json:"edges"`
	PageInfo PageInfo       `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Comment.
func (c *CommentQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*CommentConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &CommentConnection{
				Edges: []*CommentEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &CommentConnection{
				Edges: []*CommentEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		c = c.Where(comment.IDGT(after.ID))
	}
	if before != nil {
		c = c.Where(comment.IDLT(before.ID))
	}
	if first != nil {
		c = c.Order(Asc(comment.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		c = c.Order(Desc(comment.FieldID)).Limit(*last + 1)
	}
	c = c.collectConnectionFields(ctx)

	nodes, err := c.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &CommentConnection{
			Edges: []*CommentEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn CommentConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*CommentEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &CommentEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (c *CommentQuery) collectConnectionFields(ctx context.Context) *CommentQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		c = c.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return c
}

// CustomerEdge is the edge representation of Customer.
type CustomerEdge struct {
	Node   *Customer `json:"node"`
	Cursor Cursor    `json:"cursor"`
}

// CustomerConnection is the connection containing edges to Customer.
type CustomerConnection struct {
	Edges    []*CustomerEdge `json:"edges"`
	PageInfo PageInfo        `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Customer.
func (c *CustomerQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*CustomerConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &CustomerConnection{
				Edges: []*CustomerEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &CustomerConnection{
				Edges: []*CustomerEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		c = c.Where(customer.IDGT(after.ID))
	}
	if before != nil {
		c = c.Where(customer.IDLT(before.ID))
	}
	if first != nil {
		c = c.Order(Asc(customer.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		c = c.Order(Desc(customer.FieldID)).Limit(*last + 1)
	}
	c = c.collectConnectionFields(ctx)

	nodes, err := c.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &CustomerConnection{
			Edges: []*CustomerEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn CustomerConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*CustomerEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &CustomerEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (c *CustomerQuery) collectConnectionFields(ctx context.Context) *CustomerQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		c = c.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return c
}

// EquipmentEdge is the edge representation of Equipment.
type EquipmentEdge struct {
	Node   *Equipment `json:"node"`
	Cursor Cursor     `json:"cursor"`
}

// EquipmentConnection is the connection containing edges to Equipment.
type EquipmentConnection struct {
	Edges    []*EquipmentEdge `json:"edges"`
	PageInfo PageInfo         `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Equipment.
func (e *EquipmentQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*EquipmentConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &EquipmentConnection{
				Edges: []*EquipmentEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &EquipmentConnection{
				Edges: []*EquipmentEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		e = e.Where(equipment.IDGT(after.ID))
	}
	if before != nil {
		e = e.Where(equipment.IDLT(before.ID))
	}
	if first != nil {
		e = e.Order(Asc(equipment.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		e = e.Order(Desc(equipment.FieldID)).Limit(*last + 1)
	}
	e = e.collectConnectionFields(ctx)

	nodes, err := e.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &EquipmentConnection{
			Edges: []*EquipmentEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn EquipmentConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*EquipmentEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &EquipmentEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (e *EquipmentQuery) collectConnectionFields(ctx context.Context) *EquipmentQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		e = e.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return e
}

// EquipmentCategoryEdge is the edge representation of EquipmentCategory.
type EquipmentCategoryEdge struct {
	Node   *EquipmentCategory `json:"node"`
	Cursor Cursor             `json:"cursor"`
}

// EquipmentCategoryConnection is the connection containing edges to EquipmentCategory.
type EquipmentCategoryConnection struct {
	Edges    []*EquipmentCategoryEdge `json:"edges"`
	PageInfo PageInfo                 `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to EquipmentCategory.
func (ec *EquipmentCategoryQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*EquipmentCategoryConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &EquipmentCategoryConnection{
				Edges: []*EquipmentCategoryEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &EquipmentCategoryConnection{
				Edges: []*EquipmentCategoryEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		ec = ec.Where(equipmentcategory.IDGT(after.ID))
	}
	if before != nil {
		ec = ec.Where(equipmentcategory.IDLT(before.ID))
	}
	if first != nil {
		ec = ec.Order(Asc(equipmentcategory.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		ec = ec.Order(Desc(equipmentcategory.FieldID)).Limit(*last + 1)
	}
	ec = ec.collectConnectionFields(ctx)

	nodes, err := ec.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &EquipmentCategoryConnection{
			Edges: []*EquipmentCategoryEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn EquipmentCategoryConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*EquipmentCategoryEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &EquipmentCategoryEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (ec *EquipmentCategoryQuery) collectConnectionFields(ctx context.Context) *EquipmentCategoryQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		ec = ec.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return ec
}

// EquipmentPortEdge is the edge representation of EquipmentPort.
type EquipmentPortEdge struct {
	Node   *EquipmentPort `json:"node"`
	Cursor Cursor         `json:"cursor"`
}

// EquipmentPortConnection is the connection containing edges to EquipmentPort.
type EquipmentPortConnection struct {
	Edges    []*EquipmentPortEdge `json:"edges"`
	PageInfo PageInfo             `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to EquipmentPort.
func (ep *EquipmentPortQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*EquipmentPortConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &EquipmentPortConnection{
				Edges: []*EquipmentPortEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &EquipmentPortConnection{
				Edges: []*EquipmentPortEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		ep = ep.Where(equipmentport.IDGT(after.ID))
	}
	if before != nil {
		ep = ep.Where(equipmentport.IDLT(before.ID))
	}
	if first != nil {
		ep = ep.Order(Asc(equipmentport.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		ep = ep.Order(Desc(equipmentport.FieldID)).Limit(*last + 1)
	}
	ep = ep.collectConnectionFields(ctx)

	nodes, err := ep.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &EquipmentPortConnection{
			Edges: []*EquipmentPortEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn EquipmentPortConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*EquipmentPortEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &EquipmentPortEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (ep *EquipmentPortQuery) collectConnectionFields(ctx context.Context) *EquipmentPortQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		ep = ep.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return ep
}

// EquipmentPortDefinitionEdge is the edge representation of EquipmentPortDefinition.
type EquipmentPortDefinitionEdge struct {
	Node   *EquipmentPortDefinition `json:"node"`
	Cursor Cursor                   `json:"cursor"`
}

// EquipmentPortDefinitionConnection is the connection containing edges to EquipmentPortDefinition.
type EquipmentPortDefinitionConnection struct {
	Edges    []*EquipmentPortDefinitionEdge `json:"edges"`
	PageInfo PageInfo                       `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to EquipmentPortDefinition.
func (epd *EquipmentPortDefinitionQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*EquipmentPortDefinitionConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &EquipmentPortDefinitionConnection{
				Edges: []*EquipmentPortDefinitionEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &EquipmentPortDefinitionConnection{
				Edges: []*EquipmentPortDefinitionEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		epd = epd.Where(equipmentportdefinition.IDGT(after.ID))
	}
	if before != nil {
		epd = epd.Where(equipmentportdefinition.IDLT(before.ID))
	}
	if first != nil {
		epd = epd.Order(Asc(equipmentportdefinition.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		epd = epd.Order(Desc(equipmentportdefinition.FieldID)).Limit(*last + 1)
	}
	epd = epd.collectConnectionFields(ctx)

	nodes, err := epd.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &EquipmentPortDefinitionConnection{
			Edges: []*EquipmentPortDefinitionEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn EquipmentPortDefinitionConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*EquipmentPortDefinitionEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &EquipmentPortDefinitionEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (epd *EquipmentPortDefinitionQuery) collectConnectionFields(ctx context.Context) *EquipmentPortDefinitionQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		epd = epd.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return epd
}

// EquipmentPortTypeEdge is the edge representation of EquipmentPortType.
type EquipmentPortTypeEdge struct {
	Node   *EquipmentPortType `json:"node"`
	Cursor Cursor             `json:"cursor"`
}

// EquipmentPortTypeConnection is the connection containing edges to EquipmentPortType.
type EquipmentPortTypeConnection struct {
	Edges    []*EquipmentPortTypeEdge `json:"edges"`
	PageInfo PageInfo                 `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to EquipmentPortType.
func (ept *EquipmentPortTypeQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*EquipmentPortTypeConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &EquipmentPortTypeConnection{
				Edges: []*EquipmentPortTypeEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &EquipmentPortTypeConnection{
				Edges: []*EquipmentPortTypeEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		ept = ept.Where(equipmentporttype.IDGT(after.ID))
	}
	if before != nil {
		ept = ept.Where(equipmentporttype.IDLT(before.ID))
	}
	if first != nil {
		ept = ept.Order(Asc(equipmentporttype.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		ept = ept.Order(Desc(equipmentporttype.FieldID)).Limit(*last + 1)
	}
	ept = ept.collectConnectionFields(ctx)

	nodes, err := ept.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &EquipmentPortTypeConnection{
			Edges: []*EquipmentPortTypeEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn EquipmentPortTypeConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*EquipmentPortTypeEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &EquipmentPortTypeEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (ept *EquipmentPortTypeQuery) collectConnectionFields(ctx context.Context) *EquipmentPortTypeQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		ept = ept.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return ept
}

// EquipmentPositionEdge is the edge representation of EquipmentPosition.
type EquipmentPositionEdge struct {
	Node   *EquipmentPosition `json:"node"`
	Cursor Cursor             `json:"cursor"`
}

// EquipmentPositionConnection is the connection containing edges to EquipmentPosition.
type EquipmentPositionConnection struct {
	Edges    []*EquipmentPositionEdge `json:"edges"`
	PageInfo PageInfo                 `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to EquipmentPosition.
func (ep *EquipmentPositionQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*EquipmentPositionConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &EquipmentPositionConnection{
				Edges: []*EquipmentPositionEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &EquipmentPositionConnection{
				Edges: []*EquipmentPositionEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		ep = ep.Where(equipmentposition.IDGT(after.ID))
	}
	if before != nil {
		ep = ep.Where(equipmentposition.IDLT(before.ID))
	}
	if first != nil {
		ep = ep.Order(Asc(equipmentposition.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		ep = ep.Order(Desc(equipmentposition.FieldID)).Limit(*last + 1)
	}
	ep = ep.collectConnectionFields(ctx)

	nodes, err := ep.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &EquipmentPositionConnection{
			Edges: []*EquipmentPositionEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn EquipmentPositionConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*EquipmentPositionEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &EquipmentPositionEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (ep *EquipmentPositionQuery) collectConnectionFields(ctx context.Context) *EquipmentPositionQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		ep = ep.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return ep
}

// EquipmentPositionDefinitionEdge is the edge representation of EquipmentPositionDefinition.
type EquipmentPositionDefinitionEdge struct {
	Node   *EquipmentPositionDefinition `json:"node"`
	Cursor Cursor                       `json:"cursor"`
}

// EquipmentPositionDefinitionConnection is the connection containing edges to EquipmentPositionDefinition.
type EquipmentPositionDefinitionConnection struct {
	Edges    []*EquipmentPositionDefinitionEdge `json:"edges"`
	PageInfo PageInfo                           `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to EquipmentPositionDefinition.
func (epd *EquipmentPositionDefinitionQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*EquipmentPositionDefinitionConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &EquipmentPositionDefinitionConnection{
				Edges: []*EquipmentPositionDefinitionEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &EquipmentPositionDefinitionConnection{
				Edges: []*EquipmentPositionDefinitionEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		epd = epd.Where(equipmentpositiondefinition.IDGT(after.ID))
	}
	if before != nil {
		epd = epd.Where(equipmentpositiondefinition.IDLT(before.ID))
	}
	if first != nil {
		epd = epd.Order(Asc(equipmentpositiondefinition.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		epd = epd.Order(Desc(equipmentpositiondefinition.FieldID)).Limit(*last + 1)
	}
	epd = epd.collectConnectionFields(ctx)

	nodes, err := epd.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &EquipmentPositionDefinitionConnection{
			Edges: []*EquipmentPositionDefinitionEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn EquipmentPositionDefinitionConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*EquipmentPositionDefinitionEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &EquipmentPositionDefinitionEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (epd *EquipmentPositionDefinitionQuery) collectConnectionFields(ctx context.Context) *EquipmentPositionDefinitionQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		epd = epd.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return epd
}

// EquipmentTypeEdge is the edge representation of EquipmentType.
type EquipmentTypeEdge struct {
	Node   *EquipmentType `json:"node"`
	Cursor Cursor         `json:"cursor"`
}

// EquipmentTypeConnection is the connection containing edges to EquipmentType.
type EquipmentTypeConnection struct {
	Edges    []*EquipmentTypeEdge `json:"edges"`
	PageInfo PageInfo             `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to EquipmentType.
func (et *EquipmentTypeQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*EquipmentTypeConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &EquipmentTypeConnection{
				Edges: []*EquipmentTypeEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &EquipmentTypeConnection{
				Edges: []*EquipmentTypeEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		et = et.Where(equipmenttype.IDGT(after.ID))
	}
	if before != nil {
		et = et.Where(equipmenttype.IDLT(before.ID))
	}
	if first != nil {
		et = et.Order(Asc(equipmenttype.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		et = et.Order(Desc(equipmenttype.FieldID)).Limit(*last + 1)
	}
	et = et.collectConnectionFields(ctx)

	nodes, err := et.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &EquipmentTypeConnection{
			Edges: []*EquipmentTypeEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn EquipmentTypeConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*EquipmentTypeEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &EquipmentTypeEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (et *EquipmentTypeQuery) collectConnectionFields(ctx context.Context) *EquipmentTypeQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		et = et.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return et
}

// FileEdge is the edge representation of File.
type FileEdge struct {
	Node   *File  `json:"node"`
	Cursor Cursor `json:"cursor"`
}

// FileConnection is the connection containing edges to File.
type FileConnection struct {
	Edges    []*FileEdge `json:"edges"`
	PageInfo PageInfo    `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to File.
func (f *FileQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*FileConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &FileConnection{
				Edges: []*FileEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &FileConnection{
				Edges: []*FileEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		f = f.Where(file.IDGT(after.ID))
	}
	if before != nil {
		f = f.Where(file.IDLT(before.ID))
	}
	if first != nil {
		f = f.Order(Asc(file.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		f = f.Order(Desc(file.FieldID)).Limit(*last + 1)
	}
	f = f.collectConnectionFields(ctx)

	nodes, err := f.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &FileConnection{
			Edges: []*FileEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn FileConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*FileEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &FileEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (f *FileQuery) collectConnectionFields(ctx context.Context) *FileQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		f = f.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return f
}

// FloorPlanEdge is the edge representation of FloorPlan.
type FloorPlanEdge struct {
	Node   *FloorPlan `json:"node"`
	Cursor Cursor     `json:"cursor"`
}

// FloorPlanConnection is the connection containing edges to FloorPlan.
type FloorPlanConnection struct {
	Edges    []*FloorPlanEdge `json:"edges"`
	PageInfo PageInfo         `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to FloorPlan.
func (fp *FloorPlanQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*FloorPlanConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &FloorPlanConnection{
				Edges: []*FloorPlanEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &FloorPlanConnection{
				Edges: []*FloorPlanEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		fp = fp.Where(floorplan.IDGT(after.ID))
	}
	if before != nil {
		fp = fp.Where(floorplan.IDLT(before.ID))
	}
	if first != nil {
		fp = fp.Order(Asc(floorplan.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		fp = fp.Order(Desc(floorplan.FieldID)).Limit(*last + 1)
	}
	fp = fp.collectConnectionFields(ctx)

	nodes, err := fp.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &FloorPlanConnection{
			Edges: []*FloorPlanEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn FloorPlanConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*FloorPlanEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &FloorPlanEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (fp *FloorPlanQuery) collectConnectionFields(ctx context.Context) *FloorPlanQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		fp = fp.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return fp
}

// FloorPlanReferencePointEdge is the edge representation of FloorPlanReferencePoint.
type FloorPlanReferencePointEdge struct {
	Node   *FloorPlanReferencePoint `json:"node"`
	Cursor Cursor                   `json:"cursor"`
}

// FloorPlanReferencePointConnection is the connection containing edges to FloorPlanReferencePoint.
type FloorPlanReferencePointConnection struct {
	Edges    []*FloorPlanReferencePointEdge `json:"edges"`
	PageInfo PageInfo                       `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to FloorPlanReferencePoint.
func (fprp *FloorPlanReferencePointQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*FloorPlanReferencePointConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &FloorPlanReferencePointConnection{
				Edges: []*FloorPlanReferencePointEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &FloorPlanReferencePointConnection{
				Edges: []*FloorPlanReferencePointEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		fprp = fprp.Where(floorplanreferencepoint.IDGT(after.ID))
	}
	if before != nil {
		fprp = fprp.Where(floorplanreferencepoint.IDLT(before.ID))
	}
	if first != nil {
		fprp = fprp.Order(Asc(floorplanreferencepoint.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		fprp = fprp.Order(Desc(floorplanreferencepoint.FieldID)).Limit(*last + 1)
	}
	fprp = fprp.collectConnectionFields(ctx)

	nodes, err := fprp.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &FloorPlanReferencePointConnection{
			Edges: []*FloorPlanReferencePointEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn FloorPlanReferencePointConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*FloorPlanReferencePointEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &FloorPlanReferencePointEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (fprp *FloorPlanReferencePointQuery) collectConnectionFields(ctx context.Context) *FloorPlanReferencePointQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		fprp = fprp.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return fprp
}

// FloorPlanScaleEdge is the edge representation of FloorPlanScale.
type FloorPlanScaleEdge struct {
	Node   *FloorPlanScale `json:"node"`
	Cursor Cursor          `json:"cursor"`
}

// FloorPlanScaleConnection is the connection containing edges to FloorPlanScale.
type FloorPlanScaleConnection struct {
	Edges    []*FloorPlanScaleEdge `json:"edges"`
	PageInfo PageInfo              `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to FloorPlanScale.
func (fps *FloorPlanScaleQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*FloorPlanScaleConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &FloorPlanScaleConnection{
				Edges: []*FloorPlanScaleEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &FloorPlanScaleConnection{
				Edges: []*FloorPlanScaleEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		fps = fps.Where(floorplanscale.IDGT(after.ID))
	}
	if before != nil {
		fps = fps.Where(floorplanscale.IDLT(before.ID))
	}
	if first != nil {
		fps = fps.Order(Asc(floorplanscale.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		fps = fps.Order(Desc(floorplanscale.FieldID)).Limit(*last + 1)
	}
	fps = fps.collectConnectionFields(ctx)

	nodes, err := fps.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &FloorPlanScaleConnection{
			Edges: []*FloorPlanScaleEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn FloorPlanScaleConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*FloorPlanScaleEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &FloorPlanScaleEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (fps *FloorPlanScaleQuery) collectConnectionFields(ctx context.Context) *FloorPlanScaleQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		fps = fps.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return fps
}

// HyperlinkEdge is the edge representation of Hyperlink.
type HyperlinkEdge struct {
	Node   *Hyperlink `json:"node"`
	Cursor Cursor     `json:"cursor"`
}

// HyperlinkConnection is the connection containing edges to Hyperlink.
type HyperlinkConnection struct {
	Edges    []*HyperlinkEdge `json:"edges"`
	PageInfo PageInfo         `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Hyperlink.
func (h *HyperlinkQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*HyperlinkConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &HyperlinkConnection{
				Edges: []*HyperlinkEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &HyperlinkConnection{
				Edges: []*HyperlinkEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		h = h.Where(hyperlink.IDGT(after.ID))
	}
	if before != nil {
		h = h.Where(hyperlink.IDLT(before.ID))
	}
	if first != nil {
		h = h.Order(Asc(hyperlink.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		h = h.Order(Desc(hyperlink.FieldID)).Limit(*last + 1)
	}
	h = h.collectConnectionFields(ctx)

	nodes, err := h.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &HyperlinkConnection{
			Edges: []*HyperlinkEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn HyperlinkConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*HyperlinkEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &HyperlinkEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (h *HyperlinkQuery) collectConnectionFields(ctx context.Context) *HyperlinkQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		h = h.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return h
}

// LinkEdge is the edge representation of Link.
type LinkEdge struct {
	Node   *Link  `json:"node"`
	Cursor Cursor `json:"cursor"`
}

// LinkConnection is the connection containing edges to Link.
type LinkConnection struct {
	Edges    []*LinkEdge `json:"edges"`
	PageInfo PageInfo    `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Link.
func (l *LinkQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*LinkConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &LinkConnection{
				Edges: []*LinkEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &LinkConnection{
				Edges: []*LinkEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		l = l.Where(link.IDGT(after.ID))
	}
	if before != nil {
		l = l.Where(link.IDLT(before.ID))
	}
	if first != nil {
		l = l.Order(Asc(link.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		l = l.Order(Desc(link.FieldID)).Limit(*last + 1)
	}
	l = l.collectConnectionFields(ctx)

	nodes, err := l.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &LinkConnection{
			Edges: []*LinkEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn LinkConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*LinkEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &LinkEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (l *LinkQuery) collectConnectionFields(ctx context.Context) *LinkQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		l = l.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return l
}

// LocationEdge is the edge representation of Location.
type LocationEdge struct {
	Node   *Location `json:"node"`
	Cursor Cursor    `json:"cursor"`
}

// LocationConnection is the connection containing edges to Location.
type LocationConnection struct {
	Edges    []*LocationEdge `json:"edges"`
	PageInfo PageInfo        `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Location.
func (l *LocationQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*LocationConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &LocationConnection{
				Edges: []*LocationEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &LocationConnection{
				Edges: []*LocationEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		l = l.Where(location.IDGT(after.ID))
	}
	if before != nil {
		l = l.Where(location.IDLT(before.ID))
	}
	if first != nil {
		l = l.Order(Asc(location.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		l = l.Order(Desc(location.FieldID)).Limit(*last + 1)
	}
	l = l.collectConnectionFields(ctx)

	nodes, err := l.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &LocationConnection{
			Edges: []*LocationEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn LocationConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*LocationEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &LocationEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (l *LocationQuery) collectConnectionFields(ctx context.Context) *LocationQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		l = l.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return l
}

// LocationTypeEdge is the edge representation of LocationType.
type LocationTypeEdge struct {
	Node   *LocationType `json:"node"`
	Cursor Cursor        `json:"cursor"`
}

// LocationTypeConnection is the connection containing edges to LocationType.
type LocationTypeConnection struct {
	Edges    []*LocationTypeEdge `json:"edges"`
	PageInfo PageInfo            `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to LocationType.
func (lt *LocationTypeQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*LocationTypeConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &LocationTypeConnection{
				Edges: []*LocationTypeEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &LocationTypeConnection{
				Edges: []*LocationTypeEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		lt = lt.Where(locationtype.IDGT(after.ID))
	}
	if before != nil {
		lt = lt.Where(locationtype.IDLT(before.ID))
	}
	if first != nil {
		lt = lt.Order(Asc(locationtype.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		lt = lt.Order(Desc(locationtype.FieldID)).Limit(*last + 1)
	}
	lt = lt.collectConnectionFields(ctx)

	nodes, err := lt.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &LocationTypeConnection{
			Edges: []*LocationTypeEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn LocationTypeConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*LocationTypeEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &LocationTypeEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (lt *LocationTypeQuery) collectConnectionFields(ctx context.Context) *LocationTypeQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		lt = lt.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return lt
}

// PermissionsPolicyEdge is the edge representation of PermissionsPolicy.
type PermissionsPolicyEdge struct {
	Node   *PermissionsPolicy `json:"node"`
	Cursor Cursor             `json:"cursor"`
}

// PermissionsPolicyConnection is the connection containing edges to PermissionsPolicy.
type PermissionsPolicyConnection struct {
	Edges    []*PermissionsPolicyEdge `json:"edges"`
	PageInfo PageInfo                 `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to PermissionsPolicy.
func (pp *PermissionsPolicyQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*PermissionsPolicyConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &PermissionsPolicyConnection{
				Edges: []*PermissionsPolicyEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &PermissionsPolicyConnection{
				Edges: []*PermissionsPolicyEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		pp = pp.Where(permissionspolicy.IDGT(after.ID))
	}
	if before != nil {
		pp = pp.Where(permissionspolicy.IDLT(before.ID))
	}
	if first != nil {
		pp = pp.Order(Asc(permissionspolicy.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		pp = pp.Order(Desc(permissionspolicy.FieldID)).Limit(*last + 1)
	}
	pp = pp.collectConnectionFields(ctx)

	nodes, err := pp.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &PermissionsPolicyConnection{
			Edges: []*PermissionsPolicyEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn PermissionsPolicyConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*PermissionsPolicyEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &PermissionsPolicyEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (pp *PermissionsPolicyQuery) collectConnectionFields(ctx context.Context) *PermissionsPolicyQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		pp = pp.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return pp
}

// ProjectEdge is the edge representation of Project.
type ProjectEdge struct {
	Node   *Project `json:"node"`
	Cursor Cursor   `json:"cursor"`
}

// ProjectConnection is the connection containing edges to Project.
type ProjectConnection struct {
	Edges    []*ProjectEdge `json:"edges"`
	PageInfo PageInfo       `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Project.
func (pr *ProjectQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*ProjectConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &ProjectConnection{
				Edges: []*ProjectEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &ProjectConnection{
				Edges: []*ProjectEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		pr = pr.Where(project.IDGT(after.ID))
	}
	if before != nil {
		pr = pr.Where(project.IDLT(before.ID))
	}
	if first != nil {
		pr = pr.Order(Asc(project.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		pr = pr.Order(Desc(project.FieldID)).Limit(*last + 1)
	}
	pr = pr.collectConnectionFields(ctx)

	nodes, err := pr.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &ProjectConnection{
			Edges: []*ProjectEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn ProjectConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*ProjectEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &ProjectEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (pr *ProjectQuery) collectConnectionFields(ctx context.Context) *ProjectQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		pr = pr.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return pr
}

// ProjectTypeEdge is the edge representation of ProjectType.
type ProjectTypeEdge struct {
	Node   *ProjectType `json:"node"`
	Cursor Cursor       `json:"cursor"`
}

// ProjectTypeConnection is the connection containing edges to ProjectType.
type ProjectTypeConnection struct {
	Edges    []*ProjectTypeEdge `json:"edges"`
	PageInfo PageInfo           `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to ProjectType.
func (pt *ProjectTypeQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*ProjectTypeConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &ProjectTypeConnection{
				Edges: []*ProjectTypeEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &ProjectTypeConnection{
				Edges: []*ProjectTypeEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		pt = pt.Where(projecttype.IDGT(after.ID))
	}
	if before != nil {
		pt = pt.Where(projecttype.IDLT(before.ID))
	}
	if first != nil {
		pt = pt.Order(Asc(projecttype.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		pt = pt.Order(Desc(projecttype.FieldID)).Limit(*last + 1)
	}
	pt = pt.collectConnectionFields(ctx)

	nodes, err := pt.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &ProjectTypeConnection{
			Edges: []*ProjectTypeEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn ProjectTypeConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*ProjectTypeEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &ProjectTypeEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (pt *ProjectTypeQuery) collectConnectionFields(ctx context.Context) *ProjectTypeQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		pt = pt.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return pt
}

// PropertyEdge is the edge representation of Property.
type PropertyEdge struct {
	Node   *Property `json:"node"`
	Cursor Cursor    `json:"cursor"`
}

// PropertyConnection is the connection containing edges to Property.
type PropertyConnection struct {
	Edges    []*PropertyEdge `json:"edges"`
	PageInfo PageInfo        `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Property.
func (pr *PropertyQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*PropertyConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &PropertyConnection{
				Edges: []*PropertyEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &PropertyConnection{
				Edges: []*PropertyEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		pr = pr.Where(property.IDGT(after.ID))
	}
	if before != nil {
		pr = pr.Where(property.IDLT(before.ID))
	}
	if first != nil {
		pr = pr.Order(Asc(property.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		pr = pr.Order(Desc(property.FieldID)).Limit(*last + 1)
	}
	pr = pr.collectConnectionFields(ctx)

	nodes, err := pr.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &PropertyConnection{
			Edges: []*PropertyEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn PropertyConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*PropertyEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &PropertyEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (pr *PropertyQuery) collectConnectionFields(ctx context.Context) *PropertyQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		pr = pr.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return pr
}

// PropertyTypeEdge is the edge representation of PropertyType.
type PropertyTypeEdge struct {
	Node   *PropertyType `json:"node"`
	Cursor Cursor        `json:"cursor"`
}

// PropertyTypeConnection is the connection containing edges to PropertyType.
type PropertyTypeConnection struct {
	Edges    []*PropertyTypeEdge `json:"edges"`
	PageInfo PageInfo            `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to PropertyType.
func (pt *PropertyTypeQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*PropertyTypeConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &PropertyTypeConnection{
				Edges: []*PropertyTypeEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &PropertyTypeConnection{
				Edges: []*PropertyTypeEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		pt = pt.Where(propertytype.IDGT(after.ID))
	}
	if before != nil {
		pt = pt.Where(propertytype.IDLT(before.ID))
	}
	if first != nil {
		pt = pt.Order(Asc(propertytype.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		pt = pt.Order(Desc(propertytype.FieldID)).Limit(*last + 1)
	}
	pt = pt.collectConnectionFields(ctx)

	nodes, err := pt.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &PropertyTypeConnection{
			Edges: []*PropertyTypeEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn PropertyTypeConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*PropertyTypeEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &PropertyTypeEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (pt *PropertyTypeQuery) collectConnectionFields(ctx context.Context) *PropertyTypeQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		pt = pt.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return pt
}

// ReportFilterEdge is the edge representation of ReportFilter.
type ReportFilterEdge struct {
	Node   *ReportFilter `json:"node"`
	Cursor Cursor        `json:"cursor"`
}

// ReportFilterConnection is the connection containing edges to ReportFilter.
type ReportFilterConnection struct {
	Edges    []*ReportFilterEdge `json:"edges"`
	PageInfo PageInfo            `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to ReportFilter.
func (rf *ReportFilterQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*ReportFilterConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &ReportFilterConnection{
				Edges: []*ReportFilterEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &ReportFilterConnection{
				Edges: []*ReportFilterEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		rf = rf.Where(reportfilter.IDGT(after.ID))
	}
	if before != nil {
		rf = rf.Where(reportfilter.IDLT(before.ID))
	}
	if first != nil {
		rf = rf.Order(Asc(reportfilter.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		rf = rf.Order(Desc(reportfilter.FieldID)).Limit(*last + 1)
	}
	rf = rf.collectConnectionFields(ctx)

	nodes, err := rf.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &ReportFilterConnection{
			Edges: []*ReportFilterEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn ReportFilterConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*ReportFilterEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &ReportFilterEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (rf *ReportFilterQuery) collectConnectionFields(ctx context.Context) *ReportFilterQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		rf = rf.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return rf
}

// ServiceEdge is the edge representation of Service.
type ServiceEdge struct {
	Node   *Service `json:"node"`
	Cursor Cursor   `json:"cursor"`
}

// ServiceConnection is the connection containing edges to Service.
type ServiceConnection struct {
	Edges    []*ServiceEdge `json:"edges"`
	PageInfo PageInfo       `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Service.
func (s *ServiceQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*ServiceConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &ServiceConnection{
				Edges: []*ServiceEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &ServiceConnection{
				Edges: []*ServiceEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		s = s.Where(service.IDGT(after.ID))
	}
	if before != nil {
		s = s.Where(service.IDLT(before.ID))
	}
	if first != nil {
		s = s.Order(Asc(service.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		s = s.Order(Desc(service.FieldID)).Limit(*last + 1)
	}
	s = s.collectConnectionFields(ctx)

	nodes, err := s.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &ServiceConnection{
			Edges: []*ServiceEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn ServiceConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*ServiceEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &ServiceEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (s *ServiceQuery) collectConnectionFields(ctx context.Context) *ServiceQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		s = s.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return s
}

// ServiceEndpointEdge is the edge representation of ServiceEndpoint.
type ServiceEndpointEdge struct {
	Node   *ServiceEndpoint `json:"node"`
	Cursor Cursor           `json:"cursor"`
}

// ServiceEndpointConnection is the connection containing edges to ServiceEndpoint.
type ServiceEndpointConnection struct {
	Edges    []*ServiceEndpointEdge `json:"edges"`
	PageInfo PageInfo               `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to ServiceEndpoint.
func (se *ServiceEndpointQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*ServiceEndpointConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &ServiceEndpointConnection{
				Edges: []*ServiceEndpointEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &ServiceEndpointConnection{
				Edges: []*ServiceEndpointEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		se = se.Where(serviceendpoint.IDGT(after.ID))
	}
	if before != nil {
		se = se.Where(serviceendpoint.IDLT(before.ID))
	}
	if first != nil {
		se = se.Order(Asc(serviceendpoint.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		se = se.Order(Desc(serviceendpoint.FieldID)).Limit(*last + 1)
	}
	se = se.collectConnectionFields(ctx)

	nodes, err := se.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &ServiceEndpointConnection{
			Edges: []*ServiceEndpointEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn ServiceEndpointConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*ServiceEndpointEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &ServiceEndpointEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (se *ServiceEndpointQuery) collectConnectionFields(ctx context.Context) *ServiceEndpointQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		se = se.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return se
}

// ServiceEndpointDefinitionEdge is the edge representation of ServiceEndpointDefinition.
type ServiceEndpointDefinitionEdge struct {
	Node   *ServiceEndpointDefinition `json:"node"`
	Cursor Cursor                     `json:"cursor"`
}

// ServiceEndpointDefinitionConnection is the connection containing edges to ServiceEndpointDefinition.
type ServiceEndpointDefinitionConnection struct {
	Edges    []*ServiceEndpointDefinitionEdge `json:"edges"`
	PageInfo PageInfo                         `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to ServiceEndpointDefinition.
func (sed *ServiceEndpointDefinitionQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*ServiceEndpointDefinitionConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &ServiceEndpointDefinitionConnection{
				Edges: []*ServiceEndpointDefinitionEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &ServiceEndpointDefinitionConnection{
				Edges: []*ServiceEndpointDefinitionEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		sed = sed.Where(serviceendpointdefinition.IDGT(after.ID))
	}
	if before != nil {
		sed = sed.Where(serviceendpointdefinition.IDLT(before.ID))
	}
	if first != nil {
		sed = sed.Order(Asc(serviceendpointdefinition.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		sed = sed.Order(Desc(serviceendpointdefinition.FieldID)).Limit(*last + 1)
	}
	sed = sed.collectConnectionFields(ctx)

	nodes, err := sed.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &ServiceEndpointDefinitionConnection{
			Edges: []*ServiceEndpointDefinitionEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn ServiceEndpointDefinitionConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*ServiceEndpointDefinitionEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &ServiceEndpointDefinitionEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (sed *ServiceEndpointDefinitionQuery) collectConnectionFields(ctx context.Context) *ServiceEndpointDefinitionQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		sed = sed.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return sed
}

// ServiceTypeEdge is the edge representation of ServiceType.
type ServiceTypeEdge struct {
	Node   *ServiceType `json:"node"`
	Cursor Cursor       `json:"cursor"`
}

// ServiceTypeConnection is the connection containing edges to ServiceType.
type ServiceTypeConnection struct {
	Edges    []*ServiceTypeEdge `json:"edges"`
	PageInfo PageInfo           `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to ServiceType.
func (st *ServiceTypeQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*ServiceTypeConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &ServiceTypeConnection{
				Edges: []*ServiceTypeEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &ServiceTypeConnection{
				Edges: []*ServiceTypeEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		st = st.Where(servicetype.IDGT(after.ID))
	}
	if before != nil {
		st = st.Where(servicetype.IDLT(before.ID))
	}
	if first != nil {
		st = st.Order(Asc(servicetype.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		st = st.Order(Desc(servicetype.FieldID)).Limit(*last + 1)
	}
	st = st.collectConnectionFields(ctx)

	nodes, err := st.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &ServiceTypeConnection{
			Edges: []*ServiceTypeEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn ServiceTypeConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*ServiceTypeEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &ServiceTypeEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (st *ServiceTypeQuery) collectConnectionFields(ctx context.Context) *ServiceTypeQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		st = st.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return st
}

// SurveyEdge is the edge representation of Survey.
type SurveyEdge struct {
	Node   *Survey `json:"node"`
	Cursor Cursor  `json:"cursor"`
}

// SurveyConnection is the connection containing edges to Survey.
type SurveyConnection struct {
	Edges    []*SurveyEdge `json:"edges"`
	PageInfo PageInfo      `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Survey.
func (s *SurveyQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*SurveyConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &SurveyConnection{
				Edges: []*SurveyEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &SurveyConnection{
				Edges: []*SurveyEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		s = s.Where(survey.IDGT(after.ID))
	}
	if before != nil {
		s = s.Where(survey.IDLT(before.ID))
	}
	if first != nil {
		s = s.Order(Asc(survey.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		s = s.Order(Desc(survey.FieldID)).Limit(*last + 1)
	}
	s = s.collectConnectionFields(ctx)

	nodes, err := s.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &SurveyConnection{
			Edges: []*SurveyEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn SurveyConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*SurveyEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &SurveyEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (s *SurveyQuery) collectConnectionFields(ctx context.Context) *SurveyQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		s = s.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return s
}

// SurveyCellScanEdge is the edge representation of SurveyCellScan.
type SurveyCellScanEdge struct {
	Node   *SurveyCellScan `json:"node"`
	Cursor Cursor          `json:"cursor"`
}

// SurveyCellScanConnection is the connection containing edges to SurveyCellScan.
type SurveyCellScanConnection struct {
	Edges    []*SurveyCellScanEdge `json:"edges"`
	PageInfo PageInfo              `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to SurveyCellScan.
func (scs *SurveyCellScanQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*SurveyCellScanConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &SurveyCellScanConnection{
				Edges: []*SurveyCellScanEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &SurveyCellScanConnection{
				Edges: []*SurveyCellScanEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		scs = scs.Where(surveycellscan.IDGT(after.ID))
	}
	if before != nil {
		scs = scs.Where(surveycellscan.IDLT(before.ID))
	}
	if first != nil {
		scs = scs.Order(Asc(surveycellscan.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		scs = scs.Order(Desc(surveycellscan.FieldID)).Limit(*last + 1)
	}
	scs = scs.collectConnectionFields(ctx)

	nodes, err := scs.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &SurveyCellScanConnection{
			Edges: []*SurveyCellScanEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn SurveyCellScanConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*SurveyCellScanEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &SurveyCellScanEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (scs *SurveyCellScanQuery) collectConnectionFields(ctx context.Context) *SurveyCellScanQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		scs = scs.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return scs
}

// SurveyQuestionEdge is the edge representation of SurveyQuestion.
type SurveyQuestionEdge struct {
	Node   *SurveyQuestion `json:"node"`
	Cursor Cursor          `json:"cursor"`
}

// SurveyQuestionConnection is the connection containing edges to SurveyQuestion.
type SurveyQuestionConnection struct {
	Edges    []*SurveyQuestionEdge `json:"edges"`
	PageInfo PageInfo              `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to SurveyQuestion.
func (sq *SurveyQuestionQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*SurveyQuestionConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &SurveyQuestionConnection{
				Edges: []*SurveyQuestionEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &SurveyQuestionConnection{
				Edges: []*SurveyQuestionEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		sq = sq.Where(surveyquestion.IDGT(after.ID))
	}
	if before != nil {
		sq = sq.Where(surveyquestion.IDLT(before.ID))
	}
	if first != nil {
		sq = sq.Order(Asc(surveyquestion.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		sq = sq.Order(Desc(surveyquestion.FieldID)).Limit(*last + 1)
	}
	sq = sq.collectConnectionFields(ctx)

	nodes, err := sq.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &SurveyQuestionConnection{
			Edges: []*SurveyQuestionEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn SurveyQuestionConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*SurveyQuestionEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &SurveyQuestionEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (sq *SurveyQuestionQuery) collectConnectionFields(ctx context.Context) *SurveyQuestionQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		sq = sq.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return sq
}

// SurveyTemplateCategoryEdge is the edge representation of SurveyTemplateCategory.
type SurveyTemplateCategoryEdge struct {
	Node   *SurveyTemplateCategory `json:"node"`
	Cursor Cursor                  `json:"cursor"`
}

// SurveyTemplateCategoryConnection is the connection containing edges to SurveyTemplateCategory.
type SurveyTemplateCategoryConnection struct {
	Edges    []*SurveyTemplateCategoryEdge `json:"edges"`
	PageInfo PageInfo                      `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to SurveyTemplateCategory.
func (stc *SurveyTemplateCategoryQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*SurveyTemplateCategoryConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &SurveyTemplateCategoryConnection{
				Edges: []*SurveyTemplateCategoryEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &SurveyTemplateCategoryConnection{
				Edges: []*SurveyTemplateCategoryEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		stc = stc.Where(surveytemplatecategory.IDGT(after.ID))
	}
	if before != nil {
		stc = stc.Where(surveytemplatecategory.IDLT(before.ID))
	}
	if first != nil {
		stc = stc.Order(Asc(surveytemplatecategory.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		stc = stc.Order(Desc(surveytemplatecategory.FieldID)).Limit(*last + 1)
	}
	stc = stc.collectConnectionFields(ctx)

	nodes, err := stc.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &SurveyTemplateCategoryConnection{
			Edges: []*SurveyTemplateCategoryEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn SurveyTemplateCategoryConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*SurveyTemplateCategoryEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &SurveyTemplateCategoryEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (stc *SurveyTemplateCategoryQuery) collectConnectionFields(ctx context.Context) *SurveyTemplateCategoryQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		stc = stc.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return stc
}

// SurveyTemplateQuestionEdge is the edge representation of SurveyTemplateQuestion.
type SurveyTemplateQuestionEdge struct {
	Node   *SurveyTemplateQuestion `json:"node"`
	Cursor Cursor                  `json:"cursor"`
}

// SurveyTemplateQuestionConnection is the connection containing edges to SurveyTemplateQuestion.
type SurveyTemplateQuestionConnection struct {
	Edges    []*SurveyTemplateQuestionEdge `json:"edges"`
	PageInfo PageInfo                      `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to SurveyTemplateQuestion.
func (stq *SurveyTemplateQuestionQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*SurveyTemplateQuestionConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &SurveyTemplateQuestionConnection{
				Edges: []*SurveyTemplateQuestionEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &SurveyTemplateQuestionConnection{
				Edges: []*SurveyTemplateQuestionEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		stq = stq.Where(surveytemplatequestion.IDGT(after.ID))
	}
	if before != nil {
		stq = stq.Where(surveytemplatequestion.IDLT(before.ID))
	}
	if first != nil {
		stq = stq.Order(Asc(surveytemplatequestion.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		stq = stq.Order(Desc(surveytemplatequestion.FieldID)).Limit(*last + 1)
	}
	stq = stq.collectConnectionFields(ctx)

	nodes, err := stq.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &SurveyTemplateQuestionConnection{
			Edges: []*SurveyTemplateQuestionEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn SurveyTemplateQuestionConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*SurveyTemplateQuestionEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &SurveyTemplateQuestionEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (stq *SurveyTemplateQuestionQuery) collectConnectionFields(ctx context.Context) *SurveyTemplateQuestionQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		stq = stq.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return stq
}

// SurveyWiFiScanEdge is the edge representation of SurveyWiFiScan.
type SurveyWiFiScanEdge struct {
	Node   *SurveyWiFiScan `json:"node"`
	Cursor Cursor          `json:"cursor"`
}

// SurveyWiFiScanConnection is the connection containing edges to SurveyWiFiScan.
type SurveyWiFiScanConnection struct {
	Edges    []*SurveyWiFiScanEdge `json:"edges"`
	PageInfo PageInfo              `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to SurveyWiFiScan.
func (swfs *SurveyWiFiScanQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*SurveyWiFiScanConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &SurveyWiFiScanConnection{
				Edges: []*SurveyWiFiScanEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &SurveyWiFiScanConnection{
				Edges: []*SurveyWiFiScanEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		swfs = swfs.Where(surveywifiscan.IDGT(after.ID))
	}
	if before != nil {
		swfs = swfs.Where(surveywifiscan.IDLT(before.ID))
	}
	if first != nil {
		swfs = swfs.Order(Asc(surveywifiscan.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		swfs = swfs.Order(Desc(surveywifiscan.FieldID)).Limit(*last + 1)
	}
	swfs = swfs.collectConnectionFields(ctx)

	nodes, err := swfs.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &SurveyWiFiScanConnection{
			Edges: []*SurveyWiFiScanEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn SurveyWiFiScanConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*SurveyWiFiScanEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &SurveyWiFiScanEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (swfs *SurveyWiFiScanQuery) collectConnectionFields(ctx context.Context) *SurveyWiFiScanQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		swfs = swfs.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return swfs
}

// UserEdge is the edge representation of User.
type UserEdge struct {
	Node   *User  `json:"node"`
	Cursor Cursor `json:"cursor"`
}

// UserConnection is the connection containing edges to User.
type UserConnection struct {
	Edges    []*UserEdge `json:"edges"`
	PageInfo PageInfo    `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to User.
func (u *UserQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*UserConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &UserConnection{
				Edges: []*UserEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &UserConnection{
				Edges: []*UserEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		u = u.Where(user.IDGT(after.ID))
	}
	if before != nil {
		u = u.Where(user.IDLT(before.ID))
	}
	if first != nil {
		u = u.Order(Asc(user.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		u = u.Order(Desc(user.FieldID)).Limit(*last + 1)
	}
	u = u.collectConnectionFields(ctx)

	nodes, err := u.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &UserConnection{
			Edges: []*UserEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn UserConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*UserEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &UserEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (u *UserQuery) collectConnectionFields(ctx context.Context) *UserQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		u = u.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return u
}

// UsersGroupEdge is the edge representation of UsersGroup.
type UsersGroupEdge struct {
	Node   *UsersGroup `json:"node"`
	Cursor Cursor      `json:"cursor"`
}

// UsersGroupConnection is the connection containing edges to UsersGroup.
type UsersGroupConnection struct {
	Edges    []*UsersGroupEdge `json:"edges"`
	PageInfo PageInfo          `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to UsersGroup.
func (ug *UsersGroupQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*UsersGroupConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &UsersGroupConnection{
				Edges: []*UsersGroupEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &UsersGroupConnection{
				Edges: []*UsersGroupEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		ug = ug.Where(usersgroup.IDGT(after.ID))
	}
	if before != nil {
		ug = ug.Where(usersgroup.IDLT(before.ID))
	}
	if first != nil {
		ug = ug.Order(Asc(usersgroup.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		ug = ug.Order(Desc(usersgroup.FieldID)).Limit(*last + 1)
	}
	ug = ug.collectConnectionFields(ctx)

	nodes, err := ug.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &UsersGroupConnection{
			Edges: []*UsersGroupEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn UsersGroupConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*UsersGroupEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &UsersGroupEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (ug *UsersGroupQuery) collectConnectionFields(ctx context.Context) *UsersGroupQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		ug = ug.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return ug
}

// WorkOrderEdge is the edge representation of WorkOrder.
type WorkOrderEdge struct {
	Node   *WorkOrder `json:"node"`
	Cursor Cursor     `json:"cursor"`
}

// WorkOrderConnection is the connection containing edges to WorkOrder.
type WorkOrderConnection struct {
	Edges    []*WorkOrderEdge `json:"edges"`
	PageInfo PageInfo         `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to WorkOrder.
func (wo *WorkOrderQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*WorkOrderConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &WorkOrderConnection{
				Edges: []*WorkOrderEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &WorkOrderConnection{
				Edges: []*WorkOrderEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		wo = wo.Where(workorder.IDGT(after.ID))
	}
	if before != nil {
		wo = wo.Where(workorder.IDLT(before.ID))
	}
	if first != nil {
		wo = wo.Order(Asc(workorder.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		wo = wo.Order(Desc(workorder.FieldID)).Limit(*last + 1)
	}
	wo = wo.collectConnectionFields(ctx)

	nodes, err := wo.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &WorkOrderConnection{
			Edges: []*WorkOrderEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn WorkOrderConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*WorkOrderEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &WorkOrderEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (wo *WorkOrderQuery) collectConnectionFields(ctx context.Context) *WorkOrderQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		wo = wo.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return wo
}

// WorkOrderDefinitionEdge is the edge representation of WorkOrderDefinition.
type WorkOrderDefinitionEdge struct {
	Node   *WorkOrderDefinition `json:"node"`
	Cursor Cursor               `json:"cursor"`
}

// WorkOrderDefinitionConnection is the connection containing edges to WorkOrderDefinition.
type WorkOrderDefinitionConnection struct {
	Edges    []*WorkOrderDefinitionEdge `json:"edges"`
	PageInfo PageInfo                   `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to WorkOrderDefinition.
func (wod *WorkOrderDefinitionQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*WorkOrderDefinitionConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &WorkOrderDefinitionConnection{
				Edges: []*WorkOrderDefinitionEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &WorkOrderDefinitionConnection{
				Edges: []*WorkOrderDefinitionEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		wod = wod.Where(workorderdefinition.IDGT(after.ID))
	}
	if before != nil {
		wod = wod.Where(workorderdefinition.IDLT(before.ID))
	}
	if first != nil {
		wod = wod.Order(Asc(workorderdefinition.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		wod = wod.Order(Desc(workorderdefinition.FieldID)).Limit(*last + 1)
	}
	wod = wod.collectConnectionFields(ctx)

	nodes, err := wod.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &WorkOrderDefinitionConnection{
			Edges: []*WorkOrderDefinitionEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn WorkOrderDefinitionConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*WorkOrderDefinitionEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &WorkOrderDefinitionEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (wod *WorkOrderDefinitionQuery) collectConnectionFields(ctx context.Context) *WorkOrderDefinitionQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		wod = wod.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return wod
}

// WorkOrderTypeEdge is the edge representation of WorkOrderType.
type WorkOrderTypeEdge struct {
	Node   *WorkOrderType `json:"node"`
	Cursor Cursor         `json:"cursor"`
}

// WorkOrderTypeConnection is the connection containing edges to WorkOrderType.
type WorkOrderTypeConnection struct {
	Edges    []*WorkOrderTypeEdge `json:"edges"`
	PageInfo PageInfo             `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to WorkOrderType.
func (wot *WorkOrderTypeQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*WorkOrderTypeConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &WorkOrderTypeConnection{
				Edges: []*WorkOrderTypeEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &WorkOrderTypeConnection{
				Edges: []*WorkOrderTypeEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		wot = wot.Where(workordertype.IDGT(after.ID))
	}
	if before != nil {
		wot = wot.Where(workordertype.IDLT(before.ID))
	}
	if first != nil {
		wot = wot.Order(Asc(workordertype.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		wot = wot.Order(Desc(workordertype.FieldID)).Limit(*last + 1)
	}
	wot = wot.collectConnectionFields(ctx)

	nodes, err := wot.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &WorkOrderTypeConnection{
			Edges: []*WorkOrderTypeEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn WorkOrderTypeConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*WorkOrderTypeEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &WorkOrderTypeEdge{
			Node: node,
			Cursor: Cursor{
				ID: node.ID,
			},
		}
	}
	conn.PageInfo.StartCursor = &conn.Edges[0].Cursor
	conn.PageInfo.EndCursor = &conn.Edges[len(conn.Edges)-1].Cursor

	return &conn, nil
}

func (wot *WorkOrderTypeQuery) collectConnectionFields(ctx context.Context) *WorkOrderTypeQuery {
	if field := fieldForPath(ctx, "edges", "node"); field != nil {
		wot = wot.collectField(graphql.GetOperationContext(ctx), *field)
	}
	return wot
}

func fieldForPath(ctx context.Context, path ...string) *graphql.CollectedField {
	fc := graphql.GetFieldContext(ctx)
	if fc == nil {
		return nil
	}
	oc := graphql.GetOperationContext(ctx)
	field := fc.Field

walk:
	for _, name := range path {
		for _, f := range graphql.CollectFields(oc, field.Selections, nil) {
			if f.Name == name {
				field = f
				continue walk
			}
		}
		return nil
	}
	return &field
}
