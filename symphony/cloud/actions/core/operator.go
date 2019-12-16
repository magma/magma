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
	// Ex. "text", "string", "arrayString"
	DataType() string
}

type genericOperator struct {
	operatorID  string
	description string
	dataType    string
}

func (o *genericOperator) OperatorID() string {
	return o.operatorID + "-" + o.dataType
}

func (o *genericOperator) Description() string {
	return o.description
}

func (o *genericOperator) DataType() string {
	return o.dataType
}

var (
	// OperatorIsString is an implementation of Operator
	OperatorIsString = &genericOperator{"is", "is", "string"}
	// OperatorIsNotString is an implementation of Operator
	OperatorIsNotString = &genericOperator{"isNot", "is not", "string"}

	AllOperators = map[string]Operator{
		OperatorIsString.OperatorID():    OperatorIsString,
		OperatorIsNotString.OperatorID(): OperatorIsNotString,
	}
)
