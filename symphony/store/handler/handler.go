// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handler

import (
	"net/http"

	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/store/sign"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"go.opencensus.io/plugin/ochttp"
	"go.uber.org/zap"
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
	logger   log.Logger
	renderer *render.Render
	signer   sign.Signer
}

// Config is the set of handler parameters.
type Config struct {
	Logger log.Logger
	Signer sign.Signer
}

// New creates a new sign handler from config.
func New(cfg Config) *Handler {
	h := &Handler{
		logger: cfg.Logger,
		renderer: render.New(
			render.Options{StreamingJSON: true},
		),
		signer: cfg.Signer,
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
	key := h.key(r, mux.Vars(r)["key"])
	if key == "" {
		h.logger.For(r.Context()).Error("cannot resolve object key")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	url, err := h.signer.Sign(r.Context(), sign.GetObject, key, "")
	if err != nil {
		h.logger.For(r.Context()).Error("cannot sign get object operation", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handler) put(w http.ResponseWriter, r *http.Request) {
	oid := uuid.New().String()
	key := h.key(r, oid)
	url, err := h.signer.Sign(r.Context(), sign.PutObject, key, "")
	if err != nil {
		h.logger.For(r.Context()).Error("cannot sign put object operation", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	if err = h.renderer.JSON(w, http.StatusOK, map[string]string{"URL": url, "key": oid}); err != nil {
		h.logger.For(r.Context()).Error("cannot write put object response", zap.Error(err))
	}
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	key := h.key(r, mux.Vars(r)["key"])
	if key == "" {
		h.logger.For(r.Context()).Error("cannot resolve object key")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	url, err := h.signer.Sign(r.Context(), sign.DeleteObject, key, "")
	if err != nil {
		h.logger.For(r.Context()).Error("cannot sign delete object operation", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *Handler) download(w http.ResponseWriter, r *http.Request) {
	key := h.key(r, mux.Vars(r)["key"])
	if key == "" {
		h.logger.For(r.Context()).Error("cannot resolve object key")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	filename := mux.Vars(r)["filename"]
	if filename == "" {
		h.logger.For(r.Context()).Error("cannot resolve object filename")
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	url, err := h.signer.Sign(r.Context(), sign.DownloadObject, key, filename)
	if err != nil {
		h.logger.For(r.Context()).Error("cannot sign get object operation", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusSeeOther)
}
