// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package actions

import (
	"context"
	"net/http"

	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/actions/executor"
	"github.com/facebookincubator/symphony/pkg/log"
	"go.uber.org/zap"
)

// Handler adds actions framework registry to incoming requests.
func Handler(next http.Handler, logger log.Logger, registry *executor.Registry) http.Handler {
	dataLoader := executor.BasicDataLoader{
		Rules: []core.Rule{},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		exc := &executor.Executor{
			Registry:   registry,
			DataLoader: dataLoader,
			OnError: func(ctx context.Context, err error) {
				logger.For(ctx).Error("error executing action", zap.Error(err))
			},
		}
		ctx = NewContext(ctx, exc)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
