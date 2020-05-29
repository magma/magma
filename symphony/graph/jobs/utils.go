// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

type ctxKey struct{}

type servicesContext struct {
	deleted int
}

func errorReturn(w http.ResponseWriter, msg string, log *zap.Logger, err error) {
	log.Warn(msg, zap.Error(err))
	if err == nil {
		http.Error(w, msg, http.StatusBadRequest)
	} else {
		http.Error(w, fmt.Sprintf("%s %q", msg, err), http.StatusBadRequest)
	}
}

func getServicesContext(ctx context.Context) *servicesContext {
	ld, _ := ctx.Value(ctxKey{}).(*servicesContext)
	return ld
}

func newServicesContext(parent context.Context) context.Context {
	return context.WithValue(parent, ctxKey{}, &servicesContext{
		deleted: 0,
	})
}

func deleteService(ctx context.Context, s *ent.Service) error {
	client := ent.FromContext(ctx)
	ids, err := s.QueryEndpoints().IDs(ctx)
	if err != nil {
		return fmt.Errorf("query service endpoints of service: %q, %w", s.ID, err)
	}
	for _, id := range ids {
		if err := client.ServiceEndpoint.DeleteOneID(id).Exec(ctx); err != nil {
			return errors.Wrap(err, "deleting service endpoint")
		}
	}
	ids, err = s.QueryProperties().IDs(ctx)
	if err != nil {
		return fmt.Errorf("query properties of service: %q, %w", s.ID, err)
	}
	for _, id := range ids {
		if err := client.Property.DeleteOneID(id).Exec(ctx); err != nil {
			return errors.Wrap(err, "deleting property")
		}
	}
	err = client.Service.DeleteOne(s).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "deleting service")
	}
	return nil
}
