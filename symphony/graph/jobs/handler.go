// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"net/http"

	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/gorilla/mux"
)

type (
	// Config configures jobs handler.
	Config struct {
		Logger     log.Logger
		Subscriber event.Subscriber
	}

	jobs struct {
		logger log.Logger
		r      generated.ResolverRoot
	}
)

// NewHandler creates a upload http handler.
func NewHandler(cfg Config) (http.Handler, error) {
	r := resolver.New(
		resolver.Config{
			Logger:     cfg.Logger,
			Subscriber: cfg.Subscriber,
		},
		resolver.WithTransaction(false),
	)
	u := &jobs{cfg.Logger, r}
	router := mux.NewRouter()

	routes := []struct {
		name    string
		handler http.HandlerFunc
	}{
		{"sync_services", u.syncServices},
	}
	for _, route := range routes {
		router.Path("/" + route.name).
			Methods(http.MethodPost).
			HandlerFunc(route.handler).
			Name(route.name)
	}
	return router, nil
}
