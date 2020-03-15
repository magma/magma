// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/schema"
	"github.com/facebookincubator/symphony/graph/ent/actionsrule"
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
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
	"github.com/facebookincubator/symphony/graph/ent/reportfilter"
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
	"github.com/facebookincubator/symphony/graph/ent/user"
	"github.com/facebookincubator/symphony/graph/ent/workorder"
	"github.com/facebookincubator/symphony/graph/ent/workorderdefinition"
	"github.com/facebookincubator/symphony/graph/ent/workordertype"

	"golang.org/x/sync/semaphore"
)

// Noder wraps the basic Node method.
type Noder interface {
	Node(context.Context) (*Node, error)
}

// Node in the graph.
type Node struct {
	ID     int      `json:"id,omitemty"`      // node id.
	Type   string   `json:"type,omitempty"`   // node type.
	Fields []*Field `json:"fields,omitempty"` // node fields.
	Edges  []*Edge  `json:"edges,omitempty"`  // node edges.
}

// Field of a node.
type Field struct {
	Type  string `json:"type,omitempty"`  // field type.
	Name  string `json:"name,omitempty"`  // field name (as in struct).
	Value string `json:"value,omitempty"` // stringified value.
}

// Edges between two nodes.
type Edge struct {
	Type string `json:"type,omitempty"` // edge type.
	Name string `json:"name,omitempty"` // edge name.
	IDs  []int  `json:"ids,omitempty"`  // node ids (where this edge point to).
}

func (ar *ActionsRule) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     ar.ID,
		Type:   "ActionsRule",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(ar.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ar.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ar.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ar.TriggerID); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "TriggerID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ar.RuleFilters); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "[]*core.ActionsRuleFilter",
		Name:  "RuleFilters",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ar.RuleActions); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "[]*core.ActionsRuleAction",
		Name:  "RuleActions",
		Value: string(buf),
	}
	return node, nil
}

func (clc *CheckListCategory) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     clc.ID,
		Type:   "CheckListCategory",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(clc.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clc.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clc.Title); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Title",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clc.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Description",
		Value: string(buf),
	}
	var ids []int
	ids, err = clc.QueryCheckListItems().
		Select(checklistitem.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "CheckListItem",
		Name: "CheckListItems",
	}
	return node, nil
}

func (cli *CheckListItem) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     cli.ID,
		Type:   "CheckListItem",
		Fields: make([]*Field, 7),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(cli.Title); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "Title",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.Type); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "Type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.Index); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "Index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.Checked); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "bool",
		Name:  "Checked",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.StringVal); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "StringVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.EnumValues); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "EnumValues",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.HelpText); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "string",
		Name:  "HelpText",
		Value: string(buf),
	}
	var ids []int
	ids, err = cli.QueryWorkOrder().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "WorkOrder",
	}
	return node, nil
}

func (clid *CheckListItemDefinition) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     clid.ID,
		Type:   "CheckListItemDefinition",
		Fields: make([]*Field, 7),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(clid.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.Title); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Title",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.Type); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.Index); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "Index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.EnumValues); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "EnumValues",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.HelpText); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "string",
		Name:  "HelpText",
		Value: string(buf),
	}
	var ids []int
	ids, err = clid.QueryWorkOrderType().
		Select(workordertype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "WorkOrderType",
		Name: "WorkOrderType",
	}
	return node, nil
}

func (c *Comment) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     c.ID,
		Type:   "Comment",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(c.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.AuthorName); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "AuthorName",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.Text); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Text",
		Value: string(buf),
	}
	return node, nil
}

func (c *Customer) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     c.ID,
		Type:   "Customer",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(c.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.ExternalID); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "ExternalID",
		Value: string(buf),
	}
	var ids []int
	ids, err = c.QueryServices().
		Select(service.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "Service",
		Name: "Services",
	}
	return node, nil
}

func (e *Equipment) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     e.ID,
		Type:   "Equipment",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 9),
	}
	var buf []byte
	if buf, err = json.Marshal(e.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(e.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(e.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(e.FutureState); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "FutureState",
		Value: string(buf),
	}
	if buf, err = json.Marshal(e.DeviceID); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "DeviceID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(e.ExternalID); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "ExternalID",
		Value: string(buf),
	}
	var ids []int
	ids, err = e.QueryType().
		Select(equipmenttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "EquipmentType",
		Name: "Type",
	}
	ids, err = e.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Location",
	}
	ids, err = e.QueryParentPosition().
		Select(equipmentposition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "EquipmentPosition",
		Name: "ParentPosition",
	}
	ids, err = e.QueryPositions().
		Select(equipmentposition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "EquipmentPosition",
		Name: "Positions",
	}
	ids, err = e.QueryPorts().
		Select(equipmentport.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "EquipmentPort",
		Name: "Ports",
	}
	ids, err = e.QueryWorkOrder().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[5] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "WorkOrder",
	}
	ids, err = e.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[6] = &Edge{
		IDs:  ids,
		Type: "Property",
		Name: "Properties",
	}
	ids, err = e.QueryFiles().
		Select(file.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[7] = &Edge{
		IDs:  ids,
		Type: "File",
		Name: "Files",
	}
	ids, err = e.QueryHyperlinks().
		Select(hyperlink.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[8] = &Edge{
		IDs:  ids,
		Type: "Hyperlink",
		Name: "Hyperlinks",
	}
	return node, nil
}

func (ec *EquipmentCategory) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     ec.ID,
		Type:   "EquipmentCategory",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(ec.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ec.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ec.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	var ids []int
	ids, err = ec.QueryTypes().
		Select(equipmenttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "EquipmentType",
		Name: "Types",
	}
	return node, nil
}

func (ep *EquipmentPort) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     ep.ID,
		Type:   "EquipmentPort",
		Fields: make([]*Field, 2),
		Edges:  make([]*Edge, 5),
	}
	var buf []byte
	if buf, err = json.Marshal(ep.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ep.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	var ids []int
	ids, err = ep.QueryDefinition().
		Select(equipmentportdefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "EquipmentPortDefinition",
		Name: "Definition",
	}
	ids, err = ep.QueryParent().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "Parent",
	}
	ids, err = ep.QueryLink().
		Select(link.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Link",
		Name: "Link",
	}
	ids, err = ep.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "Property",
		Name: "Properties",
	}
	ids, err = ep.QueryEndpoints().
		Select(serviceendpoint.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "ServiceEndpoint",
		Name: "Endpoints",
	}
	return node, nil
}

