/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package streamer provides streamer client Go implementation for golang based gateways
package streamer

import (
	"magma/orc8r/lib/go/protos"
)

// Listener interface defines Stream Listener which will become
// the receiver of streaming updates for a registered stream
// Each received update will be unmarshalled into the Listener's update data type determined by
// the actual type returned by Listener's New() receiver method
type Listener interface {
	// GetName() returns name of the stream, the listener is getting updates on
	GetName() string
	// ReportError is going to be called by the streamer on every error.
	// If ReportError() will return nil, streamer will try to continue streaming
	// If ReportError() will return error != nil - streaming on the stream will be terminated
	ReportError(e error) error
	// Update will be called for every new update received from the stream
	// u is guaranteed to be of a type returned by New(), so - myUpdate := u.(MyDataType) should never panic
	// Update() returns bool indicating whether to continue streaming:
	//   true - continue streaming; false - stop streaming
	Update(u *protos.DataUpdateBatch) bool
}
