/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package unary

// Identity decorator bypass list is a map of RPC methods which are allowed to
// bypass Identity verification checks. For now, the only service in this category
// is Bootstrapper

var identityDecoratorBypassList = map[string]struct{}{
	// These 2 entries are here for back-compat. This may not actually be
	// necessary, as the UnaryServerInfo.FullMethod field should indicate the
	// magma.orc8r.* values even if they are on the legacy descriptor.
	"/magma.Bootstrapper/GetChallenge": {},
	"/magma.Bootstrapper/RequestSign":  {},

	"/magma.orc8r.Bootstrapper/GetChallenge": {},
	"/magma.orc8r.Bootstrapper/RequestSign":  {},
}