func (epd *EquipmentPortDefinition) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     epd.ID,
		Type:   "EquipmentPortDefinition",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(epd.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.Index); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "Index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.Bandwidth); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "Bandwidth",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.VisibilityLabel); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "VisibilityLabel",
		Value: string(buf),
	}
	var ids []int
	ids, err = epd.QueryEquipmentPortType().
		Select(equipmentporttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "EquipmentPortType",
		Name: "EquipmentPortType",
	}
	ids, err = epd.QueryPorts().
		Select(equipmentport.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "EquipmentPort",
		Name: "Ports",
	}
	ids, err = epd.QueryEquipmentType().
		Select(equipmenttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "EquipmentType",
		Name: "EquipmentType",
	}
	return node, nil
}

func (ept *EquipmentPortType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     ept.ID,
		Type:   "EquipmentPortType",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(ept.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ept.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ept.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	var ids []int
	ids, err = ept.QueryPropertyTypes().
		Select(propertytype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "PropertyType",
		Name: "PropertyTypes",
	}
	ids, err = ept.QueryLinkPropertyTypes().
		Select(propertytype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "PropertyType",
		Name: "LinkPropertyTypes",
	}
	ids, err = ept.QueryPortDefinitions().
		Select(equipmentportdefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "EquipmentPortDefinition",
		Name: "PortDefinitions",
	}
	return node, nil
}

func (ep *EquipmentPosition) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     ep.ID,
		Type:   "EquipmentPosition",
		Fields: make([]*Field, 2),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(ep.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ep.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	var ids []int
	ids, err = ep.QueryDefinition().
		Select(equipmentpositiondefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "EquipmentPositionDefinition",
		Name: "Definition",
	}
	ids, err = ep.QueryParent().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "Parent",
	}
	ids, err = ep.QueryAttachment().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "Attachment",
	}
	return node, nil
}

func (epd *EquipmentPositionDefinition) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     epd.ID,
		Type:   "EquipmentPositionDefinition",
		Fields: make([]*Field, 5),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(epd.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.Index); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "Index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.VisibilityLabel); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "VisibilityLabel",
		Value: string(buf),
	}
	var ids []int
	ids, err = epd.QueryPositions().
		Select(equipmentposition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "EquipmentPosition",
		Name: "Positions",
	}
	ids, err = epd.QueryEquipmentType().
		Select(equipmenttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "EquipmentType",
		Name: "EquipmentType",
	}
	return node, nil
}

func (et *EquipmentType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     et.ID,
		Type:   "EquipmentType",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 5),
	}
	var buf []byte
	if buf, err = json.Marshal(et.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(et.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(et.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	var ids []int
	ids, err = et.QueryPortDefinitions().
		Select(equipmentportdefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "EquipmentPortDefinition",
		Name: "PortDefinitions",
	}
	ids, err = et.QueryPositionDefinitions().
		Select(equipmentpositiondefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "EquipmentPositionDefinition",
		Name: "PositionDefinitions",
	}
	ids, err = et.QueryPropertyTypes().
		Select(propertytype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "PropertyType",
		Name: "PropertyTypes",
	}
	ids, err = et.QueryEquipment().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "Equipment",
	}
	ids, err = et.QueryCategory().
		Select(equipmentcategory.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "EquipmentCategory",
		Name: "Category",
	}
	return node, nil
}

func (f *File) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     f.ID,
		Type:   "File",
		Fields: make([]*Field, 10),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(f.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Type); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Name); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Size); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "Size",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.ModifiedAt); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "time.Time",
		Name:  "ModifiedAt",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.UploadedAt); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "time.Time",
		Name:  "UploadedAt",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.ContentType); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "ContentType",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.StoreKey); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "string",
		Name:  "StoreKey",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Category); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "Category",
		Value: string(buf),
	}
	return node, nil
}

func (fp *FloorPlan) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     fp.ID,
		Type:   "FloorPlan",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 4),
	}
	var buf []byte
	if buf, err = json.Marshal(fp.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fp.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fp.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	var ids []int
	ids, err = fp.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Location",
	}
	ids, err = fp.QueryReferencePoint().
		Select(floorplanreferencepoint.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "FloorPlanReferencePoint",
		Name: "ReferencePoint",
	}
	ids, err = fp.QueryScale().
		Select(floorplanscale.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "FloorPlanScale",
		Name: "Scale",
	}
	ids, err = fp.QueryImage().
		Select(file.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "File",
		Name: "Image",
	}
	return node, nil
}

func (fprp *FloorPlanReferencePoint) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     fprp.ID,
		Type:   "FloorPlanReferencePoint",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(fprp.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fprp.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fprp.X); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "X",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fprp.Y); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "Y",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fprp.Latitude); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "float64",
		Name:  "Latitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fprp.Longitude); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "float64",
		Name:  "Longitude",
		Value: string(buf),
	}
	return node, nil
}

