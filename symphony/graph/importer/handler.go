// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"net/http"

	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp"
)

type (
	// Config configures import handler.
	Config struct {
		Logger     log.Logger
		Emitter    event.Emitter
		Subscriber event.Subscriber
	}

	importer struct {
		logger log.Logger
		r      generated.ResolverRoot
	}
)

// NewHandler creates a upload http handler.
func NewHandler(cfg Config) (http.Handler, error) {
	r := resolver.New(
		resolver.Config{
			Logger:     cfg.Logger,
			Emitter:    cfg.Emitter,
			Subscriber: cfg.Subscriber,
		},
		resolver.WithTransaction(false),
	)
	u := &importer{cfg.Logger, r}

	router := mux.NewRouter()
	router.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := newImportContext(r.Context())
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if route := mux.CurrentRoute(r); route != nil {
					if name := route.GetName(); name != "" {
						ochttp.SetRoute(r.Context(), "import_"+name)
					}
				}
				next.ServeHTTP(w, r)
			})
		},
	)

	routes := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{"location", u.processLocationsCSV},
		{"equipment", u.processEquipmentCSV},
		{"port_def", u.processPortDefinitionsCSV},
		{"port_connect", u.processPortConnectionCSV},
		{"position_def", u.processPositionDefinitionsCSV},
		{"export_equipment", u.processExportedEquipment},
		{"export_ports", u.processExportedPorts},
		{"export_links", u.processExportedLinks},
		{"export_locations", u.processExportedLocation},
		{"export_service", u.processExportedService},
	}
	for _, route := range routes {
		router.Path("/" + route.name).
			Methods(http.MethodPost).
			HandlerFunc(route.handler).
			Name(route.name)
	}
	return router, nil
}
