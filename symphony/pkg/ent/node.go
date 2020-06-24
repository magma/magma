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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ar.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ar.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ar.TriggerID); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "triggerID",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ar.RuleFilters); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "[]*core.ActionsRuleFilter",
		Name:  "ruleFilters",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ar.RuleActions); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "[]*core.ActionsRuleAction",
		Name:  "ruleActions",
		Value: string(buf),
	}
	return node, nil
}

func (ar *ActionsRuleMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := ar.ID()
	if !exists {
		return nil, nil
	}
	ent, err := ar.Client().ActionsRule.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (a *Activity) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     a.ID,
		Type:   "Activity",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(a.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.ChangedField); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "activity.ChangedField",
		Name:  "changed_field",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.IsCreate); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "bool",
		Name:  "is_create",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.OldValue); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "old_value",
		Value: string(buf),
	}
	if buf, err = json.Marshal(a.NewValue); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "new_value",
		Value: string(buf),
	}
	var ids []int
	ids, err = a.QueryAuthor().
		Select(user.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "User",
		Name: "author",
	}
	ids, err = a.QueryWorkOrder().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "work_order",
	}
	return node, nil
}

func (a *ActivityMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := a.ID()
	if !exists {
		return nil, nil
	}
	ent, err := a.Client().Activity.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (clc *CheckListCategory) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     clc.ID,
		Type:   "CheckListCategory",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(clc.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clc.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clc.Title); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "title",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clc.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "description",
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
		Name: "check_list_items",
	}
	ids, err = clc.QueryWorkOrder().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "work_order",
	}
	return node, nil
}

func (clc *CheckListCategoryMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := clc.ID()
	if !exists {
		return nil, nil
	}
	ent, err := clc.Client().CheckListCategory.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (clcd *CheckListCategoryDefinition) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     clcd.ID,
		Type:   "CheckListCategoryDefinition",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(clcd.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clcd.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clcd.Title); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "title",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clcd.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "description",
		Value: string(buf),
	}
	var ids []int
	ids, err = clcd.QueryCheckListItemDefinitions().
		Select(checklistitemdefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "CheckListItemDefinition",
		Name: "check_list_item_definitions",
	}
	ids, err = clcd.QueryWorkOrderType().
		Select(workordertype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "WorkOrderType",
		Name: "work_order_type",
	}
	return node, nil
}

func (clcd *CheckListCategoryDefinitionMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := clcd.ID()
	if !exists {
		return nil, nil
	}
	ent, err := clcd.Client().CheckListCategoryDefinition.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (cli *CheckListItem) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     cli.ID,
		Type:   "CheckListItem",
		Fields: make([]*Field, 10),
		Edges:  make([]*Edge, 4),
	}
	var buf []byte
	if buf, err = json.Marshal(cli.Title); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "string",
		Name:  "title",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.Type); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "string",
		Name:  "type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.Index); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.Checked); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "bool",
		Name:  "checked",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.StringVal); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "string_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.EnumValues); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "enum_values",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.EnumSelectionModeValue); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "checklistitem.EnumSelectionModeValue",
		Name:  "enum_selection_mode_value",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.SelectedEnumValues); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "selected_enum_values",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.YesNoVal); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "checklistitem.YesNoVal",
		Name:  "yes_no_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(cli.HelpText); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "help_text",
		Value: string(buf),
	}
	var ids []int
	ids, err = cli.QueryFiles().
		Select(file.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "File",
		Name: "files",
	}
	ids, err = cli.QueryWifiScan().
		Select(surveywifiscan.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "SurveyWiFiScan",
		Name: "wifi_scan",
	}
	ids, err = cli.QueryCellScan().
		Select(surveycellscan.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "SurveyCellScan",
		Name: "cell_scan",
	}
	ids, err = cli.QueryCheckListCategory().
		Select(checklistcategory.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "CheckListCategory",
		Name: "check_list_category",
	}
	return node, nil
}

func (cli *CheckListItemMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := cli.ID()
	if !exists {
		return nil, nil
	}
	ent, err := cli.Client().CheckListItem.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (clid *CheckListItemDefinition) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     clid.ID,
		Type:   "CheckListItemDefinition",
		Fields: make([]*Field, 8),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(clid.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.Title); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "title",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.Type); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.Index); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.EnumValues); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "enum_values",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.EnumSelectionModeValue); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "checklistitemdefinition.EnumSelectionModeValue",
		Name:  "enum_selection_mode_value",
		Value: string(buf),
	}
	if buf, err = json.Marshal(clid.HelpText); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "help_text",
		Value: string(buf),
	}
	var ids []int
	ids, err = clid.QueryCheckListCategoryDefinition().
		Select(checklistcategorydefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "CheckListCategoryDefinition",
		Name: "check_list_category_definition",
	}
	return node, nil
}

