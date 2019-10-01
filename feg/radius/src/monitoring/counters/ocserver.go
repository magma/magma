/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package counters

var (
	// ServerInit Counterset for server initialization
	ServerInit = NewOperation("server_init")

	// FilterInit counterset for module initialization
	FilterInit = NewOperation("filter_init", FilterTag)

	// ListenerInit counterset for listener initialization
	ListenerInit = NewOperation("listener_init", ListenerTag)

	// ModuleInit counterset for module initialization
	ModuleInit = NewOperation("module_init", ListenerTag, ModuleTag)

	// DedupPacket RADIUS dedup logic counter
	DedupPacket = NewOperation("radius_dedup", ListenerTag)
)
