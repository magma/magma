// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"net/http"
	"net/http/httputil"

	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/zap"
)

type importer struct {
	log log.Logger
	r   generated.ResolverRoot
}

// newImporter is a constructor for an importer
func newImporter(log log.Logger, r generated.ResolverRoot) *importer {
	return &importer{log, r}
}

// NewHandler creates a upload http handler.
func NewHandler(log log.Logger) (http.Handler, error) {
	r, err := resolver.New(log, resolver.WithTransaction(false))
	if err != nil {
		return nil, errors.WithMessage(err, "creating resolver")
	}
	u := newImporter(log, r)

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
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if content, err := httputil.DumpRequest(r, false); err == nil {
					log.For(r.Context()).Debug(
						"http request dump",
						zap.ByteString("content", content),
					)
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
		{"ftth", u.ProcessFTTHCSV},
		{"xwfAps", u.ProcessXwfApsCSV},
		{"xwf1", u.ProcessXwf1CSV},
		{"rural_ran", u.ProcessRuralRanCSV},
		{"rural_transport", u.ProcessRuralTransportCSV},
		{"rural_legacy_locations", u.ProcessRuralLegacyLocationsCSV},
		{"rural_locations", u.ProcessRuralLocationsCSV},
		{"export_equipment", u.processExportedEquipment},
		{"export_ports", u.processExportedPorts},
		{"export_links", u.processExportedLinks},
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
