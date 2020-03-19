/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

// Package directoryd provides an API for interacting with the
// directory lookup service, which manages UE location records.
//
// Primary state
// 	- reported directly from the relevant device/gateway
// 	- managed by the state service
//	- versioned
// Secondary state
// 	- derived, in the controller, from the primary state or other information
// 	- managed by the directoryd service (DirectoryLookupServer)
//	- non-versioned, with availability and correctness provided on a best-effort basis
package directoryd
