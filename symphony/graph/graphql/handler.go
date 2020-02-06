// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphql

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/directive"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/graph/graphql/tracer"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/gorilla/websocket"

	gqlprometheus "github.com/99designs/gqlgen-contrib/prometheus"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/mux"
	"github.com/vektah/gqlparser/gqlerror"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/zap"
)

func init() { gqlprometheus.Register() }

// NewHandler creates a graphql http handler.
func NewHandler(logger log.Logger, orc8rClient *http.Client) (http.Handler, error) {
	var opts []resolver.ResolveOption
	opts = append(opts, resolver.WithOrc8rClient(orc8rClient))
	rsv, err := resolver.New(logger, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating resolver: %w", err)
	}

	router := mux.NewRouter()
	router.Use(func(handler http.Handler) http.Handler {
		timeouter := http.TimeoutHandler(handler, 30*time.Second, "request timed out")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := timeouter
			if websocket.IsWebSocketUpgrade(r) {
				h = handler
			}
			h.ServeHTTP(w, r)
		})
	})

	router.Path("/graphiql").
		MatcherFunc(func(*http.Request, *mux.RouteMatch) bool {
			_, ok := os.LookupEnv("GQL_DEBUG")
			return ok
		}).
		Handler(ochttp.WithRouteTag(
			handler.Playground("GraphIQL", "/graph/query"),
			"graphiql",
		))
	router.Path("/query").
		Handler(ochttp.WithRouteTag(
			gziphandler.GzipHandler(
				handler.GraphQL(
					generated.NewExecutableSchema(
						generated.Config{
							Resolvers:  rsv,
							Directives: directive.New(logger),
						},
					),
					handler.RequestMiddleware(gqlprometheus.RequestMiddleware()),
					handler.ResolverMiddleware(gqlprometheus.ResolverMiddleware()),
					handler.Tracer(tracer.New()),
					handler.ErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
						gqlerr := graphql.DefaultErrorPresenter(ctx, err)
						if strings.Contains(err.Error(), ent.ErrReadOnly.Error()) {
							gqlerr.Message = "Permission denied"
						} else if _, ok := err.(*gqlerror.Error); !ok {
							logger.For(ctx).Error("graphql internal error", zap.Error(err))
							gqlerr.Message = "Sorry, something went wrong"
						}
						return gqlerr
					}),
				),
			),
			"query",
		))

	return router, nil
}
