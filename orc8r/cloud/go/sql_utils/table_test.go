/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package sql_utils

import (
	"log"
	"testing"

	"github.com/lann/builder"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestColumnBuilder_ToSql(t *testing.T) {
	// psql

	actual, err := columnBuilder(postgresColumnTypeMap).
		Name("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		References("foo", "bar").
		Default("\"hello world\"").
		NotNull().
		OnDelete(ColumnOnDeleteCascade).
		ToSql()
	assert.NoError(t, err)
	expected := "pk TEXT PRIMARY KEY NOT NULL DEFAULT \"hello world\" REFERENCES foo (bar) ON DELETE CASCADE"
	assert.Equal(t, expected, actual)

	actual, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		ToSql()
	assert.NoError(t, err)
	expected = "foo BYTEA"
	assert.Equal(t, expected, actual)

	actual, err = columnBuilder(postgresColumnTypeMap).
		Name("version").
		Type(ColumnTypeInt).
		NotNull().
		Default(0).
		ToSql()
	assert.NoError(t, err)
	expected = "version INTEGER NOT NULL DEFAULT 0"
	assert.Equal(t, expected, actual)

	// mysql
	actual, err = columnBuilder(mysqlColumnTypeMap).
		Name("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		References("foo", "bar").
		Default("\"hello world\"").
		NotNull().
		OnDelete(ColumnOnDeleteCascade).
		ToSql()
	assert.NoError(t, err)
	expected = "pk TEXT PRIMARY KEY NOT NULL DEFAULT \"hello world\" REFERENCES foo (bar) ON DELETE CASCADE"
	assert.Equal(t, expected, actual)

	actual, err = columnBuilder(mysqlColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		ToSql()
	assert.NoError(t, err)
	expected = "foo LONGBLOB"
	assert.Equal(t, expected, actual)

	actual, err = columnBuilder(mysqlColumnTypeMap).
		Name("version").
		Type(ColumnTypeInt).
		NotNull().
		Default(0).
		ToSql()
	assert.NoError(t, err)
	expected = "version INT NOT NULL DEFAULT 0"
	assert.Equal(t, expected, actual)
}

func TestColumnBuilder_ToSql_Errors(t *testing.T) {
	_, err := columnBuilder(postgresColumnTypeMap).
		Type(ColumnTypeBytes).
		ToSql()
	assert.EqualError(t, err, "column name must be specified")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		ToSql()
	assert.EqualError(t, err, "column type must be specified")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes + 100).
		ToSql()
	assert.EqualError(t, err, "column type 102 not recognized")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		References("", "bar").
		ToSql()
	assert.EqualError(t, err, "reference table name and column of foreign key must not be empty")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		References("bar", "").
		ToSql()
	assert.EqualError(t, err, "reference table name and column of foreign key must not be empty")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		OnDelete(ColumnOnDeleteCascade).
		ToSql()
	assert.EqualError(t, err, "cannot specify an ON DELETE without a REFERENCES")

	_, err = columnBuilder(postgresColumnTypeMap).
		Name("foo").
		Type(ColumnTypeBytes).
		References("bar", "baz").
		OnDelete(ColumnOnDeleteCascade + 100).
		ToSql()
	assert.EqualError(t, err, "unrecognized on delete option 100")
}

func TestCreateTableBuilder_ToSql(t *testing.T) {
	// psql
	// we allow setting primary key constraint on a column and at the table
	// level even though it would be rejected by the engine
	actual, _, err := tableBuilder(postgresColumnTypeMap).
		Name("foobar").
		IfNotExists().
		StartColumn("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		EndColumn().
		StartColumn("foo").
		Type(ColumnTypeBytes).
		NotNull().
		References("barbaz", "bites").
		OnDelete(ColumnOnDeleteCascade).
		EndColumn().
		StartColumn("bar").
		Type(ColumnTypeInt).
		Default(42).
		EndColumn().
		PrimaryKey("pk", "foo").
		Unique("foo", "bar").
		ToSql()
	assert.NoError(t, err)
	expected := "CREATE TABLE IF NOT EXISTS foobar (\n" +
		"pk TEXT PRIMARY KEY,\n" +
		"foo BYTEA NOT NULL REFERENCES barbaz (bites) ON DELETE CASCADE,\n" +
		"bar INTEGER DEFAULT 42,\n" +
		"PRIMARY KEY (pk, foo),\n" +
		"UNIQUE (foo, bar)\n" +
		")"
	assert.Equal(t, expected, actual)

	actual, _, err = tableBuilder(postgresColumnTypeMap).
		Name("foobar").
		StartColumn("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		EndColumn().
		ToSql()
	assert.NoError(t, err)
	expected = "CREATE TABLE foobar (\n" +
		"pk TEXT PRIMARY KEY\n" +
		")"
	assert.Equal(t, expected, actual)

	// mysql
	actual, _, err = tableBuilder(mysqlColumnTypeMap).
		Name("foobar").
		IfNotExists().
		StartColumn("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		EndColumn().
		StartColumn("foo").
		Type(ColumnTypeBytes).
		NotNull().
		References("barbaz", "bites").
		OnDelete(ColumnOnDeleteCascade).
		EndColumn().
		StartColumn("bar").
		Type(ColumnTypeInt).
		Default(42).
		EndColumn().
		PrimaryKey("pk", "foo").
		Unique("foo", "bar").
		ToSql()
	assert.NoError(t, err)
	expected = "CREATE TABLE IF NOT EXISTS foobar (\n" +
		"pk TEXT PRIMARY KEY,\n" +
		"foo LONGBLOB NOT NULL REFERENCES barbaz (bites) ON DELETE CASCADE,\n" +
		"bar INT DEFAULT 42,\n" +
		"PRIMARY KEY (pk, foo),\n" +
		"UNIQUE (foo, bar)\n" +
		")"
	assert.Equal(t, expected, actual)

	actual, _, err = tableBuilder(mysqlColumnTypeMap).
		Name("foobar").
		StartColumn("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		EndColumn().
		ToSql()
	assert.NoError(t, err)
	expected = "CREATE TABLE foobar (\n" +
		"pk TEXT PRIMARY KEY\n" +
		")"
	assert.Equal(t, expected, actual)
}

func TestCreateTableBuilder_Exec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error oppening stub DB conn: %s", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing stub DB: %s", err)
		}
	}()

	mock.ExpectBegin()
	mock.ExpectExec(
		"CREATE TABLE IF NOT EXISTS foobar \\(\n" +
			"pk TEXT PRIMARY KEY\n" +
			"\\)",
	).WillReturnResult(sqlmock.NewResult(1, 1))

	tx, err := db.Begin()
	assert.NoError(t, err)

	_, err = NewPostgresStatementBuilder().CreateTable("foobar").
		IfNotExists().
		StartColumn("pk").
		Type(ColumnTypeText).
		PrimaryKey().
		EndColumn().
		RunWith(tx).
		Exec()
	assert.NoError(t, err)
}

func columnBuilder(mapping map[ColumnType]string) ColumnBuilder {
	return ColumnBuilder(builder.EmptyBuilder).columnTypeNames(mapping)
}

func tableBuilder(mapping map[ColumnType]string) CreateTableBuilder {
	return CreateTableBuilder(builder.EmptyBuilder).columnTypeNames(mapping)
}
