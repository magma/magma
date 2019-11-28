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
	"github.com/facebookincubator/symphony/cloud/log"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type linksFilterInput struct {
	Name          models.LinkFilterType    `json:"name"`
	Operator      models.FilterOperator    `jsons:"operator"`
	StringValue   string                   `json:"stringValue"`
	IDSet         []string                 `json:"idSet"`
	PropertyValue models.PropertyTypeInput `json:"propertyValue"`
}

type linksRower struct {
	log log.Logger
}

func (er linksRower) rows(ctx context.Context, url *url.URL) ([][]string, error) {
	log := er.log.For(ctx)

	var (
		err             error
		filterInput     []*models.LinkFilterInput
		portADataHeader = [...]string{bom + "Link ID", "Port A ID", "Port A Name", "Port A Type", "Equipment A ID", "Equipment A Name", "Equipment A Type"}
		portBDataHeader = [...]string{"Port B ID", "Port B Name", "Port B Type", "Equipment B ID", "Equipment B Name", "Equipment B Type"}
	)
	filtersParam := url.Query().Get("filters")
	if filtersParam != "" {
		filterInput, err = paramToLinkFilterInput(filtersParam)
		if err != nil {
			log.Error("cannot filter links", zap.Error(err))
			return nil, errors.Wrap(err, "cannot filter links")
		}
	}
	client := ent.FromContext(ctx)

	links, err := resolverutil.LinkSearch(ctx, client, filterInput, nil)
	if err != nil {
		log.Error("cannot query links", zap.Error(err))
		return nil, errors.Wrap(err, "cannot query links")
	}
	cg := ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))

	linksList := links.Links
	allRows := make([][]string, len(linksList)+1)

	var propertyTypes []string
	cg.Go(func(ctx context.Context) error {
		linkIDs := make([]string, len(linksList))
		for i, l := range linksList {
			linkIDs[i] = l.ID
		}
		propertyTypes, err = propertyTypesSlice(ctx, linkIDs, client, models.PropertyEntityLink)
		if err != nil {
			log.Error("cannot query property types", zap.Error(err))
			return errors.Wrap(err, "cannot query property types")
		}
		return nil
	})
	if err := cg.Wait(); err != nil {
		return nil, err
	}

	title := append(portADataHeader[:], portBDataHeader[:]...)
	title = append(title, propertyTypes...)

	allRows[0] = title
	cg = ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))
	for i, value := range linksList {
		value, i := value, i
		cg.Go(func(ctx context.Context) error {
			row, err := linkToSlice(ctx, value, propertyTypes)
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

func linkToSlice(ctx context.Context, link *ent.Link, propertyTypes []string) ([]string, error) {
	return []string{link.ID}, nil
}

func paramToLinkFilterInput(params string) ([]*models.LinkFilterInput, error) {
	var ret []*models.LinkFilterInput
	var inputs []linksFilterInput
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
		inp := models.LinkFilterInput{
			FilterType:    models.LinkFilterType(upperName),
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
