package servicers

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"

	"magma/dp/cloud/go/services/dp/storage/db"
)

func dbFloat64OrNil(field *sql.NullFloat64, val *wrappers.DoubleValue) {
	if val != nil {
		*field = db.MakeFloat(val.Value)
	}
}

func dbBoolOrNil(field *sql.NullBool, val *wrappers.BoolValue) {
	if val != nil {
		*field = db.MakeBool(val.Value)
	}
}

func dbStringOrNil(field *sql.NullString, val *wrappers.StringValue) {
	if val != nil {
		*field = db.MakeString(val.Value)
	}
}
