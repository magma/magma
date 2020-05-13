/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// Package clock provides a simple abstraction around the standard time package.
//	- Set and "freeze" the wall clock in test code, with provided wrappers for
//		- time.Now
//		- time.Since
//	- Skip sleeps in test code, with provided wrappers for
//		- time.Sleep
package clock
