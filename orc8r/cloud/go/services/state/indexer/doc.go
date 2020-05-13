/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

// Package indexer provides tools to define, use, and update state indexers.
//
// Define
// 	- The state service handles indexing and updating based on indexer
//	  implementations. Each indexer need only
//		- Implement the Indexer interface
// 		- Be registered
//
// Index
// - The state service automatically sends received states to each indexer
//	  according to its subscriptions
//
// Reindex
// - Indexer implementations indicate they need to be reindex by incrementing
//   their version
// - From there, the state service handles reindexing coordination
package indexer
