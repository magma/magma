package authz_test

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent/reportfilter"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestReportFilterCanAlwaysBeCreated(t *testing.T) {
	c := viewertest.NewTestClient(t)
	ctx := viewertest.NewContext(context.Background(), c, viewertest.WithPermissions(authz.EmptyPermissions()))
	_, err := c.ReportFilter.Create().
		SetName("ReportFilter").
		SetEntity(reportfilter.EntityWORKORDER).
		Save(ctx)
	require.NoError(t, err)
}
