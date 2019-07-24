/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package middleware

import (
	"net/http"

	"fbc/lib/go/http/header"

	"github.com/google/uuid"
)

// RequestID returns a X-Request-ID middleware.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(header.XRequestID)
		if rid == "" {
			rid = uuid.New().String()
		}
		w.Header().Set(header.XRequestID, rid)
		next.ServeHTTP(w, r)
	})
}
