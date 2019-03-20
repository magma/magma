/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// Package serde contains the definition of a SERializer-DEserializer concept.
// This package also includes a global registry of serdes for applications to
// delegate implementation-agnostic serialization and deserialization to.
// Serdes are one of the primary plugin interfaces exposed by orc8r to extend
// services with domain-specific data models and logic.
package serde
