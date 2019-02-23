/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming

import (
	"reflect"
	"strings"

	"magma/orc8r/cloud/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

// Given a struct holding a subset of mconfig fields you want to grab and an
// existing mconfig, this function will fill in the `partialOut` struct with
// either the existing value of the mconfig corresponding to each field name
// in the struct or the zero value of each field if it isn't already defined
// in the mconfig.

// `partialOut` must be a pointer to a struct

// The fields of `partialOut` must all be pointers to the concrete mconfig type
// you are expecting in each managed field of the mconfig

// All fields of `partialOut` must be exported
func GetPartialMconfig(
	configs *protos.GatewayConfigs,
	partialOut interface{},
) error {
	reflectedPartial := reflect.Indirect(reflect.ValueOf(partialOut))
	partialType := reflectedPartial.Type()

	fieldCount := partialType.NumField()
	for fieldIdx := 0; fieldIdx < fieldCount; fieldIdx++ {
		err := setConfigValueOrDefault(configs, partialType.Field(fieldIdx).Name, reflectedPartial)
		if err != nil {
			return err
		}
	}
	return nil
}

// Reflectively fill in the mconfig with the values in `partialMconfig`.

// `partialMconfig` should be a pointer to a struct. All fields should be
// exported and correspond to the fields specified in `managedFields`.
func UpdateMconfig(partialMconfig interface{}, mconfigOut *protos.GatewayConfigs) error {
	newAnys, err := marshalPartialConfigToAnys(partialMconfig)
	if err != nil {
		return err
	}

	for fieldName, fieldValue := range newAnys {
		mconfigOut.ConfigsByKey[strings.ToLower(fieldName)] = fieldValue
	}
	return nil
}

func setConfigValueOrDefault(cfg *protos.GatewayConfigs, fieldName string, partialOut reflect.Value) error {
	if cfg == nil || cfg.ConfigsByKey == nil {
		setDefaultMconfigValueInPartial(fieldName, partialOut)
		return nil
	}

	// Need to seed the field with an empty value no matter what
	_, exists := cfg.ConfigsByKey[strings.ToLower(fieldName)]
	setDefaultMconfigValueInPartial(fieldName, partialOut)
	if exists {
		protoField := partialOut.FieldByName(fieldName).Interface().(proto.Message)
		err := ptypes.UnmarshalAny(cfg.ConfigsByKey[strings.ToLower(fieldName)], protoField)
		return err
	}
	return nil
}

func setDefaultMconfigValueInPartial(fieldName string, partialOut reflect.Value) {
	field := partialOut.FieldByName(fieldName)

	// this field should be a pointer to a struct. Get the type of the struct
	// that the field points to
	expectedType := field.Type().Elem()
	field.Set(reflect.New(expectedType))
}

func marshalPartialConfigToAnys(partialConfig interface{}) (map[string]*any.Any, error) {
	reflectedValue := reflect.Indirect(reflect.ValueOf(partialConfig))
	partialType := reflectedValue.Type()

	ret := map[string]*any.Any{}
	for fieldIdx := 0; fieldIdx < partialType.NumField(); fieldIdx++ {
		fieldName := partialType.Field(fieldIdx).Name

		fieldValue := reflectedValue.FieldByName(fieldName).Interface().(proto.Message)
		fieldAny, err := ptypes.MarshalAny(fieldValue)
		if err != nil {
			return map[string]*any.Any{}, err
		}
		ret[strings.ToLower(fieldName)] = fieldAny
	}
	return ret, nil
}
