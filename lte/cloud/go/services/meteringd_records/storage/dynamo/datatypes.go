/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo

// ===========================================================================
//
// HEY! Are you planning to make changes to this file? READ THIS FIRST!
//
// Changes to existing structs will immediately be reflected in new data
// persisted to dynamoDB. As such, you should treat all existing structs as
// immutable types.
//
// If you need to make a data schema change, create a new struct and bump
// the schema version in storage.go, then extend the Encoder and Decoder
// implementations to handle the different persistence formats based on the
// schema version you read back. Remember DynamoDB is schemaless, so it will
// happily eat whatever malformed data you shove into it as long as the
// key conditions are met, for better or for worse.
//
// A better way to do this would be a map between schema version (int) and
// a type. Something to keep on the list when that urge to dig into the
// reflection library hits.
//
// ===========================================================================

// For composite keys consisting of 2 or more attributes concatenated together,
// use this delimiter to divide each component.
const CompositeKeyDelimiter = "|"

type flowRecord struct {
	Id              string
	Sid             string
	NetworkId       string
	SubNetworkId    string // Internal: subscriber and network ID concatenated with delimiter "|"
	GatewayId       string
	BytesTx         uint64
	BytesRx         uint64
	PktsTx          uint64
	PktsRx          uint64
	StartTime       int64
	LastUpdatedTime int64
	SchemaVersion   int
}