func (fps *FloorPlanScale) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     fps.ID,
		Type:   "FloorPlanScale",
		Fields: make([]*Field, 7),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(fps.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.ReferencePoint1X); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "ReferencePoint1X",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.ReferencePoint1Y); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "ReferencePoint1Y",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.ReferencePoint2X); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "ReferencePoint2X",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.ReferencePoint2Y); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "int",
		Name:  "ReferencePoint2Y",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.ScaleInMeters); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "float64",
		Name:  "ScaleInMeters",
		Value: string(buf),
	}
	return node, nil
}

func (h *Hyperlink) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     h.ID,
		Type:   "Hyperlink",
		Fields: make([]*Field, 5),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(h.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(h.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(h.URL); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "URL",
		Value: string(buf),
	}
	if buf, err = json.Marshal(h.Name); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(h.Category); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "Category",
		Value: string(buf),
	}
	return node, nil
}

func (l *Link) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     l.ID,
		Type:   "Link",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 4),
	}
	var buf []byte
	if buf, err = json.Marshal(l.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.FutureState); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "FutureState",
		Value: string(buf),
	}
	var ids []int
	ids, err = l.QueryPorts().
		Select(equipmentport.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "EquipmentPort",
		Name: "Ports",
	}
	ids, err = l.QueryWorkOrder().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "WorkOrder",
	}
	ids, err = l.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Property",
		Name: "Properties",
	}
	ids, err = l.QueryService().
		Select(service.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "Service",
		Name: "Service",
	}
	return node, nil
}

func (l *Location) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     l.ID,
		Type:   "Location",
		Fields: make([]*Field, 7),
		Edges:  make([]*Edge, 12),
	}
	var buf []byte
	if buf, err = json.Marshal(l.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.ExternalID); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "ExternalID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.Latitude); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "float64",
		Name:  "Latitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.Longitude); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "float64",
		Name:  "Longitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.SiteSurveyNeeded); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "bool",
		Name:  "SiteSurveyNeeded",
		Value: string(buf),
	}
	var ids []int
	ids, err = l.QueryType().
		Select(locationtype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "LocationType",
		Name: "Type",
	}
	ids, err = l.QueryParent().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Parent",
	}
	ids, err = l.QueryChildren().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Children",
	}
	ids, err = l.QueryFiles().
		Select(file.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "File",
		Name: "Files",
	}
	ids, err = l.QueryHyperlinks().
		Select(hyperlink.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "Hyperlink",
		Name: "Hyperlinks",
	}
	ids, err = l.QueryEquipment().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[5] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "Equipment",
	}
	ids, err = l.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[6] = &Edge{
		IDs:  ids,
		Type: "Property",
		Name: "Properties",
	}
	ids, err = l.QuerySurvey().
		Select(survey.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[7] = &Edge{
		IDs:  ids,
		Type: "Survey",
		Name: "Survey",
	}
	ids, err = l.QueryWifiScan().
		Select(surveywifiscan.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[8] = &Edge{
		IDs:  ids,
		Type: "SurveyWiFiScan",
		Name: "WifiScan",
	}
	ids, err = l.QueryCellScan().
		Select(surveycellscan.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[9] = &Edge{
		IDs:  ids,
		Type: "SurveyCellScan",
		Name: "CellScan",
	}
	ids, err = l.QueryWorkOrders().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[10] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "WorkOrders",
	}
	ids, err = l.QueryFloorPlans().
		Select(floorplan.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[11] = &Edge{
		IDs:  ids,
		Type: "FloorPlan",
		Name: "FloorPlans",
	}
	return node, nil
}

func (lt *LocationType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     lt.ID,
		Type:   "LocationType",
		Fields: make([]*Field, 7),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(lt.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.Site); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "bool",
		Name:  "Site",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.Name); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.MapType); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "MapType",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.MapZoomLevel); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "int",
		Name:  "MapZoomLevel",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.Index); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "int",
		Name:  "Index",
		Value: string(buf),
	}
	var ids []int
	ids, err = lt.QueryLocations().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Locations",
	}
	ids, err = lt.QueryPropertyTypes().
		Select(propertytype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "PropertyType",
		Name: "PropertyTypes",
	}
	ids, err = lt.QuerySurveyTemplateCategories().
		Select(surveytemplatecategory.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "SurveyTemplateCategory",
		Name: "SurveyTemplateCategories",
	}
	return node, nil
}

func (pr *Project) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pr.ID,
		Type:   "Project",
		Fields: make([]*Field, 5),
		Edges:  make([]*Edge, 5),
	}
	var buf []byte
	if buf, err = json.Marshal(pr.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Description",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.Creator); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "Creator",
		Value: string(buf),
	}
	var ids []int
	ids, err = pr.QueryType().
		Select(projecttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "ProjectType",
		Name: "Type",
	}
	ids, err = pr.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Location",
	}
	ids, err = pr.QueryComments().
		Select(comment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Comment",
		Name: "Comments",
	}
	ids, err = pr.QueryWorkOrders().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "WorkOrders",
	}
	ids, err = pr.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "Property",
		Name: "Properties",
	}
	return node, nil
}

func (pt *ProjectType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pt.ID,
		Type:   "ProjectType",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(pt.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Description",
		Value: string(buf),
	}
	var ids []int
	ids, err = pt.QueryProjects().
		Select(project.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "Project",
		Name: "Projects",
	}
	ids, err = pt.QueryProperties().
		Select(propertytype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "PropertyType",
		Name: "Properties",
	}
	ids, err = pt.QueryWorkOrders().
		Select(workorderdefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "WorkOrderDefinition",
		Name: "WorkOrders",
	}
	return node, nil
}

