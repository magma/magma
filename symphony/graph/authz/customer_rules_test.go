// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestCustomerCanAlwaysBeWritten(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c, viewertest.WithPermissions(authz.EmptyPermissions()))
	customer, err := c.Customer.Create().
		SetName("Customer").
		Save(ctx)
	require.NoError(t, err)
	err = c.Customer.UpdateOne(customer).
		SetName("NewCustomer").
		Exec(ctx)
	require.NoError(t, err)
	err = c.Customer.DeleteOne(customer).
		Exec(ctx)
	require.NoError(t, err)
}
