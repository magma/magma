// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
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

type linksFilterInput struct {
	Name          models.LinkFilterType    `json:"name"`
	Operator      models.FilterOperator    `jsons:"operator"`
	StringValue   string                   `json:"stringValue"`
	IDSet         []int                    `json:"idSet"`
	StringSet     []string                 `json:"stringSet"`
	PropertyValue models.PropertyTypeInput `json:"propertyValue"`
	MaxDepth      *int                     `json:"maxDepth"`
}

type linksRower struct {
	log log.Logger
}

func (er linksRower) rows(ctx context.Context, url *url.URL) ([][]string, error) {
	var (
		logger          = er.log.For(ctx)
		err             error
		filterInput     []*models.LinkFilterInput
		portADataHeader = [...]string{bom + "Link ID", "Port A Name", "Equipment A Name", "Equipment A Type"}
		portBDataHeader = [...]string{"Port B Name", "Equipment B Name", "Equipment B Type"}
		parentsAHeader  = [...]string{"Parent Equipment (3) A", "Position (3) A", "Parent Equipment (2) A", "Position (2) A", "Parent Equipment A", "Equipment Position A"}
		parentsBHeader  = [...]string{"Parent Equipment (3) B", "Position (3) B", "Parent Equipment (2) B", "Position (2) B", "Parent Equipment B", "Equipment Position B"}
		servicesHeader  = [...]string{"Service Names"}
	)
	filtersParam := url.Query().Get("filters")
	if filtersParam != "" {
		filterInput, err = paramToLinkFilterInput(filtersParam)
		if err != nil {
			logger.Error("cannot filter links", zap.Error(err))
			return nil, errors.Wrap(err, "cannot filter links")
		}
	}
	client := ent.FromContext(ctx)
	var orderedLocTypes, propertyTypes []string

	links, err := resolverutil.LinkSearch(ctx, client, filterInput, nil)
	if err != nil {
		logger.Error("cannot query links", zap.Error(err))
		return nil, errors.Wrap(err, "cannot query links")
	}

	linksList := links.Links
	allRows := make([][]string, len(linksList)+1)

	linkIDs := make([]int, len(linksList))
	for i, l := range linksList {
		linkIDs[i] = l.ID
	}
	cg := ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))
	cg.Go(func(ctx context.Context) error {
		if orderedLocTypes, err = locationTypeHierarchy(ctx, client); err != nil {
			logger.Error("cannot query location types", zap.Error(err))
			return errors.Wrap(err, "cannot query location types")
		}
		return nil
	})
	cg.Go(func(ctx context.Context) error {
		if propertyTypes, err = propertyTypesSlice(ctx, linkIDs, client, models.PropertyEntityLink); err != nil {
			logger.Error("cannot query property types", zap.Error(err))
			return errors.Wrap(err, "cannot query property types")
		}
		return nil
	})
	if err := cg.Wait(); err != nil {
		return nil, err
	}

	portAData := append(portADataHeader[:], orderedLocTypes...)
	portAData = append(portAData, parentsAHeader[:]...)

	portBData := append(portBDataHeader[:], orderedLocTypes...)
	portBData = append(portBData, parentsBHeader[:]...)

	title := append(portAData, portBData...)
	title = append(title, servicesHeader[:]...)
	title = append(title, propertyTypes...)

	allRows[0] = title
	cg = ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))
	for i, value := range linksList {
		value, i := value, i
		cg.Go(func(ctx context.Context) error {
			row, err := linkToSlice(ctx, value, propertyTypes, orderedLocTypes)
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

func linkToSlice(ctx context.Context, link *ent.Link, propertyTypes, orderedLocTypes []string) ([]string, error) {
	var (
		portData     = make(map[int][]string, 2)
		locationData = make(map[int][]string, 2)
		positionData = make(map[int][]string, 2)
	)
	ports, err := link.QueryPorts().All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "querying link for ports (id=%d)", link.ID)
	}
	if len(ports) != 2 {
		return nil, errors.Wrapf(err, "link must include 2 ports (link id=%d)", link.ID)
	}
	for i, port := range ports {
		portDefinition := port.QueryDefinition().OnlyX(ctx)

		portEquipment, err := port.QueryParent().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying parent for port (id=%d)", port.ID)
		}
		parentType, err := portEquipment.QueryType().Only(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "querying type for port parent (id=%d, parentID=%d)", port.ID, portEquipment.ID)
		}
		portData[i] = []string{portDefinition.Name, portEquipment.Name, parentType.Name}

		g := ctxgroup.WithContext(ctx)
		g.Go(func(ctx context.Context) (err error) {
			locationData[i], err = locationHierarchyForEquipment(ctx, portEquipment, orderedLocTypes)
			return err
		})
		g.Go(func(ctx context.Context) error {
			pos, err := portEquipment.QueryParentPosition().Only(ctx)
			if err != nil && !ent.IsNotFound(err) {
				return err
			}
			positionData[i] = make([]string, maxEquipmentParents*2)
			if pos != nil {
				positionData[i] = parentHierarchyWithAllPositions(ctx, *portEquipment)
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			return nil, err
		}
	}
	properties, err := propertiesSlice(ctx, link, propertyTypes, models.PropertyEntityLink)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create property slice for link (id=%d)", link.ID)
	}

	services, err := link.QueryService().All(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "querying link for services (id=%d)", link.ID)
	}
	var servicesList []string
	for _, service := range services {
		servicesList = append(servicesList, service.Name)
	}
	servicesStr := strings.Join(servicesList, ";")

	// Build the slice
	row := []string{strconv.Itoa(link.ID)}
	for i := 0; i < 2; i++ {
		row = append(row, portData[i]...)
		row = append(row, locationData[i]...)
		row = append(row, positionData[i]...)
	}
	row = append(row, servicesStr)
	row = append(row, properties...)
	return row, nil
}

func paramToLinkFilterInput(params string) ([]*models.LinkFilterInput, error) {
	var inputs []linksFilterInput
	err := json.Unmarshal([]byte(params), &inputs)
	if err != nil {
		return nil, err
	}

	ret := make([]*models.LinkFilterInput, 0, len(inputs))
	for _, f := range inputs {
		upperName := strings.ToUpper(f.Name.String())
		upperOp := strings.ToUpper(f.Operator.String())
		propVal := f.PropertyValue
		maxDepth := 5
		if f.MaxDepth != nil {
			maxDepth = *f.MaxDepth
		}
		inp := models.LinkFilterInput{
			FilterType:    models.LinkFilterType(upperName),
			Operator:      models.FilterOperator(upperOp),
			StringValue:   pointer.ToString(f.StringValue),
			PropertyValue: &propVal,
			IDSet:         f.IDSet,
			StringSet:     f.StringSet,
			MaxDepth:      &maxDepth,
		}
		ret = append(ret, &inp)
	}
	return ret, nil
}
