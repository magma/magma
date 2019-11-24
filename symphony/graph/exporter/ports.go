// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/facebookincubator/symphony/cloud/ctxgroup"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/facebookincubator/symphony/cloud/log"
)

type portsRower struct {
	log log.Logger
}

func (er portsRower) rows(ctx context.Context, url *url.URL) ([][]string, error) {
	log := er.log.For(ctx)

	var (
		err            error
		filterInput    []*models.PortFilterInput
		portDataHeader = [...]string{bom + "Port ID", " Port Name", "Port Type", "Equipment Name", "Equipment Type"}
		parentsHeader  = [...]string{"Parent Equipment (3)", "Parent Equipment (2)", "Parent Equipment", "Equipment Position"}
		linkHeader     = [...]string{"Linked Port ID", "Linked Port name", "Linked Port Equipment ID", "Linked Port Equipment"}
	)
	filtersParam := url.Query().Get("filters")
	if filtersParam != "" {
		filterInput, err = paramToPortFilterInput(filtersParam)
		if err != nil {
			log.Error("cannot filter ports", zap.Error(err))
			return nil, errors.Wrap(err, "cannot filter ports")
		}
	}
	client := ent.FromContext(ctx)

	ports, err := resolverutil.PortSearch(ctx, client, filterInput, nil)
	if err != nil {
		log.Error("cannot query ports", zap.Error(err))
		return nil, errors.Wrap(err, "cannot query ports")
	}
	cg := ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))

	portsList := ports.Ports
	allrows := make([][]string, len(portsList)+1)

	var orderedLocTypes, propertyTypes []string
	cg.Go(func(ctx context.Context) (err error) {
		orderedLocTypes, err = locationTypeHierarchy(ctx, client)
		if err != nil {
			log.Error("cannot query location types", zap.Error(err))
			return errors.Wrap(err, "cannot query location types")
		}
		return nil
	})
	cg.Go(func(ctx context.Context) (err error) {
		portIDs := make([]string, len(portsList))
		for i, p := range portsList {
			portIDs[i] = p.ID
		}
		propertyTypes, err = propertyTypesSlice(ctx, portIDs, client, models.PropertyEntityPort)
		if err != nil {
			log.Error("cannot query property types", zap.Error(err))
			return errors.Wrap(err, "cannot query property types")
		}
		return nil
	})
	if err := cg.Wait(); err != nil {
		return nil, err
	}

	title := append(portDataHeader[:], orderedLocTypes...)
	title = append(title, parentsHeader[:]...)
	title = append(title, linkHeader[:]...)
	title = append(title, propertyTypes...)

	allrows[0] = title
	cg = ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))
	for i, value := range portsList {
		value, i := value, i
		cg.Go(func(ctx context.Context) error {
			row, err := portToSlice(ctx, value, orderedLocTypes, propertyTypes)
			if err != nil {
				return err
			}
			allrows[i+1] = row
			return nil
		})
	}
	if err := cg.Wait(); err != nil {
		log.Error("error in wait", zap.Error(err))
		return nil, errors.WithMessage(err, "error in wait")
	}
	return allrows, nil
}

func portToSlice(ctx context.Context, port *ent.EquipmentPort, orderedLocTypes []string, propertyTypes []string) ([]string, error) {
	return []string{port.ID, port.QueryDefinition().OnlyX(ctx).Name}, nil
}

func paramToPortFilterInput(params string) ([]*models.PortFilterInput, error) {
	var ret []*models.PortFilterInput
	var inputs []filterInput
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
		inp := models.PortFilterInput{
			FilterType:    models.PortFilterType(upperName),
			Operator:      models.FilterOperator(upperOp),
			StringValue:   &StringVal,
			PropertyValue: &propVal,
			IDSet:         f.IDSet,
			MaxDepth:      &maxDepth,
		}
		ret = append(ret, &inp)
	}
	return ret, nil
}
