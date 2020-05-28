// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/enttest"
	"github.com/facebookincubator/symphony/pkg/ent/migrate"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/facebookincubator/symphony/pkg/testdb"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func newJobsTestResolver(t *testing.T) *TestJobsResolver {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	return newResolver(t, sql.OpenDB(name, db))
}

func newResolver(t *testing.T, drv dialect.Driver) *TestJobsResolver {
	client := enttest.NewClient(t,
		enttest.WithOptions(ent.Driver(drv)),
		enttest.WithMigrateOptions(migrate.WithGlobalUniqueID(true)),
	)
	r := resolver.New(resolver.Config{
		Logger:     logtest.NewTestLogger(t),
		Subscriber: event.NewNopSubscriber(),
	})
	return &TestJobsResolver{
		drv:    drv,
		client: client,
		jobsRunner: jobs{
			logger: logtest.NewTestLogger(t),
			r:      r,
		},
	}
}

func syncServicesRequest(t *testing.T, r *TestJobsResolver) *http.Response {
	h, _ := NewHandler(
		Config{
			Logger:     logtest.NewTestLogger(t),
			Subscriber: event.NewNopSubscriber(),
		},
	)

	th := viewertest.TestHandler(t, h, r.client)
	server := httptest.NewServer(th)
	defer server.Close()
	url := server.URL + "/sync_services"
	req, err := http.NewRequest(http.MethodGet, url, ioutil.NopCloser(new(bytes.Buffer)))
	require.Nil(t, err)

	viewertest.SetDefaultViewerHeaders(req)
	req.Header.Set("Content-Length", "100000")

	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
	return resp
}
