// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewertest

import (
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/viewer"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(t *testing.T) {
	want := &viewer.Viewer{
		Tenant: "facebook",
		User:   "fbuser@fb.com",
	}
	ctx := NewContext(&ent.Client{}, WithViewer(want))
	got := viewer.FromContext(ctx)
	assert.Equal(t, want, got)
}
