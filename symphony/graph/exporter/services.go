// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/jobs"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/facebookincubator/symphony/pkg/log"
)

type servicesFilterInput struct {
	Name          models.ServiceFilterType `json:"name"`
	Operator      models.FilterOperator    `jsons:"operator"`
	StringValue   string                   `json:"stringValue"`
	IDSet         []string                 `json:"idSet"`
	StringSet     []string                 `json:"stringSet"`
	PropertyValue models.PropertyTypeInput `json:"propertyValue"`
}

type servicesRower struct {
	log log.Logger
}

func (er servicesRower) rows(ctx context.Context, url *url.URL) ([][]string, error) {
	var (
		logger      = er.log.For(ctx)
		err         error
		filterInput []*models.ServiceFilterInput
		dataHeader  = [...]string{bom + "Service ID", "Service Name", "Service Type", "Discovery Method", "Service External ID", "Customer Name", "Customer External ID", "Status"}
	)
	filtersParam := url.Query().Get("filters")
	if filtersParam != "" {
		filterInput, err = paramToServiceFilterInput(filtersParam)
		if err != nil {
			logger.Error("cannot filter services", zap.Error(err))
			return nil, errors.Wrap(err, "cannot filter services")
		}
	}
	client := ent.FromContext(ctx)

	services, err := resolverutil.ServiceSearch(ctx, client, filterInput, nil)
	if err != nil {
		logger.Error("cannot query services", zap.Error(err))
		return nil, errors.Wrap(err, "cannot query services")
	}
	cg := ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))

	servicesList := services.Services
	allRows := make([][]string, len(servicesList)+1)

	var propertyTypes []string
	cg.Go(func(ctx context.Context) error {
		serviceIDs := make([]int, len(servicesList))
		for i, l := range servicesList {
			serviceIDs[i] = l.ID
		}
		propertyTypes, err = propertyTypesSlice(ctx, serviceIDs, client, models.PropertyEntityService)
		if err != nil {
			logger.Error("cannot query property types", zap.Error(err))
			return errors.Wrap(err, "cannot query property types")
		}
		return nil
	})
	if err := cg.Wait(); err != nil {
		return nil, err
	}

	endpointHeader := make([]string, jobs.MaxEndpoints*3)
	iter := 0
	for i := 0; i < len(endpointHeader); i += 3 {
		iter++
		endpointHeader[i] = "Endpoint Definition " + strconv.Itoa(iter)
		endpointHeader[i+1] = "Location " + strconv.Itoa(iter)
		endpointHeader[i+2] = "Equipment " + strconv.Itoa(iter)
	}
	title := append(dataHeader[:], endpointHeader...)
	title = append(title, propertyTypes...)

	allRows[0] = title
	cg = ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))
	for i, value := range servicesList {
		value, i := value, i
		cg.Go(func(ctx context.Context) error {
			row, err := serviceToSlice(ctx, value, propertyTypes)
			if err != nil {
				return err
			}
			allRows[i+1] = row
			return nil
		})
	}
	if err := cg.Wait(); err != nil {
		logger.Error("error in wait", zap.Error(err))
		return nil, errors.WithMessage(err, "error in wait")
	}
	return allRows, nil
}

func serviceToSlice(ctx context.Context, service *ent.Service, propertyTypes []string) ([]string, error) {
	st, err := service.QueryType().Only(ctx)
	if err != nil {
		return nil, err
	}
	serviceType := st.Name

	var customerName, customerExternalID, externalID string
	customer, err := service.QueryCustomer().Only(ctx)
	if err == nil {
		customerName = customer.Name
		if customer.ExternalID != nil {
			customerExternalID = *customer.ExternalID
		}
	}

	if service.ExternalID != nil {
		externalID = *service.ExternalID
	}

	discoveryMethod := st.DiscoveryMethod.String()
	if st.DiscoveryMethod == "" {
		discoveryMethod = models.DiscoveryMethodManual.String()
	}

	properties, err := propertiesSlice(ctx, service, propertyTypes, models.PropertyEntityService)
	if err != nil {
		return nil, err
	}
	endpoints, err := endpointsToSlice(ctx, service, st)
	if err != nil {
		return nil, err
	}

	row := []string{strconv.Itoa(service.ID), service.Name, serviceType, discoveryMethod, externalID, customerName, customerExternalID, service.Status}
	row = append(row, endpoints...)
	row = append(row, properties...)

	return row, nil
}

func endpointsToSlice(ctx context.Context, service *ent.Service, st *ent.ServiceType) ([]string, error) {
	endpointsData := make([]string, jobs.MaxEndpoints*3)
	endpointDefs, err := st.QueryEndpointDefinitions().
		Order(ent.Asc(serviceendpointdefinition.FieldIndex)).All(ctx)
	if err != nil {
		return nil, err
	}

	if (len(endpointDefs) < 2 && len(endpointDefs) != 0) || len(endpointDefs) > jobs.MaxEndpoints {
		return nil, errors.New("[SKIPPING SERVICE TYPE] either too many or not enough endpoint types ")
	}
	for i, endpointDef := range endpointDefs {
		ind := i * 3
		e, err := service.QueryEndpoints().
			Where(serviceendpoint.HasDefinitionWith(serviceendpointdefinition.ID(endpointDef.ID))).
			QueryEquipment().Only(ctx)
		if ent.MaskNotFound(err) != nil {
			return nil, err
		}
		if ent.IsNotFound(err) {
			continue
		}

		loc, err := getLastLocations(ctx, e, 3)
		if err != nil || loc == nil {
			return nil, errors.Wrap(err, "error while getting first location of equipment")
		}
		endpointsData[ind] = endpointDef.Name
		endpointsData[ind+1] = *loc
		endpointsData[ind+2] = e.Name
	}
	return endpointsData, nil
}

func paramToServiceFilterInput(params string) ([]*models.ServiceFilterInput, error) {
	var inputs []servicesFilterInput
	if err := json.Unmarshal([]byte(params), &inputs); err != nil {
		return nil, err
	}
	ret := make([]*models.ServiceFilterInput, 0, len(inputs))
	for _, f := range inputs {
		upperName := strings.ToUpper(f.Name.String())
		upperOp := strings.ToUpper(f.Operator.String())
		propertyValue := f.PropertyValue
		intIDSet, err := toIntSlice(f.IDSet)
		if err != nil {
			return nil, fmt.Errorf("wrong id set %v: %w", f.IDSet, err)
		}
		inp := models.ServiceFilterInput{
			FilterType:    models.ServiceFilterType(upperName),
			Operator:      models.FilterOperator(upperOp),
			StringValue:   pointer.ToString(f.StringValue),
			PropertyValue: &propertyValue,
			IDSet:         intIDSet,
			StringSet:     f.StringSet,
			MaxDepth:      pointer.ToInt(5),
		}
		ret = append(ret, &inp)
	}
	return ret, nil
}
