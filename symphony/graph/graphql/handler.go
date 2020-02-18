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

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/directive"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/resolver"
	"github.com/facebookincubator/symphony/graph/graphql/tracer"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/oc/ocgql"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"github.com/NYTimes/gziphandler"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/vektah/gqlparser/gqlerror"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
	"gocloud.dev/pubsub"
)

// HandlerConfig configures graphql handler.
type HandlerConfig struct {
	Logger      log.Logger
	Topic       *pubsub.Topic
	Subscribe   func(context.Context) (*pubsub.Subscription, error)
	Orc8rClient *http.Client
}

func init() {
	for _, v := range ocgql.DefaultServerViews {
		v.TagKeys = append(v.TagKeys,
			viewer.KeyTenant,
			viewer.KeyUser,
			viewer.KeyUserAgent,
		)
	}
}

// NewHandler creates a graphql http handler.
func NewHandler(cfg HandlerConfig) (http.Handler, func(), error) {
	rsv := resolver.New(
		resolver.Config{
			Logger:    cfg.Logger,
			Topic:     cfg.Topic,
			Subscribe: cfg.Subscribe,
		},
		resolver.WithOrc8rClient(
			cfg.Orc8rClient,
		),
	)

	if err := view.Register(ocgql.DefaultServerViews...); err != nil {
		return nil, nil, fmt.Errorf("registering views: %w", err)
	}
	closer := func() { view.Unregister(ocgql.DefaultServerViews...) }

	router := mux.NewRouter()
	router.Use(func(handler http.Handler) http.Handler {
		timeouter := http.TimeoutHandler(handler, 3*time.Minute, "request timed out")
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
							Directives: directive.New(cfg.Logger),
						},
					),
					ocgql.RequestMiddleware(),
					ocgql.ResolverMiddleware(),
					handler.Tracer(tracer.New()),
					handler.ErrorPresenter(errorPresenter(cfg.Logger)),
					handler.WebsocketUpgrader(websocket.Upgrader{
						CheckOrigin: func(*http.Request) bool {
							return true
						},
						ReadBufferSize:  1024,
						WriteBufferSize: 1024,
					}),
				),
			),
			"query",
		))

	return router, closer, nil
}

func errorPresenter(logger log.Logger) graphql.ErrorPresenterFunc {
	return func(ctx context.Context, err error) *gqlerror.Error {
		gqlerr := graphql.DefaultErrorPresenter(ctx, err)
		if strings.Contains(err.Error(), ent.ErrReadOnly.Error()) {
			gqlerr.Message = "Permission denied"
		} else if _, ok := err.(*gqlerror.Error); !ok {
			logger.For(ctx).Error("graphql internal error", zap.Error(err))
			gqlerr.Message = "Sorry, something went wrong"
		}
		return gqlerr
	}
}
