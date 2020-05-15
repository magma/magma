// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"

	"github.com/facebookincubator/symphony/graph/ent/service"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"

	"github.com/facebookincubator/symphony/graph/ent/servicetype"

	"go.uber.org/zap"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
)

// syncServices job syncs the services according to changes
func (m *jobs) garbageCollector(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := m.collectProperties(ctx); err != nil {
		m.logger.For(ctx).Error("collect properties", zap.Error(err))
	}
	if err := m.collectServices(ctx); err != nil {
		m.logger.For(ctx).Error("collect services", zap.Error(err))
	}
	w.WriteHeader(http.StatusOK)
}

func (m *jobs) collectProperties(ctx context.Context) error {
	client := ent.FromContext(ctx)
	m.logger.For(ctx).Info("running properties garbage collect")
	propertyTypes, err := client.PropertyType.Query().
		Where(propertytype.Deleted(true)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query properties: %w", err)
	}
	for _, pType := range propertyTypes {
		m.logger.For(ctx).Info("deleting property type",
			zap.Int("id", pType.ID),
			zap.String("name", pType.Name))
		count, err := client.Property.Delete().
			Where(property.HasTypeWith(propertytype.ID(pType.ID))).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete properties of type id: %d, %w", pType.ID, err)
		}
		m.logger.For(ctx).Info("deleted properties",
			zap.Int("id", pType.ID),
			zap.Int("count", count))
		err = client.PropertyType.DeleteOne(pType).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete property type: %d, %w", pType.ID, err)
		}
		m.logger.For(ctx).Info("deleted property type",
			zap.Int("id", pType.ID),
			zap.String("name", pType.Name))
	}
	return nil
}

func (m *jobs) collectServices(ctx context.Context) error {
	client := ent.FromContext(ctx)
	m.logger.For(ctx).Info("running services garbage collect")
	serviceTypes, err := client.ServiceType.Query().
		Where(servicetype.IsDeleted(true)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("query services: %w", err)
	}
	for _, sType := range serviceTypes {
		m.logger.For(ctx).Info("deleting service type",
			zap.Int("id", sType.ID),
			zap.String("name", sType.Name))

		count, err := client.Property.Delete().
			Where(property.HasServiceWith(service.HasTypeWith(servicetype.ID(sType.ID)))).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete properties of services. type id: %d, %w", sType.ID, err)
		}
		m.logger.For(ctx).Info("deleted properties",
			zap.Int("id", sType.ID),
			zap.Int("count", count))

		count, err = client.ServiceEndpoint.Delete().
			Where(serviceendpoint.HasServiceWith(service.HasTypeWith(servicetype.ID(sType.ID)))).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("deleting service endpoints for service type id: %d, %w", sType.ID, err)
		}

		m.logger.For(ctx).Info("deleted endpoints",
			zap.Int("id", sType.ID),
			zap.Int("count", count))

		count, err = client.Service.Delete().
			Where(service.HasTypeWith(servicetype.ID(sType.ID))).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete services of type id: %d, %w", sType.ID, err)
		}
		m.logger.For(ctx).Info("deleted services",
			zap.Int("id", sType.ID),
			zap.Int("count", count))

		count, err = client.PropertyType.Delete().
			Where(propertytype.HasServiceTypeWith(servicetype.ID(sType.ID))).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete property types of service type id: %d, %w", sType.ID, err)
		}
		m.logger.For(ctx).Info("deleted property types",
			zap.Int("id", sType.ID),
			zap.Int("count", count))

		count, err = client.ServiceEndpointDefinition.Delete().
			Where(serviceendpointdefinition.HasServiceTypeWith(servicetype.ID(sType.ID))).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("deleting service endpoints for service type id: %d, %w", sType.ID, err)
		}
		m.logger.For(ctx).Info("deleted endpoint definitions",
			zap.Int("id", sType.ID),
			zap.Int("count", count))

		err = client.ServiceType.DeleteOne(sType).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("delete service type: %d, %w", sType.ID, err)
		}
		m.logger.For(ctx).Info("deleted service type",
			zap.Int("id", sType.ID),
			zap.String("name", sType.Name))
	}
	return nil
}
