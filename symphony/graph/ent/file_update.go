// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/predicate"
)

// FileUpdate is the builder for updating File entities.
type FileUpdate struct {
	config
	hooks      []Hook
	mutation   *FileMutation
	predicates []predicate.File
}

// Where adds a new predicate for the builder.
func (fu *FileUpdate) Where(ps ...predicate.File) *FileUpdate {
	fu.predicates = append(fu.predicates, ps...)
	return fu
}

// SetType sets the type field.
func (fu *FileUpdate) SetType(s string) *FileUpdate {
	fu.mutation.SetType(s)
	return fu
}

// SetName sets the name field.
func (fu *FileUpdate) SetName(s string) *FileUpdate {
	fu.mutation.SetName(s)
	return fu
}

// SetSize sets the size field.
func (fu *FileUpdate) SetSize(i int) *FileUpdate {
	fu.mutation.ResetSize()
	fu.mutation.SetSize(i)
	return fu
}

// SetNillableSize sets the size field if the given value is not nil.
func (fu *FileUpdate) SetNillableSize(i *int) *FileUpdate {
	if i != nil {
		fu.SetSize(*i)
	}
	return fu
}

// AddSize adds i to size.
func (fu *FileUpdate) AddSize(i int) *FileUpdate {
	fu.mutation.AddSize(i)
	return fu
}

// ClearSize clears the value of size.
func (fu *FileUpdate) ClearSize() *FileUpdate {
	fu.mutation.ClearSize()
	return fu
}

// SetModifiedAt sets the modified_at field.
func (fu *FileUpdate) SetModifiedAt(t time.Time) *FileUpdate {
	fu.mutation.SetModifiedAt(t)
	return fu
}

// SetNillableModifiedAt sets the modified_at field if the given value is not nil.
func (fu *FileUpdate) SetNillableModifiedAt(t *time.Time) *FileUpdate {
	if t != nil {
		fu.SetModifiedAt(*t)
	}
	return fu
}

// ClearModifiedAt clears the value of modified_at.
func (fu *FileUpdate) ClearModifiedAt() *FileUpdate {
	fu.mutation.ClearModifiedAt()
	return fu
}

// SetUploadedAt sets the uploaded_at field.
func (fu *FileUpdate) SetUploadedAt(t time.Time) *FileUpdate {
	fu.mutation.SetUploadedAt(t)
	return fu
}

// SetNillableUploadedAt sets the uploaded_at field if the given value is not nil.
func (fu *FileUpdate) SetNillableUploadedAt(t *time.Time) *FileUpdate {
	if t != nil {
		fu.SetUploadedAt(*t)
	}
	return fu
}

// ClearUploadedAt clears the value of uploaded_at.
func (fu *FileUpdate) ClearUploadedAt() *FileUpdate {
	fu.mutation.ClearUploadedAt()
	return fu
}

// SetContentType sets the content_type field.
func (fu *FileUpdate) SetContentType(s string) *FileUpdate {
	fu.mutation.SetContentType(s)
	return fu
}

// SetStoreKey sets the store_key field.
func (fu *FileUpdate) SetStoreKey(s string) *FileUpdate {
	fu.mutation.SetStoreKey(s)
	return fu
}

// SetCategory sets the category field.
func (fu *FileUpdate) SetCategory(s string) *FileUpdate {
	fu.mutation.SetCategory(s)
	return fu
}

// SetNillableCategory sets the category field if the given value is not nil.
func (fu *FileUpdate) SetNillableCategory(s *string) *FileUpdate {
	if s != nil {
		fu.SetCategory(*s)
	}
	return fu
}

