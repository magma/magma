/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package obsidian_test

import (
	"encoding/json"
	"errors"
)

// No test cases in this file, just some common definitions

type (
	configType struct {
		Foo, Bar string
	}

	errValidateType struct {
		Msg string
	}
)

// Implementations of ConvertibleUserModel

func (e *errValidateType) MarshalBinary() (data []byte, err error) {
	return json.Marshal(e)
}

func (e *errValidateType) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *errValidateType) ValidateModel() error {
	return errors.New(e.Msg)
}

func (*errValidateType) ToServiceModel() (interface{}, error) {
	panic("implement me")
}

func (*errValidateType) FromServiceModel(serviceModel interface{}) error {
	panic("implement me")
}

func (c *configType) MarshalBinary() (data []byte, err error) {
	return json.Marshal(c)
}

func (c *configType) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}

func (*configType) ValidateModel() error {
	return nil
}

func (*configType) ToServiceModel() (interface{}, error) {
	panic("implement me")
}

func (*configType) FromServiceModel(serviceModel interface{}) error {
	panic("implement me")
}
