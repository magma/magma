// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestAddNewService(t *testing.T) {
	r := newJobsTestResolver(t)
	defer r.drv.Close()

	ctx := newServicesContext(viewertest.NewContext(context.Background(), r.client))
	eData := prepareEquipmentData(ctx, *r, "A")
	prepareServiceTypeData(ctx, *r, eData)
	prepareLinksData(ctx, *r, eData)

	sCount := r.client.Service.Query().CountX(ctx)
	require.Zero(t, sCount)
	syncServicesRequest(t, r)
	sCount = r.client.Service.Query().CountX(ctx)
	require.Equal(t, 1, sCount)
}