// ClearCategory clears the value of category.
func (fu *FileUpdate) ClearCategory() *FileUpdate {
	fu.mutation.ClearCategory()
	return fu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fu *FileUpdate) Save(ctx context.Context) (int, error) {
	if _, ok := fu.mutation.UpdateTime(); !ok {
		v := file.UpdateDefaultUpdateTime()
		fu.mutation.SetUpdateTime(v)
	}
	if v, ok := fu.mutation.Size(); ok {
		if err := file.SizeValidator(v); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"size\": %v", err)
		}
	}
	var (
		err      error
		affected int
	)
	if len(fu.hooks) == 0 {
		affected, err = fu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FileMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fu.mutation = mutation
			affected, err = fu.sqlSave(ctx)
			return affected, err
		})
		for i := len(fu.hooks) - 1; i >= 0; i-- {
			mut = fu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (fu *FileUpdate) SaveX(ctx context.Context) int {
	affected, err := fu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fu *FileUpdate) Exec(ctx context.Context) error {
	_, err := fu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fu *FileUpdate) ExecX(ctx context.Context) {
	if err := fu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fu *FileUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   file.Table,
			Columns: file.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: file.FieldID,
			},
		},
	}
	if ps := fu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fu.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUpdateTime,
		})
	}
	if value, ok := fu.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldType,
		})
	}
	if value, ok := fu.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldName,
		})
	}
	if value, ok := fu.mutation.Size(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if value, ok := fu.mutation.AddedSize(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if fu.mutation.SizeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: file.FieldSize,
		})
	}
	if value, ok := fu.mutation.ModifiedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldModifiedAt,
		})
	}
	if fu.mutation.ModifiedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldModifiedAt,
		})
	}
	if value, ok := fu.mutation.UploadedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUploadedAt,
		})
	}
	if fu.mutation.UploadedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldUploadedAt,
		})
	}
	if value, ok := fu.mutation.ContentType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldContentType,
		})
	}
	if value, ok := fu.mutation.StoreKey(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldStoreKey,
		})
	}
	if value, ok := fu.mutation.Category(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldCategory,
		})
	}
	if fu.mutation.CategoryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: file.FieldCategory,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{file.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// FileUpdateOne is the builder for updating a single File entity.
type FileUpdateOne struct {
	config
	hooks    []Hook
	mutation *FileMutation
}

// SetType sets the type field.
func (fuo *FileUpdateOne) SetType(s string) *FileUpdateOne {
	fuo.mutation.SetType(s)
	return fuo
}

// SetName sets the name field.
func (fuo *FileUpdateOne) SetName(s string) *FileUpdateOne {
	fuo.mutation.SetName(s)
	return fuo
}

// SetSize sets the size field.
func (fuo *FileUpdateOne) SetSize(i int) *FileUpdateOne {
	fuo.mutation.ResetSize()
	fuo.mutation.SetSize(i)
	return fuo
}

// SetNillableSize sets the size field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableSize(i *int) *FileUpdateOne {
	if i != nil {
		fuo.SetSize(*i)
	}
	return fuo
}

// AddSize adds i to size.
func (fuo *FileUpdateOne) AddSize(i int) *FileUpdateOne {
	fuo.mutation.AddSize(i)
	return fuo
}

// ClearSize clears the value of size.
func (fuo *FileUpdateOne) ClearSize() *FileUpdateOne {
	fuo.mutation.ClearSize()
	return fuo
}

// SetModifiedAt sets the modified_at field.
func (fuo *FileUpdateOne) SetModifiedAt(t time.Time) *FileUpdateOne {
	fuo.mutation.SetModifiedAt(t)
	return fuo
}

// SetNillableModifiedAt sets the modified_at field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableModifiedAt(t *time.Time) *FileUpdateOne {
	if t != nil {
		fuo.SetModifiedAt(*t)
	}
	return fuo
}

// ClearModifiedAt clears the value of modified_at.
func (fuo *FileUpdateOne) ClearModifiedAt() *FileUpdateOne {
	fuo.mutation.ClearModifiedAt()
	return fuo
}

// SetUploadedAt sets the uploaded_at field.
func (fuo *FileUpdateOne) SetUploadedAt(t time.Time) *FileUpdateOne {
	fuo.mutation.SetUploadedAt(t)
	return fuo
}

// SetNillableUploadedAt sets the uploaded_at field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableUploadedAt(t *time.Time) *FileUpdateOne {
	if t != nil {
		fuo.SetUploadedAt(*t)
	}
	return fuo
}

