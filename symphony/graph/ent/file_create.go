// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebookincubator/ent/dialect/sql/sqlgraph"
	"github.com/facebookincubator/ent/schema/field"
	"github.com/facebookincubator/symphony/graph/ent/file"
)

// FileCreate is the builder for creating a File entity.
type FileCreate struct {
	config
	mutation *FileMutation
	hooks    []Hook
}

// SetCreateTime sets the create_time field.
func (fc *FileCreate) SetCreateTime(t time.Time) *FileCreate {
	fc.mutation.SetCreateTime(t)
	return fc
}

// SetNillableCreateTime sets the create_time field if the given value is not nil.
func (fc *FileCreate) SetNillableCreateTime(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetCreateTime(*t)
	}
	return fc
}

// SetUpdateTime sets the update_time field.
func (fc *FileCreate) SetUpdateTime(t time.Time) *FileCreate {
	fc.mutation.SetUpdateTime(t)
	return fc
}

// SetNillableUpdateTime sets the update_time field if the given value is not nil.
func (fc *FileCreate) SetNillableUpdateTime(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetUpdateTime(*t)
	}
	return fc
}

// SetType sets the type field.
func (fc *FileCreate) SetType(s string) *FileCreate {
	fc.mutation.SetType(s)
	return fc
}

// SetName sets the name field.
func (fc *FileCreate) SetName(s string) *FileCreate {
	fc.mutation.SetName(s)
	return fc
}

// SetSize sets the size field.
func (fc *FileCreate) SetSize(i int) *FileCreate {
	fc.mutation.SetSize(i)
	return fc
}

// SetNillableSize sets the size field if the given value is not nil.
func (fc *FileCreate) SetNillableSize(i *int) *FileCreate {
	if i != nil {
		fc.SetSize(*i)
	}
	return fc
}

// SetModifiedAt sets the modified_at field.
func (fc *FileCreate) SetModifiedAt(t time.Time) *FileCreate {
	fc.mutation.SetModifiedAt(t)
	return fc
}

// SetNillableModifiedAt sets the modified_at field if the given value is not nil.
func (fc *FileCreate) SetNillableModifiedAt(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetModifiedAt(*t)
	}
	return fc
}

// SetUploadedAt sets the uploaded_at field.
func (fc *FileCreate) SetUploadedAt(t time.Time) *FileCreate {
	fc.mutation.SetUploadedAt(t)
	return fc
}

// SetNillableUploadedAt sets the uploaded_at field if the given value is not nil.
func (fc *FileCreate) SetNillableUploadedAt(t *time.Time) *FileCreate {
	if t != nil {
		fc.SetUploadedAt(*t)
	}
	return fc
}

// SetContentType sets the content_type field.
func (fc *FileCreate) SetContentType(s string) *FileCreate {
	fc.mutation.SetContentType(s)
	return fc
}

// SetStoreKey sets the store_key field.
func (fc *FileCreate) SetStoreKey(s string) *FileCreate {
	fc.mutation.SetStoreKey(s)
	return fc
}

// SetCategory sets the category field.
func (fc *FileCreate) SetCategory(s string) *FileCreate {
	fc.mutation.SetCategory(s)
	return fc
}

// SetNillableCategory sets the category field if the given value is not nil.
func (fc *FileCreate) SetNillableCategory(s *string) *FileCreate {
	if s != nil {
		fc.SetCategory(*s)
	}
	return fc
}

// Save creates the File in the database.
func (fc *FileCreate) Save(ctx context.Context) (*File, error) {
	if _, ok := fc.mutation.CreateTime(); !ok {
		v := file.DefaultCreateTime()
		fc.mutation.SetCreateTime(v)
	}
	if _, ok := fc.mutation.UpdateTime(); !ok {
		v := file.DefaultUpdateTime()
		fc.mutation.SetUpdateTime(v)
	}
	if _, ok := fc.mutation.GetType(); !ok {
		return nil, errors.New("ent: missing required field \"type\"")
	}
	if _, ok := fc.mutation.Name(); !ok {
		return nil, errors.New("ent: missing required field \"name\"")
	}
	if v, ok := fc.mutation.Size(); ok {
		if err := file.SizeValidator(v); err != nil {
			return nil, fmt.Errorf("ent: validator failed for field \"size\": %v", err)
		}
	}
	if _, ok := fc.mutation.ContentType(); !ok {
		return nil, errors.New("ent: missing required field \"content_type\"")
	}
	if _, ok := fc.mutation.StoreKey(); !ok {
		return nil, errors.New("ent: missing required field \"store_key\"")
	}
	var (
		err  error
		node *File
	)
	if len(fc.hooks) == 0 {
		node, err = fc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*FileMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			fc.mutation = mutation
			node, err = fc.sqlSave(ctx)
			return node, err
		})
		for i := len(fc.hooks); i > 0; i-- {
			mut = fc.hooks[i-1](mut)
		}
		if _, err := mut.Mutate(ctx, fc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (fc *FileCreate) SaveX(ctx context.Context) *File {
	v, err := fc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (fc *FileCreate) sqlSave(ctx context.Context) (*File, error) {
	var (
		f     = &File{config: fc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: file.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: file.FieldID,
			},
		}
	)
	if value, ok := fc.mutation.CreateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldCreateTime,
		})
		f.CreateTime = value
	}
	if value, ok := fc.mutation.UpdateTime(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUpdateTime,
		})
		f.UpdateTime = value
	}
	if value, ok := fc.mutation.GetType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldType,
		})
		f.Type = value
	}
	if value, ok := fc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldName,
		})
		f.Name = value
	}
	if value, ok := fc.mutation.Size(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: file.FieldSize,
		})
		f.Size = value
	}
	if value, ok := fc.mutation.ModifiedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldModifiedAt,
		})
		f.ModifiedAt = value
	}
	if value, ok := fc.mutation.UploadedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: file.FieldUploadedAt,
		})
		f.UploadedAt = value
	}
	if value, ok := fc.mutation.ContentType(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldContentType,
		})
		f.ContentType = value
	}
	if value, ok := fc.mutation.StoreKey(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldStoreKey,
		})
		f.StoreKey = value
	}
	if value, ok := fc.mutation.Category(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: file.FieldCategory,
		})
		f.Category = value
	}
	if err := sqlgraph.CreateNode(ctx, fc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	f.ID = int(id)
	return f, nil
}
