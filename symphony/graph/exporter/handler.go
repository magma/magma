// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"net/url"

	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// Exporter encapsulates accessing db and logging
type exporter struct {
	log log.Logger
	rower
}

type rower interface {
	rows(context.Context, *url.URL) ([][]string, error)
}

// ServerHTTP handles requests to returns an export CSV file
func (m *exporter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	filename := "export"
	rout := mux.CurrentRoute(r)
	if rout != nil {
		filename = rout.GetName()
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename+".csv")
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Transfer-Encoding", "chunked")

	writer := csv.NewWriter(w)
	log := m.log.For(ctx)

	rows, err := m.rows(ctx, r.URL)
	if err != nil {
		log.Error("error in export", zap.Error(err))
		http.Error(w, fmt.Sprintf("%q: error in export", err), http.StatusInternalServerError)
	}

	err = writer.WriteAll(rows)
	if err != nil {
		http.Error(w, fmt.Sprintf("%q: error in writing file", err), http.StatusInternalServerError)
	}
}

// NewHandler creates a upload http handler.
func NewHandler(log log.Logger) (http.Handler, error) {
	router := mux.NewRouter()
	routes := []struct {
		name    string
		handler exporter
	}{
		{"equipment", exporter{log, equipmentRower{log}}},
		{"ports", exporter{log, portsRower{log}}},
		{"work_orders", exporter{log, woRower{log}}},
		{"links", exporter{log, linksRower{log}}},
		{"locations", exporter{log, locationsRower{log}}},
		{"services", exporter{log, servicesRower{log}}},
	}

	for _, route := range routes {
		route := route
		router.Path("/" + route.name).
			Methods(http.MethodGet).
			Handler(&route.handler).
			Name(route.name)
	}
	return router, nil
}
