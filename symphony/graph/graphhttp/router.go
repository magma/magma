// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphhttp

import (
	"net/http"

	"github.com/facebookincubator/symphony/graph/exporter"
	"github.com/facebookincubator/symphony/graph/graphql"
	"github.com/facebookincubator/symphony/graph/importer"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/actions"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func newRouter(tenancy viewer.Tenancy, logger log.Logger, orc8rClient *http.Client, actionsRegistry *executor.Registry) (*mux.Router, error) {
	router := mux.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		return viewer.TenancyHandler(h, tenancy)
	})
	router.Use(func(h http.Handler) http.Handler {
		return actions.Handler(h, logger, actionsRegistry)
	})
	importHandler, err := importer.NewHandler(logger)
	if err != nil {
		return nil, errors.WithMessage(err, "creating import handler")
	}
	router.PathPrefix("/import/").
		Handler(http.StripPrefix("/import", importHandler)).
		Name("import")

	exportHandler, err := exporter.NewHandler(logger)
	if err != nil {
		return nil, errors.WithMessage(err, "creating export handler")
	}

	router.PathPrefix("/export/").
		Handler(http.StripPrefix("/export", exportHandler)).
		Name("export")

	handler, err := graphql.NewHandler(logger, orc8rClient)
	if err != nil {
		return nil, errors.WithMessage(err, "creating graphql handler")
	}
	router.PathPrefix("/").
		Handler(handler).
		Name("root")
	return router, nil
}
