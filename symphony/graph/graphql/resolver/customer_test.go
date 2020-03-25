// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"
)

func TestAddDeleteAndSearchCustomers(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	qr, mr := r.Query(), r.Mutation()
	ctx := viewertest.NewContext(r.client)

	_, err := mr.AddCustomer(ctx, models.AddCustomerInput{Name: "Donald Duck", ExternalID: pointer.ToString("S43493")})
	require.NoError(t, err)

	c, err := mr.AddCustomer(ctx, models.AddCustomerInput{Name: "Dafi Duck"})
	require.NoError(t, err)

	limit := 10
	res1, err := qr.CustomerSearch(ctx, &limit)
	require.NoError(t, err)
	require.Len(t, res1, 2)

	_, err = mr.RemoveCustomer(ctx, c.ID)
	require.NoError(t, err)

	res2, err := qr.CustomerSearch(ctx, &limit)
	require.NoError(t, err)
	require.Len(t, res2, 1)
}
