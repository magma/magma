// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xserver

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViews(t *testing.T) {
	r := regexp.MustCompile(
		`http_(request|response)_.*(total|bytes|seconds)`,
	)
	for _, v := range DefaultViews() {
		assert.Regexp(t, r, v.Name)
	}
}
