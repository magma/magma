/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package configurator

import (
	"encoding"
	"reflect"

	"magma/orc8r/cloud/go/serde"

	"github.com/pkg/errors"
)

// NewNetworkConfigSerde returns a network config domain Serde implementation
// for a pointer to a structure implementing both BinaryMarshaler and
// BinaryUnmarshaler.
// If the modelPtr argument is not a pointer to a struct matching those
// requirements, this function will panic.
func NewNetworkConfigSerde(configType string, modelPtr interface{}) serde.Serde {
	modelFactory, err := createModelFactory(modelPtr)
	if err != nil {
		panic(err)
	}

	return &binarySerde{
		domain:       NetworkConfigSerdeDomain,
		serdeType:    configType,
		modelFactory: modelFactory,
	}
}

// NewNetworkEntityConfigSerde returns a network entity config domain Serde
// implementation/ for a pointer to a structure implementing both
// BinaryMarshaler and BinaryUnmarshaler.
// If the modelPtr argument is not a pointer to a struct matching those
// requirements, this function will panic.
func NewNetworkEntityConfigSerde(configType string, modelPtr interface{}) serde.Serde {
	modelFactory, err := createModelFactory(modelPtr)
	if err != nil {
		panic(err)
	}

	return &binarySerde{
		domain:       NetworkEntitySerdeDomain,
		serdeType:    configType,
		modelFactory: modelFactory,
	}
}

func createModelFactory(modelPtr interface{}) (func() encoding.BinaryUnmarshaler, error) {
	modelPtrType := reflect.TypeOf(modelPtr)
	if modelPtrType.Kind() != reflect.Ptr {
		return nil, errors.Errorf("expected a pointer to a model, got %v instead", modelPtrType)
	}

	// local variables so we can reflect for the Type of the interfaces
	var marshalerPtr *encoding.BinaryMarshaler
	var unmarshalerPtr *encoding.BinaryUnmarshaler
	if !modelPtrType.Implements(reflect.TypeOf(marshalerPtr).Elem()) {
		return nil, errors.Errorf("model must implement encoding.BinaryMarshaler")
	}
	if !modelPtrType.Implements(reflect.TypeOf(unmarshalerPtr).Elem()) {
		return nil, errors.Errorf("model must implement encoding.BinaryUnmarshaler")
	}

	return func() encoding.BinaryUnmarshaler {
		return reflect.New(modelPtrType.Elem()).Interface().(encoding.BinaryUnmarshaler)
	}, nil
}

type binarySerde struct {
	domain       string
	serdeType    string
	modelFactory func() encoding.BinaryUnmarshaler
}

func (s *binarySerde) GetDomain() string {
	return s.domain
}

func (s *binarySerde) GetType() string {
	return s.serdeType
}

func (s *binarySerde) Serialize(in interface{}) ([]byte, error) {
	bm, ok := in.(encoding.BinaryMarshaler)
	if !ok {
		return nil, errors.Errorf("structure does not implement BinaryMarshaler")
	}
	return bm.MarshalBinary()
}

func (s *binarySerde) Deserialize(in []byte) (interface{}, error) {
	model := s.modelFactory()
	err := model.UnmarshalBinary(in)
	return model, err
}