func (clid *CheckListItemDefinitionMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := clid.ID()
	if !exists {
		return nil, nil
	}
	ent, err := clid.Client().CheckListItemDefinition.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (c *Comment) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     c.ID,
		Type:   "Comment",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(c.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.Text); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "text",
		Value: string(buf),
	}
	var ids []int
	ids, err = c.QueryAuthor().
		Select(user.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "User",
		Name: "author",
	}
	ids, err = c.QueryWorkOrder().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "work_order",
	}
	ids, err = c.QueryProject().
		Select(project.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Project",
		Name: "project",
	}
	return node, nil
}

func (c *CommentMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := c.ID()
	if !exists {
		return nil, nil
	}
	ent, err := c.Client().Comment.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(c.ExternalID); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "external_id",
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
		Name: "services",
	}
	return node, nil
}

func (c *CustomerMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := c.ID()
	if !exists {
		return nil, nil
	}
	ent, err := c.Client().Customer.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (e *Equipment) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     e.ID,
		Type:   "Equipment",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 10),
	}
	var buf []byte
	if buf, err = json.Marshal(e.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(e.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(e.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(e.FutureState); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "future_state",
		Value: string(buf),
	}
	if buf, err = json.Marshal(e.DeviceID); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "device_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(e.ExternalID); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "external_id",
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
		Name: "type",
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
		Name: "location",
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
		Name: "parent_position",
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
		Name: "positions",
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
		Name: "ports",
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
		Name: "work_order",
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
		Name: "properties",
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
		Name: "files",
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
		Name: "hyperlinks",
	}
	ids, err = e.QueryEndpoints().
		Select(serviceendpoint.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[9] = &Edge{
		IDs:  ids,
		Type: "ServiceEndpoint",
		Name: "endpoints",
	}
	return node, nil
}

func (e *EquipmentMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := e.ID()
	if !exists {
		return nil, nil
	}
	ent, err := e.Client().Equipment.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ec.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ec.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
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
		Name: "types",
	}
	return node, nil
}

func (ec *EquipmentCategoryMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := ec.ID()
	if !exists {
		return nil, nil
	}
	ent, err := ec.Client().EquipmentCategory.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ep.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
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
		Name: "definition",
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
		Name: "parent",
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
		Name: "link",
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
		Name: "properties",
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
		Name: "endpoints",
	}
	return node, nil
}

func (ep *EquipmentPortMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := ep.ID()
	if !exists {
		return nil, nil
	}
	ent, err := ep.Client().EquipmentPort.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.Index); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.Bandwidth); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "bandwidth",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.VisibilityLabel); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "visibility_label",
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
		Name: "equipment_port_type",
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
		Name: "ports",
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
		Name: "equipment_type",
	}
	return node, nil
}

func (epd *EquipmentPortDefinitionMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := epd.ID()
	if !exists {
		return nil, nil
	}
	ent, err := epd.Client().EquipmentPortDefinition.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ept.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ept.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
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
		Name: "property_types",
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
		Name: "link_property_types",
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
		Name: "port_definitions",
	}
	return node, nil
}

func (ept *EquipmentPortTypeMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := ept.ID()
	if !exists {
		return nil, nil
	}
	ent, err := ept.Client().EquipmentPortType.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ep.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
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
		Name: "definition",
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
		Name: "parent",
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
		Name: "attachment",
	}
	return node, nil
}

func (ep *EquipmentPositionMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := ep.ID()
	if !exists {
		return nil, nil
	}
	ent, err := ep.Client().EquipmentPosition.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.Index); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(epd.VisibilityLabel); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "visibility_label",
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
		Name: "positions",
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
		Name: "equipment_type",
	}
	return node, nil
}

func (epd *EquipmentPositionDefinitionMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := epd.ID()
	if !exists {
		return nil, nil
	}
	ent, err := epd.Client().EquipmentPositionDefinition.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (et *EquipmentType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     et.ID,
		Type:   "EquipmentType",
		Fields: make([]*Field, 3),
		Edges:  make([]*Edge, 6),
	}
	var buf []byte
	if buf, err = json.Marshal(et.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(et.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(et.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
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
		Name: "port_definitions",
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
		Name: "position_definitions",
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
		Name: "property_types",
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
		Name: "equipment",
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
		Name: "category",
	}
	ids, err = et.QueryServiceEndpointDefinitions().
		Select(serviceendpointdefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[5] = &Edge{
		IDs:  ids,
		Type: "ServiceEndpointDefinition",
		Name: "service_endpoint_definitions",
	}
	return node, nil
}

func (et *EquipmentTypeMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := et.ID()
	if !exists {
		return nil, nil
	}
	ent, err := et.Client().EquipmentType.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (f *File) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     f.ID,
		Type:   "File",
		Fields: make([]*Field, 11),
		Edges:  make([]*Edge, 9),
	}
	var buf []byte
	if buf, err = json.Marshal(f.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Type); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Name); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Size); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "size",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.ModifiedAt); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "time.Time",
		Name:  "modified_at",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.UploadedAt); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "time.Time",
		Name:  "uploaded_at",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.ContentType); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "content_type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.StoreKey); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "string",
		Name:  "store_key",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Category); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "category",
		Value: string(buf),
	}
	if buf, err = json.Marshal(f.Annotation); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "string",
		Name:  "annotation",
		Value: string(buf),
	}
	var ids []int
	ids, err = f.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "location",
	}
	ids, err = f.QueryEquipment().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "equipment",
	}
	ids, err = f.QueryUser().
		Select(user.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "User",
		Name: "user",
	}
	ids, err = f.QueryWorkOrder().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "work_order",
	}
	ids, err = f.QueryChecklistItem().
		Select(checklistitem.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "CheckListItem",
		Name: "checklist_item",
	}
	ids, err = f.QuerySurvey().
		Select(survey.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[5] = &Edge{
		IDs:  ids,
		Type: "Survey",
		Name: "survey",
	}
	ids, err = f.QueryFloorPlan().
		Select(floorplan.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[6] = &Edge{
		IDs:  ids,
		Type: "FloorPlan",
		Name: "floor_plan",
	}
	ids, err = f.QueryPhotoSurveyQuestion().
		Select(surveyquestion.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[7] = &Edge{
		IDs:  ids,
		Type: "SurveyQuestion",
		Name: "photo_survey_question",
	}
	ids, err = f.QuerySurveyQuestion().
		Select(surveyquestion.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[8] = &Edge{
		IDs:  ids,
		Type: "SurveyQuestion",
		Name: "survey_question",
	}
	return node, nil
}

func (f *FileMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := f.ID()
	if !exists {
		return nil, nil
	}
	ent, err := f.Client().File.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fp.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fp.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
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
		Name: "location",
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
		Name: "reference_point",
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
		Name: "scale",
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
		Name: "image",
	}
	return node, nil
}

func (fp *FloorPlanMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := fp.ID()
	if !exists {
		return nil, nil
	}
	ent, err := fp.Client().FloorPlan.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fprp.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fprp.X); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "x",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fprp.Y); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "y",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fprp.Latitude); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "float64",
		Name:  "latitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fprp.Longitude); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "float64",
		Name:  "longitude",
		Value: string(buf),
	}
	return node, nil
}

func (fprp *FloorPlanReferencePointMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := fprp.ID()
	if !exists {
		return nil, nil
	}
	ent, err := fprp.Client().FloorPlanReferencePoint.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.ReferencePoint1X); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "reference_point1_x",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.ReferencePoint1Y); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "reference_point1_y",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.ReferencePoint2X); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "reference_point2_x",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.ReferencePoint2Y); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "int",
		Name:  "reference_point2_y",
		Value: string(buf),
	}
	if buf, err = json.Marshal(fps.ScaleInMeters); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "float64",
		Name:  "scale_in_meters",
		Value: string(buf),
	}
	return node, nil
}

func (fps *FloorPlanScaleMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := fps.ID()
	if !exists {
		return nil, nil
	}
	ent, err := fps.Client().FloorPlanScale.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (h *Hyperlink) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     h.ID,
		Type:   "Hyperlink",
		Fields: make([]*Field, 5),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(h.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(h.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(h.URL); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "url",
		Value: string(buf),
	}
	if buf, err = json.Marshal(h.Name); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(h.Category); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "category",
		Value: string(buf),
	}
	var ids []int
	ids, err = h.QueryEquipment().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "equipment",
	}
	ids, err = h.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "location",
	}
	ids, err = h.QueryWorkOrder().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "work_order",
	}
	return node, nil
}

func (h *HyperlinkMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := h.ID()
	if !exists {
		return nil, nil
	}
	ent, err := h.Client().Hyperlink.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.FutureState); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "future_state",
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
		Name: "ports",
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
		Name: "work_order",
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
		Name: "properties",
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
		Name: "service",
	}
	return node, nil
}

func (l *LinkMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := l.ID()
	if !exists {
		return nil, nil
	}
	ent, err := l.Client().Link.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.ExternalID); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "external_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.Latitude); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "float64",
		Name:  "latitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.Longitude); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "float64",
		Name:  "longitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(l.SiteSurveyNeeded); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "bool",
		Name:  "site_survey_needed",
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
		Name: "type",
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
		Name: "parent",
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
		Name: "children",
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
		Name: "files",
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
		Name: "hyperlinks",
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
		Name: "equipment",
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
		Name: "properties",
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
		Name: "survey",
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
		Name: "wifi_scan",
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
		Name: "cell_scan",
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
		Name: "work_orders",
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
		Name: "floor_plans",
	}
	return node, nil
}

func (l *LocationMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := l.ID()
	if !exists {
		return nil, nil
	}
	ent, err := l.Client().Location.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.Site); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "bool",
		Name:  "site",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.Name); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.MapType); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "map_type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.MapZoomLevel); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "int",
		Name:  "map_zoom_level",
		Value: string(buf),
	}
	if buf, err = json.Marshal(lt.Index); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "int",
		Name:  "index",
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
		Name: "locations",
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
		Name: "property_types",
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
		Name: "survey_template_categories",
	}
	return node, nil
}

