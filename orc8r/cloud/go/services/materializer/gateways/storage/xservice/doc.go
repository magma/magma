/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package xservice (cross-service) contains a *read-only* storage
// implementation which is backed by the client APIs of the various services
// that the materializer gateway views are computed from. Use this
// implementation when the streamconsumers are lagging or showing symptoms of
// data inconsistencies so that user-facing data remains consistent while
// debugging.
package xservice