// ClearUploadedAt clears the value of uploaded_at.
func (fuo *FileUpdateOne) ClearUploadedAt() *FileUpdateOne {
	fuo.mutation.ClearUploadedAt()
	return fuo
}

// SetContentType sets the content_type field.
func (fuo *FileUpdateOne) SetContentType(s string) *FileUpdateOne {
	fuo.mutation.SetContentType(s)
	return fuo
}

// SetStoreKey sets the store_key field.
func (fuo *FileUpdateOne) SetStoreKey(s string) *FileUpdateOne {
	fuo.mutation.SetStoreKey(s)
	return fuo
}

// SetCategory sets the category field.
func (fuo *FileUpdateOne) SetCategory(s string) *FileUpdateOne {
	fuo.mutation.SetCategory(s)
	return fuo
}

// SetNillableCategory sets the category field if the given value is not nil.
func (fuo *FileUpdateOne) SetNillableCategory(s *string) *FileUpdateOne {
	if s != nil {
		fuo.SetCategory(*s)
	}
	return fuo
}

// ClearCategory clears the value of category.
func (fuo *FileUpdateOne) ClearCategory() *FileUpdateOne {
	fuo.mutation.ClearCategory()
	return fuo
}

// Save executes the query and returns the updated entity.
func (fuo *FileUpdateOne) Save(ctx context.Context) (*File, error) {
	if _, ok := fuo.mutation.UpdateTime(); !ok {
		v := file.UpdateDefaultUpdateTime()
		fuo.mutation.SetUpdateTime(v)
	}
	if v, ok := fuo.mutation.Size(); ok {
		if err := file.SizeValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"size\": %v", err)
		}
	}
	var (
		err  error
		node *File
	)
	if len(fuo.hooks) == 0 {
		node, err = fuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FileMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fuo.mutation = mutation
			node, err = fuo.sqlSave(ctx)
			return node, err
		})
		for i := len(fuo.hooks) - 1; i >= 0; i-- {
			mut = fuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, fuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (fuo *FileUpdateOne) SaveX(ctx context.Context) *File {
	f, err := fuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return f
}

// Exec executes the query on the entity.
func (fuo *FileUpdateOne) Exec(ctx context.Context) error {
	_, err := fuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fuo *FileUpdateOne) ExecX(ctx context.Context) {
	if err := fuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fuo *FileUpdateOne) sqlSave(ctx context.Context) (f *File, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   file.Table,
			Columns: file.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: file.FieldID,
			},
		},
	}
	id, ok := fuo.mutation.ID()
	if !ok {
		return nil, fmt.Errorf("missing File.ID for update")
	}
	_spec.Node.ID.Value = id
	if value, ok := fuo.mutation.UpdateTime(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUpdateTime,
		})
	}
	if value, ok := fuo.mutation.GetType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldType,
		})
	}
	if value, ok := fuo.mutation.Name(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldName,
		})
	}
	if value, ok := fuo.mutation.Size(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if value, ok := fuo.mutation.AddedSize(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
	}
	if fuo.mutation.SizeCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: file.FieldSize,
		})
	}
	if value, ok := fuo.mutation.ModifiedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldModifiedAt,
		})
	}
	if fuo.mutation.ModifiedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldModifiedAt,
		})
	}
	if value, ok := fuo.mutation.UploadedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUploadedAt,
		})
	}
	if fuo.mutation.UploadedAtCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldUploadedAt,
		})
	}
	if value, ok := fuo.mutation.ContentType(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldContentType,
		})
	}
	if value, ok := fuo.mutation.StoreKey(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldStoreKey,
		})
	}
	if value, ok := fuo.mutation.Category(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldCategory,
		})
	}
	if fuo.mutation.CategoryCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: file.FieldCategory,
		})
	}
	f = &File{config: fuo.config}
	_spec.Assign = f.assignValues
	_spec.ScanValues = f.scanValues()
	if err = sqlgraph.UpdateNode(ctx, fuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{file.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return f, nil
}