func (lt *LocationTypeMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := lt.ID()
	if !exists {
		return nil, nil
	}
	ent, err := lt.Client().LocationType.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (pp *PermissionsPolicy) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pp.ID,
		Type:   "PermissionsPolicy",
		Fields: make([]*Field, 7),
		Edges:  make([]*Edge, 1),
	}
	var buf []byte
	if buf, err = json.Marshal(pp.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pp.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pp.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pp.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "description",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pp.IsGlobal); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "bool",
		Name:  "is_global",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pp.InventoryPolicy); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "*models.InventoryPolicyInput",
		Name:  "inventory_policy",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pp.WorkforcePolicy); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "*models.WorkforcePolicyInput",
		Name:  "workforce_policy",
		Value: string(buf),
	}
	var ids []int
	ids, err = pp.QueryGroups().
		Select(usersgroup.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "UsersGroup",
		Name: "groups",
	}
	return node, nil
}

func (pp *PermissionsPolicyMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := pp.ID()
	if !exists {
		return nil, nil
	}
	ent, err := pp.Client().PermissionsPolicy.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (pr *Project) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pr.ID,
		Type:   "Project",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 6),
	}
	var buf []byte
	if buf, err = json.Marshal(pr.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "description",
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
		Name: "type",
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
		Name: "location",
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
		Name: "comments",
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
		Name: "work_orders",
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
		Name: "properties",
	}
	ids, err = pr.QueryCreator().
		Select(user.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[5] = &Edge{
		IDs:  ids,
		Type: "User",
		Name: "creator",
	}
	return node, nil
}

func (pr *ProjectMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := pr.ID()
	if !exists {
		return nil, nil
	}
	ent, err := pr.Client().Project.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "description",
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
		Name: "projects",
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
		Name: "properties",
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
		Name: "work_orders",
	}
	return node, nil
}

func (pt *ProjectTypeMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := pt.ID()
	if !exists {
		return nil, nil
	}
	ent, err := pt.Client().ProjectType.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (pr *Property) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pr.ID,
		Type:   "Property",
		Fields: make([]*Field, 10),
		Edges:  make([]*Edge, 13),
	}
	var buf []byte
	if buf, err = json.Marshal(pr.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.IntVal); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "int_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.BoolVal); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "bool",
		Name:  "bool_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.FloatVal); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "float64",
		Name:  "float_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.LatitudeVal); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "float64",
		Name:  "latitude_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.LongitudeVal); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "float64",
		Name:  "longitude_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.RangeFromVal); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "float64",
		Name:  "range_from_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.RangeToVal); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "float64",
		Name:  "range_to_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pr.StringVal); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "string_val",
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
		Name: "type",
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
		Name: "location",
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
		Name: "equipment",
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
		Name: "service",
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
		Name: "equipment_port",
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
		Name: "link",
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
		Name: "work_order",
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
		Name: "project",
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
		Name: "equipment_value",
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
		Name: "location_value",
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
		Name: "service_value",
	}
	ids, err = pr.QueryWorkOrderValue().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[11] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "work_order_value",
	}
	ids, err = pr.QueryUserValue().
		Select(user.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[12] = &Edge{
		IDs:  ids,
		Type: "User",
		Name: "user_value",
	}
	return node, nil
}

func (pr *PropertyMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := pr.ID()
	if !exists {
		return nil, nil
	}
	ent, err := pr.Client().Property.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (pt *PropertyType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     pt.ID,
		Type:   "PropertyType",
		Fields: make([]*Field, 20),
		Edges:  make([]*Edge, 8),
	}
	var buf []byte
	if buf, err = json.Marshal(pt.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Type); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Name); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.ExternalID); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "external_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Index); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "int",
		Name:  "index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Category); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "string",
		Name:  "category",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.IntVal); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "int",
		Name:  "int_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.BoolVal); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "bool",
		Name:  "bool_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.FloatVal); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "float64",
		Name:  "float_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.LatitudeVal); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "float64",
		Name:  "latitude_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.LongitudeVal); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "float64",
		Name:  "longitude_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.StringVal); err != nil {
		return nil, err
	}
	node.Fields[12] = &Field{
		Type:  "string",
		Name:  "string_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.RangeFromVal); err != nil {
		return nil, err
	}
	node.Fields[13] = &Field{
		Type:  "float64",
		Name:  "range_from_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.RangeToVal); err != nil {
		return nil, err
	}
	node.Fields[14] = &Field{
		Type:  "float64",
		Name:  "range_to_val",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.IsInstanceProperty); err != nil {
		return nil, err
	}
	node.Fields[15] = &Field{
		Type:  "bool",
		Name:  "is_instance_property",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Editable); err != nil {
		return nil, err
	}
	node.Fields[16] = &Field{
		Type:  "bool",
		Name:  "editable",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Mandatory); err != nil {
		return nil, err
	}
	node.Fields[17] = &Field{
		Type:  "bool",
		Name:  "mandatory",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.Deleted); err != nil {
		return nil, err
	}
	node.Fields[18] = &Field{
		Type:  "bool",
		Name:  "deleted",
		Value: string(buf),
	}
	if buf, err = json.Marshal(pt.NodeType); err != nil {
		return nil, err
	}
	node.Fields[19] = &Field{
		Type:  "string",
		Name:  "nodeType",
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
		Name: "properties",
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
		Name: "location_type",
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
		Name: "equipment_port_type",
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
		Name: "link_equipment_port_type",
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
		Name: "equipment_type",
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
		Name: "service_type",
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
		Name: "work_order_type",
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
		Name: "project_type",
	}
	return node, nil
}

