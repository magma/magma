// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type locationsFilterInput struct {
	Name          models.LocationFilterType `json:"name"`
	Operator      models.FilterOperator     `jsons:"operator"`
	StringValue   string                    `json:"stringValue"`
	IDSet         []string                  `json:"idSet"`
	PropertyValue models.PropertyTypeInput  `json:"propertyValue"`
	MaxDepth      *int                      `json:"maxDepth"`
	BoolValue     *bool                     `json:"boolValue"`
}

type locationsRower struct {
	log log.Logger
}

func (er locationsRower) rows(ctx context.Context, url *url.URL) ([][]string, error) {
	log := er.log.For(ctx)
	var (
		err              error
		filterInput      []*models.LocationFilterInput
		locationIDHeader = [...]string{bom + "Location ID"}
		fixedHeaders     = [...]string{"External ID", "Latitude", "Longitude"}
	)
	filtersParam := url.Query().Get("filters")
	if filtersParam != "" {
		filterInput, err = paramToLocationFilterInput(filtersParam)
		if err != nil {
			log.Error("cannot filter location", zap.Error(err))
			return nil, errors.Wrap(err, "cannot filter location")
		}
	}
	client := ent.FromContext(ctx)

	locations, err := resolverutil.LocationSearch(ctx, client, filterInput, nil)
	if err != nil {
		log.Error("cannot query location", zap.Error(err))
		return nil, errors.Wrap(err, "cannot query location")
	}

	locationsList := locations.Locations
	allRows := make([][]string, len(locationsList)+1)

	locationIDs := make([]string, len(locationsList))
	for i, l := range locationsList {
		locationIDs[i] = l.ID
	}

	var orderedLocTypes, propertyTypes []string
	cg := ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))
	cg.Go(func(ctx context.Context) (err error) {
		orderedLocTypes, err = locationTypeHierarchy(ctx, client)
		if err != nil {
			log.Error("cannot query location types", zap.Error(err))
			return errors.Wrap(err, "cannot query location types")
		}
		return nil
	})
	cg.Go(func(ctx context.Context) (err error) {
		locationIDs := make([]string, len(locationsList))
		for i, l := range locationsList {
			locationIDs[i] = l.ID
		}
		propertyTypes, err = propertyTypesSlice(ctx, locationIDs, client, models.PropertyEntityLocation)
		if err != nil {
			log.Error("cannot query property types", zap.Error(err))
			return errors.Wrap(err, "cannot query property types")
		}
		return nil
	})
	if err := cg.Wait(); err != nil {
		return nil, err
	}

	title := append(locationIDHeader[:], orderedLocTypes...)
	title = append(title, fixedHeaders[:]...)
	title = append(title, propertyTypes...)

	allRows[0] = title
	cg = ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))
	for i, value := range locationsList {
		value, i := value, i
		cg.Go(func(ctx context.Context) error {
			row, err := locationToSlice(ctx, value, orderedLocTypes, propertyTypes)
			if err != nil {
				return err
			}
			allRows[i+1] = row
			return nil
		})
	}
	if err := cg.Wait(); err != nil {
		log.Error("error in wait", zap.Error(err))
		return nil, errors.WithMessage(err, "error in wait")
	}
	return allRows, nil
}

func locationToSlice(ctx context.Context, location *ent.Location, orderedLocTypes []string, propertyTypes []string) ([]string, error) {
	var (
		lParents, properties []string
	)
	g := ctxgroup.WithContext(ctx)
	g.Go(func(ctx context.Context) (err error) {
		lParents, err = locationHierarchy(ctx, location, orderedLocTypes)
		return err
	})
	g.Go(func(ctx context.Context) (err error) {
		properties, err = propertiesSlice(ctx, location, propertyTypes, models.PropertyEntityLocation)
		return err
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	lat := fmt.Sprintf("%f", location.Latitude)
	long := fmt.Sprintf("%f", location.Longitude)

	fixedData := []string{location.ExternalID, lat, long}

	row := []string{location.ID}
	row = append(row, lParents...)
	row = append(row, fixedData...)
	row = append(row, properties...)

	return row, nil
}

func paramToLocationFilterInput(params string) ([]*models.LocationFilterInput, error) {
	var ret []*models.LocationFilterInput
	var inputs []locationsFilterInput
	err := json.Unmarshal([]byte(params), &inputs)
	if err != nil {
		return nil, err
	}

	for _, f := range inputs {
		upperName := strings.ToUpper(f.Name.String())
		upperOp := strings.ToUpper(f.Operator.String())
		StringVal := f.StringValue
		propVal := f.PropertyValue
		maxDepth := 5
		if f.MaxDepth != nil {
			maxDepth = *f.MaxDepth
		}
		inp := models.LocationFilterInput{
			FilterType:    models.LocationFilterType(upperName),
			Operator:      models.FilterOperator(upperOp),
			StringValue:   &StringVal,
			PropertyValue: &propVal,
			IDSet:         f.IDSet,
			MaxDepth:      &maxDepth,
			BoolValue:     f.BoolValue,
		}
		ret = append(ret, &inp)
	}
	return ret, nil
}
