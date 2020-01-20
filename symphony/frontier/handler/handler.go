// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
)

// Set is a Wire provider set that produces a handler from config.
var Set = wire.NewSet(
	NewHandler,
	wire.Struct(new(Config), "*"),
	wire.Bind(new(http.Handler), new(*mux.Router)),
)

type (
	// Config is the set of handler parameters.
	Config struct {
		ProxyTarget  ProxyTarget
		StaticTarget StaticTarget
		Logger       log.Logger
		AuthKey      []byte
	}

	// ProxyTarget wire dependency.
	ProxyTarget *url.URL
	// StaticTarget wire dependency.
	StaticTarget *url.URL
)

// NewHandler return a root http handler from config.
func NewHandler(cfg Config) *mux.Router {
	router := mux.NewRouter()
	router.Use(func(h http.Handler) http.Handler {
		csrf := nosurf.New(h)
		csrf.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg.Logger.For(r.Context()).
				Debug("failed csrf validation", zap.Error(nosurf.Reason(r)))
			w.WriteHeader(http.StatusBadRequest)
		}))
		return csrf
	})
	router.NotFoundHandler = newProxy(cfg.ProxyTarget, cfg.Logger)
	return router
}

func newProxy(target *url.URL, logger log.Logger) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &ochttp.Transport{}
	proxy.ErrorLog = zap.NewStdLog(logger.Background())
	return proxy
}

// Views for proxy metrics.
func Views() []*view.View {
	return []*view.View{
		ochttp.ClientRoundtripLatencyDistribution.
			WithName("http_proxy_request_duration_milliseconds"),
		ochttp.ClientCompletedCount.
			WithName("http_proxy_request_total"),
	}
}
