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
	var (
		builder  = sql.Dialect(fu.driver.Dialect())
		selector = builder.Select(file.FieldID).From(builder.Table(file.Table))
	)
	for _, p := range fu.predicates {
		p(selector)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = fu.driver.Query(ctx, query, args, rows); err != nil {
		return 0, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("ent: failed reading id: %v", err)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return 0, nil
	}

	tx, err := fu.driver.Tx(ctx)
	if err != nil {
		return 0, err
	}
	var (
		res     sql.Result
		updater = builder.Update(file.Table)
	)
	updater = updater.Where(sql.InInts(file.FieldID, ids...))
	if value := fu.update_time; value != nil {
		updater.Set(file.FieldUpdateTime, *value)
	}
	if value := fu._type; value != nil {
		updater.Set(file.FieldType, *value)
	}
	if value := fu.name; value != nil {
		updater.Set(file.FieldName, *value)
	}
	if value := fu.size; value != nil {
		updater.Set(file.FieldSize, *value)
	}
	if value := fu.addsize; value != nil {
		updater.Add(file.FieldSize, *value)
	}
	if fu.clearsize {
		updater.SetNull(file.FieldSize)
	}
	if value := fu.modified_at; value != nil {
		updater.Set(file.FieldModifiedAt, *value)
	}
	if fu.clearmodified_at {
		updater.SetNull(file.FieldModifiedAt)
	}
	if value := fu.uploaded_at; value != nil {
		updater.Set(file.FieldUploadedAt, *value)
	}
	if fu.clearuploaded_at {
		updater.SetNull(file.FieldUploadedAt)
	}
	if value := fu.content_type; value != nil {
		updater.Set(file.FieldContentType, *value)
	}
	if value := fu.store_key; value != nil {
		updater.Set(file.FieldStoreKey, *value)
	}
	if value := fu.category; value != nil {
		updater.Set(file.FieldCategory, *value)
	}
	if fu.clearcategory {
		updater.SetNull(file.FieldCategory)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return 0, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return 0, err
	}
	return len(ids), nil
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
	var (
		builder  = sql.Dialect(fuo.driver.Dialect())
		selector = builder.Select(file.Columns...).From(builder.Table(file.Table))
	)
	file.ID(fuo.id)(selector)
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err = fuo.driver.Query(ctx, query, args, rows); err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		f = &File{config: fuo.config}
		if err := f.FromRows(rows); err != nil {
			return nil, fmt.Errorf("ent: failed scanning row into File: %v", err)
		}
		id = f.id()
		ids = append(ids, id)
	}
	switch n := len(ids); {
	case n == 0:
		return nil, &ErrNotFound{fmt.Sprintf("File with id: %v", fuo.id)}
	case n > 1:
		return nil, fmt.Errorf("ent: more than one File with the same id: %v", fuo.id)
	}

	tx, err := fuo.driver.Tx(ctx)
	if err != nil {
		return nil, err
	}
	var (
		res     sql.Result
		updater = builder.Update(file.Table)
	)
	updater = updater.Where(sql.InInts(file.FieldID, ids...))
	if value := fuo.update_time; value != nil {
		updater.Set(file.FieldUpdateTime, *value)
		f.UpdateTime = *value
	}
	if value := fuo._type; value != nil {
		updater.Set(file.FieldType, *value)
		f.Type = *value
	}
	if value := fuo.name; value != nil {
		updater.Set(file.FieldName, *value)
		f.Name = *value
	}
	if value := fuo.size; value != nil {
		updater.Set(file.FieldSize, *value)
		f.Size = *value
	}
	if value := fuo.addsize; value != nil {
		updater.Add(file.FieldSize, *value)
		f.Size += *value
	}
	if fuo.clearsize {
		var value int
		f.Size = value
		updater.SetNull(file.FieldSize)
	}
	if value := fuo.modified_at; value != nil {
		updater.Set(file.FieldModifiedAt, *value)
		f.ModifiedAt = *value
	}
	if fuo.clearmodified_at {
		var value time.Time
		f.ModifiedAt = value
		updater.SetNull(file.FieldModifiedAt)
	}
	if value := fuo.uploaded_at; value != nil {
		updater.Set(file.FieldUploadedAt, *value)
		f.UploadedAt = *value
	}
	if fuo.clearuploaded_at {
		var value time.Time
		f.UploadedAt = value
		updater.SetNull(file.FieldUploadedAt)
	}
	if value := fuo.content_type; value != nil {
		updater.Set(file.FieldContentType, *value)
		f.ContentType = *value
	}
	if value := fuo.store_key; value != nil {
		updater.Set(file.FieldStoreKey, *value)
		f.StoreKey = *value
	}
	if value := fuo.category; value != nil {
		updater.Set(file.FieldCategory, *value)
		f.Category = *value
	}
	if fuo.clearcategory {
		var value string
		f.Category = value
		updater.SetNull(file.FieldCategory)
	}
	if !updater.Empty() {
		query, args := updater.Query()
		if err := tx.Exec(ctx, query, args, &res); err != nil {
			return nil, rollback(tx, err)
		}
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return f, nil
}