func (pt *PropertyTypeMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := pt.ID()
	if !exists {
		return nil, nil
	}
	ent, err := pt.Client().PropertyType.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rf.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rf.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rf.Entity); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "reportfilter.Entity",
		Name:  "entity",
		Value: string(buf),
	}
	if buf, err = json.Marshal(rf.Filters); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "filters",
		Value: string(buf),
	}
	return node, nil
}

func (rf *ReportFilterMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := rf.ID()
	if !exists {
		return nil, nil
	}
	ent, err := rf.Client().ReportFilter.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.ExternalID); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "external_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.Status); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "status",
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
		Name: "type",
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
		Name: "downstream",
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
		Name: "upstream",
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
		Name: "properties",
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
		Name: "links",
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
		Name: "customer",
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
		Name: "endpoints",
	}
	return node, nil
}

func (s *ServiceMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := s.ID()
	if !exists {
		return nil, nil
	}
	ent, err := s.Client().Service.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (se *ServiceEndpoint) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     se.ID,
		Type:   "ServiceEndpoint",
		Fields: make([]*Field, 2),
		Edges:  make([]*Edge, 4),
	}
	var buf []byte
	if buf, err = json.Marshal(se.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(se.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
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
		Name: "port",
	}
	ids, err = se.QueryEquipment().
		Select(equipment.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "Equipment",
		Name: "equipment",
	}
	ids, err = se.QueryService().
		Select(service.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Service",
		Name: "service",
	}
	ids, err = se.QueryDefinition().
		Select(serviceendpointdefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "ServiceEndpointDefinition",
		Name: "definition",
	}
	return node, nil
}

func (se *ServiceEndpointMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := se.ID()
	if !exists {
		return nil, nil
	}
	ent, err := se.Client().ServiceEndpoint.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (sed *ServiceEndpointDefinition) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     sed.ID,
		Type:   "ServiceEndpointDefinition",
		Fields: make([]*Field, 5),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(sed.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sed.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sed.Role); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "role",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sed.Name); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sed.Index); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "index",
		Value: string(buf),
	}
	var ids []int
	ids, err = sed.QueryEndpoints().
		Select(serviceendpoint.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "ServiceEndpoint",
		Name: "endpoints",
	}
	ids, err = sed.QueryServiceType().
		Select(servicetype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "ServiceType",
		Name: "service_type",
	}
	ids, err = sed.QueryEquipmentType().
		Select(equipmenttype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "EquipmentType",
		Name: "equipment_type",
	}
	return node, nil
}

func (sed *ServiceEndpointDefinitionMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := sed.ID()
	if !exists {
		return nil, nil
	}
	ent, err := sed.Client().ServiceEndpointDefinition.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (st *ServiceType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     st.ID,
		Type:   "ServiceType",
		Fields: make([]*Field, 6),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(st.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(st.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(st.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(st.HasCustomer); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "bool",
		Name:  "has_customer",
		Value: string(buf),
	}
	if buf, err = json.Marshal(st.IsDeleted); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "bool",
		Name:  "is_deleted",
		Value: string(buf),
	}
	if buf, err = json.Marshal(st.DiscoveryMethod); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "servicetype.DiscoveryMethod",
		Name:  "discovery_method",
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
		Name: "services",
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
		Name: "property_types",
	}
	ids, err = st.QueryEndpointDefinitions().
		Select(serviceendpointdefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "ServiceEndpointDefinition",
		Name: "endpoint_definitions",
	}
	return node, nil
}

func (st *ServiceTypeMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := st.ID()
	if !exists {
		return nil, nil
	}
	ent, err := st.Client().ServiceType.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.OwnerName); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "owner_name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.CreationTimestamp); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "time.Time",
		Name:  "creation_timestamp",
		Value: string(buf),
	}
	if buf, err = json.Marshal(s.CompletionTimestamp); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "time.Time",
		Name:  "completion_timestamp",
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
		Name: "location",
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
		Name: "source_file",
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
		Name: "questions",
	}
	return node, nil
}