func (pr *Property) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pr.ID,
		Type:   "Property",
		Fields: make([]*Field, 10),
		Edges:  make([]*Edge, 11),
	}
	var buf []byte
	if buf, err = json.Marshal(pr.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.IntVal); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "IntVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.BoolVal); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "bool",
		Name:  "BoolVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.FloatVal); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "float64",
		Name:  "FloatVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.LatitudeVal); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "float64",
		Name:  "LatitudeVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.LongitudeVal); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "float64",
		Name:  "LongitudeVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.RangeFromVal); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "float64",
		Name:  "RangeFromVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.RangeToVal); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "float64",
		Name:  "RangeToVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.StringVal); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "StringVal",
		Value: string(buf),
	}
	var ids []int
	ids, err = pr.QueryType().
		Select(propertytype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "PropertyType",
		Name: "Type",
	}
	ids, err = pr.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Location",
	}
	ids, err = pr.QueryEquipment().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "Equipment",
	}
	ids, err = pr.QueryService().
		Select(service.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "Service",
		Name: "Service",
	}
	ids, err = pr.QueryEquipmentPort().
		Select(equipmentport.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "EquipmentPort",
		Name: "EquipmentPort",
	}
	ids, err = pr.QueryLink().
		Select(link.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[5] = &Edge{
		IDs:  ids,
		Type: "Link",
		Name: "Link",
	}
	ids, err = pr.QueryWorkOrder().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[6] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "WorkOrder",
	}
	ids, err = pr.QueryProject().
		Select(project.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[7] = &Edge{
		IDs:  ids,
		Type: "Project",
		Name: "Project",
	}
	ids, err = pr.QueryEquipmentValue().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[8] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "EquipmentValue",
	}
	ids, err = pr.QueryLocationValue().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[9] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "LocationValue",
	}
	ids, err = pr.QueryServiceValue().
		Select(service.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[10] = &Edge{
		IDs:  ids,
		Type: "Service",
		Name: "ServiceValue",
	}
	return node, nil
}

func (pt *PropertyType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pt.ID,
		Type:   "PropertyType",
		Fields: make([]*Field, 18),
		Edges:  make([]*Edge, 8),
	}
	var buf []byte
	if buf, err = json.Marshal(pt.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Type); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Name); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Index); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "Index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Category); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "Category",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.IntVal); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "int",
		Name:  "IntVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.BoolVal); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "bool",
		Name:  "BoolVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.FloatVal); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "float64",
		Name:  "FloatVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.LatitudeVal); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "float64",
		Name:  "LatitudeVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.LongitudeVal); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "float64",
		Name:  "LongitudeVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.StringVal); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "string",
		Name:  "StringVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.RangeFromVal); err != nil {
		return nil, err
	}
	node.Fields[12] = &Field{
		Type:  "float64",
		Name:  "RangeFromVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.RangeToVal); err != nil {
		return nil, err
	}
	node.Fields[13] = &Field{
		Type:  "float64",
		Name:  "RangeToVal",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.IsInstanceProperty); err != nil {
		return nil, err
	}
	node.Fields[14] = &Field{
		Type:  "bool",
		Name:  "IsInstanceProperty",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Editable); err != nil {
		return nil, err
	}
	node.Fields[15] = &Field{
		Type:  "bool",
		Name:  "Editable",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Mandatory); err != nil {
		return nil, err
	}
	node.Fields[16] = &Field{
		Type:  "bool",
		Name:  "Mandatory",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Deleted); err != nil {
		return nil, err
	}
	node.Fields[17] = &Field{
		Type:  "bool",
		Name:  "Deleted",
		Value: string(buf),
	}
	var ids []int
	ids, err = pt.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "Property",
		Name: "Properties",
	}
	ids, err = pt.QueryLocationType().
		Select(locationtype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "LocationType",
		Name: "LocationType",
	}
	ids, err = pt.QueryEquipmentPortType().
		Select(equipmentporttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "EquipmentPortType",
		Name: "EquipmentPortType",
	}
	ids, err = pt.QueryLinkEquipmentPortType().
		Select(equipmentporttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "EquipmentPortType",
		Name: "LinkEquipmentPortType",
	}
	ids, err = pt.QueryEquipmentType().
		Select(equipmenttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "EquipmentType",
		Name: "EquipmentType",
	}
	ids, err = pt.QueryServiceType().
		Select(servicetype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[5] = &Edge{
		IDs:  ids,
		Type: "ServiceType",
		Name: "ServiceType",
	}
	ids, err = pt.QueryWorkOrderType().
		Select(workordertype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[6] = &Edge{
		IDs:  ids,
		Type: "WorkOrderType",
		Name: "WorkOrderType",
	}
	ids, err = pt.QueryProjectType().
		Select(projecttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[7] = &Edge{
		IDs:  ids,
		Type: "ProjectType",
		Name: "ProjectType",
	}
	return node, nil
}

func (rf *ReportFilter) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     rf.ID,
		Type:   "ReportFilter",
		Fields: make([]*Field, 5),
		Edges:  make([]*Edge, 0),
	}
	var buf []byte
	if buf, err = json.Marshal(rf.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rf.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rf.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rf.Entity); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "reportfilter.Entity",
		Name:  "Entity",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rf.Filters); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "Filters",
		Value: string(buf),
	}
	return node, nil
}

