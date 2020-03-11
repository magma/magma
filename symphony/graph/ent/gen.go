// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ent

//go:generate go run github.com/facebookincubator/ent/cmd/entc generate --storage=sql --template ./template --template ../../pkg/ent-contrib/entgqlgen/template --header "// Code generated (@generated) by entc, DO NOT EDIT." ./schema
//go:generate go run github.com/google/addlicense -c Facebook -y 2004-present -l bsd ./
