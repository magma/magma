/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package client provides a thin client for contacting the subscriberdb service.
// This can be used by apps to discover and contact the service, without knowing about
// the RPC implementation.
package subscriberdb

const EntityType = "subscriber"

const ServiceName = "SUBSCRIBERDB"
