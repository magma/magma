// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
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
	_, err := client.ServiceEndpoint.Delete().Where(serviceendpoint.HasServiceWith(service.ID(s.ID))).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "deleting service endpoints")
	}
	_, err = client.Property.Delete().Where(property.HasServiceWith(service.ID(s.ID))).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "deleting service properties")
	}
	err = client.Service.DeleteOne(s).Exec(ctx)
	if err != nil {
		return errors.Wrap(err, "deleting service")
	}
	return nil
}
