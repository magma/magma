// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
	"fmt"
	"io"
	"strconv"
)

type DataType string

func (e DataType) IsValid() bool {
	switch e {
	case DataTypeString:
	case DataTypeStringArray:
		return true
	}
	return false
}

func (e DataType) String() string {
	return string(e)
}

func (e *DataType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DataType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DataType", str)
	}
	return nil
}

func (e DataType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

const (
	DataTypeString      DataType = "string"
	DataTypeStringArray DataType = "stringArray"
)
