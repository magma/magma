// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// File is the model entity for the File schema.
type File struct {
	config `gqlgen:"-" json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Type holds the value of the "type" field.
	Type string `json:"type,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty" gqlgen:"fileName"`
	// Size holds the value of the "size" field.
	Size int `json:"size,omitempty" gqlgen:"sizeInBytes"`
	// ModifiedAt holds the value of the "modified_at" field.
	ModifiedAt time.Time `json:"modified_at,omitempty" gqlgen:"modified"`
	// UploadedAt holds the value of the "uploaded_at" field.
	UploadedAt time.Time `json:"uploaded_at,omitempty" gqlgen:"uploaded"`
	// ContentType holds the value of the "content_type" field.
	ContentType string `json:"content_type,omitempty"`
	// StoreKey holds the value of the "store_key" field.
	StoreKey string `json:"store_key,omitempty"`
	// Category holds the value of the "category" field.
	Category string `json:"category,omitempty"`
}

// FromRows scans the sql response data into File.
func (f *File) FromRows(rows *sql.Rows) error {
	var scanf struct {
		ID          int
		CreateTime  sql.NullTime
		UpdateTime  sql.NullTime
		Type        sql.NullString
		Name        sql.NullString
		Size        sql.NullInt64
		ModifiedAt  sql.NullTime
		UploadedAt  sql.NullTime
		ContentType sql.NullString
		StoreKey    sql.NullString
		Category    sql.NullString
	}
	// the order here should be the same as in the `file.Columns`.
	if err := rows.Scan(
		&scanf.ID,
		&scanf.CreateTime,
		&scanf.UpdateTime,
		&scanf.Type,
		&scanf.Name,
		&scanf.Size,
		&scanf.ModifiedAt,
		&scanf.UploadedAt,
		&scanf.ContentType,
		&scanf.StoreKey,
		&scanf.Category,
	); err != nil {
		return err
	}
	f.ID = strconv.Itoa(scanf.ID)
	f.CreateTime = scanf.CreateTime.Time
	f.UpdateTime = scanf.UpdateTime.Time
	f.Type = scanf.Type.String
	f.Name = scanf.Name.String
	f.Size = int(scanf.Size.Int64)
	f.ModifiedAt = scanf.ModifiedAt.Time
	f.UploadedAt = scanf.UploadedAt.Time
	f.ContentType = scanf.ContentType.String
	f.StoreKey = scanf.StoreKey.String
	f.Category = scanf.Category.String
	return nil
}

// Update returns a builder for updating this File.
// Note that, you need to call File.Unwrap() before calling this method, if this File
// was returned from a transaction, and the transaction was committed or rolled back.
func (f *File) Update() *FileUpdateOne {
	return (&FileClient{f.config}).UpdateOne(f)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (f *File) Unwrap() *File {
	tx, ok := f.config.driver.(*txDriver)
	if !ok {
		panic("ent: File is not a transactional entity")
	}
	f.config.driver = tx.drv
	return f
}

// String implements the fmt.Stringer.
func (f *File) String() string {
	var builder strings.Builder
	builder.WriteString("File(")
	builder.WriteString(fmt.Sprintf("id=%v", f.ID))
	builder.WriteString(", create_time=")
	builder.WriteString(f.CreateTime.Format(time.ANSIC))
	builder.WriteString(", update_time=")
	builder.WriteString(f.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", type=")
	builder.WriteString(f.Type)
	builder.WriteString(", name=")
	builder.WriteString(f.Name)
	builder.WriteString(", size=")
	builder.WriteString(fmt.Sprintf("%v", f.Size))
	builder.WriteString(", modified_at=")
	builder.WriteString(f.ModifiedAt.Format(time.ANSIC))
	builder.WriteString(", uploaded_at=")
	builder.WriteString(f.UploadedAt.Format(time.ANSIC))
	builder.WriteString(", content_type=")
	builder.WriteString(f.ContentType)
	builder.WriteString(", store_key=")
	builder.WriteString(f.StoreKey)
	builder.WriteString(", category=")
	builder.WriteString(f.Category)
	builder.WriteByte(')')
	return builder.String()
}

// id returns the int representation of the ID field.
func (f *File) id() int {
	id, _ := strconv.Atoi(f.ID)
	return id
}

// Files is a parsable slice of File.
type Files []*File

// FromRows scans the sql response data into Files.
func (f *Files) FromRows(rows *sql.Rows) error {
	for rows.Next() {
		scanf := &File{}
		if err := scanf.FromRows(rows); err != nil {
			return err
		}
		*f = append(*f, scanf)
	}
	return nil
}

func (f Files) config(cfg config) {
	for _i := range f {
		f[_i].config = cfg
	}
}
