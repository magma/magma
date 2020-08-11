/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"magma/orc8r/lib/go/protos"
)

// FillJSONConfigs sets the build response's JsonConfigsByKey to match its
// ConfigsByKey.
// TODO(T71525030): remove this file
func (m *BuildResponse) FillJSONConfigs(err error) error {
	if err != nil {
		return err
	}
	for k, v := range m.ConfigsByKey {
		b, err := protos.MarshalJSON(v)
		if err != nil {
			return err
		}
		m.JsonConfigsByKey[k] = b
	}
	return nil
}
