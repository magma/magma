/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

func (m *NetworkCellularConfigs) GetEarfcndl() uint32 {
	switch {
	case m.Ran.FddConfig != nil:
		return m.Ran.FddConfig.Earfcndl
	case m.Ran.TddConfig != nil:
		return m.Ran.TddConfig.Earfcndl
	}
	return 0
}
