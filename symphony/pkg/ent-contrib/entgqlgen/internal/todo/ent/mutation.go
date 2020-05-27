// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"sync"

	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent/todo"

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
	TypeTodo = "Todo"
)

// TodoMutation represents an operation that mutate the Todos
// nodes in the graph.
type TodoMutation struct {
	config
	op              Op
	typ             string
	id              *int
	text            *string
	clearedFields   map[string]struct{}
	parent          *int
	clearedparent   bool
	children        map[int]struct{}
	removedchildren map[int]struct{}
	done            bool
	oldValue        func(context.Context) (*Todo, error)
}

var _ ent.Mutation = (*TodoMutation)(nil)

// todoOption allows to manage the mutation configuration using functional options.
type todoOption func(*TodoMutation)

// newTodoMutation creates new mutation for $n.Name.
func newTodoMutation(c config, op Op, opts ...todoOption) *TodoMutation {
	m := &TodoMutation{
		config:        c,
		op:            op,
		typ:           TypeTodo,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withTodoID sets the id field of the mutation.
func withTodoID(id int) todoOption {
	return func(m *TodoMutation) {
		var (
			err   error
			once  sync.Once
			value *Todo
		)
		m.oldValue = func(ctx context.Context) (*Todo, error) {
			once.Do(func() {
				if m.done {
					err = fmt.Errorf("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Todo.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withTodo sets the old Todo of the mutation.
func withTodo(node *Todo) todoOption {
	return func(m *TodoMutation) {
		m.oldValue = func(context.Context) (*Todo, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m TodoMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m TodoMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, fmt.Errorf("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the id value in the mutation. Note that, the id
// is available only if it was provided to the builder.
func (m *TodoMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// SetText sets the text field.
func (m *TodoMutation) SetText(s string) {
	m.text = &s
}

// Text returns the text value in the mutation.
func (m *TodoMutation) Text() (r string, exists bool) {
	v := m.text
	if v == nil {
		return
	}
	return *v, true
}

// OldText returns the old text value of the Todo.
// If the Todo object wasn't provided to the builder, the object is fetched
// from the database.
// An error is returned if the mutation operation is not UpdateOne, or database query fails.
func (m *TodoMutation) OldText(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, fmt.Errorf("OldText is allowed only on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, fmt.Errorf("OldText requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldText: %w", err)
	}
	return oldValue.Text, nil
}

// ResetText reset all changes of the "text" field.
func (m *TodoMutation) ResetText() {
	m.text = nil
}

// SetParentID sets the parent edge to Todo by id.
func (m *TodoMutation) SetParentID(id int) {
	m.parent = &id
}

// ClearParent clears the parent edge to Todo.
func (m *TodoMutation) ClearParent() {
	m.clearedparent = true
}

// ParentCleared returns if the edge parent was cleared.
func (m *TodoMutation) ParentCleared() bool {
	return m.clearedparent
}

// ParentID returns the parent id in the mutation.
func (m *TodoMutation) ParentID() (id int, exists bool) {
	if m.parent != nil {
		return *m.parent, true
	}
	return
}

// ParentIDs returns the parent ids in the mutation.
// Note that ids always returns len(ids) <= 1 for unique edges, and you should use
// ParentID instead. It exists only for internal usage by the builders.
func (m *TodoMutation) ParentIDs() (ids []int) {
	if id := m.parent; id != nil {
		ids = append(ids, *id)
	}
	return
}

// ResetParent reset all changes of the "parent" edge.
func (m *TodoMutation) ResetParent() {
	m.parent = nil
	m.clearedparent = false
}

// AddChildIDs adds the children edge to Todo by ids.
func (m *TodoMutation) AddChildIDs(ids ...int) {
	if m.children == nil {
		m.children = make(map[int]struct{})
	}
	for i := range ids {
		m.children[ids[i]] = struct{}{}
	}
}

// RemoveChildIDs removes the children edge to Todo by ids.
func (m *TodoMutation) RemoveChildIDs(ids ...int) {
	if m.removedchildren == nil {
		m.removedchildren = make(map[int]struct{})
	}
	for i := range ids {
		m.removedchildren[ids[i]] = struct{}{}
	}
}

// RemovedChildren returns the removed ids of children.
func (m *TodoMutation) RemovedChildrenIDs() (ids []int) {
	for id := range m.removedchildren {
		ids = append(ids, id)
	}
	return
}

// ChildrenIDs returns the children ids in the mutation.
func (m *TodoMutation) ChildrenIDs() (ids []int) {
	for id := range m.children {
		ids = append(ids, id)
	}
	return
}

// ResetChildren reset all changes of the "children" edge.
func (m *TodoMutation) ResetChildren() {
	m.children = nil
	m.removedchildren = nil
}

// Op returns the operation name.
func (m *TodoMutation) Op() Op {
	return m.op
}

// Type returns the node type of this mutation (Todo).
func (m *TodoMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during
// this mutation. Note that, in order to get all numeric
// fields that were in/decremented, call AddedFields().
func (m *TodoMutation) Fields() []string {
	fields := make([]string, 0, 1)
	if m.text != nil {
		fields = append(fields, todo.FieldText)
	}
	return fields
}

// Field returns the value of a field with the given name.
// The second boolean value indicates that this field was
// not set, or was not define in the schema.
func (m *TodoMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case todo.FieldText:
		return m.Text()
	}
	return nil, false
}

// OldField returns the old value of the field from the database.
// An error is returned if the mutation operation is not UpdateOne,
// or the query to the database was failed.
func (m *TodoMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case todo.FieldText:
		return m.OldText(ctx)
	}
	return nil, fmt.Errorf("unknown Todo field %s", name)
}

// SetField sets the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *TodoMutation) SetField(name string, value ent.Value) error {
	switch name {
	case todo.FieldText:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetText(v)
		return nil
	}
	return fmt.Errorf("unknown Todo field %s", name)
}

// AddedFields returns all numeric fields that were incremented
// or decremented during this mutation.
func (m *TodoMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was in/decremented
// from a field with the given name. The second value indicates
// that this field was not set, or was not define in the schema.
func (m *TodoMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value for the given name. It returns an
// error if the field is not defined in the schema, or if the
// type mismatch the field type.
func (m *TodoMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Todo numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared
// during this mutation.
func (m *TodoMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicates if this field was
// cleared in this mutation.
func (m *TodoMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value for the given name. It returns an
// error if the field is not defined in the schema.
func (m *TodoMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Todo nullable field %s", name)
}

// ResetField resets all changes in the mutation regarding the
// given field name. It returns an error if the field is not
// defined in the schema.
func (m *TodoMutation) ResetField(name string) error {
	switch name {
	case todo.FieldText:
		m.ResetText()
		return nil
	}
	return fmt.Errorf("unknown Todo field %s", name)
}

// AddedEdges returns all edge names that were set/added in this
// mutation.
func (m *TodoMutation) AddedEdges() []string {
	edges := make([]string, 0, 2)
	if m.parent != nil {
		edges = append(edges, todo.EdgeParent)
	}
	if m.children != nil {
		edges = append(edges, todo.EdgeChildren)
	}
	return edges
}

// AddedIDs returns all ids (to other nodes) that were added for
// the given edge name.
func (m *TodoMutation) AddedIDs(name string) []ent.Value {
	switch name {
	case todo.EdgeParent:
		if id := m.parent; id != nil {
			return []ent.Value{*id}
		}
	case todo.EdgeChildren:
		ids := make([]ent.Value, 0, len(m.children))
		for id := range m.children {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// RemovedEdges returns all edge names that were removed in this
// mutation.
func (m *TodoMutation) RemovedEdges() []string {
	edges := make([]string, 0, 2)
	if m.removedchildren != nil {
		edges = append(edges, todo.EdgeChildren)
	}
	return edges
}

// RemovedIDs returns all ids (to other nodes) that were removed for
// the given edge name.
func (m *TodoMutation) RemovedIDs(name string) []ent.Value {
	switch name {
	case todo.EdgeChildren:
		ids := make([]ent.Value, 0, len(m.removedchildren))
		for id := range m.removedchildren {
			ids = append(ids, id)
		}
		return ids
	}
	return nil
}

// ClearedEdges returns all edge names that were cleared in this
// mutation.
func (m *TodoMutation) ClearedEdges() []string {
	edges := make([]string, 0, 2)
	if m.clearedparent {
		edges = append(edges, todo.EdgeParent)
	}
	return edges
}

// EdgeCleared returns a boolean indicates if this edge was
// cleared in this mutation.
func (m *TodoMutation) EdgeCleared(name string) bool {
	switch name {
	case todo.EdgeParent:
		return m.clearedparent
	}
	return false
}

// ClearEdge clears the value for the given name. It returns an
// error if the edge name is not defined in the schema.
func (m *TodoMutation) ClearEdge(name string) error {
	switch name {
	case todo.EdgeParent:
		m.ClearParent()
		return nil
	}
	return fmt.Errorf("unknown Todo unique edge %s", name)
}

// ResetEdge resets all changes in the mutation regarding the
// given edge name. It returns an error if the edge is not
// defined in the schema.
func (m *TodoMutation) ResetEdge(name string) error {
	switch name {
	case todo.EdgeParent:
		m.ResetParent()
		return nil
	case todo.EdgeChildren:
		m.ResetChildren()
		return nil
	}
	return fmt.Errorf("unknown Todo edge %s", name)
}
