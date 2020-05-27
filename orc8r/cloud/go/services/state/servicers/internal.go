/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"magma/orc8r/cloud/go/services/state"
)

// StateServiceInternal provides a cross-network DAO for local usage
// by the state service.
type StateServiceInternal interface {
	// GetAllIDs returns all IDs known to the state service, keyed
	// by network ID.
	GetAllIDs() (state.IDsByNetwork, error)
}
