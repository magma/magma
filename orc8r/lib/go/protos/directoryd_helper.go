/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"errors"
)

func (m *GetHostnameForHWIDRequest) Validate() error {
	if m == nil {
		return errors.New("request cannot be nil")
	}
	if m.Hwid == "" {
		return errors.New("request params cannot be empty")
	}
	return nil
}

func (m *MapHWIDToHostnameRequest) Validate() error {
	if m == nil {
		return errors.New("request cannot be nil")
	}
	if m.HwidToHostname == nil {
		return errors.New("request params cannot be empty")
	}
	return nil
}

func (m *GetIMSIForSessionIDRequest) Validate() error {
	if m == nil {
		return errors.New("request cannot be nil")
	}
	if m.SessionID == "" {
		return errors.New("request params cannot be empty")
	}
	return nil
}

func (m *MapSessionIDToIMSIRequest) Validate() error {
	if m == nil {
		return errors.New("request cannot be nil")
	}
	if m.NetworkID == "" {
		return errors.New("network ID cannot be empty")
	}
	if m.SessionIDToIMSI == nil {
		return errors.New("request params cannot be empty")
	}
	return nil
}
