// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewertest

import (
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/viewer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContext(t *testing.T) {
	ctx := NewContext(&ent.Client{}, WithTenant("facebook"))
	v := viewer.FromContext(ctx)
	require.NotNil(t, v)
	assert.Equal(t, "facebook", v.Tenant)
}
