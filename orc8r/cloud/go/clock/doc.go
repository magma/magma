/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// Package clock provides a simple abstraction around Golang's time.Now().
// This package exposes the ability to set and "feeze" the wall clock in
// test code.
package clock
