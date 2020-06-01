/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

// File keys.go contains the config keynames in the state service's YAML config file.

package config

const (
	// EnableAutomaticReindexing is a parameter name in the state service config.
	// When value is true, state service handles automatically reindex state indexers.
	// When value is false, reindexing must be handled by the provided CLI.
	EnableAutomaticReindexing = "enable_automatic_reindexing"
)