func (s *Service) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     s.ID,
		Type:   "Service",
		Fields: make([]*Field, 5),
		Edges:  make([]*Edge, 7),
	}
	var buf []byte
	if buf, err = json.Marshal(s.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.ExternalID); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "ExternalID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.Status); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "Status",
		Value: string(buf),
	}
	var ids []int
	ids, err = s.QueryType().
		Select(servicetype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "ServiceType",
		Name: "Type",
	}
	ids, err = s.QueryDownstream().
		Select(service.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Service",
		Name: "Downstream",
	}
	ids, err = s.QueryUpstream().
		Select(service.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Service",
		Name: "Upstream",
	}
	ids, err = s.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "Property",
		Name: "Properties",
	}
	ids, err = s.QueryLinks().
		Select(link.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "Link",
		Name: "Links",
	}
	ids, err = s.QueryCustomer().
		Select(customer.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[5] = &Edge{
		IDs:  ids,
		Type: "Customer",
		Name: "Customer",
	}
	ids, err = s.QueryEndpoints().
		Select(serviceendpoint.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[6] = &Edge{
		IDs:  ids,
		Type: "ServiceEndpoint",
		Name: "Endpoints",
	}
	return node, nil
}

func (se *ServiceEndpoint) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     se.ID,
		Type:   "ServiceEndpoint",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(se.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(se.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(se.Role); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Role",
		Value: string(buf),
	}
	var ids []int
	ids, err = se.QueryPort().
		Select(equipmentport.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "EquipmentPort",
		Name: "Port",
	}
	ids, err = se.QueryService().
		Select(service.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Service",
		Name: "Service",
	}
	return node, nil
}

func (st *ServiceType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     st.ID,
		Type:   "ServiceType",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(st.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(st.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(st.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(st.HasCustomer); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "bool",
		Name:  "HasCustomer",
		Value: string(buf),
	}
	var ids []int
	ids, err = st.QueryServices().
		Select(service.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "Service",
		Name: "Services",
	}
	ids, err = st.QueryPropertyTypes().
		Select(propertytype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "PropertyType",
		Name: "PropertyTypes",
	}
	return node, nil
}

func (s *Survey) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     s.ID,
		Type:   "Survey",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(s.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.OwnerName); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "OwnerName",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.CreationTimestamp); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "time.Time",
		Name:  "CreationTimestamp",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.CompletionTimestamp); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "time.Time",
		Name:  "CompletionTimestamp",
		Value: string(buf),
	}
	var ids []int
	ids, err = s.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Location",
	}
	ids, err = s.QuerySourceFile().
		Select(file.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "File",
		Name: "SourceFile",
	}
	ids, err = s.QueryQuestions().
		Select(surveyquestion.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "SurveyQuestion",
		Name: "Questions",
	}
	return node, nil
}

func (scs *SurveyCellScan) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     scs.ID,
		Type:   "SurveyCellScan",
		Fields: make([]*Field, 22),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(scs.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.NetworkType); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "NetworkType",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.SignalStrength); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "SignalStrength",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Timestamp); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "time.Time",
		Name:  "Timestamp",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.BaseStationID); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "BaseStationID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.NetworkID); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "string",
		Name:  "NetworkID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.SystemID); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "SystemID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.CellID); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "string",
		Name:  "CellID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.LocationAreaCode); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "LocationAreaCode",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.MobileCountryCode); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "string",
		Name:  "MobileCountryCode",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.MobileNetworkCode); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "string",
		Name:  "MobileNetworkCode",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.PrimaryScramblingCode); err != nil {
		return nil, err
	}
	node.Fields[12] = &Field{
		Type:  "string",
		Name:  "PrimaryScramblingCode",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Operator); err != nil {
		return nil, err
	}
	node.Fields[13] = &Field{
		Type:  "string",
		Name:  "Operator",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Arfcn); err != nil {
		return nil, err
	}
	node.Fields[14] = &Field{
		Type:  "int",
		Name:  "Arfcn",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.PhysicalCellID); err != nil {
		return nil, err
	}
	node.Fields[15] = &Field{
		Type:  "string",
		Name:  "PhysicalCellID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.TrackingAreaCode); err != nil {
		return nil, err
	}
	node.Fields[16] = &Field{
		Type:  "string",
		Name:  "TrackingAreaCode",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.TimingAdvance); err != nil {
		return nil, err
	}
	node.Fields[17] = &Field{
		Type:  "int",
		Name:  "TimingAdvance",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Earfcn); err != nil {
		return nil, err
	}
	node.Fields[18] = &Field{
		Type:  "int",
		Name:  "Earfcn",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Uarfcn); err != nil {
		return nil, err
	}
	node.Fields[19] = &Field{
		Type:  "int",
		Name:  "Uarfcn",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Latitude); err != nil {
		return nil, err
	}
	node.Fields[20] = &Field{
		Type:  "float64",
		Name:  "Latitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Longitude); err != nil {
		return nil, err
	}
	node.Fields[21] = &Field{
		Type:  "float64",
		Name:  "Longitude",
		Value: string(buf),
	}
	var ids []int
	ids, err = scs.QuerySurveyQuestion().
		Select(surveyquestion.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "SurveyQuestion",
		Name: "SurveyQuestion",
	}
	ids, err = scs.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Location",
	}
	return node, nil
}

