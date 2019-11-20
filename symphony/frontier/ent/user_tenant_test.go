// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserTenant(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()
	tenant, err := client.Tenant.Create().
		SetName("testing").
		SetDomains([]string{}).
		SetNetworks([]string{}).
		Save(ctx)
	require.NoError(t, err)

	user, err := client.User.Create().
		SetTenant(tenant.Name).
		SetEmail("test@example.com").
		SetPassword("random-password").
		SetNetworks([]string{}).
		Save(ctx)
	require.NoError(t, err)

	tenant, err = user.QueryTenant().Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, "testing", tenant.Name)
}
