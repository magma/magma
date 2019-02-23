/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package materializer

// StreamProcessor describes an element of a stream topology which can be
// started and stopped.
type StreamProcessor interface {
	Run() error
	Stop()
}

// Application is a collection of StreamProcessors with a name -
// it represents a streaming topology. Constructing the stream topology is left
// to the Application author - wire up your stream processors as you see fit
// with intermediary state stores or Kakfa topics.
type Application struct {
	Name       string
	Processors []StreamProcessor
}