func (sq *SurveyQuestion) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     sq.ID,
		Type:   "SurveyQuestion",
		Fields: make([]*Field, 20),
		Edges:  make([]*Edge, 4),
	}
	var buf []byte
	if buf, err = json.Marshal(sq.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.FormName); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "FormName",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.FormDescription); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "FormDescription",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.FormIndex); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "FormIndex",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.QuestionType); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "QuestionType",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.QuestionFormat); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "string",
		Name:  "QuestionFormat",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.QuestionText); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "QuestionText",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.QuestionIndex); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "int",
		Name:  "QuestionIndex",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.BoolData); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "bool",
		Name:  "BoolData",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.EmailData); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "string",
		Name:  "EmailData",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.Latitude); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "float64",
		Name:  "Latitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.Longitude); err != nil {
		return nil, err
	}
	node.Fields[12] = &Field{
		Type:  "float64",
		Name:  "Longitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.LocationAccuracy); err != nil {
		return nil, err
	}
	node.Fields[13] = &Field{
		Type:  "float64",
		Name:  "LocationAccuracy",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.Altitude); err != nil {
		return nil, err
	}
	node.Fields[14] = &Field{
		Type:  "float64",
		Name:  "Altitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.PhoneData); err != nil {
		return nil, err
	}
	node.Fields[15] = &Field{
		Type:  "string",
		Name:  "PhoneData",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.TextData); err != nil {
		return nil, err
	}
	node.Fields[16] = &Field{
		Type:  "string",
		Name:  "TextData",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.FloatData); err != nil {
		return nil, err
	}
	node.Fields[17] = &Field{
		Type:  "float64",
		Name:  "FloatData",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.IntData); err != nil {
		return nil, err
	}
	node.Fields[18] = &Field{
		Type:  "int",
		Name:  "IntData",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.DateData); err != nil {
		return nil, err
	}
	node.Fields[19] = &Field{
		Type:  "time.Time",
		Name:  "DateData",
		Value: string(buf),
	}
	var ids []int
	ids, err = sq.QuerySurvey().
		Select(survey.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "Survey",
		Name: "Survey",
	}
	ids, err = sq.QueryWifiScan().
		Select(surveywifiscan.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "SurveyWiFiScan",
		Name: "WifiScan",
	}
	ids, err = sq.QueryCellScan().
		Select(surveycellscan.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "SurveyCellScan",
		Name: "CellScan",
	}
	ids, err = sq.QueryPhotoData().
		Select(file.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "File",
		Name: "PhotoData",
	}
	return node, nil
}

func (stc *SurveyTemplateCategory) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     stc.ID,
		Type:   "SurveyTemplateCategory",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(stc.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stc.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stc.CategoryTitle); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "CategoryTitle",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stc.CategoryDescription); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "CategoryDescription",
		Value: string(buf),
	}
	var ids []int
	ids, err = stc.QuerySurveyTemplateQuestions().
		Select(surveytemplatequestion.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "SurveyTemplateQuestion",
		Name: "SurveyTemplateQuestions",
	}
	return node, nil
}

func (stq *SurveyTemplateQuestion) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     stq.ID,
		Type:   "SurveyTemplateQuestion",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(stq.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stq.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stq.QuestionTitle); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "QuestionTitle",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stq.QuestionDescription); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "QuestionDescription",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stq.QuestionType); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "QuestionType",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stq.Index); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "int",
		Name:  "Index",
		Value: string(buf),
	}
	var ids []int
	ids, err = stq.QueryCategory().
		Select(surveytemplatecategory.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "SurveyTemplateCategory",
		Name: "Category",
	}
	return node, nil
}

func (swfs *SurveyWiFiScan) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     swfs.ID,
		Type:   "SurveyWiFiScan",
		Fields: make([]*Field, 13),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(swfs.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Ssid); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Ssid",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Bssid); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Bssid",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Timestamp); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "time.Time",
		Name:  "Timestamp",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Frequency); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "int",
		Name:  "Frequency",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Channel); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "int",
		Name:  "Channel",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Band); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "Band",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.ChannelWidth); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "int",
		Name:  "ChannelWidth",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Capabilities); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "Capabilities",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Strength); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "int",
		Name:  "Strength",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Latitude); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "float64",
		Name:  "Latitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Longitude); err != nil {
		return nil, err
	}
	node.Fields[12] = &Field{
		Type:  "float64",
		Name:  "Longitude",
		Value: string(buf),
	}
	var ids []int
	ids, err = swfs.QuerySurveyQuestion().
		Select(surveyquestion.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "SurveyQuestion",
		Name: "SurveyQuestion",
	}
	ids, err = swfs.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Location",
	}
	return node, nil
}

func (t *Technician) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     t.ID,
		Type:   "Technician",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(t.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(t.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(t.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(t.Email); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Email",
		Value: string(buf),
	}
	var ids []int
	ids, err = t.QueryWorkOrders().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "WorkOrders",
	}
	return node, nil
}

func (u *User) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     u.ID,
		Type:   "User",
		Fields: make([]*Field, 8),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(u.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.AuthID); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "AuthID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.FirstName); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "FirstName",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.LastName); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "LastName",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.Email); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "Email",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.Status); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "user.Status",
		Name:  "Status",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.Role); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "user.Role",
		Name:  "Role",
		Value: string(buf),
	}
	var ids []int
	ids, err = u.QueryProfilePhoto().
		Select(file.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "File",
		Name: "ProfilePhoto",
	}
	return node, nil
}

