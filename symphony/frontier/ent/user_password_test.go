// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ent

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/frontier/ent/user"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPassword(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()
	password := "foo-bar-baz"
	u, err := client.User.Create().
		SetEmail("test@example.com").
		SetPassword(MustHashPassword(password)).
		SetNetworks([]string{}).
		Save(ctx)
	require.NoError(t, err)
	assert.NotEqual(t, password, u.Password)

	t.Run("Validate", func(t *testing.T) {
		assert.NotEqual(t, password, u.Password)
		err := u.ValidatePassword(password)
		assert.NoError(t, err)
		err = u.ValidatePassword(password + ".")
		assert.Error(t, err)
		err = u.ValidatePassword(u.Password)
		assert.Error(t, err)
	})

	t.Run("Search", func(t *testing.T) {
		exist, err := client.User.Query().
			Where(user.Or(
				user.Password(MustHashPassword(password)),
				user.Password(password),
			)).
			Exist(ctx)
		require.NoError(t, err)
		assert.False(t, exist, "user searchable by hashed password")
	})
}
