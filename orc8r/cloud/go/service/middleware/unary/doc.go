/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// Package unary provides some default RPC interceptors and a wrapper around
// GRPC's unary interceptors called Interceptor. This package maintains a
// registry of interceptors to run on RPC requests.
package unary
