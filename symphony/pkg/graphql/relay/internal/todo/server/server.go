// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/handler"
	"github.com/facebookincubator/symphony/pkg/graphql/relay/internal/todo"
	"github.com/facebookincubator/symphony/pkg/graphql/relay/internal/todo/ent"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatal("opening ent client", err)
	}
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalln("running schema migration", err)
	}

	http.Handle("/", handler.Playground("Todo", "/query"))
	http.Handle("/query", handler.GraphQL(
		todo.NewExecutableSchema(todo.New(client)),
	))
	log.Fatal(http.ListenAndServe(":8081", nil))
}
