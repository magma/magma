// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sqltx

import (
	"context"
	"database/sql"
	"fmt"

	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptors that adds *sql.Tx to the context.
func UnaryServerInterceptor(db *sql.DB) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (rsp interface{}, err error) {
		var tx *sql.Tx
		if tx, err = db.BeginTx(ctx, nil); err != nil {
			return nil, fmt.Errorf("beginning transaction: %w", err)
		}
		defer func() {
			if r := recover(); r != nil {
				_ = tx.Rollback()
				panic(r)
			}
			if err != nil {
				_ = tx.Rollback()
				return
			}
			if err = tx.Commit(); err != nil {
				err = fmt.Errorf("committing transaction: %w", err)
			}
		}()
		return handler(NewContext(ctx, tx), req)
	}
}

type contextKey struct{}

// NewContext returns a new context with the given tx attached.
func NewContext(parent context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(parent, contextKey{}, tx)
}

// FromContext returns the transaction stored in a context, or nil if there isn't one.
func FromContext(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(contextKey{}).(*sql.Tx)
	return tx
}
