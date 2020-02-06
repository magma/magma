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

	update_time      *time.Time
	_type            *string
	name             *string
	size             *int
	addsize          *int
	clearsize        bool
	modified_at      *time.Time
	clearmodified_at bool
	uploaded_at      *time.Time
	clearuploaded_at bool
	content_type     *string
	store_key        *string
	category         *string
	clearcategory    bool
	predicates       []predicate.File
}

// Where adds a new predicate for the builder.
func (fu *FileUpdate) Where(ps ...predicate.File) *FileUpdate {
	fu.predicates = append(fu.predicates, ps...)
	return fu
}

// SetType sets the type field.
func (fu *FileUpdate) SetType(s string) *FileUpdate {
	fu._type = &s
	return fu
}

// SetName sets the name field.
func (fu *FileUpdate) SetName(s string) *FileUpdate {
	fu.name = &s
	return fu
}

// SetSize sets the size field.
func (fu *FileUpdate) SetSize(i int) *FileUpdate {
	fu.size = &i
	fu.addsize = nil
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
	if fu.addsize == nil {
		fu.addsize = &i
	} else {
		*fu.addsize += i
	}
	return fu
}

// ClearSize clears the value of size.
func (fu *FileUpdate) ClearSize() *FileUpdate {
	fu.size = nil
	fu.clearsize = true
	return fu
}

// SetModifiedAt sets the modified_at field.
func (fu *FileUpdate) SetModifiedAt(t time.Time) *FileUpdate {
	fu.modified_at = &t
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
	fu.modified_at = nil
	fu.clearmodified_at = true
	return fu
}

// SetUploadedAt sets the uploaded_at field.
func (fu *FileUpdate) SetUploadedAt(t time.Time) *FileUpdate {
	fu.uploaded_at = &t
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
	fu.uploaded_at = nil
	fu.clearuploaded_at = true
	return fu
}

// SetContentType sets the content_type field.
func (fu *FileUpdate) SetContentType(s string) *FileUpdate {
	fu.content_type = &s
	return fu
}

// SetStoreKey sets the store_key field.
func (fu *FileUpdate) SetStoreKey(s string) *FileUpdate {
	fu.store_key = &s
	return fu
}

