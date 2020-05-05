/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// Package clock provides a simple abstraction around the standard time package.
//	- time.Now() -- ability to set and "feeze" the wall clock in test code
//	- time.Sleep() -- ability to skip sleeps
package clock