func (wo *WorkOrder) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     wo.ID,
		Type:   "WorkOrder",
		Fields: make([]*Field, 12),
		Edges:  make([]*Edge, 12),
	}
	var buf []byte
	if buf, err = json.Marshal(wo.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Status); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Status",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Priority); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "Priority",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Description); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "Description",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.OwnerName); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "string",
		Name:  "OwnerName",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.InstallDate); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "time.Time",
		Name:  "InstallDate",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.CreationDate); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "time.Time",
		Name:  "CreationDate",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Assignee); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "Assignee",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Index); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "int",
		Name:  "Index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.CloseDate); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "time.Time",
		Name:  "CloseDate",
		Value: string(buf),
	}
	var ids []int
	ids, err = wo.QueryType().
		Select(workordertype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "WorkOrderType",
		Name: "Type",
	}
	ids, err = wo.QueryEquipment().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "Equipment",
	}
	ids, err = wo.QueryLinks().
		Select(link.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Link",
		Name: "Links",
	}
	ids, err = wo.QueryFiles().
		Select(file.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "File",
		Name: "Files",
	}
	ids, err = wo.QueryHyperlinks().
		Select(hyperlink.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "Hyperlink",
		Name: "Hyperlinks",
	}
	ids, err = wo.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[5] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "Location",
	}
	ids, err = wo.QueryComments().
		Select(comment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[6] = &Edge{
		IDs:  ids,
		Type: "Comment",
		Name: "Comments",
	}
	ids, err = wo.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[7] = &Edge{
		IDs:  ids,
		Type: "Property",
		Name: "Properties",
	}
	ids, err = wo.QueryCheckListCategories().
		Select(checklistcategory.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[8] = &Edge{
		IDs:  ids,
		Type: "CheckListCategory",
		Name: "CheckListCategories",
	}
	ids, err = wo.QueryCheckListItems().
		Select(checklistitem.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[9] = &Edge{
		IDs:  ids,
		Type: "CheckListItem",
		Name: "CheckListItems",
	}
	ids, err = wo.QueryTechnician().
		Select(technician.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[10] = &Edge{
		IDs:  ids,
		Type: "Technician",
		Name: "Technician",
	}
	ids, err = wo.QueryProject().
		Select(project.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[11] = &Edge{
		IDs:  ids,
		Type: "Project",
		Name: "Project",
	}
	return node, nil
}

func (wod *WorkOrderDefinition) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     wod.ID,
		Type:   "WorkOrderDefinition",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(wod.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wod.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wod.Index); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "Index",
		Value: string(buf),
	}
	var ids []int
	ids, err = wod.QueryType().
		Select(workordertype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "WorkOrderType",
		Name: "Type",
	}
	ids, err = wod.QueryProjectType().
		Select(projecttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "ProjectType",
		Name: "ProjectType",
	}
	return node, nil
}

func (wot *WorkOrderType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     wot.ID,
		Type:   "WorkOrderType",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 5),
	}
	var buf []byte
	if buf, err = json.Marshal(wot.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "CreateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wot.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "UpdateTime",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wot.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "Name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wot.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "Description",
		Value: string(buf),
	}
	var ids []int
	ids, err = wot.QueryWorkOrders().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "WorkOrders",
	}
	ids, err = wot.QueryPropertyTypes().
		Select(propertytype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "PropertyType",
		Name: "PropertyTypes",
	}
	ids, err = wot.QueryDefinitions().
		Select(workorderdefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "WorkOrderDefinition",
		Name: "Definitions",
	}
	ids, err = wot.QueryCheckListCategories().
		Select(checklistcategory.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "CheckListCategory",
		Name: "CheckListCategories",
	}
	ids, err = wot.QueryCheckListDefinitions().
		Select(checklistitemdefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "CheckListItemDefinition",
		Name: "CheckListDefinitions",
	}
	return node, nil
}

func (c *Client) Node(ctx context.Context, id int) (*Node, error) {
	n, err := c.Noder(ctx, id)
	if err != nil {
		return nil, err
	}
	return n.Node(ctx)
}

func (c *Client) Noder(ctx context.Context, id int) (Noder, error) {
	tables, err := c.tables.Load(ctx, c.driver)
	if err != nil {
		return nil, err
	}
	idx := id / (1<<32 - 1)
	if idx < 0 || idx >= len(tables) {
		return nil, fmt.Errorf("cannot resolve table from id %v: %w", id, &NotFoundError{"invalid/unknown"})
	}
	return c.noder(ctx, tables[idx], id)
}

