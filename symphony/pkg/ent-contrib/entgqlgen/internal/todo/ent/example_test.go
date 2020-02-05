// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"log"

	"github.com/facebookincubator/ent/dialect/sql"
)

// dsn for the database. In order to run the tests locally, run the following command:
//
//	 ENT_INTEGRATION_ENDPOINT="root:pass@tcp(localhost:3306)/test?parseTime=True" go test -v
//
var dsn string

func ExampleTodo() {
	if dsn == "" {
		return
	}
	ctx := context.Background()
	drv, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed creating database client: %v", err)
	}
	defer drv.Close()
	client := NewClient(Driver(drv))
	// creating vertices for the todo's edges.
	t1 := client.Todo.
		Create().
		SetText("string").
		SaveX(ctx)
	log.Println("todo created:", t1)

	// create todo vertex with its edges.
	t := client.Todo.
		Create().
		SetText("string").
		AddChildren(t1).
		SaveX(ctx)
	log.Println("todo created:", t)

	// query edges.

	t1, err = t.QueryChildren().First(ctx)
	if err != nil {
		log.Fatalf("failed querying children: %v", err)
	}
	log.Println("children found:", t1)

	// Output:
}
