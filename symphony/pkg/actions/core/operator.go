// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

// Operator species the way to filtering data, and is used in conjunction
// with Filters
// Ex. "is network1" or 'is greater than 2019-09-18'
type Operator interface {
	// OperatorID is a unique identifier for an operator.  Ex "is", "isNot"
	OperatorID() string

	// Description is the name displayed to the user.  Ex "is", "is not"
	Description() string

	// DataType is the data type expected for the input to be.
	// Ex. "text", "string", "stringArray"
	DataType() DataType
}

type genericOperator struct {
	operatorID  string
	description string
	dataType    DataType
}

func (o *genericOperator) OperatorID() string {
	return o.operatorID + "-" + string(o.dataType)
}

func (o *genericOperator) Description() string {
	return o.description
}

func (o *genericOperator) DataType() DataType {
	return o.dataType
}

var (
	// OperatorIsString is an implementation of Operator
	OperatorIsString = &genericOperator{"is", "is", DataTypeString}
	// OperatorIsNotString is an implementation of Operator
	OperatorIsNotString = &genericOperator{"isNot", "is not", DataTypeString}

	AllOperators = map[string]Operator{
		OperatorIsString.OperatorID():    OperatorIsString,
		OperatorIsNotString.OperatorID(): OperatorIsNotString,
	}
)
