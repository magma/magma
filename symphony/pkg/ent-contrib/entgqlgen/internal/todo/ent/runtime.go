// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent/schema"
	"github.com/facebookincubator/symphony/pkg/ent-contrib/entgqlgen/internal/todo/ent/todo"
)

// The init function reads all schema descriptors with runtime
// code (default values, validators or hooks) and stitches it
// to their package variables.
func init() {
	todoFields := schema.Todo{}.Fields()
	_ = todoFields
	// todoDescText is the schema descriptor for text field.
	todoDescText := todoFields[0].Descriptor()
	// todo.TextValidator is a validator for the "text" field. It is called by the builders before save.
	todo.TextValidator = todoDescText.Validators[0].(func(string) error)
}