// SetCategory sets the category field.
func (fu *FileUpdate) SetCategory(s string) *FileUpdate {
	fu.category = &s
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
	fu.category = nil
	fu.clearcategory = true
	return fu
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (fu *FileUpdate) Save(ctx context.Context) (int, error) {
	if fu.update_time == nil {
		v := file.UpdateDefaultUpdateTime()
		fu.update_time = &v
	}
	if fu.size != nil {
		if err := file.SizeValidator(*fu.size); err != nil {
			return 0, fmt.Errorf("ent: validator failed for field \"size\": %v", err)
		}
	}
	return fu.sqlSave(ctx)
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
				Type:   field.TypeString,
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
	if value := fu.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: file.FieldUpdateTime,
		})
	}
	if value := fu._type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: file.FieldType,
		})
	}
	if value := fu.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: file.FieldName,
		})
	}
	if value := fu.size; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: file.FieldSize,
		})
	}
	if value := fu.addsize; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: file.FieldSize,
		})
	}
	if fu.clearsize {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: file.FieldSize,
		})
	}
	if value := fu.modified_at; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: file.FieldModifiedAt,
		})
	}
	if fu.clearmodified_at {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldModifiedAt,
		})
	}
	if value := fu.uploaded_at; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: file.FieldUploadedAt,
		})
	}
	if fu.clearuploaded_at {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldUploadedAt,
		})
	}
	if value := fu.content_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: file.FieldContentType,
		})
	}
	if value := fu.store_key; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: file.FieldStoreKey,
		})
	}
	if value := fu.category; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: file.FieldCategory,
		})
	}
	if fu.clearcategory {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: file.FieldCategory,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fu.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// FileUpdateOne is the builder for updating a single File entity.
type FileUpdateOne struct {
	config
	id string

	update_time      *time.Time
	_type            *string
	name             *string
	size             *int
	addsize          *int
	clearsize        bool
	modified_at      *time.Time
	clearmodified_at bool
	uploaded_at      *time.Time
	clearuploaded_at bool
	content_type     *string
	store_key        *string
	category         *string
	clearcategory    bool
}

// SetType sets the type field.
func (fuo *FileUpdateOne) SetType(s string) *FileUpdateOne {
	fuo._type = &s
	return fuo
}

// SetName sets the name field.
func (fuo *FileUpdateOne) SetName(s string) *FileUpdateOne {
	fuo.name = &s
	return fuo
}

// SetSize sets the size field.
func (fuo *FileUpdateOne) SetSize(i int) *FileUpdateOne {
	fuo.size = &i
	fuo.addsize = nil
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
	if fuo.addsize == nil {
		fuo.addsize = &i
	} else {
		*fuo.addsize += i
	}
	return fuo
}

// ClearSize clears the value of size.
func (fuo *FileUpdateOne) ClearSize() *FileUpdateOne {
	fuo.size = nil
	fuo.clearsize = true
	return fuo
}

// SetModifiedAt sets the modified_at field.
func (fuo *FileUpdateOne) SetModifiedAt(t time.Time) *FileUpdateOne {
	fuo.modified_at = &t
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
	fuo.modified_at = nil
	fuo.clearmodified_at = true
	return fuo
}

// SetUploadedAt sets the uploaded_at field.
func (fuo *FileUpdateOne) SetUploadedAt(t time.Time) *FileUpdateOne {
	fuo.uploaded_at = &t
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
	fuo.uploaded_at = nil
	fuo.clearuploaded_at = true
	return fuo
}

// SetContentType sets the content_type field.
func (fuo *FileUpdateOne) SetContentType(s string) *FileUpdateOne {
	fuo.content_type = &s
	return fuo
}

// SetStoreKey sets the store_key field.
func (fuo *FileUpdateOne) SetStoreKey(s string) *FileUpdateOne {
	fuo.store_key = &s
	return fuo
}

// SetCategory sets the category field.
func (fuo *FileUpdateOne) SetCategory(s string) *FileUpdateOne {
	fuo.category = &s
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
	fuo.category = nil
	fuo.clearcategory = true
	return fuo
}

// Save executes the query and returns the updated entity.
func (fuo *FileUpdateOne) Save(ctx context.Context) (*File, error) {
	if fuo.update_time == nil {
		v := file.UpdateDefaultUpdateTime()
		fuo.update_time = &v
	}
	if fuo.size != nil {
		if err := file.SizeValidator(*fuo.size); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"size\": %v", err)
		}
	}
	return fuo.sqlSave(ctx)
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
				Value:  fuo.id,
				Type:   field.TypeString,
				Column: file.FieldID,
			},
		},
	}
	if value := fuo.update_time; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: file.FieldUpdateTime,
		})
	}
	if value := fuo._type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: file.FieldType,
		})
	}
	if value := fuo.name; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: file.FieldName,
		})
	}
	if value := fuo.size; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: file.FieldSize,
		})
	}
	if value := fuo.addsize; value != nil {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  *value,
			Column: file.FieldSize,
		})
	}
	if fuo.clearsize {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Column: file.FieldSize,
		})
	}
	if value := fuo.modified_at; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: file.FieldModifiedAt,
		})
	}
	if fuo.clearmodified_at {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldModifiedAt,
		})
	}
	if value := fuo.uploaded_at; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  *value,
			Column: file.FieldUploadedAt,
		})
	}
	if fuo.clearuploaded_at {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Column: file.FieldUploadedAt,
		})
	}
	if value := fuo.content_type; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: file.FieldContentType,
		})
	}
	if value := fuo.store_key; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: file.FieldStoreKey,
		})
	}
	if value := fuo.category; value != nil {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  *value,
			Column: file.FieldCategory,
		})
	}
	if fuo.clearcategory {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: file.FieldCategory,
		})
	}
	f = &File{config: fuo.config}
	_spec.Assign = f.assignValues
	_spec.ScanValues = f.scanValues()
	if err = sqlgraph.UpdateNode(ctx, fuo.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return f, nil
}
