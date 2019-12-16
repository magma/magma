// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated (@generated) by entc, DO NOT EDIT.

package ent

import (
	"context"
	"log"
	"time"

	"github.com/facebookincubator/ent/dialect/sql"
)

// dsn for the database. In order to run the tests locally, run the following command:
//
//	 ENT_INTEGRATION_ENDPOINT="root:pass@tcp(localhost:3306)/test?parseTime=True" go test -v
//
var dsn string

func ExampleAuditLog() {
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
	// creating vertices for the auditlog's edges.

	// create auditlog vertex with its edges.
	al := client.AuditLog.
		Create().
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetActingUserID(1).
		SetOrganization("string").
		SetMutationType("string").
		SetObjectID("string").
		SetObjectType("string").
		SetObjectDisplayName("string").
		SetMutationData(nil).
		SetURL("string").
		SetIPAddress("string").
		SetStatus("string").
		SetStatusCode("string").
		SaveX(ctx)
	log.Println("auditlog created:", al)

	// query edges.

	// Output:
}
func ExampleTenant() {
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
	// creating vertices for the tenant's edges.

	// create tenant vertex with its edges.
	t := client.Tenant.
		Create().
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetName("string").
		SetDomains(nil).
		SetNetworks(nil).
		SetTabs(nil).
		SetSSOCert("string").
		SetSSOEntryPoint("string").
		SetSSOIssuer("string").
		SaveX(ctx)
	log.Println("tenant created:", t)

	// query edges.

	// Output:
}
func ExampleToken() {
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
	// creating vertices for the token's edges.

	// create token vertex with its edges.
	t := client.Token.
		Create().
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetValue("string").
		SaveX(ctx)
	log.Println("token created:", t)

	// query edges.

	// Output:
}
func ExampleUser() {
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
	// creating vertices for the user's edges.
	t0 := client.Token.
		Create().
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetValue("string").
		SaveX(ctx)
	log.Println("token created:", t0)

	// create user vertex with its edges.
	u := client.User.
		Create().
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetEmail("string").
		SetPassword("string").
		SetRole(1).
		SetTenant("string").
		SetNetworks(nil).
		SetTabs(nil).
		AddTokens(t0).
		SaveX(ctx)
	log.Println("user created:", u)

	// query edges.
	t0, err = u.QueryTokens().First(ctx)
	if err != nil {
		log.Fatalf("failed querying tokens: %v", err)
	}
	log.Println("tokens found:", t0)

	// Output:
}
