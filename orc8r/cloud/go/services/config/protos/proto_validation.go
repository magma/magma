/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import "errors"

func ValidateGetOrDeleteConfigsRequest(req *GetOrDeleteConfigsRequest) error {
	if len(req.GetNetworkId()) == 0 {
		return errors.New("Network ID must be specified")
	}
	return validateConfigFilter(req.GetFilter())
}

func ValidateGetOrDeleteConfigRequest(req *GetOrDeleteConfigRequest) error {
	if len(req.GetNetworkId()) == 0 {
		return errors.New("Network ID must be specified")
	}
	if len(req.GetType()) == 0 {
		return errors.New("Config type must be specified")
	}
	if len(req.GetKey()) == 0 {
		return errors.New("Config key must be specified")
	}
	return nil
}

func ValidateListKeysForTypeRequest(req *ListKeysForTypeRequest) error {
	if len(req.GetNetworkId()) == 0 {
		return errors.New("Network ID must be specified")
	}
	if len(req.GetType()) == 0 {
		return errors.New("Config type must be specified")
	}
	return nil
}

func ValidateCreateOrUpdateConfigRequest(req *CreateOrUpdateConfigRequest) error {
	if len(req.GetNetworkId()) == 0 {
		return errors.New("Network ID must be specified")
	}
	if len(req.GetType()) == 0 {
		return errors.New("Config type must be specified")
	}
	if len(req.GetKey()) == 0 {
		return errors.New("Config key must be specified")
	}
	if req.GetValue() == nil || len(req.GetValue()) == 0 {
		return errors.New("Config value must be specified and non-empty")
	}
	return nil
}

func validateConfigFilter(filter *ConfigFilter) error {
	if filter == nil {
		return errors.New("Config filter must be specified.")
	}

	if len(filter.GetType()) == 0 && len(filter.GetKey()) == 0 {
		return errors.New("At least one of type or key must be specified in the config filter")
	}
	return nil
}