func (s *SurveyMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := s.ID()
	if !exists {
		return nil, nil
	}
	ent, err := s.Client().Survey.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (scs *SurveyCellScan) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     scs.ID,
		Type:   "SurveyCellScan",
		Fields: make([]*Field, 22),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(scs.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.NetworkType); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "network_type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.SignalStrength); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "int",
		Name:  "signal_strength",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Timestamp); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "time.Time",
		Name:  "timestamp",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.BaseStationID); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "base_station_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.NetworkID); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "string",
		Name:  "network_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.SystemID); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "system_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.CellID); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "string",
		Name:  "cell_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.LocationAreaCode); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "location_area_code",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.MobileCountryCode); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "string",
		Name:  "mobile_country_code",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.MobileNetworkCode); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "string",
		Name:  "mobile_network_code",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.PrimaryScramblingCode); err != nil {
		return nil, err
	}
	node.Fields[12] = &Field{
		Type:  "string",
		Name:  "primary_scrambling_code",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Operator); err != nil {
		return nil, err
	}
	node.Fields[13] = &Field{
		Type:  "string",
		Name:  "operator",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Arfcn); err != nil {
		return nil, err
	}
	node.Fields[14] = &Field{
		Type:  "int",
		Name:  "arfcn",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.PhysicalCellID); err != nil {
		return nil, err
	}
	node.Fields[15] = &Field{
		Type:  "string",
		Name:  "physical_cell_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.TrackingAreaCode); err != nil {
		return nil, err
	}
	node.Fields[16] = &Field{
		Type:  "string",
		Name:  "tracking_area_code",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.TimingAdvance); err != nil {
		return nil, err
	}
	node.Fields[17] = &Field{
		Type:  "int",
		Name:  "timing_advance",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Earfcn); err != nil {
		return nil, err
	}
	node.Fields[18] = &Field{
		Type:  "int",
		Name:  "earfcn",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Uarfcn); err != nil {
		return nil, err
	}
	node.Fields[19] = &Field{
		Type:  "int",
		Name:  "uarfcn",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Latitude); err != nil {
		return nil, err
	}
	node.Fields[20] = &Field{
		Type:  "float64",
		Name:  "latitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(scs.Longitude); err != nil {
		return nil, err
	}
	node.Fields[21] = &Field{
		Type:  "float64",
		Name:  "longitude",
		Value: string(buf),
	}
	var ids []int
	ids, err = scs.QueryChecklistItem().
		Select(checklistitem.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "CheckListItem",
		Name: "checklist_item",
	}
	ids, err = scs.QuerySurveyQuestion().
		Select(surveyquestion.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "SurveyQuestion",
		Name: "survey_question",
	}
	ids, err = scs.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "location",
	}
	return node, nil
}

func (scs *SurveyCellScanMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := scs.ID()
	if !exists {
		return nil, nil
	}
	ent, err := scs.Client().SurveyCellScan.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (sq *SurveyQuestion) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     sq.ID,
		Type:   "SurveyQuestion",
		Fields: make([]*Field, 20),
		Edges:  make([]*Edge, 5),
	}
	var buf []byte
	if buf, err = json.Marshal(sq.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.FormName); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "form_name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.FormDescription); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "form_description",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.FormIndex); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "int",
		Name:  "form_index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.QuestionType); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "question_type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.QuestionFormat); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "string",
		Name:  "question_format",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.QuestionText); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "question_text",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.QuestionIndex); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "int",
		Name:  "question_index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.BoolData); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "bool",
		Name:  "bool_data",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.EmailData); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "string",
		Name:  "email_data",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.Latitude); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "float64",
		Name:  "latitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.Longitude); err != nil {
		return nil, err
	}
	node.Fields[12] = &Field{
		Type:  "float64",
		Name:  "longitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.LocationAccuracy); err != nil {
		return nil, err
	}
	node.Fields[13] = &Field{
		Type:  "float64",
		Name:  "location_accuracy",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.Altitude); err != nil {
		return nil, err
	}
	node.Fields[14] = &Field{
		Type:  "float64",
		Name:  "altitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.PhoneData); err != nil {
		return nil, err
	}
	node.Fields[15] = &Field{
		Type:  "string",
		Name:  "phone_data",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.TextData); err != nil {
		return nil, err
	}
	node.Fields[16] = &Field{
		Type:  "string",
		Name:  "text_data",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.FloatData); err != nil {
		return nil, err
	}
	node.Fields[17] = &Field{
		Type:  "float64",
		Name:  "float_data",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.IntData); err != nil {
		return nil, err
	}
	node.Fields[18] = &Field{
		Type:  "int",
		Name:  "int_data",
		Value: string(buf),
	}
	if buf, err = json.Marshal(sq.DateData); err != nil {
		return nil, err
	}
	node.Fields[19] = &Field{
		Type:  "time.Time",
		Name:  "date_data",
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
		Name: "survey",
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
		Name: "wifi_scan",
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
		Name: "cell_scan",
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
		Name: "photo_data",
	}
	ids, err = sq.QueryImages().
		Select(file.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "File",
		Name: "images",
	}
	return node, nil
}

func (sq *SurveyQuestionMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := sq.ID()
	if !exists {
		return nil, nil
	}
	ent, err := sq.Client().SurveyQuestion.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (stc *SurveyTemplateCategory) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     stc.ID,
		Type:   "SurveyTemplateCategory",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(stc.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stc.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stc.CategoryTitle); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "category_title",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stc.CategoryDescription); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "category_description",
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
		Name: "survey_template_questions",
	}
	ids, err = stc.QueryLocationType().
		Select(locationtype.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "LocationType",
		Name: "location_type",
	}
	return node, nil
}

func (stc *SurveyTemplateCategoryMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := stc.ID()
	if !exists {
		return nil, nil
	}
	ent, err := stc.Client().SurveyTemplateCategory.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stq.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stq.QuestionTitle); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "question_title",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stq.QuestionDescription); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "question_description",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stq.QuestionType); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "question_type",
		Value: string(buf),
	}
	if buf, err = json.Marshal(stq.Index); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "int",
		Name:  "index",
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
		Name: "category",
	}
	return node, nil
}

func (stq *SurveyTemplateQuestionMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := stq.ID()
	if !exists {
		return nil, nil
	}
	ent, err := stq.Client().SurveyTemplateQuestion.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (swfs *SurveyWiFiScan) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     swfs.ID,
		Type:   "SurveyWiFiScan",
		Fields: make([]*Field, 13),
		Edges:  make([]*Edge, 3),
	}
	var buf []byte
	if buf, err = json.Marshal(swfs.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Ssid); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "ssid",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Bssid); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "bssid",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Timestamp); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "time.Time",
		Name:  "timestamp",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Frequency); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "int",
		Name:  "frequency",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Channel); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "int",
		Name:  "channel",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Band); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "string",
		Name:  "band",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.ChannelWidth); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "int",
		Name:  "channel_width",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Capabilities); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "string",
		Name:  "capabilities",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Strength); err != nil {
		return nil, err
	}
	node.Fields[10] = &Field{
		Type:  "int",
		Name:  "strength",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Latitude); err != nil {
		return nil, err
	}
	node.Fields[11] = &Field{
		Type:  "float64",
		Name:  "latitude",
		Value: string(buf),
	}
	if buf, err = json.Marshal(swfs.Longitude); err != nil {
		return nil, err
	}
	node.Fields[12] = &Field{
		Type:  "float64",
		Name:  "longitude",
		Value: string(buf),
	}
	var ids []int
	ids, err = swfs.QueryChecklistItem().
		Select(checklistitem.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "CheckListItem",
		Name: "checklist_item",
	}
	ids, err = swfs.QuerySurveyQuestion().
		Select(surveyquestion.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "SurveyQuestion",
		Name: "survey_question",
	}
	ids, err = swfs.QueryLocation().
		Select(location.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "Location",
		Name: "location",
	}
	return node, nil
}

func (swfs *SurveyWiFiScanMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := swfs.ID()
	if !exists {
		return nil, nil
	}
	ent, err := swfs.Client().SurveyWiFiScan.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (u *User) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     u.ID,
		Type:   "User",
		Fields: make([]*Field, 8),
		Edges:  make([]*Edge, 5),
	}
	var buf []byte
	if buf, err = json.Marshal(u.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.AuthID); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "auth_id",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.FirstName); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "first_name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.LastName); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "last_name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.Email); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "email",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.Status); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "user.Status",
		Name:  "status",
		Value: string(buf),
	}
	if buf, err = json.Marshal(u.Role); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "user.Role",
		Name:  "role",
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
		Name: "profile_photo",
	}
	ids, err = u.QueryGroups().
		Select(usersgroup.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "UsersGroup",
		Name: "groups",
	}
	ids, err = u.QueryOwnedWorkOrders().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[2] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "owned_work_orders",
	}
	ids, err = u.QueryAssignedWorkOrders().
		Select(workorder.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "WorkOrder",
		Name: "assigned_work_orders",
	}
	ids, err = u.QueryCreatedProjects().
		Select(project.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[4] = &Edge{
		IDs:  ids,
		Type: "Project",
		Name: "created_projects",
	}
	return node, nil
}

