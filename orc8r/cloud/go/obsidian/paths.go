/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package obsidian

const (
	UrlSep = "/"

	MagmaNetworksUrlPart  = "networks"
	MagmaOperatorsUrlPart = "operators"

	// "/magma"
	RestRoot = UrlSep + "magma"
	// "/magma/networks"
	NetworksRoot = RestRoot + UrlSep + MagmaNetworksUrlPart
	// "/magma/operators"
	OperatorsRoot = RestRoot + UrlSep + MagmaOperatorsUrlPart

	// Supported API versions
	V0 = ""
	V1 = "v1"
	// Note the trailing slash (this is actually important for apidocs to render properly)
	V1Root = RestRoot + UrlSep + V1 + UrlSep
)
