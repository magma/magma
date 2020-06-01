// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/zap"
	"gocloud.dev/blob"
)

// Set is a Wire provider set that produces a handler from config.
var Set = wire.NewSet(
	New,
	wire.Struct(new(Config), "*"),
	wire.Bind(new(http.Handler), new(*Handler)),
)

// Handler implements signer endpoints.
type Handler struct {
	http.Handler
	logger log.Logger
	bucket *blob.Bucket
}

// Config is the set of handler parameters.
type Config struct {
	Logger log.Logger
	Bucket *blob.Bucket
}

// New creates a new sign handler from config.
func New(cfg Config) *Handler {
	h := &Handler{
		logger: cfg.Logger,
		bucket: cfg.Bucket,
	}

	router := mux.NewRouter()
	router.Path("/get").
		Methods(http.MethodGet).
		Queries("key", "{key}").
		Handler(ochttp.WithRouteTag(
			http.HandlerFunc(h.get), "get"),
		)
	router.Path("/put").
		Methods(http.MethodGet).
		Handler(ochttp.WithRouteTag(
			http.HandlerFunc(h.put), "put",
		))
	router.Path("/delete").
		Queries("key", "{key}").
		Methods(http.MethodDelete).
		Handler(ochttp.WithRouteTag(
			http.HandlerFunc(h.delete), "delete",
		))
	router.Path("/download").
		Methods(http.MethodGet).
		Queries("key", "{key}", "fileName", "{filename}").
		Handler(ochttp.WithRouteTag(
			http.HandlerFunc(h.download), "download",
		))
	h.Handler = router
	return h
}

func (h *Handler) key(r *http.Request, key string) string {
	isGlobal := r.Header.Get("Is-Global")
	if key != "" && isGlobal != "True" {
		if ns := r.Header.Get("x-auth-organization"); ns != "" {
			key = ns + "/" + key
		}
	}
	return key
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.For(ctx)
	key := h.key(r, mux.Vars(r)["key"])
	if key == "" {
		logger.Error("cannot resolve object key")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	u, err := h.bucket.SignedURL(ctx, key, nil)
	if err != nil {
		logger.Error("cannot sign get object url", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	logger.Debug("signed get url", zap.String("key", key))
	http.Redirect(w, r, u, http.StatusSeeOther)
}

func (h *Handler) put(w http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		logger = h.logger.For(ctx)
		rsp    struct {
			URL string `json:"URL"`
			Key string `json:"key"`
		}
		err error
	)
	rsp.Key = uuid.New().String()
	key := h.key(r, rsp.Key)
	if rsp.URL, err = h.bucket.SignedURL(ctx, key,
		&blob.SignedURLOptions{Method: http.MethodPut},
	); err != nil {
		logger.Error("cannot sign put object url", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&rsp); err != nil {
		logger.Error("cannot write put object response", zap.Error(err))
		return
	}
	logger.Debug("signed put url", zap.String("key", key))
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.For(ctx)
	key := h.key(r, mux.Vars(r)["key"])
	if key == "" {
		logger.Error("cannot resolve object key")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	u, err := h.bucket.SignedURL(ctx, key,
		&blob.SignedURLOptions{Method: http.MethodDelete},
	)
	if err != nil {
		logger.Error("cannot sign delete object url", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	logger.Debug("signed delete url", zap.String("key", key))
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (h *Handler) download(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	logger := h.logger.For(ctx)
	key := h.key(r, vars["key"])
	if key == "" {
		logger.Error("cannot resolve object key")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	filename := vars["filename"]
	if filename == "" {
		logger.Error("cannot resolve object filename")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	rawurl, err := h.bucket.SignedURL(ctx, key, nil)
	if err != nil {
		logger.Error("cannot sign get object url", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	u, _ := url.Parse(rawurl)
	q := u.Query()
	q.Set("response-content-disposition", "attachment; filename="+filename)
	u.RawQuery = q.Encode()
	logger.Debug("signed download url",
		zap.String("key", key),
		zap.String("filename", filename),
	)
	http.Redirect(w, r, u.String(), http.StatusSeeOther)
}
