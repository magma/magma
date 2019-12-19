// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
	"github.com/facebookincubator/symphony/pkg/actions/trigger/magmaalert"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	registry := executor.NewRegistry()
	registry.MustRegisterTrigger(magmaalert.New())

	h := Handler(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			exc := FromContext(r.Context())
			assert.NotNil(t, exc)
			trigger, err := exc.TriggerForID(core.MagmaAlertTriggerID)
			assert.NoError(t, err)
			assert.NotNil(t, trigger)
			_, _ = io.WriteString(w, "success")
		}),
		nil,
		registry,
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "success", rec.Body.String())
}