func (c *Client) noder(ctx context.Context, tbl string, id int) (Noder, error) {
	switch tbl {
	case actionsrule.Table:
		n, err := c.ActionsRule.Query().
			Where(actionsrule.ID(id)).
			CollectFields(ctx, "ActionsRule").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case checklistcategory.Table:
		n, err := c.CheckListCategory.Query().
			Where(checklistcategory.ID(id)).
			CollectFields(ctx, "CheckListCategory").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case checklistitem.Table:
		n, err := c.CheckListItem.Query().
			Where(checklistitem.ID(id)).
			CollectFields(ctx, "CheckListItem").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case checklistitemdefinition.Table:
		n, err := c.CheckListItemDefinition.Query().
			Where(checklistitemdefinition.ID(id)).
			CollectFields(ctx, "CheckListItemDefinition").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case comment.Table:
		n, err := c.Comment.Query().
			Where(comment.ID(id)).
			CollectFields(ctx, "Comment").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case customer.Table:
		n, err := c.Customer.Query().
			Where(customer.ID(id)).
			CollectFields(ctx, "Customer").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case equipment.Table:
		n, err := c.Equipment.Query().
			Where(equipment.ID(id)).
			CollectFields(ctx, "Equipment").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case equipmentcategory.Table:
		n, err := c.EquipmentCategory.Query().
			Where(equipmentcategory.ID(id)).
			CollectFields(ctx, "EquipmentCategory").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case equipmentport.Table:
		n, err := c.EquipmentPort.Query().
			Where(equipmentport.ID(id)).
			CollectFields(ctx, "EquipmentPort").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case equipmentportdefinition.Table:
		n, err := c.EquipmentPortDefinition.Query().
			Where(equipmentportdefinition.ID(id)).
			CollectFields(ctx, "EquipmentPortDefinition").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case equipmentporttype.Table:
		n, err := c.EquipmentPortType.Query().
			Where(equipmentporttype.ID(id)).
			CollectFields(ctx, "EquipmentPortType").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case equipmentposition.Table:
		n, err := c.EquipmentPosition.Query().
			Where(equipmentposition.ID(id)).
			CollectFields(ctx, "EquipmentPosition").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case equipmentpositiondefinition.Table:
		n, err := c.EquipmentPositionDefinition.Query().
			Where(equipmentpositiondefinition.ID(id)).
			CollectFields(ctx, "EquipmentPositionDefinition").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case equipmenttype.Table:
		n, err := c.EquipmentType.Query().
			Where(equipmenttype.ID(id)).
			CollectFields(ctx, "EquipmentType").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case file.Table:
		n, err := c.File.Query().
			Where(file.ID(id)).
			CollectFields(ctx, "File").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case floorplan.Table:
		n, err := c.FloorPlan.Query().
			Where(floorplan.ID(id)).
			CollectFields(ctx, "FloorPlan").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case floorplanreferencepoint.Table:
		n, err := c.FloorPlanReferencePoint.Query().
			Where(floorplanreferencepoint.ID(id)).
			CollectFields(ctx, "FloorPlanReferencePoint").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case floorplanscale.Table:
		n, err := c.FloorPlanScale.Query().
			Where(floorplanscale.ID(id)).
			CollectFields(ctx, "FloorPlanScale").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case hyperlink.Table:
		n, err := c.Hyperlink.Query().
			Where(hyperlink.ID(id)).
			CollectFields(ctx, "Hyperlink").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case link.Table:
		n, err := c.Link.Query().
			Where(link.ID(id)).
			CollectFields(ctx, "Link").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case location.Table:
		n, err := c.Location.Query().
			Where(location.ID(id)).
			CollectFields(ctx, "Location").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case locationtype.Table:
		n, err := c.LocationType.Query().
			Where(locationtype.ID(id)).
			CollectFields(ctx, "LocationType").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case project.Table:
		n, err := c.Project.Query().
			Where(project.ID(id)).
			CollectFields(ctx, "Project").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case projecttype.Table:
		n, err := c.ProjectType.Query().
			Where(projecttype.ID(id)).
			CollectFields(ctx, "ProjectType").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case property.Table:
		n, err := c.Property.Query().
			Where(property.ID(id)).
			CollectFields(ctx, "Property").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case propertytype.Table:
		n, err := c.PropertyType.Query().
			Where(propertytype.ID(id)).
			CollectFields(ctx, "PropertyType").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case reportfilter.Table:
		n, err := c.ReportFilter.Query().
			Where(reportfilter.ID(id)).
			CollectFields(ctx, "ReportFilter").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case service.Table:
		n, err := c.Service.Query().
			Where(service.ID(id)).
			CollectFields(ctx, "Service").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case serviceendpoint.Table:
		n, err := c.ServiceEndpoint.Query().
			Where(serviceendpoint.ID(id)).
			CollectFields(ctx, "ServiceEndpoint").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case servicetype.Table:
		n, err := c.ServiceType.Query().
			Where(servicetype.ID(id)).
			CollectFields(ctx, "ServiceType").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case survey.Table:
		n, err := c.Survey.Query().
			Where(survey.ID(id)).
			CollectFields(ctx, "Survey").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case surveycellscan.Table:
		n, err := c.SurveyCellScan.Query().
			Where(surveycellscan.ID(id)).
			CollectFields(ctx, "SurveyCellScan").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case surveyquestion.Table:
		n, err := c.SurveyQuestion.Query().
			Where(surveyquestion.ID(id)).
			CollectFields(ctx, "SurveyQuestion").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case surveytemplatecategory.Table:
		n, err := c.SurveyTemplateCategory.Query().
			Where(surveytemplatecategory.ID(id)).
			CollectFields(ctx, "SurveyTemplateCategory").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case surveytemplatequestion.Table:
		n, err := c.SurveyTemplateQuestion.Query().
			Where(surveytemplatequestion.ID(id)).
			CollectFields(ctx, "SurveyTemplateQuestion").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case surveywifiscan.Table:
		n, err := c.SurveyWiFiScan.Query().
			Where(surveywifiscan.ID(id)).
			CollectFields(ctx, "SurveyWiFiScan").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case technician.Table:
		n, err := c.Technician.Query().
			Where(technician.ID(id)).
			CollectFields(ctx, "Technician").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case user.Table:
		n, err := c.User.Query().
			Where(user.ID(id)).
			CollectFields(ctx, "User").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case workorder.Table:
		n, err := c.WorkOrder.Query().
			Where(workorder.ID(id)).
			CollectFields(ctx, "WorkOrder").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case workorderdefinition.Table:
		n, err := c.WorkOrderDefinition.Query().
			Where(workorderdefinition.ID(id)).
			CollectFields(ctx, "WorkOrderDefinition").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case workordertype.Table:
		n, err := c.WorkOrderType.Query().
			Where(workordertype.ID(id)).
			CollectFields(ctx, "WorkOrderType").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	default:
		return nil, fmt.Errorf("cannot resolve noder from table %q: %w", tbl, &NotFoundError{"invalid/unknown"})
	}
}

type (
	tables struct {
		once  sync.Once
		sem   *semaphore.Weighted
		value atomic.Value
	}

	querier interface {
		Query(ctx context.Context, query string, args, v interface{}) error
	}
)

func (t *tables) Load(ctx context.Context, querier querier) ([]string, error) {
	if tables := t.value.Load(); tables != nil {
		return tables.([]string), nil
	}
	t.once.Do(func() { t.sem = semaphore.NewWeighted(1) })
	if err := t.sem.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer t.sem.Release(1)
	if tables := t.value.Load(); tables != nil {
		return tables.([]string), nil
	}
	tables, err := t.load(ctx, querier)
	if err == nil {
		t.value.Store(tables)
	}
	return tables, err
}

func (tables) load(ctx context.Context, querier querier) ([]string, error) {
	rows := &sql.Rows{}
	query, args := sql.Select("type").
		From(sql.Table(schema.TypeTable)).
		OrderBy(sql.Asc("id")).
		Query()
	if err := querier.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	return tables, sql.ScanSlice(rows, &tables)
}
