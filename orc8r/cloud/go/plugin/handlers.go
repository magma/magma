/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package plugin

import "magma/orc8r/cloud/go/obsidian"

// FlattenHandlerLists turns a variadic list of obsidian handlers into a
// single flattened list of handlers. This is typically used to merge handlers
// from different services into a single collection to return in an impl
// of OrchestratorPlugin.
func FlattenHandlerLists(handlersIn ...[]obsidian.Handler) []obsidian.Handler {
	var ret []obsidian.Handler
	for _, h := range handlersIn {
		ret = append(ret, h...)
	}
	return ret
}
