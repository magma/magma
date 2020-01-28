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

	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/checklistitemdefinition"
	"github.com/facebookincubator/symphony/graph/ent/comment"
	"github.com/facebookincubator/symphony/graph/ent/customer"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/equipmentcategory"
	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentporttype"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/equipmenttype"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/floorplan"
	"github.com/facebookincubator/symphony/graph/ent/floorplanreferencepoint"
	"github.com/facebookincubator/symphony/graph/ent/floorplanscale"
	"github.com/facebookincubator/symphony/graph/ent/hyperlink"
	"github.com/facebookincubator/symphony/graph/ent/link"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/project"
	"github.com/facebookincubator/symphony/graph/ent/projecttype"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/servicetype"
	"github.com/facebookincubator/symphony/graph/ent/survey"
	"github.com/facebookincubator/symphony/graph/ent/surveycellscan"
	"github.com/facebookincubator/symphony/graph/ent/surveyquestion"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatecategory"
	"github.com/facebookincubator/symphony/graph/ent/surveytemplatequestion"
	"github.com/facebookincubator/symphony/graph/ent/surveywifiscan"
	"github.com/facebookincubator/symphony/graph/ent/technician"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"
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
	ID string
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

// TechnicianEdge is the edge representation of Technician.
type TechnicianEdge struct {
	Node   *Technician `json:"node"`
	Cursor Cursor      `json:"cursor"`
}

// TechnicianConnection is the connection containing edges to Technician.
type TechnicianConnection struct {
	Edges    []*TechnicianEdge `json:"edges"`
	PageInfo PageInfo          `json:"pageInfo"`
}

// Paginate executes the query and returns a relay based cursor connection to Technician.
func (t *TechnicianQuery) Paginate(ctx context.Context, after *Cursor, first *int, before *Cursor, last *int) (*TechnicianConnection, error) {
	if first != nil && last != nil {
		return nil, ErrInvalidPagination
	}
	if first != nil {
		if *first == 0 {
			return &TechnicianConnection{
				Edges: []*TechnicianEdge{},
			}, nil
		} else if *first < 0 {
			return nil, ErrInvalidPagination
		}
	}
	if last != nil {
		if *last == 0 {
			return &TechnicianConnection{
				Edges: []*TechnicianEdge{},
			}, nil
		} else if *last < 0 {
			return nil, ErrInvalidPagination
		}
	}

	if after != nil {
		t = t.Where(technician.IDGT(after.ID))
	}
	if before != nil {
		t = t.Where(technician.IDLT(before.ID))
	}
	if first != nil {
		t = t.Order(Asc(technician.FieldID)).Limit(*first + 1)
	}
	if last != nil {
		t = t.Order(Desc(technician.FieldID)).Limit(*last + 1)
	}

	nodes, err := t.All(ctx)
	if err != nil || len(nodes) == 0 {
		return &TechnicianConnection{
			Edges: []*TechnicianEdge{},
		}, err
	}
	if last != nil {
		for left, right := 0, len(nodes)-1; left < right; left, right = left+1, right-1 {
			nodes[left], nodes[right] = nodes[right], nodes[left]
		}
	}

	var conn TechnicianConnection
	if first != nil && len(nodes) > *first {
		conn.PageInfo.HasNextPage = true
		nodes = nodes[:len(nodes)-1]
	} else if last != nil && len(nodes) > *last {
		conn.PageInfo.HasPreviousPage = true
		nodes = nodes[1:]
	}
	conn.Edges = make([]*TechnicianEdge, len(nodes))
	for i, node := range nodes {
		conn.Edges[i] = &TechnicianEdge{
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
