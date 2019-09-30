/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package counters

import "go.opencensus.io/tag"

var (
	// ListenerTag The name of the listener under which the operation occurred
	ListenerTag, _ = tag.NewKey("listener")

	// ModuleTag The module in which the operation occurred
	ModuleTag, _ = tag.NewKey("module")

	// FilterTag The filter in which the operation occurred
	FilterTag, _ = tag.NewKey("filter")

	// RadiusTypeTag The RADIUS message type
	RadiusTypeTag, _ = tag.NewKey("radius_type")

	// ErrorCodeTag code describing the error
	ErrorCodeTag, _ = tag.NewKey("error_code")

	// SessionIDTag code indicating the session id used for the operation
	SessionIDTag, _ = tag.NewKey("session_id")

	// StorageTag code describing the type of storage used for the operation
	StorageTag, _ = tag.NewKey("storage")
)
