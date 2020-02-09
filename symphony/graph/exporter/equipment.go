// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/AlekSi/pointer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	bom = "\uFEFF"
)

type equipmentFilterInput struct {
	Name          models.EquipmentFilterType `json:"name"`
	Operator      models.FilterOperator      `jsons:"operator"`
	StringValue   string                     `json:"stringValue"`
	IDSet         []string                   `json:"idSet"`
	StringSet     []string                   `json:"stringSet"`
	PropertyValue models.PropertyTypeInput   `json:"propertyValue"`
}

type equipmentRower struct {
	log log.Logger
}

func (er equipmentRower) rows(ctx context.Context, url *url.URL) ([][]string, error) {
	var (
		logger          = er.log.For(ctx)
		err             error
		filterInput     []*models.EquipmentFilterInput
		equipDataHeader = [...]string{bom + "Equipment ID", "Equipment Name", "Equipment Type", "External ID"}
		parentsHeader   = [...]string{"Parent Equipment (3)", "Position (3)", "Parent Equipment (2)", "Position (2)", "Parent Equipment", "Equipment Position"}
	)
	filtersParam := url.Query().Get("filters")
	if filtersParam != "" {
		filterInput, err = paramToFilterInput(filtersParam)
		if err != nil {
			logger.Error("cannot filter equipment", zap.Error(err))
			return nil, errors.Wrap(err, "cannot filter equipment")
		}
	}
	client := ent.FromContext(ctx)

	equips, err := resolverutil.EquipmentSearch(ctx, client, filterInput, nil)
	if err != nil {
		logger.Error("cannot query equipment", zap.Error(err))
		return nil, errors.Wrap(err, "cannot query equipment")
	}
	cg := ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))

	equipList := equips.Equipment
	allrows := make([][]string, len(equipList)+1)

	var orderedLocTypes, propertyTypes []string
	cg.Go(func(ctx context.Context) (err error) {
		orderedLocTypes, err = locationTypeHierarchy(ctx, client)
		if err != nil {
			logger.Error("cannot query location types", zap.Error(err))
			return errors.Wrap(err, "cannot query location types")
		}
		return nil
	})
	cg.Go(func(ctx context.Context) (err error) {
		equipIDs := make([]string, len(equipList))
		for i, e := range equips.Equipment {
			equipIDs[i] = e.ID
		}
		propertyTypes, err = propertyTypesSlice(ctx, equipIDs, client, models.PropertyEntityEquipment)
		if err != nil {
			logger.Error("cannot query property types", zap.Error(err))
			return errors.Wrap(err, "cannot query property types")
		}
		return nil
	})
	if err := cg.Wait(); err != nil {
		return nil, err
	}

	title := append(equipDataHeader[:], orderedLocTypes...)
	title = append(title, parentsHeader[:]...)
	title = append(title, propertyTypes...)

	allrows[0] = title
	cg = ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))
	for i, value := range equipList {
		value, i := value, i
		cg.Go(func(ctx context.Context) error {
			row, err := equipToSlice(ctx, value, orderedLocTypes, propertyTypes)
			if err != nil {
				return err
			}
			allrows[i+1] = row
			return nil
		})
	}
	if err := cg.Wait(); err != nil {
		logger.Error("error in wait", zap.Error(err))
		return nil, errors.WithMessage(err, "error in wait")
	}
	return allrows, nil
}

func paramToFilterInput(params string) ([]*models.EquipmentFilterInput, error) {
	var inputs []equipmentFilterInput
	err := json.Unmarshal([]byte(params), &inputs)
	if err != nil {
		return nil, err
	}

	returnType := make([]*models.EquipmentFilterInput, 0, len(inputs))
	for _, f := range inputs {
		upperName := strings.ToUpper(f.Name.String())
		upperOp := strings.ToUpper(f.Operator.String())
		propertyValue := f.PropertyValue
		inp := models.EquipmentFilterInput{
			FilterType:    models.EquipmentFilterType(upperName),
			Operator:      models.FilterOperator(upperOp),
			StringValue:   pointer.ToString(f.StringValue),
			PropertyValue: &propertyValue,
			IDSet:         f.IDSet,
			StringSet:     f.StringSet,
			MaxDepth:      pointer.ToInt(5),
		}
		returnType = append(returnType, &inp)
	}
	return returnType, nil
}

func equipToSlice(ctx context.Context, equipment *ent.Equipment, orderedLocTypes []string, propertyTypes []string) ([]string, error) {
	var (
		lParents, properties []string
		eParents             = make([]string, maxEquipmentParents*2)
	)
	g := ctxgroup.WithContext(ctx)
	g.Go(func(ctx context.Context) (err error) {
		lParents, err = locationHierarchyForEquipment(ctx, equipment, orderedLocTypes)
		return err
	})
	g.Go(func(ctx context.Context) (err error) {
		properties, err = propertiesSlice(ctx, equipment, propertyTypes, models.PropertyEntityEquipment)
		return err
	})
	g.Go(func(ctx context.Context) (err error) {
		pos, err := equipment.QueryParentPosition().Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return err
		}
		err = nil
		if pos != nil {
			eParents = parentHierarchyWithAllPositions(ctx, *equipment)
		}
		return
	})
	if err := g.Wait(); err != nil {
		return nil, err
	}
	row := []string{equipment.ID, equipment.Name, equipment.QueryType().OnlyX(ctx).Name, equipment.ExternalID}
	row = append(row, lParents...)
	row = append(row, eParents...)
	row = append(row, properties...)

	return row, nil
}
