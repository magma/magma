package db

import (
	"database/sql"

	"magma/orc8r/cloud/go/sqorc"
)

type BaseType interface {
	baseValue() interface{}
	nullableValue() interface{}
	ptr() interface{}
	sqlType() sqorc.ColumnType
}

type IntType struct{ X *sql.NullInt64 }

func (x IntType) baseValue() interface{}     { return x.X.Int64 }
func (x IntType) nullableValue() interface{} { return *x.X }
func (x IntType) ptr() interface{}           { return x.X }
func (x IntType) sqlType() sqorc.ColumnType  { return sqorc.ColumnTypeInt }

type FloatType struct{ X *sql.NullFloat64 }

func (x FloatType) baseValue() interface{}     { return x.X.Float64 }
func (x FloatType) nullableValue() interface{} { return *x.X }
func (x FloatType) ptr() interface{}           { return x.X }
func (x FloatType) sqlType() sqorc.ColumnType  { return sqorc.ColumnTypeReal }

type StringType struct{ X *sql.NullString }

func (x StringType) baseValue() interface{}     { return x.X.String }
func (x StringType) nullableValue() interface{} { return *x.X }
func (x StringType) ptr() interface{}           { return x.X }
func (x StringType) sqlType() sqorc.ColumnType  { return sqorc.ColumnTypeText }

type BoolType struct{ X *sql.NullBool }

func (x BoolType) baseValue() interface{}     { return x.X.Bool }
func (x BoolType) nullableValue() interface{} { return *x.X }
func (x BoolType) ptr() interface{}           { return x.X }
func (x BoolType) sqlType() sqorc.ColumnType  { return sqorc.ColumnTypeBool }

type TimeType struct{ X *sql.NullTime }

func (x TimeType) baseValue() interface{}     { return x.X.Time }
func (x TimeType) nullableValue() interface{} { return *x.X }
func (x TimeType) ptr() interface{}           { return x.X }
func (x TimeType) sqlType() sqorc.ColumnType  { return sqorc.ColumnTypeDatetime }