func (u *UserMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := u.ID()
	if !exists {
		return nil, nil
	}
	ent, err := u.Client().User.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (ug *UsersGroup) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     ug.ID,
		Type:   "UsersGroup",
		Fields: make([]*Field, 5),
		Edges:  make([]*Edge, 2),
	}
	var buf []byte
	if buf, err = json.Marshal(ug.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ug.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ug.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ug.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "description",
		Value: string(buf),
	}
	if buf, err = json.Marshal(ug.Status); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "usersgroup.Status",
		Name:  "status",
		Value: string(buf),
	}
	var ids []int
	ids, err = ug.QueryMembers().
		Select(user.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[0] = &Edge{
		IDs:  ids,
		Type: "User",
		Name: "members",
	}
	ids, err = ug.QueryPolicies().
		Select(permissionspolicy.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[1] = &Edge{
		IDs:  ids,
		Type: "PermissionsPolicy",
		Name: "policies",
	}
	return node, nil
}

func (ug *UsersGroupMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := ug.ID()
	if !exists {
		return nil, nil
	}
	ent, err := ug.Client().UsersGroup.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (wo *WorkOrder) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     wo.ID,
		Type:   "WorkOrder",
		Fields: make([]*Field, 10),
		Edges:  make([]*Edge, 13),
	}
	var buf []byte
	if buf, err = json.Marshal(wo.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Status); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "status",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Priority); err != nil {
		return nil, err
	}
	node.Fields[4] = &Field{
		Type:  "string",
		Name:  "priority",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Description); err != nil {
		return nil, err
	}
	node.Fields[5] = &Field{
		Type:  "string",
		Name:  "description",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.InstallDate); err != nil {
		return nil, err
	}
	node.Fields[6] = &Field{
		Type:  "time.Time",
		Name:  "install_date",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.CreationDate); err != nil {
		return nil, err
	}
	node.Fields[7] = &Field{
		Type:  "time.Time",
		Name:  "creation_date",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.Index); err != nil {
		return nil, err
	}
	node.Fields[8] = &Field{
		Type:  "int",
		Name:  "index",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wo.CloseDate); err != nil {
		return nil, err
	}
	node.Fields[9] = &Field{
		Type:  "time.Time",
		Name:  "close_date",
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
		Name: "type",
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
		Name: "equipment",
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
		Name: "links",
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
		Name: "files",
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
		Name: "hyperlinks",
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
		Name: "location",
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
		Name: "comments",
	}
	ids, err = wo.QueryActivities().
		Select(activity.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[7] = &Edge{
		IDs:  ids,
		Type: "Activity",
		Name: "activities",
	}
	ids, err = wo.QueryProperties().
		Select(property.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[8] = &Edge{
		IDs:  ids,
		Type: "Property",
		Name: "properties",
	}
	ids, err = wo.QueryCheckListCategories().
		Select(checklistcategory.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[9] = &Edge{
		IDs:  ids,
		Type: "CheckListCategory",
		Name: "check_list_categories",
	}
	ids, err = wo.QueryProject().
		Select(project.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[10] = &Edge{
		IDs:  ids,
		Type: "Project",
		Name: "project",
	}
	ids, err = wo.QueryOwner().
		Select(user.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[11] = &Edge{
		IDs:  ids,
		Type: "User",
		Name: "owner",
	}
	ids, err = wo.QueryAssignee().
		Select(user.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[12] = &Edge{
		IDs:  ids,
		Type: "User",
		Name: "assignee",
	}
	return node, nil
}

func (wo *WorkOrderMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := wo.ID()
	if !exists {
		return nil, nil
	}
	ent, err := wo.Client().WorkOrder.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wod.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wod.Index); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "int",
		Name:  "index",
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
		Name: "type",
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
		Name: "project_type",
	}
	return node, nil
}

func (wod *WorkOrderDefinitionMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := wod.ID()
	if !exists {
		return nil, nil
	}
	ent, err := wod.Client().WorkOrderDefinition.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
}

func (wot *WorkOrderType) Node(ctx context.Context) (node *Node, err error) {
	node = &Node{
		ID:     wot.ID,
		Type:   "WorkOrderType",
		Fields: make([]*Field, 4),
		Edges:  make([]*Edge, 4),
	}
	var buf []byte
	if buf, err = json.Marshal(wot.CreateTime); err != nil {
		return nil, err
	}
	node.Fields[0] = &Field{
		Type:  "time.Time",
		Name:  "create_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wot.UpdateTime); err != nil {
		return nil, err
	}
	node.Fields[1] = &Field{
		Type:  "time.Time",
		Name:  "update_time",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wot.Name); err != nil {
		return nil, err
	}
	node.Fields[2] = &Field{
		Type:  "string",
		Name:  "name",
		Value: string(buf),
	}
	if buf, err = json.Marshal(wot.Description); err != nil {
		return nil, err
	}
	node.Fields[3] = &Field{
		Type:  "string",
		Name:  "description",
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
		Name: "work_orders",
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
		Name: "property_types",
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
		Name: "definitions",
	}
	ids, err = wot.QueryCheckListCategoryDefinitions().
		Select(checklistcategorydefinition.FieldID).
		Ints(ctx)
	if err != nil {
		return nil, err
	}
	node.Edges[3] = &Edge{
		IDs:  ids,
		Type: "CheckListCategoryDefinition",
		Name: "check_list_category_definitions",
	}
	return node, nil
}

func (wot *WorkOrderTypeMutation) Node(ctx context.Context) (node *Node, err error) {
	id, exists := wot.ID()
	if !exists {
		return nil, nil
	}
	ent, err := wot.Client().WorkOrderType.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return ent.Node(ctx)
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
	case activity.Table:
		n, err := c.Activity.Query().
			Where(activity.ID(id)).
			CollectFields(ctx, "Activity").
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
	case checklistcategorydefinition.Table:
		n, err := c.CheckListCategoryDefinition.Query().
			Where(checklistcategorydefinition.ID(id)).
			CollectFields(ctx, "CheckListCategoryDefinition").
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
	case permissionspolicy.Table:
		n, err := c.PermissionsPolicy.Query().
			Where(permissionspolicy.ID(id)).
			CollectFields(ctx, "PermissionsPolicy").
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
	case serviceendpointdefinition.Table:
		n, err := c.ServiceEndpointDefinition.Query().
			Where(serviceendpointdefinition.ID(id)).
			CollectFields(ctx, "ServiceEndpointDefinition").
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
	case user.Table:
		n, err := c.User.Query().
			Where(user.ID(id)).
			CollectFields(ctx, "User").
			Only(ctx)
		if err != nil {
			return nil, err
		}
		return n, nil
	case usersgroup.Table:
		n, err := c.UsersGroup.Query().
			Where(usersgroup.ID(id)).
			CollectFields(ctx, "UsersGroup").
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
