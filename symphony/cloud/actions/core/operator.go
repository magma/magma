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

	// DataInput is the data type expected for the input to be.
	// Ex. "text", "string", "arrayString"
	DataInput() string
}

type genericOperator struct {
	operatorID  string
	description string
	dataInput   string
}

func (o *genericOperator) OperatorID() string {
	return o.operatorID + "-" + o.dataInput
}

func (o *genericOperator) Description() string {
	return o.description
}

func (o *genericOperator) DataInput() string {
	return o.dataInput
}

var (
	// OperatorIsString is an implementation of Operator
	OperatorIsString = &genericOperator{"is", "is", "string"}
	// OperatorIsNotString is an implementation of Operator
	OperatorIsNotString = &genericOperator{"isNot", "is not", "string"}
)
