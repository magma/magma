/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package oc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opencensus.io/stats/view"
)

func TestViewerRegistration(t *testing.T) {
	want := Views{&view.View{Name: "test"}}
	err := RegisterViewer("test", want)
	assert.NoError(t, err)
	got := GetViewer("test")
	assert.Equal(t, want, got)
	assert.Panics(t, func() { MustRegisterViewer("test", nil) })
}
