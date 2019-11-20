// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package recovery

import (
	"context"
	"net/http"
)

// HandlerFunc is a function that recovers from panic by returning an `error`.
type HandlerFunc func(context.Context, interface{}) error

// Handler is an http.Handler wrapper for panic recovery.
type Handler struct {
	// Handler is the handler used to handle the incoming request.
	Handler http.Handler
	// HandlerFunc is a function that recovers from panic.
	HandlerFunc HandlerFunc
}

// ServeHTTP implements http.Handler interface.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if p := recover(); p != nil {
			var msg string
			if err := h.recoverFrom(r.Context(), p); err != nil {
				msg = err.Error()
			}
			http.Error(w, msg, http.StatusInternalServerError)
		}
	}()
	h.Handler.ServeHTTP(w, r)
}

func (h *Handler) recoverFrom(ctx context.Context, p interface{}) error {
	if h.HandlerFunc != nil {
		return h.HandlerFunc(ctx, p)
	}
	return nil
}
