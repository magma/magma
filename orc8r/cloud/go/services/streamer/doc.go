/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

// Package streamer provides a logical stream for orc8r to push updates
// to gateways.
//
// Orc8r services can implement the StreamProvider servicer interface to
// provide a named stream. Streamer forwards requests for updates under a
// specific stream name to the appropriate remote orc8r service.
//
// E.g., consider a gateway requesting subscriber updates, from the
// "subscriber" stream. This takes the form
//
//		gateway -(a)-> streamer -(b)-> lte
//
// (a) GetUpdates("subscriber")
// (*) streamer: look up service name of provider for "subscriber" stream
// (b) GetUpdates("subscriber")
package streamer

const (
	ServiceName = "STREAMER"
)
