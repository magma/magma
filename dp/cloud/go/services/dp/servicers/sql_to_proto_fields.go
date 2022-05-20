package servicers

import (
	"database/sql"

	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func protoDoubleOrNil(field sql.NullFloat64) *wrappers.DoubleValue {
	if field.Valid {
		return wrapperspb.Double(field.Float64)
	}
	return nil
}

func protoBoolOrNil(field sql.NullBool) *wrappers.BoolValue {
	if field.Valid {
		return wrapperspb.Bool(field.Bool)
	}
	return nil
}

func protoStringOrNil(field sql.NullString) *wrappers.StringValue {
	if field.Valid {
		return wrapperspb.String(field.String)
	}
	return nil
}
