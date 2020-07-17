/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package test_init

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/definitions"
	"orc8r/fbinternal/cloud/go/services/testcontroller"
	"orc8r/fbinternal/cloud/go/services/testcontroller/protos"
	"orc8r/fbinternal/cloud/go/services/testcontroller/servicers"
	"orc8r/fbinternal/cloud/go/services/testcontroller/storage"

	"github.com/golang/glog"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, testcontroller.ServiceName)

	// Connect to postgres_test
	db, err := sqorc.Open("postgres", definitions.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_test user=magma_test password=magma_test host=postgres_test sslmode=disable"))
	if err != nil {
		t.Fatalf("could not start test testcontroller service: %s", err)
	}

	// Start with a fresh DB
	tx, err := db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		t.Fatalf("could not start tx for testcontroller service cleanup")
	}
	// Don't exit the test if cleanup fails
	_, err = tx.Exec("DROP TABLE IF EXISTS testcontroller_nodes")
	if err != nil {
		glog.Errorf("drop testcontroller_nodes table error: %s", err)
	}
	_, err = tx.Exec("DROP TABLE IF EXISTS testcontroller_tests")
	if err != nil {
		glog.Errorf("drop testcontroller_tests table error: %s", err)
	}
	if err := tx.Commit(); err != nil {
		t.Fatalf("failed to commit testcontroller service cleanup Tx")
	}

	nodeStore := storage.NewSQLNodeLeasorStorage(db, &mockIDGenerator{}, sqorc.GetSqlBuilder())
	err = nodeStore.Init()
	assert.NoError(t, err)
	nodes := servicers.NewNodeLeasorServicer(nodeStore)
	protos.RegisterNodeLeasorServer(srv.GrpcServer, nodes)

	testStore := storage.NewSQLTestcontrollerStorage(db, sqorc.GetSqlBuilder())
	err = testStore.Init()
	assert.NoError(t, err)
	tests := servicers.NewTestControllerServicer(testStore)
	protos.RegisterTestControllerServer(srv.GrpcServer, tests)

	go func() {
		defer db.Close()
		srv.RunTest(lis)
	}()
}

func GetTestTestcontrollerStorage(t *testing.T) storage.TestControllerStorage {
	db, err := sqorc.Open("postgres", definitions.GetEnvWithDefault("DATABASE_SOURCE", "dbname=magma_test user=magma_test password=magma_test host=postgres_test sslmode=disable"))
	if err != nil {
		t.Fatalf("could not dial potgres_test DB %s", err)
	}
	return storage.NewSQLTestcontrollerStorage(db, sqorc.GetSqlBuilder())
}

type mockIDGenerator struct {
	current uint64
}

func (m *mockIDGenerator) New() string {
	m.current++
	return fmt.Sprintf("%d", m.current)
}
