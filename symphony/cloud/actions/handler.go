// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"net/http"

	"github.com/facebookincubator/symphony/cloud/actions/action/magmarebootnode"
	"github.com/facebookincubator/symphony/cloud/actions/core"
	"github.com/facebookincubator/symphony/cloud/actions/executor"
	"github.com/facebookincubator/symphony/cloud/actions/trigger/magmaalert"
	"github.com/facebookincubator/symphony/cloud/log"
	"go.uber.org/zap"
)

// MainRegistry is a registry that contains all actions and triggers
func MainRegistry() executor.Registry {

	registry := executor.NewRegistry()

	registry.MustRegisterAction(magmarebootnode.New())
	registry.MustRegisterTrigger(magmaalert.New())

	return registry
}

// Handler adds actions framework registry to incoming requests.
func Handler(next http.Handler, logger log.Logger) http.Handler {

	dataLoader := executor.BasicDataLoader{
		Rules: []core.Rule{},
	}

	registry := MainRegistry()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		exc := &executor.Executor{
			Context:    ctx,
			Registry:   registry,
			DataLoader: dataLoader,
			OnError: func(err error) {
				logger.For(ctx).Error("error executing action", zap.Error(err))
			},
		}
		ctx = NewContext(ctx, exc)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
