/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// Code generated (@generated) by entc, DO NOT EDIT.

package migrate

import (
	"magma/orc8r/cloud/go/blobstore/ent/blob"

	"github.com/facebookincubator/ent/dialect/sql/schema"
	"github.com/facebookincubator/ent/schema/field"
)

var (
	// StatesColumns holds the columns for the "states" table.
	StatesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "network_id", Type: field.TypeString, Size: 2147483647},
		{Name: "type", Type: field.TypeString, Size: 2147483647},
		{Name: "key", Type: field.TypeString, Size: 2147483647},
		{Name: "value", Type: field.TypeBytes, Nullable: true},
		{Name: "version", Type: field.TypeUint64, Default: blob.DefaultVersion},
	}
	// StatesTable holds the schema information for the "states" table.
	StatesTable = &schema.Table{
		Name:        "states",
		Columns:     StatesColumns,
		PrimaryKey:  []*schema.Column{StatesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		StatesTable,
	}
)

func init() {
}
