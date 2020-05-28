// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/facebookincubator/symphony/pkg/ent/servicetype"

	"go.uber.org/zap"

	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
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
		return fmt.Errorf("query property type: %w", err)
	}
	for _, pType := range propertyTypes {
		if err := m.deletePropertyType(ctx, client, pType); err != nil {
			return err
		}
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

		pTypes, err := sType.QueryPropertyTypes().All(ctx)
		if err != nil {
			return fmt.Errorf("query property types of service type: %q, %w", sType.ID, err)
		}
		for _, pType := range pTypes {
			if err := m.deletePropertyType(ctx, client, pType); err != nil {
				return err
			}
		}
		endpoints, err := sType.QueryServices().QueryEndpoints().All(ctx)
		if err != nil {
			return fmt.Errorf("query service endpoints of service type: %q, %w", sType.ID, err)
		}
		for _, endpoint := range endpoints {
			if err := client.ServiceEndpoint.DeleteOne(endpoint).Exec(ctx); err != nil {
				return fmt.Errorf("deleting service endpoint of service type: %q, %w", endpoint.ID, err)
			}
		}
		m.logger.For(ctx).Info("deleted endpoints",
			zap.Int("id", sType.ID),
			zap.Int("count", len(endpoints)))

		services, err := sType.QueryServices().All(ctx)
		if err != nil {
			return fmt.Errorf("query services of type id: %d, %w", sType.ID, err)
		}
		for _, s := range services {
			if err := client.Service.DeleteOne(s).
				Exec(ctx); err != nil {
				return fmt.Errorf("delete service of service type: %q, %w", s.ID, err)
			}
		}
		m.logger.For(ctx).Info("deleted services",
			zap.Int("id", sType.ID),
			zap.Int("count", len(services)))

		endpointDefs, err := sType.QueryEndpointDefinitions().All(ctx)
		if err != nil {
			return fmt.Errorf("query endpoint definitions of type id: %d, %w", sType.ID, err)
		}
		for _, endpointDef := range endpointDefs {
			if err := client.ServiceEndpointDefinition.DeleteOne(endpointDef).
				Exec(ctx); err != nil {
				return fmt.Errorf("delete ednpoint definition of service type: %q, %w", endpointDef.ID, err)
			}
		}
		m.logger.For(ctx).Info("deleted endpoint definitions",
			zap.Int("id", sType.ID),
			zap.Int("count", len(endpointDefs)))

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

func (m *jobs) deletePropertyType(ctx context.Context, client *ent.Client, pType *ent.PropertyType) error {
	m.logger.For(ctx).Info("deleting property type",
		zap.Int("id", pType.ID),
		zap.String("name", pType.Name))
	propsToDelete, err := pType.QueryProperties().All(ctx)
	if err != nil {
		return fmt.Errorf("query properties of type id: %d, %w", pType.ID, err)
	}
	for _, prop := range propsToDelete {
		if err := client.Property.DeleteOne(prop).Exec(ctx); err != nil {
			return fmt.Errorf("delete property: %d, %w", prop.ID, err)
		}
	}
	m.logger.For(ctx).Info("deleted properties",
		zap.Int("id", pType.ID),
		zap.Int("count", len(propsToDelete)))
	err = client.PropertyType.DeleteOne(pType).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete property type: %d, %w", pType.ID, err)
	}
	m.logger.For(ctx).Info("deleted property type",
		zap.Int("id", pType.ID),
		zap.String("name", pType.Name))
	return nil
}
