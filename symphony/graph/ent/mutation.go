// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"time"

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
	"github.com/facebookincubator/symphony/pkg/actions/core"

	"github.com/facebookincubator/ent"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeActionsRule                 = "ActionsRule"
	TypeCheckListCategory           = "CheckListCategory"
	TypeCheckListItem               = "CheckListItem"
	TypeCheckListItemDefinition     = "CheckListItemDefinition"
	TypeComment                     = "Comment"
	TypeCustomer                    = "Customer"
	TypeEquipment                   = "Equipment"
	TypeEquipmentCategory           = "EquipmentCategory"
	TypeEquipmentPort               = "EquipmentPort"
	TypeEquipmentPortDefinition     = "EquipmentPortDefinition"
	TypeEquipmentPortType           = "EquipmentPortType"
	TypeEquipmentPosition           = "EquipmentPosition"
	TypeEquipmentPositionDefinition = "EquipmentPositionDefinition"
	TypeEquipmentType               = "EquipmentType"
	TypeFile                        = "File"
	TypeFloorPlan                   = "FloorPlan"
	TypeFloorPlanReferencePoint     = "FloorPlanReferencePoint"
	TypeFloorPlanScale              = "FloorPlanScale"
	TypeHyperlink                   = "Hyperlink"
	TypeLink                        = "Link"
	TypeLocation                    = "Location"
	TypeLocationType                = "LocationType"
	TypeProject                     = "Project"
	TypeProjectType                 = "ProjectType"
	TypeProperty                    = "Property"
	TypePropertyType                = "PropertyType"
	TypeReportFilter                = "ReportFilter"
	TypeService                     = "Service"
	TypeServiceEndpoint             = "ServiceEndpoint"
	TypeServiceType                 = "ServiceType"
	TypeSurvey                      = "Survey"
	TypeSurveyCellScan              = "SurveyCellScan"
	TypeSurveyQuestion              = "SurveyQuestion"
	TypeSurveyTemplateCategory      = "SurveyTemplateCategory"
	TypeSurveyTemplateQuestion      = "SurveyTemplateQuestion"
	TypeSurveyWiFiScan              = "SurveyWiFiScan"
	TypeTechnician                  = "Technician"
	TypeUser                        = "User"
	TypeWorkOrder                   = "WorkOrder"
	TypeWorkOrderDefinition         = "WorkOrderDefinition"
	TypeWorkOrderType               = "WorkOrderType"
)

// ActionsRuleMutation represents an operation that mutate the ActionsRules
// nodes in the graph.
type ActionsRuleMutation struct {
	config
	op            Op
	typ           string
	id            *int
	create_time   *time.Time
	update_time   *time.Time
	name          *string
	triggerID     *string
	ruleFilters   *[]*core.ActionsRuleFilter
	ruleActions   *[]*core.ActionsRuleAction
	clearedFields map[string]bool
}

var _ ent.Mutation = (*ActionsRuleMutation)(nil)

// newActionsRuleMutation creates new mutation for $n.Name.
func newActionsRuleMutation(c config, op Op) *ActionsRuleMutation {
	return &ActionsRuleMutation{
		config:        c,
		op:            op,
		typ:           TypeActionsRule,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ActionsRuleMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ActionsRuleMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *ActionsRuleMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *ActionsRuleMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *ActionsRuleMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *ActionsRuleMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *ActionsRuleMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *ActionsRuleMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *ActionsRuleMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *ActionsRuleMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *ActionsRuleMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *ActionsRuleMutation) ResetName() {
	m.name = nil
}

// SetTriggerID sets the triggerID field.
func (m *ActionsRuleMutation) SetTriggerID(s string) {
	m.triggerID = &s
}

// TriggerID returns the triggerID value in the mutation.
func (m *ActionsRuleMutation) TriggerID() (r string, exists bool) {
	v := m.triggerID
	if v == nil {
		return
	}
	return *v, true
}

// ResetTriggerID reset all changes of the triggerID field.
func (m *ActionsRuleMutation) ResetTriggerID() {
	m.triggerID = nil
}

// SetRuleFilters sets the ruleFilters field.
func (m *ActionsRuleMutation) SetRuleFilters(crf []*core.ActionsRuleFilter) {
	m.ruleFilters = &crf
}

// RuleFilters returns the ruleFilters value in the mutation.
func (m *ActionsRuleMutation) RuleFilters() (r []*core.ActionsRuleFilter, exists bool) {
	v := m.ruleFilters
	if v == nil {
		return
	}
	return *v, true
}

// ResetRuleFilters reset all changes of the ruleFilters field.
func (m *ActionsRuleMutation) ResetRuleFilters() {
	m.ruleFilters = nil
}

// SetRuleActions sets the ruleActions field.
func (m *ActionsRuleMutation) SetRuleActions(cra []*core.ActionsRuleAction) {
	m.ruleActions = &cra
}

// RuleActions returns the ruleActions value in the mutation.
func (m *ActionsRuleMutation) RuleActions() (r []*core.ActionsRuleAction, exists bool) {
	v := m.ruleActions
	if v == nil {
		return
	}
	return *v, true
}

// ResetRuleActions reset all changes of the ruleActions field.
func (m *ActionsRuleMutation) ResetRuleActions() {
	m.ruleActions = nil
}

// Op returns the operation name.
func (m *ActionsRuleMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (ActionsRule).
func (m *ActionsRuleMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *ActionsRuleMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m.create_time != nil {
		fields = append(fields, actionsrule.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, actionsrule.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, actionsrule.FieldName)
	}
	if m.triggerID != nil {
		fields = append(fields, actionsrule.FieldTriggerID)
	}
	if m.ruleFilters != nil {
		fields = append(fields, actionsrule.FieldRuleFilters)
	}
	if m.ruleActions != nil {
		fields = append(fields, actionsrule.FieldRuleActions)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *ActionsRuleMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case actionsrule.FieldCreateTime:
		return m.CreateTime()
	case actionsrule.FieldUpdateTime:
		return m.UpdateTime()
	case actionsrule.FieldName:
		return m.Name()
	case actionsrule.FieldTriggerID:
		return m.TriggerID()
	case actionsrule.FieldRuleFilters:
		return m.RuleFilters()
	case actionsrule.FieldRuleActions:
		return m.RuleActions()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ActionsRuleMutation) SetField(name string, value ent.Value) error {
	switch name {
	case actionsrule.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case actionsrule.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case actionsrule.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case actionsrule.FieldTriggerID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTriggerID(v)
		return nil
	case actionsrule.FieldRuleFilters:
		v, ok := value.([]*core.ActionsRuleFilter)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRuleFilters(v)
		return nil
	case actionsrule.FieldRuleActions:
		v, ok := value.([]*core.ActionsRuleAction)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRuleActions(v)
		return nil
	}
	return fmt.Errorf("unknown ActionsRule field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *ActionsRuleMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *ActionsRuleMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ActionsRuleMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ActionsRule numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *ActionsRuleMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *ActionsRuleMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *ActionsRuleMutation) ClearField(name string) error {
	return fmt.Errorf("unknown ActionsRule nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *ActionsRuleMutation) ResetField(name string) error {
	switch name {
	case actionsrule.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case actionsrule.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case actionsrule.FieldName:
		m.ResetName()
		return nil
	case actionsrule.FieldTriggerID:
		m.ResetTriggerID()
		return nil
	case actionsrule.FieldRuleFilters:
		m.ResetRuleFilters()
		return nil
	case actionsrule.FieldRuleActions:
		m.ResetRuleActions()
		return nil
	}
	return fmt.Errorf("unknown ActionsRule field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *ActionsRuleMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *ActionsRuleMutation) AddedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *ActionsRuleMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *ActionsRuleMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *ActionsRuleMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *ActionsRuleMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *ActionsRuleMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown ActionsRule unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *ActionsRuleMutation) ResetEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown ActionsRule edge %s", name)
}

// CheckListCategoryMutation represents an operation that mutate the CheckListCategories
// nodes in the graph.
type CheckListCategoryMutation struct {
	config
	op                      Op
	typ                     string
	id                      *int
	create_time             *time.Time
	update_time             *time.Time
	title                   *string
	description             *string
	clearedFields           map[string]bool
	check_list_items        map[int]struct{}
	removedcheck_list_items map[int]struct{}
}

var _ ent.Mutation = (*CheckListCategoryMutation)(nil)

// newCheckListCategoryMutation creates new mutation for $n.Name.
func newCheckListCategoryMutation(c config, op Op) *CheckListCategoryMutation {
	return &CheckListCategoryMutation{
		config:        c,
		op:            op,
		typ:           TypeCheckListCategory,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m CheckListCategoryMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m CheckListCategoryMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *CheckListCategoryMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *CheckListCategoryMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *CheckListCategoryMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *CheckListCategoryMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *CheckListCategoryMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *CheckListCategoryMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *CheckListCategoryMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetTitle sets the title field.
func (m *CheckListCategoryMutation) SetTitle(s string) {
	m.title = &s
}

// Title returns the title value in the mutation.
func (m *CheckListCategoryMutation) Title() (r string, exists bool) {
	v := m.title
	if v == nil {
		return
	}
	return *v, true
}

// ResetTitle reset all changes of the title field.
func (m *CheckListCategoryMutation) ResetTitle() {
	m.title = nil
}

// SetDescription sets the description field.
func (m *CheckListCategoryMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the description value in the mutation.
func (m *CheckListCategoryMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// ClearDescription clears the value of description.
func (m *CheckListCategoryMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[checklistcategory.FieldDescription] = true
}

// DescriptionCleared returns if the field description was cleared in this mutation.
func (m *CheckListCategoryMutation) DescriptionCleared() bool {
	return m.clearedFields[checklistcategory.FieldDescription]
}

// ResetDescription reset all changes of the description field.
func (m *CheckListCategoryMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, checklistcategory.FieldDescription)
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (m *CheckListCategoryMutation) AddCheckListItemIDs(ids ...int) {
	if m.check_list_items == nil {
		m.check_list_items = make(map[int]struct{})
	}
	for i := range ids {
		m.check_list_items[ids[i]] = struct{}{}
	}
}

// RemoveCheckListItemIDs removes the check_list_items edge to CheckListItem by ids.
func (m *CheckListCategoryMutation) RemoveCheckListItemIDs(ids ...int) {
	if m.removedcheck_list_items == nil {
		m.removedcheck_list_items = make(map[int]struct{})
	}
	for i := range ids {
		m.removedcheck_list_items[ids[i]] = struct{}{}
	}
}

// RemovedCheckListItems returns the removed ids of check_list_items.
func (m *CheckListCategoryMutation) RemovedCheckListItemsIDs() (ids []int) {
	for id := range m.removedcheck_list_items {
		ids = append(ids, id)
	}
	return
}

// CheckListItemsIDs returns the check_list_items ids in the mutation.
func (m *CheckListCategoryMutation) CheckListItemsIDs() (ids []int) {
	for id := range m.check_list_items {
		ids = append(ids, id)
	}
	return
}

// ResetCheckListItems reset all changes of the check_list_items edge.
func (m *CheckListCategoryMutation) ResetCheckListItems() {
	m.check_list_items = nil
	m.removedcheck_list_items = nil
}

// Op returns the operation name.
func (m *CheckListCategoryMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (CheckListCategory).
func (m *CheckListCategoryMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *CheckListCategoryMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.create_time != nil {
		fields = append(fields, checklistcategory.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, checklistcategory.FieldUpdateTime)
	}
	if m.title != nil {
		fields = append(fields, checklistcategory.FieldTitle)
	}
	if m.description != nil {
		fields = append(fields, checklistcategory.FieldDescription)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *CheckListCategoryMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case checklistcategory.FieldCreateTime:
		return m.CreateTime()
	case checklistcategory.FieldUpdateTime:
		return m.UpdateTime()
	case checklistcategory.FieldTitle:
		return m.Title()
	case checklistcategory.FieldDescription:
		return m.Description()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *CheckListCategoryMutation) SetField(name string, value ent.Value) error {
	switch name {
	case checklistcategory.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case checklistcategory.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case checklistcategory.FieldTitle:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTitle(v)
		return nil
	case checklistcategory.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	}
	return fmt.Errorf("unknown CheckListCategory field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *CheckListCategoryMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *CheckListCategoryMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *CheckListCategoryMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown CheckListCategory numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *CheckListCategoryMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[checklistcategory.FieldDescription] {
		fields = append(fields, checklistcategory.FieldDescription)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *CheckListCategoryMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *CheckListCategoryMutation) ClearField(name string) error {
	switch name {
	case checklistcategory.FieldDescription:
		m.ClearDescription()
		return nil
	}
	return fmt.Errorf("unknown CheckListCategory nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *CheckListCategoryMutation) ResetField(name string) error {
	switch name {
	case checklistcategory.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case checklistcategory.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case checklistcategory.FieldTitle:
		m.ResetTitle()
		return nil
	case checklistcategory.FieldDescription:
		m.ResetDescription()
		return nil
	}
	return fmt.Errorf("unknown CheckListCategory field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *CheckListCategoryMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.check_list_items != nil {
		edges = append(edges, checklistcategory.EdgeCheckListItems)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *CheckListCategoryMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case checklistcategory.EdgeCheckListItems:
		ids := make([]ent.Value, 0, len(m.check_list_items))
		for id := range m.check_list_items {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *CheckListCategoryMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedcheck_list_items != nil {
		edges = append(edges, checklistcategory.EdgeCheckListItems)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *CheckListCategoryMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case checklistcategory.EdgeCheckListItems:
		ids := make([]ent.Value, 0, len(m.removedcheck_list_items))
		for id := range m.removedcheck_list_items {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *CheckListCategoryMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *CheckListCategoryMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *CheckListCategoryMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown CheckListCategory unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *CheckListCategoryMutation) ResetEdge(name string) error {
	switch name {
	case checklistcategory.EdgeCheckListItems:
		m.ResetCheckListItems()
		return nil
	}
	return fmt.Errorf("unknown CheckListCategory edge %s", name)
}

// CheckListItemMutation represents an operation that mutate the CheckListItems
// nodes in the graph.
type CheckListItemMutation struct {
	config
	op                   Op
	typ                  string
	id                   *int
	title                *string
	_type                *string
	index                *int
	addindex             *int
	checked              *bool
	string_val           *string
	enum_values          *string
	enum_selection_mode  *string
	selected_enum_values *string
	help_text            *string
	clearedFields        map[string]bool
	work_order           *int
	clearedwork_order    bool
}

var _ ent.Mutation = (*CheckListItemMutation)(nil)

// newCheckListItemMutation creates new mutation for $n.Name.
func newCheckListItemMutation(c config, op Op) *CheckListItemMutation {
	return &CheckListItemMutation{
		config:        c,
		op:            op,
		typ:           TypeCheckListItem,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m CheckListItemMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m CheckListItemMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *CheckListItemMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetTitle sets the title field.
func (m *CheckListItemMutation) SetTitle(s string) {
	m.title = &s
}

// Title returns the title value in the mutation.
func (m *CheckListItemMutation) Title() (r string, exists bool) {
	v := m.title
	if v == nil {
		return
	}
	return *v, true
}

// ResetTitle reset all changes of the title field.
func (m *CheckListItemMutation) ResetTitle() {
	m.title = nil
}

// SetType sets the type field.
func (m *CheckListItemMutation) SetType(s string) {
	m._type = &s
}

// GetType returns the type value in the mutation.
func (m *CheckListItemMutation) GetType() (r string, exists bool) {
	v := m._type
	if v == nil {
		return
	}
	return *v, true
}

// ResetType reset all changes of the type field.
func (m *CheckListItemMutation) ResetType() {
	m._type = nil
}

// SetIndex sets the index field.
func (m *CheckListItemMutation) SetIndex(i int) {
	m.index = &i
	m.addindex = nil
}

// Index returns the index value in the mutation.
func (m *CheckListItemMutation) Index() (r int, exists bool) {
	v := m.index
	if v == nil {
		return
	}
	return *v, true
}

// AddIndex adds i to index.
func (m *CheckListItemMutation) AddIndex(i int) {
	if m.addindex != nil {
		*m.addindex += i
	} else {
		m.addindex = &i
	}
}

// AddedIndex returns the value that was added to the index field in this mutation.
func (m *CheckListItemMutation) AddedIndex() (r int, exists bool) {
	v := m.addindex
	if v == nil {
		return
	}
	return *v, true
}

// ClearIndex clears the value of index.
func (m *CheckListItemMutation) ClearIndex() {
	m.index = nil
	m.addindex = nil
	m.clearedFields[checklistitem.FieldIndex] = true
}

// IndexCleared returns if the field index was cleared in this mutation.
func (m *CheckListItemMutation) IndexCleared() bool {
	return m.clearedFields[checklistitem.FieldIndex]
}

// ResetIndex reset all changes of the index field.
func (m *CheckListItemMutation) ResetIndex() {
	m.index = nil
	m.addindex = nil
	delete(m.clearedFields, checklistitem.FieldIndex)
}

// SetChecked sets the checked field.
func (m *CheckListItemMutation) SetChecked(b bool) {
	m.checked = &b
}

// Checked returns the checked value in the mutation.
func (m *CheckListItemMutation) Checked() (r bool, exists bool) {
	v := m.checked
	if v == nil {
		return
	}
	return *v, true
}

// ClearChecked clears the value of checked.
func (m *CheckListItemMutation) ClearChecked() {
	m.checked = nil
	m.clearedFields[checklistitem.FieldChecked] = true
}

// CheckedCleared returns if the field checked was cleared in this mutation.
func (m *CheckListItemMutation) CheckedCleared() bool {
	return m.clearedFields[checklistitem.FieldChecked]
}

// ResetChecked reset all changes of the checked field.
func (m *CheckListItemMutation) ResetChecked() {
	m.checked = nil
	delete(m.clearedFields, checklistitem.FieldChecked)
}

// SetStringVal sets the string_val field.
func (m *CheckListItemMutation) SetStringVal(s string) {
	m.string_val = &s
}

// StringVal returns the string_val value in the mutation.
func (m *CheckListItemMutation) StringVal() (r string, exists bool) {
	v := m.string_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearStringVal clears the value of string_val.
func (m *CheckListItemMutation) ClearStringVal() {
	m.string_val = nil
	m.clearedFields[checklistitem.FieldStringVal] = true
}

// StringValCleared returns if the field string_val was cleared in this mutation.
func (m *CheckListItemMutation) StringValCleared() bool {
	return m.clearedFields[checklistitem.FieldStringVal]
}

// ResetStringVal reset all changes of the string_val field.
func (m *CheckListItemMutation) ResetStringVal() {
	m.string_val = nil
	delete(m.clearedFields, checklistitem.FieldStringVal)
}

// SetEnumValues sets the enum_values field.
func (m *CheckListItemMutation) SetEnumValues(s string) {
	m.enum_values = &s
}

// EnumValues returns the enum_values value in the mutation.
func (m *CheckListItemMutation) EnumValues() (r string, exists bool) {
	v := m.enum_values
	if v == nil {
		return
	}
	return *v, true
}

// ClearEnumValues clears the value of enum_values.
func (m *CheckListItemMutation) ClearEnumValues() {
	m.enum_values = nil
	m.clearedFields[checklistitem.FieldEnumValues] = true
}

// EnumValuesCleared returns if the field enum_values was cleared in this mutation.
func (m *CheckListItemMutation) EnumValuesCleared() bool {
	return m.clearedFields[checklistitem.FieldEnumValues]
}

// ResetEnumValues reset all changes of the enum_values field.
func (m *CheckListItemMutation) ResetEnumValues() {
	m.enum_values = nil
	delete(m.clearedFields, checklistitem.FieldEnumValues)
}

// SetEnumSelectionMode sets the enum_selection_mode field.
func (m *CheckListItemMutation) SetEnumSelectionMode(s string) {
	m.enum_selection_mode = &s
}

// EnumSelectionMode returns the enum_selection_mode value in the mutation.
func (m *CheckListItemMutation) EnumSelectionMode() (r string, exists bool) {
	v := m.enum_selection_mode
	if v == nil {
		return
	}
	return *v, true
}

// ClearEnumSelectionMode clears the value of enum_selection_mode.
func (m *CheckListItemMutation) ClearEnumSelectionMode() {
	m.enum_selection_mode = nil
	m.clearedFields[checklistitem.FieldEnumSelectionMode] = true
}

// EnumSelectionModeCleared returns if the field enum_selection_mode was cleared in this mutation.
func (m *CheckListItemMutation) EnumSelectionModeCleared() bool {
	return m.clearedFields[checklistitem.FieldEnumSelectionMode]
}

// ResetEnumSelectionMode reset all changes of the enum_selection_mode field.
func (m *CheckListItemMutation) ResetEnumSelectionMode() {
	m.enum_selection_mode = nil
	delete(m.clearedFields, checklistitem.FieldEnumSelectionMode)
}

// SetSelectedEnumValues sets the selected_enum_values field.
func (m *CheckListItemMutation) SetSelectedEnumValues(s string) {
	m.selected_enum_values = &s
}

// SelectedEnumValues returns the selected_enum_values value in the mutation.
func (m *CheckListItemMutation) SelectedEnumValues() (r string, exists bool) {
	v := m.selected_enum_values
	if v == nil {
		return
	}
	return *v, true
}

// ClearSelectedEnumValues clears the value of selected_enum_values.
func (m *CheckListItemMutation) ClearSelectedEnumValues() {
	m.selected_enum_values = nil
	m.clearedFields[checklistitem.FieldSelectedEnumValues] = true
}

// SelectedEnumValuesCleared returns if the field selected_enum_values was cleared in this mutation.
func (m *CheckListItemMutation) SelectedEnumValuesCleared() bool {
	return m.clearedFields[checklistitem.FieldSelectedEnumValues]
}

// ResetSelectedEnumValues reset all changes of the selected_enum_values field.
func (m *CheckListItemMutation) ResetSelectedEnumValues() {
	m.selected_enum_values = nil
	delete(m.clearedFields, checklistitem.FieldSelectedEnumValues)
}

// SetHelpText sets the help_text field.
func (m *CheckListItemMutation) SetHelpText(s string) {
	m.help_text = &s
}

// HelpText returns the help_text value in the mutation.
func (m *CheckListItemMutation) HelpText() (r string, exists bool) {
	v := m.help_text
	if v == nil {
		return
	}
	return *v, true
}

// ClearHelpText clears the value of help_text.
func (m *CheckListItemMutation) ClearHelpText() {
	m.help_text = nil
	m.clearedFields[checklistitem.FieldHelpText] = true
}

// HelpTextCleared returns if the field help_text was cleared in this mutation.
func (m *CheckListItemMutation) HelpTextCleared() bool {
	return m.clearedFields[checklistitem.FieldHelpText]
}

// ResetHelpText reset all changes of the help_text field.
func (m *CheckListItemMutation) ResetHelpText() {
	m.help_text = nil
	delete(m.clearedFields, checklistitem.FieldHelpText)
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (m *CheckListItemMutation) SetWorkOrderID(id int) {
	m.work_order = &id
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (m *CheckListItemMutation) ClearWorkOrder() {
	m.clearedwork_order = true
}

// WorkOrderCleared returns if the edge work_order was cleared.
func (m *CheckListItemMutation) WorkOrderCleared() bool {
	return m.clearedwork_order
}

// WorkOrderID returns the work_order id in the mutation.
func (m *CheckListItemMutation) WorkOrderID() (id int, exists bool) {
	if m.work_order != nil {
		return *m.work_order, true
	}
	return
}

// WorkOrderIDs returns the work_order ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// WorkOrderID instead. It exists only for internal usage by the builders.
func (m *CheckListItemMutation) WorkOrderIDs() (ids []int) {
	if id := m.work_order; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetWorkOrder reset all changes of the work_order edge.
func (m *CheckListItemMutation) ResetWorkOrder() {
	m.work_order = nil
	m.clearedwork_order = false
}

// Op returns the operation name.
func (m *CheckListItemMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (CheckListItem).
func (m *CheckListItemMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *CheckListItemMutation) Fields() []string {
	fields := make([]string, 0, 9)
	if m.title != nil {
		fields = append(fields, checklistitem.FieldTitle)
	}
	if m._type != nil {
		fields = append(fields, checklistitem.FieldType)
	}
	if m.index != nil {
		fields = append(fields, checklistitem.FieldIndex)
	}
	if m.checked != nil {
		fields = append(fields, checklistitem.FieldChecked)
	}
	if m.string_val != nil {
		fields = append(fields, checklistitem.FieldStringVal)
	}
	if m.enum_values != nil {
		fields = append(fields, checklistitem.FieldEnumValues)
	}
	if m.enum_selection_mode != nil {
		fields = append(fields, checklistitem.FieldEnumSelectionMode)
	}
	if m.selected_enum_values != nil {
		fields = append(fields, checklistitem.FieldSelectedEnumValues)
	}
	if m.help_text != nil {
		fields = append(fields, checklistitem.FieldHelpText)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *CheckListItemMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case checklistitem.FieldTitle:
		return m.Title()
	case checklistitem.FieldType:
		return m.GetType()
	case checklistitem.FieldIndex:
		return m.Index()
	case checklistitem.FieldChecked:
		return m.Checked()
	case checklistitem.FieldStringVal:
		return m.StringVal()
	case checklistitem.FieldEnumValues:
		return m.EnumValues()
	case checklistitem.FieldEnumSelectionMode:
		return m.EnumSelectionMode()
	case checklistitem.FieldSelectedEnumValues:
		return m.SelectedEnumValues()
	case checklistitem.FieldHelpText:
		return m.HelpText()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *CheckListItemMutation) SetField(name string, value ent.Value) error {
	switch name {
	case checklistitem.FieldTitle:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTitle(v)
		return nil
	case checklistitem.FieldType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
		return nil
	case checklistitem.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIndex(v)
		return nil
	case checklistitem.FieldChecked:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetChecked(v)
		return nil
	case checklistitem.FieldStringVal:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStringVal(v)
		return nil
	case checklistitem.FieldEnumValues:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnumValues(v)
		return nil
	case checklistitem.FieldEnumSelectionMode:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnumSelectionMode(v)
		return nil
	case checklistitem.FieldSelectedEnumValues:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSelectedEnumValues(v)
		return nil
	case checklistitem.FieldHelpText:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetHelpText(v)
		return nil
	}
	return fmt.Errorf("unknown CheckListItem field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *CheckListItemMutation) AddedFields() []string {
	var fields []string
	if m.addindex != nil {
		fields = append(fields, checklistitem.FieldIndex)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *CheckListItemMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case checklistitem.FieldIndex:
		return m.AddedIndex()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *CheckListItemMutation) AddField(name string, value ent.Value) error {
	switch name {
	case checklistitem.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIndex(v)
		return nil
	}
	return fmt.Errorf("unknown CheckListItem numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *CheckListItemMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[checklistitem.FieldIndex] {
		fields = append(fields, checklistitem.FieldIndex)
	}
	if m.clearedFields[checklistitem.FieldChecked] {
		fields = append(fields, checklistitem.FieldChecked)
	}
	if m.clearedFields[checklistitem.FieldStringVal] {
		fields = append(fields, checklistitem.FieldStringVal)
	}
	if m.clearedFields[checklistitem.FieldEnumValues] {
		fields = append(fields, checklistitem.FieldEnumValues)
	}
	if m.clearedFields[checklistitem.FieldEnumSelectionMode] {
		fields = append(fields, checklistitem.FieldEnumSelectionMode)
	}
	if m.clearedFields[checklistitem.FieldSelectedEnumValues] {
		fields = append(fields, checklistitem.FieldSelectedEnumValues)
	}
	if m.clearedFields[checklistitem.FieldHelpText] {
		fields = append(fields, checklistitem.FieldHelpText)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *CheckListItemMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *CheckListItemMutation) ClearField(name string) error {
	switch name {
	case checklistitem.FieldIndex:
		m.ClearIndex()
		return nil
	case checklistitem.FieldChecked:
		m.ClearChecked()
		return nil
	case checklistitem.FieldStringVal:
		m.ClearStringVal()
		return nil
	case checklistitem.FieldEnumValues:
		m.ClearEnumValues()
		return nil
	case checklistitem.FieldEnumSelectionMode:
		m.ClearEnumSelectionMode()
		return nil
	case checklistitem.FieldSelectedEnumValues:
		m.ClearSelectedEnumValues()
		return nil
	case checklistitem.FieldHelpText:
		m.ClearHelpText()
		return nil
	}
	return fmt.Errorf("unknown CheckListItem nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *CheckListItemMutation) ResetField(name string) error {
	switch name {
	case checklistitem.FieldTitle:
		m.ResetTitle()
		return nil
	case checklistitem.FieldType:
		m.ResetType()
		return nil
	case checklistitem.FieldIndex:
		m.ResetIndex()
		return nil
	case checklistitem.FieldChecked:
		m.ResetChecked()
		return nil
	case checklistitem.FieldStringVal:
		m.ResetStringVal()
		return nil
	case checklistitem.FieldEnumValues:
		m.ResetEnumValues()
		return nil
	case checklistitem.FieldEnumSelectionMode:
		m.ResetEnumSelectionMode()
		return nil
	case checklistitem.FieldSelectedEnumValues:
		m.ResetSelectedEnumValues()
		return nil
	case checklistitem.FieldHelpText:
		m.ResetHelpText()
		return nil
	}
	return fmt.Errorf("unknown CheckListItem field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *CheckListItemMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.work_order != nil {
		edges = append(edges, checklistitem.EdgeWorkOrder)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *CheckListItemMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case checklistitem.EdgeWorkOrder:
		if id := m.work_order; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *CheckListItemMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *CheckListItemMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *CheckListItemMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedwork_order {
		edges = append(edges, checklistitem.EdgeWorkOrder)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *CheckListItemMutation) EdgeCleared(name string) bool {
	switch name {
	case checklistitem.EdgeWorkOrder:
		return m.clearedwork_order
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *CheckListItemMutation) ClearEdge(name string) error {
	switch name {
	case checklistitem.EdgeWorkOrder:
		m.ClearWorkOrder()
		return nil
	}
	return fmt.Errorf("unknown CheckListItem unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *CheckListItemMutation) ResetEdge(name string) error {
	switch name {
	case checklistitem.EdgeWorkOrder:
		m.ResetWorkOrder()
		return nil
	}
	return fmt.Errorf("unknown CheckListItem edge %s", name)
}

// CheckListItemDefinitionMutation represents an operation that mutate the CheckListItemDefinitions
// nodes in the graph.
type CheckListItemDefinitionMutation struct {
	config
	op                     Op
	typ                    string
	id                     *int
	create_time            *time.Time
	update_time            *time.Time
	title                  *string
	_type                  *string
	index                  *int
	addindex               *int
	enum_values            *string
	help_text              *string
	clearedFields          map[string]bool
	work_order_type        *int
	clearedwork_order_type bool
}

var _ ent.Mutation = (*CheckListItemDefinitionMutation)(nil)

// newCheckListItemDefinitionMutation creates new mutation for $n.Name.
func newCheckListItemDefinitionMutation(c config, op Op) *CheckListItemDefinitionMutation {
	return &CheckListItemDefinitionMutation{
		config:        c,
		op:            op,
		typ:           TypeCheckListItemDefinition,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m CheckListItemDefinitionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m CheckListItemDefinitionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *CheckListItemDefinitionMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *CheckListItemDefinitionMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *CheckListItemDefinitionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *CheckListItemDefinitionMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *CheckListItemDefinitionMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *CheckListItemDefinitionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *CheckListItemDefinitionMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetTitle sets the title field.
func (m *CheckListItemDefinitionMutation) SetTitle(s string) {
	m.title = &s
}

// Title returns the title value in the mutation.
func (m *CheckListItemDefinitionMutation) Title() (r string, exists bool) {
	v := m.title
	if v == nil {
		return
	}
	return *v, true
}

// ResetTitle reset all changes of the title field.
func (m *CheckListItemDefinitionMutation) ResetTitle() {
	m.title = nil
}

// SetType sets the type field.
func (m *CheckListItemDefinitionMutation) SetType(s string) {
	m._type = &s
}

// GetType returns the type value in the mutation.
func (m *CheckListItemDefinitionMutation) GetType() (r string, exists bool) {
	v := m._type
	if v == nil {
		return
	}
	return *v, true
}

// ResetType reset all changes of the type field.
func (m *CheckListItemDefinitionMutation) ResetType() {
	m._type = nil
}

// SetIndex sets the index field.
func (m *CheckListItemDefinitionMutation) SetIndex(i int) {
	m.index = &i
	m.addindex = nil
}

// Index returns the index value in the mutation.
func (m *CheckListItemDefinitionMutation) Index() (r int, exists bool) {
	v := m.index
	if v == nil {
		return
	}
	return *v, true
}

// AddIndex adds i to index.
func (m *CheckListItemDefinitionMutation) AddIndex(i int) {
	if m.addindex != nil {
		*m.addindex += i
	} else {
		m.addindex = &i
	}
}

// AddedIndex returns the value that was added to the index field in this mutation.
func (m *CheckListItemDefinitionMutation) AddedIndex() (r int, exists bool) {
	v := m.addindex
	if v == nil {
		return
	}
	return *v, true
}

// ClearIndex clears the value of index.
func (m *CheckListItemDefinitionMutation) ClearIndex() {
	m.index = nil
	m.addindex = nil
	m.clearedFields[checklistitemdefinition.FieldIndex] = true
}

// IndexCleared returns if the field index was cleared in this mutation.
func (m *CheckListItemDefinitionMutation) IndexCleared() bool {
	return m.clearedFields[checklistitemdefinition.FieldIndex]
}

// ResetIndex reset all changes of the index field.
func (m *CheckListItemDefinitionMutation) ResetIndex() {
	m.index = nil
	m.addindex = nil
	delete(m.clearedFields, checklistitemdefinition.FieldIndex)
}

// SetEnumValues sets the enum_values field.
func (m *CheckListItemDefinitionMutation) SetEnumValues(s string) {
	m.enum_values = &s
}

// EnumValues returns the enum_values value in the mutation.
func (m *CheckListItemDefinitionMutation) EnumValues() (r string, exists bool) {
	v := m.enum_values
	if v == nil {
		return
	}
	return *v, true
}

// ClearEnumValues clears the value of enum_values.
func (m *CheckListItemDefinitionMutation) ClearEnumValues() {
	m.enum_values = nil
	m.clearedFields[checklistitemdefinition.FieldEnumValues] = true
}

// EnumValuesCleared returns if the field enum_values was cleared in this mutation.
func (m *CheckListItemDefinitionMutation) EnumValuesCleared() bool {
	return m.clearedFields[checklistitemdefinition.FieldEnumValues]
}

// ResetEnumValues reset all changes of the enum_values field.
func (m *CheckListItemDefinitionMutation) ResetEnumValues() {
	m.enum_values = nil
	delete(m.clearedFields, checklistitemdefinition.FieldEnumValues)
}

// SetHelpText sets the help_text field.
func (m *CheckListItemDefinitionMutation) SetHelpText(s string) {
	m.help_text = &s
}

// HelpText returns the help_text value in the mutation.
func (m *CheckListItemDefinitionMutation) HelpText() (r string, exists bool) {
	v := m.help_text
	if v == nil {
		return
	}
	return *v, true
}

// ClearHelpText clears the value of help_text.
func (m *CheckListItemDefinitionMutation) ClearHelpText() {
	m.help_text = nil
	m.clearedFields[checklistitemdefinition.FieldHelpText] = true
}

// HelpTextCleared returns if the field help_text was cleared in this mutation.
func (m *CheckListItemDefinitionMutation) HelpTextCleared() bool {
	return m.clearedFields[checklistitemdefinition.FieldHelpText]
}

// ResetHelpText reset all changes of the help_text field.
func (m *CheckListItemDefinitionMutation) ResetHelpText() {
	m.help_text = nil
	delete(m.clearedFields, checklistitemdefinition.FieldHelpText)
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (m *CheckListItemDefinitionMutation) SetWorkOrderTypeID(id int) {
	m.work_order_type = &id
}

// ClearWorkOrderType clears the work_order_type edge to WorkOrderType.
func (m *CheckListItemDefinitionMutation) ClearWorkOrderType() {
	m.clearedwork_order_type = true
}

// WorkOrderTypeCleared returns if the edge work_order_type was cleared.
func (m *CheckListItemDefinitionMutation) WorkOrderTypeCleared() bool {
	return m.clearedwork_order_type
}

// WorkOrderTypeID returns the work_order_type id in the mutation.
func (m *CheckListItemDefinitionMutation) WorkOrderTypeID() (id int, exists bool) {
	if m.work_order_type != nil {
		return *m.work_order_type, true
	}
	return
}

// WorkOrderTypeIDs returns the work_order_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// WorkOrderTypeID instead. It exists only for internal usage by the builders.
func (m *CheckListItemDefinitionMutation) WorkOrderTypeIDs() (ids []int) {
	if id := m.work_order_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetWorkOrderType reset all changes of the work_order_type edge.
func (m *CheckListItemDefinitionMutation) ResetWorkOrderType() {
	m.work_order_type = nil
	m.clearedwork_order_type = false
}

// Op returns the operation name.
func (m *CheckListItemDefinitionMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (CheckListItemDefinition).
func (m *CheckListItemDefinitionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *CheckListItemDefinitionMutation) Fields() []string {
	fields := make([]string, 0, 7)
	if m.create_time != nil {
		fields = append(fields, checklistitemdefinition.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, checklistitemdefinition.FieldUpdateTime)
	}
	if m.title != nil {
		fields = append(fields, checklistitemdefinition.FieldTitle)
	}
	if m._type != nil {
		fields = append(fields, checklistitemdefinition.FieldType)
	}
	if m.index != nil {
		fields = append(fields, checklistitemdefinition.FieldIndex)
	}
	if m.enum_values != nil {
		fields = append(fields, checklistitemdefinition.FieldEnumValues)
	}
	if m.help_text != nil {
		fields = append(fields, checklistitemdefinition.FieldHelpText)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *CheckListItemDefinitionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case checklistitemdefinition.FieldCreateTime:
		return m.CreateTime()
	case checklistitemdefinition.FieldUpdateTime:
		return m.UpdateTime()
	case checklistitemdefinition.FieldTitle:
		return m.Title()
	case checklistitemdefinition.FieldType:
		return m.GetType()
	case checklistitemdefinition.FieldIndex:
		return m.Index()
	case checklistitemdefinition.FieldEnumValues:
		return m.EnumValues()
	case checklistitemdefinition.FieldHelpText:
		return m.HelpText()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *CheckListItemDefinitionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case checklistitemdefinition.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case checklistitemdefinition.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case checklistitemdefinition.FieldTitle:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTitle(v)
		return nil
	case checklistitemdefinition.FieldType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
		return nil
	case checklistitemdefinition.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIndex(v)
		return nil
	case checklistitemdefinition.FieldEnumValues:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEnumValues(v)
		return nil
	case checklistitemdefinition.FieldHelpText:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetHelpText(v)
		return nil
	}
	return fmt.Errorf("unknown CheckListItemDefinition field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *CheckListItemDefinitionMutation) AddedFields() []string {
	var fields []string
	if m.addindex != nil {
		fields = append(fields, checklistitemdefinition.FieldIndex)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *CheckListItemDefinitionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case checklistitemdefinition.FieldIndex:
		return m.AddedIndex()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *CheckListItemDefinitionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case checklistitemdefinition.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIndex(v)
		return nil
	}
	return fmt.Errorf("unknown CheckListItemDefinition numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *CheckListItemDefinitionMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[checklistitemdefinition.FieldIndex] {
		fields = append(fields, checklistitemdefinition.FieldIndex)
	}
	if m.clearedFields[checklistitemdefinition.FieldEnumValues] {
		fields = append(fields, checklistitemdefinition.FieldEnumValues)
	}
	if m.clearedFields[checklistitemdefinition.FieldHelpText] {
		fields = append(fields, checklistitemdefinition.FieldHelpText)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *CheckListItemDefinitionMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *CheckListItemDefinitionMutation) ClearField(name string) error {
	switch name {
	case checklistitemdefinition.FieldIndex:
		m.ClearIndex()
		return nil
	case checklistitemdefinition.FieldEnumValues:
		m.ClearEnumValues()
		return nil
	case checklistitemdefinition.FieldHelpText:
		m.ClearHelpText()
		return nil
	}
	return fmt.Errorf("unknown CheckListItemDefinition nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *CheckListItemDefinitionMutation) ResetField(name string) error {
	switch name {
	case checklistitemdefinition.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case checklistitemdefinition.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case checklistitemdefinition.FieldTitle:
		m.ResetTitle()
		return nil
	case checklistitemdefinition.FieldType:
		m.ResetType()
		return nil
	case checklistitemdefinition.FieldIndex:
		m.ResetIndex()
		return nil
	case checklistitemdefinition.FieldEnumValues:
		m.ResetEnumValues()
		return nil
	case checklistitemdefinition.FieldHelpText:
		m.ResetHelpText()
		return nil
	}
	return fmt.Errorf("unknown CheckListItemDefinition field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *CheckListItemDefinitionMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.work_order_type != nil {
		edges = append(edges, checklistitemdefinition.EdgeWorkOrderType)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *CheckListItemDefinitionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case checklistitemdefinition.EdgeWorkOrderType:
		if id := m.work_order_type; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *CheckListItemDefinitionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *CheckListItemDefinitionMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *CheckListItemDefinitionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedwork_order_type {
		edges = append(edges, checklistitemdefinition.EdgeWorkOrderType)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *CheckListItemDefinitionMutation) EdgeCleared(name string) bool {
	switch name {
	case checklistitemdefinition.EdgeWorkOrderType:
		return m.clearedwork_order_type
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *CheckListItemDefinitionMutation) ClearEdge(name string) error {
	switch name {
	case checklistitemdefinition.EdgeWorkOrderType:
		m.ClearWorkOrderType()
		return nil
	}
	return fmt.Errorf("unknown CheckListItemDefinition unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *CheckListItemDefinitionMutation) ResetEdge(name string) error {
	switch name {
	case checklistitemdefinition.EdgeWorkOrderType:
		m.ResetWorkOrderType()
		return nil
	}
	return fmt.Errorf("unknown CheckListItemDefinition edge %s", name)
}

// CommentMutation represents an operation that mutate the Comments
// nodes in the graph.
type CommentMutation struct {
	config
	op            Op
	typ           string
	id            *int
	create_time   *time.Time
	update_time   *time.Time
	author_name   *string
	text          *string
	clearedFields map[string]bool
}

var _ ent.Mutation = (*CommentMutation)(nil)

// newCommentMutation creates new mutation for $n.Name.
func newCommentMutation(c config, op Op) *CommentMutation {
	return &CommentMutation{
		config:        c,
		op:            op,
		typ:           TypeComment,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m CommentMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m CommentMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *CommentMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *CommentMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *CommentMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *CommentMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *CommentMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *CommentMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *CommentMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetAuthorName sets the author_name field.
func (m *CommentMutation) SetAuthorName(s string) {
	m.author_name = &s
}

// AuthorName returns the author_name value in the mutation.
func (m *CommentMutation) AuthorName() (r string, exists bool) {
	v := m.author_name
	if v == nil {
		return
	}
	return *v, true
}

// ResetAuthorName reset all changes of the author_name field.
func (m *CommentMutation) ResetAuthorName() {
	m.author_name = nil
}

// SetText sets the text field.
func (m *CommentMutation) SetText(s string) {
	m.text = &s
}

// Text returns the text value in the mutation.
func (m *CommentMutation) Text() (r string, exists bool) {
	v := m.text
	if v == nil {
		return
	}
	return *v, true
}

// ResetText reset all changes of the text field.
func (m *CommentMutation) ResetText() {
	m.text = nil
}

// Op returns the operation name.
func (m *CommentMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Comment).
func (m *CommentMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *CommentMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.create_time != nil {
		fields = append(fields, comment.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, comment.FieldUpdateTime)
	}
	if m.author_name != nil {
		fields = append(fields, comment.FieldAuthorName)
	}
	if m.text != nil {
		fields = append(fields, comment.FieldText)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *CommentMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case comment.FieldCreateTime:
		return m.CreateTime()
	case comment.FieldUpdateTime:
		return m.UpdateTime()
	case comment.FieldAuthorName:
		return m.AuthorName()
	case comment.FieldText:
		return m.Text()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *CommentMutation) SetField(name string, value ent.Value) error {
	switch name {
	case comment.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case comment.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case comment.FieldAuthorName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetAuthorName(v)
		return nil
	case comment.FieldText:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetText(v)
		return nil
	}
	return fmt.Errorf("unknown Comment field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *CommentMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *CommentMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *CommentMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Comment numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *CommentMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *CommentMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *CommentMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Comment nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *CommentMutation) ResetField(name string) error {
	switch name {
	case comment.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case comment.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case comment.FieldAuthorName:
		m.ResetAuthorName()
		return nil
	case comment.FieldText:
		m.ResetText()
		return nil
	}
	return fmt.Errorf("unknown Comment field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *CommentMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *CommentMutation) AddedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *CommentMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *CommentMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *CommentMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *CommentMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *CommentMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Comment unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *CommentMutation) ResetEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Comment edge %s", name)
}

// CustomerMutation represents an operation that mutate the Customers
// nodes in the graph.
type CustomerMutation struct {
	config
	op              Op
	typ             string
	id              *int
	create_time     *time.Time
	update_time     *time.Time
	name            *string
	external_id     *string
	clearedFields   map[string]bool
	services        map[int]struct{}
	removedservices map[int]struct{}
}

var _ ent.Mutation = (*CustomerMutation)(nil)

// newCustomerMutation creates new mutation for $n.Name.
func newCustomerMutation(c config, op Op) *CustomerMutation {
	return &CustomerMutation{
		config:        c,
		op:            op,
		typ:           TypeCustomer,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m CustomerMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m CustomerMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *CustomerMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *CustomerMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *CustomerMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *CustomerMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *CustomerMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *CustomerMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *CustomerMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *CustomerMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *CustomerMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *CustomerMutation) ResetName() {
	m.name = nil
}

// SetExternalID sets the external_id field.
func (m *CustomerMutation) SetExternalID(s string) {
	m.external_id = &s
}

// ExternalID returns the external_id value in the mutation.
func (m *CustomerMutation) ExternalID() (r string, exists bool) {
	v := m.external_id
	if v == nil {
		return
	}
	return *v, true
}

// ClearExternalID clears the value of external_id.
func (m *CustomerMutation) ClearExternalID() {
	m.external_id = nil
	m.clearedFields[customer.FieldExternalID] = true
}

// ExternalIDCleared returns if the field external_id was cleared in this mutation.
func (m *CustomerMutation) ExternalIDCleared() bool {
	return m.clearedFields[customer.FieldExternalID]
}

// ResetExternalID reset all changes of the external_id field.
func (m *CustomerMutation) ResetExternalID() {
	m.external_id = nil
	delete(m.clearedFields, customer.FieldExternalID)
}

// AddServiceIDs adds the services edge to Service by ids.
func (m *CustomerMutation) AddServiceIDs(ids ...int) {
	if m.services == nil {
		m.services = make(map[int]struct{})
	}
	for i := range ids {
		m.services[ids[i]] = struct{}{}
	}
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (m *CustomerMutation) RemoveServiceIDs(ids ...int) {
	if m.removedservices == nil {
		m.removedservices = make(map[int]struct{})
	}
	for i := range ids {
		m.removedservices[ids[i]] = struct{}{}
	}
}

// RemovedServices returns the removed ids of services.
func (m *CustomerMutation) RemovedServicesIDs() (ids []int) {
	for id := range m.removedservices {
		ids = append(ids, id)
	}
	return
}

// ServicesIDs returns the services ids in the mutation.
func (m *CustomerMutation) ServicesIDs() (ids []int) {
	for id := range m.services {
		ids = append(ids, id)
	}
	return
}

// ResetServices reset all changes of the services edge.
func (m *CustomerMutation) ResetServices() {
	m.services = nil
	m.removedservices = nil
}

// Op returns the operation name.
func (m *CustomerMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Customer).
func (m *CustomerMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *CustomerMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.create_time != nil {
		fields = append(fields, customer.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, customer.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, customer.FieldName)
	}
	if m.external_id != nil {
		fields = append(fields, customer.FieldExternalID)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *CustomerMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case customer.FieldCreateTime:
		return m.CreateTime()
	case customer.FieldUpdateTime:
		return m.UpdateTime()
	case customer.FieldName:
		return m.Name()
	case customer.FieldExternalID:
		return m.ExternalID()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *CustomerMutation) SetField(name string, value ent.Value) error {
	switch name {
	case customer.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case customer.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case customer.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case customer.FieldExternalID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetExternalID(v)
		return nil
	}
	return fmt.Errorf("unknown Customer field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *CustomerMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *CustomerMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *CustomerMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Customer numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *CustomerMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[customer.FieldExternalID] {
		fields = append(fields, customer.FieldExternalID)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *CustomerMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *CustomerMutation) ClearField(name string) error {
	switch name {
	case customer.FieldExternalID:
		m.ClearExternalID()
		return nil
	}
	return fmt.Errorf("unknown Customer nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *CustomerMutation) ResetField(name string) error {
	switch name {
	case customer.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case customer.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case customer.FieldName:
		m.ResetName()
		return nil
	case customer.FieldExternalID:
		m.ResetExternalID()
		return nil
	}
	return fmt.Errorf("unknown Customer field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *CustomerMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.services != nil {
		edges = append(edges, customer.EdgeServices)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *CustomerMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case customer.EdgeServices:
		ids := make([]ent.Value, 0, len(m.services))
		for id := range m.services {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *CustomerMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedservices != nil {
		edges = append(edges, customer.EdgeServices)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *CustomerMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case customer.EdgeServices:
		ids := make([]ent.Value, 0, len(m.removedservices))
		for id := range m.removedservices {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *CustomerMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *CustomerMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *CustomerMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Customer unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *CustomerMutation) ResetEdge(name string) error {
	switch name {
	case customer.EdgeServices:
		m.ResetServices()
		return nil
	}
	return fmt.Errorf("unknown Customer edge %s", name)
}

// EquipmentMutation represents an operation that mutate the EquipmentSlice
// nodes in the graph.
type EquipmentMutation struct {
	config
	op                     Op
	typ                    string
	id                     *int
	create_time            *time.Time
	update_time            *time.Time
	name                   *string
	future_state           *string
	device_id              *string
	external_id            *string
	clearedFields          map[string]bool
	_type                  *int
	cleared_type           bool
	location               *int
	clearedlocation        bool
	parent_position        *int
	clearedparent_position bool
	positions              map[int]struct{}
	removedpositions       map[int]struct{}
	ports                  map[int]struct{}
	removedports           map[int]struct{}
	work_order             *int
	clearedwork_order      bool
	properties             map[int]struct{}
	removedproperties      map[int]struct{}
	files                  map[int]struct{}
	removedfiles           map[int]struct{}
	hyperlinks             map[int]struct{}
	removedhyperlinks      map[int]struct{}
}

var _ ent.Mutation = (*EquipmentMutation)(nil)

// newEquipmentMutation creates new mutation for $n.Name.
func newEquipmentMutation(c config, op Op) *EquipmentMutation {
	return &EquipmentMutation{
		config:        c,
		op:            op,
		typ:           TypeEquipment,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EquipmentMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EquipmentMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *EquipmentMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *EquipmentMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *EquipmentMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *EquipmentMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *EquipmentMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *EquipmentMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *EquipmentMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *EquipmentMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *EquipmentMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *EquipmentMutation) ResetName() {
	m.name = nil
}

// SetFutureState sets the future_state field.
func (m *EquipmentMutation) SetFutureState(s string) {
	m.future_state = &s
}

// FutureState returns the future_state value in the mutation.
func (m *EquipmentMutation) FutureState() (r string, exists bool) {
	v := m.future_state
	if v == nil {
		return
	}
	return *v, true
}

// ClearFutureState clears the value of future_state.
func (m *EquipmentMutation) ClearFutureState() {
	m.future_state = nil
	m.clearedFields[equipment.FieldFutureState] = true
}

// FutureStateCleared returns if the field future_state was cleared in this mutation.
func (m *EquipmentMutation) FutureStateCleared() bool {
	return m.clearedFields[equipment.FieldFutureState]
}

// ResetFutureState reset all changes of the future_state field.
func (m *EquipmentMutation) ResetFutureState() {
	m.future_state = nil
	delete(m.clearedFields, equipment.FieldFutureState)
}

// SetDeviceID sets the device_id field.
func (m *EquipmentMutation) SetDeviceID(s string) {
	m.device_id = &s
}

// DeviceID returns the device_id value in the mutation.
func (m *EquipmentMutation) DeviceID() (r string, exists bool) {
	v := m.device_id
	if v == nil {
		return
	}
	return *v, true
}

// ClearDeviceID clears the value of device_id.
func (m *EquipmentMutation) ClearDeviceID() {
	m.device_id = nil
	m.clearedFields[equipment.FieldDeviceID] = true
}

// DeviceIDCleared returns if the field device_id was cleared in this mutation.
func (m *EquipmentMutation) DeviceIDCleared() bool {
	return m.clearedFields[equipment.FieldDeviceID]
}

// ResetDeviceID reset all changes of the device_id field.
func (m *EquipmentMutation) ResetDeviceID() {
	m.device_id = nil
	delete(m.clearedFields, equipment.FieldDeviceID)
}

// SetExternalID sets the external_id field.
func (m *EquipmentMutation) SetExternalID(s string) {
	m.external_id = &s
}

// ExternalID returns the external_id value in the mutation.
func (m *EquipmentMutation) ExternalID() (r string, exists bool) {
	v := m.external_id
	if v == nil {
		return
	}
	return *v, true
}

// ClearExternalID clears the value of external_id.
func (m *EquipmentMutation) ClearExternalID() {
	m.external_id = nil
	m.clearedFields[equipment.FieldExternalID] = true
}

// ExternalIDCleared returns if the field external_id was cleared in this mutation.
func (m *EquipmentMutation) ExternalIDCleared() bool {
	return m.clearedFields[equipment.FieldExternalID]
}

// ResetExternalID reset all changes of the external_id field.
func (m *EquipmentMutation) ResetExternalID() {
	m.external_id = nil
	delete(m.clearedFields, equipment.FieldExternalID)
}

// SetTypeID sets the type edge to EquipmentType by id.
func (m *EquipmentMutation) SetTypeID(id int) {
	m._type = &id
}

// ClearType clears the type edge to EquipmentType.
func (m *EquipmentMutation) ClearType() {
	m.cleared_type = true
}

// TypeCleared returns if the edge type was cleared.
func (m *EquipmentMutation) TypeCleared() bool {
	return m.cleared_type
}

// TypeID returns the type id in the mutation.
func (m *EquipmentMutation) TypeID() (id int, exists bool) {
	if m._type != nil {
		return *m._type, true
	}
	return
}

// TypeIDs returns the type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// TypeID instead. It exists only for internal usage by the builders.
func (m *EquipmentMutation) TypeIDs() (ids []int) {
	if id := m._type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetType reset all changes of the type edge.
func (m *EquipmentMutation) ResetType() {
	m._type = nil
	m.cleared_type = false
}

// SetLocationID sets the location edge to Location by id.
func (m *EquipmentMutation) SetLocationID(id int) {
	m.location = &id
}

// ClearLocation clears the location edge to Location.
func (m *EquipmentMutation) ClearLocation() {
	m.clearedlocation = true
}

// LocationCleared returns if the edge location was cleared.
func (m *EquipmentMutation) LocationCleared() bool {
	return m.clearedlocation
}

// LocationID returns the location id in the mutation.
func (m *EquipmentMutation) LocationID() (id int, exists bool) {
	if m.location != nil {
		return *m.location, true
	}
	return
}

// LocationIDs returns the location ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LocationID instead. It exists only for internal usage by the builders.
func (m *EquipmentMutation) LocationIDs() (ids []int) {
	if id := m.location; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLocation reset all changes of the location edge.
func (m *EquipmentMutation) ResetLocation() {
	m.location = nil
	m.clearedlocation = false
}

// SetParentPositionID sets the parent_position edge to EquipmentPosition by id.
func (m *EquipmentMutation) SetParentPositionID(id int) {
	m.parent_position = &id
}

// ClearParentPosition clears the parent_position edge to EquipmentPosition.
func (m *EquipmentMutation) ClearParentPosition() {
	m.clearedparent_position = true
}

// ParentPositionCleared returns if the edge parent_position was cleared.
func (m *EquipmentMutation) ParentPositionCleared() bool {
	return m.clearedparent_position
}

// ParentPositionID returns the parent_position id in the mutation.
func (m *EquipmentMutation) ParentPositionID() (id int, exists bool) {
	if m.parent_position != nil {
		return *m.parent_position, true
	}
	return
}

// ParentPositionIDs returns the parent_position ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ParentPositionID instead. It exists only for internal usage by the builders.
func (m *EquipmentMutation) ParentPositionIDs() (ids []int) {
	if id := m.parent_position; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetParentPosition reset all changes of the parent_position edge.
func (m *EquipmentMutation) ResetParentPosition() {
	m.parent_position = nil
	m.clearedparent_position = false
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (m *EquipmentMutation) AddPositionIDs(ids ...int) {
	if m.positions == nil {
		m.positions = make(map[int]struct{})
	}
	for i := range ids {
		m.positions[ids[i]] = struct{}{}
	}
}

// RemovePositionIDs removes the positions edge to EquipmentPosition by ids.
func (m *EquipmentMutation) RemovePositionIDs(ids ...int) {
	if m.removedpositions == nil {
		m.removedpositions = make(map[int]struct{})
	}
	for i := range ids {
		m.removedpositions[ids[i]] = struct{}{}
	}
}

// RemovedPositions returns the removed ids of positions.
func (m *EquipmentMutation) RemovedPositionsIDs() (ids []int) {
	for id := range m.removedpositions {
		ids = append(ids, id)
	}
	return
}

// PositionsIDs returns the positions ids in the mutation.
func (m *EquipmentMutation) PositionsIDs() (ids []int) {
	for id := range m.positions {
		ids = append(ids, id)
	}
	return
}

// ResetPositions reset all changes of the positions edge.
func (m *EquipmentMutation) ResetPositions() {
	m.positions = nil
	m.removedpositions = nil
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (m *EquipmentMutation) AddPortIDs(ids ...int) {
	if m.ports == nil {
		m.ports = make(map[int]struct{})
	}
	for i := range ids {
		m.ports[ids[i]] = struct{}{}
	}
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (m *EquipmentMutation) RemovePortIDs(ids ...int) {
	if m.removedports == nil {
		m.removedports = make(map[int]struct{})
	}
	for i := range ids {
		m.removedports[ids[i]] = struct{}{}
	}
}

// RemovedPorts returns the removed ids of ports.
func (m *EquipmentMutation) RemovedPortsIDs() (ids []int) {
	for id := range m.removedports {
		ids = append(ids, id)
	}
	return
}

// PortsIDs returns the ports ids in the mutation.
func (m *EquipmentMutation) PortsIDs() (ids []int) {
	for id := range m.ports {
		ids = append(ids, id)
	}
	return
}

// ResetPorts reset all changes of the ports edge.
func (m *EquipmentMutation) ResetPorts() {
	m.ports = nil
	m.removedports = nil
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (m *EquipmentMutation) SetWorkOrderID(id int) {
	m.work_order = &id
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (m *EquipmentMutation) ClearWorkOrder() {
	m.clearedwork_order = true
}

// WorkOrderCleared returns if the edge work_order was cleared.
func (m *EquipmentMutation) WorkOrderCleared() bool {
	return m.clearedwork_order
}

// WorkOrderID returns the work_order id in the mutation.
func (m *EquipmentMutation) WorkOrderID() (id int, exists bool) {
	if m.work_order != nil {
		return *m.work_order, true
	}
	return
}

// WorkOrderIDs returns the work_order ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// WorkOrderID instead. It exists only for internal usage by the builders.
func (m *EquipmentMutation) WorkOrderIDs() (ids []int) {
	if id := m.work_order; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetWorkOrder reset all changes of the work_order edge.
func (m *EquipmentMutation) ResetWorkOrder() {
	m.work_order = nil
	m.clearedwork_order = false
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (m *EquipmentMutation) AddPropertyIDs(ids ...int) {
	if m.properties == nil {
		m.properties = make(map[int]struct{})
	}
	for i := range ids {
		m.properties[ids[i]] = struct{}{}
	}
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (m *EquipmentMutation) RemovePropertyIDs(ids ...int) {
	if m.removedproperties == nil {
		m.removedproperties = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperties[ids[i]] = struct{}{}
	}
}

// RemovedProperties returns the removed ids of properties.
func (m *EquipmentMutation) RemovedPropertiesIDs() (ids []int) {
	for id := range m.removedproperties {
		ids = append(ids, id)
	}
	return
}

// PropertiesIDs returns the properties ids in the mutation.
func (m *EquipmentMutation) PropertiesIDs() (ids []int) {
	for id := range m.properties {
		ids = append(ids, id)
	}
	return
}

// ResetProperties reset all changes of the properties edge.
func (m *EquipmentMutation) ResetProperties() {
	m.properties = nil
	m.removedproperties = nil
}

// AddFileIDs adds the files edge to File by ids.
func (m *EquipmentMutation) AddFileIDs(ids ...int) {
	if m.files == nil {
		m.files = make(map[int]struct{})
	}
	for i := range ids {
		m.files[ids[i]] = struct{}{}
	}
}

// RemoveFileIDs removes the files edge to File by ids.
func (m *EquipmentMutation) RemoveFileIDs(ids ...int) {
	if m.removedfiles == nil {
		m.removedfiles = make(map[int]struct{})
	}
	for i := range ids {
		m.removedfiles[ids[i]] = struct{}{}
	}
}

// RemovedFiles returns the removed ids of files.
func (m *EquipmentMutation) RemovedFilesIDs() (ids []int) {
	for id := range m.removedfiles {
		ids = append(ids, id)
	}
	return
}

// FilesIDs returns the files ids in the mutation.
func (m *EquipmentMutation) FilesIDs() (ids []int) {
	for id := range m.files {
		ids = append(ids, id)
	}
	return
}

// ResetFiles reset all changes of the files edge.
func (m *EquipmentMutation) ResetFiles() {
	m.files = nil
	m.removedfiles = nil
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (m *EquipmentMutation) AddHyperlinkIDs(ids ...int) {
	if m.hyperlinks == nil {
		m.hyperlinks = make(map[int]struct{})
	}
	for i := range ids {
		m.hyperlinks[ids[i]] = struct{}{}
	}
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (m *EquipmentMutation) RemoveHyperlinkIDs(ids ...int) {
	if m.removedhyperlinks == nil {
		m.removedhyperlinks = make(map[int]struct{})
	}
	for i := range ids {
		m.removedhyperlinks[ids[i]] = struct{}{}
	}
}

// RemovedHyperlinks returns the removed ids of hyperlinks.
func (m *EquipmentMutation) RemovedHyperlinksIDs() (ids []int) {
	for id := range m.removedhyperlinks {
		ids = append(ids, id)
	}
	return
}

// HyperlinksIDs returns the hyperlinks ids in the mutation.
func (m *EquipmentMutation) HyperlinksIDs() (ids []int) {
	for id := range m.hyperlinks {
		ids = append(ids, id)
	}
	return
}

// ResetHyperlinks reset all changes of the hyperlinks edge.
func (m *EquipmentMutation) ResetHyperlinks() {
	m.hyperlinks = nil
	m.removedhyperlinks = nil
}

// Op returns the operation name.
func (m *EquipmentMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Equipment).
func (m *EquipmentMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *EquipmentMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m.create_time != nil {
		fields = append(fields, equipment.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, equipment.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, equipment.FieldName)
	}
	if m.future_state != nil {
		fields = append(fields, equipment.FieldFutureState)
	}
	if m.device_id != nil {
		fields = append(fields, equipment.FieldDeviceID)
	}
	if m.external_id != nil {
		fields = append(fields, equipment.FieldExternalID)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *EquipmentMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case equipment.FieldCreateTime:
		return m.CreateTime()
	case equipment.FieldUpdateTime:
		return m.UpdateTime()
	case equipment.FieldName:
		return m.Name()
	case equipment.FieldFutureState:
		return m.FutureState()
	case equipment.FieldDeviceID:
		return m.DeviceID()
	case equipment.FieldExternalID:
		return m.ExternalID()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentMutation) SetField(name string, value ent.Value) error {
	switch name {
	case equipment.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case equipment.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case equipment.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case equipment.FieldFutureState:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFutureState(v)
		return nil
	case equipment.FieldDeviceID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDeviceID(v)
		return nil
	case equipment.FieldExternalID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetExternalID(v)
		return nil
	}
	return fmt.Errorf("unknown Equipment field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *EquipmentMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *EquipmentMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Equipment numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *EquipmentMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[equipment.FieldFutureState] {
		fields = append(fields, equipment.FieldFutureState)
	}
	if m.clearedFields[equipment.FieldDeviceID] {
		fields = append(fields, equipment.FieldDeviceID)
	}
	if m.clearedFields[equipment.FieldExternalID] {
		fields = append(fields, equipment.FieldExternalID)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *EquipmentMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *EquipmentMutation) ClearField(name string) error {
	switch name {
	case equipment.FieldFutureState:
		m.ClearFutureState()
		return nil
	case equipment.FieldDeviceID:
		m.ClearDeviceID()
		return nil
	case equipment.FieldExternalID:
		m.ClearExternalID()
		return nil
	}
	return fmt.Errorf("unknown Equipment nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *EquipmentMutation) ResetField(name string) error {
	switch name {
	case equipment.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case equipment.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case equipment.FieldName:
		m.ResetName()
		return nil
	case equipment.FieldFutureState:
		m.ResetFutureState()
		return nil
	case equipment.FieldDeviceID:
		m.ResetDeviceID()
		return nil
	case equipment.FieldExternalID:
		m.ResetExternalID()
		return nil
	}
	return fmt.Errorf("unknown Equipment field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *EquipmentMutation) AddedEdges() []string {
	edges := make([]string, 0, 9)
	if m._type != nil {
		edges = append(edges, equipment.EdgeType)
	}
	if m.location != nil {
		edges = append(edges, equipment.EdgeLocation)
	}
	if m.parent_position != nil {
		edges = append(edges, equipment.EdgeParentPosition)
	}
	if m.positions != nil {
		edges = append(edges, equipment.EdgePositions)
	}
	if m.ports != nil {
		edges = append(edges, equipment.EdgePorts)
	}
	if m.work_order != nil {
		edges = append(edges, equipment.EdgeWorkOrder)
	}
	if m.properties != nil {
		edges = append(edges, equipment.EdgeProperties)
	}
	if m.files != nil {
		edges = append(edges, equipment.EdgeFiles)
	}
	if m.hyperlinks != nil {
		edges = append(edges, equipment.EdgeHyperlinks)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *EquipmentMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case equipment.EdgeType:
		if id := m._type; id != nil {
			return []ent.Value{*id}
		}
	case equipment.EdgeLocation:
		if id := m.location; id != nil {
			return []ent.Value{*id}
		}
	case equipment.EdgeParentPosition:
		if id := m.parent_position; id != nil {
			return []ent.Value{*id}
		}
	case equipment.EdgePositions:
		ids := make([]ent.Value, 0, len(m.positions))
		for id := range m.positions {
			ids = append(ids, id)
		}
		return ids
	case equipment.EdgePorts:
		ids := make([]ent.Value, 0, len(m.ports))
		for id := range m.ports {
			ids = append(ids, id)
		}
		return ids
	case equipment.EdgeWorkOrder:
		if id := m.work_order; id != nil {
			return []ent.Value{*id}
		}
	case equipment.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.properties))
		for id := range m.properties {
			ids = append(ids, id)
		}
		return ids
	case equipment.EdgeFiles:
		ids := make([]ent.Value, 0, len(m.files))
		for id := range m.files {
			ids = append(ids, id)
		}
		return ids
	case equipment.EdgeHyperlinks:
		ids := make([]ent.Value, 0, len(m.hyperlinks))
		for id := range m.hyperlinks {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *EquipmentMutation) RemovedEdges() []string {
	edges := make([]string, 0, 9)
	if m.removedpositions != nil {
		edges = append(edges, equipment.EdgePositions)
	}
	if m.removedports != nil {
		edges = append(edges, equipment.EdgePorts)
	}
	if m.removedproperties != nil {
		edges = append(edges, equipment.EdgeProperties)
	}
	if m.removedfiles != nil {
		edges = append(edges, equipment.EdgeFiles)
	}
	if m.removedhyperlinks != nil {
		edges = append(edges, equipment.EdgeHyperlinks)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *EquipmentMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case equipment.EdgePositions:
		ids := make([]ent.Value, 0, len(m.removedpositions))
		for id := range m.removedpositions {
			ids = append(ids, id)
		}
		return ids
	case equipment.EdgePorts:
		ids := make([]ent.Value, 0, len(m.removedports))
		for id := range m.removedports {
			ids = append(ids, id)
		}
		return ids
	case equipment.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.removedproperties))
		for id := range m.removedproperties {
			ids = append(ids, id)
		}
		return ids
	case equipment.EdgeFiles:
		ids := make([]ent.Value, 0, len(m.removedfiles))
		for id := range m.removedfiles {
			ids = append(ids, id)
		}
		return ids
	case equipment.EdgeHyperlinks:
		ids := make([]ent.Value, 0, len(m.removedhyperlinks))
		for id := range m.removedhyperlinks {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *EquipmentMutation) ClearedEdges() []string {
	edges := make([]string, 0, 9)
	if m.cleared_type {
		edges = append(edges, equipment.EdgeType)
	}
	if m.clearedlocation {
		edges = append(edges, equipment.EdgeLocation)
	}
	if m.clearedparent_position {
		edges = append(edges, equipment.EdgeParentPosition)
	}
	if m.clearedwork_order {
		edges = append(edges, equipment.EdgeWorkOrder)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *EquipmentMutation) EdgeCleared(name string) bool {
	switch name {
	case equipment.EdgeType:
		return m.cleared_type
	case equipment.EdgeLocation:
		return m.clearedlocation
	case equipment.EdgeParentPosition:
		return m.clearedparent_position
	case equipment.EdgeWorkOrder:
		return m.clearedwork_order
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *EquipmentMutation) ClearEdge(name string) error {
	switch name {
	case equipment.EdgeType:
		m.ClearType()
		return nil
	case equipment.EdgeLocation:
		m.ClearLocation()
		return nil
	case equipment.EdgeParentPosition:
		m.ClearParentPosition()
		return nil
	case equipment.EdgeWorkOrder:
		m.ClearWorkOrder()
		return nil
	}
	return fmt.Errorf("unknown Equipment unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *EquipmentMutation) ResetEdge(name string) error {
	switch name {
	case equipment.EdgeType:
		m.ResetType()
		return nil
	case equipment.EdgeLocation:
		m.ResetLocation()
		return nil
	case equipment.EdgeParentPosition:
		m.ResetParentPosition()
		return nil
	case equipment.EdgePositions:
		m.ResetPositions()
		return nil
	case equipment.EdgePorts:
		m.ResetPorts()
		return nil
	case equipment.EdgeWorkOrder:
		m.ResetWorkOrder()
		return nil
	case equipment.EdgeProperties:
		m.ResetProperties()
		return nil
	case equipment.EdgeFiles:
		m.ResetFiles()
		return nil
	case equipment.EdgeHyperlinks:
		m.ResetHyperlinks()
		return nil
	}
	return fmt.Errorf("unknown Equipment edge %s", name)
}

// EquipmentCategoryMutation represents an operation that mutate the EquipmentCategories
// nodes in the graph.
type EquipmentCategoryMutation struct {
	config
	op            Op
	typ           string
	id            *int
	create_time   *time.Time
	update_time   *time.Time
	name          *string
	clearedFields map[string]bool
	types         map[int]struct{}
	removedtypes  map[int]struct{}
}

var _ ent.Mutation = (*EquipmentCategoryMutation)(nil)

// newEquipmentCategoryMutation creates new mutation for $n.Name.
func newEquipmentCategoryMutation(c config, op Op) *EquipmentCategoryMutation {
	return &EquipmentCategoryMutation{
		config:        c,
		op:            op,
		typ:           TypeEquipmentCategory,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EquipmentCategoryMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EquipmentCategoryMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *EquipmentCategoryMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *EquipmentCategoryMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *EquipmentCategoryMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *EquipmentCategoryMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *EquipmentCategoryMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *EquipmentCategoryMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *EquipmentCategoryMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *EquipmentCategoryMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *EquipmentCategoryMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *EquipmentCategoryMutation) ResetName() {
	m.name = nil
}

// AddTypeIDs adds the types edge to EquipmentType by ids.
func (m *EquipmentCategoryMutation) AddTypeIDs(ids ...int) {
	if m.types == nil {
		m.types = make(map[int]struct{})
	}
	for i := range ids {
		m.types[ids[i]] = struct{}{}
	}
}

// RemoveTypeIDs removes the types edge to EquipmentType by ids.
func (m *EquipmentCategoryMutation) RemoveTypeIDs(ids ...int) {
	if m.removedtypes == nil {
		m.removedtypes = make(map[int]struct{})
	}
	for i := range ids {
		m.removedtypes[ids[i]] = struct{}{}
	}
}

// RemovedTypes returns the removed ids of types.
func (m *EquipmentCategoryMutation) RemovedTypesIDs() (ids []int) {
	for id := range m.removedtypes {
		ids = append(ids, id)
	}
	return
}

// TypesIDs returns the types ids in the mutation.
func (m *EquipmentCategoryMutation) TypesIDs() (ids []int) {
	for id := range m.types {
		ids = append(ids, id)
	}
	return
}

// ResetTypes reset all changes of the types edge.
func (m *EquipmentCategoryMutation) ResetTypes() {
	m.types = nil
	m.removedtypes = nil
}

// Op returns the operation name.
func (m *EquipmentCategoryMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (EquipmentCategory).
func (m *EquipmentCategoryMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *EquipmentCategoryMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.create_time != nil {
		fields = append(fields, equipmentcategory.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, equipmentcategory.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, equipmentcategory.FieldName)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *EquipmentCategoryMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case equipmentcategory.FieldCreateTime:
		return m.CreateTime()
	case equipmentcategory.FieldUpdateTime:
		return m.UpdateTime()
	case equipmentcategory.FieldName:
		return m.Name()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentCategoryMutation) SetField(name string, value ent.Value) error {
	switch name {
	case equipmentcategory.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case equipmentcategory.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case equipmentcategory.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	}
	return fmt.Errorf("unknown EquipmentCategory field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *EquipmentCategoryMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *EquipmentCategoryMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentCategoryMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown EquipmentCategory numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *EquipmentCategoryMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *EquipmentCategoryMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *EquipmentCategoryMutation) ClearField(name string) error {
	return fmt.Errorf("unknown EquipmentCategory nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *EquipmentCategoryMutation) ResetField(name string) error {
	switch name {
	case equipmentcategory.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case equipmentcategory.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case equipmentcategory.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown EquipmentCategory field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *EquipmentCategoryMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.types != nil {
		edges = append(edges, equipmentcategory.EdgeTypes)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *EquipmentCategoryMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case equipmentcategory.EdgeTypes:
		ids := make([]ent.Value, 0, len(m.types))
		for id := range m.types {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *EquipmentCategoryMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedtypes != nil {
		edges = append(edges, equipmentcategory.EdgeTypes)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *EquipmentCategoryMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case equipmentcategory.EdgeTypes:
		ids := make([]ent.Value, 0, len(m.removedtypes))
		for id := range m.removedtypes {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *EquipmentCategoryMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *EquipmentCategoryMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *EquipmentCategoryMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown EquipmentCategory unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *EquipmentCategoryMutation) ResetEdge(name string) error {
	switch name {
	case equipmentcategory.EdgeTypes:
		m.ResetTypes()
		return nil
	}
	return fmt.Errorf("unknown EquipmentCategory edge %s", name)
}

// EquipmentPortMutation represents an operation that mutate the EquipmentPorts
// nodes in the graph.
type EquipmentPortMutation struct {
	config
	op                Op
	typ               string
	id                *int
	create_time       *time.Time
	update_time       *time.Time
	clearedFields     map[string]bool
	definition        *int
	cleareddefinition bool
	parent            *int
	clearedparent     bool
	link              *int
	clearedlink       bool
	properties        map[int]struct{}
	removedproperties map[int]struct{}
	endpoints         map[int]struct{}
	removedendpoints  map[int]struct{}
}

var _ ent.Mutation = (*EquipmentPortMutation)(nil)

// newEquipmentPortMutation creates new mutation for $n.Name.
func newEquipmentPortMutation(c config, op Op) *EquipmentPortMutation {
	return &EquipmentPortMutation{
		config:        c,
		op:            op,
		typ:           TypeEquipmentPort,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EquipmentPortMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EquipmentPortMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *EquipmentPortMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *EquipmentPortMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *EquipmentPortMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *EquipmentPortMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *EquipmentPortMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *EquipmentPortMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *EquipmentPortMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetDefinitionID sets the definition edge to EquipmentPortDefinition by id.
func (m *EquipmentPortMutation) SetDefinitionID(id int) {
	m.definition = &id
}

// ClearDefinition clears the definition edge to EquipmentPortDefinition.
func (m *EquipmentPortMutation) ClearDefinition() {
	m.cleareddefinition = true
}

// DefinitionCleared returns if the edge definition was cleared.
func (m *EquipmentPortMutation) DefinitionCleared() bool {
	return m.cleareddefinition
}

// DefinitionID returns the definition id in the mutation.
func (m *EquipmentPortMutation) DefinitionID() (id int, exists bool) {
	if m.definition != nil {
		return *m.definition, true
	}
	return
}

// DefinitionIDs returns the definition ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// DefinitionID instead. It exists only for internal usage by the builders.
func (m *EquipmentPortMutation) DefinitionIDs() (ids []int) {
	if id := m.definition; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetDefinition reset all changes of the definition edge.
func (m *EquipmentPortMutation) ResetDefinition() {
	m.definition = nil
	m.cleareddefinition = false
}

// SetParentID sets the parent edge to Equipment by id.
func (m *EquipmentPortMutation) SetParentID(id int) {
	m.parent = &id
}

// ClearParent clears the parent edge to Equipment.
func (m *EquipmentPortMutation) ClearParent() {
	m.clearedparent = true
}

// ParentCleared returns if the edge parent was cleared.
func (m *EquipmentPortMutation) ParentCleared() bool {
	return m.clearedparent
}

// ParentID returns the parent id in the mutation.
func (m *EquipmentPortMutation) ParentID() (id int, exists bool) {
	if m.parent != nil {
		return *m.parent, true
	}
	return
}

// ParentIDs returns the parent ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ParentID instead. It exists only for internal usage by the builders.
func (m *EquipmentPortMutation) ParentIDs() (ids []int) {
	if id := m.parent; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetParent reset all changes of the parent edge.
func (m *EquipmentPortMutation) ResetParent() {
	m.parent = nil
	m.clearedparent = false
}

// SetLinkID sets the link edge to Link by id.
func (m *EquipmentPortMutation) SetLinkID(id int) {
	m.link = &id
}

// ClearLink clears the link edge to Link.
func (m *EquipmentPortMutation) ClearLink() {
	m.clearedlink = true
}

// LinkCleared returns if the edge link was cleared.
func (m *EquipmentPortMutation) LinkCleared() bool {
	return m.clearedlink
}

// LinkID returns the link id in the mutation.
func (m *EquipmentPortMutation) LinkID() (id int, exists bool) {
	if m.link != nil {
		return *m.link, true
	}
	return
}

// LinkIDs returns the link ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LinkID instead. It exists only for internal usage by the builders.
func (m *EquipmentPortMutation) LinkIDs() (ids []int) {
	if id := m.link; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLink reset all changes of the link edge.
func (m *EquipmentPortMutation) ResetLink() {
	m.link = nil
	m.clearedlink = false
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (m *EquipmentPortMutation) AddPropertyIDs(ids ...int) {
	if m.properties == nil {
		m.properties = make(map[int]struct{})
	}
	for i := range ids {
		m.properties[ids[i]] = struct{}{}
	}
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (m *EquipmentPortMutation) RemovePropertyIDs(ids ...int) {
	if m.removedproperties == nil {
		m.removedproperties = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperties[ids[i]] = struct{}{}
	}
}

// RemovedProperties returns the removed ids of properties.
func (m *EquipmentPortMutation) RemovedPropertiesIDs() (ids []int) {
	for id := range m.removedproperties {
		ids = append(ids, id)
	}
	return
}

// PropertiesIDs returns the properties ids in the mutation.
func (m *EquipmentPortMutation) PropertiesIDs() (ids []int) {
	for id := range m.properties {
		ids = append(ids, id)
	}
	return
}

// ResetProperties reset all changes of the properties edge.
func (m *EquipmentPortMutation) ResetProperties() {
	m.properties = nil
	m.removedproperties = nil
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (m *EquipmentPortMutation) AddEndpointIDs(ids ...int) {
	if m.endpoints == nil {
		m.endpoints = make(map[int]struct{})
	}
	for i := range ids {
		m.endpoints[ids[i]] = struct{}{}
	}
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (m *EquipmentPortMutation) RemoveEndpointIDs(ids ...int) {
	if m.removedendpoints == nil {
		m.removedendpoints = make(map[int]struct{})
	}
	for i := range ids {
		m.removedendpoints[ids[i]] = struct{}{}
	}
}

// RemovedEndpoints returns the removed ids of endpoints.
func (m *EquipmentPortMutation) RemovedEndpointsIDs() (ids []int) {
	for id := range m.removedendpoints {
		ids = append(ids, id)
	}
	return
}

// EndpointsIDs returns the endpoints ids in the mutation.
func (m *EquipmentPortMutation) EndpointsIDs() (ids []int) {
	for id := range m.endpoints {
		ids = append(ids, id)
	}
	return
}

// ResetEndpoints reset all changes of the endpoints edge.
func (m *EquipmentPortMutation) ResetEndpoints() {
	m.endpoints = nil
	m.removedendpoints = nil
}

// Op returns the operation name.
func (m *EquipmentPortMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (EquipmentPort).
func (m *EquipmentPortMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *EquipmentPortMutation) Fields() []string {
	fields := make([]string, 0, 2)
	if m.create_time != nil {
		fields = append(fields, equipmentport.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, equipmentport.FieldUpdateTime)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *EquipmentPortMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case equipmentport.FieldCreateTime:
		return m.CreateTime()
	case equipmentport.FieldUpdateTime:
		return m.UpdateTime()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentPortMutation) SetField(name string, value ent.Value) error {
	switch name {
	case equipmentport.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case equipmentport.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	}
	return fmt.Errorf("unknown EquipmentPort field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *EquipmentPortMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *EquipmentPortMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentPortMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown EquipmentPort numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *EquipmentPortMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *EquipmentPortMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *EquipmentPortMutation) ClearField(name string) error {
	return fmt.Errorf("unknown EquipmentPort nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *EquipmentPortMutation) ResetField(name string) error {
	switch name {
	case equipmentport.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case equipmentport.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPort field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *EquipmentPortMutation) AddedEdges() []string {
	edges := make([]string, 0, 5)
	if m.definition != nil {
		edges = append(edges, equipmentport.EdgeDefinition)
	}
	if m.parent != nil {
		edges = append(edges, equipmentport.EdgeParent)
	}
	if m.link != nil {
		edges = append(edges, equipmentport.EdgeLink)
	}
	if m.properties != nil {
		edges = append(edges, equipmentport.EdgeProperties)
	}
	if m.endpoints != nil {
		edges = append(edges, equipmentport.EdgeEndpoints)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *EquipmentPortMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case equipmentport.EdgeDefinition:
		if id := m.definition; id != nil {
			return []ent.Value{*id}
		}
	case equipmentport.EdgeParent:
		if id := m.parent; id != nil {
			return []ent.Value{*id}
		}
	case equipmentport.EdgeLink:
		if id := m.link; id != nil {
			return []ent.Value{*id}
		}
	case equipmentport.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.properties))
		for id := range m.properties {
			ids = append(ids, id)
		}
		return ids
	case equipmentport.EdgeEndpoints:
		ids := make([]ent.Value, 0, len(m.endpoints))
		for id := range m.endpoints {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *EquipmentPortMutation) RemovedEdges() []string {
	edges := make([]string, 0, 5)
	if m.removedproperties != nil {
		edges = append(edges, equipmentport.EdgeProperties)
	}
	if m.removedendpoints != nil {
		edges = append(edges, equipmentport.EdgeEndpoints)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *EquipmentPortMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case equipmentport.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.removedproperties))
		for id := range m.removedproperties {
			ids = append(ids, id)
		}
		return ids
	case equipmentport.EdgeEndpoints:
		ids := make([]ent.Value, 0, len(m.removedendpoints))
		for id := range m.removedendpoints {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *EquipmentPortMutation) ClearedEdges() []string {
	edges := make([]string, 0, 5)
	if m.cleareddefinition {
		edges = append(edges, equipmentport.EdgeDefinition)
	}
	if m.clearedparent {
		edges = append(edges, equipmentport.EdgeParent)
	}
	if m.clearedlink {
		edges = append(edges, equipmentport.EdgeLink)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *EquipmentPortMutation) EdgeCleared(name string) bool {
	switch name {
	case equipmentport.EdgeDefinition:
		return m.cleareddefinition
	case equipmentport.EdgeParent:
		return m.clearedparent
	case equipmentport.EdgeLink:
		return m.clearedlink
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *EquipmentPortMutation) ClearEdge(name string) error {
	switch name {
	case equipmentport.EdgeDefinition:
		m.ClearDefinition()
		return nil
	case equipmentport.EdgeParent:
		m.ClearParent()
		return nil
	case equipmentport.EdgeLink:
		m.ClearLink()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPort unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *EquipmentPortMutation) ResetEdge(name string) error {
	switch name {
	case equipmentport.EdgeDefinition:
		m.ResetDefinition()
		return nil
	case equipmentport.EdgeParent:
		m.ResetParent()
		return nil
	case equipmentport.EdgeLink:
		m.ResetLink()
		return nil
	case equipmentport.EdgeProperties:
		m.ResetProperties()
		return nil
	case equipmentport.EdgeEndpoints:
		m.ResetEndpoints()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPort edge %s", name)
}

// EquipmentPortDefinitionMutation represents an operation that mutate the EquipmentPortDefinitions
// nodes in the graph.
type EquipmentPortDefinitionMutation struct {
	config
	op                         Op
	typ                        string
	id                         *int
	create_time                *time.Time
	update_time                *time.Time
	name                       *string
	index                      *int
	addindex                   *int
	bandwidth                  *string
	visibility_label           *string
	clearedFields              map[string]bool
	equipment_port_type        *int
	clearedequipment_port_type bool
	ports                      map[int]struct{}
	removedports               map[int]struct{}
	equipment_type             *int
	clearedequipment_type      bool
}

var _ ent.Mutation = (*EquipmentPortDefinitionMutation)(nil)

// newEquipmentPortDefinitionMutation creates new mutation for $n.Name.
func newEquipmentPortDefinitionMutation(c config, op Op) *EquipmentPortDefinitionMutation {
	return &EquipmentPortDefinitionMutation{
		config:        c,
		op:            op,
		typ:           TypeEquipmentPortDefinition,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EquipmentPortDefinitionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EquipmentPortDefinitionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *EquipmentPortDefinitionMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *EquipmentPortDefinitionMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *EquipmentPortDefinitionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *EquipmentPortDefinitionMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *EquipmentPortDefinitionMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *EquipmentPortDefinitionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *EquipmentPortDefinitionMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *EquipmentPortDefinitionMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *EquipmentPortDefinitionMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *EquipmentPortDefinitionMutation) ResetName() {
	m.name = nil
}

// SetIndex sets the index field.
func (m *EquipmentPortDefinitionMutation) SetIndex(i int) {
	m.index = &i
	m.addindex = nil
}

// Index returns the index value in the mutation.
func (m *EquipmentPortDefinitionMutation) Index() (r int, exists bool) {
	v := m.index
	if v == nil {
		return
	}
	return *v, true
}

// AddIndex adds i to index.
func (m *EquipmentPortDefinitionMutation) AddIndex(i int) {
	if m.addindex != nil {
		*m.addindex += i
	} else {
		m.addindex = &i
	}
}

// AddedIndex returns the value that was added to the index field in this mutation.
func (m *EquipmentPortDefinitionMutation) AddedIndex() (r int, exists bool) {
	v := m.addindex
	if v == nil {
		return
	}
	return *v, true
}

// ClearIndex clears the value of index.
func (m *EquipmentPortDefinitionMutation) ClearIndex() {
	m.index = nil
	m.addindex = nil
	m.clearedFields[equipmentportdefinition.FieldIndex] = true
}

// IndexCleared returns if the field index was cleared in this mutation.
func (m *EquipmentPortDefinitionMutation) IndexCleared() bool {
	return m.clearedFields[equipmentportdefinition.FieldIndex]
}

// ResetIndex reset all changes of the index field.
func (m *EquipmentPortDefinitionMutation) ResetIndex() {
	m.index = nil
	m.addindex = nil
	delete(m.clearedFields, equipmentportdefinition.FieldIndex)
}

// SetBandwidth sets the bandwidth field.
func (m *EquipmentPortDefinitionMutation) SetBandwidth(s string) {
	m.bandwidth = &s
}

// Bandwidth returns the bandwidth value in the mutation.
func (m *EquipmentPortDefinitionMutation) Bandwidth() (r string, exists bool) {
	v := m.bandwidth
	if v == nil {
		return
	}
	return *v, true
}

// ClearBandwidth clears the value of bandwidth.
func (m *EquipmentPortDefinitionMutation) ClearBandwidth() {
	m.bandwidth = nil
	m.clearedFields[equipmentportdefinition.FieldBandwidth] = true
}

// BandwidthCleared returns if the field bandwidth was cleared in this mutation.
func (m *EquipmentPortDefinitionMutation) BandwidthCleared() bool {
	return m.clearedFields[equipmentportdefinition.FieldBandwidth]
}

// ResetBandwidth reset all changes of the bandwidth field.
func (m *EquipmentPortDefinitionMutation) ResetBandwidth() {
	m.bandwidth = nil
	delete(m.clearedFields, equipmentportdefinition.FieldBandwidth)
}

// SetVisibilityLabel sets the visibility_label field.
func (m *EquipmentPortDefinitionMutation) SetVisibilityLabel(s string) {
	m.visibility_label = &s
}

// VisibilityLabel returns the visibility_label value in the mutation.
func (m *EquipmentPortDefinitionMutation) VisibilityLabel() (r string, exists bool) {
	v := m.visibility_label
	if v == nil {
		return
	}
	return *v, true
}

// ClearVisibilityLabel clears the value of visibility_label.
func (m *EquipmentPortDefinitionMutation) ClearVisibilityLabel() {
	m.visibility_label = nil
	m.clearedFields[equipmentportdefinition.FieldVisibilityLabel] = true
}

// VisibilityLabelCleared returns if the field visibility_label was cleared in this mutation.
func (m *EquipmentPortDefinitionMutation) VisibilityLabelCleared() bool {
	return m.clearedFields[equipmentportdefinition.FieldVisibilityLabel]
}

// ResetVisibilityLabel reset all changes of the visibility_label field.
func (m *EquipmentPortDefinitionMutation) ResetVisibilityLabel() {
	m.visibility_label = nil
	delete(m.clearedFields, equipmentportdefinition.FieldVisibilityLabel)
}

// SetEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id.
func (m *EquipmentPortDefinitionMutation) SetEquipmentPortTypeID(id int) {
	m.equipment_port_type = &id
}

// ClearEquipmentPortType clears the equipment_port_type edge to EquipmentPortType.
func (m *EquipmentPortDefinitionMutation) ClearEquipmentPortType() {
	m.clearedequipment_port_type = true
}

// EquipmentPortTypeCleared returns if the edge equipment_port_type was cleared.
func (m *EquipmentPortDefinitionMutation) EquipmentPortTypeCleared() bool {
	return m.clearedequipment_port_type
}

// EquipmentPortTypeID returns the equipment_port_type id in the mutation.
func (m *EquipmentPortDefinitionMutation) EquipmentPortTypeID() (id int, exists bool) {
	if m.equipment_port_type != nil {
		return *m.equipment_port_type, true
	}
	return
}

// EquipmentPortTypeIDs returns the equipment_port_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// EquipmentPortTypeID instead. It exists only for internal usage by the builders.
func (m *EquipmentPortDefinitionMutation) EquipmentPortTypeIDs() (ids []int) {
	if id := m.equipment_port_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEquipmentPortType reset all changes of the equipment_port_type edge.
func (m *EquipmentPortDefinitionMutation) ResetEquipmentPortType() {
	m.equipment_port_type = nil
	m.clearedequipment_port_type = false
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (m *EquipmentPortDefinitionMutation) AddPortIDs(ids ...int) {
	if m.ports == nil {
		m.ports = make(map[int]struct{})
	}
	for i := range ids {
		m.ports[ids[i]] = struct{}{}
	}
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (m *EquipmentPortDefinitionMutation) RemovePortIDs(ids ...int) {
	if m.removedports == nil {
		m.removedports = make(map[int]struct{})
	}
	for i := range ids {
		m.removedports[ids[i]] = struct{}{}
	}
}

// RemovedPorts returns the removed ids of ports.
func (m *EquipmentPortDefinitionMutation) RemovedPortsIDs() (ids []int) {
	for id := range m.removedports {
		ids = append(ids, id)
	}
	return
}

// PortsIDs returns the ports ids in the mutation.
func (m *EquipmentPortDefinitionMutation) PortsIDs() (ids []int) {
	for id := range m.ports {
		ids = append(ids, id)
	}
	return
}

// ResetPorts reset all changes of the ports edge.
func (m *EquipmentPortDefinitionMutation) ResetPorts() {
	m.ports = nil
	m.removedports = nil
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (m *EquipmentPortDefinitionMutation) SetEquipmentTypeID(id int) {
	m.equipment_type = &id
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (m *EquipmentPortDefinitionMutation) ClearEquipmentType() {
	m.clearedequipment_type = true
}

// EquipmentTypeCleared returns if the edge equipment_type was cleared.
func (m *EquipmentPortDefinitionMutation) EquipmentTypeCleared() bool {
	return m.clearedequipment_type
}

// EquipmentTypeID returns the equipment_type id in the mutation.
func (m *EquipmentPortDefinitionMutation) EquipmentTypeID() (id int, exists bool) {
	if m.equipment_type != nil {
		return *m.equipment_type, true
	}
	return
}

// EquipmentTypeIDs returns the equipment_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// EquipmentTypeID instead. It exists only for internal usage by the builders.
func (m *EquipmentPortDefinitionMutation) EquipmentTypeIDs() (ids []int) {
	if id := m.equipment_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEquipmentType reset all changes of the equipment_type edge.
func (m *EquipmentPortDefinitionMutation) ResetEquipmentType() {
	m.equipment_type = nil
	m.clearedequipment_type = false
}

// Op returns the operation name.
func (m *EquipmentPortDefinitionMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (EquipmentPortDefinition).
func (m *EquipmentPortDefinitionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *EquipmentPortDefinitionMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m.create_time != nil {
		fields = append(fields, equipmentportdefinition.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, equipmentportdefinition.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, equipmentportdefinition.FieldName)
	}
	if m.index != nil {
		fields = append(fields, equipmentportdefinition.FieldIndex)
	}
	if m.bandwidth != nil {
		fields = append(fields, equipmentportdefinition.FieldBandwidth)
	}
	if m.visibility_label != nil {
		fields = append(fields, equipmentportdefinition.FieldVisibilityLabel)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *EquipmentPortDefinitionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case equipmentportdefinition.FieldCreateTime:
		return m.CreateTime()
	case equipmentportdefinition.FieldUpdateTime:
		return m.UpdateTime()
	case equipmentportdefinition.FieldName:
		return m.Name()
	case equipmentportdefinition.FieldIndex:
		return m.Index()
	case equipmentportdefinition.FieldBandwidth:
		return m.Bandwidth()
	case equipmentportdefinition.FieldVisibilityLabel:
		return m.VisibilityLabel()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentPortDefinitionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case equipmentportdefinition.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case equipmentportdefinition.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case equipmentportdefinition.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case equipmentportdefinition.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIndex(v)
		return nil
	case equipmentportdefinition.FieldBandwidth:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBandwidth(v)
		return nil
	case equipmentportdefinition.FieldVisibilityLabel:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetVisibilityLabel(v)
		return nil
	}
	return fmt.Errorf("unknown EquipmentPortDefinition field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *EquipmentPortDefinitionMutation) AddedFields() []string {
	var fields []string
	if m.addindex != nil {
		fields = append(fields, equipmentportdefinition.FieldIndex)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *EquipmentPortDefinitionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case equipmentportdefinition.FieldIndex:
		return m.AddedIndex()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentPortDefinitionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case equipmentportdefinition.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIndex(v)
		return nil
	}
	return fmt.Errorf("unknown EquipmentPortDefinition numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *EquipmentPortDefinitionMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[equipmentportdefinition.FieldIndex] {
		fields = append(fields, equipmentportdefinition.FieldIndex)
	}
	if m.clearedFields[equipmentportdefinition.FieldBandwidth] {
		fields = append(fields, equipmentportdefinition.FieldBandwidth)
	}
	if m.clearedFields[equipmentportdefinition.FieldVisibilityLabel] {
		fields = append(fields, equipmentportdefinition.FieldVisibilityLabel)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *EquipmentPortDefinitionMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *EquipmentPortDefinitionMutation) ClearField(name string) error {
	switch name {
	case equipmentportdefinition.FieldIndex:
		m.ClearIndex()
		return nil
	case equipmentportdefinition.FieldBandwidth:
		m.ClearBandwidth()
		return nil
	case equipmentportdefinition.FieldVisibilityLabel:
		m.ClearVisibilityLabel()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPortDefinition nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *EquipmentPortDefinitionMutation) ResetField(name string) error {
	switch name {
	case equipmentportdefinition.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case equipmentportdefinition.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case equipmentportdefinition.FieldName:
		m.ResetName()
		return nil
	case equipmentportdefinition.FieldIndex:
		m.ResetIndex()
		return nil
	case equipmentportdefinition.FieldBandwidth:
		m.ResetBandwidth()
		return nil
	case equipmentportdefinition.FieldVisibilityLabel:
		m.ResetVisibilityLabel()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPortDefinition field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *EquipmentPortDefinitionMutation) AddedEdges() []string {
	edges := make([]string, 0, 3)
	if m.equipment_port_type != nil {
		edges = append(edges, equipmentportdefinition.EdgeEquipmentPortType)
	}
	if m.ports != nil {
		edges = append(edges, equipmentportdefinition.EdgePorts)
	}
	if m.equipment_type != nil {
		edges = append(edges, equipmentportdefinition.EdgeEquipmentType)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *EquipmentPortDefinitionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case equipmentportdefinition.EdgeEquipmentPortType:
		if id := m.equipment_port_type; id != nil {
			return []ent.Value{*id}
		}
	case equipmentportdefinition.EdgePorts:
		ids := make([]ent.Value, 0, len(m.ports))
		for id := range m.ports {
			ids = append(ids, id)
		}
		return ids
	case equipmentportdefinition.EdgeEquipmentType:
		if id := m.equipment_type; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *EquipmentPortDefinitionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 3)
	if m.removedports != nil {
		edges = append(edges, equipmentportdefinition.EdgePorts)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *EquipmentPortDefinitionMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case equipmentportdefinition.EdgePorts:
		ids := make([]ent.Value, 0, len(m.removedports))
		for id := range m.removedports {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *EquipmentPortDefinitionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 3)
	if m.clearedequipment_port_type {
		edges = append(edges, equipmentportdefinition.EdgeEquipmentPortType)
	}
	if m.clearedequipment_type {
		edges = append(edges, equipmentportdefinition.EdgeEquipmentType)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *EquipmentPortDefinitionMutation) EdgeCleared(name string) bool {
	switch name {
	case equipmentportdefinition.EdgeEquipmentPortType:
		return m.clearedequipment_port_type
	case equipmentportdefinition.EdgeEquipmentType:
		return m.clearedequipment_type
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *EquipmentPortDefinitionMutation) ClearEdge(name string) error {
	switch name {
	case equipmentportdefinition.EdgeEquipmentPortType:
		m.ClearEquipmentPortType()
		return nil
	case equipmentportdefinition.EdgeEquipmentType:
		m.ClearEquipmentType()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPortDefinition unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *EquipmentPortDefinitionMutation) ResetEdge(name string) error {
	switch name {
	case equipmentportdefinition.EdgeEquipmentPortType:
		m.ResetEquipmentPortType()
		return nil
	case equipmentportdefinition.EdgePorts:
		m.ResetPorts()
		return nil
	case equipmentportdefinition.EdgeEquipmentType:
		m.ResetEquipmentType()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPortDefinition edge %s", name)
}

// EquipmentPortTypeMutation represents an operation that mutate the EquipmentPortTypes
// nodes in the graph.
type EquipmentPortTypeMutation struct {
	config
	op                         Op
	typ                        string
	id                         *int
	create_time                *time.Time
	update_time                *time.Time
	name                       *string
	clearedFields              map[string]bool
	property_types             map[int]struct{}
	removedproperty_types      map[int]struct{}
	link_property_types        map[int]struct{}
	removedlink_property_types map[int]struct{}
	port_definitions           map[int]struct{}
	removedport_definitions    map[int]struct{}
}

var _ ent.Mutation = (*EquipmentPortTypeMutation)(nil)

// newEquipmentPortTypeMutation creates new mutation for $n.Name.
func newEquipmentPortTypeMutation(c config, op Op) *EquipmentPortTypeMutation {
	return &EquipmentPortTypeMutation{
		config:        c,
		op:            op,
		typ:           TypeEquipmentPortType,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EquipmentPortTypeMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EquipmentPortTypeMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *EquipmentPortTypeMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *EquipmentPortTypeMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *EquipmentPortTypeMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *EquipmentPortTypeMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *EquipmentPortTypeMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *EquipmentPortTypeMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *EquipmentPortTypeMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *EquipmentPortTypeMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *EquipmentPortTypeMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *EquipmentPortTypeMutation) ResetName() {
	m.name = nil
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (m *EquipmentPortTypeMutation) AddPropertyTypeIDs(ids ...int) {
	if m.property_types == nil {
		m.property_types = make(map[int]struct{})
	}
	for i := range ids {
		m.property_types[ids[i]] = struct{}{}
	}
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (m *EquipmentPortTypeMutation) RemovePropertyTypeIDs(ids ...int) {
	if m.removedproperty_types == nil {
		m.removedproperty_types = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperty_types[ids[i]] = struct{}{}
	}
}

// RemovedPropertyTypes returns the removed ids of property_types.
func (m *EquipmentPortTypeMutation) RemovedPropertyTypesIDs() (ids []int) {
	for id := range m.removedproperty_types {
		ids = append(ids, id)
	}
	return
}

// PropertyTypesIDs returns the property_types ids in the mutation.
func (m *EquipmentPortTypeMutation) PropertyTypesIDs() (ids []int) {
	for id := range m.property_types {
		ids = append(ids, id)
	}
	return
}

// ResetPropertyTypes reset all changes of the property_types edge.
func (m *EquipmentPortTypeMutation) ResetPropertyTypes() {
	m.property_types = nil
	m.removedproperty_types = nil
}

// AddLinkPropertyTypeIDs adds the link_property_types edge to PropertyType by ids.
func (m *EquipmentPortTypeMutation) AddLinkPropertyTypeIDs(ids ...int) {
	if m.link_property_types == nil {
		m.link_property_types = make(map[int]struct{})
	}
	for i := range ids {
		m.link_property_types[ids[i]] = struct{}{}
	}
}

// RemoveLinkPropertyTypeIDs removes the link_property_types edge to PropertyType by ids.
func (m *EquipmentPortTypeMutation) RemoveLinkPropertyTypeIDs(ids ...int) {
	if m.removedlink_property_types == nil {
		m.removedlink_property_types = make(map[int]struct{})
	}
	for i := range ids {
		m.removedlink_property_types[ids[i]] = struct{}{}
	}
}

// RemovedLinkPropertyTypes returns the removed ids of link_property_types.
func (m *EquipmentPortTypeMutation) RemovedLinkPropertyTypesIDs() (ids []int) {
	for id := range m.removedlink_property_types {
		ids = append(ids, id)
	}
	return
}

// LinkPropertyTypesIDs returns the link_property_types ids in the mutation.
func (m *EquipmentPortTypeMutation) LinkPropertyTypesIDs() (ids []int) {
	for id := range m.link_property_types {
		ids = append(ids, id)
	}
	return
}

// ResetLinkPropertyTypes reset all changes of the link_property_types edge.
func (m *EquipmentPortTypeMutation) ResetLinkPropertyTypes() {
	m.link_property_types = nil
	m.removedlink_property_types = nil
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (m *EquipmentPortTypeMutation) AddPortDefinitionIDs(ids ...int) {
	if m.port_definitions == nil {
		m.port_definitions = make(map[int]struct{})
	}
	for i := range ids {
		m.port_definitions[ids[i]] = struct{}{}
	}
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (m *EquipmentPortTypeMutation) RemovePortDefinitionIDs(ids ...int) {
	if m.removedport_definitions == nil {
		m.removedport_definitions = make(map[int]struct{})
	}
	for i := range ids {
		m.removedport_definitions[ids[i]] = struct{}{}
	}
}

// RemovedPortDefinitions returns the removed ids of port_definitions.
func (m *EquipmentPortTypeMutation) RemovedPortDefinitionsIDs() (ids []int) {
	for id := range m.removedport_definitions {
		ids = append(ids, id)
	}
	return
}

// PortDefinitionsIDs returns the port_definitions ids in the mutation.
func (m *EquipmentPortTypeMutation) PortDefinitionsIDs() (ids []int) {
	for id := range m.port_definitions {
		ids = append(ids, id)
	}
	return
}

// ResetPortDefinitions reset all changes of the port_definitions edge.
func (m *EquipmentPortTypeMutation) ResetPortDefinitions() {
	m.port_definitions = nil
	m.removedport_definitions = nil
}

// Op returns the operation name.
func (m *EquipmentPortTypeMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (EquipmentPortType).
func (m *EquipmentPortTypeMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *EquipmentPortTypeMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.create_time != nil {
		fields = append(fields, equipmentporttype.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, equipmentporttype.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, equipmentporttype.FieldName)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *EquipmentPortTypeMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case equipmentporttype.FieldCreateTime:
		return m.CreateTime()
	case equipmentporttype.FieldUpdateTime:
		return m.UpdateTime()
	case equipmentporttype.FieldName:
		return m.Name()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentPortTypeMutation) SetField(name string, value ent.Value) error {
	switch name {
	case equipmentporttype.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case equipmentporttype.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case equipmentporttype.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	}
	return fmt.Errorf("unknown EquipmentPortType field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *EquipmentPortTypeMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *EquipmentPortTypeMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentPortTypeMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown EquipmentPortType numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *EquipmentPortTypeMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *EquipmentPortTypeMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *EquipmentPortTypeMutation) ClearField(name string) error {
	return fmt.Errorf("unknown EquipmentPortType nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *EquipmentPortTypeMutation) ResetField(name string) error {
	switch name {
	case equipmentporttype.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case equipmentporttype.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case equipmentporttype.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPortType field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *EquipmentPortTypeMutation) AddedEdges() []string {
	edges := make([]string, 0, 3)
	if m.property_types != nil {
		edges = append(edges, equipmentporttype.EdgePropertyTypes)
	}
	if m.link_property_types != nil {
		edges = append(edges, equipmentporttype.EdgeLinkPropertyTypes)
	}
	if m.port_definitions != nil {
		edges = append(edges, equipmentporttype.EdgePortDefinitions)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *EquipmentPortTypeMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case equipmentporttype.EdgePropertyTypes:
		ids := make([]ent.Value, 0, len(m.property_types))
		for id := range m.property_types {
			ids = append(ids, id)
		}
		return ids
	case equipmentporttype.EdgeLinkPropertyTypes:
		ids := make([]ent.Value, 0, len(m.link_property_types))
		for id := range m.link_property_types {
			ids = append(ids, id)
		}
		return ids
	case equipmentporttype.EdgePortDefinitions:
		ids := make([]ent.Value, 0, len(m.port_definitions))
		for id := range m.port_definitions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *EquipmentPortTypeMutation) RemovedEdges() []string {
	edges := make([]string, 0, 3)
	if m.removedproperty_types != nil {
		edges = append(edges, equipmentporttype.EdgePropertyTypes)
	}
	if m.removedlink_property_types != nil {
		edges = append(edges, equipmentporttype.EdgeLinkPropertyTypes)
	}
	if m.removedport_definitions != nil {
		edges = append(edges, equipmentporttype.EdgePortDefinitions)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *EquipmentPortTypeMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case equipmentporttype.EdgePropertyTypes:
		ids := make([]ent.Value, 0, len(m.removedproperty_types))
		for id := range m.removedproperty_types {
			ids = append(ids, id)
		}
		return ids
	case equipmentporttype.EdgeLinkPropertyTypes:
		ids := make([]ent.Value, 0, len(m.removedlink_property_types))
		for id := range m.removedlink_property_types {
			ids = append(ids, id)
		}
		return ids
	case equipmentporttype.EdgePortDefinitions:
		ids := make([]ent.Value, 0, len(m.removedport_definitions))
		for id := range m.removedport_definitions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *EquipmentPortTypeMutation) ClearedEdges() []string {
	edges := make([]string, 0, 3)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *EquipmentPortTypeMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *EquipmentPortTypeMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown EquipmentPortType unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *EquipmentPortTypeMutation) ResetEdge(name string) error {
	switch name {
	case equipmentporttype.EdgePropertyTypes:
		m.ResetPropertyTypes()
		return nil
	case equipmentporttype.EdgeLinkPropertyTypes:
		m.ResetLinkPropertyTypes()
		return nil
	case equipmentporttype.EdgePortDefinitions:
		m.ResetPortDefinitions()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPortType edge %s", name)
}

// EquipmentPositionMutation represents an operation that mutate the EquipmentPositions
// nodes in the graph.
type EquipmentPositionMutation struct {
	config
	op                Op
	typ               string
	id                *int
	create_time       *time.Time
	update_time       *time.Time
	clearedFields     map[string]bool
	definition        *int
	cleareddefinition bool
	parent            *int
	clearedparent     bool
	attachment        *int
	clearedattachment bool
}

var _ ent.Mutation = (*EquipmentPositionMutation)(nil)

// newEquipmentPositionMutation creates new mutation for $n.Name.
func newEquipmentPositionMutation(c config, op Op) *EquipmentPositionMutation {
	return &EquipmentPositionMutation{
		config:        c,
		op:            op,
		typ:           TypeEquipmentPosition,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EquipmentPositionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EquipmentPositionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *EquipmentPositionMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *EquipmentPositionMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *EquipmentPositionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *EquipmentPositionMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *EquipmentPositionMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *EquipmentPositionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *EquipmentPositionMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetDefinitionID sets the definition edge to EquipmentPositionDefinition by id.
func (m *EquipmentPositionMutation) SetDefinitionID(id int) {
	m.definition = &id
}

// ClearDefinition clears the definition edge to EquipmentPositionDefinition.
func (m *EquipmentPositionMutation) ClearDefinition() {
	m.cleareddefinition = true
}

// DefinitionCleared returns if the edge definition was cleared.
func (m *EquipmentPositionMutation) DefinitionCleared() bool {
	return m.cleareddefinition
}

// DefinitionID returns the definition id in the mutation.
func (m *EquipmentPositionMutation) DefinitionID() (id int, exists bool) {
	if m.definition != nil {
		return *m.definition, true
	}
	return
}

// DefinitionIDs returns the definition ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// DefinitionID instead. It exists only for internal usage by the builders.
func (m *EquipmentPositionMutation) DefinitionIDs() (ids []int) {
	if id := m.definition; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetDefinition reset all changes of the definition edge.
func (m *EquipmentPositionMutation) ResetDefinition() {
	m.definition = nil
	m.cleareddefinition = false
}

// SetParentID sets the parent edge to Equipment by id.
func (m *EquipmentPositionMutation) SetParentID(id int) {
	m.parent = &id
}

// ClearParent clears the parent edge to Equipment.
func (m *EquipmentPositionMutation) ClearParent() {
	m.clearedparent = true
}

// ParentCleared returns if the edge parent was cleared.
func (m *EquipmentPositionMutation) ParentCleared() bool {
	return m.clearedparent
}

// ParentID returns the parent id in the mutation.
func (m *EquipmentPositionMutation) ParentID() (id int, exists bool) {
	if m.parent != nil {
		return *m.parent, true
	}
	return
}

// ParentIDs returns the parent ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ParentID instead. It exists only for internal usage by the builders.
func (m *EquipmentPositionMutation) ParentIDs() (ids []int) {
	if id := m.parent; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetParent reset all changes of the parent edge.
func (m *EquipmentPositionMutation) ResetParent() {
	m.parent = nil
	m.clearedparent = false
}

// SetAttachmentID sets the attachment edge to Equipment by id.
func (m *EquipmentPositionMutation) SetAttachmentID(id int) {
	m.attachment = &id
}

// ClearAttachment clears the attachment edge to Equipment.
func (m *EquipmentPositionMutation) ClearAttachment() {
	m.clearedattachment = true
}

// AttachmentCleared returns if the edge attachment was cleared.
func (m *EquipmentPositionMutation) AttachmentCleared() bool {
	return m.clearedattachment
}

// AttachmentID returns the attachment id in the mutation.
func (m *EquipmentPositionMutation) AttachmentID() (id int, exists bool) {
	if m.attachment != nil {
		return *m.attachment, true
	}
	return
}

// AttachmentIDs returns the attachment ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// AttachmentID instead. It exists only for internal usage by the builders.
func (m *EquipmentPositionMutation) AttachmentIDs() (ids []int) {
	if id := m.attachment; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetAttachment reset all changes of the attachment edge.
func (m *EquipmentPositionMutation) ResetAttachment() {
	m.attachment = nil
	m.clearedattachment = false
}

// Op returns the operation name.
func (m *EquipmentPositionMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (EquipmentPosition).
func (m *EquipmentPositionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *EquipmentPositionMutation) Fields() []string {
	fields := make([]string, 0, 2)
	if m.create_time != nil {
		fields = append(fields, equipmentposition.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, equipmentposition.FieldUpdateTime)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *EquipmentPositionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case equipmentposition.FieldCreateTime:
		return m.CreateTime()
	case equipmentposition.FieldUpdateTime:
		return m.UpdateTime()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentPositionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case equipmentposition.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case equipmentposition.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	}
	return fmt.Errorf("unknown EquipmentPosition field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *EquipmentPositionMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *EquipmentPositionMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentPositionMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown EquipmentPosition numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *EquipmentPositionMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *EquipmentPositionMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *EquipmentPositionMutation) ClearField(name string) error {
	return fmt.Errorf("unknown EquipmentPosition nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *EquipmentPositionMutation) ResetField(name string) error {
	switch name {
	case equipmentposition.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case equipmentposition.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPosition field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *EquipmentPositionMutation) AddedEdges() []string {
	edges := make([]string, 0, 3)
	if m.definition != nil {
		edges = append(edges, equipmentposition.EdgeDefinition)
	}
	if m.parent != nil {
		edges = append(edges, equipmentposition.EdgeParent)
	}
	if m.attachment != nil {
		edges = append(edges, equipmentposition.EdgeAttachment)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *EquipmentPositionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case equipmentposition.EdgeDefinition:
		if id := m.definition; id != nil {
			return []ent.Value{*id}
		}
	case equipmentposition.EdgeParent:
		if id := m.parent; id != nil {
			return []ent.Value{*id}
		}
	case equipmentposition.EdgeAttachment:
		if id := m.attachment; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *EquipmentPositionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 3)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *EquipmentPositionMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *EquipmentPositionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 3)
	if m.cleareddefinition {
		edges = append(edges, equipmentposition.EdgeDefinition)
	}
	if m.clearedparent {
		edges = append(edges, equipmentposition.EdgeParent)
	}
	if m.clearedattachment {
		edges = append(edges, equipmentposition.EdgeAttachment)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *EquipmentPositionMutation) EdgeCleared(name string) bool {
	switch name {
	case equipmentposition.EdgeDefinition:
		return m.cleareddefinition
	case equipmentposition.EdgeParent:
		return m.clearedparent
	case equipmentposition.EdgeAttachment:
		return m.clearedattachment
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *EquipmentPositionMutation) ClearEdge(name string) error {
	switch name {
	case equipmentposition.EdgeDefinition:
		m.ClearDefinition()
		return nil
	case equipmentposition.EdgeParent:
		m.ClearParent()
		return nil
	case equipmentposition.EdgeAttachment:
		m.ClearAttachment()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPosition unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *EquipmentPositionMutation) ResetEdge(name string) error {
	switch name {
	case equipmentposition.EdgeDefinition:
		m.ResetDefinition()
		return nil
	case equipmentposition.EdgeParent:
		m.ResetParent()
		return nil
	case equipmentposition.EdgeAttachment:
		m.ResetAttachment()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPosition edge %s", name)
}

// EquipmentPositionDefinitionMutation represents an operation that mutate the EquipmentPositionDefinitions
// nodes in the graph.
type EquipmentPositionDefinitionMutation struct {
	config
	op                    Op
	typ                   string
	id                    *int
	create_time           *time.Time
	update_time           *time.Time
	name                  *string
	index                 *int
	addindex              *int
	visibility_label      *string
	clearedFields         map[string]bool
	positions             map[int]struct{}
	removedpositions      map[int]struct{}
	equipment_type        *int
	clearedequipment_type bool
}

var _ ent.Mutation = (*EquipmentPositionDefinitionMutation)(nil)

// newEquipmentPositionDefinitionMutation creates new mutation for $n.Name.
func newEquipmentPositionDefinitionMutation(c config, op Op) *EquipmentPositionDefinitionMutation {
	return &EquipmentPositionDefinitionMutation{
		config:        c,
		op:            op,
		typ:           TypeEquipmentPositionDefinition,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EquipmentPositionDefinitionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EquipmentPositionDefinitionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *EquipmentPositionDefinitionMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *EquipmentPositionDefinitionMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *EquipmentPositionDefinitionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *EquipmentPositionDefinitionMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *EquipmentPositionDefinitionMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *EquipmentPositionDefinitionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *EquipmentPositionDefinitionMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *EquipmentPositionDefinitionMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *EquipmentPositionDefinitionMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *EquipmentPositionDefinitionMutation) ResetName() {
	m.name = nil
}

// SetIndex sets the index field.
func (m *EquipmentPositionDefinitionMutation) SetIndex(i int) {
	m.index = &i
	m.addindex = nil
}

// Index returns the index value in the mutation.
func (m *EquipmentPositionDefinitionMutation) Index() (r int, exists bool) {
	v := m.index
	if v == nil {
		return
	}
	return *v, true
}

// AddIndex adds i to index.
func (m *EquipmentPositionDefinitionMutation) AddIndex(i int) {
	if m.addindex != nil {
		*m.addindex += i
	} else {
		m.addindex = &i
	}
}

// AddedIndex returns the value that was added to the index field in this mutation.
func (m *EquipmentPositionDefinitionMutation) AddedIndex() (r int, exists bool) {
	v := m.addindex
	if v == nil {
		return
	}
	return *v, true
}

// ClearIndex clears the value of index.
func (m *EquipmentPositionDefinitionMutation) ClearIndex() {
	m.index = nil
	m.addindex = nil
	m.clearedFields[equipmentpositiondefinition.FieldIndex] = true
}

// IndexCleared returns if the field index was cleared in this mutation.
func (m *EquipmentPositionDefinitionMutation) IndexCleared() bool {
	return m.clearedFields[equipmentpositiondefinition.FieldIndex]
}

// ResetIndex reset all changes of the index field.
func (m *EquipmentPositionDefinitionMutation) ResetIndex() {
	m.index = nil
	m.addindex = nil
	delete(m.clearedFields, equipmentpositiondefinition.FieldIndex)
}

// SetVisibilityLabel sets the visibility_label field.
func (m *EquipmentPositionDefinitionMutation) SetVisibilityLabel(s string) {
	m.visibility_label = &s
}

// VisibilityLabel returns the visibility_label value in the mutation.
func (m *EquipmentPositionDefinitionMutation) VisibilityLabel() (r string, exists bool) {
	v := m.visibility_label
	if v == nil {
		return
	}
	return *v, true
}

// ClearVisibilityLabel clears the value of visibility_label.
func (m *EquipmentPositionDefinitionMutation) ClearVisibilityLabel() {
	m.visibility_label = nil
	m.clearedFields[equipmentpositiondefinition.FieldVisibilityLabel] = true
}

// VisibilityLabelCleared returns if the field visibility_label was cleared in this mutation.
func (m *EquipmentPositionDefinitionMutation) VisibilityLabelCleared() bool {
	return m.clearedFields[equipmentpositiondefinition.FieldVisibilityLabel]
}

// ResetVisibilityLabel reset all changes of the visibility_label field.
func (m *EquipmentPositionDefinitionMutation) ResetVisibilityLabel() {
	m.visibility_label = nil
	delete(m.clearedFields, equipmentpositiondefinition.FieldVisibilityLabel)
}

// AddPositionIDs adds the positions edge to EquipmentPosition by ids.
func (m *EquipmentPositionDefinitionMutation) AddPositionIDs(ids ...int) {
	if m.positions == nil {
		m.positions = make(map[int]struct{})
	}
	for i := range ids {
		m.positions[ids[i]] = struct{}{}
	}
}

// RemovePositionIDs removes the positions edge to EquipmentPosition by ids.
func (m *EquipmentPositionDefinitionMutation) RemovePositionIDs(ids ...int) {
	if m.removedpositions == nil {
		m.removedpositions = make(map[int]struct{})
	}
	for i := range ids {
		m.removedpositions[ids[i]] = struct{}{}
	}
}

// RemovedPositions returns the removed ids of positions.
func (m *EquipmentPositionDefinitionMutation) RemovedPositionsIDs() (ids []int) {
	for id := range m.removedpositions {
		ids = append(ids, id)
	}
	return
}

// PositionsIDs returns the positions ids in the mutation.
func (m *EquipmentPositionDefinitionMutation) PositionsIDs() (ids []int) {
	for id := range m.positions {
		ids = append(ids, id)
	}
	return
}

// ResetPositions reset all changes of the positions edge.
func (m *EquipmentPositionDefinitionMutation) ResetPositions() {
	m.positions = nil
	m.removedpositions = nil
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (m *EquipmentPositionDefinitionMutation) SetEquipmentTypeID(id int) {
	m.equipment_type = &id
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (m *EquipmentPositionDefinitionMutation) ClearEquipmentType() {
	m.clearedequipment_type = true
}

// EquipmentTypeCleared returns if the edge equipment_type was cleared.
func (m *EquipmentPositionDefinitionMutation) EquipmentTypeCleared() bool {
	return m.clearedequipment_type
}

// EquipmentTypeID returns the equipment_type id in the mutation.
func (m *EquipmentPositionDefinitionMutation) EquipmentTypeID() (id int, exists bool) {
	if m.equipment_type != nil {
		return *m.equipment_type, true
	}
	return
}

// EquipmentTypeIDs returns the equipment_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// EquipmentTypeID instead. It exists only for internal usage by the builders.
func (m *EquipmentPositionDefinitionMutation) EquipmentTypeIDs() (ids []int) {
	if id := m.equipment_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEquipmentType reset all changes of the equipment_type edge.
func (m *EquipmentPositionDefinitionMutation) ResetEquipmentType() {
	m.equipment_type = nil
	m.clearedequipment_type = false
}

// Op returns the operation name.
func (m *EquipmentPositionDefinitionMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (EquipmentPositionDefinition).
func (m *EquipmentPositionDefinitionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *EquipmentPositionDefinitionMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.create_time != nil {
		fields = append(fields, equipmentpositiondefinition.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, equipmentpositiondefinition.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, equipmentpositiondefinition.FieldName)
	}
	if m.index != nil {
		fields = append(fields, equipmentpositiondefinition.FieldIndex)
	}
	if m.visibility_label != nil {
		fields = append(fields, equipmentpositiondefinition.FieldVisibilityLabel)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *EquipmentPositionDefinitionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case equipmentpositiondefinition.FieldCreateTime:
		return m.CreateTime()
	case equipmentpositiondefinition.FieldUpdateTime:
		return m.UpdateTime()
	case equipmentpositiondefinition.FieldName:
		return m.Name()
	case equipmentpositiondefinition.FieldIndex:
		return m.Index()
	case equipmentpositiondefinition.FieldVisibilityLabel:
		return m.VisibilityLabel()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentPositionDefinitionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case equipmentpositiondefinition.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case equipmentpositiondefinition.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case equipmentpositiondefinition.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case equipmentpositiondefinition.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIndex(v)
		return nil
	case equipmentpositiondefinition.FieldVisibilityLabel:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetVisibilityLabel(v)
		return nil
	}
	return fmt.Errorf("unknown EquipmentPositionDefinition field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *EquipmentPositionDefinitionMutation) AddedFields() []string {
	var fields []string
	if m.addindex != nil {
		fields = append(fields, equipmentpositiondefinition.FieldIndex)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *EquipmentPositionDefinitionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case equipmentpositiondefinition.FieldIndex:
		return m.AddedIndex()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentPositionDefinitionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case equipmentpositiondefinition.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIndex(v)
		return nil
	}
	return fmt.Errorf("unknown EquipmentPositionDefinition numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *EquipmentPositionDefinitionMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[equipmentpositiondefinition.FieldIndex] {
		fields = append(fields, equipmentpositiondefinition.FieldIndex)
	}
	if m.clearedFields[equipmentpositiondefinition.FieldVisibilityLabel] {
		fields = append(fields, equipmentpositiondefinition.FieldVisibilityLabel)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *EquipmentPositionDefinitionMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *EquipmentPositionDefinitionMutation) ClearField(name string) error {
	switch name {
	case equipmentpositiondefinition.FieldIndex:
		m.ClearIndex()
		return nil
	case equipmentpositiondefinition.FieldVisibilityLabel:
		m.ClearVisibilityLabel()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPositionDefinition nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *EquipmentPositionDefinitionMutation) ResetField(name string) error {
	switch name {
	case equipmentpositiondefinition.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case equipmentpositiondefinition.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case equipmentpositiondefinition.FieldName:
		m.ResetName()
		return nil
	case equipmentpositiondefinition.FieldIndex:
		m.ResetIndex()
		return nil
	case equipmentpositiondefinition.FieldVisibilityLabel:
		m.ResetVisibilityLabel()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPositionDefinition field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *EquipmentPositionDefinitionMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.positions != nil {
		edges = append(edges, equipmentpositiondefinition.EdgePositions)
	}
	if m.equipment_type != nil {
		edges = append(edges, equipmentpositiondefinition.EdgeEquipmentType)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *EquipmentPositionDefinitionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case equipmentpositiondefinition.EdgePositions:
		ids := make([]ent.Value, 0, len(m.positions))
		for id := range m.positions {
			ids = append(ids, id)
		}
		return ids
	case equipmentpositiondefinition.EdgeEquipmentType:
		if id := m.equipment_type; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *EquipmentPositionDefinitionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	if m.removedpositions != nil {
		edges = append(edges, equipmentpositiondefinition.EdgePositions)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *EquipmentPositionDefinitionMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case equipmentpositiondefinition.EdgePositions:
		ids := make([]ent.Value, 0, len(m.removedpositions))
		for id := range m.removedpositions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *EquipmentPositionDefinitionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedequipment_type {
		edges = append(edges, equipmentpositiondefinition.EdgeEquipmentType)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *EquipmentPositionDefinitionMutation) EdgeCleared(name string) bool {
	switch name {
	case equipmentpositiondefinition.EdgeEquipmentType:
		return m.clearedequipment_type
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *EquipmentPositionDefinitionMutation) ClearEdge(name string) error {
	switch name {
	case equipmentpositiondefinition.EdgeEquipmentType:
		m.ClearEquipmentType()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPositionDefinition unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *EquipmentPositionDefinitionMutation) ResetEdge(name string) error {
	switch name {
	case equipmentpositiondefinition.EdgePositions:
		m.ResetPositions()
		return nil
	case equipmentpositiondefinition.EdgeEquipmentType:
		m.ResetEquipmentType()
		return nil
	}
	return fmt.Errorf("unknown EquipmentPositionDefinition edge %s", name)
}

// EquipmentTypeMutation represents an operation that mutate the EquipmentTypes
// nodes in the graph.
type EquipmentTypeMutation struct {
	config
	op                          Op
	typ                         string
	id                          *int
	create_time                 *time.Time
	update_time                 *time.Time
	name                        *string
	clearedFields               map[string]bool
	port_definitions            map[int]struct{}
	removedport_definitions     map[int]struct{}
	position_definitions        map[int]struct{}
	removedposition_definitions map[int]struct{}
	property_types              map[int]struct{}
	removedproperty_types       map[int]struct{}
	equipment                   map[int]struct{}
	removedequipment            map[int]struct{}
	category                    *int
	clearedcategory             bool
}

var _ ent.Mutation = (*EquipmentTypeMutation)(nil)

// newEquipmentTypeMutation creates new mutation for $n.Name.
func newEquipmentTypeMutation(c config, op Op) *EquipmentTypeMutation {
	return &EquipmentTypeMutation{
		config:        c,
		op:            op,
		typ:           TypeEquipmentType,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m EquipmentTypeMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m EquipmentTypeMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *EquipmentTypeMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *EquipmentTypeMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *EquipmentTypeMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *EquipmentTypeMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *EquipmentTypeMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *EquipmentTypeMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *EquipmentTypeMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *EquipmentTypeMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *EquipmentTypeMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *EquipmentTypeMutation) ResetName() {
	m.name = nil
}

// AddPortDefinitionIDs adds the port_definitions edge to EquipmentPortDefinition by ids.
func (m *EquipmentTypeMutation) AddPortDefinitionIDs(ids ...int) {
	if m.port_definitions == nil {
		m.port_definitions = make(map[int]struct{})
	}
	for i := range ids {
		m.port_definitions[ids[i]] = struct{}{}
	}
}

// RemovePortDefinitionIDs removes the port_definitions edge to EquipmentPortDefinition by ids.
func (m *EquipmentTypeMutation) RemovePortDefinitionIDs(ids ...int) {
	if m.removedport_definitions == nil {
		m.removedport_definitions = make(map[int]struct{})
	}
	for i := range ids {
		m.removedport_definitions[ids[i]] = struct{}{}
	}
}

// RemovedPortDefinitions returns the removed ids of port_definitions.
func (m *EquipmentTypeMutation) RemovedPortDefinitionsIDs() (ids []int) {
	for id := range m.removedport_definitions {
		ids = append(ids, id)
	}
	return
}

// PortDefinitionsIDs returns the port_definitions ids in the mutation.
func (m *EquipmentTypeMutation) PortDefinitionsIDs() (ids []int) {
	for id := range m.port_definitions {
		ids = append(ids, id)
	}
	return
}

// ResetPortDefinitions reset all changes of the port_definitions edge.
func (m *EquipmentTypeMutation) ResetPortDefinitions() {
	m.port_definitions = nil
	m.removedport_definitions = nil
}

// AddPositionDefinitionIDs adds the position_definitions edge to EquipmentPositionDefinition by ids.
func (m *EquipmentTypeMutation) AddPositionDefinitionIDs(ids ...int) {
	if m.position_definitions == nil {
		m.position_definitions = make(map[int]struct{})
	}
	for i := range ids {
		m.position_definitions[ids[i]] = struct{}{}
	}
}

// RemovePositionDefinitionIDs removes the position_definitions edge to EquipmentPositionDefinition by ids.
func (m *EquipmentTypeMutation) RemovePositionDefinitionIDs(ids ...int) {
	if m.removedposition_definitions == nil {
		m.removedposition_definitions = make(map[int]struct{})
	}
	for i := range ids {
		m.removedposition_definitions[ids[i]] = struct{}{}
	}
}

// RemovedPositionDefinitions returns the removed ids of position_definitions.
func (m *EquipmentTypeMutation) RemovedPositionDefinitionsIDs() (ids []int) {
	for id := range m.removedposition_definitions {
		ids = append(ids, id)
	}
	return
}

// PositionDefinitionsIDs returns the position_definitions ids in the mutation.
func (m *EquipmentTypeMutation) PositionDefinitionsIDs() (ids []int) {
	for id := range m.position_definitions {
		ids = append(ids, id)
	}
	return
}

// ResetPositionDefinitions reset all changes of the position_definitions edge.
func (m *EquipmentTypeMutation) ResetPositionDefinitions() {
	m.position_definitions = nil
	m.removedposition_definitions = nil
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (m *EquipmentTypeMutation) AddPropertyTypeIDs(ids ...int) {
	if m.property_types == nil {
		m.property_types = make(map[int]struct{})
	}
	for i := range ids {
		m.property_types[ids[i]] = struct{}{}
	}
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (m *EquipmentTypeMutation) RemovePropertyTypeIDs(ids ...int) {
	if m.removedproperty_types == nil {
		m.removedproperty_types = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperty_types[ids[i]] = struct{}{}
	}
}

// RemovedPropertyTypes returns the removed ids of property_types.
func (m *EquipmentTypeMutation) RemovedPropertyTypesIDs() (ids []int) {
	for id := range m.removedproperty_types {
		ids = append(ids, id)
	}
	return
}

// PropertyTypesIDs returns the property_types ids in the mutation.
func (m *EquipmentTypeMutation) PropertyTypesIDs() (ids []int) {
	for id := range m.property_types {
		ids = append(ids, id)
	}
	return
}

// ResetPropertyTypes reset all changes of the property_types edge.
func (m *EquipmentTypeMutation) ResetPropertyTypes() {
	m.property_types = nil
	m.removedproperty_types = nil
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (m *EquipmentTypeMutation) AddEquipmentIDs(ids ...int) {
	if m.equipment == nil {
		m.equipment = make(map[int]struct{})
	}
	for i := range ids {
		m.equipment[ids[i]] = struct{}{}
	}
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (m *EquipmentTypeMutation) RemoveEquipmentIDs(ids ...int) {
	if m.removedequipment == nil {
		m.removedequipment = make(map[int]struct{})
	}
	for i := range ids {
		m.removedequipment[ids[i]] = struct{}{}
	}
}

// RemovedEquipment returns the removed ids of equipment.
func (m *EquipmentTypeMutation) RemovedEquipmentIDs() (ids []int) {
	for id := range m.removedequipment {
		ids = append(ids, id)
	}
	return
}

// EquipmentIDs returns the equipment ids in the mutation.
func (m *EquipmentTypeMutation) EquipmentIDs() (ids []int) {
	for id := range m.equipment {
		ids = append(ids, id)
	}
	return
}

// ResetEquipment reset all changes of the equipment edge.
func (m *EquipmentTypeMutation) ResetEquipment() {
	m.equipment = nil
	m.removedequipment = nil
}

// SetCategoryID sets the category edge to EquipmentCategory by id.
func (m *EquipmentTypeMutation) SetCategoryID(id int) {
	m.category = &id
}

// ClearCategory clears the category edge to EquipmentCategory.
func (m *EquipmentTypeMutation) ClearCategory() {
	m.clearedcategory = true
}

// CategoryCleared returns if the edge category was cleared.
func (m *EquipmentTypeMutation) CategoryCleared() bool {
	return m.clearedcategory
}

// CategoryID returns the category id in the mutation.
func (m *EquipmentTypeMutation) CategoryID() (id int, exists bool) {
	if m.category != nil {
		return *m.category, true
	}
	return
}

// CategoryIDs returns the category ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// CategoryID instead. It exists only for internal usage by the builders.
func (m *EquipmentTypeMutation) CategoryIDs() (ids []int) {
	if id := m.category; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetCategory reset all changes of the category edge.
func (m *EquipmentTypeMutation) ResetCategory() {
	m.category = nil
	m.clearedcategory = false
}

// Op returns the operation name.
func (m *EquipmentTypeMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (EquipmentType).
func (m *EquipmentTypeMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *EquipmentTypeMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.create_time != nil {
		fields = append(fields, equipmenttype.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, equipmenttype.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, equipmenttype.FieldName)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *EquipmentTypeMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case equipmenttype.FieldCreateTime:
		return m.CreateTime()
	case equipmenttype.FieldUpdateTime:
		return m.UpdateTime()
	case equipmenttype.FieldName:
		return m.Name()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentTypeMutation) SetField(name string, value ent.Value) error {
	switch name {
	case equipmenttype.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case equipmenttype.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case equipmenttype.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	}
	return fmt.Errorf("unknown EquipmentType field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *EquipmentTypeMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *EquipmentTypeMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *EquipmentTypeMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown EquipmentType numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *EquipmentTypeMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *EquipmentTypeMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *EquipmentTypeMutation) ClearField(name string) error {
	return fmt.Errorf("unknown EquipmentType nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *EquipmentTypeMutation) ResetField(name string) error {
	switch name {
	case equipmenttype.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case equipmenttype.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case equipmenttype.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown EquipmentType field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *EquipmentTypeMutation) AddedEdges() []string {
	edges := make([]string, 0, 5)
	if m.port_definitions != nil {
		edges = append(edges, equipmenttype.EdgePortDefinitions)
	}
	if m.position_definitions != nil {
		edges = append(edges, equipmenttype.EdgePositionDefinitions)
	}
	if m.property_types != nil {
		edges = append(edges, equipmenttype.EdgePropertyTypes)
	}
	if m.equipment != nil {
		edges = append(edges, equipmenttype.EdgeEquipment)
	}
	if m.category != nil {
		edges = append(edges, equipmenttype.EdgeCategory)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *EquipmentTypeMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case equipmenttype.EdgePortDefinitions:
		ids := make([]ent.Value, 0, len(m.port_definitions))
		for id := range m.port_definitions {
			ids = append(ids, id)
		}
		return ids
	case equipmenttype.EdgePositionDefinitions:
		ids := make([]ent.Value, 0, len(m.position_definitions))
		for id := range m.position_definitions {
			ids = append(ids, id)
		}
		return ids
	case equipmenttype.EdgePropertyTypes:
		ids := make([]ent.Value, 0, len(m.property_types))
		for id := range m.property_types {
			ids = append(ids, id)
		}
		return ids
	case equipmenttype.EdgeEquipment:
		ids := make([]ent.Value, 0, len(m.equipment))
		for id := range m.equipment {
			ids = append(ids, id)
		}
		return ids
	case equipmenttype.EdgeCategory:
		if id := m.category; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *EquipmentTypeMutation) RemovedEdges() []string {
	edges := make([]string, 0, 5)
	if m.removedport_definitions != nil {
		edges = append(edges, equipmenttype.EdgePortDefinitions)
	}
	if m.removedposition_definitions != nil {
		edges = append(edges, equipmenttype.EdgePositionDefinitions)
	}
	if m.removedproperty_types != nil {
		edges = append(edges, equipmenttype.EdgePropertyTypes)
	}
	if m.removedequipment != nil {
		edges = append(edges, equipmenttype.EdgeEquipment)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *EquipmentTypeMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case equipmenttype.EdgePortDefinitions:
		ids := make([]ent.Value, 0, len(m.removedport_definitions))
		for id := range m.removedport_definitions {
			ids = append(ids, id)
		}
		return ids
	case equipmenttype.EdgePositionDefinitions:
		ids := make([]ent.Value, 0, len(m.removedposition_definitions))
		for id := range m.removedposition_definitions {
			ids = append(ids, id)
		}
		return ids
	case equipmenttype.EdgePropertyTypes:
		ids := make([]ent.Value, 0, len(m.removedproperty_types))
		for id := range m.removedproperty_types {
			ids = append(ids, id)
		}
		return ids
	case equipmenttype.EdgeEquipment:
		ids := make([]ent.Value, 0, len(m.removedequipment))
		for id := range m.removedequipment {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *EquipmentTypeMutation) ClearedEdges() []string {
	edges := make([]string, 0, 5)
	if m.clearedcategory {
		edges = append(edges, equipmenttype.EdgeCategory)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *EquipmentTypeMutation) EdgeCleared(name string) bool {
	switch name {
	case equipmenttype.EdgeCategory:
		return m.clearedcategory
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *EquipmentTypeMutation) ClearEdge(name string) error {
	switch name {
	case equipmenttype.EdgeCategory:
		m.ClearCategory()
		return nil
	}
	return fmt.Errorf("unknown EquipmentType unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *EquipmentTypeMutation) ResetEdge(name string) error {
	switch name {
	case equipmenttype.EdgePortDefinitions:
		m.ResetPortDefinitions()
		return nil
	case equipmenttype.EdgePositionDefinitions:
		m.ResetPositionDefinitions()
		return nil
	case equipmenttype.EdgePropertyTypes:
		m.ResetPropertyTypes()
		return nil
	case equipmenttype.EdgeEquipment:
		m.ResetEquipment()
		return nil
	case equipmenttype.EdgeCategory:
		m.ResetCategory()
		return nil
	}
	return fmt.Errorf("unknown EquipmentType edge %s", name)
}

// FileMutation represents an operation that mutate the Files
// nodes in the graph.
type FileMutation struct {
	config
	op            Op
	typ           string
	id            *int
	create_time   *time.Time
	update_time   *time.Time
	_type         *string
	name          *string
	size          *int
	addsize       *int
	modified_at   *time.Time
	uploaded_at   *time.Time
	content_type  *string
	store_key     *string
	category      *string
	clearedFields map[string]bool
}

var _ ent.Mutation = (*FileMutation)(nil)

// newFileMutation creates new mutation for $n.Name.
func newFileMutation(c config, op Op) *FileMutation {
	return &FileMutation{
		config:        c,
		op:            op,
		typ:           TypeFile,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m FileMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m FileMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *FileMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *FileMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *FileMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *FileMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *FileMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *FileMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *FileMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetType sets the type field.
func (m *FileMutation) SetType(s string) {
	m._type = &s
}

// GetType returns the type value in the mutation.
func (m *FileMutation) GetType() (r string, exists bool) {
	v := m._type
	if v == nil {
		return
	}
	return *v, true
}

// ResetType reset all changes of the type field.
func (m *FileMutation) ResetType() {
	m._type = nil
}

// SetName sets the name field.
func (m *FileMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *FileMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *FileMutation) ResetName() {
	m.name = nil
}

// SetSize sets the size field.
func (m *FileMutation) SetSize(i int) {
	m.size = &i
	m.addsize = nil
}

// Size returns the size value in the mutation.
func (m *FileMutation) Size() (r int, exists bool) {
	v := m.size
	if v == nil {
		return
	}
	return *v, true
}

// AddSize adds i to size.
func (m *FileMutation) AddSize(i int) {
	if m.addsize != nil {
		*m.addsize += i
	} else {
		m.addsize = &i
	}
}

// AddedSize returns the value that was added to the size field in this mutation.
func (m *FileMutation) AddedSize() (r int, exists bool) {
	v := m.addsize
	if v == nil {
		return
	}
	return *v, true
}

// ClearSize clears the value of size.
func (m *FileMutation) ClearSize() {
	m.size = nil
	m.addsize = nil
	m.clearedFields[file.FieldSize] = true
}

// SizeCleared returns if the field size was cleared in this mutation.
func (m *FileMutation) SizeCleared() bool {
	return m.clearedFields[file.FieldSize]
}

// ResetSize reset all changes of the size field.
func (m *FileMutation) ResetSize() {
	m.size = nil
	m.addsize = nil
	delete(m.clearedFields, file.FieldSize)
}

// SetModifiedAt sets the modified_at field.
func (m *FileMutation) SetModifiedAt(t time.Time) {
	m.modified_at = &t
}

// ModifiedAt returns the modified_at value in the mutation.
func (m *FileMutation) ModifiedAt() (r time.Time, exists bool) {
	v := m.modified_at
	if v == nil {
		return
	}
	return *v, true
}

// ClearModifiedAt clears the value of modified_at.
func (m *FileMutation) ClearModifiedAt() {
	m.modified_at = nil
	m.clearedFields[file.FieldModifiedAt] = true
}

// ModifiedAtCleared returns if the field modified_at was cleared in this mutation.
func (m *FileMutation) ModifiedAtCleared() bool {
	return m.clearedFields[file.FieldModifiedAt]
}

// ResetModifiedAt reset all changes of the modified_at field.
func (m *FileMutation) ResetModifiedAt() {
	m.modified_at = nil
	delete(m.clearedFields, file.FieldModifiedAt)
}

// SetUploadedAt sets the uploaded_at field.
func (m *FileMutation) SetUploadedAt(t time.Time) {
	m.uploaded_at = &t
}

// UploadedAt returns the uploaded_at value in the mutation.
func (m *FileMutation) UploadedAt() (r time.Time, exists bool) {
	v := m.uploaded_at
	if v == nil {
		return
	}
	return *v, true
}

// ClearUploadedAt clears the value of uploaded_at.
func (m *FileMutation) ClearUploadedAt() {
	m.uploaded_at = nil
	m.clearedFields[file.FieldUploadedAt] = true
}

// UploadedAtCleared returns if the field uploaded_at was cleared in this mutation.
func (m *FileMutation) UploadedAtCleared() bool {
	return m.clearedFields[file.FieldUploadedAt]
}

// ResetUploadedAt reset all changes of the uploaded_at field.
func (m *FileMutation) ResetUploadedAt() {
	m.uploaded_at = nil
	delete(m.clearedFields, file.FieldUploadedAt)
}

// SetContentType sets the content_type field.
func (m *FileMutation) SetContentType(s string) {
	m.content_type = &s
}

// ContentType returns the content_type value in the mutation.
func (m *FileMutation) ContentType() (r string, exists bool) {
	v := m.content_type
	if v == nil {
		return
	}
	return *v, true
}

// ResetContentType reset all changes of the content_type field.
func (m *FileMutation) ResetContentType() {
	m.content_type = nil
}

// SetStoreKey sets the store_key field.
func (m *FileMutation) SetStoreKey(s string) {
	m.store_key = &s
}

// StoreKey returns the store_key value in the mutation.
func (m *FileMutation) StoreKey() (r string, exists bool) {
	v := m.store_key
	if v == nil {
		return
	}
	return *v, true
}

// ResetStoreKey reset all changes of the store_key field.
func (m *FileMutation) ResetStoreKey() {
	m.store_key = nil
}

// SetCategory sets the category field.
func (m *FileMutation) SetCategory(s string) {
	m.category = &s
}

// Category returns the category value in the mutation.
func (m *FileMutation) Category() (r string, exists bool) {
	v := m.category
	if v == nil {
		return
	}
	return *v, true
}

// ClearCategory clears the value of category.
func (m *FileMutation) ClearCategory() {
	m.category = nil
	m.clearedFields[file.FieldCategory] = true
}

// CategoryCleared returns if the field category was cleared in this mutation.
func (m *FileMutation) CategoryCleared() bool {
	return m.clearedFields[file.FieldCategory]
}

// ResetCategory reset all changes of the category field.
func (m *FileMutation) ResetCategory() {
	m.category = nil
	delete(m.clearedFields, file.FieldCategory)
}

// Op returns the operation name.
func (m *FileMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (File).
func (m *FileMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *FileMutation) Fields() []string {
	fields := make([]string, 0, 10)
	if m.create_time != nil {
		fields = append(fields, file.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, file.FieldUpdateTime)
	}
	if m._type != nil {
		fields = append(fields, file.FieldType)
	}
	if m.name != nil {
		fields = append(fields, file.FieldName)
	}
	if m.size != nil {
		fields = append(fields, file.FieldSize)
	}
	if m.modified_at != nil {
		fields = append(fields, file.FieldModifiedAt)
	}
	if m.uploaded_at != nil {
		fields = append(fields, file.FieldUploadedAt)
	}
	if m.content_type != nil {
		fields = append(fields, file.FieldContentType)
	}
	if m.store_key != nil {
		fields = append(fields, file.FieldStoreKey)
	}
	if m.category != nil {
		fields = append(fields, file.FieldCategory)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *FileMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case file.FieldCreateTime:
		return m.CreateTime()
	case file.FieldUpdateTime:
		return m.UpdateTime()
	case file.FieldType:
		return m.GetType()
	case file.FieldName:
		return m.Name()
	case file.FieldSize:
		return m.Size()
	case file.FieldModifiedAt:
		return m.ModifiedAt()
	case file.FieldUploadedAt:
		return m.UploadedAt()
	case file.FieldContentType:
		return m.ContentType()
	case file.FieldStoreKey:
		return m.StoreKey()
	case file.FieldCategory:
		return m.Category()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *FileMutation) SetField(name string, value ent.Value) error {
	switch name {
	case file.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case file.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case file.FieldType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
		return nil
	case file.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case file.FieldSize:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSize(v)
		return nil
	case file.FieldModifiedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetModifiedAt(v)
		return nil
	case file.FieldUploadedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUploadedAt(v)
		return nil
	case file.FieldContentType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetContentType(v)
		return nil
	case file.FieldStoreKey:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStoreKey(v)
		return nil
	case file.FieldCategory:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCategory(v)
		return nil
	}
	return fmt.Errorf("unknown File field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *FileMutation) AddedFields() []string {
	var fields []string
	if m.addsize != nil {
		fields = append(fields, file.FieldSize)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *FileMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case file.FieldSize:
		return m.AddedSize()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *FileMutation) AddField(name string, value ent.Value) error {
	switch name {
	case file.FieldSize:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddSize(v)
		return nil
	}
	return fmt.Errorf("unknown File numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *FileMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[file.FieldSize] {
		fields = append(fields, file.FieldSize)
	}
	if m.clearedFields[file.FieldModifiedAt] {
		fields = append(fields, file.FieldModifiedAt)
	}
	if m.clearedFields[file.FieldUploadedAt] {
		fields = append(fields, file.FieldUploadedAt)
	}
	if m.clearedFields[file.FieldCategory] {
		fields = append(fields, file.FieldCategory)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *FileMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *FileMutation) ClearField(name string) error {
	switch name {
	case file.FieldSize:
		m.ClearSize()
		return nil
	case file.FieldModifiedAt:
		m.ClearModifiedAt()
		return nil
	case file.FieldUploadedAt:
		m.ClearUploadedAt()
		return nil
	case file.FieldCategory:
		m.ClearCategory()
		return nil
	}
	return fmt.Errorf("unknown File nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *FileMutation) ResetField(name string) error {
	switch name {
	case file.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case file.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case file.FieldType:
		m.ResetType()
		return nil
	case file.FieldName:
		m.ResetName()
		return nil
	case file.FieldSize:
		m.ResetSize()
		return nil
	case file.FieldModifiedAt:
		m.ResetModifiedAt()
		return nil
	case file.FieldUploadedAt:
		m.ResetUploadedAt()
		return nil
	case file.FieldContentType:
		m.ResetContentType()
		return nil
	case file.FieldStoreKey:
		m.ResetStoreKey()
		return nil
	case file.FieldCategory:
		m.ResetCategory()
		return nil
	}
	return fmt.Errorf("unknown File field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *FileMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *FileMutation) AddedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *FileMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *FileMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *FileMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *FileMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *FileMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown File unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *FileMutation) ResetEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown File edge %s", name)
}

// FloorPlanMutation represents an operation that mutate the FloorPlans
// nodes in the graph.
type FloorPlanMutation struct {
	config
	op                     Op
	typ                    string
	id                     *int
	create_time            *time.Time
	update_time            *time.Time
	name                   *string
	clearedFields          map[string]bool
	location               *int
	clearedlocation        bool
	reference_point        *int
	clearedreference_point bool
	scale                  *int
	clearedscale           bool
	image                  *int
	clearedimage           bool
}

var _ ent.Mutation = (*FloorPlanMutation)(nil)

// newFloorPlanMutation creates new mutation for $n.Name.
func newFloorPlanMutation(c config, op Op) *FloorPlanMutation {
	return &FloorPlanMutation{
		config:        c,
		op:            op,
		typ:           TypeFloorPlan,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m FloorPlanMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m FloorPlanMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *FloorPlanMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *FloorPlanMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *FloorPlanMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *FloorPlanMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *FloorPlanMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *FloorPlanMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *FloorPlanMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *FloorPlanMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *FloorPlanMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *FloorPlanMutation) ResetName() {
	m.name = nil
}

// SetLocationID sets the location edge to Location by id.
func (m *FloorPlanMutation) SetLocationID(id int) {
	m.location = &id
}

// ClearLocation clears the location edge to Location.
func (m *FloorPlanMutation) ClearLocation() {
	m.clearedlocation = true
}

// LocationCleared returns if the edge location was cleared.
func (m *FloorPlanMutation) LocationCleared() bool {
	return m.clearedlocation
}

// LocationID returns the location id in the mutation.
func (m *FloorPlanMutation) LocationID() (id int, exists bool) {
	if m.location != nil {
		return *m.location, true
	}
	return
}

// LocationIDs returns the location ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LocationID instead. It exists only for internal usage by the builders.
func (m *FloorPlanMutation) LocationIDs() (ids []int) {
	if id := m.location; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLocation reset all changes of the location edge.
func (m *FloorPlanMutation) ResetLocation() {
	m.location = nil
	m.clearedlocation = false
}

// SetReferencePointID sets the reference_point edge to FloorPlanReferencePoint by id.
func (m *FloorPlanMutation) SetReferencePointID(id int) {
	m.reference_point = &id
}

// ClearReferencePoint clears the reference_point edge to FloorPlanReferencePoint.
func (m *FloorPlanMutation) ClearReferencePoint() {
	m.clearedreference_point = true
}

// ReferencePointCleared returns if the edge reference_point was cleared.
func (m *FloorPlanMutation) ReferencePointCleared() bool {
	return m.clearedreference_point
}

// ReferencePointID returns the reference_point id in the mutation.
func (m *FloorPlanMutation) ReferencePointID() (id int, exists bool) {
	if m.reference_point != nil {
		return *m.reference_point, true
	}
	return
}

// ReferencePointIDs returns the reference_point ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ReferencePointID instead. It exists only for internal usage by the builders.
func (m *FloorPlanMutation) ReferencePointIDs() (ids []int) {
	if id := m.reference_point; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetReferencePoint reset all changes of the reference_point edge.
func (m *FloorPlanMutation) ResetReferencePoint() {
	m.reference_point = nil
	m.clearedreference_point = false
}

// SetScaleID sets the scale edge to FloorPlanScale by id.
func (m *FloorPlanMutation) SetScaleID(id int) {
	m.scale = &id
}

// ClearScale clears the scale edge to FloorPlanScale.
func (m *FloorPlanMutation) ClearScale() {
	m.clearedscale = true
}

// ScaleCleared returns if the edge scale was cleared.
func (m *FloorPlanMutation) ScaleCleared() bool {
	return m.clearedscale
}

// ScaleID returns the scale id in the mutation.
func (m *FloorPlanMutation) ScaleID() (id int, exists bool) {
	if m.scale != nil {
		return *m.scale, true
	}
	return
}

// ScaleIDs returns the scale ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ScaleID instead. It exists only for internal usage by the builders.
func (m *FloorPlanMutation) ScaleIDs() (ids []int) {
	if id := m.scale; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetScale reset all changes of the scale edge.
func (m *FloorPlanMutation) ResetScale() {
	m.scale = nil
	m.clearedscale = false
}

// SetImageID sets the image edge to File by id.
func (m *FloorPlanMutation) SetImageID(id int) {
	m.image = &id
}

// ClearImage clears the image edge to File.
func (m *FloorPlanMutation) ClearImage() {
	m.clearedimage = true
}

// ImageCleared returns if the edge image was cleared.
func (m *FloorPlanMutation) ImageCleared() bool {
	return m.clearedimage
}

// ImageID returns the image id in the mutation.
func (m *FloorPlanMutation) ImageID() (id int, exists bool) {
	if m.image != nil {
		return *m.image, true
	}
	return
}

// ImageIDs returns the image ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ImageID instead. It exists only for internal usage by the builders.
func (m *FloorPlanMutation) ImageIDs() (ids []int) {
	if id := m.image; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetImage reset all changes of the image edge.
func (m *FloorPlanMutation) ResetImage() {
	m.image = nil
	m.clearedimage = false
}

// Op returns the operation name.
func (m *FloorPlanMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (FloorPlan).
func (m *FloorPlanMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *FloorPlanMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.create_time != nil {
		fields = append(fields, floorplan.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, floorplan.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, floorplan.FieldName)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *FloorPlanMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case floorplan.FieldCreateTime:
		return m.CreateTime()
	case floorplan.FieldUpdateTime:
		return m.UpdateTime()
	case floorplan.FieldName:
		return m.Name()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *FloorPlanMutation) SetField(name string, value ent.Value) error {
	switch name {
	case floorplan.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case floorplan.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case floorplan.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	}
	return fmt.Errorf("unknown FloorPlan field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *FloorPlanMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *FloorPlanMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *FloorPlanMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown FloorPlan numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *FloorPlanMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *FloorPlanMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *FloorPlanMutation) ClearField(name string) error {
	return fmt.Errorf("unknown FloorPlan nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *FloorPlanMutation) ResetField(name string) error {
	switch name {
	case floorplan.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case floorplan.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case floorplan.FieldName:
		m.ResetName()
		return nil
	}
	return fmt.Errorf("unknown FloorPlan field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *FloorPlanMutation) AddedEdges() []string {
	edges := make([]string, 0, 4)
	if m.location != nil {
		edges = append(edges, floorplan.EdgeLocation)
	}
	if m.reference_point != nil {
		edges = append(edges, floorplan.EdgeReferencePoint)
	}
	if m.scale != nil {
		edges = append(edges, floorplan.EdgeScale)
	}
	if m.image != nil {
		edges = append(edges, floorplan.EdgeImage)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *FloorPlanMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case floorplan.EdgeLocation:
		if id := m.location; id != nil {
			return []ent.Value{*id}
		}
	case floorplan.EdgeReferencePoint:
		if id := m.reference_point; id != nil {
			return []ent.Value{*id}
		}
	case floorplan.EdgeScale:
		if id := m.scale; id != nil {
			return []ent.Value{*id}
		}
	case floorplan.EdgeImage:
		if id := m.image; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *FloorPlanMutation) RemovedEdges() []string {
	edges := make([]string, 0, 4)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *FloorPlanMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *FloorPlanMutation) ClearedEdges() []string {
	edges := make([]string, 0, 4)
	if m.clearedlocation {
		edges = append(edges, floorplan.EdgeLocation)
	}
	if m.clearedreference_point {
		edges = append(edges, floorplan.EdgeReferencePoint)
	}
	if m.clearedscale {
		edges = append(edges, floorplan.EdgeScale)
	}
	if m.clearedimage {
		edges = append(edges, floorplan.EdgeImage)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *FloorPlanMutation) EdgeCleared(name string) bool {
	switch name {
	case floorplan.EdgeLocation:
		return m.clearedlocation
	case floorplan.EdgeReferencePoint:
		return m.clearedreference_point
	case floorplan.EdgeScale:
		return m.clearedscale
	case floorplan.EdgeImage:
		return m.clearedimage
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *FloorPlanMutation) ClearEdge(name string) error {
	switch name {
	case floorplan.EdgeLocation:
		m.ClearLocation()
		return nil
	case floorplan.EdgeReferencePoint:
		m.ClearReferencePoint()
		return nil
	case floorplan.EdgeScale:
		m.ClearScale()
		return nil
	case floorplan.EdgeImage:
		m.ClearImage()
		return nil
	}
	return fmt.Errorf("unknown FloorPlan unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *FloorPlanMutation) ResetEdge(name string) error {
	switch name {
	case floorplan.EdgeLocation:
		m.ResetLocation()
		return nil
	case floorplan.EdgeReferencePoint:
		m.ResetReferencePoint()
		return nil
	case floorplan.EdgeScale:
		m.ResetScale()
		return nil
	case floorplan.EdgeImage:
		m.ResetImage()
		return nil
	}
	return fmt.Errorf("unknown FloorPlan edge %s", name)
}

// FloorPlanReferencePointMutation represents an operation that mutate the FloorPlanReferencePoints
// nodes in the graph.
type FloorPlanReferencePointMutation struct {
	config
	op            Op
	typ           string
	id            *int
	create_time   *time.Time
	update_time   *time.Time
	x             *int
	addx          *int
	y             *int
	addy          *int
	latitude      *float64
	addlatitude   *float64
	longitude     *float64
	addlongitude  *float64
	clearedFields map[string]bool
}

var _ ent.Mutation = (*FloorPlanReferencePointMutation)(nil)

// newFloorPlanReferencePointMutation creates new mutation for $n.Name.
func newFloorPlanReferencePointMutation(c config, op Op) *FloorPlanReferencePointMutation {
	return &FloorPlanReferencePointMutation{
		config:        c,
		op:            op,
		typ:           TypeFloorPlanReferencePoint,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m FloorPlanReferencePointMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m FloorPlanReferencePointMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *FloorPlanReferencePointMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *FloorPlanReferencePointMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *FloorPlanReferencePointMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *FloorPlanReferencePointMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *FloorPlanReferencePointMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *FloorPlanReferencePointMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *FloorPlanReferencePointMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetX sets the x field.
func (m *FloorPlanReferencePointMutation) SetX(i int) {
	m.x = &i
	m.addx = nil
}

// X returns the x value in the mutation.
func (m *FloorPlanReferencePointMutation) X() (r int, exists bool) {
	v := m.x
	if v == nil {
		return
	}
	return *v, true
}

// AddX adds i to x.
func (m *FloorPlanReferencePointMutation) AddX(i int) {
	if m.addx != nil {
		*m.addx += i
	} else {
		m.addx = &i
	}
}

// AddedX returns the value that was added to the x field in this mutation.
func (m *FloorPlanReferencePointMutation) AddedX() (r int, exists bool) {
	v := m.addx
	if v == nil {
		return
	}
	return *v, true
}

// ResetX reset all changes of the x field.
func (m *FloorPlanReferencePointMutation) ResetX() {
	m.x = nil
	m.addx = nil
}

// SetY sets the y field.
func (m *FloorPlanReferencePointMutation) SetY(i int) {
	m.y = &i
	m.addy = nil
}

// Y returns the y value in the mutation.
func (m *FloorPlanReferencePointMutation) Y() (r int, exists bool) {
	v := m.y
	if v == nil {
		return
	}
	return *v, true
}

// AddY adds i to y.
func (m *FloorPlanReferencePointMutation) AddY(i int) {
	if m.addy != nil {
		*m.addy += i
	} else {
		m.addy = &i
	}
}

// AddedY returns the value that was added to the y field in this mutation.
func (m *FloorPlanReferencePointMutation) AddedY() (r int, exists bool) {
	v := m.addy
	if v == nil {
		return
	}
	return *v, true
}

// ResetY reset all changes of the y field.
func (m *FloorPlanReferencePointMutation) ResetY() {
	m.y = nil
	m.addy = nil
}

// SetLatitude sets the latitude field.
func (m *FloorPlanReferencePointMutation) SetLatitude(f float64) {
	m.latitude = &f
	m.addlatitude = nil
}

// Latitude returns the latitude value in the mutation.
func (m *FloorPlanReferencePointMutation) Latitude() (r float64, exists bool) {
	v := m.latitude
	if v == nil {
		return
	}
	return *v, true
}

// AddLatitude adds f to latitude.
func (m *FloorPlanReferencePointMutation) AddLatitude(f float64) {
	if m.addlatitude != nil {
		*m.addlatitude += f
	} else {
		m.addlatitude = &f
	}
}

// AddedLatitude returns the value that was added to the latitude field in this mutation.
func (m *FloorPlanReferencePointMutation) AddedLatitude() (r float64, exists bool) {
	v := m.addlatitude
	if v == nil {
		return
	}
	return *v, true
}

// ResetLatitude reset all changes of the latitude field.
func (m *FloorPlanReferencePointMutation) ResetLatitude() {
	m.latitude = nil
	m.addlatitude = nil
}

// SetLongitude sets the longitude field.
func (m *FloorPlanReferencePointMutation) SetLongitude(f float64) {
	m.longitude = &f
	m.addlongitude = nil
}

// Longitude returns the longitude value in the mutation.
func (m *FloorPlanReferencePointMutation) Longitude() (r float64, exists bool) {
	v := m.longitude
	if v == nil {
		return
	}
	return *v, true
}

// AddLongitude adds f to longitude.
func (m *FloorPlanReferencePointMutation) AddLongitude(f float64) {
	if m.addlongitude != nil {
		*m.addlongitude += f
	} else {
		m.addlongitude = &f
	}
}

// AddedLongitude returns the value that was added to the longitude field in this mutation.
func (m *FloorPlanReferencePointMutation) AddedLongitude() (r float64, exists bool) {
	v := m.addlongitude
	if v == nil {
		return
	}
	return *v, true
}

// ResetLongitude reset all changes of the longitude field.
func (m *FloorPlanReferencePointMutation) ResetLongitude() {
	m.longitude = nil
	m.addlongitude = nil
}

// Op returns the operation name.
func (m *FloorPlanReferencePointMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (FloorPlanReferencePoint).
func (m *FloorPlanReferencePointMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *FloorPlanReferencePointMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m.create_time != nil {
		fields = append(fields, floorplanreferencepoint.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, floorplanreferencepoint.FieldUpdateTime)
	}
	if m.x != nil {
		fields = append(fields, floorplanreferencepoint.FieldX)
	}
	if m.y != nil {
		fields = append(fields, floorplanreferencepoint.FieldY)
	}
	if m.latitude != nil {
		fields = append(fields, floorplanreferencepoint.FieldLatitude)
	}
	if m.longitude != nil {
		fields = append(fields, floorplanreferencepoint.FieldLongitude)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *FloorPlanReferencePointMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case floorplanreferencepoint.FieldCreateTime:
		return m.CreateTime()
	case floorplanreferencepoint.FieldUpdateTime:
		return m.UpdateTime()
	case floorplanreferencepoint.FieldX:
		return m.X()
	case floorplanreferencepoint.FieldY:
		return m.Y()
	case floorplanreferencepoint.FieldLatitude:
		return m.Latitude()
	case floorplanreferencepoint.FieldLongitude:
		return m.Longitude()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *FloorPlanReferencePointMutation) SetField(name string, value ent.Value) error {
	switch name {
	case floorplanreferencepoint.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case floorplanreferencepoint.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case floorplanreferencepoint.FieldX:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetX(v)
		return nil
	case floorplanreferencepoint.FieldY:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetY(v)
		return nil
	case floorplanreferencepoint.FieldLatitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLatitude(v)
		return nil
	case floorplanreferencepoint.FieldLongitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLongitude(v)
		return nil
	}
	return fmt.Errorf("unknown FloorPlanReferencePoint field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *FloorPlanReferencePointMutation) AddedFields() []string {
	var fields []string
	if m.addx != nil {
		fields = append(fields, floorplanreferencepoint.FieldX)
	}
	if m.addy != nil {
		fields = append(fields, floorplanreferencepoint.FieldY)
	}
	if m.addlatitude != nil {
		fields = append(fields, floorplanreferencepoint.FieldLatitude)
	}
	if m.addlongitude != nil {
		fields = append(fields, floorplanreferencepoint.FieldLongitude)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *FloorPlanReferencePointMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case floorplanreferencepoint.FieldX:
		return m.AddedX()
	case floorplanreferencepoint.FieldY:
		return m.AddedY()
	case floorplanreferencepoint.FieldLatitude:
		return m.AddedLatitude()
	case floorplanreferencepoint.FieldLongitude:
		return m.AddedLongitude()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *FloorPlanReferencePointMutation) AddField(name string, value ent.Value) error {
	switch name {
	case floorplanreferencepoint.FieldX:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddX(v)
		return nil
	case floorplanreferencepoint.FieldY:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddY(v)
		return nil
	case floorplanreferencepoint.FieldLatitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLatitude(v)
		return nil
	case floorplanreferencepoint.FieldLongitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLongitude(v)
		return nil
	}
	return fmt.Errorf("unknown FloorPlanReferencePoint numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *FloorPlanReferencePointMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *FloorPlanReferencePointMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *FloorPlanReferencePointMutation) ClearField(name string) error {
	return fmt.Errorf("unknown FloorPlanReferencePoint nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *FloorPlanReferencePointMutation) ResetField(name string) error {
	switch name {
	case floorplanreferencepoint.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case floorplanreferencepoint.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case floorplanreferencepoint.FieldX:
		m.ResetX()
		return nil
	case floorplanreferencepoint.FieldY:
		m.ResetY()
		return nil
	case floorplanreferencepoint.FieldLatitude:
		m.ResetLatitude()
		return nil
	case floorplanreferencepoint.FieldLongitude:
		m.ResetLongitude()
		return nil
	}
	return fmt.Errorf("unknown FloorPlanReferencePoint field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *FloorPlanReferencePointMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *FloorPlanReferencePointMutation) AddedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *FloorPlanReferencePointMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *FloorPlanReferencePointMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *FloorPlanReferencePointMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *FloorPlanReferencePointMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *FloorPlanReferencePointMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown FloorPlanReferencePoint unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *FloorPlanReferencePointMutation) ResetEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown FloorPlanReferencePoint edge %s", name)
}

// FloorPlanScaleMutation represents an operation that mutate the FloorPlanScales
// nodes in the graph.
type FloorPlanScaleMutation struct {
	config
	op                    Op
	typ                   string
	id                    *int
	create_time           *time.Time
	update_time           *time.Time
	reference_point1_x    *int
	addreference_point1_x *int
	reference_point1_y    *int
	addreference_point1_y *int
	reference_point2_x    *int
	addreference_point2_x *int
	reference_point2_y    *int
	addreference_point2_y *int
	scale_in_meters       *float64
	addscale_in_meters    *float64
	clearedFields         map[string]bool
}

var _ ent.Mutation = (*FloorPlanScaleMutation)(nil)

// newFloorPlanScaleMutation creates new mutation for $n.Name.
func newFloorPlanScaleMutation(c config, op Op) *FloorPlanScaleMutation {
	return &FloorPlanScaleMutation{
		config:        c,
		op:            op,
		typ:           TypeFloorPlanScale,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m FloorPlanScaleMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m FloorPlanScaleMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *FloorPlanScaleMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *FloorPlanScaleMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *FloorPlanScaleMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *FloorPlanScaleMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *FloorPlanScaleMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *FloorPlanScaleMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *FloorPlanScaleMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetReferencePoint1X sets the reference_point1_x field.
func (m *FloorPlanScaleMutation) SetReferencePoint1X(i int) {
	m.reference_point1_x = &i
	m.addreference_point1_x = nil
}

// ReferencePoint1X returns the reference_point1_x value in the mutation.
func (m *FloorPlanScaleMutation) ReferencePoint1X() (r int, exists bool) {
	v := m.reference_point1_x
	if v == nil {
		return
	}
	return *v, true
}

// AddReferencePoint1X adds i to reference_point1_x.
func (m *FloorPlanScaleMutation) AddReferencePoint1X(i int) {
	if m.addreference_point1_x != nil {
		*m.addreference_point1_x += i
	} else {
		m.addreference_point1_x = &i
	}
}

// AddedReferencePoint1X returns the value that was added to the reference_point1_x field in this mutation.
func (m *FloorPlanScaleMutation) AddedReferencePoint1X() (r int, exists bool) {
	v := m.addreference_point1_x
	if v == nil {
		return
	}
	return *v, true
}

// ResetReferencePoint1X reset all changes of the reference_point1_x field.
func (m *FloorPlanScaleMutation) ResetReferencePoint1X() {
	m.reference_point1_x = nil
	m.addreference_point1_x = nil
}

// SetReferencePoint1Y sets the reference_point1_y field.
func (m *FloorPlanScaleMutation) SetReferencePoint1Y(i int) {
	m.reference_point1_y = &i
	m.addreference_point1_y = nil
}

// ReferencePoint1Y returns the reference_point1_y value in the mutation.
func (m *FloorPlanScaleMutation) ReferencePoint1Y() (r int, exists bool) {
	v := m.reference_point1_y
	if v == nil {
		return
	}
	return *v, true
}

// AddReferencePoint1Y adds i to reference_point1_y.
func (m *FloorPlanScaleMutation) AddReferencePoint1Y(i int) {
	if m.addreference_point1_y != nil {
		*m.addreference_point1_y += i
	} else {
		m.addreference_point1_y = &i
	}
}

// AddedReferencePoint1Y returns the value that was added to the reference_point1_y field in this mutation.
func (m *FloorPlanScaleMutation) AddedReferencePoint1Y() (r int, exists bool) {
	v := m.addreference_point1_y
	if v == nil {
		return
	}
	return *v, true
}

// ResetReferencePoint1Y reset all changes of the reference_point1_y field.
func (m *FloorPlanScaleMutation) ResetReferencePoint1Y() {
	m.reference_point1_y = nil
	m.addreference_point1_y = nil
}

// SetReferencePoint2X sets the reference_point2_x field.
func (m *FloorPlanScaleMutation) SetReferencePoint2X(i int) {
	m.reference_point2_x = &i
	m.addreference_point2_x = nil
}

// ReferencePoint2X returns the reference_point2_x value in the mutation.
func (m *FloorPlanScaleMutation) ReferencePoint2X() (r int, exists bool) {
	v := m.reference_point2_x
	if v == nil {
		return
	}
	return *v, true
}

// AddReferencePoint2X adds i to reference_point2_x.
func (m *FloorPlanScaleMutation) AddReferencePoint2X(i int) {
	if m.addreference_point2_x != nil {
		*m.addreference_point2_x += i
	} else {
		m.addreference_point2_x = &i
	}
}

// AddedReferencePoint2X returns the value that was added to the reference_point2_x field in this mutation.
func (m *FloorPlanScaleMutation) AddedReferencePoint2X() (r int, exists bool) {
	v := m.addreference_point2_x
	if v == nil {
		return
	}
	return *v, true
}

// ResetReferencePoint2X reset all changes of the reference_point2_x field.
func (m *FloorPlanScaleMutation) ResetReferencePoint2X() {
	m.reference_point2_x = nil
	m.addreference_point2_x = nil
}

// SetReferencePoint2Y sets the reference_point2_y field.
func (m *FloorPlanScaleMutation) SetReferencePoint2Y(i int) {
	m.reference_point2_y = &i
	m.addreference_point2_y = nil
}

// ReferencePoint2Y returns the reference_point2_y value in the mutation.
func (m *FloorPlanScaleMutation) ReferencePoint2Y() (r int, exists bool) {
	v := m.reference_point2_y
	if v == nil {
		return
	}
	return *v, true
}

// AddReferencePoint2Y adds i to reference_point2_y.
func (m *FloorPlanScaleMutation) AddReferencePoint2Y(i int) {
	if m.addreference_point2_y != nil {
		*m.addreference_point2_y += i
	} else {
		m.addreference_point2_y = &i
	}
}

// AddedReferencePoint2Y returns the value that was added to the reference_point2_y field in this mutation.
func (m *FloorPlanScaleMutation) AddedReferencePoint2Y() (r int, exists bool) {
	v := m.addreference_point2_y
	if v == nil {
		return
	}
	return *v, true
}

// ResetReferencePoint2Y reset all changes of the reference_point2_y field.
func (m *FloorPlanScaleMutation) ResetReferencePoint2Y() {
	m.reference_point2_y = nil
	m.addreference_point2_y = nil
}

// SetScaleInMeters sets the scale_in_meters field.
func (m *FloorPlanScaleMutation) SetScaleInMeters(f float64) {
	m.scale_in_meters = &f
	m.addscale_in_meters = nil
}

// ScaleInMeters returns the scale_in_meters value in the mutation.
func (m *FloorPlanScaleMutation) ScaleInMeters() (r float64, exists bool) {
	v := m.scale_in_meters
	if v == nil {
		return
	}
	return *v, true
}

// AddScaleInMeters adds f to scale_in_meters.
func (m *FloorPlanScaleMutation) AddScaleInMeters(f float64) {
	if m.addscale_in_meters != nil {
		*m.addscale_in_meters += f
	} else {
		m.addscale_in_meters = &f
	}
}

// AddedScaleInMeters returns the value that was added to the scale_in_meters field in this mutation.
func (m *FloorPlanScaleMutation) AddedScaleInMeters() (r float64, exists bool) {
	v := m.addscale_in_meters
	if v == nil {
		return
	}
	return *v, true
}

// ResetScaleInMeters reset all changes of the scale_in_meters field.
func (m *FloorPlanScaleMutation) ResetScaleInMeters() {
	m.scale_in_meters = nil
	m.addscale_in_meters = nil
}

// Op returns the operation name.
func (m *FloorPlanScaleMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (FloorPlanScale).
func (m *FloorPlanScaleMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *FloorPlanScaleMutation) Fields() []string {
	fields := make([]string, 0, 7)
	if m.create_time != nil {
		fields = append(fields, floorplanscale.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, floorplanscale.FieldUpdateTime)
	}
	if m.reference_point1_x != nil {
		fields = append(fields, floorplanscale.FieldReferencePoint1X)
	}
	if m.reference_point1_y != nil {
		fields = append(fields, floorplanscale.FieldReferencePoint1Y)
	}
	if m.reference_point2_x != nil {
		fields = append(fields, floorplanscale.FieldReferencePoint2X)
	}
	if m.reference_point2_y != nil {
		fields = append(fields, floorplanscale.FieldReferencePoint2Y)
	}
	if m.scale_in_meters != nil {
		fields = append(fields, floorplanscale.FieldScaleInMeters)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *FloorPlanScaleMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case floorplanscale.FieldCreateTime:
		return m.CreateTime()
	case floorplanscale.FieldUpdateTime:
		return m.UpdateTime()
	case floorplanscale.FieldReferencePoint1X:
		return m.ReferencePoint1X()
	case floorplanscale.FieldReferencePoint1Y:
		return m.ReferencePoint1Y()
	case floorplanscale.FieldReferencePoint2X:
		return m.ReferencePoint2X()
	case floorplanscale.FieldReferencePoint2Y:
		return m.ReferencePoint2Y()
	case floorplanscale.FieldScaleInMeters:
		return m.ScaleInMeters()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *FloorPlanScaleMutation) SetField(name string, value ent.Value) error {
	switch name {
	case floorplanscale.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case floorplanscale.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case floorplanscale.FieldReferencePoint1X:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetReferencePoint1X(v)
		return nil
	case floorplanscale.FieldReferencePoint1Y:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetReferencePoint1Y(v)
		return nil
	case floorplanscale.FieldReferencePoint2X:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetReferencePoint2X(v)
		return nil
	case floorplanscale.FieldReferencePoint2Y:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetReferencePoint2Y(v)
		return nil
	case floorplanscale.FieldScaleInMeters:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetScaleInMeters(v)
		return nil
	}
	return fmt.Errorf("unknown FloorPlanScale field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *FloorPlanScaleMutation) AddedFields() []string {
	var fields []string
	if m.addreference_point1_x != nil {
		fields = append(fields, floorplanscale.FieldReferencePoint1X)
	}
	if m.addreference_point1_y != nil {
		fields = append(fields, floorplanscale.FieldReferencePoint1Y)
	}
	if m.addreference_point2_x != nil {
		fields = append(fields, floorplanscale.FieldReferencePoint2X)
	}
	if m.addreference_point2_y != nil {
		fields = append(fields, floorplanscale.FieldReferencePoint2Y)
	}
	if m.addscale_in_meters != nil {
		fields = append(fields, floorplanscale.FieldScaleInMeters)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *FloorPlanScaleMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case floorplanscale.FieldReferencePoint1X:
		return m.AddedReferencePoint1X()
	case floorplanscale.FieldReferencePoint1Y:
		return m.AddedReferencePoint1Y()
	case floorplanscale.FieldReferencePoint2X:
		return m.AddedReferencePoint2X()
	case floorplanscale.FieldReferencePoint2Y:
		return m.AddedReferencePoint2Y()
	case floorplanscale.FieldScaleInMeters:
		return m.AddedScaleInMeters()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *FloorPlanScaleMutation) AddField(name string, value ent.Value) error {
	switch name {
	case floorplanscale.FieldReferencePoint1X:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddReferencePoint1X(v)
		return nil
	case floorplanscale.FieldReferencePoint1Y:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddReferencePoint1Y(v)
		return nil
	case floorplanscale.FieldReferencePoint2X:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddReferencePoint2X(v)
		return nil
	case floorplanscale.FieldReferencePoint2Y:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddReferencePoint2Y(v)
		return nil
	case floorplanscale.FieldScaleInMeters:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddScaleInMeters(v)
		return nil
	}
	return fmt.Errorf("unknown FloorPlanScale numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *FloorPlanScaleMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *FloorPlanScaleMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *FloorPlanScaleMutation) ClearField(name string) error {
	return fmt.Errorf("unknown FloorPlanScale nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *FloorPlanScaleMutation) ResetField(name string) error {
	switch name {
	case floorplanscale.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case floorplanscale.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case floorplanscale.FieldReferencePoint1X:
		m.ResetReferencePoint1X()
		return nil
	case floorplanscale.FieldReferencePoint1Y:
		m.ResetReferencePoint1Y()
		return nil
	case floorplanscale.FieldReferencePoint2X:
		m.ResetReferencePoint2X()
		return nil
	case floorplanscale.FieldReferencePoint2Y:
		m.ResetReferencePoint2Y()
		return nil
	case floorplanscale.FieldScaleInMeters:
		m.ResetScaleInMeters()
		return nil
	}
	return fmt.Errorf("unknown FloorPlanScale field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *FloorPlanScaleMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *FloorPlanScaleMutation) AddedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *FloorPlanScaleMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *FloorPlanScaleMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *FloorPlanScaleMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *FloorPlanScaleMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *FloorPlanScaleMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown FloorPlanScale unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *FloorPlanScaleMutation) ResetEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown FloorPlanScale edge %s", name)
}

// HyperlinkMutation represents an operation that mutate the Hyperlinks
// nodes in the graph.
type HyperlinkMutation struct {
	config
	op            Op
	typ           string
	id            *int
	create_time   *time.Time
	update_time   *time.Time
	url           *string
	name          *string
	category      *string
	clearedFields map[string]bool
}

var _ ent.Mutation = (*HyperlinkMutation)(nil)

// newHyperlinkMutation creates new mutation for $n.Name.
func newHyperlinkMutation(c config, op Op) *HyperlinkMutation {
	return &HyperlinkMutation{
		config:        c,
		op:            op,
		typ:           TypeHyperlink,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m HyperlinkMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m HyperlinkMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *HyperlinkMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *HyperlinkMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *HyperlinkMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *HyperlinkMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *HyperlinkMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *HyperlinkMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *HyperlinkMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetURL sets the url field.
func (m *HyperlinkMutation) SetURL(s string) {
	m.url = &s
}

// URL returns the url value in the mutation.
func (m *HyperlinkMutation) URL() (r string, exists bool) {
	v := m.url
	if v == nil {
		return
	}
	return *v, true
}

// ResetURL reset all changes of the url field.
func (m *HyperlinkMutation) ResetURL() {
	m.url = nil
}

// SetName sets the name field.
func (m *HyperlinkMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *HyperlinkMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ClearName clears the value of name.
func (m *HyperlinkMutation) ClearName() {
	m.name = nil
	m.clearedFields[hyperlink.FieldName] = true
}

// NameCleared returns if the field name was cleared in this mutation.
func (m *HyperlinkMutation) NameCleared() bool {
	return m.clearedFields[hyperlink.FieldName]
}

// ResetName reset all changes of the name field.
func (m *HyperlinkMutation) ResetName() {
	m.name = nil
	delete(m.clearedFields, hyperlink.FieldName)
}

// SetCategory sets the category field.
func (m *HyperlinkMutation) SetCategory(s string) {
	m.category = &s
}

// Category returns the category value in the mutation.
func (m *HyperlinkMutation) Category() (r string, exists bool) {
	v := m.category
	if v == nil {
		return
	}
	return *v, true
}

// ClearCategory clears the value of category.
func (m *HyperlinkMutation) ClearCategory() {
	m.category = nil
	m.clearedFields[hyperlink.FieldCategory] = true
}

// CategoryCleared returns if the field category was cleared in this mutation.
func (m *HyperlinkMutation) CategoryCleared() bool {
	return m.clearedFields[hyperlink.FieldCategory]
}

// ResetCategory reset all changes of the category field.
func (m *HyperlinkMutation) ResetCategory() {
	m.category = nil
	delete(m.clearedFields, hyperlink.FieldCategory)
}

// Op returns the operation name.
func (m *HyperlinkMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Hyperlink).
func (m *HyperlinkMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *HyperlinkMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.create_time != nil {
		fields = append(fields, hyperlink.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, hyperlink.FieldUpdateTime)
	}
	if m.url != nil {
		fields = append(fields, hyperlink.FieldURL)
	}
	if m.name != nil {
		fields = append(fields, hyperlink.FieldName)
	}
	if m.category != nil {
		fields = append(fields, hyperlink.FieldCategory)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *HyperlinkMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case hyperlink.FieldCreateTime:
		return m.CreateTime()
	case hyperlink.FieldUpdateTime:
		return m.UpdateTime()
	case hyperlink.FieldURL:
		return m.URL()
	case hyperlink.FieldName:
		return m.Name()
	case hyperlink.FieldCategory:
		return m.Category()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *HyperlinkMutation) SetField(name string, value ent.Value) error {
	switch name {
	case hyperlink.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case hyperlink.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case hyperlink.FieldURL:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetURL(v)
		return nil
	case hyperlink.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case hyperlink.FieldCategory:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCategory(v)
		return nil
	}
	return fmt.Errorf("unknown Hyperlink field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *HyperlinkMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *HyperlinkMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *HyperlinkMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Hyperlink numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *HyperlinkMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[hyperlink.FieldName] {
		fields = append(fields, hyperlink.FieldName)
	}
	if m.clearedFields[hyperlink.FieldCategory] {
		fields = append(fields, hyperlink.FieldCategory)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *HyperlinkMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *HyperlinkMutation) ClearField(name string) error {
	switch name {
	case hyperlink.FieldName:
		m.ClearName()
		return nil
	case hyperlink.FieldCategory:
		m.ClearCategory()
		return nil
	}
	return fmt.Errorf("unknown Hyperlink nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *HyperlinkMutation) ResetField(name string) error {
	switch name {
	case hyperlink.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case hyperlink.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case hyperlink.FieldURL:
		m.ResetURL()
		return nil
	case hyperlink.FieldName:
		m.ResetName()
		return nil
	case hyperlink.FieldCategory:
		m.ResetCategory()
		return nil
	}
	return fmt.Errorf("unknown Hyperlink field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *HyperlinkMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *HyperlinkMutation) AddedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *HyperlinkMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *HyperlinkMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *HyperlinkMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *HyperlinkMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *HyperlinkMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Hyperlink unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *HyperlinkMutation) ResetEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Hyperlink edge %s", name)
}

// LinkMutation represents an operation that mutate the Links
// nodes in the graph.
type LinkMutation struct {
	config
	op                Op
	typ               string
	id                *int
	create_time       *time.Time
	update_time       *time.Time
	future_state      *string
	clearedFields     map[string]bool
	ports             map[int]struct{}
	removedports      map[int]struct{}
	work_order        *int
	clearedwork_order bool
	properties        map[int]struct{}
	removedproperties map[int]struct{}
	service           map[int]struct{}
	removedservice    map[int]struct{}
}

var _ ent.Mutation = (*LinkMutation)(nil)

// newLinkMutation creates new mutation for $n.Name.
func newLinkMutation(c config, op Op) *LinkMutation {
	return &LinkMutation{
		config:        c,
		op:            op,
		typ:           TypeLink,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m LinkMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m LinkMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *LinkMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *LinkMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *LinkMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *LinkMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *LinkMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *LinkMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *LinkMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetFutureState sets the future_state field.
func (m *LinkMutation) SetFutureState(s string) {
	m.future_state = &s
}

// FutureState returns the future_state value in the mutation.
func (m *LinkMutation) FutureState() (r string, exists bool) {
	v := m.future_state
	if v == nil {
		return
	}
	return *v, true
}

// ClearFutureState clears the value of future_state.
func (m *LinkMutation) ClearFutureState() {
	m.future_state = nil
	m.clearedFields[link.FieldFutureState] = true
}

// FutureStateCleared returns if the field future_state was cleared in this mutation.
func (m *LinkMutation) FutureStateCleared() bool {
	return m.clearedFields[link.FieldFutureState]
}

// ResetFutureState reset all changes of the future_state field.
func (m *LinkMutation) ResetFutureState() {
	m.future_state = nil
	delete(m.clearedFields, link.FieldFutureState)
}

// AddPortIDs adds the ports edge to EquipmentPort by ids.
func (m *LinkMutation) AddPortIDs(ids ...int) {
	if m.ports == nil {
		m.ports = make(map[int]struct{})
	}
	for i := range ids {
		m.ports[ids[i]] = struct{}{}
	}
}

// RemovePortIDs removes the ports edge to EquipmentPort by ids.
func (m *LinkMutation) RemovePortIDs(ids ...int) {
	if m.removedports == nil {
		m.removedports = make(map[int]struct{})
	}
	for i := range ids {
		m.removedports[ids[i]] = struct{}{}
	}
}

// RemovedPorts returns the removed ids of ports.
func (m *LinkMutation) RemovedPortsIDs() (ids []int) {
	for id := range m.removedports {
		ids = append(ids, id)
	}
	return
}

// PortsIDs returns the ports ids in the mutation.
func (m *LinkMutation) PortsIDs() (ids []int) {
	for id := range m.ports {
		ids = append(ids, id)
	}
	return
}

// ResetPorts reset all changes of the ports edge.
func (m *LinkMutation) ResetPorts() {
	m.ports = nil
	m.removedports = nil
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (m *LinkMutation) SetWorkOrderID(id int) {
	m.work_order = &id
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (m *LinkMutation) ClearWorkOrder() {
	m.clearedwork_order = true
}

// WorkOrderCleared returns if the edge work_order was cleared.
func (m *LinkMutation) WorkOrderCleared() bool {
	return m.clearedwork_order
}

// WorkOrderID returns the work_order id in the mutation.
func (m *LinkMutation) WorkOrderID() (id int, exists bool) {
	if m.work_order != nil {
		return *m.work_order, true
	}
	return
}

// WorkOrderIDs returns the work_order ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// WorkOrderID instead. It exists only for internal usage by the builders.
func (m *LinkMutation) WorkOrderIDs() (ids []int) {
	if id := m.work_order; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetWorkOrder reset all changes of the work_order edge.
func (m *LinkMutation) ResetWorkOrder() {
	m.work_order = nil
	m.clearedwork_order = false
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (m *LinkMutation) AddPropertyIDs(ids ...int) {
	if m.properties == nil {
		m.properties = make(map[int]struct{})
	}
	for i := range ids {
		m.properties[ids[i]] = struct{}{}
	}
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (m *LinkMutation) RemovePropertyIDs(ids ...int) {
	if m.removedproperties == nil {
		m.removedproperties = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperties[ids[i]] = struct{}{}
	}
}

// RemovedProperties returns the removed ids of properties.
func (m *LinkMutation) RemovedPropertiesIDs() (ids []int) {
	for id := range m.removedproperties {
		ids = append(ids, id)
	}
	return
}

// PropertiesIDs returns the properties ids in the mutation.
func (m *LinkMutation) PropertiesIDs() (ids []int) {
	for id := range m.properties {
		ids = append(ids, id)
	}
	return
}

// ResetProperties reset all changes of the properties edge.
func (m *LinkMutation) ResetProperties() {
	m.properties = nil
	m.removedproperties = nil
}

// AddServiceIDs adds the service edge to Service by ids.
func (m *LinkMutation) AddServiceIDs(ids ...int) {
	if m.service == nil {
		m.service = make(map[int]struct{})
	}
	for i := range ids {
		m.service[ids[i]] = struct{}{}
	}
}

// RemoveServiceIDs removes the service edge to Service by ids.
func (m *LinkMutation) RemoveServiceIDs(ids ...int) {
	if m.removedservice == nil {
		m.removedservice = make(map[int]struct{})
	}
	for i := range ids {
		m.removedservice[ids[i]] = struct{}{}
	}
}

// RemovedService returns the removed ids of service.
func (m *LinkMutation) RemovedServiceIDs() (ids []int) {
	for id := range m.removedservice {
		ids = append(ids, id)
	}
	return
}

// ServiceIDs returns the service ids in the mutation.
func (m *LinkMutation) ServiceIDs() (ids []int) {
	for id := range m.service {
		ids = append(ids, id)
	}
	return
}

// ResetService reset all changes of the service edge.
func (m *LinkMutation) ResetService() {
	m.service = nil
	m.removedservice = nil
}

// Op returns the operation name.
func (m *LinkMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Link).
func (m *LinkMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *LinkMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.create_time != nil {
		fields = append(fields, link.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, link.FieldUpdateTime)
	}
	if m.future_state != nil {
		fields = append(fields, link.FieldFutureState)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *LinkMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case link.FieldCreateTime:
		return m.CreateTime()
	case link.FieldUpdateTime:
		return m.UpdateTime()
	case link.FieldFutureState:
		return m.FutureState()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *LinkMutation) SetField(name string, value ent.Value) error {
	switch name {
	case link.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case link.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case link.FieldFutureState:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFutureState(v)
		return nil
	}
	return fmt.Errorf("unknown Link field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *LinkMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *LinkMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *LinkMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Link numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *LinkMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[link.FieldFutureState] {
		fields = append(fields, link.FieldFutureState)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *LinkMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *LinkMutation) ClearField(name string) error {
	switch name {
	case link.FieldFutureState:
		m.ClearFutureState()
		return nil
	}
	return fmt.Errorf("unknown Link nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *LinkMutation) ResetField(name string) error {
	switch name {
	case link.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case link.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case link.FieldFutureState:
		m.ResetFutureState()
		return nil
	}
	return fmt.Errorf("unknown Link field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *LinkMutation) AddedEdges() []string {
	edges := make([]string, 0, 4)
	if m.ports != nil {
		edges = append(edges, link.EdgePorts)
	}
	if m.work_order != nil {
		edges = append(edges, link.EdgeWorkOrder)
	}
	if m.properties != nil {
		edges = append(edges, link.EdgeProperties)
	}
	if m.service != nil {
		edges = append(edges, link.EdgeService)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *LinkMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case link.EdgePorts:
		ids := make([]ent.Value, 0, len(m.ports))
		for id := range m.ports {
			ids = append(ids, id)
		}
		return ids
	case link.EdgeWorkOrder:
		if id := m.work_order; id != nil {
			return []ent.Value{*id}
		}
	case link.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.properties))
		for id := range m.properties {
			ids = append(ids, id)
		}
		return ids
	case link.EdgeService:
		ids := make([]ent.Value, 0, len(m.service))
		for id := range m.service {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *LinkMutation) RemovedEdges() []string {
	edges := make([]string, 0, 4)
	if m.removedports != nil {
		edges = append(edges, link.EdgePorts)
	}
	if m.removedproperties != nil {
		edges = append(edges, link.EdgeProperties)
	}
	if m.removedservice != nil {
		edges = append(edges, link.EdgeService)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *LinkMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case link.EdgePorts:
		ids := make([]ent.Value, 0, len(m.removedports))
		for id := range m.removedports {
			ids = append(ids, id)
		}
		return ids
	case link.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.removedproperties))
		for id := range m.removedproperties {
			ids = append(ids, id)
		}
		return ids
	case link.EdgeService:
		ids := make([]ent.Value, 0, len(m.removedservice))
		for id := range m.removedservice {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *LinkMutation) ClearedEdges() []string {
	edges := make([]string, 0, 4)
	if m.clearedwork_order {
		edges = append(edges, link.EdgeWorkOrder)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *LinkMutation) EdgeCleared(name string) bool {
	switch name {
	case link.EdgeWorkOrder:
		return m.clearedwork_order
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *LinkMutation) ClearEdge(name string) error {
	switch name {
	case link.EdgeWorkOrder:
		m.ClearWorkOrder()
		return nil
	}
	return fmt.Errorf("unknown Link unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *LinkMutation) ResetEdge(name string) error {
	switch name {
	case link.EdgePorts:
		m.ResetPorts()
		return nil
	case link.EdgeWorkOrder:
		m.ResetWorkOrder()
		return nil
	case link.EdgeProperties:
		m.ResetProperties()
		return nil
	case link.EdgeService:
		m.ResetService()
		return nil
	}
	return fmt.Errorf("unknown Link edge %s", name)
}

// LocationMutation represents an operation that mutate the Locations
// nodes in the graph.
type LocationMutation struct {
	config
	op                 Op
	typ                string
	id                 *int
	create_time        *time.Time
	update_time        *time.Time
	name               *string
	external_id        *string
	latitude           *float64
	addlatitude        *float64
	longitude          *float64
	addlongitude       *float64
	site_survey_needed *bool
	clearedFields      map[string]bool
	_type              *int
	cleared_type       bool
	parent             *int
	clearedparent      bool
	children           map[int]struct{}
	removedchildren    map[int]struct{}
	files              map[int]struct{}
	removedfiles       map[int]struct{}
	hyperlinks         map[int]struct{}
	removedhyperlinks  map[int]struct{}
	equipment          map[int]struct{}
	removedequipment   map[int]struct{}
	properties         map[int]struct{}
	removedproperties  map[int]struct{}
	survey             map[int]struct{}
	removedsurvey      map[int]struct{}
	wifi_scan          map[int]struct{}
	removedwifi_scan   map[int]struct{}
	cell_scan          map[int]struct{}
	removedcell_scan   map[int]struct{}
	work_orders        map[int]struct{}
	removedwork_orders map[int]struct{}
	floor_plans        map[int]struct{}
	removedfloor_plans map[int]struct{}
}

var _ ent.Mutation = (*LocationMutation)(nil)

// newLocationMutation creates new mutation for $n.Name.
func newLocationMutation(c config, op Op) *LocationMutation {
	return &LocationMutation{
		config:        c,
		op:            op,
		typ:           TypeLocation,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m LocationMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m LocationMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *LocationMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *LocationMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *LocationMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *LocationMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *LocationMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *LocationMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *LocationMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *LocationMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *LocationMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *LocationMutation) ResetName() {
	m.name = nil
}

// SetExternalID sets the external_id field.
func (m *LocationMutation) SetExternalID(s string) {
	m.external_id = &s
}

// ExternalID returns the external_id value in the mutation.
func (m *LocationMutation) ExternalID() (r string, exists bool) {
	v := m.external_id
	if v == nil {
		return
	}
	return *v, true
}

// ClearExternalID clears the value of external_id.
func (m *LocationMutation) ClearExternalID() {
	m.external_id = nil
	m.clearedFields[location.FieldExternalID] = true
}

// ExternalIDCleared returns if the field external_id was cleared in this mutation.
func (m *LocationMutation) ExternalIDCleared() bool {
	return m.clearedFields[location.FieldExternalID]
}

// ResetExternalID reset all changes of the external_id field.
func (m *LocationMutation) ResetExternalID() {
	m.external_id = nil
	delete(m.clearedFields, location.FieldExternalID)
}

// SetLatitude sets the latitude field.
func (m *LocationMutation) SetLatitude(f float64) {
	m.latitude = &f
	m.addlatitude = nil
}

// Latitude returns the latitude value in the mutation.
func (m *LocationMutation) Latitude() (r float64, exists bool) {
	v := m.latitude
	if v == nil {
		return
	}
	return *v, true
}

// AddLatitude adds f to latitude.
func (m *LocationMutation) AddLatitude(f float64) {
	if m.addlatitude != nil {
		*m.addlatitude += f
	} else {
		m.addlatitude = &f
	}
}

// AddedLatitude returns the value that was added to the latitude field in this mutation.
func (m *LocationMutation) AddedLatitude() (r float64, exists bool) {
	v := m.addlatitude
	if v == nil {
		return
	}
	return *v, true
}

// ResetLatitude reset all changes of the latitude field.
func (m *LocationMutation) ResetLatitude() {
	m.latitude = nil
	m.addlatitude = nil
}

// SetLongitude sets the longitude field.
func (m *LocationMutation) SetLongitude(f float64) {
	m.longitude = &f
	m.addlongitude = nil
}

// Longitude returns the longitude value in the mutation.
func (m *LocationMutation) Longitude() (r float64, exists bool) {
	v := m.longitude
	if v == nil {
		return
	}
	return *v, true
}

// AddLongitude adds f to longitude.
func (m *LocationMutation) AddLongitude(f float64) {
	if m.addlongitude != nil {
		*m.addlongitude += f
	} else {
		m.addlongitude = &f
	}
}

// AddedLongitude returns the value that was added to the longitude field in this mutation.
func (m *LocationMutation) AddedLongitude() (r float64, exists bool) {
	v := m.addlongitude
	if v == nil {
		return
	}
	return *v, true
}

// ResetLongitude reset all changes of the longitude field.
func (m *LocationMutation) ResetLongitude() {
	m.longitude = nil
	m.addlongitude = nil
}

// SetSiteSurveyNeeded sets the site_survey_needed field.
func (m *LocationMutation) SetSiteSurveyNeeded(b bool) {
	m.site_survey_needed = &b
}

// SiteSurveyNeeded returns the site_survey_needed value in the mutation.
func (m *LocationMutation) SiteSurveyNeeded() (r bool, exists bool) {
	v := m.site_survey_needed
	if v == nil {
		return
	}
	return *v, true
}

// ClearSiteSurveyNeeded clears the value of site_survey_needed.
func (m *LocationMutation) ClearSiteSurveyNeeded() {
	m.site_survey_needed = nil
	m.clearedFields[location.FieldSiteSurveyNeeded] = true
}

// SiteSurveyNeededCleared returns if the field site_survey_needed was cleared in this mutation.
func (m *LocationMutation) SiteSurveyNeededCleared() bool {
	return m.clearedFields[location.FieldSiteSurveyNeeded]
}

// ResetSiteSurveyNeeded reset all changes of the site_survey_needed field.
func (m *LocationMutation) ResetSiteSurveyNeeded() {
	m.site_survey_needed = nil
	delete(m.clearedFields, location.FieldSiteSurveyNeeded)
}

// SetTypeID sets the type edge to LocationType by id.
func (m *LocationMutation) SetTypeID(id int) {
	m._type = &id
}

// ClearType clears the type edge to LocationType.
func (m *LocationMutation) ClearType() {
	m.cleared_type = true
}

// TypeCleared returns if the edge type was cleared.
func (m *LocationMutation) TypeCleared() bool {
	return m.cleared_type
}

// TypeID returns the type id in the mutation.
func (m *LocationMutation) TypeID() (id int, exists bool) {
	if m._type != nil {
		return *m._type, true
	}
	return
}

// TypeIDs returns the type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// TypeID instead. It exists only for internal usage by the builders.
func (m *LocationMutation) TypeIDs() (ids []int) {
	if id := m._type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetType reset all changes of the type edge.
func (m *LocationMutation) ResetType() {
	m._type = nil
	m.cleared_type = false
}

// SetParentID sets the parent edge to Location by id.
func (m *LocationMutation) SetParentID(id int) {
	m.parent = &id
}

// ClearParent clears the parent edge to Location.
func (m *LocationMutation) ClearParent() {
	m.clearedparent = true
}

// ParentCleared returns if the edge parent was cleared.
func (m *LocationMutation) ParentCleared() bool {
	return m.clearedparent
}

// ParentID returns the parent id in the mutation.
func (m *LocationMutation) ParentID() (id int, exists bool) {
	if m.parent != nil {
		return *m.parent, true
	}
	return
}

// ParentIDs returns the parent ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ParentID instead. It exists only for internal usage by the builders.
func (m *LocationMutation) ParentIDs() (ids []int) {
	if id := m.parent; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetParent reset all changes of the parent edge.
func (m *LocationMutation) ResetParent() {
	m.parent = nil
	m.clearedparent = false
}

// AddChildIDs adds the children edge to Location by ids.
func (m *LocationMutation) AddChildIDs(ids ...int) {
	if m.children == nil {
		m.children = make(map[int]struct{})
	}
	for i := range ids {
		m.children[ids[i]] = struct{}{}
	}
}

// RemoveChildIDs removes the children edge to Location by ids.
func (m *LocationMutation) RemoveChildIDs(ids ...int) {
	if m.removedchildren == nil {
		m.removedchildren = make(map[int]struct{})
	}
	for i := range ids {
		m.removedchildren[ids[i]] = struct{}{}
	}
}

// RemovedChildren returns the removed ids of children.
func (m *LocationMutation) RemovedChildrenIDs() (ids []int) {
	for id := range m.removedchildren {
		ids = append(ids, id)
	}
	return
}

// ChildrenIDs returns the children ids in the mutation.
func (m *LocationMutation) ChildrenIDs() (ids []int) {
	for id := range m.children {
		ids = append(ids, id)
	}
	return
}

// ResetChildren reset all changes of the children edge.
func (m *LocationMutation) ResetChildren() {
	m.children = nil
	m.removedchildren = nil
}

// AddFileIDs adds the files edge to File by ids.
func (m *LocationMutation) AddFileIDs(ids ...int) {
	if m.files == nil {
		m.files = make(map[int]struct{})
	}
	for i := range ids {
		m.files[ids[i]] = struct{}{}
	}
}

// RemoveFileIDs removes the files edge to File by ids.
func (m *LocationMutation) RemoveFileIDs(ids ...int) {
	if m.removedfiles == nil {
		m.removedfiles = make(map[int]struct{})
	}
	for i := range ids {
		m.removedfiles[ids[i]] = struct{}{}
	}
}

// RemovedFiles returns the removed ids of files.
func (m *LocationMutation) RemovedFilesIDs() (ids []int) {
	for id := range m.removedfiles {
		ids = append(ids, id)
	}
	return
}

// FilesIDs returns the files ids in the mutation.
func (m *LocationMutation) FilesIDs() (ids []int) {
	for id := range m.files {
		ids = append(ids, id)
	}
	return
}

// ResetFiles reset all changes of the files edge.
func (m *LocationMutation) ResetFiles() {
	m.files = nil
	m.removedfiles = nil
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (m *LocationMutation) AddHyperlinkIDs(ids ...int) {
	if m.hyperlinks == nil {
		m.hyperlinks = make(map[int]struct{})
	}
	for i := range ids {
		m.hyperlinks[ids[i]] = struct{}{}
	}
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (m *LocationMutation) RemoveHyperlinkIDs(ids ...int) {
	if m.removedhyperlinks == nil {
		m.removedhyperlinks = make(map[int]struct{})
	}
	for i := range ids {
		m.removedhyperlinks[ids[i]] = struct{}{}
	}
}

// RemovedHyperlinks returns the removed ids of hyperlinks.
func (m *LocationMutation) RemovedHyperlinksIDs() (ids []int) {
	for id := range m.removedhyperlinks {
		ids = append(ids, id)
	}
	return
}

// HyperlinksIDs returns the hyperlinks ids in the mutation.
func (m *LocationMutation) HyperlinksIDs() (ids []int) {
	for id := range m.hyperlinks {
		ids = append(ids, id)
	}
	return
}

// ResetHyperlinks reset all changes of the hyperlinks edge.
func (m *LocationMutation) ResetHyperlinks() {
	m.hyperlinks = nil
	m.removedhyperlinks = nil
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (m *LocationMutation) AddEquipmentIDs(ids ...int) {
	if m.equipment == nil {
		m.equipment = make(map[int]struct{})
	}
	for i := range ids {
		m.equipment[ids[i]] = struct{}{}
	}
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (m *LocationMutation) RemoveEquipmentIDs(ids ...int) {
	if m.removedequipment == nil {
		m.removedequipment = make(map[int]struct{})
	}
	for i := range ids {
		m.removedequipment[ids[i]] = struct{}{}
	}
}

// RemovedEquipment returns the removed ids of equipment.
func (m *LocationMutation) RemovedEquipmentIDs() (ids []int) {
	for id := range m.removedequipment {
		ids = append(ids, id)
	}
	return
}

// EquipmentIDs returns the equipment ids in the mutation.
func (m *LocationMutation) EquipmentIDs() (ids []int) {
	for id := range m.equipment {
		ids = append(ids, id)
	}
	return
}

// ResetEquipment reset all changes of the equipment edge.
func (m *LocationMutation) ResetEquipment() {
	m.equipment = nil
	m.removedequipment = nil
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (m *LocationMutation) AddPropertyIDs(ids ...int) {
	if m.properties == nil {
		m.properties = make(map[int]struct{})
	}
	for i := range ids {
		m.properties[ids[i]] = struct{}{}
	}
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (m *LocationMutation) RemovePropertyIDs(ids ...int) {
	if m.removedproperties == nil {
		m.removedproperties = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperties[ids[i]] = struct{}{}
	}
}

// RemovedProperties returns the removed ids of properties.
func (m *LocationMutation) RemovedPropertiesIDs() (ids []int) {
	for id := range m.removedproperties {
		ids = append(ids, id)
	}
	return
}

// PropertiesIDs returns the properties ids in the mutation.
func (m *LocationMutation) PropertiesIDs() (ids []int) {
	for id := range m.properties {
		ids = append(ids, id)
	}
	return
}

// ResetProperties reset all changes of the properties edge.
func (m *LocationMutation) ResetProperties() {
	m.properties = nil
	m.removedproperties = nil
}

// AddSurveyIDs adds the survey edge to Survey by ids.
func (m *LocationMutation) AddSurveyIDs(ids ...int) {
	if m.survey == nil {
		m.survey = make(map[int]struct{})
	}
	for i := range ids {
		m.survey[ids[i]] = struct{}{}
	}
}

// RemoveSurveyIDs removes the survey edge to Survey by ids.
func (m *LocationMutation) RemoveSurveyIDs(ids ...int) {
	if m.removedsurvey == nil {
		m.removedsurvey = make(map[int]struct{})
	}
	for i := range ids {
		m.removedsurvey[ids[i]] = struct{}{}
	}
}

// RemovedSurvey returns the removed ids of survey.
func (m *LocationMutation) RemovedSurveyIDs() (ids []int) {
	for id := range m.removedsurvey {
		ids = append(ids, id)
	}
	return
}

// SurveyIDs returns the survey ids in the mutation.
func (m *LocationMutation) SurveyIDs() (ids []int) {
	for id := range m.survey {
		ids = append(ids, id)
	}
	return
}

// ResetSurvey reset all changes of the survey edge.
func (m *LocationMutation) ResetSurvey() {
	m.survey = nil
	m.removedsurvey = nil
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (m *LocationMutation) AddWifiScanIDs(ids ...int) {
	if m.wifi_scan == nil {
		m.wifi_scan = make(map[int]struct{})
	}
	for i := range ids {
		m.wifi_scan[ids[i]] = struct{}{}
	}
}

// RemoveWifiScanIDs removes the wifi_scan edge to SurveyWiFiScan by ids.
func (m *LocationMutation) RemoveWifiScanIDs(ids ...int) {
	if m.removedwifi_scan == nil {
		m.removedwifi_scan = make(map[int]struct{})
	}
	for i := range ids {
		m.removedwifi_scan[ids[i]] = struct{}{}
	}
}

// RemovedWifiScan returns the removed ids of wifi_scan.
func (m *LocationMutation) RemovedWifiScanIDs() (ids []int) {
	for id := range m.removedwifi_scan {
		ids = append(ids, id)
	}
	return
}

// WifiScanIDs returns the wifi_scan ids in the mutation.
func (m *LocationMutation) WifiScanIDs() (ids []int) {
	for id := range m.wifi_scan {
		ids = append(ids, id)
	}
	return
}

// ResetWifiScan reset all changes of the wifi_scan edge.
func (m *LocationMutation) ResetWifiScan() {
	m.wifi_scan = nil
	m.removedwifi_scan = nil
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (m *LocationMutation) AddCellScanIDs(ids ...int) {
	if m.cell_scan == nil {
		m.cell_scan = make(map[int]struct{})
	}
	for i := range ids {
		m.cell_scan[ids[i]] = struct{}{}
	}
}

// RemoveCellScanIDs removes the cell_scan edge to SurveyCellScan by ids.
func (m *LocationMutation) RemoveCellScanIDs(ids ...int) {
	if m.removedcell_scan == nil {
		m.removedcell_scan = make(map[int]struct{})
	}
	for i := range ids {
		m.removedcell_scan[ids[i]] = struct{}{}
	}
}

// RemovedCellScan returns the removed ids of cell_scan.
func (m *LocationMutation) RemovedCellScanIDs() (ids []int) {
	for id := range m.removedcell_scan {
		ids = append(ids, id)
	}
	return
}

// CellScanIDs returns the cell_scan ids in the mutation.
func (m *LocationMutation) CellScanIDs() (ids []int) {
	for id := range m.cell_scan {
		ids = append(ids, id)
	}
	return
}

// ResetCellScan reset all changes of the cell_scan edge.
func (m *LocationMutation) ResetCellScan() {
	m.cell_scan = nil
	m.removedcell_scan = nil
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (m *LocationMutation) AddWorkOrderIDs(ids ...int) {
	if m.work_orders == nil {
		m.work_orders = make(map[int]struct{})
	}
	for i := range ids {
		m.work_orders[ids[i]] = struct{}{}
	}
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (m *LocationMutation) RemoveWorkOrderIDs(ids ...int) {
	if m.removedwork_orders == nil {
		m.removedwork_orders = make(map[int]struct{})
	}
	for i := range ids {
		m.removedwork_orders[ids[i]] = struct{}{}
	}
}

// RemovedWorkOrders returns the removed ids of work_orders.
func (m *LocationMutation) RemovedWorkOrdersIDs() (ids []int) {
	for id := range m.removedwork_orders {
		ids = append(ids, id)
	}
	return
}

// WorkOrdersIDs returns the work_orders ids in the mutation.
func (m *LocationMutation) WorkOrdersIDs() (ids []int) {
	for id := range m.work_orders {
		ids = append(ids, id)
	}
	return
}

// ResetWorkOrders reset all changes of the work_orders edge.
func (m *LocationMutation) ResetWorkOrders() {
	m.work_orders = nil
	m.removedwork_orders = nil
}

// AddFloorPlanIDs adds the floor_plans edge to FloorPlan by ids.
func (m *LocationMutation) AddFloorPlanIDs(ids ...int) {
	if m.floor_plans == nil {
		m.floor_plans = make(map[int]struct{})
	}
	for i := range ids {
		m.floor_plans[ids[i]] = struct{}{}
	}
}

// RemoveFloorPlanIDs removes the floor_plans edge to FloorPlan by ids.
func (m *LocationMutation) RemoveFloorPlanIDs(ids ...int) {
	if m.removedfloor_plans == nil {
		m.removedfloor_plans = make(map[int]struct{})
	}
	for i := range ids {
		m.removedfloor_plans[ids[i]] = struct{}{}
	}
}

// RemovedFloorPlans returns the removed ids of floor_plans.
func (m *LocationMutation) RemovedFloorPlansIDs() (ids []int) {
	for id := range m.removedfloor_plans {
		ids = append(ids, id)
	}
	return
}

// FloorPlansIDs returns the floor_plans ids in the mutation.
func (m *LocationMutation) FloorPlansIDs() (ids []int) {
	for id := range m.floor_plans {
		ids = append(ids, id)
	}
	return
}

// ResetFloorPlans reset all changes of the floor_plans edge.
func (m *LocationMutation) ResetFloorPlans() {
	m.floor_plans = nil
	m.removedfloor_plans = nil
}

// Op returns the operation name.
func (m *LocationMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Location).
func (m *LocationMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *LocationMutation) Fields() []string {
	fields := make([]string, 0, 7)
	if m.create_time != nil {
		fields = append(fields, location.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, location.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, location.FieldName)
	}
	if m.external_id != nil {
		fields = append(fields, location.FieldExternalID)
	}
	if m.latitude != nil {
		fields = append(fields, location.FieldLatitude)
	}
	if m.longitude != nil {
		fields = append(fields, location.FieldLongitude)
	}
	if m.site_survey_needed != nil {
		fields = append(fields, location.FieldSiteSurveyNeeded)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *LocationMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case location.FieldCreateTime:
		return m.CreateTime()
	case location.FieldUpdateTime:
		return m.UpdateTime()
	case location.FieldName:
		return m.Name()
	case location.FieldExternalID:
		return m.ExternalID()
	case location.FieldLatitude:
		return m.Latitude()
	case location.FieldLongitude:
		return m.Longitude()
	case location.FieldSiteSurveyNeeded:
		return m.SiteSurveyNeeded()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *LocationMutation) SetField(name string, value ent.Value) error {
	switch name {
	case location.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case location.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case location.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case location.FieldExternalID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetExternalID(v)
		return nil
	case location.FieldLatitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLatitude(v)
		return nil
	case location.FieldLongitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLongitude(v)
		return nil
	case location.FieldSiteSurveyNeeded:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSiteSurveyNeeded(v)
		return nil
	}
	return fmt.Errorf("unknown Location field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *LocationMutation) AddedFields() []string {
	var fields []string
	if m.addlatitude != nil {
		fields = append(fields, location.FieldLatitude)
	}
	if m.addlongitude != nil {
		fields = append(fields, location.FieldLongitude)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *LocationMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case location.FieldLatitude:
		return m.AddedLatitude()
	case location.FieldLongitude:
		return m.AddedLongitude()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *LocationMutation) AddField(name string, value ent.Value) error {
	switch name {
	case location.FieldLatitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLatitude(v)
		return nil
	case location.FieldLongitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLongitude(v)
		return nil
	}
	return fmt.Errorf("unknown Location numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *LocationMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[location.FieldExternalID] {
		fields = append(fields, location.FieldExternalID)
	}
	if m.clearedFields[location.FieldSiteSurveyNeeded] {
		fields = append(fields, location.FieldSiteSurveyNeeded)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *LocationMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *LocationMutation) ClearField(name string) error {
	switch name {
	case location.FieldExternalID:
		m.ClearExternalID()
		return nil
	case location.FieldSiteSurveyNeeded:
		m.ClearSiteSurveyNeeded()
		return nil
	}
	return fmt.Errorf("unknown Location nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *LocationMutation) ResetField(name string) error {
	switch name {
	case location.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case location.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case location.FieldName:
		m.ResetName()
		return nil
	case location.FieldExternalID:
		m.ResetExternalID()
		return nil
	case location.FieldLatitude:
		m.ResetLatitude()
		return nil
	case location.FieldLongitude:
		m.ResetLongitude()
		return nil
	case location.FieldSiteSurveyNeeded:
		m.ResetSiteSurveyNeeded()
		return nil
	}
	return fmt.Errorf("unknown Location field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *LocationMutation) AddedEdges() []string {
	edges := make([]string, 0, 12)
	if m._type != nil {
		edges = append(edges, location.EdgeType)
	}
	if m.parent != nil {
		edges = append(edges, location.EdgeParent)
	}
	if m.children != nil {
		edges = append(edges, location.EdgeChildren)
	}
	if m.files != nil {
		edges = append(edges, location.EdgeFiles)
	}
	if m.hyperlinks != nil {
		edges = append(edges, location.EdgeHyperlinks)
	}
	if m.equipment != nil {
		edges = append(edges, location.EdgeEquipment)
	}
	if m.properties != nil {
		edges = append(edges, location.EdgeProperties)
	}
	if m.survey != nil {
		edges = append(edges, location.EdgeSurvey)
	}
	if m.wifi_scan != nil {
		edges = append(edges, location.EdgeWifiScan)
	}
	if m.cell_scan != nil {
		edges = append(edges, location.EdgeCellScan)
	}
	if m.work_orders != nil {
		edges = append(edges, location.EdgeWorkOrders)
	}
	if m.floor_plans != nil {
		edges = append(edges, location.EdgeFloorPlans)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *LocationMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case location.EdgeType:
		if id := m._type; id != nil {
			return []ent.Value{*id}
		}
	case location.EdgeParent:
		if id := m.parent; id != nil {
			return []ent.Value{*id}
		}
	case location.EdgeChildren:
		ids := make([]ent.Value, 0, len(m.children))
		for id := range m.children {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeFiles:
		ids := make([]ent.Value, 0, len(m.files))
		for id := range m.files {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeHyperlinks:
		ids := make([]ent.Value, 0, len(m.hyperlinks))
		for id := range m.hyperlinks {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeEquipment:
		ids := make([]ent.Value, 0, len(m.equipment))
		for id := range m.equipment {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.properties))
		for id := range m.properties {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeSurvey:
		ids := make([]ent.Value, 0, len(m.survey))
		for id := range m.survey {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeWifiScan:
		ids := make([]ent.Value, 0, len(m.wifi_scan))
		for id := range m.wifi_scan {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeCellScan:
		ids := make([]ent.Value, 0, len(m.cell_scan))
		for id := range m.cell_scan {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeWorkOrders:
		ids := make([]ent.Value, 0, len(m.work_orders))
		for id := range m.work_orders {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeFloorPlans:
		ids := make([]ent.Value, 0, len(m.floor_plans))
		for id := range m.floor_plans {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *LocationMutation) RemovedEdges() []string {
	edges := make([]string, 0, 12)
	if m.removedchildren != nil {
		edges = append(edges, location.EdgeChildren)
	}
	if m.removedfiles != nil {
		edges = append(edges, location.EdgeFiles)
	}
	if m.removedhyperlinks != nil {
		edges = append(edges, location.EdgeHyperlinks)
	}
	if m.removedequipment != nil {
		edges = append(edges, location.EdgeEquipment)
	}
	if m.removedproperties != nil {
		edges = append(edges, location.EdgeProperties)
	}
	if m.removedsurvey != nil {
		edges = append(edges, location.EdgeSurvey)
	}
	if m.removedwifi_scan != nil {
		edges = append(edges, location.EdgeWifiScan)
	}
	if m.removedcell_scan != nil {
		edges = append(edges, location.EdgeCellScan)
	}
	if m.removedwork_orders != nil {
		edges = append(edges, location.EdgeWorkOrders)
	}
	if m.removedfloor_plans != nil {
		edges = append(edges, location.EdgeFloorPlans)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *LocationMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case location.EdgeChildren:
		ids := make([]ent.Value, 0, len(m.removedchildren))
		for id := range m.removedchildren {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeFiles:
		ids := make([]ent.Value, 0, len(m.removedfiles))
		for id := range m.removedfiles {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeHyperlinks:
		ids := make([]ent.Value, 0, len(m.removedhyperlinks))
		for id := range m.removedhyperlinks {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeEquipment:
		ids := make([]ent.Value, 0, len(m.removedequipment))
		for id := range m.removedequipment {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.removedproperties))
		for id := range m.removedproperties {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeSurvey:
		ids := make([]ent.Value, 0, len(m.removedsurvey))
		for id := range m.removedsurvey {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeWifiScan:
		ids := make([]ent.Value, 0, len(m.removedwifi_scan))
		for id := range m.removedwifi_scan {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeCellScan:
		ids := make([]ent.Value, 0, len(m.removedcell_scan))
		for id := range m.removedcell_scan {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeWorkOrders:
		ids := make([]ent.Value, 0, len(m.removedwork_orders))
		for id := range m.removedwork_orders {
			ids = append(ids, id)
		}
		return ids
	case location.EdgeFloorPlans:
		ids := make([]ent.Value, 0, len(m.removedfloor_plans))
		for id := range m.removedfloor_plans {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *LocationMutation) ClearedEdges() []string {
	edges := make([]string, 0, 12)
	if m.cleared_type {
		edges = append(edges, location.EdgeType)
	}
	if m.clearedparent {
		edges = append(edges, location.EdgeParent)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *LocationMutation) EdgeCleared(name string) bool {
	switch name {
	case location.EdgeType:
		return m.cleared_type
	case location.EdgeParent:
		return m.clearedparent
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *LocationMutation) ClearEdge(name string) error {
	switch name {
	case location.EdgeType:
		m.ClearType()
		return nil
	case location.EdgeParent:
		m.ClearParent()
		return nil
	}
	return fmt.Errorf("unknown Location unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *LocationMutation) ResetEdge(name string) error {
	switch name {
	case location.EdgeType:
		m.ResetType()
		return nil
	case location.EdgeParent:
		m.ResetParent()
		return nil
	case location.EdgeChildren:
		m.ResetChildren()
		return nil
	case location.EdgeFiles:
		m.ResetFiles()
		return nil
	case location.EdgeHyperlinks:
		m.ResetHyperlinks()
		return nil
	case location.EdgeEquipment:
		m.ResetEquipment()
		return nil
	case location.EdgeProperties:
		m.ResetProperties()
		return nil
	case location.EdgeSurvey:
		m.ResetSurvey()
		return nil
	case location.EdgeWifiScan:
		m.ResetWifiScan()
		return nil
	case location.EdgeCellScan:
		m.ResetCellScan()
		return nil
	case location.EdgeWorkOrders:
		m.ResetWorkOrders()
		return nil
	case location.EdgeFloorPlans:
		m.ResetFloorPlans()
		return nil
	}
	return fmt.Errorf("unknown Location edge %s", name)
}

// LocationTypeMutation represents an operation that mutate the LocationTypes
// nodes in the graph.
type LocationTypeMutation struct {
	config
	op                                Op
	typ                               string
	id                                *int
	create_time                       *time.Time
	update_time                       *time.Time
	site                              *bool
	name                              *string
	map_type                          *string
	map_zoom_level                    *int
	addmap_zoom_level                 *int
	index                             *int
	addindex                          *int
	clearedFields                     map[string]bool
	locations                         map[int]struct{}
	removedlocations                  map[int]struct{}
	property_types                    map[int]struct{}
	removedproperty_types             map[int]struct{}
	survey_template_categories        map[int]struct{}
	removedsurvey_template_categories map[int]struct{}
}

var _ ent.Mutation = (*LocationTypeMutation)(nil)

// newLocationTypeMutation creates new mutation for $n.Name.
func newLocationTypeMutation(c config, op Op) *LocationTypeMutation {
	return &LocationTypeMutation{
		config:        c,
		op:            op,
		typ:           TypeLocationType,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m LocationTypeMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m LocationTypeMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *LocationTypeMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *LocationTypeMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *LocationTypeMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *LocationTypeMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *LocationTypeMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *LocationTypeMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *LocationTypeMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetSite sets the site field.
func (m *LocationTypeMutation) SetSite(b bool) {
	m.site = &b
}

// Site returns the site value in the mutation.
func (m *LocationTypeMutation) Site() (r bool, exists bool) {
	v := m.site
	if v == nil {
		return
	}
	return *v, true
}

// ResetSite reset all changes of the site field.
func (m *LocationTypeMutation) ResetSite() {
	m.site = nil
}

// SetName sets the name field.
func (m *LocationTypeMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *LocationTypeMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *LocationTypeMutation) ResetName() {
	m.name = nil
}

// SetMapType sets the map_type field.
func (m *LocationTypeMutation) SetMapType(s string) {
	m.map_type = &s
}

// MapType returns the map_type value in the mutation.
func (m *LocationTypeMutation) MapType() (r string, exists bool) {
	v := m.map_type
	if v == nil {
		return
	}
	return *v, true
}

// ClearMapType clears the value of map_type.
func (m *LocationTypeMutation) ClearMapType() {
	m.map_type = nil
	m.clearedFields[locationtype.FieldMapType] = true
}

// MapTypeCleared returns if the field map_type was cleared in this mutation.
func (m *LocationTypeMutation) MapTypeCleared() bool {
	return m.clearedFields[locationtype.FieldMapType]
}

// ResetMapType reset all changes of the map_type field.
func (m *LocationTypeMutation) ResetMapType() {
	m.map_type = nil
	delete(m.clearedFields, locationtype.FieldMapType)
}

// SetMapZoomLevel sets the map_zoom_level field.
func (m *LocationTypeMutation) SetMapZoomLevel(i int) {
	m.map_zoom_level = &i
	m.addmap_zoom_level = nil
}

// MapZoomLevel returns the map_zoom_level value in the mutation.
func (m *LocationTypeMutation) MapZoomLevel() (r int, exists bool) {
	v := m.map_zoom_level
	if v == nil {
		return
	}
	return *v, true
}

// AddMapZoomLevel adds i to map_zoom_level.
func (m *LocationTypeMutation) AddMapZoomLevel(i int) {
	if m.addmap_zoom_level != nil {
		*m.addmap_zoom_level += i
	} else {
		m.addmap_zoom_level = &i
	}
}

// AddedMapZoomLevel returns the value that was added to the map_zoom_level field in this mutation.
func (m *LocationTypeMutation) AddedMapZoomLevel() (r int, exists bool) {
	v := m.addmap_zoom_level
	if v == nil {
		return
	}
	return *v, true
}

// ClearMapZoomLevel clears the value of map_zoom_level.
func (m *LocationTypeMutation) ClearMapZoomLevel() {
	m.map_zoom_level = nil
	m.addmap_zoom_level = nil
	m.clearedFields[locationtype.FieldMapZoomLevel] = true
}

// MapZoomLevelCleared returns if the field map_zoom_level was cleared in this mutation.
func (m *LocationTypeMutation) MapZoomLevelCleared() bool {
	return m.clearedFields[locationtype.FieldMapZoomLevel]
}

// ResetMapZoomLevel reset all changes of the map_zoom_level field.
func (m *LocationTypeMutation) ResetMapZoomLevel() {
	m.map_zoom_level = nil
	m.addmap_zoom_level = nil
	delete(m.clearedFields, locationtype.FieldMapZoomLevel)
}

// SetIndex sets the index field.
func (m *LocationTypeMutation) SetIndex(i int) {
	m.index = &i
	m.addindex = nil
}

// Index returns the index value in the mutation.
func (m *LocationTypeMutation) Index() (r int, exists bool) {
	v := m.index
	if v == nil {
		return
	}
	return *v, true
}

// AddIndex adds i to index.
func (m *LocationTypeMutation) AddIndex(i int) {
	if m.addindex != nil {
		*m.addindex += i
	} else {
		m.addindex = &i
	}
}

// AddedIndex returns the value that was added to the index field in this mutation.
func (m *LocationTypeMutation) AddedIndex() (r int, exists bool) {
	v := m.addindex
	if v == nil {
		return
	}
	return *v, true
}

// ResetIndex reset all changes of the index field.
func (m *LocationTypeMutation) ResetIndex() {
	m.index = nil
	m.addindex = nil
}

// AddLocationIDs adds the locations edge to Location by ids.
func (m *LocationTypeMutation) AddLocationIDs(ids ...int) {
	if m.locations == nil {
		m.locations = make(map[int]struct{})
	}
	for i := range ids {
		m.locations[ids[i]] = struct{}{}
	}
}

// RemoveLocationIDs removes the locations edge to Location by ids.
func (m *LocationTypeMutation) RemoveLocationIDs(ids ...int) {
	if m.removedlocations == nil {
		m.removedlocations = make(map[int]struct{})
	}
	for i := range ids {
		m.removedlocations[ids[i]] = struct{}{}
	}
}

// RemovedLocations returns the removed ids of locations.
func (m *LocationTypeMutation) RemovedLocationsIDs() (ids []int) {
	for id := range m.removedlocations {
		ids = append(ids, id)
	}
	return
}

// LocationsIDs returns the locations ids in the mutation.
func (m *LocationTypeMutation) LocationsIDs() (ids []int) {
	for id := range m.locations {
		ids = append(ids, id)
	}
	return
}

// ResetLocations reset all changes of the locations edge.
func (m *LocationTypeMutation) ResetLocations() {
	m.locations = nil
	m.removedlocations = nil
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (m *LocationTypeMutation) AddPropertyTypeIDs(ids ...int) {
	if m.property_types == nil {
		m.property_types = make(map[int]struct{})
	}
	for i := range ids {
		m.property_types[ids[i]] = struct{}{}
	}
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (m *LocationTypeMutation) RemovePropertyTypeIDs(ids ...int) {
	if m.removedproperty_types == nil {
		m.removedproperty_types = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperty_types[ids[i]] = struct{}{}
	}
}

// RemovedPropertyTypes returns the removed ids of property_types.
func (m *LocationTypeMutation) RemovedPropertyTypesIDs() (ids []int) {
	for id := range m.removedproperty_types {
		ids = append(ids, id)
	}
	return
}

// PropertyTypesIDs returns the property_types ids in the mutation.
func (m *LocationTypeMutation) PropertyTypesIDs() (ids []int) {
	for id := range m.property_types {
		ids = append(ids, id)
	}
	return
}

// ResetPropertyTypes reset all changes of the property_types edge.
func (m *LocationTypeMutation) ResetPropertyTypes() {
	m.property_types = nil
	m.removedproperty_types = nil
}

// AddSurveyTemplateCategoryIDs adds the survey_template_categories edge to SurveyTemplateCategory by ids.
func (m *LocationTypeMutation) AddSurveyTemplateCategoryIDs(ids ...int) {
	if m.survey_template_categories == nil {
		m.survey_template_categories = make(map[int]struct{})
	}
	for i := range ids {
		m.survey_template_categories[ids[i]] = struct{}{}
	}
}

// RemoveSurveyTemplateCategoryIDs removes the survey_template_categories edge to SurveyTemplateCategory by ids.
func (m *LocationTypeMutation) RemoveSurveyTemplateCategoryIDs(ids ...int) {
	if m.removedsurvey_template_categories == nil {
		m.removedsurvey_template_categories = make(map[int]struct{})
	}
	for i := range ids {
		m.removedsurvey_template_categories[ids[i]] = struct{}{}
	}
}

// RemovedSurveyTemplateCategories returns the removed ids of survey_template_categories.
func (m *LocationTypeMutation) RemovedSurveyTemplateCategoriesIDs() (ids []int) {
	for id := range m.removedsurvey_template_categories {
		ids = append(ids, id)
	}
	return
}

// SurveyTemplateCategoriesIDs returns the survey_template_categories ids in the mutation.
func (m *LocationTypeMutation) SurveyTemplateCategoriesIDs() (ids []int) {
	for id := range m.survey_template_categories {
		ids = append(ids, id)
	}
	return
}

// ResetSurveyTemplateCategories reset all changes of the survey_template_categories edge.
func (m *LocationTypeMutation) ResetSurveyTemplateCategories() {
	m.survey_template_categories = nil
	m.removedsurvey_template_categories = nil
}

// Op returns the operation name.
func (m *LocationTypeMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (LocationType).
func (m *LocationTypeMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *LocationTypeMutation) Fields() []string {
	fields := make([]string, 0, 7)
	if m.create_time != nil {
		fields = append(fields, locationtype.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, locationtype.FieldUpdateTime)
	}
	if m.site != nil {
		fields = append(fields, locationtype.FieldSite)
	}
	if m.name != nil {
		fields = append(fields, locationtype.FieldName)
	}
	if m.map_type != nil {
		fields = append(fields, locationtype.FieldMapType)
	}
	if m.map_zoom_level != nil {
		fields = append(fields, locationtype.FieldMapZoomLevel)
	}
	if m.index != nil {
		fields = append(fields, locationtype.FieldIndex)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *LocationTypeMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case locationtype.FieldCreateTime:
		return m.CreateTime()
	case locationtype.FieldUpdateTime:
		return m.UpdateTime()
	case locationtype.FieldSite:
		return m.Site()
	case locationtype.FieldName:
		return m.Name()
	case locationtype.FieldMapType:
		return m.MapType()
	case locationtype.FieldMapZoomLevel:
		return m.MapZoomLevel()
	case locationtype.FieldIndex:
		return m.Index()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *LocationTypeMutation) SetField(name string, value ent.Value) error {
	switch name {
	case locationtype.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case locationtype.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case locationtype.FieldSite:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSite(v)
		return nil
	case locationtype.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case locationtype.FieldMapType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMapType(v)
		return nil
	case locationtype.FieldMapZoomLevel:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMapZoomLevel(v)
		return nil
	case locationtype.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIndex(v)
		return nil
	}
	return fmt.Errorf("unknown LocationType field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *LocationTypeMutation) AddedFields() []string {
	var fields []string
	if m.addmap_zoom_level != nil {
		fields = append(fields, locationtype.FieldMapZoomLevel)
	}
	if m.addindex != nil {
		fields = append(fields, locationtype.FieldIndex)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *LocationTypeMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case locationtype.FieldMapZoomLevel:
		return m.AddedMapZoomLevel()
	case locationtype.FieldIndex:
		return m.AddedIndex()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *LocationTypeMutation) AddField(name string, value ent.Value) error {
	switch name {
	case locationtype.FieldMapZoomLevel:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddMapZoomLevel(v)
		return nil
	case locationtype.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIndex(v)
		return nil
	}
	return fmt.Errorf("unknown LocationType numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *LocationTypeMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[locationtype.FieldMapType] {
		fields = append(fields, locationtype.FieldMapType)
	}
	if m.clearedFields[locationtype.FieldMapZoomLevel] {
		fields = append(fields, locationtype.FieldMapZoomLevel)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *LocationTypeMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *LocationTypeMutation) ClearField(name string) error {
	switch name {
	case locationtype.FieldMapType:
		m.ClearMapType()
		return nil
	case locationtype.FieldMapZoomLevel:
		m.ClearMapZoomLevel()
		return nil
	}
	return fmt.Errorf("unknown LocationType nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *LocationTypeMutation) ResetField(name string) error {
	switch name {
	case locationtype.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case locationtype.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case locationtype.FieldSite:
		m.ResetSite()
		return nil
	case locationtype.FieldName:
		m.ResetName()
		return nil
	case locationtype.FieldMapType:
		m.ResetMapType()
		return nil
	case locationtype.FieldMapZoomLevel:
		m.ResetMapZoomLevel()
		return nil
	case locationtype.FieldIndex:
		m.ResetIndex()
		return nil
	}
	return fmt.Errorf("unknown LocationType field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *LocationTypeMutation) AddedEdges() []string {
	edges := make([]string, 0, 3)
	if m.locations != nil {
		edges = append(edges, locationtype.EdgeLocations)
	}
	if m.property_types != nil {
		edges = append(edges, locationtype.EdgePropertyTypes)
	}
	if m.survey_template_categories != nil {
		edges = append(edges, locationtype.EdgeSurveyTemplateCategories)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *LocationTypeMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case locationtype.EdgeLocations:
		ids := make([]ent.Value, 0, len(m.locations))
		for id := range m.locations {
			ids = append(ids, id)
		}
		return ids
	case locationtype.EdgePropertyTypes:
		ids := make([]ent.Value, 0, len(m.property_types))
		for id := range m.property_types {
			ids = append(ids, id)
		}
		return ids
	case locationtype.EdgeSurveyTemplateCategories:
		ids := make([]ent.Value, 0, len(m.survey_template_categories))
		for id := range m.survey_template_categories {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *LocationTypeMutation) RemovedEdges() []string {
	edges := make([]string, 0, 3)
	if m.removedlocations != nil {
		edges = append(edges, locationtype.EdgeLocations)
	}
	if m.removedproperty_types != nil {
		edges = append(edges, locationtype.EdgePropertyTypes)
	}
	if m.removedsurvey_template_categories != nil {
		edges = append(edges, locationtype.EdgeSurveyTemplateCategories)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *LocationTypeMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case locationtype.EdgeLocations:
		ids := make([]ent.Value, 0, len(m.removedlocations))
		for id := range m.removedlocations {
			ids = append(ids, id)
		}
		return ids
	case locationtype.EdgePropertyTypes:
		ids := make([]ent.Value, 0, len(m.removedproperty_types))
		for id := range m.removedproperty_types {
			ids = append(ids, id)
		}
		return ids
	case locationtype.EdgeSurveyTemplateCategories:
		ids := make([]ent.Value, 0, len(m.removedsurvey_template_categories))
		for id := range m.removedsurvey_template_categories {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *LocationTypeMutation) ClearedEdges() []string {
	edges := make([]string, 0, 3)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *LocationTypeMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *LocationTypeMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown LocationType unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *LocationTypeMutation) ResetEdge(name string) error {
	switch name {
	case locationtype.EdgeLocations:
		m.ResetLocations()
		return nil
	case locationtype.EdgePropertyTypes:
		m.ResetPropertyTypes()
		return nil
	case locationtype.EdgeSurveyTemplateCategories:
		m.ResetSurveyTemplateCategories()
		return nil
	}
	return fmt.Errorf("unknown LocationType edge %s", name)
}

// ProjectMutation represents an operation that mutate the Projects
// nodes in the graph.
type ProjectMutation struct {
	config
	op                 Op
	typ                string
	id                 *int
	create_time        *time.Time
	update_time        *time.Time
	name               *string
	description        *string
	clearedFields      map[string]bool
	_type              *int
	cleared_type       bool
	location           *int
	clearedlocation    bool
	comments           map[int]struct{}
	removedcomments    map[int]struct{}
	work_orders        map[int]struct{}
	removedwork_orders map[int]struct{}
	properties         map[int]struct{}
	removedproperties  map[int]struct{}
	creator            *int
	clearedcreator     bool
}

var _ ent.Mutation = (*ProjectMutation)(nil)

// newProjectMutation creates new mutation for $n.Name.
func newProjectMutation(c config, op Op) *ProjectMutation {
	return &ProjectMutation{
		config:        c,
		op:            op,
		typ:           TypeProject,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ProjectMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ProjectMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *ProjectMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *ProjectMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *ProjectMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *ProjectMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *ProjectMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *ProjectMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *ProjectMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *ProjectMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *ProjectMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *ProjectMutation) ResetName() {
	m.name = nil
}

// SetDescription sets the description field.
func (m *ProjectMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the description value in the mutation.
func (m *ProjectMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// ClearDescription clears the value of description.
func (m *ProjectMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[project.FieldDescription] = true
}

// DescriptionCleared returns if the field description was cleared in this mutation.
func (m *ProjectMutation) DescriptionCleared() bool {
	return m.clearedFields[project.FieldDescription]
}

// ResetDescription reset all changes of the description field.
func (m *ProjectMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, project.FieldDescription)
}

// SetTypeID sets the type edge to ProjectType by id.
func (m *ProjectMutation) SetTypeID(id int) {
	m._type = &id
}

// ClearType clears the type edge to ProjectType.
func (m *ProjectMutation) ClearType() {
	m.cleared_type = true
}

// TypeCleared returns if the edge type was cleared.
func (m *ProjectMutation) TypeCleared() bool {
	return m.cleared_type
}

// TypeID returns the type id in the mutation.
func (m *ProjectMutation) TypeID() (id int, exists bool) {
	if m._type != nil {
		return *m._type, true
	}
	return
}

// TypeIDs returns the type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// TypeID instead. It exists only for internal usage by the builders.
func (m *ProjectMutation) TypeIDs() (ids []int) {
	if id := m._type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetType reset all changes of the type edge.
func (m *ProjectMutation) ResetType() {
	m._type = nil
	m.cleared_type = false
}

// SetLocationID sets the location edge to Location by id.
func (m *ProjectMutation) SetLocationID(id int) {
	m.location = &id
}

// ClearLocation clears the location edge to Location.
func (m *ProjectMutation) ClearLocation() {
	m.clearedlocation = true
}

// LocationCleared returns if the edge location was cleared.
func (m *ProjectMutation) LocationCleared() bool {
	return m.clearedlocation
}

// LocationID returns the location id in the mutation.
func (m *ProjectMutation) LocationID() (id int, exists bool) {
	if m.location != nil {
		return *m.location, true
	}
	return
}

// LocationIDs returns the location ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LocationID instead. It exists only for internal usage by the builders.
func (m *ProjectMutation) LocationIDs() (ids []int) {
	if id := m.location; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLocation reset all changes of the location edge.
func (m *ProjectMutation) ResetLocation() {
	m.location = nil
	m.clearedlocation = false
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (m *ProjectMutation) AddCommentIDs(ids ...int) {
	if m.comments == nil {
		m.comments = make(map[int]struct{})
	}
	for i := range ids {
		m.comments[ids[i]] = struct{}{}
	}
}

// RemoveCommentIDs removes the comments edge to Comment by ids.
func (m *ProjectMutation) RemoveCommentIDs(ids ...int) {
	if m.removedcomments == nil {
		m.removedcomments = make(map[int]struct{})
	}
	for i := range ids {
		m.removedcomments[ids[i]] = struct{}{}
	}
}

// RemovedComments returns the removed ids of comments.
func (m *ProjectMutation) RemovedCommentsIDs() (ids []int) {
	for id := range m.removedcomments {
		ids = append(ids, id)
	}
	return
}

// CommentsIDs returns the comments ids in the mutation.
func (m *ProjectMutation) CommentsIDs() (ids []int) {
	for id := range m.comments {
		ids = append(ids, id)
	}
	return
}

// ResetComments reset all changes of the comments edge.
func (m *ProjectMutation) ResetComments() {
	m.comments = nil
	m.removedcomments = nil
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (m *ProjectMutation) AddWorkOrderIDs(ids ...int) {
	if m.work_orders == nil {
		m.work_orders = make(map[int]struct{})
	}
	for i := range ids {
		m.work_orders[ids[i]] = struct{}{}
	}
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (m *ProjectMutation) RemoveWorkOrderIDs(ids ...int) {
	if m.removedwork_orders == nil {
		m.removedwork_orders = make(map[int]struct{})
	}
	for i := range ids {
		m.removedwork_orders[ids[i]] = struct{}{}
	}
}

// RemovedWorkOrders returns the removed ids of work_orders.
func (m *ProjectMutation) RemovedWorkOrdersIDs() (ids []int) {
	for id := range m.removedwork_orders {
		ids = append(ids, id)
	}
	return
}

// WorkOrdersIDs returns the work_orders ids in the mutation.
func (m *ProjectMutation) WorkOrdersIDs() (ids []int) {
	for id := range m.work_orders {
		ids = append(ids, id)
	}
	return
}

// ResetWorkOrders reset all changes of the work_orders edge.
func (m *ProjectMutation) ResetWorkOrders() {
	m.work_orders = nil
	m.removedwork_orders = nil
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (m *ProjectMutation) AddPropertyIDs(ids ...int) {
	if m.properties == nil {
		m.properties = make(map[int]struct{})
	}
	for i := range ids {
		m.properties[ids[i]] = struct{}{}
	}
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (m *ProjectMutation) RemovePropertyIDs(ids ...int) {
	if m.removedproperties == nil {
		m.removedproperties = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperties[ids[i]] = struct{}{}
	}
}

// RemovedProperties returns the removed ids of properties.
func (m *ProjectMutation) RemovedPropertiesIDs() (ids []int) {
	for id := range m.removedproperties {
		ids = append(ids, id)
	}
	return
}

// PropertiesIDs returns the properties ids in the mutation.
func (m *ProjectMutation) PropertiesIDs() (ids []int) {
	for id := range m.properties {
		ids = append(ids, id)
	}
	return
}

// ResetProperties reset all changes of the properties edge.
func (m *ProjectMutation) ResetProperties() {
	m.properties = nil
	m.removedproperties = nil
}

// SetCreatorID sets the creator edge to User by id.
func (m *ProjectMutation) SetCreatorID(id int) {
	m.creator = &id
}

// ClearCreator clears the creator edge to User.
func (m *ProjectMutation) ClearCreator() {
	m.clearedcreator = true
}

// CreatorCleared returns if the edge creator was cleared.
func (m *ProjectMutation) CreatorCleared() bool {
	return m.clearedcreator
}

// CreatorID returns the creator id in the mutation.
func (m *ProjectMutation) CreatorID() (id int, exists bool) {
	if m.creator != nil {
		return *m.creator, true
	}
	return
}

// CreatorIDs returns the creator ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// CreatorID instead. It exists only for internal usage by the builders.
func (m *ProjectMutation) CreatorIDs() (ids []int) {
	if id := m.creator; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetCreator reset all changes of the creator edge.
func (m *ProjectMutation) ResetCreator() {
	m.creator = nil
	m.clearedcreator = false
}

// Op returns the operation name.
func (m *ProjectMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Project).
func (m *ProjectMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *ProjectMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.create_time != nil {
		fields = append(fields, project.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, project.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, project.FieldName)
	}
	if m.description != nil {
		fields = append(fields, project.FieldDescription)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *ProjectMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case project.FieldCreateTime:
		return m.CreateTime()
	case project.FieldUpdateTime:
		return m.UpdateTime()
	case project.FieldName:
		return m.Name()
	case project.FieldDescription:
		return m.Description()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ProjectMutation) SetField(name string, value ent.Value) error {
	switch name {
	case project.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case project.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case project.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case project.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	}
	return fmt.Errorf("unknown Project field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *ProjectMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *ProjectMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ProjectMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Project numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *ProjectMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[project.FieldDescription] {
		fields = append(fields, project.FieldDescription)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *ProjectMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *ProjectMutation) ClearField(name string) error {
	switch name {
	case project.FieldDescription:
		m.ClearDescription()
		return nil
	}
	return fmt.Errorf("unknown Project nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *ProjectMutation) ResetField(name string) error {
	switch name {
	case project.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case project.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case project.FieldName:
		m.ResetName()
		return nil
	case project.FieldDescription:
		m.ResetDescription()
		return nil
	}
	return fmt.Errorf("unknown Project field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *ProjectMutation) AddedEdges() []string {
	edges := make([]string, 0, 6)
	if m._type != nil {
		edges = append(edges, project.EdgeType)
	}
	if m.location != nil {
		edges = append(edges, project.EdgeLocation)
	}
	if m.comments != nil {
		edges = append(edges, project.EdgeComments)
	}
	if m.work_orders != nil {
		edges = append(edges, project.EdgeWorkOrders)
	}
	if m.properties != nil {
		edges = append(edges, project.EdgeProperties)
	}
	if m.creator != nil {
		edges = append(edges, project.EdgeCreator)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *ProjectMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case project.EdgeType:
		if id := m._type; id != nil {
			return []ent.Value{*id}
		}
	case project.EdgeLocation:
		if id := m.location; id != nil {
			return []ent.Value{*id}
		}
	case project.EdgeComments:
		ids := make([]ent.Value, 0, len(m.comments))
		for id := range m.comments {
			ids = append(ids, id)
		}
		return ids
	case project.EdgeWorkOrders:
		ids := make([]ent.Value, 0, len(m.work_orders))
		for id := range m.work_orders {
			ids = append(ids, id)
		}
		return ids
	case project.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.properties))
		for id := range m.properties {
			ids = append(ids, id)
		}
		return ids
	case project.EdgeCreator:
		if id := m.creator; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *ProjectMutation) RemovedEdges() []string {
	edges := make([]string, 0, 6)
	if m.removedcomments != nil {
		edges = append(edges, project.EdgeComments)
	}
	if m.removedwork_orders != nil {
		edges = append(edges, project.EdgeWorkOrders)
	}
	if m.removedproperties != nil {
		edges = append(edges, project.EdgeProperties)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *ProjectMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case project.EdgeComments:
		ids := make([]ent.Value, 0, len(m.removedcomments))
		for id := range m.removedcomments {
			ids = append(ids, id)
		}
		return ids
	case project.EdgeWorkOrders:
		ids := make([]ent.Value, 0, len(m.removedwork_orders))
		for id := range m.removedwork_orders {
			ids = append(ids, id)
		}
		return ids
	case project.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.removedproperties))
		for id := range m.removedproperties {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *ProjectMutation) ClearedEdges() []string {
	edges := make([]string, 0, 6)
	if m.cleared_type {
		edges = append(edges, project.EdgeType)
	}
	if m.clearedlocation {
		edges = append(edges, project.EdgeLocation)
	}
	if m.clearedcreator {
		edges = append(edges, project.EdgeCreator)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *ProjectMutation) EdgeCleared(name string) bool {
	switch name {
	case project.EdgeType:
		return m.cleared_type
	case project.EdgeLocation:
		return m.clearedlocation
	case project.EdgeCreator:
		return m.clearedcreator
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *ProjectMutation) ClearEdge(name string) error {
	switch name {
	case project.EdgeType:
		m.ClearType()
		return nil
	case project.EdgeLocation:
		m.ClearLocation()
		return nil
	case project.EdgeCreator:
		m.ClearCreator()
		return nil
	}
	return fmt.Errorf("unknown Project unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *ProjectMutation) ResetEdge(name string) error {
	switch name {
	case project.EdgeType:
		m.ResetType()
		return nil
	case project.EdgeLocation:
		m.ResetLocation()
		return nil
	case project.EdgeComments:
		m.ResetComments()
		return nil
	case project.EdgeWorkOrders:
		m.ResetWorkOrders()
		return nil
	case project.EdgeProperties:
		m.ResetProperties()
		return nil
	case project.EdgeCreator:
		m.ResetCreator()
		return nil
	}
	return fmt.Errorf("unknown Project edge %s", name)
}

// ProjectTypeMutation represents an operation that mutate the ProjectTypes
// nodes in the graph.
type ProjectTypeMutation struct {
	config
	op                 Op
	typ                string
	id                 *int
	create_time        *time.Time
	update_time        *time.Time
	name               *string
	description        *string
	clearedFields      map[string]bool
	projects           map[int]struct{}
	removedprojects    map[int]struct{}
	properties         map[int]struct{}
	removedproperties  map[int]struct{}
	work_orders        map[int]struct{}
	removedwork_orders map[int]struct{}
}

var _ ent.Mutation = (*ProjectTypeMutation)(nil)

// newProjectTypeMutation creates new mutation for $n.Name.
func newProjectTypeMutation(c config, op Op) *ProjectTypeMutation {
	return &ProjectTypeMutation{
		config:        c,
		op:            op,
		typ:           TypeProjectType,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ProjectTypeMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ProjectTypeMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *ProjectTypeMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *ProjectTypeMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *ProjectTypeMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *ProjectTypeMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *ProjectTypeMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *ProjectTypeMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *ProjectTypeMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *ProjectTypeMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *ProjectTypeMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *ProjectTypeMutation) ResetName() {
	m.name = nil
}

// SetDescription sets the description field.
func (m *ProjectTypeMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the description value in the mutation.
func (m *ProjectTypeMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// ClearDescription clears the value of description.
func (m *ProjectTypeMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[projecttype.FieldDescription] = true
}

// DescriptionCleared returns if the field description was cleared in this mutation.
func (m *ProjectTypeMutation) DescriptionCleared() bool {
	return m.clearedFields[projecttype.FieldDescription]
}

// ResetDescription reset all changes of the description field.
func (m *ProjectTypeMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, projecttype.FieldDescription)
}

// AddProjectIDs adds the projects edge to Project by ids.
func (m *ProjectTypeMutation) AddProjectIDs(ids ...int) {
	if m.projects == nil {
		m.projects = make(map[int]struct{})
	}
	for i := range ids {
		m.projects[ids[i]] = struct{}{}
	}
}

// RemoveProjectIDs removes the projects edge to Project by ids.
func (m *ProjectTypeMutation) RemoveProjectIDs(ids ...int) {
	if m.removedprojects == nil {
		m.removedprojects = make(map[int]struct{})
	}
	for i := range ids {
		m.removedprojects[ids[i]] = struct{}{}
	}
}

// RemovedProjects returns the removed ids of projects.
func (m *ProjectTypeMutation) RemovedProjectsIDs() (ids []int) {
	for id := range m.removedprojects {
		ids = append(ids, id)
	}
	return
}

// ProjectsIDs returns the projects ids in the mutation.
func (m *ProjectTypeMutation) ProjectsIDs() (ids []int) {
	for id := range m.projects {
		ids = append(ids, id)
	}
	return
}

// ResetProjects reset all changes of the projects edge.
func (m *ProjectTypeMutation) ResetProjects() {
	m.projects = nil
	m.removedprojects = nil
}

// AddPropertyIDs adds the properties edge to PropertyType by ids.
func (m *ProjectTypeMutation) AddPropertyIDs(ids ...int) {
	if m.properties == nil {
		m.properties = make(map[int]struct{})
	}
	for i := range ids {
		m.properties[ids[i]] = struct{}{}
	}
}

// RemovePropertyIDs removes the properties edge to PropertyType by ids.
func (m *ProjectTypeMutation) RemovePropertyIDs(ids ...int) {
	if m.removedproperties == nil {
		m.removedproperties = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperties[ids[i]] = struct{}{}
	}
}

// RemovedProperties returns the removed ids of properties.
func (m *ProjectTypeMutation) RemovedPropertiesIDs() (ids []int) {
	for id := range m.removedproperties {
		ids = append(ids, id)
	}
	return
}

// PropertiesIDs returns the properties ids in the mutation.
func (m *ProjectTypeMutation) PropertiesIDs() (ids []int) {
	for id := range m.properties {
		ids = append(ids, id)
	}
	return
}

// ResetProperties reset all changes of the properties edge.
func (m *ProjectTypeMutation) ResetProperties() {
	m.properties = nil
	m.removedproperties = nil
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrderDefinition by ids.
func (m *ProjectTypeMutation) AddWorkOrderIDs(ids ...int) {
	if m.work_orders == nil {
		m.work_orders = make(map[int]struct{})
	}
	for i := range ids {
		m.work_orders[ids[i]] = struct{}{}
	}
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrderDefinition by ids.
func (m *ProjectTypeMutation) RemoveWorkOrderIDs(ids ...int) {
	if m.removedwork_orders == nil {
		m.removedwork_orders = make(map[int]struct{})
	}
	for i := range ids {
		m.removedwork_orders[ids[i]] = struct{}{}
	}
}

// RemovedWorkOrders returns the removed ids of work_orders.
func (m *ProjectTypeMutation) RemovedWorkOrdersIDs() (ids []int) {
	for id := range m.removedwork_orders {
		ids = append(ids, id)
	}
	return
}

// WorkOrdersIDs returns the work_orders ids in the mutation.
func (m *ProjectTypeMutation) WorkOrdersIDs() (ids []int) {
	for id := range m.work_orders {
		ids = append(ids, id)
	}
	return
}

// ResetWorkOrders reset all changes of the work_orders edge.
func (m *ProjectTypeMutation) ResetWorkOrders() {
	m.work_orders = nil
	m.removedwork_orders = nil
}

// Op returns the operation name.
func (m *ProjectTypeMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (ProjectType).
func (m *ProjectTypeMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *ProjectTypeMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.create_time != nil {
		fields = append(fields, projecttype.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, projecttype.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, projecttype.FieldName)
	}
	if m.description != nil {
		fields = append(fields, projecttype.FieldDescription)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *ProjectTypeMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case projecttype.FieldCreateTime:
		return m.CreateTime()
	case projecttype.FieldUpdateTime:
		return m.UpdateTime()
	case projecttype.FieldName:
		return m.Name()
	case projecttype.FieldDescription:
		return m.Description()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ProjectTypeMutation) SetField(name string, value ent.Value) error {
	switch name {
	case projecttype.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case projecttype.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case projecttype.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case projecttype.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	}
	return fmt.Errorf("unknown ProjectType field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *ProjectTypeMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *ProjectTypeMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ProjectTypeMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ProjectType numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *ProjectTypeMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[projecttype.FieldDescription] {
		fields = append(fields, projecttype.FieldDescription)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *ProjectTypeMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *ProjectTypeMutation) ClearField(name string) error {
	switch name {
	case projecttype.FieldDescription:
		m.ClearDescription()
		return nil
	}
	return fmt.Errorf("unknown ProjectType nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *ProjectTypeMutation) ResetField(name string) error {
	switch name {
	case projecttype.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case projecttype.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case projecttype.FieldName:
		m.ResetName()
		return nil
	case projecttype.FieldDescription:
		m.ResetDescription()
		return nil
	}
	return fmt.Errorf("unknown ProjectType field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *ProjectTypeMutation) AddedEdges() []string {
	edges := make([]string, 0, 3)
	if m.projects != nil {
		edges = append(edges, projecttype.EdgeProjects)
	}
	if m.properties != nil {
		edges = append(edges, projecttype.EdgeProperties)
	}
	if m.work_orders != nil {
		edges = append(edges, projecttype.EdgeWorkOrders)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *ProjectTypeMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case projecttype.EdgeProjects:
		ids := make([]ent.Value, 0, len(m.projects))
		for id := range m.projects {
			ids = append(ids, id)
		}
		return ids
	case projecttype.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.properties))
		for id := range m.properties {
			ids = append(ids, id)
		}
		return ids
	case projecttype.EdgeWorkOrders:
		ids := make([]ent.Value, 0, len(m.work_orders))
		for id := range m.work_orders {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *ProjectTypeMutation) RemovedEdges() []string {
	edges := make([]string, 0, 3)
	if m.removedprojects != nil {
		edges = append(edges, projecttype.EdgeProjects)
	}
	if m.removedproperties != nil {
		edges = append(edges, projecttype.EdgeProperties)
	}
	if m.removedwork_orders != nil {
		edges = append(edges, projecttype.EdgeWorkOrders)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *ProjectTypeMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case projecttype.EdgeProjects:
		ids := make([]ent.Value, 0, len(m.removedprojects))
		for id := range m.removedprojects {
			ids = append(ids, id)
		}
		return ids
	case projecttype.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.removedproperties))
		for id := range m.removedproperties {
			ids = append(ids, id)
		}
		return ids
	case projecttype.EdgeWorkOrders:
		ids := make([]ent.Value, 0, len(m.removedwork_orders))
		for id := range m.removedwork_orders {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *ProjectTypeMutation) ClearedEdges() []string {
	edges := make([]string, 0, 3)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *ProjectTypeMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *ProjectTypeMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown ProjectType unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *ProjectTypeMutation) ResetEdge(name string) error {
	switch name {
	case projecttype.EdgeProjects:
		m.ResetProjects()
		return nil
	case projecttype.EdgeProperties:
		m.ResetProperties()
		return nil
	case projecttype.EdgeWorkOrders:
		m.ResetWorkOrders()
		return nil
	}
	return fmt.Errorf("unknown ProjectType edge %s", name)
}

// PropertyMutation represents an operation that mutate the Properties
// nodes in the graph.
type PropertyMutation struct {
	config
	op                     Op
	typ                    string
	id                     *int
	create_time            *time.Time
	update_time            *time.Time
	int_val                *int
	addint_val             *int
	bool_val               *bool
	float_val              *float64
	addfloat_val           *float64
	latitude_val           *float64
	addlatitude_val        *float64
	longitude_val          *float64
	addlongitude_val       *float64
	range_from_val         *float64
	addrange_from_val      *float64
	range_to_val           *float64
	addrange_to_val        *float64
	string_val             *string
	clearedFields          map[string]bool
	_type                  *int
	cleared_type           bool
	location               *int
	clearedlocation        bool
	equipment              *int
	clearedequipment       bool
	service                *int
	clearedservice         bool
	equipment_port         *int
	clearedequipment_port  bool
	link                   *int
	clearedlink            bool
	work_order             *int
	clearedwork_order      bool
	project                *int
	clearedproject         bool
	equipment_value        *int
	clearedequipment_value bool
	location_value         *int
	clearedlocation_value  bool
	service_value          *int
	clearedservice_value   bool
}

var _ ent.Mutation = (*PropertyMutation)(nil)

// newPropertyMutation creates new mutation for $n.Name.
func newPropertyMutation(c config, op Op) *PropertyMutation {
	return &PropertyMutation{
		config:        c,
		op:            op,
		typ:           TypeProperty,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m PropertyMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m PropertyMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *PropertyMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *PropertyMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *PropertyMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *PropertyMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *PropertyMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *PropertyMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *PropertyMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetIntVal sets the int_val field.
func (m *PropertyMutation) SetIntVal(i int) {
	m.int_val = &i
	m.addint_val = nil
}

// IntVal returns the int_val value in the mutation.
func (m *PropertyMutation) IntVal() (r int, exists bool) {
	v := m.int_val
	if v == nil {
		return
	}
	return *v, true
}

// AddIntVal adds i to int_val.
func (m *PropertyMutation) AddIntVal(i int) {
	if m.addint_val != nil {
		*m.addint_val += i
	} else {
		m.addint_val = &i
	}
}

// AddedIntVal returns the value that was added to the int_val field in this mutation.
func (m *PropertyMutation) AddedIntVal() (r int, exists bool) {
	v := m.addint_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearIntVal clears the value of int_val.
func (m *PropertyMutation) ClearIntVal() {
	m.int_val = nil
	m.addint_val = nil
	m.clearedFields[property.FieldIntVal] = true
}

// IntValCleared returns if the field int_val was cleared in this mutation.
func (m *PropertyMutation) IntValCleared() bool {
	return m.clearedFields[property.FieldIntVal]
}

// ResetIntVal reset all changes of the int_val field.
func (m *PropertyMutation) ResetIntVal() {
	m.int_val = nil
	m.addint_val = nil
	delete(m.clearedFields, property.FieldIntVal)
}

// SetBoolVal sets the bool_val field.
func (m *PropertyMutation) SetBoolVal(b bool) {
	m.bool_val = &b
}

// BoolVal returns the bool_val value in the mutation.
func (m *PropertyMutation) BoolVal() (r bool, exists bool) {
	v := m.bool_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearBoolVal clears the value of bool_val.
func (m *PropertyMutation) ClearBoolVal() {
	m.bool_val = nil
	m.clearedFields[property.FieldBoolVal] = true
}

// BoolValCleared returns if the field bool_val was cleared in this mutation.
func (m *PropertyMutation) BoolValCleared() bool {
	return m.clearedFields[property.FieldBoolVal]
}

// ResetBoolVal reset all changes of the bool_val field.
func (m *PropertyMutation) ResetBoolVal() {
	m.bool_val = nil
	delete(m.clearedFields, property.FieldBoolVal)
}

// SetFloatVal sets the float_val field.
func (m *PropertyMutation) SetFloatVal(f float64) {
	m.float_val = &f
	m.addfloat_val = nil
}

// FloatVal returns the float_val value in the mutation.
func (m *PropertyMutation) FloatVal() (r float64, exists bool) {
	v := m.float_val
	if v == nil {
		return
	}
	return *v, true
}

// AddFloatVal adds f to float_val.
func (m *PropertyMutation) AddFloatVal(f float64) {
	if m.addfloat_val != nil {
		*m.addfloat_val += f
	} else {
		m.addfloat_val = &f
	}
}

// AddedFloatVal returns the value that was added to the float_val field in this mutation.
func (m *PropertyMutation) AddedFloatVal() (r float64, exists bool) {
	v := m.addfloat_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearFloatVal clears the value of float_val.
func (m *PropertyMutation) ClearFloatVal() {
	m.float_val = nil
	m.addfloat_val = nil
	m.clearedFields[property.FieldFloatVal] = true
}

// FloatValCleared returns if the field float_val was cleared in this mutation.
func (m *PropertyMutation) FloatValCleared() bool {
	return m.clearedFields[property.FieldFloatVal]
}

// ResetFloatVal reset all changes of the float_val field.
func (m *PropertyMutation) ResetFloatVal() {
	m.float_val = nil
	m.addfloat_val = nil
	delete(m.clearedFields, property.FieldFloatVal)
}

// SetLatitudeVal sets the latitude_val field.
func (m *PropertyMutation) SetLatitudeVal(f float64) {
	m.latitude_val = &f
	m.addlatitude_val = nil
}

// LatitudeVal returns the latitude_val value in the mutation.
func (m *PropertyMutation) LatitudeVal() (r float64, exists bool) {
	v := m.latitude_val
	if v == nil {
		return
	}
	return *v, true
}

// AddLatitudeVal adds f to latitude_val.
func (m *PropertyMutation) AddLatitudeVal(f float64) {
	if m.addlatitude_val != nil {
		*m.addlatitude_val += f
	} else {
		m.addlatitude_val = &f
	}
}

// AddedLatitudeVal returns the value that was added to the latitude_val field in this mutation.
func (m *PropertyMutation) AddedLatitudeVal() (r float64, exists bool) {
	v := m.addlatitude_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearLatitudeVal clears the value of latitude_val.
func (m *PropertyMutation) ClearLatitudeVal() {
	m.latitude_val = nil
	m.addlatitude_val = nil
	m.clearedFields[property.FieldLatitudeVal] = true
}

// LatitudeValCleared returns if the field latitude_val was cleared in this mutation.
func (m *PropertyMutation) LatitudeValCleared() bool {
	return m.clearedFields[property.FieldLatitudeVal]
}

// ResetLatitudeVal reset all changes of the latitude_val field.
func (m *PropertyMutation) ResetLatitudeVal() {
	m.latitude_val = nil
	m.addlatitude_val = nil
	delete(m.clearedFields, property.FieldLatitudeVal)
}

// SetLongitudeVal sets the longitude_val field.
func (m *PropertyMutation) SetLongitudeVal(f float64) {
	m.longitude_val = &f
	m.addlongitude_val = nil
}

// LongitudeVal returns the longitude_val value in the mutation.
func (m *PropertyMutation) LongitudeVal() (r float64, exists bool) {
	v := m.longitude_val
	if v == nil {
		return
	}
	return *v, true
}

// AddLongitudeVal adds f to longitude_val.
func (m *PropertyMutation) AddLongitudeVal(f float64) {
	if m.addlongitude_val != nil {
		*m.addlongitude_val += f
	} else {
		m.addlongitude_val = &f
	}
}

// AddedLongitudeVal returns the value that was added to the longitude_val field in this mutation.
func (m *PropertyMutation) AddedLongitudeVal() (r float64, exists bool) {
	v := m.addlongitude_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearLongitudeVal clears the value of longitude_val.
func (m *PropertyMutation) ClearLongitudeVal() {
	m.longitude_val = nil
	m.addlongitude_val = nil
	m.clearedFields[property.FieldLongitudeVal] = true
}

// LongitudeValCleared returns if the field longitude_val was cleared in this mutation.
func (m *PropertyMutation) LongitudeValCleared() bool {
	return m.clearedFields[property.FieldLongitudeVal]
}

// ResetLongitudeVal reset all changes of the longitude_val field.
func (m *PropertyMutation) ResetLongitudeVal() {
	m.longitude_val = nil
	m.addlongitude_val = nil
	delete(m.clearedFields, property.FieldLongitudeVal)
}

// SetRangeFromVal sets the range_from_val field.
func (m *PropertyMutation) SetRangeFromVal(f float64) {
	m.range_from_val = &f
	m.addrange_from_val = nil
}

// RangeFromVal returns the range_from_val value in the mutation.
func (m *PropertyMutation) RangeFromVal() (r float64, exists bool) {
	v := m.range_from_val
	if v == nil {
		return
	}
	return *v, true
}

// AddRangeFromVal adds f to range_from_val.
func (m *PropertyMutation) AddRangeFromVal(f float64) {
	if m.addrange_from_val != nil {
		*m.addrange_from_val += f
	} else {
		m.addrange_from_val = &f
	}
}

// AddedRangeFromVal returns the value that was added to the range_from_val field in this mutation.
func (m *PropertyMutation) AddedRangeFromVal() (r float64, exists bool) {
	v := m.addrange_from_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearRangeFromVal clears the value of range_from_val.
func (m *PropertyMutation) ClearRangeFromVal() {
	m.range_from_val = nil
	m.addrange_from_val = nil
	m.clearedFields[property.FieldRangeFromVal] = true
}

// RangeFromValCleared returns if the field range_from_val was cleared in this mutation.
func (m *PropertyMutation) RangeFromValCleared() bool {
	return m.clearedFields[property.FieldRangeFromVal]
}

// ResetRangeFromVal reset all changes of the range_from_val field.
func (m *PropertyMutation) ResetRangeFromVal() {
	m.range_from_val = nil
	m.addrange_from_val = nil
	delete(m.clearedFields, property.FieldRangeFromVal)
}

// SetRangeToVal sets the range_to_val field.
func (m *PropertyMutation) SetRangeToVal(f float64) {
	m.range_to_val = &f
	m.addrange_to_val = nil
}

// RangeToVal returns the range_to_val value in the mutation.
func (m *PropertyMutation) RangeToVal() (r float64, exists bool) {
	v := m.range_to_val
	if v == nil {
		return
	}
	return *v, true
}

// AddRangeToVal adds f to range_to_val.
func (m *PropertyMutation) AddRangeToVal(f float64) {
	if m.addrange_to_val != nil {
		*m.addrange_to_val += f
	} else {
		m.addrange_to_val = &f
	}
}

// AddedRangeToVal returns the value that was added to the range_to_val field in this mutation.
func (m *PropertyMutation) AddedRangeToVal() (r float64, exists bool) {
	v := m.addrange_to_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearRangeToVal clears the value of range_to_val.
func (m *PropertyMutation) ClearRangeToVal() {
	m.range_to_val = nil
	m.addrange_to_val = nil
	m.clearedFields[property.FieldRangeToVal] = true
}

// RangeToValCleared returns if the field range_to_val was cleared in this mutation.
func (m *PropertyMutation) RangeToValCleared() bool {
	return m.clearedFields[property.FieldRangeToVal]
}

// ResetRangeToVal reset all changes of the range_to_val field.
func (m *PropertyMutation) ResetRangeToVal() {
	m.range_to_val = nil
	m.addrange_to_val = nil
	delete(m.clearedFields, property.FieldRangeToVal)
}

// SetStringVal sets the string_val field.
func (m *PropertyMutation) SetStringVal(s string) {
	m.string_val = &s
}

// StringVal returns the string_val value in the mutation.
func (m *PropertyMutation) StringVal() (r string, exists bool) {
	v := m.string_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearStringVal clears the value of string_val.
func (m *PropertyMutation) ClearStringVal() {
	m.string_val = nil
	m.clearedFields[property.FieldStringVal] = true
}

// StringValCleared returns if the field string_val was cleared in this mutation.
func (m *PropertyMutation) StringValCleared() bool {
	return m.clearedFields[property.FieldStringVal]
}

// ResetStringVal reset all changes of the string_val field.
func (m *PropertyMutation) ResetStringVal() {
	m.string_val = nil
	delete(m.clearedFields, property.FieldStringVal)
}

// SetTypeID sets the type edge to PropertyType by id.
func (m *PropertyMutation) SetTypeID(id int) {
	m._type = &id
}

// ClearType clears the type edge to PropertyType.
func (m *PropertyMutation) ClearType() {
	m.cleared_type = true
}

// TypeCleared returns if the edge type was cleared.
func (m *PropertyMutation) TypeCleared() bool {
	return m.cleared_type
}

// TypeID returns the type id in the mutation.
func (m *PropertyMutation) TypeID() (id int, exists bool) {
	if m._type != nil {
		return *m._type, true
	}
	return
}

// TypeIDs returns the type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// TypeID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) TypeIDs() (ids []int) {
	if id := m._type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetType reset all changes of the type edge.
func (m *PropertyMutation) ResetType() {
	m._type = nil
	m.cleared_type = false
}

// SetLocationID sets the location edge to Location by id.
func (m *PropertyMutation) SetLocationID(id int) {
	m.location = &id
}

// ClearLocation clears the location edge to Location.
func (m *PropertyMutation) ClearLocation() {
	m.clearedlocation = true
}

// LocationCleared returns if the edge location was cleared.
func (m *PropertyMutation) LocationCleared() bool {
	return m.clearedlocation
}

// LocationID returns the location id in the mutation.
func (m *PropertyMutation) LocationID() (id int, exists bool) {
	if m.location != nil {
		return *m.location, true
	}
	return
}

// LocationIDs returns the location ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LocationID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) LocationIDs() (ids []int) {
	if id := m.location; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLocation reset all changes of the location edge.
func (m *PropertyMutation) ResetLocation() {
	m.location = nil
	m.clearedlocation = false
}

// SetEquipmentID sets the equipment edge to Equipment by id.
func (m *PropertyMutation) SetEquipmentID(id int) {
	m.equipment = &id
}

// ClearEquipment clears the equipment edge to Equipment.
func (m *PropertyMutation) ClearEquipment() {
	m.clearedequipment = true
}

// EquipmentCleared returns if the edge equipment was cleared.
func (m *PropertyMutation) EquipmentCleared() bool {
	return m.clearedequipment
}

// EquipmentID returns the equipment id in the mutation.
func (m *PropertyMutation) EquipmentID() (id int, exists bool) {
	if m.equipment != nil {
		return *m.equipment, true
	}
	return
}

// EquipmentIDs returns the equipment ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// EquipmentID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) EquipmentIDs() (ids []int) {
	if id := m.equipment; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEquipment reset all changes of the equipment edge.
func (m *PropertyMutation) ResetEquipment() {
	m.equipment = nil
	m.clearedequipment = false
}

// SetServiceID sets the service edge to Service by id.
func (m *PropertyMutation) SetServiceID(id int) {
	m.service = &id
}

// ClearService clears the service edge to Service.
func (m *PropertyMutation) ClearService() {
	m.clearedservice = true
}

// ServiceCleared returns if the edge service was cleared.
func (m *PropertyMutation) ServiceCleared() bool {
	return m.clearedservice
}

// ServiceID returns the service id in the mutation.
func (m *PropertyMutation) ServiceID() (id int, exists bool) {
	if m.service != nil {
		return *m.service, true
	}
	return
}

// ServiceIDs returns the service ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ServiceID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) ServiceIDs() (ids []int) {
	if id := m.service; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetService reset all changes of the service edge.
func (m *PropertyMutation) ResetService() {
	m.service = nil
	m.clearedservice = false
}

// SetEquipmentPortID sets the equipment_port edge to EquipmentPort by id.
func (m *PropertyMutation) SetEquipmentPortID(id int) {
	m.equipment_port = &id
}

// ClearEquipmentPort clears the equipment_port edge to EquipmentPort.
func (m *PropertyMutation) ClearEquipmentPort() {
	m.clearedequipment_port = true
}

// EquipmentPortCleared returns if the edge equipment_port was cleared.
func (m *PropertyMutation) EquipmentPortCleared() bool {
	return m.clearedequipment_port
}

// EquipmentPortID returns the equipment_port id in the mutation.
func (m *PropertyMutation) EquipmentPortID() (id int, exists bool) {
	if m.equipment_port != nil {
		return *m.equipment_port, true
	}
	return
}

// EquipmentPortIDs returns the equipment_port ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// EquipmentPortID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) EquipmentPortIDs() (ids []int) {
	if id := m.equipment_port; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEquipmentPort reset all changes of the equipment_port edge.
func (m *PropertyMutation) ResetEquipmentPort() {
	m.equipment_port = nil
	m.clearedequipment_port = false
}

// SetLinkID sets the link edge to Link by id.
func (m *PropertyMutation) SetLinkID(id int) {
	m.link = &id
}

// ClearLink clears the link edge to Link.
func (m *PropertyMutation) ClearLink() {
	m.clearedlink = true
}

// LinkCleared returns if the edge link was cleared.
func (m *PropertyMutation) LinkCleared() bool {
	return m.clearedlink
}

// LinkID returns the link id in the mutation.
func (m *PropertyMutation) LinkID() (id int, exists bool) {
	if m.link != nil {
		return *m.link, true
	}
	return
}

// LinkIDs returns the link ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LinkID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) LinkIDs() (ids []int) {
	if id := m.link; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLink reset all changes of the link edge.
func (m *PropertyMutation) ResetLink() {
	m.link = nil
	m.clearedlink = false
}

// SetWorkOrderID sets the work_order edge to WorkOrder by id.
func (m *PropertyMutation) SetWorkOrderID(id int) {
	m.work_order = &id
}

// ClearWorkOrder clears the work_order edge to WorkOrder.
func (m *PropertyMutation) ClearWorkOrder() {
	m.clearedwork_order = true
}

// WorkOrderCleared returns if the edge work_order was cleared.
func (m *PropertyMutation) WorkOrderCleared() bool {
	return m.clearedwork_order
}

// WorkOrderID returns the work_order id in the mutation.
func (m *PropertyMutation) WorkOrderID() (id int, exists bool) {
	if m.work_order != nil {
		return *m.work_order, true
	}
	return
}

// WorkOrderIDs returns the work_order ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// WorkOrderID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) WorkOrderIDs() (ids []int) {
	if id := m.work_order; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetWorkOrder reset all changes of the work_order edge.
func (m *PropertyMutation) ResetWorkOrder() {
	m.work_order = nil
	m.clearedwork_order = false
}

// SetProjectID sets the project edge to Project by id.
func (m *PropertyMutation) SetProjectID(id int) {
	m.project = &id
}

// ClearProject clears the project edge to Project.
func (m *PropertyMutation) ClearProject() {
	m.clearedproject = true
}

// ProjectCleared returns if the edge project was cleared.
func (m *PropertyMutation) ProjectCleared() bool {
	return m.clearedproject
}

// ProjectID returns the project id in the mutation.
func (m *PropertyMutation) ProjectID() (id int, exists bool) {
	if m.project != nil {
		return *m.project, true
	}
	return
}

// ProjectIDs returns the project ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ProjectID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) ProjectIDs() (ids []int) {
	if id := m.project; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetProject reset all changes of the project edge.
func (m *PropertyMutation) ResetProject() {
	m.project = nil
	m.clearedproject = false
}

// SetEquipmentValueID sets the equipment_value edge to Equipment by id.
func (m *PropertyMutation) SetEquipmentValueID(id int) {
	m.equipment_value = &id
}

// ClearEquipmentValue clears the equipment_value edge to Equipment.
func (m *PropertyMutation) ClearEquipmentValue() {
	m.clearedequipment_value = true
}

// EquipmentValueCleared returns if the edge equipment_value was cleared.
func (m *PropertyMutation) EquipmentValueCleared() bool {
	return m.clearedequipment_value
}

// EquipmentValueID returns the equipment_value id in the mutation.
func (m *PropertyMutation) EquipmentValueID() (id int, exists bool) {
	if m.equipment_value != nil {
		return *m.equipment_value, true
	}
	return
}

// EquipmentValueIDs returns the equipment_value ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// EquipmentValueID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) EquipmentValueIDs() (ids []int) {
	if id := m.equipment_value; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEquipmentValue reset all changes of the equipment_value edge.
func (m *PropertyMutation) ResetEquipmentValue() {
	m.equipment_value = nil
	m.clearedequipment_value = false
}

// SetLocationValueID sets the location_value edge to Location by id.
func (m *PropertyMutation) SetLocationValueID(id int) {
	m.location_value = &id
}

// ClearLocationValue clears the location_value edge to Location.
func (m *PropertyMutation) ClearLocationValue() {
	m.clearedlocation_value = true
}

// LocationValueCleared returns if the edge location_value was cleared.
func (m *PropertyMutation) LocationValueCleared() bool {
	return m.clearedlocation_value
}

// LocationValueID returns the location_value id in the mutation.
func (m *PropertyMutation) LocationValueID() (id int, exists bool) {
	if m.location_value != nil {
		return *m.location_value, true
	}
	return
}

// LocationValueIDs returns the location_value ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LocationValueID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) LocationValueIDs() (ids []int) {
	if id := m.location_value; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLocationValue reset all changes of the location_value edge.
func (m *PropertyMutation) ResetLocationValue() {
	m.location_value = nil
	m.clearedlocation_value = false
}

// SetServiceValueID sets the service_value edge to Service by id.
func (m *PropertyMutation) SetServiceValueID(id int) {
	m.service_value = &id
}

// ClearServiceValue clears the service_value edge to Service.
func (m *PropertyMutation) ClearServiceValue() {
	m.clearedservice_value = true
}

// ServiceValueCleared returns if the edge service_value was cleared.
func (m *PropertyMutation) ServiceValueCleared() bool {
	return m.clearedservice_value
}

// ServiceValueID returns the service_value id in the mutation.
func (m *PropertyMutation) ServiceValueID() (id int, exists bool) {
	if m.service_value != nil {
		return *m.service_value, true
	}
	return
}

// ServiceValueIDs returns the service_value ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ServiceValueID instead. It exists only for internal usage by the builders.
func (m *PropertyMutation) ServiceValueIDs() (ids []int) {
	if id := m.service_value; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetServiceValue reset all changes of the service_value edge.
func (m *PropertyMutation) ResetServiceValue() {
	m.service_value = nil
	m.clearedservice_value = false
}

// Op returns the operation name.
func (m *PropertyMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Property).
func (m *PropertyMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *PropertyMutation) Fields() []string {
	fields := make([]string, 0, 10)
	if m.create_time != nil {
		fields = append(fields, property.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, property.FieldUpdateTime)
	}
	if m.int_val != nil {
		fields = append(fields, property.FieldIntVal)
	}
	if m.bool_val != nil {
		fields = append(fields, property.FieldBoolVal)
	}
	if m.float_val != nil {
		fields = append(fields, property.FieldFloatVal)
	}
	if m.latitude_val != nil {
		fields = append(fields, property.FieldLatitudeVal)
	}
	if m.longitude_val != nil {
		fields = append(fields, property.FieldLongitudeVal)
	}
	if m.range_from_val != nil {
		fields = append(fields, property.FieldRangeFromVal)
	}
	if m.range_to_val != nil {
		fields = append(fields, property.FieldRangeToVal)
	}
	if m.string_val != nil {
		fields = append(fields, property.FieldStringVal)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *PropertyMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case property.FieldCreateTime:
		return m.CreateTime()
	case property.FieldUpdateTime:
		return m.UpdateTime()
	case property.FieldIntVal:
		return m.IntVal()
	case property.FieldBoolVal:
		return m.BoolVal()
	case property.FieldFloatVal:
		return m.FloatVal()
	case property.FieldLatitudeVal:
		return m.LatitudeVal()
	case property.FieldLongitudeVal:
		return m.LongitudeVal()
	case property.FieldRangeFromVal:
		return m.RangeFromVal()
	case property.FieldRangeToVal:
		return m.RangeToVal()
	case property.FieldStringVal:
		return m.StringVal()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *PropertyMutation) SetField(name string, value ent.Value) error {
	switch name {
	case property.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case property.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case property.FieldIntVal:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIntVal(v)
		return nil
	case property.FieldBoolVal:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBoolVal(v)
		return nil
	case property.FieldFloatVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFloatVal(v)
		return nil
	case property.FieldLatitudeVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLatitudeVal(v)
		return nil
	case property.FieldLongitudeVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLongitudeVal(v)
		return nil
	case property.FieldRangeFromVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRangeFromVal(v)
		return nil
	case property.FieldRangeToVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRangeToVal(v)
		return nil
	case property.FieldStringVal:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStringVal(v)
		return nil
	}
	return fmt.Errorf("unknown Property field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *PropertyMutation) AddedFields() []string {
	var fields []string
	if m.addint_val != nil {
		fields = append(fields, property.FieldIntVal)
	}
	if m.addfloat_val != nil {
		fields = append(fields, property.FieldFloatVal)
	}
	if m.addlatitude_val != nil {
		fields = append(fields, property.FieldLatitudeVal)
	}
	if m.addlongitude_val != nil {
		fields = append(fields, property.FieldLongitudeVal)
	}
	if m.addrange_from_val != nil {
		fields = append(fields, property.FieldRangeFromVal)
	}
	if m.addrange_to_val != nil {
		fields = append(fields, property.FieldRangeToVal)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *PropertyMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case property.FieldIntVal:
		return m.AddedIntVal()
	case property.FieldFloatVal:
		return m.AddedFloatVal()
	case property.FieldLatitudeVal:
		return m.AddedLatitudeVal()
	case property.FieldLongitudeVal:
		return m.AddedLongitudeVal()
	case property.FieldRangeFromVal:
		return m.AddedRangeFromVal()
	case property.FieldRangeToVal:
		return m.AddedRangeToVal()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *PropertyMutation) AddField(name string, value ent.Value) error {
	switch name {
	case property.FieldIntVal:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIntVal(v)
		return nil
	case property.FieldFloatVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddFloatVal(v)
		return nil
	case property.FieldLatitudeVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLatitudeVal(v)
		return nil
	case property.FieldLongitudeVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLongitudeVal(v)
		return nil
	case property.FieldRangeFromVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddRangeFromVal(v)
		return nil
	case property.FieldRangeToVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddRangeToVal(v)
		return nil
	}
	return fmt.Errorf("unknown Property numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *PropertyMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[property.FieldIntVal] {
		fields = append(fields, property.FieldIntVal)
	}
	if m.clearedFields[property.FieldBoolVal] {
		fields = append(fields, property.FieldBoolVal)
	}
	if m.clearedFields[property.FieldFloatVal] {
		fields = append(fields, property.FieldFloatVal)
	}
	if m.clearedFields[property.FieldLatitudeVal] {
		fields = append(fields, property.FieldLatitudeVal)
	}
	if m.clearedFields[property.FieldLongitudeVal] {
		fields = append(fields, property.FieldLongitudeVal)
	}
	if m.clearedFields[property.FieldRangeFromVal] {
		fields = append(fields, property.FieldRangeFromVal)
	}
	if m.clearedFields[property.FieldRangeToVal] {
		fields = append(fields, property.FieldRangeToVal)
	}
	if m.clearedFields[property.FieldStringVal] {
		fields = append(fields, property.FieldStringVal)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *PropertyMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *PropertyMutation) ClearField(name string) error {
	switch name {
	case property.FieldIntVal:
		m.ClearIntVal()
		return nil
	case property.FieldBoolVal:
		m.ClearBoolVal()
		return nil
	case property.FieldFloatVal:
		m.ClearFloatVal()
		return nil
	case property.FieldLatitudeVal:
		m.ClearLatitudeVal()
		return nil
	case property.FieldLongitudeVal:
		m.ClearLongitudeVal()
		return nil
	case property.FieldRangeFromVal:
		m.ClearRangeFromVal()
		return nil
	case property.FieldRangeToVal:
		m.ClearRangeToVal()
		return nil
	case property.FieldStringVal:
		m.ClearStringVal()
		return nil
	}
	return fmt.Errorf("unknown Property nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *PropertyMutation) ResetField(name string) error {
	switch name {
	case property.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case property.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case property.FieldIntVal:
		m.ResetIntVal()
		return nil
	case property.FieldBoolVal:
		m.ResetBoolVal()
		return nil
	case property.FieldFloatVal:
		m.ResetFloatVal()
		return nil
	case property.FieldLatitudeVal:
		m.ResetLatitudeVal()
		return nil
	case property.FieldLongitudeVal:
		m.ResetLongitudeVal()
		return nil
	case property.FieldRangeFromVal:
		m.ResetRangeFromVal()
		return nil
	case property.FieldRangeToVal:
		m.ResetRangeToVal()
		return nil
	case property.FieldStringVal:
		m.ResetStringVal()
		return nil
	}
	return fmt.Errorf("unknown Property field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *PropertyMutation) AddedEdges() []string {
	edges := make([]string, 0, 11)
	if m._type != nil {
		edges = append(edges, property.EdgeType)
	}
	if m.location != nil {
		edges = append(edges, property.EdgeLocation)
	}
	if m.equipment != nil {
		edges = append(edges, property.EdgeEquipment)
	}
	if m.service != nil {
		edges = append(edges, property.EdgeService)
	}
	if m.equipment_port != nil {
		edges = append(edges, property.EdgeEquipmentPort)
	}
	if m.link != nil {
		edges = append(edges, property.EdgeLink)
	}
	if m.work_order != nil {
		edges = append(edges, property.EdgeWorkOrder)
	}
	if m.project != nil {
		edges = append(edges, property.EdgeProject)
	}
	if m.equipment_value != nil {
		edges = append(edges, property.EdgeEquipmentValue)
	}
	if m.location_value != nil {
		edges = append(edges, property.EdgeLocationValue)
	}
	if m.service_value != nil {
		edges = append(edges, property.EdgeServiceValue)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *PropertyMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case property.EdgeType:
		if id := m._type; id != nil {
			return []ent.Value{*id}
		}
	case property.EdgeLocation:
		if id := m.location; id != nil {
			return []ent.Value{*id}
		}
	case property.EdgeEquipment:
		if id := m.equipment; id != nil {
			return []ent.Value{*id}
		}
	case property.EdgeService:
		if id := m.service; id != nil {
			return []ent.Value{*id}
		}
	case property.EdgeEquipmentPort:
		if id := m.equipment_port; id != nil {
			return []ent.Value{*id}
		}
	case property.EdgeLink:
		if id := m.link; id != nil {
			return []ent.Value{*id}
		}
	case property.EdgeWorkOrder:
		if id := m.work_order; id != nil {
			return []ent.Value{*id}
		}
	case property.EdgeProject:
		if id := m.project; id != nil {
			return []ent.Value{*id}
		}
	case property.EdgeEquipmentValue:
		if id := m.equipment_value; id != nil {
			return []ent.Value{*id}
		}
	case property.EdgeLocationValue:
		if id := m.location_value; id != nil {
			return []ent.Value{*id}
		}
	case property.EdgeServiceValue:
		if id := m.service_value; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *PropertyMutation) RemovedEdges() []string {
	edges := make([]string, 0, 11)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *PropertyMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *PropertyMutation) ClearedEdges() []string {
	edges := make([]string, 0, 11)
	if m.cleared_type {
		edges = append(edges, property.EdgeType)
	}
	if m.clearedlocation {
		edges = append(edges, property.EdgeLocation)
	}
	if m.clearedequipment {
		edges = append(edges, property.EdgeEquipment)
	}
	if m.clearedservice {
		edges = append(edges, property.EdgeService)
	}
	if m.clearedequipment_port {
		edges = append(edges, property.EdgeEquipmentPort)
	}
	if m.clearedlink {
		edges = append(edges, property.EdgeLink)
	}
	if m.clearedwork_order {
		edges = append(edges, property.EdgeWorkOrder)
	}
	if m.clearedproject {
		edges = append(edges, property.EdgeProject)
	}
	if m.clearedequipment_value {
		edges = append(edges, property.EdgeEquipmentValue)
	}
	if m.clearedlocation_value {
		edges = append(edges, property.EdgeLocationValue)
	}
	if m.clearedservice_value {
		edges = append(edges, property.EdgeServiceValue)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *PropertyMutation) EdgeCleared(name string) bool {
	switch name {
	case property.EdgeType:
		return m.cleared_type
	case property.EdgeLocation:
		return m.clearedlocation
	case property.EdgeEquipment:
		return m.clearedequipment
	case property.EdgeService:
		return m.clearedservice
	case property.EdgeEquipmentPort:
		return m.clearedequipment_port
	case property.EdgeLink:
		return m.clearedlink
	case property.EdgeWorkOrder:
		return m.clearedwork_order
	case property.EdgeProject:
		return m.clearedproject
	case property.EdgeEquipmentValue:
		return m.clearedequipment_value
	case property.EdgeLocationValue:
		return m.clearedlocation_value
	case property.EdgeServiceValue:
		return m.clearedservice_value
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *PropertyMutation) ClearEdge(name string) error {
	switch name {
	case property.EdgeType:
		m.ClearType()
		return nil
	case property.EdgeLocation:
		m.ClearLocation()
		return nil
	case property.EdgeEquipment:
		m.ClearEquipment()
		return nil
	case property.EdgeService:
		m.ClearService()
		return nil
	case property.EdgeEquipmentPort:
		m.ClearEquipmentPort()
		return nil
	case property.EdgeLink:
		m.ClearLink()
		return nil
	case property.EdgeWorkOrder:
		m.ClearWorkOrder()
		return nil
	case property.EdgeProject:
		m.ClearProject()
		return nil
	case property.EdgeEquipmentValue:
		m.ClearEquipmentValue()
		return nil
	case property.EdgeLocationValue:
		m.ClearLocationValue()
		return nil
	case property.EdgeServiceValue:
		m.ClearServiceValue()
		return nil
	}
	return fmt.Errorf("unknown Property unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *PropertyMutation) ResetEdge(name string) error {
	switch name {
	case property.EdgeType:
		m.ResetType()
		return nil
	case property.EdgeLocation:
		m.ResetLocation()
		return nil
	case property.EdgeEquipment:
		m.ResetEquipment()
		return nil
	case property.EdgeService:
		m.ResetService()
		return nil
	case property.EdgeEquipmentPort:
		m.ResetEquipmentPort()
		return nil
	case property.EdgeLink:
		m.ResetLink()
		return nil
	case property.EdgeWorkOrder:
		m.ResetWorkOrder()
		return nil
	case property.EdgeProject:
		m.ResetProject()
		return nil
	case property.EdgeEquipmentValue:
		m.ResetEquipmentValue()
		return nil
	case property.EdgeLocationValue:
		m.ResetLocationValue()
		return nil
	case property.EdgeServiceValue:
		m.ResetServiceValue()
		return nil
	}
	return fmt.Errorf("unknown Property edge %s", name)
}

// PropertyTypeMutation represents an operation that mutate the PropertyTypes
// nodes in the graph.
type PropertyTypeMutation struct {
	config
	op                              Op
	typ                             string
	id                              *int
	create_time                     *time.Time
	update_time                     *time.Time
	_type                           *string
	name                            *string
	index                           *int
	addindex                        *int
	category                        *string
	int_val                         *int
	addint_val                      *int
	bool_val                        *bool
	float_val                       *float64
	addfloat_val                    *float64
	latitude_val                    *float64
	addlatitude_val                 *float64
	longitude_val                   *float64
	addlongitude_val                *float64
	string_val                      *string
	range_from_val                  *float64
	addrange_from_val               *float64
	range_to_val                    *float64
	addrange_to_val                 *float64
	is_instance_property            *bool
	editable                        *bool
	mandatory                       *bool
	deleted                         *bool
	clearedFields                   map[string]bool
	properties                      map[int]struct{}
	removedproperties               map[int]struct{}
	location_type                   *int
	clearedlocation_type            bool
	equipment_port_type             *int
	clearedequipment_port_type      bool
	link_equipment_port_type        *int
	clearedlink_equipment_port_type bool
	equipment_type                  *int
	clearedequipment_type           bool
	service_type                    *int
	clearedservice_type             bool
	work_order_type                 *int
	clearedwork_order_type          bool
	project_type                    *int
	clearedproject_type             bool
}

var _ ent.Mutation = (*PropertyTypeMutation)(nil)

// newPropertyTypeMutation creates new mutation for $n.Name.
func newPropertyTypeMutation(c config, op Op) *PropertyTypeMutation {
	return &PropertyTypeMutation{
		config:        c,
		op:            op,
		typ:           TypePropertyType,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m PropertyTypeMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m PropertyTypeMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *PropertyTypeMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *PropertyTypeMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *PropertyTypeMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *PropertyTypeMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *PropertyTypeMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *PropertyTypeMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *PropertyTypeMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetType sets the type field.
func (m *PropertyTypeMutation) SetType(s string) {
	m._type = &s
}

// GetType returns the type value in the mutation.
func (m *PropertyTypeMutation) GetType() (r string, exists bool) {
	v := m._type
	if v == nil {
		return
	}
	return *v, true
}

// ResetType reset all changes of the type field.
func (m *PropertyTypeMutation) ResetType() {
	m._type = nil
}

// SetName sets the name field.
func (m *PropertyTypeMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *PropertyTypeMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *PropertyTypeMutation) ResetName() {
	m.name = nil
}

// SetIndex sets the index field.
func (m *PropertyTypeMutation) SetIndex(i int) {
	m.index = &i
	m.addindex = nil
}

// Index returns the index value in the mutation.
func (m *PropertyTypeMutation) Index() (r int, exists bool) {
	v := m.index
	if v == nil {
		return
	}
	return *v, true
}

// AddIndex adds i to index.
func (m *PropertyTypeMutation) AddIndex(i int) {
	if m.addindex != nil {
		*m.addindex += i
	} else {
		m.addindex = &i
	}
}

// AddedIndex returns the value that was added to the index field in this mutation.
func (m *PropertyTypeMutation) AddedIndex() (r int, exists bool) {
	v := m.addindex
	if v == nil {
		return
	}
	return *v, true
}

// ClearIndex clears the value of index.
func (m *PropertyTypeMutation) ClearIndex() {
	m.index = nil
	m.addindex = nil
	m.clearedFields[propertytype.FieldIndex] = true
}

// IndexCleared returns if the field index was cleared in this mutation.
func (m *PropertyTypeMutation) IndexCleared() bool {
	return m.clearedFields[propertytype.FieldIndex]
}

// ResetIndex reset all changes of the index field.
func (m *PropertyTypeMutation) ResetIndex() {
	m.index = nil
	m.addindex = nil
	delete(m.clearedFields, propertytype.FieldIndex)
}

// SetCategory sets the category field.
func (m *PropertyTypeMutation) SetCategory(s string) {
	m.category = &s
}

// Category returns the category value in the mutation.
func (m *PropertyTypeMutation) Category() (r string, exists bool) {
	v := m.category
	if v == nil {
		return
	}
	return *v, true
}

// ClearCategory clears the value of category.
func (m *PropertyTypeMutation) ClearCategory() {
	m.category = nil
	m.clearedFields[propertytype.FieldCategory] = true
}

// CategoryCleared returns if the field category was cleared in this mutation.
func (m *PropertyTypeMutation) CategoryCleared() bool {
	return m.clearedFields[propertytype.FieldCategory]
}

// ResetCategory reset all changes of the category field.
func (m *PropertyTypeMutation) ResetCategory() {
	m.category = nil
	delete(m.clearedFields, propertytype.FieldCategory)
}

// SetIntVal sets the int_val field.
func (m *PropertyTypeMutation) SetIntVal(i int) {
	m.int_val = &i
	m.addint_val = nil
}

// IntVal returns the int_val value in the mutation.
func (m *PropertyTypeMutation) IntVal() (r int, exists bool) {
	v := m.int_val
	if v == nil {
		return
	}
	return *v, true
}

// AddIntVal adds i to int_val.
func (m *PropertyTypeMutation) AddIntVal(i int) {
	if m.addint_val != nil {
		*m.addint_val += i
	} else {
		m.addint_val = &i
	}
}

// AddedIntVal returns the value that was added to the int_val field in this mutation.
func (m *PropertyTypeMutation) AddedIntVal() (r int, exists bool) {
	v := m.addint_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearIntVal clears the value of int_val.
func (m *PropertyTypeMutation) ClearIntVal() {
	m.int_val = nil
	m.addint_val = nil
	m.clearedFields[propertytype.FieldIntVal] = true
}

// IntValCleared returns if the field int_val was cleared in this mutation.
func (m *PropertyTypeMutation) IntValCleared() bool {
	return m.clearedFields[propertytype.FieldIntVal]
}

// ResetIntVal reset all changes of the int_val field.
func (m *PropertyTypeMutation) ResetIntVal() {
	m.int_val = nil
	m.addint_val = nil
	delete(m.clearedFields, propertytype.FieldIntVal)
}

// SetBoolVal sets the bool_val field.
func (m *PropertyTypeMutation) SetBoolVal(b bool) {
	m.bool_val = &b
}

// BoolVal returns the bool_val value in the mutation.
func (m *PropertyTypeMutation) BoolVal() (r bool, exists bool) {
	v := m.bool_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearBoolVal clears the value of bool_val.
func (m *PropertyTypeMutation) ClearBoolVal() {
	m.bool_val = nil
	m.clearedFields[propertytype.FieldBoolVal] = true
}

// BoolValCleared returns if the field bool_val was cleared in this mutation.
func (m *PropertyTypeMutation) BoolValCleared() bool {
	return m.clearedFields[propertytype.FieldBoolVal]
}

// ResetBoolVal reset all changes of the bool_val field.
func (m *PropertyTypeMutation) ResetBoolVal() {
	m.bool_val = nil
	delete(m.clearedFields, propertytype.FieldBoolVal)
}

// SetFloatVal sets the float_val field.
func (m *PropertyTypeMutation) SetFloatVal(f float64) {
	m.float_val = &f
	m.addfloat_val = nil
}

// FloatVal returns the float_val value in the mutation.
func (m *PropertyTypeMutation) FloatVal() (r float64, exists bool) {
	v := m.float_val
	if v == nil {
		return
	}
	return *v, true
}

// AddFloatVal adds f to float_val.
func (m *PropertyTypeMutation) AddFloatVal(f float64) {
	if m.addfloat_val != nil {
		*m.addfloat_val += f
	} else {
		m.addfloat_val = &f
	}
}

// AddedFloatVal returns the value that was added to the float_val field in this mutation.
func (m *PropertyTypeMutation) AddedFloatVal() (r float64, exists bool) {
	v := m.addfloat_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearFloatVal clears the value of float_val.
func (m *PropertyTypeMutation) ClearFloatVal() {
	m.float_val = nil
	m.addfloat_val = nil
	m.clearedFields[propertytype.FieldFloatVal] = true
}

// FloatValCleared returns if the field float_val was cleared in this mutation.
func (m *PropertyTypeMutation) FloatValCleared() bool {
	return m.clearedFields[propertytype.FieldFloatVal]
}

// ResetFloatVal reset all changes of the float_val field.
func (m *PropertyTypeMutation) ResetFloatVal() {
	m.float_val = nil
	m.addfloat_val = nil
	delete(m.clearedFields, propertytype.FieldFloatVal)
}

// SetLatitudeVal sets the latitude_val field.
func (m *PropertyTypeMutation) SetLatitudeVal(f float64) {
	m.latitude_val = &f
	m.addlatitude_val = nil
}

// LatitudeVal returns the latitude_val value in the mutation.
func (m *PropertyTypeMutation) LatitudeVal() (r float64, exists bool) {
	v := m.latitude_val
	if v == nil {
		return
	}
	return *v, true
}

// AddLatitudeVal adds f to latitude_val.
func (m *PropertyTypeMutation) AddLatitudeVal(f float64) {
	if m.addlatitude_val != nil {
		*m.addlatitude_val += f
	} else {
		m.addlatitude_val = &f
	}
}

// AddedLatitudeVal returns the value that was added to the latitude_val field in this mutation.
func (m *PropertyTypeMutation) AddedLatitudeVal() (r float64, exists bool) {
	v := m.addlatitude_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearLatitudeVal clears the value of latitude_val.
func (m *PropertyTypeMutation) ClearLatitudeVal() {
	m.latitude_val = nil
	m.addlatitude_val = nil
	m.clearedFields[propertytype.FieldLatitudeVal] = true
}

// LatitudeValCleared returns if the field latitude_val was cleared in this mutation.
func (m *PropertyTypeMutation) LatitudeValCleared() bool {
	return m.clearedFields[propertytype.FieldLatitudeVal]
}

// ResetLatitudeVal reset all changes of the latitude_val field.
func (m *PropertyTypeMutation) ResetLatitudeVal() {
	m.latitude_val = nil
	m.addlatitude_val = nil
	delete(m.clearedFields, propertytype.FieldLatitudeVal)
}

// SetLongitudeVal sets the longitude_val field.
func (m *PropertyTypeMutation) SetLongitudeVal(f float64) {
	m.longitude_val = &f
	m.addlongitude_val = nil
}

// LongitudeVal returns the longitude_val value in the mutation.
func (m *PropertyTypeMutation) LongitudeVal() (r float64, exists bool) {
	v := m.longitude_val
	if v == nil {
		return
	}
	return *v, true
}

// AddLongitudeVal adds f to longitude_val.
func (m *PropertyTypeMutation) AddLongitudeVal(f float64) {
	if m.addlongitude_val != nil {
		*m.addlongitude_val += f
	} else {
		m.addlongitude_val = &f
	}
}

// AddedLongitudeVal returns the value that was added to the longitude_val field in this mutation.
func (m *PropertyTypeMutation) AddedLongitudeVal() (r float64, exists bool) {
	v := m.addlongitude_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearLongitudeVal clears the value of longitude_val.
func (m *PropertyTypeMutation) ClearLongitudeVal() {
	m.longitude_val = nil
	m.addlongitude_val = nil
	m.clearedFields[propertytype.FieldLongitudeVal] = true
}

// LongitudeValCleared returns if the field longitude_val was cleared in this mutation.
func (m *PropertyTypeMutation) LongitudeValCleared() bool {
	return m.clearedFields[propertytype.FieldLongitudeVal]
}

// ResetLongitudeVal reset all changes of the longitude_val field.
func (m *PropertyTypeMutation) ResetLongitudeVal() {
	m.longitude_val = nil
	m.addlongitude_val = nil
	delete(m.clearedFields, propertytype.FieldLongitudeVal)
}

// SetStringVal sets the string_val field.
func (m *PropertyTypeMutation) SetStringVal(s string) {
	m.string_val = &s
}

// StringVal returns the string_val value in the mutation.
func (m *PropertyTypeMutation) StringVal() (r string, exists bool) {
	v := m.string_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearStringVal clears the value of string_val.
func (m *PropertyTypeMutation) ClearStringVal() {
	m.string_val = nil
	m.clearedFields[propertytype.FieldStringVal] = true
}

// StringValCleared returns if the field string_val was cleared in this mutation.
func (m *PropertyTypeMutation) StringValCleared() bool {
	return m.clearedFields[propertytype.FieldStringVal]
}

// ResetStringVal reset all changes of the string_val field.
func (m *PropertyTypeMutation) ResetStringVal() {
	m.string_val = nil
	delete(m.clearedFields, propertytype.FieldStringVal)
}

// SetRangeFromVal sets the range_from_val field.
func (m *PropertyTypeMutation) SetRangeFromVal(f float64) {
	m.range_from_val = &f
	m.addrange_from_val = nil
}

// RangeFromVal returns the range_from_val value in the mutation.
func (m *PropertyTypeMutation) RangeFromVal() (r float64, exists bool) {
	v := m.range_from_val
	if v == nil {
		return
	}
	return *v, true
}

// AddRangeFromVal adds f to range_from_val.
func (m *PropertyTypeMutation) AddRangeFromVal(f float64) {
	if m.addrange_from_val != nil {
		*m.addrange_from_val += f
	} else {
		m.addrange_from_val = &f
	}
}

// AddedRangeFromVal returns the value that was added to the range_from_val field in this mutation.
func (m *PropertyTypeMutation) AddedRangeFromVal() (r float64, exists bool) {
	v := m.addrange_from_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearRangeFromVal clears the value of range_from_val.
func (m *PropertyTypeMutation) ClearRangeFromVal() {
	m.range_from_val = nil
	m.addrange_from_val = nil
	m.clearedFields[propertytype.FieldRangeFromVal] = true
}

// RangeFromValCleared returns if the field range_from_val was cleared in this mutation.
func (m *PropertyTypeMutation) RangeFromValCleared() bool {
	return m.clearedFields[propertytype.FieldRangeFromVal]
}

// ResetRangeFromVal reset all changes of the range_from_val field.
func (m *PropertyTypeMutation) ResetRangeFromVal() {
	m.range_from_val = nil
	m.addrange_from_val = nil
	delete(m.clearedFields, propertytype.FieldRangeFromVal)
}

// SetRangeToVal sets the range_to_val field.
func (m *PropertyTypeMutation) SetRangeToVal(f float64) {
	m.range_to_val = &f
	m.addrange_to_val = nil
}

// RangeToVal returns the range_to_val value in the mutation.
func (m *PropertyTypeMutation) RangeToVal() (r float64, exists bool) {
	v := m.range_to_val
	if v == nil {
		return
	}
	return *v, true
}

// AddRangeToVal adds f to range_to_val.
func (m *PropertyTypeMutation) AddRangeToVal(f float64) {
	if m.addrange_to_val != nil {
		*m.addrange_to_val += f
	} else {
		m.addrange_to_val = &f
	}
}

// AddedRangeToVal returns the value that was added to the range_to_val field in this mutation.
func (m *PropertyTypeMutation) AddedRangeToVal() (r float64, exists bool) {
	v := m.addrange_to_val
	if v == nil {
		return
	}
	return *v, true
}

// ClearRangeToVal clears the value of range_to_val.
func (m *PropertyTypeMutation) ClearRangeToVal() {
	m.range_to_val = nil
	m.addrange_to_val = nil
	m.clearedFields[propertytype.FieldRangeToVal] = true
}

// RangeToValCleared returns if the field range_to_val was cleared in this mutation.
func (m *PropertyTypeMutation) RangeToValCleared() bool {
	return m.clearedFields[propertytype.FieldRangeToVal]
}

// ResetRangeToVal reset all changes of the range_to_val field.
func (m *PropertyTypeMutation) ResetRangeToVal() {
	m.range_to_val = nil
	m.addrange_to_val = nil
	delete(m.clearedFields, propertytype.FieldRangeToVal)
}

// SetIsInstanceProperty sets the is_instance_property field.
func (m *PropertyTypeMutation) SetIsInstanceProperty(b bool) {
	m.is_instance_property = &b
}

// IsInstanceProperty returns the is_instance_property value in the mutation.
func (m *PropertyTypeMutation) IsInstanceProperty() (r bool, exists bool) {
	v := m.is_instance_property
	if v == nil {
		return
	}
	return *v, true
}

// ResetIsInstanceProperty reset all changes of the is_instance_property field.
func (m *PropertyTypeMutation) ResetIsInstanceProperty() {
	m.is_instance_property = nil
}

// SetEditable sets the editable field.
func (m *PropertyTypeMutation) SetEditable(b bool) {
	m.editable = &b
}

// Editable returns the editable value in the mutation.
func (m *PropertyTypeMutation) Editable() (r bool, exists bool) {
	v := m.editable
	if v == nil {
		return
	}
	return *v, true
}

// ResetEditable reset all changes of the editable field.
func (m *PropertyTypeMutation) ResetEditable() {
	m.editable = nil
}

// SetMandatory sets the mandatory field.
func (m *PropertyTypeMutation) SetMandatory(b bool) {
	m.mandatory = &b
}

// Mandatory returns the mandatory value in the mutation.
func (m *PropertyTypeMutation) Mandatory() (r bool, exists bool) {
	v := m.mandatory
	if v == nil {
		return
	}
	return *v, true
}

// ResetMandatory reset all changes of the mandatory field.
func (m *PropertyTypeMutation) ResetMandatory() {
	m.mandatory = nil
}

// SetDeleted sets the deleted field.
func (m *PropertyTypeMutation) SetDeleted(b bool) {
	m.deleted = &b
}

// Deleted returns the deleted value in the mutation.
func (m *PropertyTypeMutation) Deleted() (r bool, exists bool) {
	v := m.deleted
	if v == nil {
		return
	}
	return *v, true
}

// ResetDeleted reset all changes of the deleted field.
func (m *PropertyTypeMutation) ResetDeleted() {
	m.deleted = nil
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (m *PropertyTypeMutation) AddPropertyIDs(ids ...int) {
	if m.properties == nil {
		m.properties = make(map[int]struct{})
	}
	for i := range ids {
		m.properties[ids[i]] = struct{}{}
	}
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (m *PropertyTypeMutation) RemovePropertyIDs(ids ...int) {
	if m.removedproperties == nil {
		m.removedproperties = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperties[ids[i]] = struct{}{}
	}
}

// RemovedProperties returns the removed ids of properties.
func (m *PropertyTypeMutation) RemovedPropertiesIDs() (ids []int) {
	for id := range m.removedproperties {
		ids = append(ids, id)
	}
	return
}

// PropertiesIDs returns the properties ids in the mutation.
func (m *PropertyTypeMutation) PropertiesIDs() (ids []int) {
	for id := range m.properties {
		ids = append(ids, id)
	}
	return
}

// ResetProperties reset all changes of the properties edge.
func (m *PropertyTypeMutation) ResetProperties() {
	m.properties = nil
	m.removedproperties = nil
}

// SetLocationTypeID sets the location_type edge to LocationType by id.
func (m *PropertyTypeMutation) SetLocationTypeID(id int) {
	m.location_type = &id
}

// ClearLocationType clears the location_type edge to LocationType.
func (m *PropertyTypeMutation) ClearLocationType() {
	m.clearedlocation_type = true
}

// LocationTypeCleared returns if the edge location_type was cleared.
func (m *PropertyTypeMutation) LocationTypeCleared() bool {
	return m.clearedlocation_type
}

// LocationTypeID returns the location_type id in the mutation.
func (m *PropertyTypeMutation) LocationTypeID() (id int, exists bool) {
	if m.location_type != nil {
		return *m.location_type, true
	}
	return
}

// LocationTypeIDs returns the location_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LocationTypeID instead. It exists only for internal usage by the builders.
func (m *PropertyTypeMutation) LocationTypeIDs() (ids []int) {
	if id := m.location_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLocationType reset all changes of the location_type edge.
func (m *PropertyTypeMutation) ResetLocationType() {
	m.location_type = nil
	m.clearedlocation_type = false
}

// SetEquipmentPortTypeID sets the equipment_port_type edge to EquipmentPortType by id.
func (m *PropertyTypeMutation) SetEquipmentPortTypeID(id int) {
	m.equipment_port_type = &id
}

// ClearEquipmentPortType clears the equipment_port_type edge to EquipmentPortType.
func (m *PropertyTypeMutation) ClearEquipmentPortType() {
	m.clearedequipment_port_type = true
}

// EquipmentPortTypeCleared returns if the edge equipment_port_type was cleared.
func (m *PropertyTypeMutation) EquipmentPortTypeCleared() bool {
	return m.clearedequipment_port_type
}

// EquipmentPortTypeID returns the equipment_port_type id in the mutation.
func (m *PropertyTypeMutation) EquipmentPortTypeID() (id int, exists bool) {
	if m.equipment_port_type != nil {
		return *m.equipment_port_type, true
	}
	return
}

// EquipmentPortTypeIDs returns the equipment_port_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// EquipmentPortTypeID instead. It exists only for internal usage by the builders.
func (m *PropertyTypeMutation) EquipmentPortTypeIDs() (ids []int) {
	if id := m.equipment_port_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEquipmentPortType reset all changes of the equipment_port_type edge.
func (m *PropertyTypeMutation) ResetEquipmentPortType() {
	m.equipment_port_type = nil
	m.clearedequipment_port_type = false
}

// SetLinkEquipmentPortTypeID sets the link_equipment_port_type edge to EquipmentPortType by id.
func (m *PropertyTypeMutation) SetLinkEquipmentPortTypeID(id int) {
	m.link_equipment_port_type = &id
}

// ClearLinkEquipmentPortType clears the link_equipment_port_type edge to EquipmentPortType.
func (m *PropertyTypeMutation) ClearLinkEquipmentPortType() {
	m.clearedlink_equipment_port_type = true
}

// LinkEquipmentPortTypeCleared returns if the edge link_equipment_port_type was cleared.
func (m *PropertyTypeMutation) LinkEquipmentPortTypeCleared() bool {
	return m.clearedlink_equipment_port_type
}

// LinkEquipmentPortTypeID returns the link_equipment_port_type id in the mutation.
func (m *PropertyTypeMutation) LinkEquipmentPortTypeID() (id int, exists bool) {
	if m.link_equipment_port_type != nil {
		return *m.link_equipment_port_type, true
	}
	return
}

// LinkEquipmentPortTypeIDs returns the link_equipment_port_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LinkEquipmentPortTypeID instead. It exists only for internal usage by the builders.
func (m *PropertyTypeMutation) LinkEquipmentPortTypeIDs() (ids []int) {
	if id := m.link_equipment_port_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLinkEquipmentPortType reset all changes of the link_equipment_port_type edge.
func (m *PropertyTypeMutation) ResetLinkEquipmentPortType() {
	m.link_equipment_port_type = nil
	m.clearedlink_equipment_port_type = false
}

// SetEquipmentTypeID sets the equipment_type edge to EquipmentType by id.
func (m *PropertyTypeMutation) SetEquipmentTypeID(id int) {
	m.equipment_type = &id
}

// ClearEquipmentType clears the equipment_type edge to EquipmentType.
func (m *PropertyTypeMutation) ClearEquipmentType() {
	m.clearedequipment_type = true
}

// EquipmentTypeCleared returns if the edge equipment_type was cleared.
func (m *PropertyTypeMutation) EquipmentTypeCleared() bool {
	return m.clearedequipment_type
}

// EquipmentTypeID returns the equipment_type id in the mutation.
func (m *PropertyTypeMutation) EquipmentTypeID() (id int, exists bool) {
	if m.equipment_type != nil {
		return *m.equipment_type, true
	}
	return
}

// EquipmentTypeIDs returns the equipment_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// EquipmentTypeID instead. It exists only for internal usage by the builders.
func (m *PropertyTypeMutation) EquipmentTypeIDs() (ids []int) {
	if id := m.equipment_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetEquipmentType reset all changes of the equipment_type edge.
func (m *PropertyTypeMutation) ResetEquipmentType() {
	m.equipment_type = nil
	m.clearedequipment_type = false
}

// SetServiceTypeID sets the service_type edge to ServiceType by id.
func (m *PropertyTypeMutation) SetServiceTypeID(id int) {
	m.service_type = &id
}

// ClearServiceType clears the service_type edge to ServiceType.
func (m *PropertyTypeMutation) ClearServiceType() {
	m.clearedservice_type = true
}

// ServiceTypeCleared returns if the edge service_type was cleared.
func (m *PropertyTypeMutation) ServiceTypeCleared() bool {
	return m.clearedservice_type
}

// ServiceTypeID returns the service_type id in the mutation.
func (m *PropertyTypeMutation) ServiceTypeID() (id int, exists bool) {
	if m.service_type != nil {
		return *m.service_type, true
	}
	return
}

// ServiceTypeIDs returns the service_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ServiceTypeID instead. It exists only for internal usage by the builders.
func (m *PropertyTypeMutation) ServiceTypeIDs() (ids []int) {
	if id := m.service_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetServiceType reset all changes of the service_type edge.
func (m *PropertyTypeMutation) ResetServiceType() {
	m.service_type = nil
	m.clearedservice_type = false
}

// SetWorkOrderTypeID sets the work_order_type edge to WorkOrderType by id.
func (m *PropertyTypeMutation) SetWorkOrderTypeID(id int) {
	m.work_order_type = &id
}

// ClearWorkOrderType clears the work_order_type edge to WorkOrderType.
func (m *PropertyTypeMutation) ClearWorkOrderType() {
	m.clearedwork_order_type = true
}

// WorkOrderTypeCleared returns if the edge work_order_type was cleared.
func (m *PropertyTypeMutation) WorkOrderTypeCleared() bool {
	return m.clearedwork_order_type
}

// WorkOrderTypeID returns the work_order_type id in the mutation.
func (m *PropertyTypeMutation) WorkOrderTypeID() (id int, exists bool) {
	if m.work_order_type != nil {
		return *m.work_order_type, true
	}
	return
}

// WorkOrderTypeIDs returns the work_order_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// WorkOrderTypeID instead. It exists only for internal usage by the builders.
func (m *PropertyTypeMutation) WorkOrderTypeIDs() (ids []int) {
	if id := m.work_order_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetWorkOrderType reset all changes of the work_order_type edge.
func (m *PropertyTypeMutation) ResetWorkOrderType() {
	m.work_order_type = nil
	m.clearedwork_order_type = false
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (m *PropertyTypeMutation) SetProjectTypeID(id int) {
	m.project_type = &id
}

// ClearProjectType clears the project_type edge to ProjectType.
func (m *PropertyTypeMutation) ClearProjectType() {
	m.clearedproject_type = true
}

// ProjectTypeCleared returns if the edge project_type was cleared.
func (m *PropertyTypeMutation) ProjectTypeCleared() bool {
	return m.clearedproject_type
}

// ProjectTypeID returns the project_type id in the mutation.
func (m *PropertyTypeMutation) ProjectTypeID() (id int, exists bool) {
	if m.project_type != nil {
		return *m.project_type, true
	}
	return
}

// ProjectTypeIDs returns the project_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ProjectTypeID instead. It exists only for internal usage by the builders.
func (m *PropertyTypeMutation) ProjectTypeIDs() (ids []int) {
	if id := m.project_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetProjectType reset all changes of the project_type edge.
func (m *PropertyTypeMutation) ResetProjectType() {
	m.project_type = nil
	m.clearedproject_type = false
}

// Op returns the operation name.
func (m *PropertyTypeMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (PropertyType).
func (m *PropertyTypeMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *PropertyTypeMutation) Fields() []string {
	fields := make([]string, 0, 18)
	if m.create_time != nil {
		fields = append(fields, propertytype.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, propertytype.FieldUpdateTime)
	}
	if m._type != nil {
		fields = append(fields, propertytype.FieldType)
	}
	if m.name != nil {
		fields = append(fields, propertytype.FieldName)
	}
	if m.index != nil {
		fields = append(fields, propertytype.FieldIndex)
	}
	if m.category != nil {
		fields = append(fields, propertytype.FieldCategory)
	}
	if m.int_val != nil {
		fields = append(fields, propertytype.FieldIntVal)
	}
	if m.bool_val != nil {
		fields = append(fields, propertytype.FieldBoolVal)
	}
	if m.float_val != nil {
		fields = append(fields, propertytype.FieldFloatVal)
	}
	if m.latitude_val != nil {
		fields = append(fields, propertytype.FieldLatitudeVal)
	}
	if m.longitude_val != nil {
		fields = append(fields, propertytype.FieldLongitudeVal)
	}
	if m.string_val != nil {
		fields = append(fields, propertytype.FieldStringVal)
	}
	if m.range_from_val != nil {
		fields = append(fields, propertytype.FieldRangeFromVal)
	}
	if m.range_to_val != nil {
		fields = append(fields, propertytype.FieldRangeToVal)
	}
	if m.is_instance_property != nil {
		fields = append(fields, propertytype.FieldIsInstanceProperty)
	}
	if m.editable != nil {
		fields = append(fields, propertytype.FieldEditable)
	}
	if m.mandatory != nil {
		fields = append(fields, propertytype.FieldMandatory)
	}
	if m.deleted != nil {
		fields = append(fields, propertytype.FieldDeleted)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *PropertyTypeMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case propertytype.FieldCreateTime:
		return m.CreateTime()
	case propertytype.FieldUpdateTime:
		return m.UpdateTime()
	case propertytype.FieldType:
		return m.GetType()
	case propertytype.FieldName:
		return m.Name()
	case propertytype.FieldIndex:
		return m.Index()
	case propertytype.FieldCategory:
		return m.Category()
	case propertytype.FieldIntVal:
		return m.IntVal()
	case propertytype.FieldBoolVal:
		return m.BoolVal()
	case propertytype.FieldFloatVal:
		return m.FloatVal()
	case propertytype.FieldLatitudeVal:
		return m.LatitudeVal()
	case propertytype.FieldLongitudeVal:
		return m.LongitudeVal()
	case propertytype.FieldStringVal:
		return m.StringVal()
	case propertytype.FieldRangeFromVal:
		return m.RangeFromVal()
	case propertytype.FieldRangeToVal:
		return m.RangeToVal()
	case propertytype.FieldIsInstanceProperty:
		return m.IsInstanceProperty()
	case propertytype.FieldEditable:
		return m.Editable()
	case propertytype.FieldMandatory:
		return m.Mandatory()
	case propertytype.FieldDeleted:
		return m.Deleted()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *PropertyTypeMutation) SetField(name string, value ent.Value) error {
	switch name {
	case propertytype.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case propertytype.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case propertytype.FieldType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetType(v)
		return nil
	case propertytype.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case propertytype.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIndex(v)
		return nil
	case propertytype.FieldCategory:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCategory(v)
		return nil
	case propertytype.FieldIntVal:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIntVal(v)
		return nil
	case propertytype.FieldBoolVal:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBoolVal(v)
		return nil
	case propertytype.FieldFloatVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFloatVal(v)
		return nil
	case propertytype.FieldLatitudeVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLatitudeVal(v)
		return nil
	case propertytype.FieldLongitudeVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLongitudeVal(v)
		return nil
	case propertytype.FieldStringVal:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStringVal(v)
		return nil
	case propertytype.FieldRangeFromVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRangeFromVal(v)
		return nil
	case propertytype.FieldRangeToVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRangeToVal(v)
		return nil
	case propertytype.FieldIsInstanceProperty:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIsInstanceProperty(v)
		return nil
	case propertytype.FieldEditable:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEditable(v)
		return nil
	case propertytype.FieldMandatory:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMandatory(v)
		return nil
	case propertytype.FieldDeleted:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDeleted(v)
		return nil
	}
	return fmt.Errorf("unknown PropertyType field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *PropertyTypeMutation) AddedFields() []string {
	var fields []string
	if m.addindex != nil {
		fields = append(fields, propertytype.FieldIndex)
	}
	if m.addint_val != nil {
		fields = append(fields, propertytype.FieldIntVal)
	}
	if m.addfloat_val != nil {
		fields = append(fields, propertytype.FieldFloatVal)
	}
	if m.addlatitude_val != nil {
		fields = append(fields, propertytype.FieldLatitudeVal)
	}
	if m.addlongitude_val != nil {
		fields = append(fields, propertytype.FieldLongitudeVal)
	}
	if m.addrange_from_val != nil {
		fields = append(fields, propertytype.FieldRangeFromVal)
	}
	if m.addrange_to_val != nil {
		fields = append(fields, propertytype.FieldRangeToVal)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *PropertyTypeMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case propertytype.FieldIndex:
		return m.AddedIndex()
	case propertytype.FieldIntVal:
		return m.AddedIntVal()
	case propertytype.FieldFloatVal:
		return m.AddedFloatVal()
	case propertytype.FieldLatitudeVal:
		return m.AddedLatitudeVal()
	case propertytype.FieldLongitudeVal:
		return m.AddedLongitudeVal()
	case propertytype.FieldRangeFromVal:
		return m.AddedRangeFromVal()
	case propertytype.FieldRangeToVal:
		return m.AddedRangeToVal()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *PropertyTypeMutation) AddField(name string, value ent.Value) error {
	switch name {
	case propertytype.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIndex(v)
		return nil
	case propertytype.FieldIntVal:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIntVal(v)
		return nil
	case propertytype.FieldFloatVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddFloatVal(v)
		return nil
	case propertytype.FieldLatitudeVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLatitudeVal(v)
		return nil
	case propertytype.FieldLongitudeVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLongitudeVal(v)
		return nil
	case propertytype.FieldRangeFromVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddRangeFromVal(v)
		return nil
	case propertytype.FieldRangeToVal:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddRangeToVal(v)
		return nil
	}
	return fmt.Errorf("unknown PropertyType numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *PropertyTypeMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[propertytype.FieldIndex] {
		fields = append(fields, propertytype.FieldIndex)
	}
	if m.clearedFields[propertytype.FieldCategory] {
		fields = append(fields, propertytype.FieldCategory)
	}
	if m.clearedFields[propertytype.FieldIntVal] {
		fields = append(fields, propertytype.FieldIntVal)
	}
	if m.clearedFields[propertytype.FieldBoolVal] {
		fields = append(fields, propertytype.FieldBoolVal)
	}
	if m.clearedFields[propertytype.FieldFloatVal] {
		fields = append(fields, propertytype.FieldFloatVal)
	}
	if m.clearedFields[propertytype.FieldLatitudeVal] {
		fields = append(fields, propertytype.FieldLatitudeVal)
	}
	if m.clearedFields[propertytype.FieldLongitudeVal] {
		fields = append(fields, propertytype.FieldLongitudeVal)
	}
	if m.clearedFields[propertytype.FieldStringVal] {
		fields = append(fields, propertytype.FieldStringVal)
	}
	if m.clearedFields[propertytype.FieldRangeFromVal] {
		fields = append(fields, propertytype.FieldRangeFromVal)
	}
	if m.clearedFields[propertytype.FieldRangeToVal] {
		fields = append(fields, propertytype.FieldRangeToVal)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *PropertyTypeMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *PropertyTypeMutation) ClearField(name string) error {
	switch name {
	case propertytype.FieldIndex:
		m.ClearIndex()
		return nil
	case propertytype.FieldCategory:
		m.ClearCategory()
		return nil
	case propertytype.FieldIntVal:
		m.ClearIntVal()
		return nil
	case propertytype.FieldBoolVal:
		m.ClearBoolVal()
		return nil
	case propertytype.FieldFloatVal:
		m.ClearFloatVal()
		return nil
	case propertytype.FieldLatitudeVal:
		m.ClearLatitudeVal()
		return nil
	case propertytype.FieldLongitudeVal:
		m.ClearLongitudeVal()
		return nil
	case propertytype.FieldStringVal:
		m.ClearStringVal()
		return nil
	case propertytype.FieldRangeFromVal:
		m.ClearRangeFromVal()
		return nil
	case propertytype.FieldRangeToVal:
		m.ClearRangeToVal()
		return nil
	}
	return fmt.Errorf("unknown PropertyType nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *PropertyTypeMutation) ResetField(name string) error {
	switch name {
	case propertytype.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case propertytype.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case propertytype.FieldType:
		m.ResetType()
		return nil
	case propertytype.FieldName:
		m.ResetName()
		return nil
	case propertytype.FieldIndex:
		m.ResetIndex()
		return nil
	case propertytype.FieldCategory:
		m.ResetCategory()
		return nil
	case propertytype.FieldIntVal:
		m.ResetIntVal()
		return nil
	case propertytype.FieldBoolVal:
		m.ResetBoolVal()
		return nil
	case propertytype.FieldFloatVal:
		m.ResetFloatVal()
		return nil
	case propertytype.FieldLatitudeVal:
		m.ResetLatitudeVal()
		return nil
	case propertytype.FieldLongitudeVal:
		m.ResetLongitudeVal()
		return nil
	case propertytype.FieldStringVal:
		m.ResetStringVal()
		return nil
	case propertytype.FieldRangeFromVal:
		m.ResetRangeFromVal()
		return nil
	case propertytype.FieldRangeToVal:
		m.ResetRangeToVal()
		return nil
	case propertytype.FieldIsInstanceProperty:
		m.ResetIsInstanceProperty()
		return nil
	case propertytype.FieldEditable:
		m.ResetEditable()
		return nil
	case propertytype.FieldMandatory:
		m.ResetMandatory()
		return nil
	case propertytype.FieldDeleted:
		m.ResetDeleted()
		return nil
	}
	return fmt.Errorf("unknown PropertyType field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *PropertyTypeMutation) AddedEdges() []string {
	edges := make([]string, 0, 8)
	if m.properties != nil {
		edges = append(edges, propertytype.EdgeProperties)
	}
	if m.location_type != nil {
		edges = append(edges, propertytype.EdgeLocationType)
	}
	if m.equipment_port_type != nil {
		edges = append(edges, propertytype.EdgeEquipmentPortType)
	}
	if m.link_equipment_port_type != nil {
		edges = append(edges, propertytype.EdgeLinkEquipmentPortType)
	}
	if m.equipment_type != nil {
		edges = append(edges, propertytype.EdgeEquipmentType)
	}
	if m.service_type != nil {
		edges = append(edges, propertytype.EdgeServiceType)
	}
	if m.work_order_type != nil {
		edges = append(edges, propertytype.EdgeWorkOrderType)
	}
	if m.project_type != nil {
		edges = append(edges, propertytype.EdgeProjectType)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *PropertyTypeMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case propertytype.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.properties))
		for id := range m.properties {
			ids = append(ids, id)
		}
		return ids
	case propertytype.EdgeLocationType:
		if id := m.location_type; id != nil {
			return []ent.Value{*id}
		}
	case propertytype.EdgeEquipmentPortType:
		if id := m.equipment_port_type; id != nil {
			return []ent.Value{*id}
		}
	case propertytype.EdgeLinkEquipmentPortType:
		if id := m.link_equipment_port_type; id != nil {
			return []ent.Value{*id}
		}
	case propertytype.EdgeEquipmentType:
		if id := m.equipment_type; id != nil {
			return []ent.Value{*id}
		}
	case propertytype.EdgeServiceType:
		if id := m.service_type; id != nil {
			return []ent.Value{*id}
		}
	case propertytype.EdgeWorkOrderType:
		if id := m.work_order_type; id != nil {
			return []ent.Value{*id}
		}
	case propertytype.EdgeProjectType:
		if id := m.project_type; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *PropertyTypeMutation) RemovedEdges() []string {
	edges := make([]string, 0, 8)
	if m.removedproperties != nil {
		edges = append(edges, propertytype.EdgeProperties)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *PropertyTypeMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case propertytype.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.removedproperties))
		for id := range m.removedproperties {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *PropertyTypeMutation) ClearedEdges() []string {
	edges := make([]string, 0, 8)
	if m.clearedlocation_type {
		edges = append(edges, propertytype.EdgeLocationType)
	}
	if m.clearedequipment_port_type {
		edges = append(edges, propertytype.EdgeEquipmentPortType)
	}
	if m.clearedlink_equipment_port_type {
		edges = append(edges, propertytype.EdgeLinkEquipmentPortType)
	}
	if m.clearedequipment_type {
		edges = append(edges, propertytype.EdgeEquipmentType)
	}
	if m.clearedservice_type {
		edges = append(edges, propertytype.EdgeServiceType)
	}
	if m.clearedwork_order_type {
		edges = append(edges, propertytype.EdgeWorkOrderType)
	}
	if m.clearedproject_type {
		edges = append(edges, propertytype.EdgeProjectType)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *PropertyTypeMutation) EdgeCleared(name string) bool {
	switch name {
	case propertytype.EdgeLocationType:
		return m.clearedlocation_type
	case propertytype.EdgeEquipmentPortType:
		return m.clearedequipment_port_type
	case propertytype.EdgeLinkEquipmentPortType:
		return m.clearedlink_equipment_port_type
	case propertytype.EdgeEquipmentType:
		return m.clearedequipment_type
	case propertytype.EdgeServiceType:
		return m.clearedservice_type
	case propertytype.EdgeWorkOrderType:
		return m.clearedwork_order_type
	case propertytype.EdgeProjectType:
		return m.clearedproject_type
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *PropertyTypeMutation) ClearEdge(name string) error {
	switch name {
	case propertytype.EdgeLocationType:
		m.ClearLocationType()
		return nil
	case propertytype.EdgeEquipmentPortType:
		m.ClearEquipmentPortType()
		return nil
	case propertytype.EdgeLinkEquipmentPortType:
		m.ClearLinkEquipmentPortType()
		return nil
	case propertytype.EdgeEquipmentType:
		m.ClearEquipmentType()
		return nil
	case propertytype.EdgeServiceType:
		m.ClearServiceType()
		return nil
	case propertytype.EdgeWorkOrderType:
		m.ClearWorkOrderType()
		return nil
	case propertytype.EdgeProjectType:
		m.ClearProjectType()
		return nil
	}
	return fmt.Errorf("unknown PropertyType unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *PropertyTypeMutation) ResetEdge(name string) error {
	switch name {
	case propertytype.EdgeProperties:
		m.ResetProperties()
		return nil
	case propertytype.EdgeLocationType:
		m.ResetLocationType()
		return nil
	case propertytype.EdgeEquipmentPortType:
		m.ResetEquipmentPortType()
		return nil
	case propertytype.EdgeLinkEquipmentPortType:
		m.ResetLinkEquipmentPortType()
		return nil
	case propertytype.EdgeEquipmentType:
		m.ResetEquipmentType()
		return nil
	case propertytype.EdgeServiceType:
		m.ResetServiceType()
		return nil
	case propertytype.EdgeWorkOrderType:
		m.ResetWorkOrderType()
		return nil
	case propertytype.EdgeProjectType:
		m.ResetProjectType()
		return nil
	}
	return fmt.Errorf("unknown PropertyType edge %s", name)
}

// ReportFilterMutation represents an operation that mutate the ReportFilters
// nodes in the graph.
type ReportFilterMutation struct {
	config
	op            Op
	typ           string
	id            *int
	create_time   *time.Time
	update_time   *time.Time
	name          *string
	entity        *reportfilter.Entity
	filters       *string
	clearedFields map[string]bool
}

var _ ent.Mutation = (*ReportFilterMutation)(nil)

// newReportFilterMutation creates new mutation for $n.Name.
func newReportFilterMutation(c config, op Op) *ReportFilterMutation {
	return &ReportFilterMutation{
		config:        c,
		op:            op,
		typ:           TypeReportFilter,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ReportFilterMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ReportFilterMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *ReportFilterMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *ReportFilterMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *ReportFilterMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *ReportFilterMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *ReportFilterMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *ReportFilterMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *ReportFilterMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *ReportFilterMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *ReportFilterMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *ReportFilterMutation) ResetName() {
	m.name = nil
}

// SetEntity sets the entity field.
func (m *ReportFilterMutation) SetEntity(r reportfilter.Entity) {
	m.entity = &r
}

// Entity returns the entity value in the mutation.
func (m *ReportFilterMutation) Entity() (r reportfilter.Entity, exists bool) {
	v := m.entity
	if v == nil {
		return
	}
	return *v, true
}

// ResetEntity reset all changes of the entity field.
func (m *ReportFilterMutation) ResetEntity() {
	m.entity = nil
}

// SetFilters sets the filters field.
func (m *ReportFilterMutation) SetFilters(s string) {
	m.filters = &s
}

// Filters returns the filters value in the mutation.
func (m *ReportFilterMutation) Filters() (r string, exists bool) {
	v := m.filters
	if v == nil {
		return
	}
	return *v, true
}

// ResetFilters reset all changes of the filters field.
func (m *ReportFilterMutation) ResetFilters() {
	m.filters = nil
}

// Op returns the operation name.
func (m *ReportFilterMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (ReportFilter).
func (m *ReportFilterMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *ReportFilterMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.create_time != nil {
		fields = append(fields, reportfilter.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, reportfilter.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, reportfilter.FieldName)
	}
	if m.entity != nil {
		fields = append(fields, reportfilter.FieldEntity)
	}
	if m.filters != nil {
		fields = append(fields, reportfilter.FieldFilters)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *ReportFilterMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case reportfilter.FieldCreateTime:
		return m.CreateTime()
	case reportfilter.FieldUpdateTime:
		return m.UpdateTime()
	case reportfilter.FieldName:
		return m.Name()
	case reportfilter.FieldEntity:
		return m.Entity()
	case reportfilter.FieldFilters:
		return m.Filters()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ReportFilterMutation) SetField(name string, value ent.Value) error {
	switch name {
	case reportfilter.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case reportfilter.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case reportfilter.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case reportfilter.FieldEntity:
		v, ok := value.(reportfilter.Entity)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEntity(v)
		return nil
	case reportfilter.FieldFilters:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFilters(v)
		return nil
	}
	return fmt.Errorf("unknown ReportFilter field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *ReportFilterMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *ReportFilterMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ReportFilterMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ReportFilter numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *ReportFilterMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *ReportFilterMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *ReportFilterMutation) ClearField(name string) error {
	return fmt.Errorf("unknown ReportFilter nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *ReportFilterMutation) ResetField(name string) error {
	switch name {
	case reportfilter.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case reportfilter.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case reportfilter.FieldName:
		m.ResetName()
		return nil
	case reportfilter.FieldEntity:
		m.ResetEntity()
		return nil
	case reportfilter.FieldFilters:
		m.ResetFilters()
		return nil
	}
	return fmt.Errorf("unknown ReportFilter field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *ReportFilterMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *ReportFilterMutation) AddedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *ReportFilterMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *ReportFilterMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *ReportFilterMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *ReportFilterMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *ReportFilterMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown ReportFilter unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *ReportFilterMutation) ResetEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown ReportFilter edge %s", name)
}

// ServiceMutation represents an operation that mutate the Services
// nodes in the graph.
type ServiceMutation struct {
	config
	op                Op
	typ               string
	id                *int
	create_time       *time.Time
	update_time       *time.Time
	name              *string
	external_id       *string
	status            *string
	clearedFields     map[string]bool
	_type             *int
	cleared_type      bool
	downstream        map[int]struct{}
	removeddownstream map[int]struct{}
	upstream          map[int]struct{}
	removedupstream   map[int]struct{}
	properties        map[int]struct{}
	removedproperties map[int]struct{}
	links             map[int]struct{}
	removedlinks      map[int]struct{}
	customer          map[int]struct{}
	removedcustomer   map[int]struct{}
	endpoints         map[int]struct{}
	removedendpoints  map[int]struct{}
}

var _ ent.Mutation = (*ServiceMutation)(nil)

// newServiceMutation creates new mutation for $n.Name.
func newServiceMutation(c config, op Op) *ServiceMutation {
	return &ServiceMutation{
		config:        c,
		op:            op,
		typ:           TypeService,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ServiceMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ServiceMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *ServiceMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *ServiceMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *ServiceMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *ServiceMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *ServiceMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *ServiceMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *ServiceMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *ServiceMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *ServiceMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *ServiceMutation) ResetName() {
	m.name = nil
}

// SetExternalID sets the external_id field.
func (m *ServiceMutation) SetExternalID(s string) {
	m.external_id = &s
}

// ExternalID returns the external_id value in the mutation.
func (m *ServiceMutation) ExternalID() (r string, exists bool) {
	v := m.external_id
	if v == nil {
		return
	}
	return *v, true
}

// ClearExternalID clears the value of external_id.
func (m *ServiceMutation) ClearExternalID() {
	m.external_id = nil
	m.clearedFields[service.FieldExternalID] = true
}

// ExternalIDCleared returns if the field external_id was cleared in this mutation.
func (m *ServiceMutation) ExternalIDCleared() bool {
	return m.clearedFields[service.FieldExternalID]
}

// ResetExternalID reset all changes of the external_id field.
func (m *ServiceMutation) ResetExternalID() {
	m.external_id = nil
	delete(m.clearedFields, service.FieldExternalID)
}

// SetStatus sets the status field.
func (m *ServiceMutation) SetStatus(s string) {
	m.status = &s
}

// Status returns the status value in the mutation.
func (m *ServiceMutation) Status() (r string, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// ResetStatus reset all changes of the status field.
func (m *ServiceMutation) ResetStatus() {
	m.status = nil
}

// SetTypeID sets the type edge to ServiceType by id.
func (m *ServiceMutation) SetTypeID(id int) {
	m._type = &id
}

// ClearType clears the type edge to ServiceType.
func (m *ServiceMutation) ClearType() {
	m.cleared_type = true
}

// TypeCleared returns if the edge type was cleared.
func (m *ServiceMutation) TypeCleared() bool {
	return m.cleared_type
}

// TypeID returns the type id in the mutation.
func (m *ServiceMutation) TypeID() (id int, exists bool) {
	if m._type != nil {
		return *m._type, true
	}
	return
}

// TypeIDs returns the type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// TypeID instead. It exists only for internal usage by the builders.
func (m *ServiceMutation) TypeIDs() (ids []int) {
	if id := m._type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetType reset all changes of the type edge.
func (m *ServiceMutation) ResetType() {
	m._type = nil
	m.cleared_type = false
}

// AddDownstreamIDs adds the downstream edge to Service by ids.
func (m *ServiceMutation) AddDownstreamIDs(ids ...int) {
	if m.downstream == nil {
		m.downstream = make(map[int]struct{})
	}
	for i := range ids {
		m.downstream[ids[i]] = struct{}{}
	}
}

// RemoveDownstreamIDs removes the downstream edge to Service by ids.
func (m *ServiceMutation) RemoveDownstreamIDs(ids ...int) {
	if m.removeddownstream == nil {
		m.removeddownstream = make(map[int]struct{})
	}
	for i := range ids {
		m.removeddownstream[ids[i]] = struct{}{}
	}
}

// RemovedDownstream returns the removed ids of downstream.
func (m *ServiceMutation) RemovedDownstreamIDs() (ids []int) {
	for id := range m.removeddownstream {
		ids = append(ids, id)
	}
	return
}

// DownstreamIDs returns the downstream ids in the mutation.
func (m *ServiceMutation) DownstreamIDs() (ids []int) {
	for id := range m.downstream {
		ids = append(ids, id)
	}
	return
}

// ResetDownstream reset all changes of the downstream edge.
func (m *ServiceMutation) ResetDownstream() {
	m.downstream = nil
	m.removeddownstream = nil
}

// AddUpstreamIDs adds the upstream edge to Service by ids.
func (m *ServiceMutation) AddUpstreamIDs(ids ...int) {
	if m.upstream == nil {
		m.upstream = make(map[int]struct{})
	}
	for i := range ids {
		m.upstream[ids[i]] = struct{}{}
	}
}

// RemoveUpstreamIDs removes the upstream edge to Service by ids.
func (m *ServiceMutation) RemoveUpstreamIDs(ids ...int) {
	if m.removedupstream == nil {
		m.removedupstream = make(map[int]struct{})
	}
	for i := range ids {
		m.removedupstream[ids[i]] = struct{}{}
	}
}

// RemovedUpstream returns the removed ids of upstream.
func (m *ServiceMutation) RemovedUpstreamIDs() (ids []int) {
	for id := range m.removedupstream {
		ids = append(ids, id)
	}
	return
}

// UpstreamIDs returns the upstream ids in the mutation.
func (m *ServiceMutation) UpstreamIDs() (ids []int) {
	for id := range m.upstream {
		ids = append(ids, id)
	}
	return
}

// ResetUpstream reset all changes of the upstream edge.
func (m *ServiceMutation) ResetUpstream() {
	m.upstream = nil
	m.removedupstream = nil
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (m *ServiceMutation) AddPropertyIDs(ids ...int) {
	if m.properties == nil {
		m.properties = make(map[int]struct{})
	}
	for i := range ids {
		m.properties[ids[i]] = struct{}{}
	}
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (m *ServiceMutation) RemovePropertyIDs(ids ...int) {
	if m.removedproperties == nil {
		m.removedproperties = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperties[ids[i]] = struct{}{}
	}
}

// RemovedProperties returns the removed ids of properties.
func (m *ServiceMutation) RemovedPropertiesIDs() (ids []int) {
	for id := range m.removedproperties {
		ids = append(ids, id)
	}
	return
}

// PropertiesIDs returns the properties ids in the mutation.
func (m *ServiceMutation) PropertiesIDs() (ids []int) {
	for id := range m.properties {
		ids = append(ids, id)
	}
	return
}

// ResetProperties reset all changes of the properties edge.
func (m *ServiceMutation) ResetProperties() {
	m.properties = nil
	m.removedproperties = nil
}

// AddLinkIDs adds the links edge to Link by ids.
func (m *ServiceMutation) AddLinkIDs(ids ...int) {
	if m.links == nil {
		m.links = make(map[int]struct{})
	}
	for i := range ids {
		m.links[ids[i]] = struct{}{}
	}
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (m *ServiceMutation) RemoveLinkIDs(ids ...int) {
	if m.removedlinks == nil {
		m.removedlinks = make(map[int]struct{})
	}
	for i := range ids {
		m.removedlinks[ids[i]] = struct{}{}
	}
}

// RemovedLinks returns the removed ids of links.
func (m *ServiceMutation) RemovedLinksIDs() (ids []int) {
	for id := range m.removedlinks {
		ids = append(ids, id)
	}
	return
}

// LinksIDs returns the links ids in the mutation.
func (m *ServiceMutation) LinksIDs() (ids []int) {
	for id := range m.links {
		ids = append(ids, id)
	}
	return
}

// ResetLinks reset all changes of the links edge.
func (m *ServiceMutation) ResetLinks() {
	m.links = nil
	m.removedlinks = nil
}

// AddCustomerIDs adds the customer edge to Customer by ids.
func (m *ServiceMutation) AddCustomerIDs(ids ...int) {
	if m.customer == nil {
		m.customer = make(map[int]struct{})
	}
	for i := range ids {
		m.customer[ids[i]] = struct{}{}
	}
}

// RemoveCustomerIDs removes the customer edge to Customer by ids.
func (m *ServiceMutation) RemoveCustomerIDs(ids ...int) {
	if m.removedcustomer == nil {
		m.removedcustomer = make(map[int]struct{})
	}
	for i := range ids {
		m.removedcustomer[ids[i]] = struct{}{}
	}
}

// RemovedCustomer returns the removed ids of customer.
func (m *ServiceMutation) RemovedCustomerIDs() (ids []int) {
	for id := range m.removedcustomer {
		ids = append(ids, id)
	}
	return
}

// CustomerIDs returns the customer ids in the mutation.
func (m *ServiceMutation) CustomerIDs() (ids []int) {
	for id := range m.customer {
		ids = append(ids, id)
	}
	return
}

// ResetCustomer reset all changes of the customer edge.
func (m *ServiceMutation) ResetCustomer() {
	m.customer = nil
	m.removedcustomer = nil
}

// AddEndpointIDs adds the endpoints edge to ServiceEndpoint by ids.
func (m *ServiceMutation) AddEndpointIDs(ids ...int) {
	if m.endpoints == nil {
		m.endpoints = make(map[int]struct{})
	}
	for i := range ids {
		m.endpoints[ids[i]] = struct{}{}
	}
}

// RemoveEndpointIDs removes the endpoints edge to ServiceEndpoint by ids.
func (m *ServiceMutation) RemoveEndpointIDs(ids ...int) {
	if m.removedendpoints == nil {
		m.removedendpoints = make(map[int]struct{})
	}
	for i := range ids {
		m.removedendpoints[ids[i]] = struct{}{}
	}
}

// RemovedEndpoints returns the removed ids of endpoints.
func (m *ServiceMutation) RemovedEndpointsIDs() (ids []int) {
	for id := range m.removedendpoints {
		ids = append(ids, id)
	}
	return
}

// EndpointsIDs returns the endpoints ids in the mutation.
func (m *ServiceMutation) EndpointsIDs() (ids []int) {
	for id := range m.endpoints {
		ids = append(ids, id)
	}
	return
}

// ResetEndpoints reset all changes of the endpoints edge.
func (m *ServiceMutation) ResetEndpoints() {
	m.endpoints = nil
	m.removedendpoints = nil
}

// Op returns the operation name.
func (m *ServiceMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Service).
func (m *ServiceMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *ServiceMutation) Fields() []string {
	fields := make([]string, 0, 5)
	if m.create_time != nil {
		fields = append(fields, service.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, service.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, service.FieldName)
	}
	if m.external_id != nil {
		fields = append(fields, service.FieldExternalID)
	}
	if m.status != nil {
		fields = append(fields, service.FieldStatus)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *ServiceMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case service.FieldCreateTime:
		return m.CreateTime()
	case service.FieldUpdateTime:
		return m.UpdateTime()
	case service.FieldName:
		return m.Name()
	case service.FieldExternalID:
		return m.ExternalID()
	case service.FieldStatus:
		return m.Status()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ServiceMutation) SetField(name string, value ent.Value) error {
	switch name {
	case service.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case service.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case service.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case service.FieldExternalID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetExternalID(v)
		return nil
	case service.FieldStatus:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	}
	return fmt.Errorf("unknown Service field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *ServiceMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *ServiceMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ServiceMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Service numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *ServiceMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[service.FieldExternalID] {
		fields = append(fields, service.FieldExternalID)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *ServiceMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *ServiceMutation) ClearField(name string) error {
	switch name {
	case service.FieldExternalID:
		m.ClearExternalID()
		return nil
	}
	return fmt.Errorf("unknown Service nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *ServiceMutation) ResetField(name string) error {
	switch name {
	case service.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case service.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case service.FieldName:
		m.ResetName()
		return nil
	case service.FieldExternalID:
		m.ResetExternalID()
		return nil
	case service.FieldStatus:
		m.ResetStatus()
		return nil
	}
	return fmt.Errorf("unknown Service field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *ServiceMutation) AddedEdges() []string {
	edges := make([]string, 0, 7)
	if m._type != nil {
		edges = append(edges, service.EdgeType)
	}
	if m.downstream != nil {
		edges = append(edges, service.EdgeDownstream)
	}
	if m.upstream != nil {
		edges = append(edges, service.EdgeUpstream)
	}
	if m.properties != nil {
		edges = append(edges, service.EdgeProperties)
	}
	if m.links != nil {
		edges = append(edges, service.EdgeLinks)
	}
	if m.customer != nil {
		edges = append(edges, service.EdgeCustomer)
	}
	if m.endpoints != nil {
		edges = append(edges, service.EdgeEndpoints)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *ServiceMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case service.EdgeType:
		if id := m._type; id != nil {
			return []ent.Value{*id}
		}
	case service.EdgeDownstream:
		ids := make([]ent.Value, 0, len(m.downstream))
		for id := range m.downstream {
			ids = append(ids, id)
		}
		return ids
	case service.EdgeUpstream:
		ids := make([]ent.Value, 0, len(m.upstream))
		for id := range m.upstream {
			ids = append(ids, id)
		}
		return ids
	case service.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.properties))
		for id := range m.properties {
			ids = append(ids, id)
		}
		return ids
	case service.EdgeLinks:
		ids := make([]ent.Value, 0, len(m.links))
		for id := range m.links {
			ids = append(ids, id)
		}
		return ids
	case service.EdgeCustomer:
		ids := make([]ent.Value, 0, len(m.customer))
		for id := range m.customer {
			ids = append(ids, id)
		}
		return ids
	case service.EdgeEndpoints:
		ids := make([]ent.Value, 0, len(m.endpoints))
		for id := range m.endpoints {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *ServiceMutation) RemovedEdges() []string {
	edges := make([]string, 0, 7)
	if m.removeddownstream != nil {
		edges = append(edges, service.EdgeDownstream)
	}
	if m.removedupstream != nil {
		edges = append(edges, service.EdgeUpstream)
	}
	if m.removedproperties != nil {
		edges = append(edges, service.EdgeProperties)
	}
	if m.removedlinks != nil {
		edges = append(edges, service.EdgeLinks)
	}
	if m.removedcustomer != nil {
		edges = append(edges, service.EdgeCustomer)
	}
	if m.removedendpoints != nil {
		edges = append(edges, service.EdgeEndpoints)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *ServiceMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case service.EdgeDownstream:
		ids := make([]ent.Value, 0, len(m.removeddownstream))
		for id := range m.removeddownstream {
			ids = append(ids, id)
		}
		return ids
	case service.EdgeUpstream:
		ids := make([]ent.Value, 0, len(m.removedupstream))
		for id := range m.removedupstream {
			ids = append(ids, id)
		}
		return ids
	case service.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.removedproperties))
		for id := range m.removedproperties {
			ids = append(ids, id)
		}
		return ids
	case service.EdgeLinks:
		ids := make([]ent.Value, 0, len(m.removedlinks))
		for id := range m.removedlinks {
			ids = append(ids, id)
		}
		return ids
	case service.EdgeCustomer:
		ids := make([]ent.Value, 0, len(m.removedcustomer))
		for id := range m.removedcustomer {
			ids = append(ids, id)
		}
		return ids
	case service.EdgeEndpoints:
		ids := make([]ent.Value, 0, len(m.removedendpoints))
		for id := range m.removedendpoints {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *ServiceMutation) ClearedEdges() []string {
	edges := make([]string, 0, 7)
	if m.cleared_type {
		edges = append(edges, service.EdgeType)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *ServiceMutation) EdgeCleared(name string) bool {
	switch name {
	case service.EdgeType:
		return m.cleared_type
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *ServiceMutation) ClearEdge(name string) error {
	switch name {
	case service.EdgeType:
		m.ClearType()
		return nil
	}
	return fmt.Errorf("unknown Service unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *ServiceMutation) ResetEdge(name string) error {
	switch name {
	case service.EdgeType:
		m.ResetType()
		return nil
	case service.EdgeDownstream:
		m.ResetDownstream()
		return nil
	case service.EdgeUpstream:
		m.ResetUpstream()
		return nil
	case service.EdgeProperties:
		m.ResetProperties()
		return nil
	case service.EdgeLinks:
		m.ResetLinks()
		return nil
	case service.EdgeCustomer:
		m.ResetCustomer()
		return nil
	case service.EdgeEndpoints:
		m.ResetEndpoints()
		return nil
	}
	return fmt.Errorf("unknown Service edge %s", name)
}

// ServiceEndpointMutation represents an operation that mutate the ServiceEndpoints
// nodes in the graph.
type ServiceEndpointMutation struct {
	config
	op             Op
	typ            string
	id             *int
	create_time    *time.Time
	update_time    *time.Time
	role           *string
	clearedFields  map[string]bool
	port           *int
	clearedport    bool
	service        *int
	clearedservice bool
}

var _ ent.Mutation = (*ServiceEndpointMutation)(nil)

// newServiceEndpointMutation creates new mutation for $n.Name.
func newServiceEndpointMutation(c config, op Op) *ServiceEndpointMutation {
	return &ServiceEndpointMutation{
		config:        c,
		op:            op,
		typ:           TypeServiceEndpoint,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ServiceEndpointMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ServiceEndpointMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *ServiceEndpointMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *ServiceEndpointMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *ServiceEndpointMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *ServiceEndpointMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *ServiceEndpointMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *ServiceEndpointMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *ServiceEndpointMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetRole sets the role field.
func (m *ServiceEndpointMutation) SetRole(s string) {
	m.role = &s
}

// Role returns the role value in the mutation.
func (m *ServiceEndpointMutation) Role() (r string, exists bool) {
	v := m.role
	if v == nil {
		return
	}
	return *v, true
}

// ResetRole reset all changes of the role field.
func (m *ServiceEndpointMutation) ResetRole() {
	m.role = nil
}

// SetPortID sets the port edge to EquipmentPort by id.
func (m *ServiceEndpointMutation) SetPortID(id int) {
	m.port = &id
}

// ClearPort clears the port edge to EquipmentPort.
func (m *ServiceEndpointMutation) ClearPort() {
	m.clearedport = true
}

// PortCleared returns if the edge port was cleared.
func (m *ServiceEndpointMutation) PortCleared() bool {
	return m.clearedport
}

// PortID returns the port id in the mutation.
func (m *ServiceEndpointMutation) PortID() (id int, exists bool) {
	if m.port != nil {
		return *m.port, true
	}
	return
}

// PortIDs returns the port ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// PortID instead. It exists only for internal usage by the builders.
func (m *ServiceEndpointMutation) PortIDs() (ids []int) {
	if id := m.port; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetPort reset all changes of the port edge.
func (m *ServiceEndpointMutation) ResetPort() {
	m.port = nil
	m.clearedport = false
}

// SetServiceID sets the service edge to Service by id.
func (m *ServiceEndpointMutation) SetServiceID(id int) {
	m.service = &id
}

// ClearService clears the service edge to Service.
func (m *ServiceEndpointMutation) ClearService() {
	m.clearedservice = true
}

// ServiceCleared returns if the edge service was cleared.
func (m *ServiceEndpointMutation) ServiceCleared() bool {
	return m.clearedservice
}

// ServiceID returns the service id in the mutation.
func (m *ServiceEndpointMutation) ServiceID() (id int, exists bool) {
	if m.service != nil {
		return *m.service, true
	}
	return
}

// ServiceIDs returns the service ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ServiceID instead. It exists only for internal usage by the builders.
func (m *ServiceEndpointMutation) ServiceIDs() (ids []int) {
	if id := m.service; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetService reset all changes of the service edge.
func (m *ServiceEndpointMutation) ResetService() {
	m.service = nil
	m.clearedservice = false
}

// Op returns the operation name.
func (m *ServiceEndpointMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (ServiceEndpoint).
func (m *ServiceEndpointMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *ServiceEndpointMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.create_time != nil {
		fields = append(fields, serviceendpoint.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, serviceendpoint.FieldUpdateTime)
	}
	if m.role != nil {
		fields = append(fields, serviceendpoint.FieldRole)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *ServiceEndpointMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case serviceendpoint.FieldCreateTime:
		return m.CreateTime()
	case serviceendpoint.FieldUpdateTime:
		return m.UpdateTime()
	case serviceendpoint.FieldRole:
		return m.Role()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ServiceEndpointMutation) SetField(name string, value ent.Value) error {
	switch name {
	case serviceendpoint.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case serviceendpoint.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case serviceendpoint.FieldRole:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRole(v)
		return nil
	}
	return fmt.Errorf("unknown ServiceEndpoint field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *ServiceEndpointMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *ServiceEndpointMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ServiceEndpointMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ServiceEndpoint numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *ServiceEndpointMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *ServiceEndpointMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *ServiceEndpointMutation) ClearField(name string) error {
	return fmt.Errorf("unknown ServiceEndpoint nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *ServiceEndpointMutation) ResetField(name string) error {
	switch name {
	case serviceendpoint.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case serviceendpoint.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case serviceendpoint.FieldRole:
		m.ResetRole()
		return nil
	}
	return fmt.Errorf("unknown ServiceEndpoint field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *ServiceEndpointMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.port != nil {
		edges = append(edges, serviceendpoint.EdgePort)
	}
	if m.service != nil {
		edges = append(edges, serviceendpoint.EdgeService)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *ServiceEndpointMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case serviceendpoint.EdgePort:
		if id := m.port; id != nil {
			return []ent.Value{*id}
		}
	case serviceendpoint.EdgeService:
		if id := m.service; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *ServiceEndpointMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *ServiceEndpointMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *ServiceEndpointMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedport {
		edges = append(edges, serviceendpoint.EdgePort)
	}
	if m.clearedservice {
		edges = append(edges, serviceendpoint.EdgeService)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *ServiceEndpointMutation) EdgeCleared(name string) bool {
	switch name {
	case serviceendpoint.EdgePort:
		return m.clearedport
	case serviceendpoint.EdgeService:
		return m.clearedservice
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *ServiceEndpointMutation) ClearEdge(name string) error {
	switch name {
	case serviceendpoint.EdgePort:
		m.ClearPort()
		return nil
	case serviceendpoint.EdgeService:
		m.ClearService()
		return nil
	}
	return fmt.Errorf("unknown ServiceEndpoint unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *ServiceEndpointMutation) ResetEdge(name string) error {
	switch name {
	case serviceendpoint.EdgePort:
		m.ResetPort()
		return nil
	case serviceendpoint.EdgeService:
		m.ResetService()
		return nil
	}
	return fmt.Errorf("unknown ServiceEndpoint edge %s", name)
}

// ServiceTypeMutation represents an operation that mutate the ServiceTypes
// nodes in the graph.
type ServiceTypeMutation struct {
	config
	op                    Op
	typ                   string
	id                    *int
	create_time           *time.Time
	update_time           *time.Time
	name                  *string
	has_customer          *bool
	clearedFields         map[string]bool
	services              map[int]struct{}
	removedservices       map[int]struct{}
	property_types        map[int]struct{}
	removedproperty_types map[int]struct{}
}

var _ ent.Mutation = (*ServiceTypeMutation)(nil)

// newServiceTypeMutation creates new mutation for $n.Name.
func newServiceTypeMutation(c config, op Op) *ServiceTypeMutation {
	return &ServiceTypeMutation{
		config:        c,
		op:            op,
		typ:           TypeServiceType,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m ServiceTypeMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m ServiceTypeMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *ServiceTypeMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *ServiceTypeMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *ServiceTypeMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *ServiceTypeMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *ServiceTypeMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *ServiceTypeMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *ServiceTypeMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *ServiceTypeMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *ServiceTypeMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *ServiceTypeMutation) ResetName() {
	m.name = nil
}

// SetHasCustomer sets the has_customer field.
func (m *ServiceTypeMutation) SetHasCustomer(b bool) {
	m.has_customer = &b
}

// HasCustomer returns the has_customer value in the mutation.
func (m *ServiceTypeMutation) HasCustomer() (r bool, exists bool) {
	v := m.has_customer
	if v == nil {
		return
	}
	return *v, true
}

// ResetHasCustomer reset all changes of the has_customer field.
func (m *ServiceTypeMutation) ResetHasCustomer() {
	m.has_customer = nil
}

// AddServiceIDs adds the services edge to Service by ids.
func (m *ServiceTypeMutation) AddServiceIDs(ids ...int) {
	if m.services == nil {
		m.services = make(map[int]struct{})
	}
	for i := range ids {
		m.services[ids[i]] = struct{}{}
	}
}

// RemoveServiceIDs removes the services edge to Service by ids.
func (m *ServiceTypeMutation) RemoveServiceIDs(ids ...int) {
	if m.removedservices == nil {
		m.removedservices = make(map[int]struct{})
	}
	for i := range ids {
		m.removedservices[ids[i]] = struct{}{}
	}
}

// RemovedServices returns the removed ids of services.
func (m *ServiceTypeMutation) RemovedServicesIDs() (ids []int) {
	for id := range m.removedservices {
		ids = append(ids, id)
	}
	return
}

// ServicesIDs returns the services ids in the mutation.
func (m *ServiceTypeMutation) ServicesIDs() (ids []int) {
	for id := range m.services {
		ids = append(ids, id)
	}
	return
}

// ResetServices reset all changes of the services edge.
func (m *ServiceTypeMutation) ResetServices() {
	m.services = nil
	m.removedservices = nil
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (m *ServiceTypeMutation) AddPropertyTypeIDs(ids ...int) {
	if m.property_types == nil {
		m.property_types = make(map[int]struct{})
	}
	for i := range ids {
		m.property_types[ids[i]] = struct{}{}
	}
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (m *ServiceTypeMutation) RemovePropertyTypeIDs(ids ...int) {
	if m.removedproperty_types == nil {
		m.removedproperty_types = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperty_types[ids[i]] = struct{}{}
	}
}

// RemovedPropertyTypes returns the removed ids of property_types.
func (m *ServiceTypeMutation) RemovedPropertyTypesIDs() (ids []int) {
	for id := range m.removedproperty_types {
		ids = append(ids, id)
	}
	return
}

// PropertyTypesIDs returns the property_types ids in the mutation.
func (m *ServiceTypeMutation) PropertyTypesIDs() (ids []int) {
	for id := range m.property_types {
		ids = append(ids, id)
	}
	return
}

// ResetPropertyTypes reset all changes of the property_types edge.
func (m *ServiceTypeMutation) ResetPropertyTypes() {
	m.property_types = nil
	m.removedproperty_types = nil
}

// Op returns the operation name.
func (m *ServiceTypeMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (ServiceType).
func (m *ServiceTypeMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *ServiceTypeMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.create_time != nil {
		fields = append(fields, servicetype.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, servicetype.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, servicetype.FieldName)
	}
	if m.has_customer != nil {
		fields = append(fields, servicetype.FieldHasCustomer)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *ServiceTypeMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case servicetype.FieldCreateTime:
		return m.CreateTime()
	case servicetype.FieldUpdateTime:
		return m.UpdateTime()
	case servicetype.FieldName:
		return m.Name()
	case servicetype.FieldHasCustomer:
		return m.HasCustomer()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ServiceTypeMutation) SetField(name string, value ent.Value) error {
	switch name {
	case servicetype.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case servicetype.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case servicetype.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case servicetype.FieldHasCustomer:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetHasCustomer(v)
		return nil
	}
	return fmt.Errorf("unknown ServiceType field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *ServiceTypeMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *ServiceTypeMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *ServiceTypeMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown ServiceType numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *ServiceTypeMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *ServiceTypeMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *ServiceTypeMutation) ClearField(name string) error {
	return fmt.Errorf("unknown ServiceType nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *ServiceTypeMutation) ResetField(name string) error {
	switch name {
	case servicetype.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case servicetype.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case servicetype.FieldName:
		m.ResetName()
		return nil
	case servicetype.FieldHasCustomer:
		m.ResetHasCustomer()
		return nil
	}
	return fmt.Errorf("unknown ServiceType field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *ServiceTypeMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.services != nil {
		edges = append(edges, servicetype.EdgeServices)
	}
	if m.property_types != nil {
		edges = append(edges, servicetype.EdgePropertyTypes)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *ServiceTypeMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case servicetype.EdgeServices:
		ids := make([]ent.Value, 0, len(m.services))
		for id := range m.services {
			ids = append(ids, id)
		}
		return ids
	case servicetype.EdgePropertyTypes:
		ids := make([]ent.Value, 0, len(m.property_types))
		for id := range m.property_types {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *ServiceTypeMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	if m.removedservices != nil {
		edges = append(edges, servicetype.EdgeServices)
	}
	if m.removedproperty_types != nil {
		edges = append(edges, servicetype.EdgePropertyTypes)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *ServiceTypeMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case servicetype.EdgeServices:
		ids := make([]ent.Value, 0, len(m.removedservices))
		for id := range m.removedservices {
			ids = append(ids, id)
		}
		return ids
	case servicetype.EdgePropertyTypes:
		ids := make([]ent.Value, 0, len(m.removedproperty_types))
		for id := range m.removedproperty_types {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *ServiceTypeMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *ServiceTypeMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *ServiceTypeMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown ServiceType unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *ServiceTypeMutation) ResetEdge(name string) error {
	switch name {
	case servicetype.EdgeServices:
		m.ResetServices()
		return nil
	case servicetype.EdgePropertyTypes:
		m.ResetPropertyTypes()
		return nil
	}
	return fmt.Errorf("unknown ServiceType edge %s", name)
}

// SurveyMutation represents an operation that mutate the Surveys
// nodes in the graph.
type SurveyMutation struct {
	config
	op                   Op
	typ                  string
	id                   *int
	create_time          *time.Time
	update_time          *time.Time
	name                 *string
	owner_name           *string
	creation_timestamp   *time.Time
	completion_timestamp *time.Time
	clearedFields        map[string]bool
	location             *int
	clearedlocation      bool
	source_file          *int
	clearedsource_file   bool
	questions            map[int]struct{}
	removedquestions     map[int]struct{}
}

var _ ent.Mutation = (*SurveyMutation)(nil)

// newSurveyMutation creates new mutation for $n.Name.
func newSurveyMutation(c config, op Op) *SurveyMutation {
	return &SurveyMutation{
		config:        c,
		op:            op,
		typ:           TypeSurvey,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SurveyMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SurveyMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *SurveyMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *SurveyMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *SurveyMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *SurveyMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *SurveyMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *SurveyMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *SurveyMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *SurveyMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *SurveyMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *SurveyMutation) ResetName() {
	m.name = nil
}

// SetOwnerName sets the owner_name field.
func (m *SurveyMutation) SetOwnerName(s string) {
	m.owner_name = &s
}

// OwnerName returns the owner_name value in the mutation.
func (m *SurveyMutation) OwnerName() (r string, exists bool) {
	v := m.owner_name
	if v == nil {
		return
	}
	return *v, true
}

// ClearOwnerName clears the value of owner_name.
func (m *SurveyMutation) ClearOwnerName() {
	m.owner_name = nil
	m.clearedFields[survey.FieldOwnerName] = true
}

// OwnerNameCleared returns if the field owner_name was cleared in this mutation.
func (m *SurveyMutation) OwnerNameCleared() bool {
	return m.clearedFields[survey.FieldOwnerName]
}

// ResetOwnerName reset all changes of the owner_name field.
func (m *SurveyMutation) ResetOwnerName() {
	m.owner_name = nil
	delete(m.clearedFields, survey.FieldOwnerName)
}

// SetCreationTimestamp sets the creation_timestamp field.
func (m *SurveyMutation) SetCreationTimestamp(t time.Time) {
	m.creation_timestamp = &t
}

// CreationTimestamp returns the creation_timestamp value in the mutation.
func (m *SurveyMutation) CreationTimestamp() (r time.Time, exists bool) {
	v := m.creation_timestamp
	if v == nil {
		return
	}
	return *v, true
}

// ClearCreationTimestamp clears the value of creation_timestamp.
func (m *SurveyMutation) ClearCreationTimestamp() {
	m.creation_timestamp = nil
	m.clearedFields[survey.FieldCreationTimestamp] = true
}

// CreationTimestampCleared returns if the field creation_timestamp was cleared in this mutation.
func (m *SurveyMutation) CreationTimestampCleared() bool {
	return m.clearedFields[survey.FieldCreationTimestamp]
}

// ResetCreationTimestamp reset all changes of the creation_timestamp field.
func (m *SurveyMutation) ResetCreationTimestamp() {
	m.creation_timestamp = nil
	delete(m.clearedFields, survey.FieldCreationTimestamp)
}

// SetCompletionTimestamp sets the completion_timestamp field.
func (m *SurveyMutation) SetCompletionTimestamp(t time.Time) {
	m.completion_timestamp = &t
}

// CompletionTimestamp returns the completion_timestamp value in the mutation.
func (m *SurveyMutation) CompletionTimestamp() (r time.Time, exists bool) {
	v := m.completion_timestamp
	if v == nil {
		return
	}
	return *v, true
}

// ResetCompletionTimestamp reset all changes of the completion_timestamp field.
func (m *SurveyMutation) ResetCompletionTimestamp() {
	m.completion_timestamp = nil
}

// SetLocationID sets the location edge to Location by id.
func (m *SurveyMutation) SetLocationID(id int) {
	m.location = &id
}

// ClearLocation clears the location edge to Location.
func (m *SurveyMutation) ClearLocation() {
	m.clearedlocation = true
}

// LocationCleared returns if the edge location was cleared.
func (m *SurveyMutation) LocationCleared() bool {
	return m.clearedlocation
}

// LocationID returns the location id in the mutation.
func (m *SurveyMutation) LocationID() (id int, exists bool) {
	if m.location != nil {
		return *m.location, true
	}
	return
}

// LocationIDs returns the location ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LocationID instead. It exists only for internal usage by the builders.
func (m *SurveyMutation) LocationIDs() (ids []int) {
	if id := m.location; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLocation reset all changes of the location edge.
func (m *SurveyMutation) ResetLocation() {
	m.location = nil
	m.clearedlocation = false
}

// SetSourceFileID sets the source_file edge to File by id.
func (m *SurveyMutation) SetSourceFileID(id int) {
	m.source_file = &id
}

// ClearSourceFile clears the source_file edge to File.
func (m *SurveyMutation) ClearSourceFile() {
	m.clearedsource_file = true
}

// SourceFileCleared returns if the edge source_file was cleared.
func (m *SurveyMutation) SourceFileCleared() bool {
	return m.clearedsource_file
}

// SourceFileID returns the source_file id in the mutation.
func (m *SurveyMutation) SourceFileID() (id int, exists bool) {
	if m.source_file != nil {
		return *m.source_file, true
	}
	return
}

// SourceFileIDs returns the source_file ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// SourceFileID instead. It exists only for internal usage by the builders.
func (m *SurveyMutation) SourceFileIDs() (ids []int) {
	if id := m.source_file; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetSourceFile reset all changes of the source_file edge.
func (m *SurveyMutation) ResetSourceFile() {
	m.source_file = nil
	m.clearedsource_file = false
}

// AddQuestionIDs adds the questions edge to SurveyQuestion by ids.
func (m *SurveyMutation) AddQuestionIDs(ids ...int) {
	if m.questions == nil {
		m.questions = make(map[int]struct{})
	}
	for i := range ids {
		m.questions[ids[i]] = struct{}{}
	}
}

// RemoveQuestionIDs removes the questions edge to SurveyQuestion by ids.
func (m *SurveyMutation) RemoveQuestionIDs(ids ...int) {
	if m.removedquestions == nil {
		m.removedquestions = make(map[int]struct{})
	}
	for i := range ids {
		m.removedquestions[ids[i]] = struct{}{}
	}
}

// RemovedQuestions returns the removed ids of questions.
func (m *SurveyMutation) RemovedQuestionsIDs() (ids []int) {
	for id := range m.removedquestions {
		ids = append(ids, id)
	}
	return
}

// QuestionsIDs returns the questions ids in the mutation.
func (m *SurveyMutation) QuestionsIDs() (ids []int) {
	for id := range m.questions {
		ids = append(ids, id)
	}
	return
}

// ResetQuestions reset all changes of the questions edge.
func (m *SurveyMutation) ResetQuestions() {
	m.questions = nil
	m.removedquestions = nil
}

// Op returns the operation name.
func (m *SurveyMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Survey).
func (m *SurveyMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *SurveyMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m.create_time != nil {
		fields = append(fields, survey.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, survey.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, survey.FieldName)
	}
	if m.owner_name != nil {
		fields = append(fields, survey.FieldOwnerName)
	}
	if m.creation_timestamp != nil {
		fields = append(fields, survey.FieldCreationTimestamp)
	}
	if m.completion_timestamp != nil {
		fields = append(fields, survey.FieldCompletionTimestamp)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *SurveyMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case survey.FieldCreateTime:
		return m.CreateTime()
	case survey.FieldUpdateTime:
		return m.UpdateTime()
	case survey.FieldName:
		return m.Name()
	case survey.FieldOwnerName:
		return m.OwnerName()
	case survey.FieldCreationTimestamp:
		return m.CreationTimestamp()
	case survey.FieldCompletionTimestamp:
		return m.CompletionTimestamp()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyMutation) SetField(name string, value ent.Value) error {
	switch name {
	case survey.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case survey.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case survey.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case survey.FieldOwnerName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetOwnerName(v)
		return nil
	case survey.FieldCreationTimestamp:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreationTimestamp(v)
		return nil
	case survey.FieldCompletionTimestamp:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCompletionTimestamp(v)
		return nil
	}
	return fmt.Errorf("unknown Survey field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *SurveyMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *SurveyMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Survey numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *SurveyMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[survey.FieldOwnerName] {
		fields = append(fields, survey.FieldOwnerName)
	}
	if m.clearedFields[survey.FieldCreationTimestamp] {
		fields = append(fields, survey.FieldCreationTimestamp)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *SurveyMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *SurveyMutation) ClearField(name string) error {
	switch name {
	case survey.FieldOwnerName:
		m.ClearOwnerName()
		return nil
	case survey.FieldCreationTimestamp:
		m.ClearCreationTimestamp()
		return nil
	}
	return fmt.Errorf("unknown Survey nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *SurveyMutation) ResetField(name string) error {
	switch name {
	case survey.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case survey.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case survey.FieldName:
		m.ResetName()
		return nil
	case survey.FieldOwnerName:
		m.ResetOwnerName()
		return nil
	case survey.FieldCreationTimestamp:
		m.ResetCreationTimestamp()
		return nil
	case survey.FieldCompletionTimestamp:
		m.ResetCompletionTimestamp()
		return nil
	}
	return fmt.Errorf("unknown Survey field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *SurveyMutation) AddedEdges() []string {
	edges := make([]string, 0, 3)
	if m.location != nil {
		edges = append(edges, survey.EdgeLocation)
	}
	if m.source_file != nil {
		edges = append(edges, survey.EdgeSourceFile)
	}
	if m.questions != nil {
		edges = append(edges, survey.EdgeQuestions)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *SurveyMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case survey.EdgeLocation:
		if id := m.location; id != nil {
			return []ent.Value{*id}
		}
	case survey.EdgeSourceFile:
		if id := m.source_file; id != nil {
			return []ent.Value{*id}
		}
	case survey.EdgeQuestions:
		ids := make([]ent.Value, 0, len(m.questions))
		for id := range m.questions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *SurveyMutation) RemovedEdges() []string {
	edges := make([]string, 0, 3)
	if m.removedquestions != nil {
		edges = append(edges, survey.EdgeQuestions)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *SurveyMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case survey.EdgeQuestions:
		ids := make([]ent.Value, 0, len(m.removedquestions))
		for id := range m.removedquestions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *SurveyMutation) ClearedEdges() []string {
	edges := make([]string, 0, 3)
	if m.clearedlocation {
		edges = append(edges, survey.EdgeLocation)
	}
	if m.clearedsource_file {
		edges = append(edges, survey.EdgeSourceFile)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *SurveyMutation) EdgeCleared(name string) bool {
	switch name {
	case survey.EdgeLocation:
		return m.clearedlocation
	case survey.EdgeSourceFile:
		return m.clearedsource_file
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *SurveyMutation) ClearEdge(name string) error {
	switch name {
	case survey.EdgeLocation:
		m.ClearLocation()
		return nil
	case survey.EdgeSourceFile:
		m.ClearSourceFile()
		return nil
	}
	return fmt.Errorf("unknown Survey unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *SurveyMutation) ResetEdge(name string) error {
	switch name {
	case survey.EdgeLocation:
		m.ResetLocation()
		return nil
	case survey.EdgeSourceFile:
		m.ResetSourceFile()
		return nil
	case survey.EdgeQuestions:
		m.ResetQuestions()
		return nil
	}
	return fmt.Errorf("unknown Survey edge %s", name)
}

// SurveyCellScanMutation represents an operation that mutate the SurveyCellScans
// nodes in the graph.
type SurveyCellScanMutation struct {
	config
	op                      Op
	typ                     string
	id                      *int
	create_time             *time.Time
	update_time             *time.Time
	network_type            *string
	signal_strength         *int
	addsignal_strength      *int
	timestamp               *time.Time
	base_station_id         *string
	network_id              *string
	system_id               *string
	cell_id                 *string
	location_area_code      *string
	mobile_country_code     *string
	mobile_network_code     *string
	primary_scrambling_code *string
	operator                *string
	arfcn                   *int
	addarfcn                *int
	physical_cell_id        *string
	tracking_area_code      *string
	timing_advance          *int
	addtiming_advance       *int
	earfcn                  *int
	addearfcn               *int
	uarfcn                  *int
	adduarfcn               *int
	latitude                *float64
	addlatitude             *float64
	longitude               *float64
	addlongitude            *float64
	clearedFields           map[string]bool
	survey_question         *int
	clearedsurvey_question  bool
	location                *int
	clearedlocation         bool
}

var _ ent.Mutation = (*SurveyCellScanMutation)(nil)

// newSurveyCellScanMutation creates new mutation for $n.Name.
func newSurveyCellScanMutation(c config, op Op) *SurveyCellScanMutation {
	return &SurveyCellScanMutation{
		config:        c,
		op:            op,
		typ:           TypeSurveyCellScan,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SurveyCellScanMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SurveyCellScanMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *SurveyCellScanMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *SurveyCellScanMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *SurveyCellScanMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *SurveyCellScanMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *SurveyCellScanMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *SurveyCellScanMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *SurveyCellScanMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetNetworkType sets the network_type field.
func (m *SurveyCellScanMutation) SetNetworkType(s string) {
	m.network_type = &s
}

// NetworkType returns the network_type value in the mutation.
func (m *SurveyCellScanMutation) NetworkType() (r string, exists bool) {
	v := m.network_type
	if v == nil {
		return
	}
	return *v, true
}

// ResetNetworkType reset all changes of the network_type field.
func (m *SurveyCellScanMutation) ResetNetworkType() {
	m.network_type = nil
}

// SetSignalStrength sets the signal_strength field.
func (m *SurveyCellScanMutation) SetSignalStrength(i int) {
	m.signal_strength = &i
	m.addsignal_strength = nil
}

// SignalStrength returns the signal_strength value in the mutation.
func (m *SurveyCellScanMutation) SignalStrength() (r int, exists bool) {
	v := m.signal_strength
	if v == nil {
		return
	}
	return *v, true
}

// AddSignalStrength adds i to signal_strength.
func (m *SurveyCellScanMutation) AddSignalStrength(i int) {
	if m.addsignal_strength != nil {
		*m.addsignal_strength += i
	} else {
		m.addsignal_strength = &i
	}
}

// AddedSignalStrength returns the value that was added to the signal_strength field in this mutation.
func (m *SurveyCellScanMutation) AddedSignalStrength() (r int, exists bool) {
	v := m.addsignal_strength
	if v == nil {
		return
	}
	return *v, true
}

// ResetSignalStrength reset all changes of the signal_strength field.
func (m *SurveyCellScanMutation) ResetSignalStrength() {
	m.signal_strength = nil
	m.addsignal_strength = nil
}

// SetTimestamp sets the timestamp field.
func (m *SurveyCellScanMutation) SetTimestamp(t time.Time) {
	m.timestamp = &t
}

// Timestamp returns the timestamp value in the mutation.
func (m *SurveyCellScanMutation) Timestamp() (r time.Time, exists bool) {
	v := m.timestamp
	if v == nil {
		return
	}
	return *v, true
}

// ClearTimestamp clears the value of timestamp.
func (m *SurveyCellScanMutation) ClearTimestamp() {
	m.timestamp = nil
	m.clearedFields[surveycellscan.FieldTimestamp] = true
}

// TimestampCleared returns if the field timestamp was cleared in this mutation.
func (m *SurveyCellScanMutation) TimestampCleared() bool {
	return m.clearedFields[surveycellscan.FieldTimestamp]
}

// ResetTimestamp reset all changes of the timestamp field.
func (m *SurveyCellScanMutation) ResetTimestamp() {
	m.timestamp = nil
	delete(m.clearedFields, surveycellscan.FieldTimestamp)
}

// SetBaseStationID sets the base_station_id field.
func (m *SurveyCellScanMutation) SetBaseStationID(s string) {
	m.base_station_id = &s
}

// BaseStationID returns the base_station_id value in the mutation.
func (m *SurveyCellScanMutation) BaseStationID() (r string, exists bool) {
	v := m.base_station_id
	if v == nil {
		return
	}
	return *v, true
}

// ClearBaseStationID clears the value of base_station_id.
func (m *SurveyCellScanMutation) ClearBaseStationID() {
	m.base_station_id = nil
	m.clearedFields[surveycellscan.FieldBaseStationID] = true
}

// BaseStationIDCleared returns if the field base_station_id was cleared in this mutation.
func (m *SurveyCellScanMutation) BaseStationIDCleared() bool {
	return m.clearedFields[surveycellscan.FieldBaseStationID]
}

// ResetBaseStationID reset all changes of the base_station_id field.
func (m *SurveyCellScanMutation) ResetBaseStationID() {
	m.base_station_id = nil
	delete(m.clearedFields, surveycellscan.FieldBaseStationID)
}

// SetNetworkID sets the network_id field.
func (m *SurveyCellScanMutation) SetNetworkID(s string) {
	m.network_id = &s
}

// NetworkID returns the network_id value in the mutation.
func (m *SurveyCellScanMutation) NetworkID() (r string, exists bool) {
	v := m.network_id
	if v == nil {
		return
	}
	return *v, true
}

// ClearNetworkID clears the value of network_id.
func (m *SurveyCellScanMutation) ClearNetworkID() {
	m.network_id = nil
	m.clearedFields[surveycellscan.FieldNetworkID] = true
}

// NetworkIDCleared returns if the field network_id was cleared in this mutation.
func (m *SurveyCellScanMutation) NetworkIDCleared() bool {
	return m.clearedFields[surveycellscan.FieldNetworkID]
}

// ResetNetworkID reset all changes of the network_id field.
func (m *SurveyCellScanMutation) ResetNetworkID() {
	m.network_id = nil
	delete(m.clearedFields, surveycellscan.FieldNetworkID)
}

// SetSystemID sets the system_id field.
func (m *SurveyCellScanMutation) SetSystemID(s string) {
	m.system_id = &s
}

// SystemID returns the system_id value in the mutation.
func (m *SurveyCellScanMutation) SystemID() (r string, exists bool) {
	v := m.system_id
	if v == nil {
		return
	}
	return *v, true
}

// ClearSystemID clears the value of system_id.
func (m *SurveyCellScanMutation) ClearSystemID() {
	m.system_id = nil
	m.clearedFields[surveycellscan.FieldSystemID] = true
}

// SystemIDCleared returns if the field system_id was cleared in this mutation.
func (m *SurveyCellScanMutation) SystemIDCleared() bool {
	return m.clearedFields[surveycellscan.FieldSystemID]
}

// ResetSystemID reset all changes of the system_id field.
func (m *SurveyCellScanMutation) ResetSystemID() {
	m.system_id = nil
	delete(m.clearedFields, surveycellscan.FieldSystemID)
}

// SetCellID sets the cell_id field.
func (m *SurveyCellScanMutation) SetCellID(s string) {
	m.cell_id = &s
}

// CellID returns the cell_id value in the mutation.
func (m *SurveyCellScanMutation) CellID() (r string, exists bool) {
	v := m.cell_id
	if v == nil {
		return
	}
	return *v, true
}

// ClearCellID clears the value of cell_id.
func (m *SurveyCellScanMutation) ClearCellID() {
	m.cell_id = nil
	m.clearedFields[surveycellscan.FieldCellID] = true
}

// CellIDCleared returns if the field cell_id was cleared in this mutation.
func (m *SurveyCellScanMutation) CellIDCleared() bool {
	return m.clearedFields[surveycellscan.FieldCellID]
}

// ResetCellID reset all changes of the cell_id field.
func (m *SurveyCellScanMutation) ResetCellID() {
	m.cell_id = nil
	delete(m.clearedFields, surveycellscan.FieldCellID)
}

// SetLocationAreaCode sets the location_area_code field.
func (m *SurveyCellScanMutation) SetLocationAreaCode(s string) {
	m.location_area_code = &s
}

// LocationAreaCode returns the location_area_code value in the mutation.
func (m *SurveyCellScanMutation) LocationAreaCode() (r string, exists bool) {
	v := m.location_area_code
	if v == nil {
		return
	}
	return *v, true
}

// ClearLocationAreaCode clears the value of location_area_code.
func (m *SurveyCellScanMutation) ClearLocationAreaCode() {
	m.location_area_code = nil
	m.clearedFields[surveycellscan.FieldLocationAreaCode] = true
}

// LocationAreaCodeCleared returns if the field location_area_code was cleared in this mutation.
func (m *SurveyCellScanMutation) LocationAreaCodeCleared() bool {
	return m.clearedFields[surveycellscan.FieldLocationAreaCode]
}

// ResetLocationAreaCode reset all changes of the location_area_code field.
func (m *SurveyCellScanMutation) ResetLocationAreaCode() {
	m.location_area_code = nil
	delete(m.clearedFields, surveycellscan.FieldLocationAreaCode)
}

// SetMobileCountryCode sets the mobile_country_code field.
func (m *SurveyCellScanMutation) SetMobileCountryCode(s string) {
	m.mobile_country_code = &s
}

// MobileCountryCode returns the mobile_country_code value in the mutation.
func (m *SurveyCellScanMutation) MobileCountryCode() (r string, exists bool) {
	v := m.mobile_country_code
	if v == nil {
		return
	}
	return *v, true
}

// ClearMobileCountryCode clears the value of mobile_country_code.
func (m *SurveyCellScanMutation) ClearMobileCountryCode() {
	m.mobile_country_code = nil
	m.clearedFields[surveycellscan.FieldMobileCountryCode] = true
}

// MobileCountryCodeCleared returns if the field mobile_country_code was cleared in this mutation.
func (m *SurveyCellScanMutation) MobileCountryCodeCleared() bool {
	return m.clearedFields[surveycellscan.FieldMobileCountryCode]
}

// ResetMobileCountryCode reset all changes of the mobile_country_code field.
func (m *SurveyCellScanMutation) ResetMobileCountryCode() {
	m.mobile_country_code = nil
	delete(m.clearedFields, surveycellscan.FieldMobileCountryCode)
}

// SetMobileNetworkCode sets the mobile_network_code field.
func (m *SurveyCellScanMutation) SetMobileNetworkCode(s string) {
	m.mobile_network_code = &s
}

// MobileNetworkCode returns the mobile_network_code value in the mutation.
func (m *SurveyCellScanMutation) MobileNetworkCode() (r string, exists bool) {
	v := m.mobile_network_code
	if v == nil {
		return
	}
	return *v, true
}

// ClearMobileNetworkCode clears the value of mobile_network_code.
func (m *SurveyCellScanMutation) ClearMobileNetworkCode() {
	m.mobile_network_code = nil
	m.clearedFields[surveycellscan.FieldMobileNetworkCode] = true
}

// MobileNetworkCodeCleared returns if the field mobile_network_code was cleared in this mutation.
func (m *SurveyCellScanMutation) MobileNetworkCodeCleared() bool {
	return m.clearedFields[surveycellscan.FieldMobileNetworkCode]
}

// ResetMobileNetworkCode reset all changes of the mobile_network_code field.
func (m *SurveyCellScanMutation) ResetMobileNetworkCode() {
	m.mobile_network_code = nil
	delete(m.clearedFields, surveycellscan.FieldMobileNetworkCode)
}

// SetPrimaryScramblingCode sets the primary_scrambling_code field.
func (m *SurveyCellScanMutation) SetPrimaryScramblingCode(s string) {
	m.primary_scrambling_code = &s
}

// PrimaryScramblingCode returns the primary_scrambling_code value in the mutation.
func (m *SurveyCellScanMutation) PrimaryScramblingCode() (r string, exists bool) {
	v := m.primary_scrambling_code
	if v == nil {
		return
	}
	return *v, true
}

// ClearPrimaryScramblingCode clears the value of primary_scrambling_code.
func (m *SurveyCellScanMutation) ClearPrimaryScramblingCode() {
	m.primary_scrambling_code = nil
	m.clearedFields[surveycellscan.FieldPrimaryScramblingCode] = true
}

// PrimaryScramblingCodeCleared returns if the field primary_scrambling_code was cleared in this mutation.
func (m *SurveyCellScanMutation) PrimaryScramblingCodeCleared() bool {
	return m.clearedFields[surveycellscan.FieldPrimaryScramblingCode]
}

// ResetPrimaryScramblingCode reset all changes of the primary_scrambling_code field.
func (m *SurveyCellScanMutation) ResetPrimaryScramblingCode() {
	m.primary_scrambling_code = nil
	delete(m.clearedFields, surveycellscan.FieldPrimaryScramblingCode)
}

// SetOperator sets the operator field.
func (m *SurveyCellScanMutation) SetOperator(s string) {
	m.operator = &s
}

// Operator returns the operator value in the mutation.
func (m *SurveyCellScanMutation) Operator() (r string, exists bool) {
	v := m.operator
	if v == nil {
		return
	}
	return *v, true
}

// ClearOperator clears the value of operator.
func (m *SurveyCellScanMutation) ClearOperator() {
	m.operator = nil
	m.clearedFields[surveycellscan.FieldOperator] = true
}

// OperatorCleared returns if the field operator was cleared in this mutation.
func (m *SurveyCellScanMutation) OperatorCleared() bool {
	return m.clearedFields[surveycellscan.FieldOperator]
}

// ResetOperator reset all changes of the operator field.
func (m *SurveyCellScanMutation) ResetOperator() {
	m.operator = nil
	delete(m.clearedFields, surveycellscan.FieldOperator)
}

// SetArfcn sets the arfcn field.
func (m *SurveyCellScanMutation) SetArfcn(i int) {
	m.arfcn = &i
	m.addarfcn = nil
}

// Arfcn returns the arfcn value in the mutation.
func (m *SurveyCellScanMutation) Arfcn() (r int, exists bool) {
	v := m.arfcn
	if v == nil {
		return
	}
	return *v, true
}

// AddArfcn adds i to arfcn.
func (m *SurveyCellScanMutation) AddArfcn(i int) {
	if m.addarfcn != nil {
		*m.addarfcn += i
	} else {
		m.addarfcn = &i
	}
}

// AddedArfcn returns the value that was added to the arfcn field in this mutation.
func (m *SurveyCellScanMutation) AddedArfcn() (r int, exists bool) {
	v := m.addarfcn
	if v == nil {
		return
	}
	return *v, true
}

// ClearArfcn clears the value of arfcn.
func (m *SurveyCellScanMutation) ClearArfcn() {
	m.arfcn = nil
	m.addarfcn = nil
	m.clearedFields[surveycellscan.FieldArfcn] = true
}

// ArfcnCleared returns if the field arfcn was cleared in this mutation.
func (m *SurveyCellScanMutation) ArfcnCleared() bool {
	return m.clearedFields[surveycellscan.FieldArfcn]
}

// ResetArfcn reset all changes of the arfcn field.
func (m *SurveyCellScanMutation) ResetArfcn() {
	m.arfcn = nil
	m.addarfcn = nil
	delete(m.clearedFields, surveycellscan.FieldArfcn)
}

// SetPhysicalCellID sets the physical_cell_id field.
func (m *SurveyCellScanMutation) SetPhysicalCellID(s string) {
	m.physical_cell_id = &s
}

// PhysicalCellID returns the physical_cell_id value in the mutation.
func (m *SurveyCellScanMutation) PhysicalCellID() (r string, exists bool) {
	v := m.physical_cell_id
	if v == nil {
		return
	}
	return *v, true
}

// ClearPhysicalCellID clears the value of physical_cell_id.
func (m *SurveyCellScanMutation) ClearPhysicalCellID() {
	m.physical_cell_id = nil
	m.clearedFields[surveycellscan.FieldPhysicalCellID] = true
}

// PhysicalCellIDCleared returns if the field physical_cell_id was cleared in this mutation.
func (m *SurveyCellScanMutation) PhysicalCellIDCleared() bool {
	return m.clearedFields[surveycellscan.FieldPhysicalCellID]
}

// ResetPhysicalCellID reset all changes of the physical_cell_id field.
func (m *SurveyCellScanMutation) ResetPhysicalCellID() {
	m.physical_cell_id = nil
	delete(m.clearedFields, surveycellscan.FieldPhysicalCellID)
}

// SetTrackingAreaCode sets the tracking_area_code field.
func (m *SurveyCellScanMutation) SetTrackingAreaCode(s string) {
	m.tracking_area_code = &s
}

// TrackingAreaCode returns the tracking_area_code value in the mutation.
func (m *SurveyCellScanMutation) TrackingAreaCode() (r string, exists bool) {
	v := m.tracking_area_code
	if v == nil {
		return
	}
	return *v, true
}

// ClearTrackingAreaCode clears the value of tracking_area_code.
func (m *SurveyCellScanMutation) ClearTrackingAreaCode() {
	m.tracking_area_code = nil
	m.clearedFields[surveycellscan.FieldTrackingAreaCode] = true
}

// TrackingAreaCodeCleared returns if the field tracking_area_code was cleared in this mutation.
func (m *SurveyCellScanMutation) TrackingAreaCodeCleared() bool {
	return m.clearedFields[surveycellscan.FieldTrackingAreaCode]
}

// ResetTrackingAreaCode reset all changes of the tracking_area_code field.
func (m *SurveyCellScanMutation) ResetTrackingAreaCode() {
	m.tracking_area_code = nil
	delete(m.clearedFields, surveycellscan.FieldTrackingAreaCode)
}

// SetTimingAdvance sets the timing_advance field.
func (m *SurveyCellScanMutation) SetTimingAdvance(i int) {
	m.timing_advance = &i
	m.addtiming_advance = nil
}

// TimingAdvance returns the timing_advance value in the mutation.
func (m *SurveyCellScanMutation) TimingAdvance() (r int, exists bool) {
	v := m.timing_advance
	if v == nil {
		return
	}
	return *v, true
}

// AddTimingAdvance adds i to timing_advance.
func (m *SurveyCellScanMutation) AddTimingAdvance(i int) {
	if m.addtiming_advance != nil {
		*m.addtiming_advance += i
	} else {
		m.addtiming_advance = &i
	}
}

// AddedTimingAdvance returns the value that was added to the timing_advance field in this mutation.
func (m *SurveyCellScanMutation) AddedTimingAdvance() (r int, exists bool) {
	v := m.addtiming_advance
	if v == nil {
		return
	}
	return *v, true
}

// ClearTimingAdvance clears the value of timing_advance.
func (m *SurveyCellScanMutation) ClearTimingAdvance() {
	m.timing_advance = nil
	m.addtiming_advance = nil
	m.clearedFields[surveycellscan.FieldTimingAdvance] = true
}

// TimingAdvanceCleared returns if the field timing_advance was cleared in this mutation.
func (m *SurveyCellScanMutation) TimingAdvanceCleared() bool {
	return m.clearedFields[surveycellscan.FieldTimingAdvance]
}

// ResetTimingAdvance reset all changes of the timing_advance field.
func (m *SurveyCellScanMutation) ResetTimingAdvance() {
	m.timing_advance = nil
	m.addtiming_advance = nil
	delete(m.clearedFields, surveycellscan.FieldTimingAdvance)
}

// SetEarfcn sets the earfcn field.
func (m *SurveyCellScanMutation) SetEarfcn(i int) {
	m.earfcn = &i
	m.addearfcn = nil
}

// Earfcn returns the earfcn value in the mutation.
func (m *SurveyCellScanMutation) Earfcn() (r int, exists bool) {
	v := m.earfcn
	if v == nil {
		return
	}
	return *v, true
}

// AddEarfcn adds i to earfcn.
func (m *SurveyCellScanMutation) AddEarfcn(i int) {
	if m.addearfcn != nil {
		*m.addearfcn += i
	} else {
		m.addearfcn = &i
	}
}

// AddedEarfcn returns the value that was added to the earfcn field in this mutation.
func (m *SurveyCellScanMutation) AddedEarfcn() (r int, exists bool) {
	v := m.addearfcn
	if v == nil {
		return
	}
	return *v, true
}

// ClearEarfcn clears the value of earfcn.
func (m *SurveyCellScanMutation) ClearEarfcn() {
	m.earfcn = nil
	m.addearfcn = nil
	m.clearedFields[surveycellscan.FieldEarfcn] = true
}

// EarfcnCleared returns if the field earfcn was cleared in this mutation.
func (m *SurveyCellScanMutation) EarfcnCleared() bool {
	return m.clearedFields[surveycellscan.FieldEarfcn]
}

// ResetEarfcn reset all changes of the earfcn field.
func (m *SurveyCellScanMutation) ResetEarfcn() {
	m.earfcn = nil
	m.addearfcn = nil
	delete(m.clearedFields, surveycellscan.FieldEarfcn)
}

// SetUarfcn sets the uarfcn field.
func (m *SurveyCellScanMutation) SetUarfcn(i int) {
	m.uarfcn = &i
	m.adduarfcn = nil
}

// Uarfcn returns the uarfcn value in the mutation.
func (m *SurveyCellScanMutation) Uarfcn() (r int, exists bool) {
	v := m.uarfcn
	if v == nil {
		return
	}
	return *v, true
}

// AddUarfcn adds i to uarfcn.
func (m *SurveyCellScanMutation) AddUarfcn(i int) {
	if m.adduarfcn != nil {
		*m.adduarfcn += i
	} else {
		m.adduarfcn = &i
	}
}

// AddedUarfcn returns the value that was added to the uarfcn field in this mutation.
func (m *SurveyCellScanMutation) AddedUarfcn() (r int, exists bool) {
	v := m.adduarfcn
	if v == nil {
		return
	}
	return *v, true
}

// ClearUarfcn clears the value of uarfcn.
func (m *SurveyCellScanMutation) ClearUarfcn() {
	m.uarfcn = nil
	m.adduarfcn = nil
	m.clearedFields[surveycellscan.FieldUarfcn] = true
}

// UarfcnCleared returns if the field uarfcn was cleared in this mutation.
func (m *SurveyCellScanMutation) UarfcnCleared() bool {
	return m.clearedFields[surveycellscan.FieldUarfcn]
}

// ResetUarfcn reset all changes of the uarfcn field.
func (m *SurveyCellScanMutation) ResetUarfcn() {
	m.uarfcn = nil
	m.adduarfcn = nil
	delete(m.clearedFields, surveycellscan.FieldUarfcn)
}

// SetLatitude sets the latitude field.
func (m *SurveyCellScanMutation) SetLatitude(f float64) {
	m.latitude = &f
	m.addlatitude = nil
}

// Latitude returns the latitude value in the mutation.
func (m *SurveyCellScanMutation) Latitude() (r float64, exists bool) {
	v := m.latitude
	if v == nil {
		return
	}
	return *v, true
}

// AddLatitude adds f to latitude.
func (m *SurveyCellScanMutation) AddLatitude(f float64) {
	if m.addlatitude != nil {
		*m.addlatitude += f
	} else {
		m.addlatitude = &f
	}
}

// AddedLatitude returns the value that was added to the latitude field in this mutation.
func (m *SurveyCellScanMutation) AddedLatitude() (r float64, exists bool) {
	v := m.addlatitude
	if v == nil {
		return
	}
	return *v, true
}

// ClearLatitude clears the value of latitude.
func (m *SurveyCellScanMutation) ClearLatitude() {
	m.latitude = nil
	m.addlatitude = nil
	m.clearedFields[surveycellscan.FieldLatitude] = true
}

// LatitudeCleared returns if the field latitude was cleared in this mutation.
func (m *SurveyCellScanMutation) LatitudeCleared() bool {
	return m.clearedFields[surveycellscan.FieldLatitude]
}

// ResetLatitude reset all changes of the latitude field.
func (m *SurveyCellScanMutation) ResetLatitude() {
	m.latitude = nil
	m.addlatitude = nil
	delete(m.clearedFields, surveycellscan.FieldLatitude)
}

// SetLongitude sets the longitude field.
func (m *SurveyCellScanMutation) SetLongitude(f float64) {
	m.longitude = &f
	m.addlongitude = nil
}

// Longitude returns the longitude value in the mutation.
func (m *SurveyCellScanMutation) Longitude() (r float64, exists bool) {
	v := m.longitude
	if v == nil {
		return
	}
	return *v, true
}

// AddLongitude adds f to longitude.
func (m *SurveyCellScanMutation) AddLongitude(f float64) {
	if m.addlongitude != nil {
		*m.addlongitude += f
	} else {
		m.addlongitude = &f
	}
}

// AddedLongitude returns the value that was added to the longitude field in this mutation.
func (m *SurveyCellScanMutation) AddedLongitude() (r float64, exists bool) {
	v := m.addlongitude
	if v == nil {
		return
	}
	return *v, true
}

// ClearLongitude clears the value of longitude.
func (m *SurveyCellScanMutation) ClearLongitude() {
	m.longitude = nil
	m.addlongitude = nil
	m.clearedFields[surveycellscan.FieldLongitude] = true
}

// LongitudeCleared returns if the field longitude was cleared in this mutation.
func (m *SurveyCellScanMutation) LongitudeCleared() bool {
	return m.clearedFields[surveycellscan.FieldLongitude]
}

// ResetLongitude reset all changes of the longitude field.
func (m *SurveyCellScanMutation) ResetLongitude() {
	m.longitude = nil
	m.addlongitude = nil
	delete(m.clearedFields, surveycellscan.FieldLongitude)
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (m *SurveyCellScanMutation) SetSurveyQuestionID(id int) {
	m.survey_question = &id
}

// ClearSurveyQuestion clears the survey_question edge to SurveyQuestion.
func (m *SurveyCellScanMutation) ClearSurveyQuestion() {
	m.clearedsurvey_question = true
}

// SurveyQuestionCleared returns if the edge survey_question was cleared.
func (m *SurveyCellScanMutation) SurveyQuestionCleared() bool {
	return m.clearedsurvey_question
}

// SurveyQuestionID returns the survey_question id in the mutation.
func (m *SurveyCellScanMutation) SurveyQuestionID() (id int, exists bool) {
	if m.survey_question != nil {
		return *m.survey_question, true
	}
	return
}

// SurveyQuestionIDs returns the survey_question ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// SurveyQuestionID instead. It exists only for internal usage by the builders.
func (m *SurveyCellScanMutation) SurveyQuestionIDs() (ids []int) {
	if id := m.survey_question; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetSurveyQuestion reset all changes of the survey_question edge.
func (m *SurveyCellScanMutation) ResetSurveyQuestion() {
	m.survey_question = nil
	m.clearedsurvey_question = false
}

// SetLocationID sets the location edge to Location by id.
func (m *SurveyCellScanMutation) SetLocationID(id int) {
	m.location = &id
}

// ClearLocation clears the location edge to Location.
func (m *SurveyCellScanMutation) ClearLocation() {
	m.clearedlocation = true
}

// LocationCleared returns if the edge location was cleared.
func (m *SurveyCellScanMutation) LocationCleared() bool {
	return m.clearedlocation
}

// LocationID returns the location id in the mutation.
func (m *SurveyCellScanMutation) LocationID() (id int, exists bool) {
	if m.location != nil {
		return *m.location, true
	}
	return
}

// LocationIDs returns the location ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LocationID instead. It exists only for internal usage by the builders.
func (m *SurveyCellScanMutation) LocationIDs() (ids []int) {
	if id := m.location; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLocation reset all changes of the location edge.
func (m *SurveyCellScanMutation) ResetLocation() {
	m.location = nil
	m.clearedlocation = false
}

// Op returns the operation name.
func (m *SurveyCellScanMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (SurveyCellScan).
func (m *SurveyCellScanMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *SurveyCellScanMutation) Fields() []string {
	fields := make([]string, 0, 22)
	if m.create_time != nil {
		fields = append(fields, surveycellscan.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, surveycellscan.FieldUpdateTime)
	}
	if m.network_type != nil {
		fields = append(fields, surveycellscan.FieldNetworkType)
	}
	if m.signal_strength != nil {
		fields = append(fields, surveycellscan.FieldSignalStrength)
	}
	if m.timestamp != nil {
		fields = append(fields, surveycellscan.FieldTimestamp)
	}
	if m.base_station_id != nil {
		fields = append(fields, surveycellscan.FieldBaseStationID)
	}
	if m.network_id != nil {
		fields = append(fields, surveycellscan.FieldNetworkID)
	}
	if m.system_id != nil {
		fields = append(fields, surveycellscan.FieldSystemID)
	}
	if m.cell_id != nil {
		fields = append(fields, surveycellscan.FieldCellID)
	}
	if m.location_area_code != nil {
		fields = append(fields, surveycellscan.FieldLocationAreaCode)
	}
	if m.mobile_country_code != nil {
		fields = append(fields, surveycellscan.FieldMobileCountryCode)
	}
	if m.mobile_network_code != nil {
		fields = append(fields, surveycellscan.FieldMobileNetworkCode)
	}
	if m.primary_scrambling_code != nil {
		fields = append(fields, surveycellscan.FieldPrimaryScramblingCode)
	}
	if m.operator != nil {
		fields = append(fields, surveycellscan.FieldOperator)
	}
	if m.arfcn != nil {
		fields = append(fields, surveycellscan.FieldArfcn)
	}
	if m.physical_cell_id != nil {
		fields = append(fields, surveycellscan.FieldPhysicalCellID)
	}
	if m.tracking_area_code != nil {
		fields = append(fields, surveycellscan.FieldTrackingAreaCode)
	}
	if m.timing_advance != nil {
		fields = append(fields, surveycellscan.FieldTimingAdvance)
	}
	if m.earfcn != nil {
		fields = append(fields, surveycellscan.FieldEarfcn)
	}
	if m.uarfcn != nil {
		fields = append(fields, surveycellscan.FieldUarfcn)
	}
	if m.latitude != nil {
		fields = append(fields, surveycellscan.FieldLatitude)
	}
	if m.longitude != nil {
		fields = append(fields, surveycellscan.FieldLongitude)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *SurveyCellScanMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case surveycellscan.FieldCreateTime:
		return m.CreateTime()
	case surveycellscan.FieldUpdateTime:
		return m.UpdateTime()
	case surveycellscan.FieldNetworkType:
		return m.NetworkType()
	case surveycellscan.FieldSignalStrength:
		return m.SignalStrength()
	case surveycellscan.FieldTimestamp:
		return m.Timestamp()
	case surveycellscan.FieldBaseStationID:
		return m.BaseStationID()
	case surveycellscan.FieldNetworkID:
		return m.NetworkID()
	case surveycellscan.FieldSystemID:
		return m.SystemID()
	case surveycellscan.FieldCellID:
		return m.CellID()
	case surveycellscan.FieldLocationAreaCode:
		return m.LocationAreaCode()
	case surveycellscan.FieldMobileCountryCode:
		return m.MobileCountryCode()
	case surveycellscan.FieldMobileNetworkCode:
		return m.MobileNetworkCode()
	case surveycellscan.FieldPrimaryScramblingCode:
		return m.PrimaryScramblingCode()
	case surveycellscan.FieldOperator:
		return m.Operator()
	case surveycellscan.FieldArfcn:
		return m.Arfcn()
	case surveycellscan.FieldPhysicalCellID:
		return m.PhysicalCellID()
	case surveycellscan.FieldTrackingAreaCode:
		return m.TrackingAreaCode()
	case surveycellscan.FieldTimingAdvance:
		return m.TimingAdvance()
	case surveycellscan.FieldEarfcn:
		return m.Earfcn()
	case surveycellscan.FieldUarfcn:
		return m.Uarfcn()
	case surveycellscan.FieldLatitude:
		return m.Latitude()
	case surveycellscan.FieldLongitude:
		return m.Longitude()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyCellScanMutation) SetField(name string, value ent.Value) error {
	switch name {
	case surveycellscan.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case surveycellscan.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case surveycellscan.FieldNetworkType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetNetworkType(v)
		return nil
	case surveycellscan.FieldSignalStrength:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSignalStrength(v)
		return nil
	case surveycellscan.FieldTimestamp:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTimestamp(v)
		return nil
	case surveycellscan.FieldBaseStationID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBaseStationID(v)
		return nil
	case surveycellscan.FieldNetworkID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetNetworkID(v)
		return nil
	case surveycellscan.FieldSystemID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSystemID(v)
		return nil
	case surveycellscan.FieldCellID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCellID(v)
		return nil
	case surveycellscan.FieldLocationAreaCode:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLocationAreaCode(v)
		return nil
	case surveycellscan.FieldMobileCountryCode:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMobileCountryCode(v)
		return nil
	case surveycellscan.FieldMobileNetworkCode:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetMobileNetworkCode(v)
		return nil
	case surveycellscan.FieldPrimaryScramblingCode:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPrimaryScramblingCode(v)
		return nil
	case surveycellscan.FieldOperator:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetOperator(v)
		return nil
	case surveycellscan.FieldArfcn:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetArfcn(v)
		return nil
	case surveycellscan.FieldPhysicalCellID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPhysicalCellID(v)
		return nil
	case surveycellscan.FieldTrackingAreaCode:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTrackingAreaCode(v)
		return nil
	case surveycellscan.FieldTimingAdvance:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTimingAdvance(v)
		return nil
	case surveycellscan.FieldEarfcn:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEarfcn(v)
		return nil
	case surveycellscan.FieldUarfcn:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUarfcn(v)
		return nil
	case surveycellscan.FieldLatitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLatitude(v)
		return nil
	case surveycellscan.FieldLongitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLongitude(v)
		return nil
	}
	return fmt.Errorf("unknown SurveyCellScan field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *SurveyCellScanMutation) AddedFields() []string {
	var fields []string
	if m.addsignal_strength != nil {
		fields = append(fields, surveycellscan.FieldSignalStrength)
	}
	if m.addarfcn != nil {
		fields = append(fields, surveycellscan.FieldArfcn)
	}
	if m.addtiming_advance != nil {
		fields = append(fields, surveycellscan.FieldTimingAdvance)
	}
	if m.addearfcn != nil {
		fields = append(fields, surveycellscan.FieldEarfcn)
	}
	if m.adduarfcn != nil {
		fields = append(fields, surveycellscan.FieldUarfcn)
	}
	if m.addlatitude != nil {
		fields = append(fields, surveycellscan.FieldLatitude)
	}
	if m.addlongitude != nil {
		fields = append(fields, surveycellscan.FieldLongitude)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *SurveyCellScanMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case surveycellscan.FieldSignalStrength:
		return m.AddedSignalStrength()
	case surveycellscan.FieldArfcn:
		return m.AddedArfcn()
	case surveycellscan.FieldTimingAdvance:
		return m.AddedTimingAdvance()
	case surveycellscan.FieldEarfcn:
		return m.AddedEarfcn()
	case surveycellscan.FieldUarfcn:
		return m.AddedUarfcn()
	case surveycellscan.FieldLatitude:
		return m.AddedLatitude()
	case surveycellscan.FieldLongitude:
		return m.AddedLongitude()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyCellScanMutation) AddField(name string, value ent.Value) error {
	switch name {
	case surveycellscan.FieldSignalStrength:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddSignalStrength(v)
		return nil
	case surveycellscan.FieldArfcn:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddArfcn(v)
		return nil
	case surveycellscan.FieldTimingAdvance:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddTimingAdvance(v)
		return nil
	case surveycellscan.FieldEarfcn:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddEarfcn(v)
		return nil
	case surveycellscan.FieldUarfcn:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddUarfcn(v)
		return nil
	case surveycellscan.FieldLatitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLatitude(v)
		return nil
	case surveycellscan.FieldLongitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLongitude(v)
		return nil
	}
	return fmt.Errorf("unknown SurveyCellScan numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *SurveyCellScanMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[surveycellscan.FieldTimestamp] {
		fields = append(fields, surveycellscan.FieldTimestamp)
	}
	if m.clearedFields[surveycellscan.FieldBaseStationID] {
		fields = append(fields, surveycellscan.FieldBaseStationID)
	}
	if m.clearedFields[surveycellscan.FieldNetworkID] {
		fields = append(fields, surveycellscan.FieldNetworkID)
	}
	if m.clearedFields[surveycellscan.FieldSystemID] {
		fields = append(fields, surveycellscan.FieldSystemID)
	}
	if m.clearedFields[surveycellscan.FieldCellID] {
		fields = append(fields, surveycellscan.FieldCellID)
	}
	if m.clearedFields[surveycellscan.FieldLocationAreaCode] {
		fields = append(fields, surveycellscan.FieldLocationAreaCode)
	}
	if m.clearedFields[surveycellscan.FieldMobileCountryCode] {
		fields = append(fields, surveycellscan.FieldMobileCountryCode)
	}
	if m.clearedFields[surveycellscan.FieldMobileNetworkCode] {
		fields = append(fields, surveycellscan.FieldMobileNetworkCode)
	}
	if m.clearedFields[surveycellscan.FieldPrimaryScramblingCode] {
		fields = append(fields, surveycellscan.FieldPrimaryScramblingCode)
	}
	if m.clearedFields[surveycellscan.FieldOperator] {
		fields = append(fields, surveycellscan.FieldOperator)
	}
	if m.clearedFields[surveycellscan.FieldArfcn] {
		fields = append(fields, surveycellscan.FieldArfcn)
	}
	if m.clearedFields[surveycellscan.FieldPhysicalCellID] {
		fields = append(fields, surveycellscan.FieldPhysicalCellID)
	}
	if m.clearedFields[surveycellscan.FieldTrackingAreaCode] {
		fields = append(fields, surveycellscan.FieldTrackingAreaCode)
	}
	if m.clearedFields[surveycellscan.FieldTimingAdvance] {
		fields = append(fields, surveycellscan.FieldTimingAdvance)
	}
	if m.clearedFields[surveycellscan.FieldEarfcn] {
		fields = append(fields, surveycellscan.FieldEarfcn)
	}
	if m.clearedFields[surveycellscan.FieldUarfcn] {
		fields = append(fields, surveycellscan.FieldUarfcn)
	}
	if m.clearedFields[surveycellscan.FieldLatitude] {
		fields = append(fields, surveycellscan.FieldLatitude)
	}
	if m.clearedFields[surveycellscan.FieldLongitude] {
		fields = append(fields, surveycellscan.FieldLongitude)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *SurveyCellScanMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *SurveyCellScanMutation) ClearField(name string) error {
	switch name {
	case surveycellscan.FieldTimestamp:
		m.ClearTimestamp()
		return nil
	case surveycellscan.FieldBaseStationID:
		m.ClearBaseStationID()
		return nil
	case surveycellscan.FieldNetworkID:
		m.ClearNetworkID()
		return nil
	case surveycellscan.FieldSystemID:
		m.ClearSystemID()
		return nil
	case surveycellscan.FieldCellID:
		m.ClearCellID()
		return nil
	case surveycellscan.FieldLocationAreaCode:
		m.ClearLocationAreaCode()
		return nil
	case surveycellscan.FieldMobileCountryCode:
		m.ClearMobileCountryCode()
		return nil
	case surveycellscan.FieldMobileNetworkCode:
		m.ClearMobileNetworkCode()
		return nil
	case surveycellscan.FieldPrimaryScramblingCode:
		m.ClearPrimaryScramblingCode()
		return nil
	case surveycellscan.FieldOperator:
		m.ClearOperator()
		return nil
	case surveycellscan.FieldArfcn:
		m.ClearArfcn()
		return nil
	case surveycellscan.FieldPhysicalCellID:
		m.ClearPhysicalCellID()
		return nil
	case surveycellscan.FieldTrackingAreaCode:
		m.ClearTrackingAreaCode()
		return nil
	case surveycellscan.FieldTimingAdvance:
		m.ClearTimingAdvance()
		return nil
	case surveycellscan.FieldEarfcn:
		m.ClearEarfcn()
		return nil
	case surveycellscan.FieldUarfcn:
		m.ClearUarfcn()
		return nil
	case surveycellscan.FieldLatitude:
		m.ClearLatitude()
		return nil
	case surveycellscan.FieldLongitude:
		m.ClearLongitude()
		return nil
	}
	return fmt.Errorf("unknown SurveyCellScan nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *SurveyCellScanMutation) ResetField(name string) error {
	switch name {
	case surveycellscan.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case surveycellscan.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case surveycellscan.FieldNetworkType:
		m.ResetNetworkType()
		return nil
	case surveycellscan.FieldSignalStrength:
		m.ResetSignalStrength()
		return nil
	case surveycellscan.FieldTimestamp:
		m.ResetTimestamp()
		return nil
	case surveycellscan.FieldBaseStationID:
		m.ResetBaseStationID()
		return nil
	case surveycellscan.FieldNetworkID:
		m.ResetNetworkID()
		return nil
	case surveycellscan.FieldSystemID:
		m.ResetSystemID()
		return nil
	case surveycellscan.FieldCellID:
		m.ResetCellID()
		return nil
	case surveycellscan.FieldLocationAreaCode:
		m.ResetLocationAreaCode()
		return nil
	case surveycellscan.FieldMobileCountryCode:
		m.ResetMobileCountryCode()
		return nil
	case surveycellscan.FieldMobileNetworkCode:
		m.ResetMobileNetworkCode()
		return nil
	case surveycellscan.FieldPrimaryScramblingCode:
		m.ResetPrimaryScramblingCode()
		return nil
	case surveycellscan.FieldOperator:
		m.ResetOperator()
		return nil
	case surveycellscan.FieldArfcn:
		m.ResetArfcn()
		return nil
	case surveycellscan.FieldPhysicalCellID:
		m.ResetPhysicalCellID()
		return nil
	case surveycellscan.FieldTrackingAreaCode:
		m.ResetTrackingAreaCode()
		return nil
	case surveycellscan.FieldTimingAdvance:
		m.ResetTimingAdvance()
		return nil
	case surveycellscan.FieldEarfcn:
		m.ResetEarfcn()
		return nil
	case surveycellscan.FieldUarfcn:
		m.ResetUarfcn()
		return nil
	case surveycellscan.FieldLatitude:
		m.ResetLatitude()
		return nil
	case surveycellscan.FieldLongitude:
		m.ResetLongitude()
		return nil
	}
	return fmt.Errorf("unknown SurveyCellScan field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *SurveyCellScanMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.survey_question != nil {
		edges = append(edges, surveycellscan.EdgeSurveyQuestion)
	}
	if m.location != nil {
		edges = append(edges, surveycellscan.EdgeLocation)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *SurveyCellScanMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case surveycellscan.EdgeSurveyQuestion:
		if id := m.survey_question; id != nil {
			return []ent.Value{*id}
		}
	case surveycellscan.EdgeLocation:
		if id := m.location; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *SurveyCellScanMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *SurveyCellScanMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *SurveyCellScanMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedsurvey_question {
		edges = append(edges, surveycellscan.EdgeSurveyQuestion)
	}
	if m.clearedlocation {
		edges = append(edges, surveycellscan.EdgeLocation)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *SurveyCellScanMutation) EdgeCleared(name string) bool {
	switch name {
	case surveycellscan.EdgeSurveyQuestion:
		return m.clearedsurvey_question
	case surveycellscan.EdgeLocation:
		return m.clearedlocation
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *SurveyCellScanMutation) ClearEdge(name string) error {
	switch name {
	case surveycellscan.EdgeSurveyQuestion:
		m.ClearSurveyQuestion()
		return nil
	case surveycellscan.EdgeLocation:
		m.ClearLocation()
		return nil
	}
	return fmt.Errorf("unknown SurveyCellScan unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *SurveyCellScanMutation) ResetEdge(name string) error {
	switch name {
	case surveycellscan.EdgeSurveyQuestion:
		m.ResetSurveyQuestion()
		return nil
	case surveycellscan.EdgeLocation:
		m.ResetLocation()
		return nil
	}
	return fmt.Errorf("unknown SurveyCellScan edge %s", name)
}

// SurveyQuestionMutation represents an operation that mutate the SurveyQuestions
// nodes in the graph.
type SurveyQuestionMutation struct {
	config
	op                   Op
	typ                  string
	id                   *int
	create_time          *time.Time
	update_time          *time.Time
	form_name            *string
	form_description     *string
	form_index           *int
	addform_index        *int
	question_type        *string
	question_format      *string
	question_text        *string
	question_index       *int
	addquestion_index    *int
	bool_data            *bool
	email_data           *string
	latitude             *float64
	addlatitude          *float64
	longitude            *float64
	addlongitude         *float64
	location_accuracy    *float64
	addlocation_accuracy *float64
	altitude             *float64
	addaltitude          *float64
	phone_data           *string
	text_data            *string
	float_data           *float64
	addfloat_data        *float64
	int_data             *int
	addint_data          *int
	date_data            *time.Time
	clearedFields        map[string]bool
	survey               *int
	clearedsurvey        bool
	wifi_scan            map[int]struct{}
	removedwifi_scan     map[int]struct{}
	cell_scan            map[int]struct{}
	removedcell_scan     map[int]struct{}
	photo_data           map[int]struct{}
	removedphoto_data    map[int]struct{}
}

var _ ent.Mutation = (*SurveyQuestionMutation)(nil)

// newSurveyQuestionMutation creates new mutation for $n.Name.
func newSurveyQuestionMutation(c config, op Op) *SurveyQuestionMutation {
	return &SurveyQuestionMutation{
		config:        c,
		op:            op,
		typ:           TypeSurveyQuestion,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SurveyQuestionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SurveyQuestionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *SurveyQuestionMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *SurveyQuestionMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *SurveyQuestionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *SurveyQuestionMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *SurveyQuestionMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *SurveyQuestionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *SurveyQuestionMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetFormName sets the form_name field.
func (m *SurveyQuestionMutation) SetFormName(s string) {
	m.form_name = &s
}

// FormName returns the form_name value in the mutation.
func (m *SurveyQuestionMutation) FormName() (r string, exists bool) {
	v := m.form_name
	if v == nil {
		return
	}
	return *v, true
}

// ClearFormName clears the value of form_name.
func (m *SurveyQuestionMutation) ClearFormName() {
	m.form_name = nil
	m.clearedFields[surveyquestion.FieldFormName] = true
}

// FormNameCleared returns if the field form_name was cleared in this mutation.
func (m *SurveyQuestionMutation) FormNameCleared() bool {
	return m.clearedFields[surveyquestion.FieldFormName]
}

// ResetFormName reset all changes of the form_name field.
func (m *SurveyQuestionMutation) ResetFormName() {
	m.form_name = nil
	delete(m.clearedFields, surveyquestion.FieldFormName)
}

// SetFormDescription sets the form_description field.
func (m *SurveyQuestionMutation) SetFormDescription(s string) {
	m.form_description = &s
}

// FormDescription returns the form_description value in the mutation.
func (m *SurveyQuestionMutation) FormDescription() (r string, exists bool) {
	v := m.form_description
	if v == nil {
		return
	}
	return *v, true
}

// ClearFormDescription clears the value of form_description.
func (m *SurveyQuestionMutation) ClearFormDescription() {
	m.form_description = nil
	m.clearedFields[surveyquestion.FieldFormDescription] = true
}

// FormDescriptionCleared returns if the field form_description was cleared in this mutation.
func (m *SurveyQuestionMutation) FormDescriptionCleared() bool {
	return m.clearedFields[surveyquestion.FieldFormDescription]
}

// ResetFormDescription reset all changes of the form_description field.
func (m *SurveyQuestionMutation) ResetFormDescription() {
	m.form_description = nil
	delete(m.clearedFields, surveyquestion.FieldFormDescription)
}

// SetFormIndex sets the form_index field.
func (m *SurveyQuestionMutation) SetFormIndex(i int) {
	m.form_index = &i
	m.addform_index = nil
}

// FormIndex returns the form_index value in the mutation.
func (m *SurveyQuestionMutation) FormIndex() (r int, exists bool) {
	v := m.form_index
	if v == nil {
		return
	}
	return *v, true
}

// AddFormIndex adds i to form_index.
func (m *SurveyQuestionMutation) AddFormIndex(i int) {
	if m.addform_index != nil {
		*m.addform_index += i
	} else {
		m.addform_index = &i
	}
}

// AddedFormIndex returns the value that was added to the form_index field in this mutation.
func (m *SurveyQuestionMutation) AddedFormIndex() (r int, exists bool) {
	v := m.addform_index
	if v == nil {
		return
	}
	return *v, true
}

// ResetFormIndex reset all changes of the form_index field.
func (m *SurveyQuestionMutation) ResetFormIndex() {
	m.form_index = nil
	m.addform_index = nil
}

// SetQuestionType sets the question_type field.
func (m *SurveyQuestionMutation) SetQuestionType(s string) {
	m.question_type = &s
}

// QuestionType returns the question_type value in the mutation.
func (m *SurveyQuestionMutation) QuestionType() (r string, exists bool) {
	v := m.question_type
	if v == nil {
		return
	}
	return *v, true
}

// ClearQuestionType clears the value of question_type.
func (m *SurveyQuestionMutation) ClearQuestionType() {
	m.question_type = nil
	m.clearedFields[surveyquestion.FieldQuestionType] = true
}

// QuestionTypeCleared returns if the field question_type was cleared in this mutation.
func (m *SurveyQuestionMutation) QuestionTypeCleared() bool {
	return m.clearedFields[surveyquestion.FieldQuestionType]
}

// ResetQuestionType reset all changes of the question_type field.
func (m *SurveyQuestionMutation) ResetQuestionType() {
	m.question_type = nil
	delete(m.clearedFields, surveyquestion.FieldQuestionType)
}

// SetQuestionFormat sets the question_format field.
func (m *SurveyQuestionMutation) SetQuestionFormat(s string) {
	m.question_format = &s
}

// QuestionFormat returns the question_format value in the mutation.
func (m *SurveyQuestionMutation) QuestionFormat() (r string, exists bool) {
	v := m.question_format
	if v == nil {
		return
	}
	return *v, true
}

// ClearQuestionFormat clears the value of question_format.
func (m *SurveyQuestionMutation) ClearQuestionFormat() {
	m.question_format = nil
	m.clearedFields[surveyquestion.FieldQuestionFormat] = true
}

// QuestionFormatCleared returns if the field question_format was cleared in this mutation.
func (m *SurveyQuestionMutation) QuestionFormatCleared() bool {
	return m.clearedFields[surveyquestion.FieldQuestionFormat]
}

// ResetQuestionFormat reset all changes of the question_format field.
func (m *SurveyQuestionMutation) ResetQuestionFormat() {
	m.question_format = nil
	delete(m.clearedFields, surveyquestion.FieldQuestionFormat)
}

// SetQuestionText sets the question_text field.
func (m *SurveyQuestionMutation) SetQuestionText(s string) {
	m.question_text = &s
}

// QuestionText returns the question_text value in the mutation.
func (m *SurveyQuestionMutation) QuestionText() (r string, exists bool) {
	v := m.question_text
	if v == nil {
		return
	}
	return *v, true
}

// ClearQuestionText clears the value of question_text.
func (m *SurveyQuestionMutation) ClearQuestionText() {
	m.question_text = nil
	m.clearedFields[surveyquestion.FieldQuestionText] = true
}

// QuestionTextCleared returns if the field question_text was cleared in this mutation.
func (m *SurveyQuestionMutation) QuestionTextCleared() bool {
	return m.clearedFields[surveyquestion.FieldQuestionText]
}

// ResetQuestionText reset all changes of the question_text field.
func (m *SurveyQuestionMutation) ResetQuestionText() {
	m.question_text = nil
	delete(m.clearedFields, surveyquestion.FieldQuestionText)
}

// SetQuestionIndex sets the question_index field.
func (m *SurveyQuestionMutation) SetQuestionIndex(i int) {
	m.question_index = &i
	m.addquestion_index = nil
}

// QuestionIndex returns the question_index value in the mutation.
func (m *SurveyQuestionMutation) QuestionIndex() (r int, exists bool) {
	v := m.question_index
	if v == nil {
		return
	}
	return *v, true
}

// AddQuestionIndex adds i to question_index.
func (m *SurveyQuestionMutation) AddQuestionIndex(i int) {
	if m.addquestion_index != nil {
		*m.addquestion_index += i
	} else {
		m.addquestion_index = &i
	}
}

// AddedQuestionIndex returns the value that was added to the question_index field in this mutation.
func (m *SurveyQuestionMutation) AddedQuestionIndex() (r int, exists bool) {
	v := m.addquestion_index
	if v == nil {
		return
	}
	return *v, true
}

// ResetQuestionIndex reset all changes of the question_index field.
func (m *SurveyQuestionMutation) ResetQuestionIndex() {
	m.question_index = nil
	m.addquestion_index = nil
}

// SetBoolData sets the bool_data field.
func (m *SurveyQuestionMutation) SetBoolData(b bool) {
	m.bool_data = &b
}

// BoolData returns the bool_data value in the mutation.
func (m *SurveyQuestionMutation) BoolData() (r bool, exists bool) {
	v := m.bool_data
	if v == nil {
		return
	}
	return *v, true
}

// ClearBoolData clears the value of bool_data.
func (m *SurveyQuestionMutation) ClearBoolData() {
	m.bool_data = nil
	m.clearedFields[surveyquestion.FieldBoolData] = true
}

// BoolDataCleared returns if the field bool_data was cleared in this mutation.
func (m *SurveyQuestionMutation) BoolDataCleared() bool {
	return m.clearedFields[surveyquestion.FieldBoolData]
}

// ResetBoolData reset all changes of the bool_data field.
func (m *SurveyQuestionMutation) ResetBoolData() {
	m.bool_data = nil
	delete(m.clearedFields, surveyquestion.FieldBoolData)
}

// SetEmailData sets the email_data field.
func (m *SurveyQuestionMutation) SetEmailData(s string) {
	m.email_data = &s
}

// EmailData returns the email_data value in the mutation.
func (m *SurveyQuestionMutation) EmailData() (r string, exists bool) {
	v := m.email_data
	if v == nil {
		return
	}
	return *v, true
}

// ClearEmailData clears the value of email_data.
func (m *SurveyQuestionMutation) ClearEmailData() {
	m.email_data = nil
	m.clearedFields[surveyquestion.FieldEmailData] = true
}

// EmailDataCleared returns if the field email_data was cleared in this mutation.
func (m *SurveyQuestionMutation) EmailDataCleared() bool {
	return m.clearedFields[surveyquestion.FieldEmailData]
}

// ResetEmailData reset all changes of the email_data field.
func (m *SurveyQuestionMutation) ResetEmailData() {
	m.email_data = nil
	delete(m.clearedFields, surveyquestion.FieldEmailData)
}

// SetLatitude sets the latitude field.
func (m *SurveyQuestionMutation) SetLatitude(f float64) {
	m.latitude = &f
	m.addlatitude = nil
}

// Latitude returns the latitude value in the mutation.
func (m *SurveyQuestionMutation) Latitude() (r float64, exists bool) {
	v := m.latitude
	if v == nil {
		return
	}
	return *v, true
}

// AddLatitude adds f to latitude.
func (m *SurveyQuestionMutation) AddLatitude(f float64) {
	if m.addlatitude != nil {
		*m.addlatitude += f
	} else {
		m.addlatitude = &f
	}
}

// AddedLatitude returns the value that was added to the latitude field in this mutation.
func (m *SurveyQuestionMutation) AddedLatitude() (r float64, exists bool) {
	v := m.addlatitude
	if v == nil {
		return
	}
	return *v, true
}

// ClearLatitude clears the value of latitude.
func (m *SurveyQuestionMutation) ClearLatitude() {
	m.latitude = nil
	m.addlatitude = nil
	m.clearedFields[surveyquestion.FieldLatitude] = true
}

// LatitudeCleared returns if the field latitude was cleared in this mutation.
func (m *SurveyQuestionMutation) LatitudeCleared() bool {
	return m.clearedFields[surveyquestion.FieldLatitude]
}

// ResetLatitude reset all changes of the latitude field.
func (m *SurveyQuestionMutation) ResetLatitude() {
	m.latitude = nil
	m.addlatitude = nil
	delete(m.clearedFields, surveyquestion.FieldLatitude)
}

// SetLongitude sets the longitude field.
func (m *SurveyQuestionMutation) SetLongitude(f float64) {
	m.longitude = &f
	m.addlongitude = nil
}

// Longitude returns the longitude value in the mutation.
func (m *SurveyQuestionMutation) Longitude() (r float64, exists bool) {
	v := m.longitude
	if v == nil {
		return
	}
	return *v, true
}

// AddLongitude adds f to longitude.
func (m *SurveyQuestionMutation) AddLongitude(f float64) {
	if m.addlongitude != nil {
		*m.addlongitude += f
	} else {
		m.addlongitude = &f
	}
}

// AddedLongitude returns the value that was added to the longitude field in this mutation.
func (m *SurveyQuestionMutation) AddedLongitude() (r float64, exists bool) {
	v := m.addlongitude
	if v == nil {
		return
	}
	return *v, true
}

// ClearLongitude clears the value of longitude.
func (m *SurveyQuestionMutation) ClearLongitude() {
	m.longitude = nil
	m.addlongitude = nil
	m.clearedFields[surveyquestion.FieldLongitude] = true
}

// LongitudeCleared returns if the field longitude was cleared in this mutation.
func (m *SurveyQuestionMutation) LongitudeCleared() bool {
	return m.clearedFields[surveyquestion.FieldLongitude]
}

// ResetLongitude reset all changes of the longitude field.
func (m *SurveyQuestionMutation) ResetLongitude() {
	m.longitude = nil
	m.addlongitude = nil
	delete(m.clearedFields, surveyquestion.FieldLongitude)
}

// SetLocationAccuracy sets the location_accuracy field.
func (m *SurveyQuestionMutation) SetLocationAccuracy(f float64) {
	m.location_accuracy = &f
	m.addlocation_accuracy = nil
}

// LocationAccuracy returns the location_accuracy value in the mutation.
func (m *SurveyQuestionMutation) LocationAccuracy() (r float64, exists bool) {
	v := m.location_accuracy
	if v == nil {
		return
	}
	return *v, true
}

// AddLocationAccuracy adds f to location_accuracy.
func (m *SurveyQuestionMutation) AddLocationAccuracy(f float64) {
	if m.addlocation_accuracy != nil {
		*m.addlocation_accuracy += f
	} else {
		m.addlocation_accuracy = &f
	}
}

// AddedLocationAccuracy returns the value that was added to the location_accuracy field in this mutation.
func (m *SurveyQuestionMutation) AddedLocationAccuracy() (r float64, exists bool) {
	v := m.addlocation_accuracy
	if v == nil {
		return
	}
	return *v, true
}

// ClearLocationAccuracy clears the value of location_accuracy.
func (m *SurveyQuestionMutation) ClearLocationAccuracy() {
	m.location_accuracy = nil
	m.addlocation_accuracy = nil
	m.clearedFields[surveyquestion.FieldLocationAccuracy] = true
}

// LocationAccuracyCleared returns if the field location_accuracy was cleared in this mutation.
func (m *SurveyQuestionMutation) LocationAccuracyCleared() bool {
	return m.clearedFields[surveyquestion.FieldLocationAccuracy]
}

// ResetLocationAccuracy reset all changes of the location_accuracy field.
func (m *SurveyQuestionMutation) ResetLocationAccuracy() {
	m.location_accuracy = nil
	m.addlocation_accuracy = nil
	delete(m.clearedFields, surveyquestion.FieldLocationAccuracy)
}

// SetAltitude sets the altitude field.
func (m *SurveyQuestionMutation) SetAltitude(f float64) {
	m.altitude = &f
	m.addaltitude = nil
}

// Altitude returns the altitude value in the mutation.
func (m *SurveyQuestionMutation) Altitude() (r float64, exists bool) {
	v := m.altitude
	if v == nil {
		return
	}
	return *v, true
}

// AddAltitude adds f to altitude.
func (m *SurveyQuestionMutation) AddAltitude(f float64) {
	if m.addaltitude != nil {
		*m.addaltitude += f
	} else {
		m.addaltitude = &f
	}
}

// AddedAltitude returns the value that was added to the altitude field in this mutation.
func (m *SurveyQuestionMutation) AddedAltitude() (r float64, exists bool) {
	v := m.addaltitude
	if v == nil {
		return
	}
	return *v, true
}

// ClearAltitude clears the value of altitude.
func (m *SurveyQuestionMutation) ClearAltitude() {
	m.altitude = nil
	m.addaltitude = nil
	m.clearedFields[surveyquestion.FieldAltitude] = true
}

// AltitudeCleared returns if the field altitude was cleared in this mutation.
func (m *SurveyQuestionMutation) AltitudeCleared() bool {
	return m.clearedFields[surveyquestion.FieldAltitude]
}

// ResetAltitude reset all changes of the altitude field.
func (m *SurveyQuestionMutation) ResetAltitude() {
	m.altitude = nil
	m.addaltitude = nil
	delete(m.clearedFields, surveyquestion.FieldAltitude)
}

// SetPhoneData sets the phone_data field.
func (m *SurveyQuestionMutation) SetPhoneData(s string) {
	m.phone_data = &s
}

// PhoneData returns the phone_data value in the mutation.
func (m *SurveyQuestionMutation) PhoneData() (r string, exists bool) {
	v := m.phone_data
	if v == nil {
		return
	}
	return *v, true
}

// ClearPhoneData clears the value of phone_data.
func (m *SurveyQuestionMutation) ClearPhoneData() {
	m.phone_data = nil
	m.clearedFields[surveyquestion.FieldPhoneData] = true
}

// PhoneDataCleared returns if the field phone_data was cleared in this mutation.
func (m *SurveyQuestionMutation) PhoneDataCleared() bool {
	return m.clearedFields[surveyquestion.FieldPhoneData]
}

// ResetPhoneData reset all changes of the phone_data field.
func (m *SurveyQuestionMutation) ResetPhoneData() {
	m.phone_data = nil
	delete(m.clearedFields, surveyquestion.FieldPhoneData)
}

// SetTextData sets the text_data field.
func (m *SurveyQuestionMutation) SetTextData(s string) {
	m.text_data = &s
}

// TextData returns the text_data value in the mutation.
func (m *SurveyQuestionMutation) TextData() (r string, exists bool) {
	v := m.text_data
	if v == nil {
		return
	}
	return *v, true
}

// ClearTextData clears the value of text_data.
func (m *SurveyQuestionMutation) ClearTextData() {
	m.text_data = nil
	m.clearedFields[surveyquestion.FieldTextData] = true
}

// TextDataCleared returns if the field text_data was cleared in this mutation.
func (m *SurveyQuestionMutation) TextDataCleared() bool {
	return m.clearedFields[surveyquestion.FieldTextData]
}

// ResetTextData reset all changes of the text_data field.
func (m *SurveyQuestionMutation) ResetTextData() {
	m.text_data = nil
	delete(m.clearedFields, surveyquestion.FieldTextData)
}

// SetFloatData sets the float_data field.
func (m *SurveyQuestionMutation) SetFloatData(f float64) {
	m.float_data = &f
	m.addfloat_data = nil
}

// FloatData returns the float_data value in the mutation.
func (m *SurveyQuestionMutation) FloatData() (r float64, exists bool) {
	v := m.float_data
	if v == nil {
		return
	}
	return *v, true
}

// AddFloatData adds f to float_data.
func (m *SurveyQuestionMutation) AddFloatData(f float64) {
	if m.addfloat_data != nil {
		*m.addfloat_data += f
	} else {
		m.addfloat_data = &f
	}
}

// AddedFloatData returns the value that was added to the float_data field in this mutation.
func (m *SurveyQuestionMutation) AddedFloatData() (r float64, exists bool) {
	v := m.addfloat_data
	if v == nil {
		return
	}
	return *v, true
}

// ClearFloatData clears the value of float_data.
func (m *SurveyQuestionMutation) ClearFloatData() {
	m.float_data = nil
	m.addfloat_data = nil
	m.clearedFields[surveyquestion.FieldFloatData] = true
}

// FloatDataCleared returns if the field float_data was cleared in this mutation.
func (m *SurveyQuestionMutation) FloatDataCleared() bool {
	return m.clearedFields[surveyquestion.FieldFloatData]
}

// ResetFloatData reset all changes of the float_data field.
func (m *SurveyQuestionMutation) ResetFloatData() {
	m.float_data = nil
	m.addfloat_data = nil
	delete(m.clearedFields, surveyquestion.FieldFloatData)
}

// SetIntData sets the int_data field.
func (m *SurveyQuestionMutation) SetIntData(i int) {
	m.int_data = &i
	m.addint_data = nil
}

// IntData returns the int_data value in the mutation.
func (m *SurveyQuestionMutation) IntData() (r int, exists bool) {
	v := m.int_data
	if v == nil {
		return
	}
	return *v, true
}

// AddIntData adds i to int_data.
func (m *SurveyQuestionMutation) AddIntData(i int) {
	if m.addint_data != nil {
		*m.addint_data += i
	} else {
		m.addint_data = &i
	}
}

// AddedIntData returns the value that was added to the int_data field in this mutation.
func (m *SurveyQuestionMutation) AddedIntData() (r int, exists bool) {
	v := m.addint_data
	if v == nil {
		return
	}
	return *v, true
}

// ClearIntData clears the value of int_data.
func (m *SurveyQuestionMutation) ClearIntData() {
	m.int_data = nil
	m.addint_data = nil
	m.clearedFields[surveyquestion.FieldIntData] = true
}

// IntDataCleared returns if the field int_data was cleared in this mutation.
func (m *SurveyQuestionMutation) IntDataCleared() bool {
	return m.clearedFields[surveyquestion.FieldIntData]
}

// ResetIntData reset all changes of the int_data field.
func (m *SurveyQuestionMutation) ResetIntData() {
	m.int_data = nil
	m.addint_data = nil
	delete(m.clearedFields, surveyquestion.FieldIntData)
}

// SetDateData sets the date_data field.
func (m *SurveyQuestionMutation) SetDateData(t time.Time) {
	m.date_data = &t
}

// DateData returns the date_data value in the mutation.
func (m *SurveyQuestionMutation) DateData() (r time.Time, exists bool) {
	v := m.date_data
	if v == nil {
		return
	}
	return *v, true
}

// ClearDateData clears the value of date_data.
func (m *SurveyQuestionMutation) ClearDateData() {
	m.date_data = nil
	m.clearedFields[surveyquestion.FieldDateData] = true
}

// DateDataCleared returns if the field date_data was cleared in this mutation.
func (m *SurveyQuestionMutation) DateDataCleared() bool {
	return m.clearedFields[surveyquestion.FieldDateData]
}

// ResetDateData reset all changes of the date_data field.
func (m *SurveyQuestionMutation) ResetDateData() {
	m.date_data = nil
	delete(m.clearedFields, surveyquestion.FieldDateData)
}

// SetSurveyID sets the survey edge to Survey by id.
func (m *SurveyQuestionMutation) SetSurveyID(id int) {
	m.survey = &id
}

// ClearSurvey clears the survey edge to Survey.
func (m *SurveyQuestionMutation) ClearSurvey() {
	m.clearedsurvey = true
}

// SurveyCleared returns if the edge survey was cleared.
func (m *SurveyQuestionMutation) SurveyCleared() bool {
	return m.clearedsurvey
}

// SurveyID returns the survey id in the mutation.
func (m *SurveyQuestionMutation) SurveyID() (id int, exists bool) {
	if m.survey != nil {
		return *m.survey, true
	}
	return
}

// SurveyIDs returns the survey ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// SurveyID instead. It exists only for internal usage by the builders.
func (m *SurveyQuestionMutation) SurveyIDs() (ids []int) {
	if id := m.survey; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetSurvey reset all changes of the survey edge.
func (m *SurveyQuestionMutation) ResetSurvey() {
	m.survey = nil
	m.clearedsurvey = false
}

// AddWifiScanIDs adds the wifi_scan edge to SurveyWiFiScan by ids.
func (m *SurveyQuestionMutation) AddWifiScanIDs(ids ...int) {
	if m.wifi_scan == nil {
		m.wifi_scan = make(map[int]struct{})
	}
	for i := range ids {
		m.wifi_scan[ids[i]] = struct{}{}
	}
}

// RemoveWifiScanIDs removes the wifi_scan edge to SurveyWiFiScan by ids.
func (m *SurveyQuestionMutation) RemoveWifiScanIDs(ids ...int) {
	if m.removedwifi_scan == nil {
		m.removedwifi_scan = make(map[int]struct{})
	}
	for i := range ids {
		m.removedwifi_scan[ids[i]] = struct{}{}
	}
}

// RemovedWifiScan returns the removed ids of wifi_scan.
func (m *SurveyQuestionMutation) RemovedWifiScanIDs() (ids []int) {
	for id := range m.removedwifi_scan {
		ids = append(ids, id)
	}
	return
}

// WifiScanIDs returns the wifi_scan ids in the mutation.
func (m *SurveyQuestionMutation) WifiScanIDs() (ids []int) {
	for id := range m.wifi_scan {
		ids = append(ids, id)
	}
	return
}

// ResetWifiScan reset all changes of the wifi_scan edge.
func (m *SurveyQuestionMutation) ResetWifiScan() {
	m.wifi_scan = nil
	m.removedwifi_scan = nil
}

// AddCellScanIDs adds the cell_scan edge to SurveyCellScan by ids.
func (m *SurveyQuestionMutation) AddCellScanIDs(ids ...int) {
	if m.cell_scan == nil {
		m.cell_scan = make(map[int]struct{})
	}
	for i := range ids {
		m.cell_scan[ids[i]] = struct{}{}
	}
}

// RemoveCellScanIDs removes the cell_scan edge to SurveyCellScan by ids.
func (m *SurveyQuestionMutation) RemoveCellScanIDs(ids ...int) {
	if m.removedcell_scan == nil {
		m.removedcell_scan = make(map[int]struct{})
	}
	for i := range ids {
		m.removedcell_scan[ids[i]] = struct{}{}
	}
}

// RemovedCellScan returns the removed ids of cell_scan.
func (m *SurveyQuestionMutation) RemovedCellScanIDs() (ids []int) {
	for id := range m.removedcell_scan {
		ids = append(ids, id)
	}
	return
}

// CellScanIDs returns the cell_scan ids in the mutation.
func (m *SurveyQuestionMutation) CellScanIDs() (ids []int) {
	for id := range m.cell_scan {
		ids = append(ids, id)
	}
	return
}

// ResetCellScan reset all changes of the cell_scan edge.
func (m *SurveyQuestionMutation) ResetCellScan() {
	m.cell_scan = nil
	m.removedcell_scan = nil
}

// AddPhotoDatumIDs adds the photo_data edge to File by ids.
func (m *SurveyQuestionMutation) AddPhotoDatumIDs(ids ...int) {
	if m.photo_data == nil {
		m.photo_data = make(map[int]struct{})
	}
	for i := range ids {
		m.photo_data[ids[i]] = struct{}{}
	}
}

// RemovePhotoDatumIDs removes the photo_data edge to File by ids.
func (m *SurveyQuestionMutation) RemovePhotoDatumIDs(ids ...int) {
	if m.removedphoto_data == nil {
		m.removedphoto_data = make(map[int]struct{})
	}
	for i := range ids {
		m.removedphoto_data[ids[i]] = struct{}{}
	}
}

// RemovedPhotoData returns the removed ids of photo_data.
func (m *SurveyQuestionMutation) RemovedPhotoDataIDs() (ids []int) {
	for id := range m.removedphoto_data {
		ids = append(ids, id)
	}
	return
}

// PhotoDataIDs returns the photo_data ids in the mutation.
func (m *SurveyQuestionMutation) PhotoDataIDs() (ids []int) {
	for id := range m.photo_data {
		ids = append(ids, id)
	}
	return
}

// ResetPhotoData reset all changes of the photo_data edge.
func (m *SurveyQuestionMutation) ResetPhotoData() {
	m.photo_data = nil
	m.removedphoto_data = nil
}

// Op returns the operation name.
func (m *SurveyQuestionMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (SurveyQuestion).
func (m *SurveyQuestionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *SurveyQuestionMutation) Fields() []string {
	fields := make([]string, 0, 20)
	if m.create_time != nil {
		fields = append(fields, surveyquestion.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, surveyquestion.FieldUpdateTime)
	}
	if m.form_name != nil {
		fields = append(fields, surveyquestion.FieldFormName)
	}
	if m.form_description != nil {
		fields = append(fields, surveyquestion.FieldFormDescription)
	}
	if m.form_index != nil {
		fields = append(fields, surveyquestion.FieldFormIndex)
	}
	if m.question_type != nil {
		fields = append(fields, surveyquestion.FieldQuestionType)
	}
	if m.question_format != nil {
		fields = append(fields, surveyquestion.FieldQuestionFormat)
	}
	if m.question_text != nil {
		fields = append(fields, surveyquestion.FieldQuestionText)
	}
	if m.question_index != nil {
		fields = append(fields, surveyquestion.FieldQuestionIndex)
	}
	if m.bool_data != nil {
		fields = append(fields, surveyquestion.FieldBoolData)
	}
	if m.email_data != nil {
		fields = append(fields, surveyquestion.FieldEmailData)
	}
	if m.latitude != nil {
		fields = append(fields, surveyquestion.FieldLatitude)
	}
	if m.longitude != nil {
		fields = append(fields, surveyquestion.FieldLongitude)
	}
	if m.location_accuracy != nil {
		fields = append(fields, surveyquestion.FieldLocationAccuracy)
	}
	if m.altitude != nil {
		fields = append(fields, surveyquestion.FieldAltitude)
	}
	if m.phone_data != nil {
		fields = append(fields, surveyquestion.FieldPhoneData)
	}
	if m.text_data != nil {
		fields = append(fields, surveyquestion.FieldTextData)
	}
	if m.float_data != nil {
		fields = append(fields, surveyquestion.FieldFloatData)
	}
	if m.int_data != nil {
		fields = append(fields, surveyquestion.FieldIntData)
	}
	if m.date_data != nil {
		fields = append(fields, surveyquestion.FieldDateData)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *SurveyQuestionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case surveyquestion.FieldCreateTime:
		return m.CreateTime()
	case surveyquestion.FieldUpdateTime:
		return m.UpdateTime()
	case surveyquestion.FieldFormName:
		return m.FormName()
	case surveyquestion.FieldFormDescription:
		return m.FormDescription()
	case surveyquestion.FieldFormIndex:
		return m.FormIndex()
	case surveyquestion.FieldQuestionType:
		return m.QuestionType()
	case surveyquestion.FieldQuestionFormat:
		return m.QuestionFormat()
	case surveyquestion.FieldQuestionText:
		return m.QuestionText()
	case surveyquestion.FieldQuestionIndex:
		return m.QuestionIndex()
	case surveyquestion.FieldBoolData:
		return m.BoolData()
	case surveyquestion.FieldEmailData:
		return m.EmailData()
	case surveyquestion.FieldLatitude:
		return m.Latitude()
	case surveyquestion.FieldLongitude:
		return m.Longitude()
	case surveyquestion.FieldLocationAccuracy:
		return m.LocationAccuracy()
	case surveyquestion.FieldAltitude:
		return m.Altitude()
	case surveyquestion.FieldPhoneData:
		return m.PhoneData()
	case surveyquestion.FieldTextData:
		return m.TextData()
	case surveyquestion.FieldFloatData:
		return m.FloatData()
	case surveyquestion.FieldIntData:
		return m.IntData()
	case surveyquestion.FieldDateData:
		return m.DateData()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyQuestionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case surveyquestion.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case surveyquestion.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case surveyquestion.FieldFormName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFormName(v)
		return nil
	case surveyquestion.FieldFormDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFormDescription(v)
		return nil
	case surveyquestion.FieldFormIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFormIndex(v)
		return nil
	case surveyquestion.FieldQuestionType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetQuestionType(v)
		return nil
	case surveyquestion.FieldQuestionFormat:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetQuestionFormat(v)
		return nil
	case surveyquestion.FieldQuestionText:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetQuestionText(v)
		return nil
	case surveyquestion.FieldQuestionIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetQuestionIndex(v)
		return nil
	case surveyquestion.FieldBoolData:
		v, ok := value.(bool)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBoolData(v)
		return nil
	case surveyquestion.FieldEmailData:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEmailData(v)
		return nil
	case surveyquestion.FieldLatitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLatitude(v)
		return nil
	case surveyquestion.FieldLongitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLongitude(v)
		return nil
	case surveyquestion.FieldLocationAccuracy:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLocationAccuracy(v)
		return nil
	case surveyquestion.FieldAltitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetAltitude(v)
		return nil
	case surveyquestion.FieldPhoneData:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPhoneData(v)
		return nil
	case surveyquestion.FieldTextData:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTextData(v)
		return nil
	case surveyquestion.FieldFloatData:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFloatData(v)
		return nil
	case surveyquestion.FieldIntData:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIntData(v)
		return nil
	case surveyquestion.FieldDateData:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDateData(v)
		return nil
	}
	return fmt.Errorf("unknown SurveyQuestion field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *SurveyQuestionMutation) AddedFields() []string {
	var fields []string
	if m.addform_index != nil {
		fields = append(fields, surveyquestion.FieldFormIndex)
	}
	if m.addquestion_index != nil {
		fields = append(fields, surveyquestion.FieldQuestionIndex)
	}
	if m.addlatitude != nil {
		fields = append(fields, surveyquestion.FieldLatitude)
	}
	if m.addlongitude != nil {
		fields = append(fields, surveyquestion.FieldLongitude)
	}
	if m.addlocation_accuracy != nil {
		fields = append(fields, surveyquestion.FieldLocationAccuracy)
	}
	if m.addaltitude != nil {
		fields = append(fields, surveyquestion.FieldAltitude)
	}
	if m.addfloat_data != nil {
		fields = append(fields, surveyquestion.FieldFloatData)
	}
	if m.addint_data != nil {
		fields = append(fields, surveyquestion.FieldIntData)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *SurveyQuestionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case surveyquestion.FieldFormIndex:
		return m.AddedFormIndex()
	case surveyquestion.FieldQuestionIndex:
		return m.AddedQuestionIndex()
	case surveyquestion.FieldLatitude:
		return m.AddedLatitude()
	case surveyquestion.FieldLongitude:
		return m.AddedLongitude()
	case surveyquestion.FieldLocationAccuracy:
		return m.AddedLocationAccuracy()
	case surveyquestion.FieldAltitude:
		return m.AddedAltitude()
	case surveyquestion.FieldFloatData:
		return m.AddedFloatData()
	case surveyquestion.FieldIntData:
		return m.AddedIntData()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyQuestionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case surveyquestion.FieldFormIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddFormIndex(v)
		return nil
	case surveyquestion.FieldQuestionIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddQuestionIndex(v)
		return nil
	case surveyquestion.FieldLatitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLatitude(v)
		return nil
	case surveyquestion.FieldLongitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLongitude(v)
		return nil
	case surveyquestion.FieldLocationAccuracy:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLocationAccuracy(v)
		return nil
	case surveyquestion.FieldAltitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddAltitude(v)
		return nil
	case surveyquestion.FieldFloatData:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddFloatData(v)
		return nil
	case surveyquestion.FieldIntData:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIntData(v)
		return nil
	}
	return fmt.Errorf("unknown SurveyQuestion numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *SurveyQuestionMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[surveyquestion.FieldFormName] {
		fields = append(fields, surveyquestion.FieldFormName)
	}
	if m.clearedFields[surveyquestion.FieldFormDescription] {
		fields = append(fields, surveyquestion.FieldFormDescription)
	}
	if m.clearedFields[surveyquestion.FieldQuestionType] {
		fields = append(fields, surveyquestion.FieldQuestionType)
	}
	if m.clearedFields[surveyquestion.FieldQuestionFormat] {
		fields = append(fields, surveyquestion.FieldQuestionFormat)
	}
	if m.clearedFields[surveyquestion.FieldQuestionText] {
		fields = append(fields, surveyquestion.FieldQuestionText)
	}
	if m.clearedFields[surveyquestion.FieldBoolData] {
		fields = append(fields, surveyquestion.FieldBoolData)
	}
	if m.clearedFields[surveyquestion.FieldEmailData] {
		fields = append(fields, surveyquestion.FieldEmailData)
	}
	if m.clearedFields[surveyquestion.FieldLatitude] {
		fields = append(fields, surveyquestion.FieldLatitude)
	}
	if m.clearedFields[surveyquestion.FieldLongitude] {
		fields = append(fields, surveyquestion.FieldLongitude)
	}
	if m.clearedFields[surveyquestion.FieldLocationAccuracy] {
		fields = append(fields, surveyquestion.FieldLocationAccuracy)
	}
	if m.clearedFields[surveyquestion.FieldAltitude] {
		fields = append(fields, surveyquestion.FieldAltitude)
	}
	if m.clearedFields[surveyquestion.FieldPhoneData] {
		fields = append(fields, surveyquestion.FieldPhoneData)
	}
	if m.clearedFields[surveyquestion.FieldTextData] {
		fields = append(fields, surveyquestion.FieldTextData)
	}
	if m.clearedFields[surveyquestion.FieldFloatData] {
		fields = append(fields, surveyquestion.FieldFloatData)
	}
	if m.clearedFields[surveyquestion.FieldIntData] {
		fields = append(fields, surveyquestion.FieldIntData)
	}
	if m.clearedFields[surveyquestion.FieldDateData] {
		fields = append(fields, surveyquestion.FieldDateData)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *SurveyQuestionMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *SurveyQuestionMutation) ClearField(name string) error {
	switch name {
	case surveyquestion.FieldFormName:
		m.ClearFormName()
		return nil
	case surveyquestion.FieldFormDescription:
		m.ClearFormDescription()
		return nil
	case surveyquestion.FieldQuestionType:
		m.ClearQuestionType()
		return nil
	case surveyquestion.FieldQuestionFormat:
		m.ClearQuestionFormat()
		return nil
	case surveyquestion.FieldQuestionText:
		m.ClearQuestionText()
		return nil
	case surveyquestion.FieldBoolData:
		m.ClearBoolData()
		return nil
	case surveyquestion.FieldEmailData:
		m.ClearEmailData()
		return nil
	case surveyquestion.FieldLatitude:
		m.ClearLatitude()
		return nil
	case surveyquestion.FieldLongitude:
		m.ClearLongitude()
		return nil
	case surveyquestion.FieldLocationAccuracy:
		m.ClearLocationAccuracy()
		return nil
	case surveyquestion.FieldAltitude:
		m.ClearAltitude()
		return nil
	case surveyquestion.FieldPhoneData:
		m.ClearPhoneData()
		return nil
	case surveyquestion.FieldTextData:
		m.ClearTextData()
		return nil
	case surveyquestion.FieldFloatData:
		m.ClearFloatData()
		return nil
	case surveyquestion.FieldIntData:
		m.ClearIntData()
		return nil
	case surveyquestion.FieldDateData:
		m.ClearDateData()
		return nil
	}
	return fmt.Errorf("unknown SurveyQuestion nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *SurveyQuestionMutation) ResetField(name string) error {
	switch name {
	case surveyquestion.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case surveyquestion.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case surveyquestion.FieldFormName:
		m.ResetFormName()
		return nil
	case surveyquestion.FieldFormDescription:
		m.ResetFormDescription()
		return nil
	case surveyquestion.FieldFormIndex:
		m.ResetFormIndex()
		return nil
	case surveyquestion.FieldQuestionType:
		m.ResetQuestionType()
		return nil
	case surveyquestion.FieldQuestionFormat:
		m.ResetQuestionFormat()
		return nil
	case surveyquestion.FieldQuestionText:
		m.ResetQuestionText()
		return nil
	case surveyquestion.FieldQuestionIndex:
		m.ResetQuestionIndex()
		return nil
	case surveyquestion.FieldBoolData:
		m.ResetBoolData()
		return nil
	case surveyquestion.FieldEmailData:
		m.ResetEmailData()
		return nil
	case surveyquestion.FieldLatitude:
		m.ResetLatitude()
		return nil
	case surveyquestion.FieldLongitude:
		m.ResetLongitude()
		return nil
	case surveyquestion.FieldLocationAccuracy:
		m.ResetLocationAccuracy()
		return nil
	case surveyquestion.FieldAltitude:
		m.ResetAltitude()
		return nil
	case surveyquestion.FieldPhoneData:
		m.ResetPhoneData()
		return nil
	case surveyquestion.FieldTextData:
		m.ResetTextData()
		return nil
	case surveyquestion.FieldFloatData:
		m.ResetFloatData()
		return nil
	case surveyquestion.FieldIntData:
		m.ResetIntData()
		return nil
	case surveyquestion.FieldDateData:
		m.ResetDateData()
		return nil
	}
	return fmt.Errorf("unknown SurveyQuestion field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *SurveyQuestionMutation) AddedEdges() []string {
	edges := make([]string, 0, 4)
	if m.survey != nil {
		edges = append(edges, surveyquestion.EdgeSurvey)
	}
	if m.wifi_scan != nil {
		edges = append(edges, surveyquestion.EdgeWifiScan)
	}
	if m.cell_scan != nil {
		edges = append(edges, surveyquestion.EdgeCellScan)
	}
	if m.photo_data != nil {
		edges = append(edges, surveyquestion.EdgePhotoData)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *SurveyQuestionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case surveyquestion.EdgeSurvey:
		if id := m.survey; id != nil {
			return []ent.Value{*id}
		}
	case surveyquestion.EdgeWifiScan:
		ids := make([]ent.Value, 0, len(m.wifi_scan))
		for id := range m.wifi_scan {
			ids = append(ids, id)
		}
		return ids
	case surveyquestion.EdgeCellScan:
		ids := make([]ent.Value, 0, len(m.cell_scan))
		for id := range m.cell_scan {
			ids = append(ids, id)
		}
		return ids
	case surveyquestion.EdgePhotoData:
		ids := make([]ent.Value, 0, len(m.photo_data))
		for id := range m.photo_data {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *SurveyQuestionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 4)
	if m.removedwifi_scan != nil {
		edges = append(edges, surveyquestion.EdgeWifiScan)
	}
	if m.removedcell_scan != nil {
		edges = append(edges, surveyquestion.EdgeCellScan)
	}
	if m.removedphoto_data != nil {
		edges = append(edges, surveyquestion.EdgePhotoData)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *SurveyQuestionMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case surveyquestion.EdgeWifiScan:
		ids := make([]ent.Value, 0, len(m.removedwifi_scan))
		for id := range m.removedwifi_scan {
			ids = append(ids, id)
		}
		return ids
	case surveyquestion.EdgeCellScan:
		ids := make([]ent.Value, 0, len(m.removedcell_scan))
		for id := range m.removedcell_scan {
			ids = append(ids, id)
		}
		return ids
	case surveyquestion.EdgePhotoData:
		ids := make([]ent.Value, 0, len(m.removedphoto_data))
		for id := range m.removedphoto_data {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *SurveyQuestionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 4)
	if m.clearedsurvey {
		edges = append(edges, surveyquestion.EdgeSurvey)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *SurveyQuestionMutation) EdgeCleared(name string) bool {
	switch name {
	case surveyquestion.EdgeSurvey:
		return m.clearedsurvey
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *SurveyQuestionMutation) ClearEdge(name string) error {
	switch name {
	case surveyquestion.EdgeSurvey:
		m.ClearSurvey()
		return nil
	}
	return fmt.Errorf("unknown SurveyQuestion unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *SurveyQuestionMutation) ResetEdge(name string) error {
	switch name {
	case surveyquestion.EdgeSurvey:
		m.ResetSurvey()
		return nil
	case surveyquestion.EdgeWifiScan:
		m.ResetWifiScan()
		return nil
	case surveyquestion.EdgeCellScan:
		m.ResetCellScan()
		return nil
	case surveyquestion.EdgePhotoData:
		m.ResetPhotoData()
		return nil
	}
	return fmt.Errorf("unknown SurveyQuestion edge %s", name)
}

// SurveyTemplateCategoryMutation represents an operation that mutate the SurveyTemplateCategories
// nodes in the graph.
type SurveyTemplateCategoryMutation struct {
	config
	op                               Op
	typ                              string
	id                               *int
	create_time                      *time.Time
	update_time                      *time.Time
	category_title                   *string
	category_description             *string
	clearedFields                    map[string]bool
	survey_template_questions        map[int]struct{}
	removedsurvey_template_questions map[int]struct{}
}

var _ ent.Mutation = (*SurveyTemplateCategoryMutation)(nil)

// newSurveyTemplateCategoryMutation creates new mutation for $n.Name.
func newSurveyTemplateCategoryMutation(c config, op Op) *SurveyTemplateCategoryMutation {
	return &SurveyTemplateCategoryMutation{
		config:        c,
		op:            op,
		typ:           TypeSurveyTemplateCategory,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SurveyTemplateCategoryMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SurveyTemplateCategoryMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *SurveyTemplateCategoryMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *SurveyTemplateCategoryMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *SurveyTemplateCategoryMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *SurveyTemplateCategoryMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *SurveyTemplateCategoryMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *SurveyTemplateCategoryMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *SurveyTemplateCategoryMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetCategoryTitle sets the category_title field.
func (m *SurveyTemplateCategoryMutation) SetCategoryTitle(s string) {
	m.category_title = &s
}

// CategoryTitle returns the category_title value in the mutation.
func (m *SurveyTemplateCategoryMutation) CategoryTitle() (r string, exists bool) {
	v := m.category_title
	if v == nil {
		return
	}
	return *v, true
}

// ResetCategoryTitle reset all changes of the category_title field.
func (m *SurveyTemplateCategoryMutation) ResetCategoryTitle() {
	m.category_title = nil
}

// SetCategoryDescription sets the category_description field.
func (m *SurveyTemplateCategoryMutation) SetCategoryDescription(s string) {
	m.category_description = &s
}

// CategoryDescription returns the category_description value in the mutation.
func (m *SurveyTemplateCategoryMutation) CategoryDescription() (r string, exists bool) {
	v := m.category_description
	if v == nil {
		return
	}
	return *v, true
}

// ResetCategoryDescription reset all changes of the category_description field.
func (m *SurveyTemplateCategoryMutation) ResetCategoryDescription() {
	m.category_description = nil
}

// AddSurveyTemplateQuestionIDs adds the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (m *SurveyTemplateCategoryMutation) AddSurveyTemplateQuestionIDs(ids ...int) {
	if m.survey_template_questions == nil {
		m.survey_template_questions = make(map[int]struct{})
	}
	for i := range ids {
		m.survey_template_questions[ids[i]] = struct{}{}
	}
}

// RemoveSurveyTemplateQuestionIDs removes the survey_template_questions edge to SurveyTemplateQuestion by ids.
func (m *SurveyTemplateCategoryMutation) RemoveSurveyTemplateQuestionIDs(ids ...int) {
	if m.removedsurvey_template_questions == nil {
		m.removedsurvey_template_questions = make(map[int]struct{})
	}
	for i := range ids {
		m.removedsurvey_template_questions[ids[i]] = struct{}{}
	}
}

// RemovedSurveyTemplateQuestions returns the removed ids of survey_template_questions.
func (m *SurveyTemplateCategoryMutation) RemovedSurveyTemplateQuestionsIDs() (ids []int) {
	for id := range m.removedsurvey_template_questions {
		ids = append(ids, id)
	}
	return
}

// SurveyTemplateQuestionsIDs returns the survey_template_questions ids in the mutation.
func (m *SurveyTemplateCategoryMutation) SurveyTemplateQuestionsIDs() (ids []int) {
	for id := range m.survey_template_questions {
		ids = append(ids, id)
	}
	return
}

// ResetSurveyTemplateQuestions reset all changes of the survey_template_questions edge.
func (m *SurveyTemplateCategoryMutation) ResetSurveyTemplateQuestions() {
	m.survey_template_questions = nil
	m.removedsurvey_template_questions = nil
}

// Op returns the operation name.
func (m *SurveyTemplateCategoryMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (SurveyTemplateCategory).
func (m *SurveyTemplateCategoryMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *SurveyTemplateCategoryMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.create_time != nil {
		fields = append(fields, surveytemplatecategory.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, surveytemplatecategory.FieldUpdateTime)
	}
	if m.category_title != nil {
		fields = append(fields, surveytemplatecategory.FieldCategoryTitle)
	}
	if m.category_description != nil {
		fields = append(fields, surveytemplatecategory.FieldCategoryDescription)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *SurveyTemplateCategoryMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case surveytemplatecategory.FieldCreateTime:
		return m.CreateTime()
	case surveytemplatecategory.FieldUpdateTime:
		return m.UpdateTime()
	case surveytemplatecategory.FieldCategoryTitle:
		return m.CategoryTitle()
	case surveytemplatecategory.FieldCategoryDescription:
		return m.CategoryDescription()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyTemplateCategoryMutation) SetField(name string, value ent.Value) error {
	switch name {
	case surveytemplatecategory.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case surveytemplatecategory.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case surveytemplatecategory.FieldCategoryTitle:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCategoryTitle(v)
		return nil
	case surveytemplatecategory.FieldCategoryDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCategoryDescription(v)
		return nil
	}
	return fmt.Errorf("unknown SurveyTemplateCategory field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *SurveyTemplateCategoryMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *SurveyTemplateCategoryMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyTemplateCategoryMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown SurveyTemplateCategory numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *SurveyTemplateCategoryMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *SurveyTemplateCategoryMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *SurveyTemplateCategoryMutation) ClearField(name string) error {
	return fmt.Errorf("unknown SurveyTemplateCategory nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *SurveyTemplateCategoryMutation) ResetField(name string) error {
	switch name {
	case surveytemplatecategory.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case surveytemplatecategory.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case surveytemplatecategory.FieldCategoryTitle:
		m.ResetCategoryTitle()
		return nil
	case surveytemplatecategory.FieldCategoryDescription:
		m.ResetCategoryDescription()
		return nil
	}
	return fmt.Errorf("unknown SurveyTemplateCategory field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *SurveyTemplateCategoryMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.survey_template_questions != nil {
		edges = append(edges, surveytemplatecategory.EdgeSurveyTemplateQuestions)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *SurveyTemplateCategoryMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case surveytemplatecategory.EdgeSurveyTemplateQuestions:
		ids := make([]ent.Value, 0, len(m.survey_template_questions))
		for id := range m.survey_template_questions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *SurveyTemplateCategoryMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedsurvey_template_questions != nil {
		edges = append(edges, surveytemplatecategory.EdgeSurveyTemplateQuestions)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *SurveyTemplateCategoryMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case surveytemplatecategory.EdgeSurveyTemplateQuestions:
		ids := make([]ent.Value, 0, len(m.removedsurvey_template_questions))
		for id := range m.removedsurvey_template_questions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *SurveyTemplateCategoryMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *SurveyTemplateCategoryMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *SurveyTemplateCategoryMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown SurveyTemplateCategory unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *SurveyTemplateCategoryMutation) ResetEdge(name string) error {
	switch name {
	case surveytemplatecategory.EdgeSurveyTemplateQuestions:
		m.ResetSurveyTemplateQuestions()
		return nil
	}
	return fmt.Errorf("unknown SurveyTemplateCategory edge %s", name)
}

// SurveyTemplateQuestionMutation represents an operation that mutate the SurveyTemplateQuestions
// nodes in the graph.
type SurveyTemplateQuestionMutation struct {
	config
	op                   Op
	typ                  string
	id                   *int
	create_time          *time.Time
	update_time          *time.Time
	question_title       *string
	question_description *string
	question_type        *string
	index                *int
	addindex             *int
	clearedFields        map[string]bool
	category             *int
	clearedcategory      bool
}

var _ ent.Mutation = (*SurveyTemplateQuestionMutation)(nil)

// newSurveyTemplateQuestionMutation creates new mutation for $n.Name.
func newSurveyTemplateQuestionMutation(c config, op Op) *SurveyTemplateQuestionMutation {
	return &SurveyTemplateQuestionMutation{
		config:        c,
		op:            op,
		typ:           TypeSurveyTemplateQuestion,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SurveyTemplateQuestionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SurveyTemplateQuestionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *SurveyTemplateQuestionMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *SurveyTemplateQuestionMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *SurveyTemplateQuestionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *SurveyTemplateQuestionMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *SurveyTemplateQuestionMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *SurveyTemplateQuestionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *SurveyTemplateQuestionMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetQuestionTitle sets the question_title field.
func (m *SurveyTemplateQuestionMutation) SetQuestionTitle(s string) {
	m.question_title = &s
}

// QuestionTitle returns the question_title value in the mutation.
func (m *SurveyTemplateQuestionMutation) QuestionTitle() (r string, exists bool) {
	v := m.question_title
	if v == nil {
		return
	}
	return *v, true
}

// ResetQuestionTitle reset all changes of the question_title field.
func (m *SurveyTemplateQuestionMutation) ResetQuestionTitle() {
	m.question_title = nil
}

// SetQuestionDescription sets the question_description field.
func (m *SurveyTemplateQuestionMutation) SetQuestionDescription(s string) {
	m.question_description = &s
}

// QuestionDescription returns the question_description value in the mutation.
func (m *SurveyTemplateQuestionMutation) QuestionDescription() (r string, exists bool) {
	v := m.question_description
	if v == nil {
		return
	}
	return *v, true
}

// ResetQuestionDescription reset all changes of the question_description field.
func (m *SurveyTemplateQuestionMutation) ResetQuestionDescription() {
	m.question_description = nil
}

// SetQuestionType sets the question_type field.
func (m *SurveyTemplateQuestionMutation) SetQuestionType(s string) {
	m.question_type = &s
}

// QuestionType returns the question_type value in the mutation.
func (m *SurveyTemplateQuestionMutation) QuestionType() (r string, exists bool) {
	v := m.question_type
	if v == nil {
		return
	}
	return *v, true
}

// ResetQuestionType reset all changes of the question_type field.
func (m *SurveyTemplateQuestionMutation) ResetQuestionType() {
	m.question_type = nil
}

// SetIndex sets the index field.
func (m *SurveyTemplateQuestionMutation) SetIndex(i int) {
	m.index = &i
	m.addindex = nil
}

// Index returns the index value in the mutation.
func (m *SurveyTemplateQuestionMutation) Index() (r int, exists bool) {
	v := m.index
	if v == nil {
		return
	}
	return *v, true
}

// AddIndex adds i to index.
func (m *SurveyTemplateQuestionMutation) AddIndex(i int) {
	if m.addindex != nil {
		*m.addindex += i
	} else {
		m.addindex = &i
	}
}

// AddedIndex returns the value that was added to the index field in this mutation.
func (m *SurveyTemplateQuestionMutation) AddedIndex() (r int, exists bool) {
	v := m.addindex
	if v == nil {
		return
	}
	return *v, true
}

// ResetIndex reset all changes of the index field.
func (m *SurveyTemplateQuestionMutation) ResetIndex() {
	m.index = nil
	m.addindex = nil
}

// SetCategoryID sets the category edge to SurveyTemplateCategory by id.
func (m *SurveyTemplateQuestionMutation) SetCategoryID(id int) {
	m.category = &id
}

// ClearCategory clears the category edge to SurveyTemplateCategory.
func (m *SurveyTemplateQuestionMutation) ClearCategory() {
	m.clearedcategory = true
}

// CategoryCleared returns if the edge category was cleared.
func (m *SurveyTemplateQuestionMutation) CategoryCleared() bool {
	return m.clearedcategory
}

// CategoryID returns the category id in the mutation.
func (m *SurveyTemplateQuestionMutation) CategoryID() (id int, exists bool) {
	if m.category != nil {
		return *m.category, true
	}
	return
}

// CategoryIDs returns the category ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// CategoryID instead. It exists only for internal usage by the builders.
func (m *SurveyTemplateQuestionMutation) CategoryIDs() (ids []int) {
	if id := m.category; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetCategory reset all changes of the category edge.
func (m *SurveyTemplateQuestionMutation) ResetCategory() {
	m.category = nil
	m.clearedcategory = false
}

// Op returns the operation name.
func (m *SurveyTemplateQuestionMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (SurveyTemplateQuestion).
func (m *SurveyTemplateQuestionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *SurveyTemplateQuestionMutation) Fields() []string {
	fields := make([]string, 0, 6)
	if m.create_time != nil {
		fields = append(fields, surveytemplatequestion.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, surveytemplatequestion.FieldUpdateTime)
	}
	if m.question_title != nil {
		fields = append(fields, surveytemplatequestion.FieldQuestionTitle)
	}
	if m.question_description != nil {
		fields = append(fields, surveytemplatequestion.FieldQuestionDescription)
	}
	if m.question_type != nil {
		fields = append(fields, surveytemplatequestion.FieldQuestionType)
	}
	if m.index != nil {
		fields = append(fields, surveytemplatequestion.FieldIndex)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *SurveyTemplateQuestionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case surveytemplatequestion.FieldCreateTime:
		return m.CreateTime()
	case surveytemplatequestion.FieldUpdateTime:
		return m.UpdateTime()
	case surveytemplatequestion.FieldQuestionTitle:
		return m.QuestionTitle()
	case surveytemplatequestion.FieldQuestionDescription:
		return m.QuestionDescription()
	case surveytemplatequestion.FieldQuestionType:
		return m.QuestionType()
	case surveytemplatequestion.FieldIndex:
		return m.Index()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyTemplateQuestionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case surveytemplatequestion.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case surveytemplatequestion.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case surveytemplatequestion.FieldQuestionTitle:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetQuestionTitle(v)
		return nil
	case surveytemplatequestion.FieldQuestionDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetQuestionDescription(v)
		return nil
	case surveytemplatequestion.FieldQuestionType:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetQuestionType(v)
		return nil
	case surveytemplatequestion.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIndex(v)
		return nil
	}
	return fmt.Errorf("unknown SurveyTemplateQuestion field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *SurveyTemplateQuestionMutation) AddedFields() []string {
	var fields []string
	if m.addindex != nil {
		fields = append(fields, surveytemplatequestion.FieldIndex)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *SurveyTemplateQuestionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case surveytemplatequestion.FieldIndex:
		return m.AddedIndex()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyTemplateQuestionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case surveytemplatequestion.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIndex(v)
		return nil
	}
	return fmt.Errorf("unknown SurveyTemplateQuestion numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *SurveyTemplateQuestionMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *SurveyTemplateQuestionMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *SurveyTemplateQuestionMutation) ClearField(name string) error {
	return fmt.Errorf("unknown SurveyTemplateQuestion nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *SurveyTemplateQuestionMutation) ResetField(name string) error {
	switch name {
	case surveytemplatequestion.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case surveytemplatequestion.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case surveytemplatequestion.FieldQuestionTitle:
		m.ResetQuestionTitle()
		return nil
	case surveytemplatequestion.FieldQuestionDescription:
		m.ResetQuestionDescription()
		return nil
	case surveytemplatequestion.FieldQuestionType:
		m.ResetQuestionType()
		return nil
	case surveytemplatequestion.FieldIndex:
		m.ResetIndex()
		return nil
	}
	return fmt.Errorf("unknown SurveyTemplateQuestion field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *SurveyTemplateQuestionMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.category != nil {
		edges = append(edges, surveytemplatequestion.EdgeCategory)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *SurveyTemplateQuestionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case surveytemplatequestion.EdgeCategory:
		if id := m.category; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *SurveyTemplateQuestionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *SurveyTemplateQuestionMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *SurveyTemplateQuestionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedcategory {
		edges = append(edges, surveytemplatequestion.EdgeCategory)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *SurveyTemplateQuestionMutation) EdgeCleared(name string) bool {
	switch name {
	case surveytemplatequestion.EdgeCategory:
		return m.clearedcategory
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *SurveyTemplateQuestionMutation) ClearEdge(name string) error {
	switch name {
	case surveytemplatequestion.EdgeCategory:
		m.ClearCategory()
		return nil
	}
	return fmt.Errorf("unknown SurveyTemplateQuestion unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *SurveyTemplateQuestionMutation) ResetEdge(name string) error {
	switch name {
	case surveytemplatequestion.EdgeCategory:
		m.ResetCategory()
		return nil
	}
	return fmt.Errorf("unknown SurveyTemplateQuestion edge %s", name)
}

// SurveyWiFiScanMutation represents an operation that mutate the SurveyWiFiScans
// nodes in the graph.
type SurveyWiFiScanMutation struct {
	config
	op                     Op
	typ                    string
	id                     *int
	create_time            *time.Time
	update_time            *time.Time
	ssid                   *string
	bssid                  *string
	timestamp              *time.Time
	frequency              *int
	addfrequency           *int
	channel                *int
	addchannel             *int
	band                   *string
	channel_width          *int
	addchannel_width       *int
	capabilities           *string
	strength               *int
	addstrength            *int
	latitude               *float64
	addlatitude            *float64
	longitude              *float64
	addlongitude           *float64
	clearedFields          map[string]bool
	survey_question        *int
	clearedsurvey_question bool
	location               *int
	clearedlocation        bool
}

var _ ent.Mutation = (*SurveyWiFiScanMutation)(nil)

// newSurveyWiFiScanMutation creates new mutation for $n.Name.
func newSurveyWiFiScanMutation(c config, op Op) *SurveyWiFiScanMutation {
	return &SurveyWiFiScanMutation{
		config:        c,
		op:            op,
		typ:           TypeSurveyWiFiScan,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SurveyWiFiScanMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SurveyWiFiScanMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *SurveyWiFiScanMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *SurveyWiFiScanMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *SurveyWiFiScanMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *SurveyWiFiScanMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *SurveyWiFiScanMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *SurveyWiFiScanMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *SurveyWiFiScanMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetSsid sets the ssid field.
func (m *SurveyWiFiScanMutation) SetSsid(s string) {
	m.ssid = &s
}

// Ssid returns the ssid value in the mutation.
func (m *SurveyWiFiScanMutation) Ssid() (r string, exists bool) {
	v := m.ssid
	if v == nil {
		return
	}
	return *v, true
}

// ClearSsid clears the value of ssid.
func (m *SurveyWiFiScanMutation) ClearSsid() {
	m.ssid = nil
	m.clearedFields[surveywifiscan.FieldSsid] = true
}

// SsidCleared returns if the field ssid was cleared in this mutation.
func (m *SurveyWiFiScanMutation) SsidCleared() bool {
	return m.clearedFields[surveywifiscan.FieldSsid]
}

// ResetSsid reset all changes of the ssid field.
func (m *SurveyWiFiScanMutation) ResetSsid() {
	m.ssid = nil
	delete(m.clearedFields, surveywifiscan.FieldSsid)
}

// SetBssid sets the bssid field.
func (m *SurveyWiFiScanMutation) SetBssid(s string) {
	m.bssid = &s
}

// Bssid returns the bssid value in the mutation.
func (m *SurveyWiFiScanMutation) Bssid() (r string, exists bool) {
	v := m.bssid
	if v == nil {
		return
	}
	return *v, true
}

// ResetBssid reset all changes of the bssid field.
func (m *SurveyWiFiScanMutation) ResetBssid() {
	m.bssid = nil
}

// SetTimestamp sets the timestamp field.
func (m *SurveyWiFiScanMutation) SetTimestamp(t time.Time) {
	m.timestamp = &t
}

// Timestamp returns the timestamp value in the mutation.
func (m *SurveyWiFiScanMutation) Timestamp() (r time.Time, exists bool) {
	v := m.timestamp
	if v == nil {
		return
	}
	return *v, true
}

// ResetTimestamp reset all changes of the timestamp field.
func (m *SurveyWiFiScanMutation) ResetTimestamp() {
	m.timestamp = nil
}

// SetFrequency sets the frequency field.
func (m *SurveyWiFiScanMutation) SetFrequency(i int) {
	m.frequency = &i
	m.addfrequency = nil
}

// Frequency returns the frequency value in the mutation.
func (m *SurveyWiFiScanMutation) Frequency() (r int, exists bool) {
	v := m.frequency
	if v == nil {
		return
	}
	return *v, true
}

// AddFrequency adds i to frequency.
func (m *SurveyWiFiScanMutation) AddFrequency(i int) {
	if m.addfrequency != nil {
		*m.addfrequency += i
	} else {
		m.addfrequency = &i
	}
}

// AddedFrequency returns the value that was added to the frequency field in this mutation.
func (m *SurveyWiFiScanMutation) AddedFrequency() (r int, exists bool) {
	v := m.addfrequency
	if v == nil {
		return
	}
	return *v, true
}

// ResetFrequency reset all changes of the frequency field.
func (m *SurveyWiFiScanMutation) ResetFrequency() {
	m.frequency = nil
	m.addfrequency = nil
}

// SetChannel sets the channel field.
func (m *SurveyWiFiScanMutation) SetChannel(i int) {
	m.channel = &i
	m.addchannel = nil
}

// Channel returns the channel value in the mutation.
func (m *SurveyWiFiScanMutation) Channel() (r int, exists bool) {
	v := m.channel
	if v == nil {
		return
	}
	return *v, true
}

// AddChannel adds i to channel.
func (m *SurveyWiFiScanMutation) AddChannel(i int) {
	if m.addchannel != nil {
		*m.addchannel += i
	} else {
		m.addchannel = &i
	}
}

// AddedChannel returns the value that was added to the channel field in this mutation.
func (m *SurveyWiFiScanMutation) AddedChannel() (r int, exists bool) {
	v := m.addchannel
	if v == nil {
		return
	}
	return *v, true
}

// ResetChannel reset all changes of the channel field.
func (m *SurveyWiFiScanMutation) ResetChannel() {
	m.channel = nil
	m.addchannel = nil
}

// SetBand sets the band field.
func (m *SurveyWiFiScanMutation) SetBand(s string) {
	m.band = &s
}

// Band returns the band value in the mutation.
func (m *SurveyWiFiScanMutation) Band() (r string, exists bool) {
	v := m.band
	if v == nil {
		return
	}
	return *v, true
}

// ClearBand clears the value of band.
func (m *SurveyWiFiScanMutation) ClearBand() {
	m.band = nil
	m.clearedFields[surveywifiscan.FieldBand] = true
}

// BandCleared returns if the field band was cleared in this mutation.
func (m *SurveyWiFiScanMutation) BandCleared() bool {
	return m.clearedFields[surveywifiscan.FieldBand]
}

// ResetBand reset all changes of the band field.
func (m *SurveyWiFiScanMutation) ResetBand() {
	m.band = nil
	delete(m.clearedFields, surveywifiscan.FieldBand)
}

// SetChannelWidth sets the channel_width field.
func (m *SurveyWiFiScanMutation) SetChannelWidth(i int) {
	m.channel_width = &i
	m.addchannel_width = nil
}

// ChannelWidth returns the channel_width value in the mutation.
func (m *SurveyWiFiScanMutation) ChannelWidth() (r int, exists bool) {
	v := m.channel_width
	if v == nil {
		return
	}
	return *v, true
}

// AddChannelWidth adds i to channel_width.
func (m *SurveyWiFiScanMutation) AddChannelWidth(i int) {
	if m.addchannel_width != nil {
		*m.addchannel_width += i
	} else {
		m.addchannel_width = &i
	}
}

// AddedChannelWidth returns the value that was added to the channel_width field in this mutation.
func (m *SurveyWiFiScanMutation) AddedChannelWidth() (r int, exists bool) {
	v := m.addchannel_width
	if v == nil {
		return
	}
	return *v, true
}

// ClearChannelWidth clears the value of channel_width.
func (m *SurveyWiFiScanMutation) ClearChannelWidth() {
	m.channel_width = nil
	m.addchannel_width = nil
	m.clearedFields[surveywifiscan.FieldChannelWidth] = true
}

// ChannelWidthCleared returns if the field channel_width was cleared in this mutation.
func (m *SurveyWiFiScanMutation) ChannelWidthCleared() bool {
	return m.clearedFields[surveywifiscan.FieldChannelWidth]
}

// ResetChannelWidth reset all changes of the channel_width field.
func (m *SurveyWiFiScanMutation) ResetChannelWidth() {
	m.channel_width = nil
	m.addchannel_width = nil
	delete(m.clearedFields, surveywifiscan.FieldChannelWidth)
}

// SetCapabilities sets the capabilities field.
func (m *SurveyWiFiScanMutation) SetCapabilities(s string) {
	m.capabilities = &s
}

// Capabilities returns the capabilities value in the mutation.
func (m *SurveyWiFiScanMutation) Capabilities() (r string, exists bool) {
	v := m.capabilities
	if v == nil {
		return
	}
	return *v, true
}

// ClearCapabilities clears the value of capabilities.
func (m *SurveyWiFiScanMutation) ClearCapabilities() {
	m.capabilities = nil
	m.clearedFields[surveywifiscan.FieldCapabilities] = true
}

// CapabilitiesCleared returns if the field capabilities was cleared in this mutation.
func (m *SurveyWiFiScanMutation) CapabilitiesCleared() bool {
	return m.clearedFields[surveywifiscan.FieldCapabilities]
}

// ResetCapabilities reset all changes of the capabilities field.
func (m *SurveyWiFiScanMutation) ResetCapabilities() {
	m.capabilities = nil
	delete(m.clearedFields, surveywifiscan.FieldCapabilities)
}

// SetStrength sets the strength field.
func (m *SurveyWiFiScanMutation) SetStrength(i int) {
	m.strength = &i
	m.addstrength = nil
}

// Strength returns the strength value in the mutation.
func (m *SurveyWiFiScanMutation) Strength() (r int, exists bool) {
	v := m.strength
	if v == nil {
		return
	}
	return *v, true
}

// AddStrength adds i to strength.
func (m *SurveyWiFiScanMutation) AddStrength(i int) {
	if m.addstrength != nil {
		*m.addstrength += i
	} else {
		m.addstrength = &i
	}
}

// AddedStrength returns the value that was added to the strength field in this mutation.
func (m *SurveyWiFiScanMutation) AddedStrength() (r int, exists bool) {
	v := m.addstrength
	if v == nil {
		return
	}
	return *v, true
}

// ResetStrength reset all changes of the strength field.
func (m *SurveyWiFiScanMutation) ResetStrength() {
	m.strength = nil
	m.addstrength = nil
}

// SetLatitude sets the latitude field.
func (m *SurveyWiFiScanMutation) SetLatitude(f float64) {
	m.latitude = &f
	m.addlatitude = nil
}

// Latitude returns the latitude value in the mutation.
func (m *SurveyWiFiScanMutation) Latitude() (r float64, exists bool) {
	v := m.latitude
	if v == nil {
		return
	}
	return *v, true
}

// AddLatitude adds f to latitude.
func (m *SurveyWiFiScanMutation) AddLatitude(f float64) {
	if m.addlatitude != nil {
		*m.addlatitude += f
	} else {
		m.addlatitude = &f
	}
}

// AddedLatitude returns the value that was added to the latitude field in this mutation.
func (m *SurveyWiFiScanMutation) AddedLatitude() (r float64, exists bool) {
	v := m.addlatitude
	if v == nil {
		return
	}
	return *v, true
}

// ClearLatitude clears the value of latitude.
func (m *SurveyWiFiScanMutation) ClearLatitude() {
	m.latitude = nil
	m.addlatitude = nil
	m.clearedFields[surveywifiscan.FieldLatitude] = true
}

// LatitudeCleared returns if the field latitude was cleared in this mutation.
func (m *SurveyWiFiScanMutation) LatitudeCleared() bool {
	return m.clearedFields[surveywifiscan.FieldLatitude]
}

// ResetLatitude reset all changes of the latitude field.
func (m *SurveyWiFiScanMutation) ResetLatitude() {
	m.latitude = nil
	m.addlatitude = nil
	delete(m.clearedFields, surveywifiscan.FieldLatitude)
}

// SetLongitude sets the longitude field.
func (m *SurveyWiFiScanMutation) SetLongitude(f float64) {
	m.longitude = &f
	m.addlongitude = nil
}

// Longitude returns the longitude value in the mutation.
func (m *SurveyWiFiScanMutation) Longitude() (r float64, exists bool) {
	v := m.longitude
	if v == nil {
		return
	}
	return *v, true
}

// AddLongitude adds f to longitude.
func (m *SurveyWiFiScanMutation) AddLongitude(f float64) {
	if m.addlongitude != nil {
		*m.addlongitude += f
	} else {
		m.addlongitude = &f
	}
}

// AddedLongitude returns the value that was added to the longitude field in this mutation.
func (m *SurveyWiFiScanMutation) AddedLongitude() (r float64, exists bool) {
	v := m.addlongitude
	if v == nil {
		return
	}
	return *v, true
}

// ClearLongitude clears the value of longitude.
func (m *SurveyWiFiScanMutation) ClearLongitude() {
	m.longitude = nil
	m.addlongitude = nil
	m.clearedFields[surveywifiscan.FieldLongitude] = true
}

// LongitudeCleared returns if the field longitude was cleared in this mutation.
func (m *SurveyWiFiScanMutation) LongitudeCleared() bool {
	return m.clearedFields[surveywifiscan.FieldLongitude]
}

// ResetLongitude reset all changes of the longitude field.
func (m *SurveyWiFiScanMutation) ResetLongitude() {
	m.longitude = nil
	m.addlongitude = nil
	delete(m.clearedFields, surveywifiscan.FieldLongitude)
}

// SetSurveyQuestionID sets the survey_question edge to SurveyQuestion by id.
func (m *SurveyWiFiScanMutation) SetSurveyQuestionID(id int) {
	m.survey_question = &id
}

// ClearSurveyQuestion clears the survey_question edge to SurveyQuestion.
func (m *SurveyWiFiScanMutation) ClearSurveyQuestion() {
	m.clearedsurvey_question = true
}

// SurveyQuestionCleared returns if the edge survey_question was cleared.
func (m *SurveyWiFiScanMutation) SurveyQuestionCleared() bool {
	return m.clearedsurvey_question
}

// SurveyQuestionID returns the survey_question id in the mutation.
func (m *SurveyWiFiScanMutation) SurveyQuestionID() (id int, exists bool) {
	if m.survey_question != nil {
		return *m.survey_question, true
	}
	return
}

// SurveyQuestionIDs returns the survey_question ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// SurveyQuestionID instead. It exists only for internal usage by the builders.
func (m *SurveyWiFiScanMutation) SurveyQuestionIDs() (ids []int) {
	if id := m.survey_question; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetSurveyQuestion reset all changes of the survey_question edge.
func (m *SurveyWiFiScanMutation) ResetSurveyQuestion() {
	m.survey_question = nil
	m.clearedsurvey_question = false
}

// SetLocationID sets the location edge to Location by id.
func (m *SurveyWiFiScanMutation) SetLocationID(id int) {
	m.location = &id
}

// ClearLocation clears the location edge to Location.
func (m *SurveyWiFiScanMutation) ClearLocation() {
	m.clearedlocation = true
}

// LocationCleared returns if the edge location was cleared.
func (m *SurveyWiFiScanMutation) LocationCleared() bool {
	return m.clearedlocation
}

// LocationID returns the location id in the mutation.
func (m *SurveyWiFiScanMutation) LocationID() (id int, exists bool) {
	if m.location != nil {
		return *m.location, true
	}
	return
}

// LocationIDs returns the location ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LocationID instead. It exists only for internal usage by the builders.
func (m *SurveyWiFiScanMutation) LocationIDs() (ids []int) {
	if id := m.location; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLocation reset all changes of the location edge.
func (m *SurveyWiFiScanMutation) ResetLocation() {
	m.location = nil
	m.clearedlocation = false
}

// Op returns the operation name.
func (m *SurveyWiFiScanMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (SurveyWiFiScan).
func (m *SurveyWiFiScanMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *SurveyWiFiScanMutation) Fields() []string {
	fields := make([]string, 0, 13)
	if m.create_time != nil {
		fields = append(fields, surveywifiscan.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, surveywifiscan.FieldUpdateTime)
	}
	if m.ssid != nil {
		fields = append(fields, surveywifiscan.FieldSsid)
	}
	if m.bssid != nil {
		fields = append(fields, surveywifiscan.FieldBssid)
	}
	if m.timestamp != nil {
		fields = append(fields, surveywifiscan.FieldTimestamp)
	}
	if m.frequency != nil {
		fields = append(fields, surveywifiscan.FieldFrequency)
	}
	if m.channel != nil {
		fields = append(fields, surveywifiscan.FieldChannel)
	}
	if m.band != nil {
		fields = append(fields, surveywifiscan.FieldBand)
	}
	if m.channel_width != nil {
		fields = append(fields, surveywifiscan.FieldChannelWidth)
	}
	if m.capabilities != nil {
		fields = append(fields, surveywifiscan.FieldCapabilities)
	}
	if m.strength != nil {
		fields = append(fields, surveywifiscan.FieldStrength)
	}
	if m.latitude != nil {
		fields = append(fields, surveywifiscan.FieldLatitude)
	}
	if m.longitude != nil {
		fields = append(fields, surveywifiscan.FieldLongitude)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *SurveyWiFiScanMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case surveywifiscan.FieldCreateTime:
		return m.CreateTime()
	case surveywifiscan.FieldUpdateTime:
		return m.UpdateTime()
	case surveywifiscan.FieldSsid:
		return m.Ssid()
	case surveywifiscan.FieldBssid:
		return m.Bssid()
	case surveywifiscan.FieldTimestamp:
		return m.Timestamp()
	case surveywifiscan.FieldFrequency:
		return m.Frequency()
	case surveywifiscan.FieldChannel:
		return m.Channel()
	case surveywifiscan.FieldBand:
		return m.Band()
	case surveywifiscan.FieldChannelWidth:
		return m.ChannelWidth()
	case surveywifiscan.FieldCapabilities:
		return m.Capabilities()
	case surveywifiscan.FieldStrength:
		return m.Strength()
	case surveywifiscan.FieldLatitude:
		return m.Latitude()
	case surveywifiscan.FieldLongitude:
		return m.Longitude()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyWiFiScanMutation) SetField(name string, value ent.Value) error {
	switch name {
	case surveywifiscan.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case surveywifiscan.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case surveywifiscan.FieldSsid:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSsid(v)
		return nil
	case surveywifiscan.FieldBssid:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBssid(v)
		return nil
	case surveywifiscan.FieldTimestamp:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTimestamp(v)
		return nil
	case surveywifiscan.FieldFrequency:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFrequency(v)
		return nil
	case surveywifiscan.FieldChannel:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetChannel(v)
		return nil
	case surveywifiscan.FieldBand:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetBand(v)
		return nil
	case surveywifiscan.FieldChannelWidth:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetChannelWidth(v)
		return nil
	case surveywifiscan.FieldCapabilities:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCapabilities(v)
		return nil
	case surveywifiscan.FieldStrength:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStrength(v)
		return nil
	case surveywifiscan.FieldLatitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLatitude(v)
		return nil
	case surveywifiscan.FieldLongitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLongitude(v)
		return nil
	}
	return fmt.Errorf("unknown SurveyWiFiScan field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *SurveyWiFiScanMutation) AddedFields() []string {
	var fields []string
	if m.addfrequency != nil {
		fields = append(fields, surveywifiscan.FieldFrequency)
	}
	if m.addchannel != nil {
		fields = append(fields, surveywifiscan.FieldChannel)
	}
	if m.addchannel_width != nil {
		fields = append(fields, surveywifiscan.FieldChannelWidth)
	}
	if m.addstrength != nil {
		fields = append(fields, surveywifiscan.FieldStrength)
	}
	if m.addlatitude != nil {
		fields = append(fields, surveywifiscan.FieldLatitude)
	}
	if m.addlongitude != nil {
		fields = append(fields, surveywifiscan.FieldLongitude)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *SurveyWiFiScanMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case surveywifiscan.FieldFrequency:
		return m.AddedFrequency()
	case surveywifiscan.FieldChannel:
		return m.AddedChannel()
	case surveywifiscan.FieldChannelWidth:
		return m.AddedChannelWidth()
	case surveywifiscan.FieldStrength:
		return m.AddedStrength()
	case surveywifiscan.FieldLatitude:
		return m.AddedLatitude()
	case surveywifiscan.FieldLongitude:
		return m.AddedLongitude()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *SurveyWiFiScanMutation) AddField(name string, value ent.Value) error {
	switch name {
	case surveywifiscan.FieldFrequency:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddFrequency(v)
		return nil
	case surveywifiscan.FieldChannel:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddChannel(v)
		return nil
	case surveywifiscan.FieldChannelWidth:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddChannelWidth(v)
		return nil
	case surveywifiscan.FieldStrength:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddStrength(v)
		return nil
	case surveywifiscan.FieldLatitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLatitude(v)
		return nil
	case surveywifiscan.FieldLongitude:
		v, ok := value.(float64)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddLongitude(v)
		return nil
	}
	return fmt.Errorf("unknown SurveyWiFiScan numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *SurveyWiFiScanMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[surveywifiscan.FieldSsid] {
		fields = append(fields, surveywifiscan.FieldSsid)
	}
	if m.clearedFields[surveywifiscan.FieldBand] {
		fields = append(fields, surveywifiscan.FieldBand)
	}
	if m.clearedFields[surveywifiscan.FieldChannelWidth] {
		fields = append(fields, surveywifiscan.FieldChannelWidth)
	}
	if m.clearedFields[surveywifiscan.FieldCapabilities] {
		fields = append(fields, surveywifiscan.FieldCapabilities)
	}
	if m.clearedFields[surveywifiscan.FieldLatitude] {
		fields = append(fields, surveywifiscan.FieldLatitude)
	}
	if m.clearedFields[surveywifiscan.FieldLongitude] {
		fields = append(fields, surveywifiscan.FieldLongitude)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *SurveyWiFiScanMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *SurveyWiFiScanMutation) ClearField(name string) error {
	switch name {
	case surveywifiscan.FieldSsid:
		m.ClearSsid()
		return nil
	case surveywifiscan.FieldBand:
		m.ClearBand()
		return nil
	case surveywifiscan.FieldChannelWidth:
		m.ClearChannelWidth()
		return nil
	case surveywifiscan.FieldCapabilities:
		m.ClearCapabilities()
		return nil
	case surveywifiscan.FieldLatitude:
		m.ClearLatitude()
		return nil
	case surveywifiscan.FieldLongitude:
		m.ClearLongitude()
		return nil
	}
	return fmt.Errorf("unknown SurveyWiFiScan nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *SurveyWiFiScanMutation) ResetField(name string) error {
	switch name {
	case surveywifiscan.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case surveywifiscan.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case surveywifiscan.FieldSsid:
		m.ResetSsid()
		return nil
	case surveywifiscan.FieldBssid:
		m.ResetBssid()
		return nil
	case surveywifiscan.FieldTimestamp:
		m.ResetTimestamp()
		return nil
	case surveywifiscan.FieldFrequency:
		m.ResetFrequency()
		return nil
	case surveywifiscan.FieldChannel:
		m.ResetChannel()
		return nil
	case surveywifiscan.FieldBand:
		m.ResetBand()
		return nil
	case surveywifiscan.FieldChannelWidth:
		m.ResetChannelWidth()
		return nil
	case surveywifiscan.FieldCapabilities:
		m.ResetCapabilities()
		return nil
	case surveywifiscan.FieldStrength:
		m.ResetStrength()
		return nil
	case surveywifiscan.FieldLatitude:
		m.ResetLatitude()
		return nil
	case surveywifiscan.FieldLongitude:
		m.ResetLongitude()
		return nil
	}
	return fmt.Errorf("unknown SurveyWiFiScan field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *SurveyWiFiScanMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.survey_question != nil {
		edges = append(edges, surveywifiscan.EdgeSurveyQuestion)
	}
	if m.location != nil {
		edges = append(edges, surveywifiscan.EdgeLocation)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *SurveyWiFiScanMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case surveywifiscan.EdgeSurveyQuestion:
		if id := m.survey_question; id != nil {
			return []ent.Value{*id}
		}
	case surveywifiscan.EdgeLocation:
		if id := m.location; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *SurveyWiFiScanMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *SurveyWiFiScanMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *SurveyWiFiScanMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedsurvey_question {
		edges = append(edges, surveywifiscan.EdgeSurveyQuestion)
	}
	if m.clearedlocation {
		edges = append(edges, surveywifiscan.EdgeLocation)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *SurveyWiFiScanMutation) EdgeCleared(name string) bool {
	switch name {
	case surveywifiscan.EdgeSurveyQuestion:
		return m.clearedsurvey_question
	case surveywifiscan.EdgeLocation:
		return m.clearedlocation
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *SurveyWiFiScanMutation) ClearEdge(name string) error {
	switch name {
	case surveywifiscan.EdgeSurveyQuestion:
		m.ClearSurveyQuestion()
		return nil
	case surveywifiscan.EdgeLocation:
		m.ClearLocation()
		return nil
	}
	return fmt.Errorf("unknown SurveyWiFiScan unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *SurveyWiFiScanMutation) ResetEdge(name string) error {
	switch name {
	case surveywifiscan.EdgeSurveyQuestion:
		m.ResetSurveyQuestion()
		return nil
	case surveywifiscan.EdgeLocation:
		m.ResetLocation()
		return nil
	}
	return fmt.Errorf("unknown SurveyWiFiScan edge %s", name)
}

// TechnicianMutation represents an operation that mutate the Technicians
// nodes in the graph.
type TechnicianMutation struct {
	config
	op                 Op
	typ                string
	id                 *int
	create_time        *time.Time
	update_time        *time.Time
	name               *string
	email              *string
	clearedFields      map[string]bool
	work_orders        map[int]struct{}
	removedwork_orders map[int]struct{}
}

var _ ent.Mutation = (*TechnicianMutation)(nil)

// newTechnicianMutation creates new mutation for $n.Name.
func newTechnicianMutation(c config, op Op) *TechnicianMutation {
	return &TechnicianMutation{
		config:        c,
		op:            op,
		typ:           TypeTechnician,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m TechnicianMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m TechnicianMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *TechnicianMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *TechnicianMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *TechnicianMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *TechnicianMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *TechnicianMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *TechnicianMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *TechnicianMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *TechnicianMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *TechnicianMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *TechnicianMutation) ResetName() {
	m.name = nil
}

// SetEmail sets the email field.
func (m *TechnicianMutation) SetEmail(s string) {
	m.email = &s
}

// Email returns the email value in the mutation.
func (m *TechnicianMutation) Email() (r string, exists bool) {
	v := m.email
	if v == nil {
		return
	}
	return *v, true
}

// ResetEmail reset all changes of the email field.
func (m *TechnicianMutation) ResetEmail() {
	m.email = nil
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (m *TechnicianMutation) AddWorkOrderIDs(ids ...int) {
	if m.work_orders == nil {
		m.work_orders = make(map[int]struct{})
	}
	for i := range ids {
		m.work_orders[ids[i]] = struct{}{}
	}
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (m *TechnicianMutation) RemoveWorkOrderIDs(ids ...int) {
	if m.removedwork_orders == nil {
		m.removedwork_orders = make(map[int]struct{})
	}
	for i := range ids {
		m.removedwork_orders[ids[i]] = struct{}{}
	}
}

// RemovedWorkOrders returns the removed ids of work_orders.
func (m *TechnicianMutation) RemovedWorkOrdersIDs() (ids []int) {
	for id := range m.removedwork_orders {
		ids = append(ids, id)
	}
	return
}

// WorkOrdersIDs returns the work_orders ids in the mutation.
func (m *TechnicianMutation) WorkOrdersIDs() (ids []int) {
	for id := range m.work_orders {
		ids = append(ids, id)
	}
	return
}

// ResetWorkOrders reset all changes of the work_orders edge.
func (m *TechnicianMutation) ResetWorkOrders() {
	m.work_orders = nil
	m.removedwork_orders = nil
}

// Op returns the operation name.
func (m *TechnicianMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Technician).
func (m *TechnicianMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *TechnicianMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.create_time != nil {
		fields = append(fields, technician.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, technician.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, technician.FieldName)
	}
	if m.email != nil {
		fields = append(fields, technician.FieldEmail)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *TechnicianMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case technician.FieldCreateTime:
		return m.CreateTime()
	case technician.FieldUpdateTime:
		return m.UpdateTime()
	case technician.FieldName:
		return m.Name()
	case technician.FieldEmail:
		return m.Email()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *TechnicianMutation) SetField(name string, value ent.Value) error {
	switch name {
	case technician.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case technician.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case technician.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case technician.FieldEmail:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEmail(v)
		return nil
	}
	return fmt.Errorf("unknown Technician field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *TechnicianMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *TechnicianMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *TechnicianMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Technician numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *TechnicianMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *TechnicianMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *TechnicianMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Technician nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *TechnicianMutation) ResetField(name string) error {
	switch name {
	case technician.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case technician.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case technician.FieldName:
		m.ResetName()
		return nil
	case technician.FieldEmail:
		m.ResetEmail()
		return nil
	}
	return fmt.Errorf("unknown Technician field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *TechnicianMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.work_orders != nil {
		edges = append(edges, technician.EdgeWorkOrders)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *TechnicianMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case technician.EdgeWorkOrders:
		ids := make([]ent.Value, 0, len(m.work_orders))
		for id := range m.work_orders {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *TechnicianMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	if m.removedwork_orders != nil {
		edges = append(edges, technician.EdgeWorkOrders)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *TechnicianMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case technician.EdgeWorkOrders:
		ids := make([]ent.Value, 0, len(m.removedwork_orders))
		for id := range m.removedwork_orders {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *TechnicianMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *TechnicianMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *TechnicianMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown Technician unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *TechnicianMutation) ResetEdge(name string) error {
	switch name {
	case technician.EdgeWorkOrders:
		m.ResetWorkOrders()
		return nil
	}
	return fmt.Errorf("unknown Technician edge %s", name)
}

// UserMutation represents an operation that mutate the Users
// nodes in the graph.
type UserMutation struct {
	config
	op                   Op
	typ                  string
	id                   *int
	create_time          *time.Time
	update_time          *time.Time
	auth_id              *string
	first_name           *string
	last_name            *string
	email                *string
	status               *user.Status
	role                 *user.Role
	clearedFields        map[string]bool
	profile_photo        *int
	clearedprofile_photo bool
}

var _ ent.Mutation = (*UserMutation)(nil)

// newUserMutation creates new mutation for $n.Name.
func newUserMutation(c config, op Op) *UserMutation {
	return &UserMutation{
		config:        c,
		op:            op,
		typ:           TypeUser,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m UserMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m UserMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *UserMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *UserMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *UserMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *UserMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *UserMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *UserMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *UserMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetAuthID sets the auth_id field.
func (m *UserMutation) SetAuthID(s string) {
	m.auth_id = &s
}

// AuthID returns the auth_id value in the mutation.
func (m *UserMutation) AuthID() (r string, exists bool) {
	v := m.auth_id
	if v == nil {
		return
	}
	return *v, true
}

// ResetAuthID reset all changes of the auth_id field.
func (m *UserMutation) ResetAuthID() {
	m.auth_id = nil
}

// SetFirstName sets the first_name field.
func (m *UserMutation) SetFirstName(s string) {
	m.first_name = &s
}

// FirstName returns the first_name value in the mutation.
func (m *UserMutation) FirstName() (r string, exists bool) {
	v := m.first_name
	if v == nil {
		return
	}
	return *v, true
}

// ClearFirstName clears the value of first_name.
func (m *UserMutation) ClearFirstName() {
	m.first_name = nil
	m.clearedFields[user.FieldFirstName] = true
}

// FirstNameCleared returns if the field first_name was cleared in this mutation.
func (m *UserMutation) FirstNameCleared() bool {
	return m.clearedFields[user.FieldFirstName]
}

// ResetFirstName reset all changes of the first_name field.
func (m *UserMutation) ResetFirstName() {
	m.first_name = nil
	delete(m.clearedFields, user.FieldFirstName)
}

// SetLastName sets the last_name field.
func (m *UserMutation) SetLastName(s string) {
	m.last_name = &s
}

// LastName returns the last_name value in the mutation.
func (m *UserMutation) LastName() (r string, exists bool) {
	v := m.last_name
	if v == nil {
		return
	}
	return *v, true
}

// ClearLastName clears the value of last_name.
func (m *UserMutation) ClearLastName() {
	m.last_name = nil
	m.clearedFields[user.FieldLastName] = true
}

// LastNameCleared returns if the field last_name was cleared in this mutation.
func (m *UserMutation) LastNameCleared() bool {
	return m.clearedFields[user.FieldLastName]
}

// ResetLastName reset all changes of the last_name field.
func (m *UserMutation) ResetLastName() {
	m.last_name = nil
	delete(m.clearedFields, user.FieldLastName)
}

// SetEmail sets the email field.
func (m *UserMutation) SetEmail(s string) {
	m.email = &s
}

// Email returns the email value in the mutation.
func (m *UserMutation) Email() (r string, exists bool) {
	v := m.email
	if v == nil {
		return
	}
	return *v, true
}

// ClearEmail clears the value of email.
func (m *UserMutation) ClearEmail() {
	m.email = nil
	m.clearedFields[user.FieldEmail] = true
}

// EmailCleared returns if the field email was cleared in this mutation.
func (m *UserMutation) EmailCleared() bool {
	return m.clearedFields[user.FieldEmail]
}

// ResetEmail reset all changes of the email field.
func (m *UserMutation) ResetEmail() {
	m.email = nil
	delete(m.clearedFields, user.FieldEmail)
}

// SetStatus sets the status field.
func (m *UserMutation) SetStatus(u user.Status) {
	m.status = &u
}

// Status returns the status value in the mutation.
func (m *UserMutation) Status() (r user.Status, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// ResetStatus reset all changes of the status field.
func (m *UserMutation) ResetStatus() {
	m.status = nil
}

// SetRole sets the role field.
func (m *UserMutation) SetRole(u user.Role) {
	m.role = &u
}

// Role returns the role value in the mutation.
func (m *UserMutation) Role() (r user.Role, exists bool) {
	v := m.role
	if v == nil {
		return
	}
	return *v, true
}

// ResetRole reset all changes of the role field.
func (m *UserMutation) ResetRole() {
	m.role = nil
}

// SetProfilePhotoID sets the profile_photo edge to File by id.
func (m *UserMutation) SetProfilePhotoID(id int) {
	m.profile_photo = &id
}

// ClearProfilePhoto clears the profile_photo edge to File.
func (m *UserMutation) ClearProfilePhoto() {
	m.clearedprofile_photo = true
}

// ProfilePhotoCleared returns if the edge profile_photo was cleared.
func (m *UserMutation) ProfilePhotoCleared() bool {
	return m.clearedprofile_photo
}

// ProfilePhotoID returns the profile_photo id in the mutation.
func (m *UserMutation) ProfilePhotoID() (id int, exists bool) {
	if m.profile_photo != nil {
		return *m.profile_photo, true
	}
	return
}

// ProfilePhotoIDs returns the profile_photo ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ProfilePhotoID instead. It exists only for internal usage by the builders.
func (m *UserMutation) ProfilePhotoIDs() (ids []int) {
	if id := m.profile_photo; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetProfilePhoto reset all changes of the profile_photo edge.
func (m *UserMutation) ResetProfilePhoto() {
	m.profile_photo = nil
	m.clearedprofile_photo = false
}

// Op returns the operation name.
func (m *UserMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (User).
func (m *UserMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *UserMutation) Fields() []string {
	fields := make([]string, 0, 8)
	if m.create_time != nil {
		fields = append(fields, user.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, user.FieldUpdateTime)
	}
	if m.auth_id != nil {
		fields = append(fields, user.FieldAuthID)
	}
	if m.first_name != nil {
		fields = append(fields, user.FieldFirstName)
	}
	if m.last_name != nil {
		fields = append(fields, user.FieldLastName)
	}
	if m.email != nil {
		fields = append(fields, user.FieldEmail)
	}
	if m.status != nil {
		fields = append(fields, user.FieldStatus)
	}
	if m.role != nil {
		fields = append(fields, user.FieldRole)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *UserMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case user.FieldCreateTime:
		return m.CreateTime()
	case user.FieldUpdateTime:
		return m.UpdateTime()
	case user.FieldAuthID:
		return m.AuthID()
	case user.FieldFirstName:
		return m.FirstName()
	case user.FieldLastName:
		return m.LastName()
	case user.FieldEmail:
		return m.Email()
	case user.FieldStatus:
		return m.Status()
	case user.FieldRole:
		return m.Role()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *UserMutation) SetField(name string, value ent.Value) error {
	switch name {
	case user.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case user.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case user.FieldAuthID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetAuthID(v)
		return nil
	case user.FieldFirstName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFirstName(v)
		return nil
	case user.FieldLastName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetLastName(v)
		return nil
	case user.FieldEmail:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetEmail(v)
		return nil
	case user.FieldStatus:
		v, ok := value.(user.Status)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case user.FieldRole:
		v, ok := value.(user.Role)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRole(v)
		return nil
	}
	return fmt.Errorf("unknown User field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *UserMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *UserMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *UserMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown User numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *UserMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[user.FieldFirstName] {
		fields = append(fields, user.FieldFirstName)
	}
	if m.clearedFields[user.FieldLastName] {
		fields = append(fields, user.FieldLastName)
	}
	if m.clearedFields[user.FieldEmail] {
		fields = append(fields, user.FieldEmail)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *UserMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *UserMutation) ClearField(name string) error {
	switch name {
	case user.FieldFirstName:
		m.ClearFirstName()
		return nil
	case user.FieldLastName:
		m.ClearLastName()
		return nil
	case user.FieldEmail:
		m.ClearEmail()
		return nil
	}
	return fmt.Errorf("unknown User nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *UserMutation) ResetField(name string) error {
	switch name {
	case user.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case user.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case user.FieldAuthID:
		m.ResetAuthID()
		return nil
	case user.FieldFirstName:
		m.ResetFirstName()
		return nil
	case user.FieldLastName:
		m.ResetLastName()
		return nil
	case user.FieldEmail:
		m.ResetEmail()
		return nil
	case user.FieldStatus:
		m.ResetStatus()
		return nil
	case user.FieldRole:
		m.ResetRole()
		return nil
	}
	return fmt.Errorf("unknown User field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *UserMutation) AddedEdges() []string {
	edges := make([]string, 0, 1)
	if m.profile_photo != nil {
		edges = append(edges, user.EdgeProfilePhoto)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *UserMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case user.EdgeProfilePhoto:
		if id := m.profile_photo; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *UserMutation) RemovedEdges() []string {
	edges := make([]string, 0, 1)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *UserMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *UserMutation) ClearedEdges() []string {
	edges := make([]string, 0, 1)
	if m.clearedprofile_photo {
		edges = append(edges, user.EdgeProfilePhoto)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *UserMutation) EdgeCleared(name string) bool {
	switch name {
	case user.EdgeProfilePhoto:
		return m.clearedprofile_photo
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *UserMutation) ClearEdge(name string) error {
	switch name {
	case user.EdgeProfilePhoto:
		m.ClearProfilePhoto()
		return nil
	}
	return fmt.Errorf("unknown User unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *UserMutation) ResetEdge(name string) error {
	switch name {
	case user.EdgeProfilePhoto:
		m.ResetProfilePhoto()
		return nil
	}
	return fmt.Errorf("unknown User edge %s", name)
}

// WorkOrderMutation represents an operation that mutate the WorkOrders
// nodes in the graph.
type WorkOrderMutation struct {
	config
	op                           Op
	typ                          string
	id                           *int
	create_time                  *time.Time
	update_time                  *time.Time
	name                         *string
	status                       *string
	priority                     *string
	description                  *string
	install_date                 *time.Time
	creation_date                *time.Time
	index                        *int
	addindex                     *int
	close_date                   *time.Time
	clearedFields                map[string]bool
	_type                        *int
	cleared_type                 bool
	equipment                    map[int]struct{}
	removedequipment             map[int]struct{}
	links                        map[int]struct{}
	removedlinks                 map[int]struct{}
	files                        map[int]struct{}
	removedfiles                 map[int]struct{}
	hyperlinks                   map[int]struct{}
	removedhyperlinks            map[int]struct{}
	location                     *int
	clearedlocation              bool
	comments                     map[int]struct{}
	removedcomments              map[int]struct{}
	properties                   map[int]struct{}
	removedproperties            map[int]struct{}
	check_list_categories        map[int]struct{}
	removedcheck_list_categories map[int]struct{}
	check_list_items             map[int]struct{}
	removedcheck_list_items      map[int]struct{}
	technician                   *int
	clearedtechnician            bool
	project                      *int
	clearedproject               bool
	owner                        *int
	clearedowner                 bool
	assignee                     *int
	clearedassignee              bool
}

var _ ent.Mutation = (*WorkOrderMutation)(nil)

// newWorkOrderMutation creates new mutation for $n.Name.
func newWorkOrderMutation(c config, op Op) *WorkOrderMutation {
	return &WorkOrderMutation{
		config:        c,
		op:            op,
		typ:           TypeWorkOrder,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m WorkOrderMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m WorkOrderMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *WorkOrderMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *WorkOrderMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *WorkOrderMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *WorkOrderMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *WorkOrderMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *WorkOrderMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *WorkOrderMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *WorkOrderMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *WorkOrderMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *WorkOrderMutation) ResetName() {
	m.name = nil
}

// SetStatus sets the status field.
func (m *WorkOrderMutation) SetStatus(s string) {
	m.status = &s
}

// Status returns the status value in the mutation.
func (m *WorkOrderMutation) Status() (r string, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// ResetStatus reset all changes of the status field.
func (m *WorkOrderMutation) ResetStatus() {
	m.status = nil
}

// SetPriority sets the priority field.
func (m *WorkOrderMutation) SetPriority(s string) {
	m.priority = &s
}

// Priority returns the priority value in the mutation.
func (m *WorkOrderMutation) Priority() (r string, exists bool) {
	v := m.priority
	if v == nil {
		return
	}
	return *v, true
}

// ResetPriority reset all changes of the priority field.
func (m *WorkOrderMutation) ResetPriority() {
	m.priority = nil
}

// SetDescription sets the description field.
func (m *WorkOrderMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the description value in the mutation.
func (m *WorkOrderMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// ClearDescription clears the value of description.
func (m *WorkOrderMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[workorder.FieldDescription] = true
}

// DescriptionCleared returns if the field description was cleared in this mutation.
func (m *WorkOrderMutation) DescriptionCleared() bool {
	return m.clearedFields[workorder.FieldDescription]
}

// ResetDescription reset all changes of the description field.
func (m *WorkOrderMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, workorder.FieldDescription)
}

// SetInstallDate sets the install_date field.
func (m *WorkOrderMutation) SetInstallDate(t time.Time) {
	m.install_date = &t
}

// InstallDate returns the install_date value in the mutation.
func (m *WorkOrderMutation) InstallDate() (r time.Time, exists bool) {
	v := m.install_date
	if v == nil {
		return
	}
	return *v, true
}

// ClearInstallDate clears the value of install_date.
func (m *WorkOrderMutation) ClearInstallDate() {
	m.install_date = nil
	m.clearedFields[workorder.FieldInstallDate] = true
}

// InstallDateCleared returns if the field install_date was cleared in this mutation.
func (m *WorkOrderMutation) InstallDateCleared() bool {
	return m.clearedFields[workorder.FieldInstallDate]
}

// ResetInstallDate reset all changes of the install_date field.
func (m *WorkOrderMutation) ResetInstallDate() {
	m.install_date = nil
	delete(m.clearedFields, workorder.FieldInstallDate)
}

// SetCreationDate sets the creation_date field.
func (m *WorkOrderMutation) SetCreationDate(t time.Time) {
	m.creation_date = &t
}

// CreationDate returns the creation_date value in the mutation.
func (m *WorkOrderMutation) CreationDate() (r time.Time, exists bool) {
	v := m.creation_date
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreationDate reset all changes of the creation_date field.
func (m *WorkOrderMutation) ResetCreationDate() {
	m.creation_date = nil
}

// SetIndex sets the index field.
func (m *WorkOrderMutation) SetIndex(i int) {
	m.index = &i
	m.addindex = nil
}

// Index returns the index value in the mutation.
func (m *WorkOrderMutation) Index() (r int, exists bool) {
	v := m.index
	if v == nil {
		return
	}
	return *v, true
}

// AddIndex adds i to index.
func (m *WorkOrderMutation) AddIndex(i int) {
	if m.addindex != nil {
		*m.addindex += i
	} else {
		m.addindex = &i
	}
}

// AddedIndex returns the value that was added to the index field in this mutation.
func (m *WorkOrderMutation) AddedIndex() (r int, exists bool) {
	v := m.addindex
	if v == nil {
		return
	}
	return *v, true
}

// ClearIndex clears the value of index.
func (m *WorkOrderMutation) ClearIndex() {
	m.index = nil
	m.addindex = nil
	m.clearedFields[workorder.FieldIndex] = true
}

// IndexCleared returns if the field index was cleared in this mutation.
func (m *WorkOrderMutation) IndexCleared() bool {
	return m.clearedFields[workorder.FieldIndex]
}

// ResetIndex reset all changes of the index field.
func (m *WorkOrderMutation) ResetIndex() {
	m.index = nil
	m.addindex = nil
	delete(m.clearedFields, workorder.FieldIndex)
}

// SetCloseDate sets the close_date field.
func (m *WorkOrderMutation) SetCloseDate(t time.Time) {
	m.close_date = &t
}

// CloseDate returns the close_date value in the mutation.
func (m *WorkOrderMutation) CloseDate() (r time.Time, exists bool) {
	v := m.close_date
	if v == nil {
		return
	}
	return *v, true
}

// ClearCloseDate clears the value of close_date.
func (m *WorkOrderMutation) ClearCloseDate() {
	m.close_date = nil
	m.clearedFields[workorder.FieldCloseDate] = true
}

// CloseDateCleared returns if the field close_date was cleared in this mutation.
func (m *WorkOrderMutation) CloseDateCleared() bool {
	return m.clearedFields[workorder.FieldCloseDate]
}

// ResetCloseDate reset all changes of the close_date field.
func (m *WorkOrderMutation) ResetCloseDate() {
	m.close_date = nil
	delete(m.clearedFields, workorder.FieldCloseDate)
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (m *WorkOrderMutation) SetTypeID(id int) {
	m._type = &id
}

// ClearType clears the type edge to WorkOrderType.
func (m *WorkOrderMutation) ClearType() {
	m.cleared_type = true
}

// TypeCleared returns if the edge type was cleared.
func (m *WorkOrderMutation) TypeCleared() bool {
	return m.cleared_type
}

// TypeID returns the type id in the mutation.
func (m *WorkOrderMutation) TypeID() (id int, exists bool) {
	if m._type != nil {
		return *m._type, true
	}
	return
}

// TypeIDs returns the type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// TypeID instead. It exists only for internal usage by the builders.
func (m *WorkOrderMutation) TypeIDs() (ids []int) {
	if id := m._type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetType reset all changes of the type edge.
func (m *WorkOrderMutation) ResetType() {
	m._type = nil
	m.cleared_type = false
}

// AddEquipmentIDs adds the equipment edge to Equipment by ids.
func (m *WorkOrderMutation) AddEquipmentIDs(ids ...int) {
	if m.equipment == nil {
		m.equipment = make(map[int]struct{})
	}
	for i := range ids {
		m.equipment[ids[i]] = struct{}{}
	}
}

// RemoveEquipmentIDs removes the equipment edge to Equipment by ids.
func (m *WorkOrderMutation) RemoveEquipmentIDs(ids ...int) {
	if m.removedequipment == nil {
		m.removedequipment = make(map[int]struct{})
	}
	for i := range ids {
		m.removedequipment[ids[i]] = struct{}{}
	}
}

// RemovedEquipment returns the removed ids of equipment.
func (m *WorkOrderMutation) RemovedEquipmentIDs() (ids []int) {
	for id := range m.removedequipment {
		ids = append(ids, id)
	}
	return
}

// EquipmentIDs returns the equipment ids in the mutation.
func (m *WorkOrderMutation) EquipmentIDs() (ids []int) {
	for id := range m.equipment {
		ids = append(ids, id)
	}
	return
}

// ResetEquipment reset all changes of the equipment edge.
func (m *WorkOrderMutation) ResetEquipment() {
	m.equipment = nil
	m.removedequipment = nil
}

// AddLinkIDs adds the links edge to Link by ids.
func (m *WorkOrderMutation) AddLinkIDs(ids ...int) {
	if m.links == nil {
		m.links = make(map[int]struct{})
	}
	for i := range ids {
		m.links[ids[i]] = struct{}{}
	}
}

// RemoveLinkIDs removes the links edge to Link by ids.
func (m *WorkOrderMutation) RemoveLinkIDs(ids ...int) {
	if m.removedlinks == nil {
		m.removedlinks = make(map[int]struct{})
	}
	for i := range ids {
		m.removedlinks[ids[i]] = struct{}{}
	}
}

// RemovedLinks returns the removed ids of links.
func (m *WorkOrderMutation) RemovedLinksIDs() (ids []int) {
	for id := range m.removedlinks {
		ids = append(ids, id)
	}
	return
}

// LinksIDs returns the links ids in the mutation.
func (m *WorkOrderMutation) LinksIDs() (ids []int) {
	for id := range m.links {
		ids = append(ids, id)
	}
	return
}

// ResetLinks reset all changes of the links edge.
func (m *WorkOrderMutation) ResetLinks() {
	m.links = nil
	m.removedlinks = nil
}

// AddFileIDs adds the files edge to File by ids.
func (m *WorkOrderMutation) AddFileIDs(ids ...int) {
	if m.files == nil {
		m.files = make(map[int]struct{})
	}
	for i := range ids {
		m.files[ids[i]] = struct{}{}
	}
}

// RemoveFileIDs removes the files edge to File by ids.
func (m *WorkOrderMutation) RemoveFileIDs(ids ...int) {
	if m.removedfiles == nil {
		m.removedfiles = make(map[int]struct{})
	}
	for i := range ids {
		m.removedfiles[ids[i]] = struct{}{}
	}
}

// RemovedFiles returns the removed ids of files.
func (m *WorkOrderMutation) RemovedFilesIDs() (ids []int) {
	for id := range m.removedfiles {
		ids = append(ids, id)
	}
	return
}

// FilesIDs returns the files ids in the mutation.
func (m *WorkOrderMutation) FilesIDs() (ids []int) {
	for id := range m.files {
		ids = append(ids, id)
	}
	return
}

// ResetFiles reset all changes of the files edge.
func (m *WorkOrderMutation) ResetFiles() {
	m.files = nil
	m.removedfiles = nil
}

// AddHyperlinkIDs adds the hyperlinks edge to Hyperlink by ids.
func (m *WorkOrderMutation) AddHyperlinkIDs(ids ...int) {
	if m.hyperlinks == nil {
		m.hyperlinks = make(map[int]struct{})
	}
	for i := range ids {
		m.hyperlinks[ids[i]] = struct{}{}
	}
}

// RemoveHyperlinkIDs removes the hyperlinks edge to Hyperlink by ids.
func (m *WorkOrderMutation) RemoveHyperlinkIDs(ids ...int) {
	if m.removedhyperlinks == nil {
		m.removedhyperlinks = make(map[int]struct{})
	}
	for i := range ids {
		m.removedhyperlinks[ids[i]] = struct{}{}
	}
}

// RemovedHyperlinks returns the removed ids of hyperlinks.
func (m *WorkOrderMutation) RemovedHyperlinksIDs() (ids []int) {
	for id := range m.removedhyperlinks {
		ids = append(ids, id)
	}
	return
}

// HyperlinksIDs returns the hyperlinks ids in the mutation.
func (m *WorkOrderMutation) HyperlinksIDs() (ids []int) {
	for id := range m.hyperlinks {
		ids = append(ids, id)
	}
	return
}

// ResetHyperlinks reset all changes of the hyperlinks edge.
func (m *WorkOrderMutation) ResetHyperlinks() {
	m.hyperlinks = nil
	m.removedhyperlinks = nil
}

// SetLocationID sets the location edge to Location by id.
func (m *WorkOrderMutation) SetLocationID(id int) {
	m.location = &id
}

// ClearLocation clears the location edge to Location.
func (m *WorkOrderMutation) ClearLocation() {
	m.clearedlocation = true
}

// LocationCleared returns if the edge location was cleared.
func (m *WorkOrderMutation) LocationCleared() bool {
	return m.clearedlocation
}

// LocationID returns the location id in the mutation.
func (m *WorkOrderMutation) LocationID() (id int, exists bool) {
	if m.location != nil {
		return *m.location, true
	}
	return
}

// LocationIDs returns the location ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// LocationID instead. It exists only for internal usage by the builders.
func (m *WorkOrderMutation) LocationIDs() (ids []int) {
	if id := m.location; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetLocation reset all changes of the location edge.
func (m *WorkOrderMutation) ResetLocation() {
	m.location = nil
	m.clearedlocation = false
}

// AddCommentIDs adds the comments edge to Comment by ids.
func (m *WorkOrderMutation) AddCommentIDs(ids ...int) {
	if m.comments == nil {
		m.comments = make(map[int]struct{})
	}
	for i := range ids {
		m.comments[ids[i]] = struct{}{}
	}
}

// RemoveCommentIDs removes the comments edge to Comment by ids.
func (m *WorkOrderMutation) RemoveCommentIDs(ids ...int) {
	if m.removedcomments == nil {
		m.removedcomments = make(map[int]struct{})
	}
	for i := range ids {
		m.removedcomments[ids[i]] = struct{}{}
	}
}

// RemovedComments returns the removed ids of comments.
func (m *WorkOrderMutation) RemovedCommentsIDs() (ids []int) {
	for id := range m.removedcomments {
		ids = append(ids, id)
	}
	return
}

// CommentsIDs returns the comments ids in the mutation.
func (m *WorkOrderMutation) CommentsIDs() (ids []int) {
	for id := range m.comments {
		ids = append(ids, id)
	}
	return
}

// ResetComments reset all changes of the comments edge.
func (m *WorkOrderMutation) ResetComments() {
	m.comments = nil
	m.removedcomments = nil
}

// AddPropertyIDs adds the properties edge to Property by ids.
func (m *WorkOrderMutation) AddPropertyIDs(ids ...int) {
	if m.properties == nil {
		m.properties = make(map[int]struct{})
	}
	for i := range ids {
		m.properties[ids[i]] = struct{}{}
	}
}

// RemovePropertyIDs removes the properties edge to Property by ids.
func (m *WorkOrderMutation) RemovePropertyIDs(ids ...int) {
	if m.removedproperties == nil {
		m.removedproperties = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperties[ids[i]] = struct{}{}
	}
}

// RemovedProperties returns the removed ids of properties.
func (m *WorkOrderMutation) RemovedPropertiesIDs() (ids []int) {
	for id := range m.removedproperties {
		ids = append(ids, id)
	}
	return
}

// PropertiesIDs returns the properties ids in the mutation.
func (m *WorkOrderMutation) PropertiesIDs() (ids []int) {
	for id := range m.properties {
		ids = append(ids, id)
	}
	return
}

// ResetProperties reset all changes of the properties edge.
func (m *WorkOrderMutation) ResetProperties() {
	m.properties = nil
	m.removedproperties = nil
}

// AddCheckListCategoryIDs adds the check_list_categories edge to CheckListCategory by ids.
func (m *WorkOrderMutation) AddCheckListCategoryIDs(ids ...int) {
	if m.check_list_categories == nil {
		m.check_list_categories = make(map[int]struct{})
	}
	for i := range ids {
		m.check_list_categories[ids[i]] = struct{}{}
	}
}

// RemoveCheckListCategoryIDs removes the check_list_categories edge to CheckListCategory by ids.
func (m *WorkOrderMutation) RemoveCheckListCategoryIDs(ids ...int) {
	if m.removedcheck_list_categories == nil {
		m.removedcheck_list_categories = make(map[int]struct{})
	}
	for i := range ids {
		m.removedcheck_list_categories[ids[i]] = struct{}{}
	}
}

// RemovedCheckListCategories returns the removed ids of check_list_categories.
func (m *WorkOrderMutation) RemovedCheckListCategoriesIDs() (ids []int) {
	for id := range m.removedcheck_list_categories {
		ids = append(ids, id)
	}
	return
}

// CheckListCategoriesIDs returns the check_list_categories ids in the mutation.
func (m *WorkOrderMutation) CheckListCategoriesIDs() (ids []int) {
	for id := range m.check_list_categories {
		ids = append(ids, id)
	}
	return
}

// ResetCheckListCategories reset all changes of the check_list_categories edge.
func (m *WorkOrderMutation) ResetCheckListCategories() {
	m.check_list_categories = nil
	m.removedcheck_list_categories = nil
}

// AddCheckListItemIDs adds the check_list_items edge to CheckListItem by ids.
func (m *WorkOrderMutation) AddCheckListItemIDs(ids ...int) {
	if m.check_list_items == nil {
		m.check_list_items = make(map[int]struct{})
	}
	for i := range ids {
		m.check_list_items[ids[i]] = struct{}{}
	}
}

// RemoveCheckListItemIDs removes the check_list_items edge to CheckListItem by ids.
func (m *WorkOrderMutation) RemoveCheckListItemIDs(ids ...int) {
	if m.removedcheck_list_items == nil {
		m.removedcheck_list_items = make(map[int]struct{})
	}
	for i := range ids {
		m.removedcheck_list_items[ids[i]] = struct{}{}
	}
}

// RemovedCheckListItems returns the removed ids of check_list_items.
func (m *WorkOrderMutation) RemovedCheckListItemsIDs() (ids []int) {
	for id := range m.removedcheck_list_items {
		ids = append(ids, id)
	}
	return
}

// CheckListItemsIDs returns the check_list_items ids in the mutation.
func (m *WorkOrderMutation) CheckListItemsIDs() (ids []int) {
	for id := range m.check_list_items {
		ids = append(ids, id)
	}
	return
}

// ResetCheckListItems reset all changes of the check_list_items edge.
func (m *WorkOrderMutation) ResetCheckListItems() {
	m.check_list_items = nil
	m.removedcheck_list_items = nil
}

// SetTechnicianID sets the technician edge to Technician by id.
func (m *WorkOrderMutation) SetTechnicianID(id int) {
	m.technician = &id
}

// ClearTechnician clears the technician edge to Technician.
func (m *WorkOrderMutation) ClearTechnician() {
	m.clearedtechnician = true
}

// TechnicianCleared returns if the edge technician was cleared.
func (m *WorkOrderMutation) TechnicianCleared() bool {
	return m.clearedtechnician
}

// TechnicianID returns the technician id in the mutation.
func (m *WorkOrderMutation) TechnicianID() (id int, exists bool) {
	if m.technician != nil {
		return *m.technician, true
	}
	return
}

// TechnicianIDs returns the technician ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// TechnicianID instead. It exists only for internal usage by the builders.
func (m *WorkOrderMutation) TechnicianIDs() (ids []int) {
	if id := m.technician; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetTechnician reset all changes of the technician edge.
func (m *WorkOrderMutation) ResetTechnician() {
	m.technician = nil
	m.clearedtechnician = false
}

// SetProjectID sets the project edge to Project by id.
func (m *WorkOrderMutation) SetProjectID(id int) {
	m.project = &id
}

// ClearProject clears the project edge to Project.
func (m *WorkOrderMutation) ClearProject() {
	m.clearedproject = true
}

// ProjectCleared returns if the edge project was cleared.
func (m *WorkOrderMutation) ProjectCleared() bool {
	return m.clearedproject
}

// ProjectID returns the project id in the mutation.
func (m *WorkOrderMutation) ProjectID() (id int, exists bool) {
	if m.project != nil {
		return *m.project, true
	}
	return
}

// ProjectIDs returns the project ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ProjectID instead. It exists only for internal usage by the builders.
func (m *WorkOrderMutation) ProjectIDs() (ids []int) {
	if id := m.project; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetProject reset all changes of the project edge.
func (m *WorkOrderMutation) ResetProject() {
	m.project = nil
	m.clearedproject = false
}

// SetOwnerID sets the owner edge to User by id.
func (m *WorkOrderMutation) SetOwnerID(id int) {
	m.owner = &id
}

// ClearOwner clears the owner edge to User.
func (m *WorkOrderMutation) ClearOwner() {
	m.clearedowner = true
}

// OwnerCleared returns if the edge owner was cleared.
func (m *WorkOrderMutation) OwnerCleared() bool {
	return m.clearedowner
}

// OwnerID returns the owner id in the mutation.
func (m *WorkOrderMutation) OwnerID() (id int, exists bool) {
	if m.owner != nil {
		return *m.owner, true
	}
	return
}

// OwnerIDs returns the owner ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// OwnerID instead. It exists only for internal usage by the builders.
func (m *WorkOrderMutation) OwnerIDs() (ids []int) {
	if id := m.owner; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetOwner reset all changes of the owner edge.
func (m *WorkOrderMutation) ResetOwner() {
	m.owner = nil
	m.clearedowner = false
}

// SetAssigneeID sets the assignee edge to User by id.
func (m *WorkOrderMutation) SetAssigneeID(id int) {
	m.assignee = &id
}

// ClearAssignee clears the assignee edge to User.
func (m *WorkOrderMutation) ClearAssignee() {
	m.clearedassignee = true
}

// AssigneeCleared returns if the edge assignee was cleared.
func (m *WorkOrderMutation) AssigneeCleared() bool {
	return m.clearedassignee
}

// AssigneeID returns the assignee id in the mutation.
func (m *WorkOrderMutation) AssigneeID() (id int, exists bool) {
	if m.assignee != nil {
		return *m.assignee, true
	}
	return
}

// AssigneeIDs returns the assignee ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// AssigneeID instead. It exists only for internal usage by the builders.
func (m *WorkOrderMutation) AssigneeIDs() (ids []int) {
	if id := m.assignee; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetAssignee reset all changes of the assignee edge.
func (m *WorkOrderMutation) ResetAssignee() {
	m.assignee = nil
	m.clearedassignee = false
}

// Op returns the operation name.
func (m *WorkOrderMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (WorkOrder).
func (m *WorkOrderMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *WorkOrderMutation) Fields() []string {
	fields := make([]string, 0, 10)
	if m.create_time != nil {
		fields = append(fields, workorder.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, workorder.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, workorder.FieldName)
	}
	if m.status != nil {
		fields = append(fields, workorder.FieldStatus)
	}
	if m.priority != nil {
		fields = append(fields, workorder.FieldPriority)
	}
	if m.description != nil {
		fields = append(fields, workorder.FieldDescription)
	}
	if m.install_date != nil {
		fields = append(fields, workorder.FieldInstallDate)
	}
	if m.creation_date != nil {
		fields = append(fields, workorder.FieldCreationDate)
	}
	if m.index != nil {
		fields = append(fields, workorder.FieldIndex)
	}
	if m.close_date != nil {
		fields = append(fields, workorder.FieldCloseDate)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *WorkOrderMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case workorder.FieldCreateTime:
		return m.CreateTime()
	case workorder.FieldUpdateTime:
		return m.UpdateTime()
	case workorder.FieldName:
		return m.Name()
	case workorder.FieldStatus:
		return m.Status()
	case workorder.FieldPriority:
		return m.Priority()
	case workorder.FieldDescription:
		return m.Description()
	case workorder.FieldInstallDate:
		return m.InstallDate()
	case workorder.FieldCreationDate:
		return m.CreationDate()
	case workorder.FieldIndex:
		return m.Index()
	case workorder.FieldCloseDate:
		return m.CloseDate()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *WorkOrderMutation) SetField(name string, value ent.Value) error {
	switch name {
	case workorder.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case workorder.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case workorder.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case workorder.FieldStatus:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case workorder.FieldPriority:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPriority(v)
		return nil
	case workorder.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	case workorder.FieldInstallDate:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetInstallDate(v)
		return nil
	case workorder.FieldCreationDate:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreationDate(v)
		return nil
	case workorder.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIndex(v)
		return nil
	case workorder.FieldCloseDate:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCloseDate(v)
		return nil
	}
	return fmt.Errorf("unknown WorkOrder field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *WorkOrderMutation) AddedFields() []string {
	var fields []string
	if m.addindex != nil {
		fields = append(fields, workorder.FieldIndex)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *WorkOrderMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case workorder.FieldIndex:
		return m.AddedIndex()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *WorkOrderMutation) AddField(name string, value ent.Value) error {
	switch name {
	case workorder.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIndex(v)
		return nil
	}
	return fmt.Errorf("unknown WorkOrder numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *WorkOrderMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[workorder.FieldDescription] {
		fields = append(fields, workorder.FieldDescription)
	}
	if m.clearedFields[workorder.FieldInstallDate] {
		fields = append(fields, workorder.FieldInstallDate)
	}
	if m.clearedFields[workorder.FieldIndex] {
		fields = append(fields, workorder.FieldIndex)
	}
	if m.clearedFields[workorder.FieldCloseDate] {
		fields = append(fields, workorder.FieldCloseDate)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *WorkOrderMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *WorkOrderMutation) ClearField(name string) error {
	switch name {
	case workorder.FieldDescription:
		m.ClearDescription()
		return nil
	case workorder.FieldInstallDate:
		m.ClearInstallDate()
		return nil
	case workorder.FieldIndex:
		m.ClearIndex()
		return nil
	case workorder.FieldCloseDate:
		m.ClearCloseDate()
		return nil
	}
	return fmt.Errorf("unknown WorkOrder nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *WorkOrderMutation) ResetField(name string) error {
	switch name {
	case workorder.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case workorder.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case workorder.FieldName:
		m.ResetName()
		return nil
	case workorder.FieldStatus:
		m.ResetStatus()
		return nil
	case workorder.FieldPriority:
		m.ResetPriority()
		return nil
	case workorder.FieldDescription:
		m.ResetDescription()
		return nil
	case workorder.FieldInstallDate:
		m.ResetInstallDate()
		return nil
	case workorder.FieldCreationDate:
		m.ResetCreationDate()
		return nil
	case workorder.FieldIndex:
		m.ResetIndex()
		return nil
	case workorder.FieldCloseDate:
		m.ResetCloseDate()
		return nil
	}
	return fmt.Errorf("unknown WorkOrder field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *WorkOrderMutation) AddedEdges() []string {
	edges := make([]string, 0, 14)
	if m._type != nil {
		edges = append(edges, workorder.EdgeType)
	}
	if m.equipment != nil {
		edges = append(edges, workorder.EdgeEquipment)
	}
	if m.links != nil {
		edges = append(edges, workorder.EdgeLinks)
	}
	if m.files != nil {
		edges = append(edges, workorder.EdgeFiles)
	}
	if m.hyperlinks != nil {
		edges = append(edges, workorder.EdgeHyperlinks)
	}
	if m.location != nil {
		edges = append(edges, workorder.EdgeLocation)
	}
	if m.comments != nil {
		edges = append(edges, workorder.EdgeComments)
	}
	if m.properties != nil {
		edges = append(edges, workorder.EdgeProperties)
	}
	if m.check_list_categories != nil {
		edges = append(edges, workorder.EdgeCheckListCategories)
	}
	if m.check_list_items != nil {
		edges = append(edges, workorder.EdgeCheckListItems)
	}
	if m.technician != nil {
		edges = append(edges, workorder.EdgeTechnician)
	}
	if m.project != nil {
		edges = append(edges, workorder.EdgeProject)
	}
	if m.owner != nil {
		edges = append(edges, workorder.EdgeOwner)
	}
	if m.assignee != nil {
		edges = append(edges, workorder.EdgeAssignee)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *WorkOrderMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case workorder.EdgeType:
		if id := m._type; id != nil {
			return []ent.Value{*id}
		}
	case workorder.EdgeEquipment:
		ids := make([]ent.Value, 0, len(m.equipment))
		for id := range m.equipment {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeLinks:
		ids := make([]ent.Value, 0, len(m.links))
		for id := range m.links {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeFiles:
		ids := make([]ent.Value, 0, len(m.files))
		for id := range m.files {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeHyperlinks:
		ids := make([]ent.Value, 0, len(m.hyperlinks))
		for id := range m.hyperlinks {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeLocation:
		if id := m.location; id != nil {
			return []ent.Value{*id}
		}
	case workorder.EdgeComments:
		ids := make([]ent.Value, 0, len(m.comments))
		for id := range m.comments {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.properties))
		for id := range m.properties {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeCheckListCategories:
		ids := make([]ent.Value, 0, len(m.check_list_categories))
		for id := range m.check_list_categories {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeCheckListItems:
		ids := make([]ent.Value, 0, len(m.check_list_items))
		for id := range m.check_list_items {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeTechnician:
		if id := m.technician; id != nil {
			return []ent.Value{*id}
		}
	case workorder.EdgeProject:
		if id := m.project; id != nil {
			return []ent.Value{*id}
		}
	case workorder.EdgeOwner:
		if id := m.owner; id != nil {
			return []ent.Value{*id}
		}
	case workorder.EdgeAssignee:
		if id := m.assignee; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *WorkOrderMutation) RemovedEdges() []string {
	edges := make([]string, 0, 14)
	if m.removedequipment != nil {
		edges = append(edges, workorder.EdgeEquipment)
	}
	if m.removedlinks != nil {
		edges = append(edges, workorder.EdgeLinks)
	}
	if m.removedfiles != nil {
		edges = append(edges, workorder.EdgeFiles)
	}
	if m.removedhyperlinks != nil {
		edges = append(edges, workorder.EdgeHyperlinks)
	}
	if m.removedcomments != nil {
		edges = append(edges, workorder.EdgeComments)
	}
	if m.removedproperties != nil {
		edges = append(edges, workorder.EdgeProperties)
	}
	if m.removedcheck_list_categories != nil {
		edges = append(edges, workorder.EdgeCheckListCategories)
	}
	if m.removedcheck_list_items != nil {
		edges = append(edges, workorder.EdgeCheckListItems)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *WorkOrderMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case workorder.EdgeEquipment:
		ids := make([]ent.Value, 0, len(m.removedequipment))
		for id := range m.removedequipment {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeLinks:
		ids := make([]ent.Value, 0, len(m.removedlinks))
		for id := range m.removedlinks {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeFiles:
		ids := make([]ent.Value, 0, len(m.removedfiles))
		for id := range m.removedfiles {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeHyperlinks:
		ids := make([]ent.Value, 0, len(m.removedhyperlinks))
		for id := range m.removedhyperlinks {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeComments:
		ids := make([]ent.Value, 0, len(m.removedcomments))
		for id := range m.removedcomments {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeProperties:
		ids := make([]ent.Value, 0, len(m.removedproperties))
		for id := range m.removedproperties {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeCheckListCategories:
		ids := make([]ent.Value, 0, len(m.removedcheck_list_categories))
		for id := range m.removedcheck_list_categories {
			ids = append(ids, id)
		}
		return ids
	case workorder.EdgeCheckListItems:
		ids := make([]ent.Value, 0, len(m.removedcheck_list_items))
		for id := range m.removedcheck_list_items {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *WorkOrderMutation) ClearedEdges() []string {
	edges := make([]string, 0, 14)
	if m.cleared_type {
		edges = append(edges, workorder.EdgeType)
	}
	if m.clearedlocation {
		edges = append(edges, workorder.EdgeLocation)
	}
	if m.clearedtechnician {
		edges = append(edges, workorder.EdgeTechnician)
	}
	if m.clearedproject {
		edges = append(edges, workorder.EdgeProject)
	}
	if m.clearedowner {
		edges = append(edges, workorder.EdgeOwner)
	}
	if m.clearedassignee {
		edges = append(edges, workorder.EdgeAssignee)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *WorkOrderMutation) EdgeCleared(name string) bool {
	switch name {
	case workorder.EdgeType:
		return m.cleared_type
	case workorder.EdgeLocation:
		return m.clearedlocation
	case workorder.EdgeTechnician:
		return m.clearedtechnician
	case workorder.EdgeProject:
		return m.clearedproject
	case workorder.EdgeOwner:
		return m.clearedowner
	case workorder.EdgeAssignee:
		return m.clearedassignee
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *WorkOrderMutation) ClearEdge(name string) error {
	switch name {
	case workorder.EdgeType:
		m.ClearType()
		return nil
	case workorder.EdgeLocation:
		m.ClearLocation()
		return nil
	case workorder.EdgeTechnician:
		m.ClearTechnician()
		return nil
	case workorder.EdgeProject:
		m.ClearProject()
		return nil
	case workorder.EdgeOwner:
		m.ClearOwner()
		return nil
	case workorder.EdgeAssignee:
		m.ClearAssignee()
		return nil
	}
	return fmt.Errorf("unknown WorkOrder unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *WorkOrderMutation) ResetEdge(name string) error {
	switch name {
	case workorder.EdgeType:
		m.ResetType()
		return nil
	case workorder.EdgeEquipment:
		m.ResetEquipment()
		return nil
	case workorder.EdgeLinks:
		m.ResetLinks()
		return nil
	case workorder.EdgeFiles:
		m.ResetFiles()
		return nil
	case workorder.EdgeHyperlinks:
		m.ResetHyperlinks()
		return nil
	case workorder.EdgeLocation:
		m.ResetLocation()
		return nil
	case workorder.EdgeComments:
		m.ResetComments()
		return nil
	case workorder.EdgeProperties:
		m.ResetProperties()
		return nil
	case workorder.EdgeCheckListCategories:
		m.ResetCheckListCategories()
		return nil
	case workorder.EdgeCheckListItems:
		m.ResetCheckListItems()
		return nil
	case workorder.EdgeTechnician:
		m.ResetTechnician()
		return nil
	case workorder.EdgeProject:
		m.ResetProject()
		return nil
	case workorder.EdgeOwner:
		m.ResetOwner()
		return nil
	case workorder.EdgeAssignee:
		m.ResetAssignee()
		return nil
	}
	return fmt.Errorf("unknown WorkOrder edge %s", name)
}

// WorkOrderDefinitionMutation represents an operation that mutate the WorkOrderDefinitions
// nodes in the graph.
type WorkOrderDefinitionMutation struct {
	config
	op                  Op
	typ                 string
	id                  *int
	create_time         *time.Time
	update_time         *time.Time
	index               *int
	addindex            *int
	clearedFields       map[string]bool
	_type               *int
	cleared_type        bool
	project_type        *int
	clearedproject_type bool
}

var _ ent.Mutation = (*WorkOrderDefinitionMutation)(nil)

// newWorkOrderDefinitionMutation creates new mutation for $n.Name.
func newWorkOrderDefinitionMutation(c config, op Op) *WorkOrderDefinitionMutation {
	return &WorkOrderDefinitionMutation{
		config:        c,
		op:            op,
		typ:           TypeWorkOrderDefinition,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m WorkOrderDefinitionMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m WorkOrderDefinitionMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *WorkOrderDefinitionMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *WorkOrderDefinitionMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *WorkOrderDefinitionMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *WorkOrderDefinitionMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *WorkOrderDefinitionMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *WorkOrderDefinitionMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *WorkOrderDefinitionMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetIndex sets the index field.
func (m *WorkOrderDefinitionMutation) SetIndex(i int) {
	m.index = &i
	m.addindex = nil
}

// Index returns the index value in the mutation.
func (m *WorkOrderDefinitionMutation) Index() (r int, exists bool) {
	v := m.index
	if v == nil {
		return
	}
	return *v, true
}

// AddIndex adds i to index.
func (m *WorkOrderDefinitionMutation) AddIndex(i int) {
	if m.addindex != nil {
		*m.addindex += i
	} else {
		m.addindex = &i
	}
}

// AddedIndex returns the value that was added to the index field in this mutation.
func (m *WorkOrderDefinitionMutation) AddedIndex() (r int, exists bool) {
	v := m.addindex
	if v == nil {
		return
	}
	return *v, true
}

// ClearIndex clears the value of index.
func (m *WorkOrderDefinitionMutation) ClearIndex() {
	m.index = nil
	m.addindex = nil
	m.clearedFields[workorderdefinition.FieldIndex] = true
}

// IndexCleared returns if the field index was cleared in this mutation.
func (m *WorkOrderDefinitionMutation) IndexCleared() bool {
	return m.clearedFields[workorderdefinition.FieldIndex]
}

// ResetIndex reset all changes of the index field.
func (m *WorkOrderDefinitionMutation) ResetIndex() {
	m.index = nil
	m.addindex = nil
	delete(m.clearedFields, workorderdefinition.FieldIndex)
}

// SetTypeID sets the type edge to WorkOrderType by id.
func (m *WorkOrderDefinitionMutation) SetTypeID(id int) {
	m._type = &id
}

// ClearType clears the type edge to WorkOrderType.
func (m *WorkOrderDefinitionMutation) ClearType() {
	m.cleared_type = true
}

// TypeCleared returns if the edge type was cleared.
func (m *WorkOrderDefinitionMutation) TypeCleared() bool {
	return m.cleared_type
}

// TypeID returns the type id in the mutation.
func (m *WorkOrderDefinitionMutation) TypeID() (id int, exists bool) {
	if m._type != nil {
		return *m._type, true
	}
	return
}

// TypeIDs returns the type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// TypeID instead. It exists only for internal usage by the builders.
func (m *WorkOrderDefinitionMutation) TypeIDs() (ids []int) {
	if id := m._type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetType reset all changes of the type edge.
func (m *WorkOrderDefinitionMutation) ResetType() {
	m._type = nil
	m.cleared_type = false
}

// SetProjectTypeID sets the project_type edge to ProjectType by id.
func (m *WorkOrderDefinitionMutation) SetProjectTypeID(id int) {
	m.project_type = &id
}

// ClearProjectType clears the project_type edge to ProjectType.
func (m *WorkOrderDefinitionMutation) ClearProjectType() {
	m.clearedproject_type = true
}

// ProjectTypeCleared returns if the edge project_type was cleared.
func (m *WorkOrderDefinitionMutation) ProjectTypeCleared() bool {
	return m.clearedproject_type
}

// ProjectTypeID returns the project_type id in the mutation.
func (m *WorkOrderDefinitionMutation) ProjectTypeID() (id int, exists bool) {
	if m.project_type != nil {
		return *m.project_type, true
	}
	return
}

// ProjectTypeIDs returns the project_type ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ProjectTypeID instead. It exists only for internal usage by the builders.
func (m *WorkOrderDefinitionMutation) ProjectTypeIDs() (ids []int) {
	if id := m.project_type; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetProjectType reset all changes of the project_type edge.
func (m *WorkOrderDefinitionMutation) ResetProjectType() {
	m.project_type = nil
	m.clearedproject_type = false
}

// Op returns the operation name.
func (m *WorkOrderDefinitionMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (WorkOrderDefinition).
func (m *WorkOrderDefinitionMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *WorkOrderDefinitionMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.create_time != nil {
		fields = append(fields, workorderdefinition.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, workorderdefinition.FieldUpdateTime)
	}
	if m.index != nil {
		fields = append(fields, workorderdefinition.FieldIndex)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *WorkOrderDefinitionMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case workorderdefinition.FieldCreateTime:
		return m.CreateTime()
	case workorderdefinition.FieldUpdateTime:
		return m.UpdateTime()
	case workorderdefinition.FieldIndex:
		return m.Index()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *WorkOrderDefinitionMutation) SetField(name string, value ent.Value) error {
	switch name {
	case workorderdefinition.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case workorderdefinition.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case workorderdefinition.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetIndex(v)
		return nil
	}
	return fmt.Errorf("unknown WorkOrderDefinition field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *WorkOrderDefinitionMutation) AddedFields() []string {
	var fields []string
	if m.addindex != nil {
		fields = append(fields, workorderdefinition.FieldIndex)
	}
	return fields
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *WorkOrderDefinitionMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case workorderdefinition.FieldIndex:
		return m.AddedIndex()
	}
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *WorkOrderDefinitionMutation) AddField(name string, value ent.Value) error {
	switch name {
	case workorderdefinition.FieldIndex:
		v, ok := value.(int)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddIndex(v)
		return nil
	}
	return fmt.Errorf("unknown WorkOrderDefinition numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *WorkOrderDefinitionMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[workorderdefinition.FieldIndex] {
		fields = append(fields, workorderdefinition.FieldIndex)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *WorkOrderDefinitionMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *WorkOrderDefinitionMutation) ClearField(name string) error {
	switch name {
	case workorderdefinition.FieldIndex:
		m.ClearIndex()
		return nil
	}
	return fmt.Errorf("unknown WorkOrderDefinition nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *WorkOrderDefinitionMutation) ResetField(name string) error {
	switch name {
	case workorderdefinition.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case workorderdefinition.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case workorderdefinition.FieldIndex:
		m.ResetIndex()
		return nil
	}
	return fmt.Errorf("unknown WorkOrderDefinition field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *WorkOrderDefinitionMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m._type != nil {
		edges = append(edges, workorderdefinition.EdgeType)
	}
	if m.project_type != nil {
		edges = append(edges, workorderdefinition.EdgeProjectType)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *WorkOrderDefinitionMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case workorderdefinition.EdgeType:
		if id := m._type; id != nil {
			return []ent.Value{*id}
		}
	case workorderdefinition.EdgeProjectType:
		if id := m.project_type; id != nil {
			return []ent.Value{*id}
		}
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *WorkOrderDefinitionMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *WorkOrderDefinitionMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *WorkOrderDefinitionMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.cleared_type {
		edges = append(edges, workorderdefinition.EdgeType)
	}
	if m.clearedproject_type {
		edges = append(edges, workorderdefinition.EdgeProjectType)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *WorkOrderDefinitionMutation) EdgeCleared(name string) bool {
	switch name {
	case workorderdefinition.EdgeType:
		return m.cleared_type
	case workorderdefinition.EdgeProjectType:
		return m.clearedproject_type
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *WorkOrderDefinitionMutation) ClearEdge(name string) error {
	switch name {
	case workorderdefinition.EdgeType:
		m.ClearType()
		return nil
	case workorderdefinition.EdgeProjectType:
		m.ClearProjectType()
		return nil
	}
	return fmt.Errorf("unknown WorkOrderDefinition unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *WorkOrderDefinitionMutation) ResetEdge(name string) error {
	switch name {
	case workorderdefinition.EdgeType:
		m.ResetType()
		return nil
	case workorderdefinition.EdgeProjectType:
		m.ResetProjectType()
		return nil
	}
	return fmt.Errorf("unknown WorkOrderDefinition edge %s", name)
}

// WorkOrderTypeMutation represents an operation that mutate the WorkOrderTypes
// nodes in the graph.
type WorkOrderTypeMutation struct {
	config
	op                            Op
	typ                           string
	id                            *int
	create_time                   *time.Time
	update_time                   *time.Time
	name                          *string
	description                   *string
	clearedFields                 map[string]bool
	work_orders                   map[int]struct{}
	removedwork_orders            map[int]struct{}
	property_types                map[int]struct{}
	removedproperty_types         map[int]struct{}
	definitions                   map[int]struct{}
	removeddefinitions            map[int]struct{}
	check_list_categories         map[int]struct{}
	removedcheck_list_categories  map[int]struct{}
	check_list_definitions        map[int]struct{}
	removedcheck_list_definitions map[int]struct{}
}

var _ ent.Mutation = (*WorkOrderTypeMutation)(nil)

// newWorkOrderTypeMutation creates new mutation for $n.Name.
func newWorkOrderTypeMutation(c config, op Op) *WorkOrderTypeMutation {
	return &WorkOrderTypeMutation{
		config:        c,
		op:            op,
		typ:           TypeWorkOrderType,
		clearedFields: make(map[string]bool),
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m WorkOrderTypeMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m WorkOrderTypeMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *WorkOrderTypeMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetCreateTime sets the create_time field.
func (m *WorkOrderTypeMutation) SetCreateTime(t time.Time) {
	m.create_time = &t
}

// CreateTime returns the create_time value in the mutation.
func (m *WorkOrderTypeMutation) CreateTime() (r time.Time, exists bool) {
	v := m.create_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetCreateTime reset all changes of the create_time field.
func (m *WorkOrderTypeMutation) ResetCreateTime() {
	m.create_time = nil
}

// SetUpdateTime sets the update_time field.
func (m *WorkOrderTypeMutation) SetUpdateTime(t time.Time) {
	m.update_time = &t
}

// UpdateTime returns the update_time value in the mutation.
func (m *WorkOrderTypeMutation) UpdateTime() (r time.Time, exists bool) {
	v := m.update_time
	if v == nil {
		return
	}
	return *v, true
}

// ResetUpdateTime reset all changes of the update_time field.
func (m *WorkOrderTypeMutation) ResetUpdateTime() {
	m.update_time = nil
}

// SetName sets the name field.
func (m *WorkOrderTypeMutation) SetName(s string) {
	m.name = &s
}

// Name returns the name value in the mutation.
func (m *WorkOrderTypeMutation) Name() (r string, exists bool) {
	v := m.name
	if v == nil {
		return
	}
	return *v, true
}

// ResetName reset all changes of the name field.
func (m *WorkOrderTypeMutation) ResetName() {
	m.name = nil
}

// SetDescription sets the description field.
func (m *WorkOrderTypeMutation) SetDescription(s string) {
	m.description = &s
}

// Description returns the description value in the mutation.
func (m *WorkOrderTypeMutation) Description() (r string, exists bool) {
	v := m.description
	if v == nil {
		return
	}
	return *v, true
}

// ClearDescription clears the value of description.
func (m *WorkOrderTypeMutation) ClearDescription() {
	m.description = nil
	m.clearedFields[workordertype.FieldDescription] = true
}

// DescriptionCleared returns if the field description was cleared in this mutation.
func (m *WorkOrderTypeMutation) DescriptionCleared() bool {
	return m.clearedFields[workordertype.FieldDescription]
}

// ResetDescription reset all changes of the description field.
func (m *WorkOrderTypeMutation) ResetDescription() {
	m.description = nil
	delete(m.clearedFields, workordertype.FieldDescription)
}

// AddWorkOrderIDs adds the work_orders edge to WorkOrder by ids.
func (m *WorkOrderTypeMutation) AddWorkOrderIDs(ids ...int) {
	if m.work_orders == nil {
		m.work_orders = make(map[int]struct{})
	}
	for i := range ids {
		m.work_orders[ids[i]] = struct{}{}
	}
}

// RemoveWorkOrderIDs removes the work_orders edge to WorkOrder by ids.
func (m *WorkOrderTypeMutation) RemoveWorkOrderIDs(ids ...int) {
	if m.removedwork_orders == nil {
		m.removedwork_orders = make(map[int]struct{})
	}
	for i := range ids {
		m.removedwork_orders[ids[i]] = struct{}{}
	}
}

// RemovedWorkOrders returns the removed ids of work_orders.
func (m *WorkOrderTypeMutation) RemovedWorkOrdersIDs() (ids []int) {
	for id := range m.removedwork_orders {
		ids = append(ids, id)
	}
	return
}

// WorkOrdersIDs returns the work_orders ids in the mutation.
func (m *WorkOrderTypeMutation) WorkOrdersIDs() (ids []int) {
	for id := range m.work_orders {
		ids = append(ids, id)
	}
	return
}

// ResetWorkOrders reset all changes of the work_orders edge.
func (m *WorkOrderTypeMutation) ResetWorkOrders() {
	m.work_orders = nil
	m.removedwork_orders = nil
}

// AddPropertyTypeIDs adds the property_types edge to PropertyType by ids.
func (m *WorkOrderTypeMutation) AddPropertyTypeIDs(ids ...int) {
	if m.property_types == nil {
		m.property_types = make(map[int]struct{})
	}
	for i := range ids {
		m.property_types[ids[i]] = struct{}{}
	}
}

// RemovePropertyTypeIDs removes the property_types edge to PropertyType by ids.
func (m *WorkOrderTypeMutation) RemovePropertyTypeIDs(ids ...int) {
	if m.removedproperty_types == nil {
		m.removedproperty_types = make(map[int]struct{})
	}
	for i := range ids {
		m.removedproperty_types[ids[i]] = struct{}{}
	}
}

// RemovedPropertyTypes returns the removed ids of property_types.
func (m *WorkOrderTypeMutation) RemovedPropertyTypesIDs() (ids []int) {
	for id := range m.removedproperty_types {
		ids = append(ids, id)
	}
	return
}

// PropertyTypesIDs returns the property_types ids in the mutation.
func (m *WorkOrderTypeMutation) PropertyTypesIDs() (ids []int) {
	for id := range m.property_types {
		ids = append(ids, id)
	}
	return
}

// ResetPropertyTypes reset all changes of the property_types edge.
func (m *WorkOrderTypeMutation) ResetPropertyTypes() {
	m.property_types = nil
	m.removedproperty_types = nil
}

// AddDefinitionIDs adds the definitions edge to WorkOrderDefinition by ids.
func (m *WorkOrderTypeMutation) AddDefinitionIDs(ids ...int) {
	if m.definitions == nil {
		m.definitions = make(map[int]struct{})
	}
	for i := range ids {
		m.definitions[ids[i]] = struct{}{}
	}
}

// RemoveDefinitionIDs removes the definitions edge to WorkOrderDefinition by ids.
func (m *WorkOrderTypeMutation) RemoveDefinitionIDs(ids ...int) {
	if m.removeddefinitions == nil {
		m.removeddefinitions = make(map[int]struct{})
	}
	for i := range ids {
		m.removeddefinitions[ids[i]] = struct{}{}
	}
}

// RemovedDefinitions returns the removed ids of definitions.
func (m *WorkOrderTypeMutation) RemovedDefinitionsIDs() (ids []int) {
	for id := range m.removeddefinitions {
		ids = append(ids, id)
	}
	return
}

// DefinitionsIDs returns the definitions ids in the mutation.
func (m *WorkOrderTypeMutation) DefinitionsIDs() (ids []int) {
	for id := range m.definitions {
		ids = append(ids, id)
	}
	return
}

// ResetDefinitions reset all changes of the definitions edge.
func (m *WorkOrderTypeMutation) ResetDefinitions() {
	m.definitions = nil
	m.removeddefinitions = nil
}

// AddCheckListCategoryIDs adds the check_list_categories edge to CheckListCategory by ids.
func (m *WorkOrderTypeMutation) AddCheckListCategoryIDs(ids ...int) {
	if m.check_list_categories == nil {
		m.check_list_categories = make(map[int]struct{})
	}
	for i := range ids {
		m.check_list_categories[ids[i]] = struct{}{}
	}
}

// RemoveCheckListCategoryIDs removes the check_list_categories edge to CheckListCategory by ids.
func (m *WorkOrderTypeMutation) RemoveCheckListCategoryIDs(ids ...int) {
	if m.removedcheck_list_categories == nil {
		m.removedcheck_list_categories = make(map[int]struct{})
	}
	for i := range ids {
		m.removedcheck_list_categories[ids[i]] = struct{}{}
	}
}

// RemovedCheckListCategories returns the removed ids of check_list_categories.
func (m *WorkOrderTypeMutation) RemovedCheckListCategoriesIDs() (ids []int) {
	for id := range m.removedcheck_list_categories {
		ids = append(ids, id)
	}
	return
}

// CheckListCategoriesIDs returns the check_list_categories ids in the mutation.
func (m *WorkOrderTypeMutation) CheckListCategoriesIDs() (ids []int) {
	for id := range m.check_list_categories {
		ids = append(ids, id)
	}
	return
}

// ResetCheckListCategories reset all changes of the check_list_categories edge.
func (m *WorkOrderTypeMutation) ResetCheckListCategories() {
	m.check_list_categories = nil
	m.removedcheck_list_categories = nil
}

// AddCheckListDefinitionIDs adds the check_list_definitions edge to CheckListItemDefinition by ids.
func (m *WorkOrderTypeMutation) AddCheckListDefinitionIDs(ids ...int) {
	if m.check_list_definitions == nil {
		m.check_list_definitions = make(map[int]struct{})
	}
	for i := range ids {
		m.check_list_definitions[ids[i]] = struct{}{}
	}
}

// RemoveCheckListDefinitionIDs removes the check_list_definitions edge to CheckListItemDefinition by ids.
func (m *WorkOrderTypeMutation) RemoveCheckListDefinitionIDs(ids ...int) {
	if m.removedcheck_list_definitions == nil {
		m.removedcheck_list_definitions = make(map[int]struct{})
	}
	for i := range ids {
		m.removedcheck_list_definitions[ids[i]] = struct{}{}
	}
}

// RemovedCheckListDefinitions returns the removed ids of check_list_definitions.
func (m *WorkOrderTypeMutation) RemovedCheckListDefinitionsIDs() (ids []int) {
	for id := range m.removedcheck_list_definitions {
		ids = append(ids, id)
	}
	return
}

// CheckListDefinitionsIDs returns the check_list_definitions ids in the mutation.
func (m *WorkOrderTypeMutation) CheckListDefinitionsIDs() (ids []int) {
	for id := range m.check_list_definitions {
		ids = append(ids, id)
	}
	return
}

// ResetCheckListDefinitions reset all changes of the check_list_definitions edge.
func (m *WorkOrderTypeMutation) ResetCheckListDefinitions() {
	m.check_list_definitions = nil
	m.removedcheck_list_definitions = nil
}

// Op returns the operation name.
func (m *WorkOrderTypeMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (WorkOrderType).
func (m *WorkOrderTypeMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *WorkOrderTypeMutation) Fields() []string {
	fields := make([]string, 0, 4)
	if m.create_time != nil {
		fields = append(fields, workordertype.FieldCreateTime)
	}
	if m.update_time != nil {
		fields = append(fields, workordertype.FieldUpdateTime)
	}
	if m.name != nil {
		fields = append(fields, workordertype.FieldName)
	}
	if m.description != nil {
		fields = append(fields, workordertype.FieldDescription)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *WorkOrderTypeMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case workordertype.FieldCreateTime:
		return m.CreateTime()
	case workordertype.FieldUpdateTime:
		return m.UpdateTime()
	case workordertype.FieldName:
		return m.Name()
	case workordertype.FieldDescription:
		return m.Description()
	}
	return nil, false
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *WorkOrderTypeMutation) SetField(name string, value ent.Value) error {
	switch name {
	case workordertype.FieldCreateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetCreateTime(v)
		return nil
	case workordertype.FieldUpdateTime:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetUpdateTime(v)
		return nil
	case workordertype.FieldName:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetName(v)
		return nil
	case workordertype.FieldDescription:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetDescription(v)
		return nil
	}
	return fmt.Errorf("unknown WorkOrderType field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *WorkOrderTypeMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *WorkOrderTypeMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *WorkOrderTypeMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown WorkOrderType numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *WorkOrderTypeMutation) ClearedFields() []string {
	var fields []string
	if m.clearedFields[workordertype.FieldDescription] {
		fields = append(fields, workordertype.FieldDescription)
	}
	return fields
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *WorkOrderTypeMutation) FieldCleared(name string) bool {
	return m.clearedFields[name]
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *WorkOrderTypeMutation) ClearField(name string) error {
	switch name {
	case workordertype.FieldDescription:
		m.ClearDescription()
		return nil
	}
	return fmt.Errorf("unknown WorkOrderType nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *WorkOrderTypeMutation) ResetField(name string) error {
	switch name {
	case workordertype.FieldCreateTime:
		m.ResetCreateTime()
		return nil
	case workordertype.FieldUpdateTime:
		m.ResetUpdateTime()
		return nil
	case workordertype.FieldName:
		m.ResetName()
		return nil
	case workordertype.FieldDescription:
		m.ResetDescription()
		return nil
	}
	return fmt.Errorf("unknown WorkOrderType field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *WorkOrderTypeMutation) AddedEdges() []string {
	edges := make([]string, 0, 5)
	if m.work_orders != nil {
		edges = append(edges, workordertype.EdgeWorkOrders)
	}
	if m.property_types != nil {
		edges = append(edges, workordertype.EdgePropertyTypes)
	}
	if m.definitions != nil {
		edges = append(edges, workordertype.EdgeDefinitions)
	}
	if m.check_list_categories != nil {
		edges = append(edges, workordertype.EdgeCheckListCategories)
	}
	if m.check_list_definitions != nil {
		edges = append(edges, workordertype.EdgeCheckListDefinitions)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *WorkOrderTypeMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case workordertype.EdgeWorkOrders:
		ids := make([]ent.Value, 0, len(m.work_orders))
		for id := range m.work_orders {
			ids = append(ids, id)
		}
		return ids
	case workordertype.EdgePropertyTypes:
		ids := make([]ent.Value, 0, len(m.property_types))
		for id := range m.property_types {
			ids = append(ids, id)
		}
		return ids
	case workordertype.EdgeDefinitions:
		ids := make([]ent.Value, 0, len(m.definitions))
		for id := range m.definitions {
			ids = append(ids, id)
		}
		return ids
	case workordertype.EdgeCheckListCategories:
		ids := make([]ent.Value, 0, len(m.check_list_categories))
		for id := range m.check_list_categories {
			ids = append(ids, id)
		}
		return ids
	case workordertype.EdgeCheckListDefinitions:
		ids := make([]ent.Value, 0, len(m.check_list_definitions))
		for id := range m.check_list_definitions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *WorkOrderTypeMutation) RemovedEdges() []string {
	edges := make([]string, 0, 5)
	if m.removedwork_orders != nil {
		edges = append(edges, workordertype.EdgeWorkOrders)
	}
	if m.removedproperty_types != nil {
		edges = append(edges, workordertype.EdgePropertyTypes)
	}
	if m.removeddefinitions != nil {
		edges = append(edges, workordertype.EdgeDefinitions)
	}
	if m.removedcheck_list_categories != nil {
		edges = append(edges, workordertype.EdgeCheckListCategories)
	}
	if m.removedcheck_list_definitions != nil {
		edges = append(edges, workordertype.EdgeCheckListDefinitions)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *WorkOrderTypeMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case workordertype.EdgeWorkOrders:
		ids := make([]ent.Value, 0, len(m.removedwork_orders))
		for id := range m.removedwork_orders {
			ids = append(ids, id)
		}
		return ids
	case workordertype.EdgePropertyTypes:
		ids := make([]ent.Value, 0, len(m.removedproperty_types))
		for id := range m.removedproperty_types {
			ids = append(ids, id)
		}
		return ids
	case workordertype.EdgeDefinitions:
		ids := make([]ent.Value, 0, len(m.removeddefinitions))
		for id := range m.removeddefinitions {
			ids = append(ids, id)
		}
		return ids
	case workordertype.EdgeCheckListCategories:
		ids := make([]ent.Value, 0, len(m.removedcheck_list_categories))
		for id := range m.removedcheck_list_categories {
			ids = append(ids, id)
		}
		return ids
	case workordertype.EdgeCheckListDefinitions:
		ids := make([]ent.Value, 0, len(m.removedcheck_list_definitions))
		for id := range m.removedcheck_list_definitions {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *WorkOrderTypeMutation) ClearedEdges() []string {
	edges := make([]string, 0, 5)
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *WorkOrderTypeMutation) EdgeCleared(name string) bool {
	switch name {
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *WorkOrderTypeMutation) ClearEdge(name string) error {
	switch name {
	}
	return fmt.Errorf("unknown WorkOrderType unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *WorkOrderTypeMutation) ResetEdge(name string) error {
	switch name {
	case workordertype.EdgeWorkOrders:
		m.ResetWorkOrders()
		return nil
	case workordertype.EdgePropertyTypes:
		m.ResetPropertyTypes()
		return nil
	case workordertype.EdgeDefinitions:
		m.ResetDefinitions()
		return nil
	case workordertype.EdgeCheckListCategories:
		m.ResetCheckListCategories()
		return nil
	case workordertype.EdgeCheckListDefinitions:
		m.ResetCheckListDefinitions()
		return nil
	}
	return fmt.Errorf("unknown WorkOrderType edge %s", name)
}
