// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"testing"

	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/assert"
)

func TestQueryMe(t *testing.T) {
	resolver := newTestResolver(t)
	defer resolver.drv.Close()
	c := newGraphClient(t, resolver)

	var rsp struct {
		Me struct {
			Tenant string
			Email  string
			User   struct {
				AuthID string
			}
		}
	}
	c.MustPost("query { me { tenant, email user { authID } } }", &rsp)
	assert.Equal(t, viewertest.DefaultViewer.Tenant, rsp.Me.Tenant)
	assert.Equal(t, viewertest.DefaultViewer.User, rsp.Me.Email)
	assert.Equal(t, viewertest.DefaultViewer.User, rsp.Me.User.AuthID)
}
