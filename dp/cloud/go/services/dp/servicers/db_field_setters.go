package servicers

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"

	"magma/dp/cloud/go/services/dp/storage/db"
)

func dbFloat64OrNil(val *wrappers.DoubleValue) sql.NullFloat64 {
	if val == nil {
		return sql.NullFloat64{}
	}
	return db.MakeFloat(val.Value)
}

func dbBoolOrNil(val *wrappers.BoolValue) sql.NullBool {
	if val == nil {
		return sql.NullBool{}
	}
	return db.MakeBool(val.Value)
}

func dbStringOrNil(val *wrappers.StringValue) sql.NullString {
	if val == nil {
		return sql.NullString{}
	}
	return db.MakeString(val.Value)
}
