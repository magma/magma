/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package serde

import (
	"encoding"
	"reflect"

	"github.com/pkg/errors"
)

// Serde (SERializer-DEserializer) implements logic to serialize/deserialize
// a specific piece of data.
type Serde interface {
	// GetDomain returns a globally unique key which represents the domain of
	// this Serde. Serde types are unique within each domain but not across
	// domains.
	GetDomain() string

	// GetType returns a unique key within the domain for the specific Serde
	// implementation. This represents the type of data that the Serde will be
	// responsible for serializing and deserialing.
	GetType() string

	// Serialize a piece of data
	Serialize(in interface{}) ([]byte, error)

	// Deserialize a piece of data
	Deserialize(in []byte) (interface{}, error)
}

// ValidateableModel implements a ValidateModel() function that returns whether
// the instance is valid.
type ValidatableModel interface {
	ValidateModel() error
}

// ValidateableBinaryConvertible wraps both BinaryConvertible, for generic
// serde factory functions, and ValidateableModel for validations.
type ValidateableBinaryConvertible interface {
	BinaryConvertible
	ValidatableModel
}

// BinaryConvertible wraps encoding.BinaryMarshaler and
// encoding.BinaryUnmarshaler for use in generic serde factory functions.
type BinaryConvertible interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// NewBinarySerde returns a Serde implementation for a structure which
// implements BinaryConvertible. `dataInstance` is expected to be a pointer.
func NewBinarySerde(domain string, serdeType string, dataInstance BinaryConvertible) Serde {
	return &binarySerde{domain: domain, serdeType: serdeType, dataInstance: dataInstance}
}

type binarySerde struct {
	domain       string
	serdeType    string
	dataInstance BinaryConvertible
}

func (s *binarySerde) GetDomain() string {
	return s.domain
}

func (s *binarySerde) GetType() string {
	return s.serdeType
}

func (s *binarySerde) Serialize(in interface{}) ([]byte, error) {
	bm, ok := in.(BinaryConvertible)
	if !ok {
		return nil, errors.Errorf("structure does not implement BinaryConvertible")
	}
	return bm.MarshalBinary()
}

func (s *binarySerde) Deserialize(in []byte) (interface{}, error) {
	model := reflect.New(reflect.TypeOf(s.dataInstance).Elem()).Interface().(BinaryConvertible)
	err := model.UnmarshalBinary(in)
	return model, err
}
