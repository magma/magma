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
	"time"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/resolverutil"
	"github.com/facebookincubator/symphony/pkg/ctxgroup"
	"github.com/facebookincubator/symphony/pkg/log"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type woFilterInput struct {
	Name          models.EquipmentFilterType `json:"name"`
	Operator      models.FilterOperator      `jsons:"operator"`
	StringValue   string                     `json:"stringValue"`
	IDSet         []string                   `json:"idSet"`
	StringSet     []string                   `json:"stringSet"`
	PropertyValue models.PropertyTypeInput   `json:"propertyValue"`
	BoolValue     bool                       `json:"boolValue"`
}

type woRower struct {
	log log.Logger
}

func (er woRower) rows(ctx context.Context, url *url.URL) ([][]string, error) {
	var (
		logger       = er.log.For(ctx)
		err          error
		filterInput  []*models.WorkOrderFilterInput
		woDataHeader = [...]string{bom + "Work Order ID", "Work Order Name", "Project Name", "Status", "Assignee", "Owner", "Priority", "Created date", "Target date", "Location"}
	)
	filtersParam := url.Query().Get("filters")
	if filtersParam != "" {
		filterInput, err = paramToWOFilterInput(filtersParam)
		if err != nil {
			logger.Error("cannot filter work orders", zap.Error(err))
			return nil, errors.Wrap(err, "cannot filter work orders")
		}
	}
	client := ent.FromContext(ctx)

	searchResult, err := resolverutil.WorkOrderSearch(ctx, client, filterInput, nil)
	if err != nil {
		logger.Error("cannot query work orders", zap.Error(err))
		return nil, errors.Wrap(err, "cannot query work orders")
	}

	wosList := searchResult.WorkOrders
	allrows := make([][]string, len(wosList)+1)

	woIDs := make([]string, len(wosList))
	for i, w := range wosList {
		woIDs[i] = w.ID
	}
	propertyTypes, err := propertyTypesSlice(ctx, woIDs, client, models.PropertyEntityWorkOrders)
	if err != nil {
		logger.Error("cannot query property types", zap.Error(err))
		return nil, errors.Wrap(err, "cannot query property types")
	}

	title := append(woDataHeader[:], propertyTypes...)

	allrows[0] = title
	cg := ctxgroup.WithContext(ctx, ctxgroup.MaxConcurrency(32))
	for i, wo := range wosList {
		wo, i := wo, i

		cg.Go(func(ctx context.Context) error {
			row, err := woToSlice(ctx, wo, propertyTypes)
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

func woToSlice(ctx context.Context, wo *ent.WorkOrder, propertyTypes []string) ([]string, error) {
	properties, err := propertiesSlice(ctx, wo, propertyTypes, models.PropertyEntityWorkOrders)
	if err != nil {
		return nil, err
	}
	var projName, locName string

	proj, err := wo.QueryProject().Only(ctx)
	if ent.MaskNotFound(err) != nil {
		return nil, err
	}
	if proj != nil {
		projName = proj.Name
	}

	loc, err := wo.QueryLocation().Only(ctx)
	if ent.MaskNotFound(err) != nil {
		return nil, err
	}
	if loc != nil {
		locName = loc.Name
		parent, err := loc.QueryParent().Only(ctx)
		if err == nil && parent != nil {
			locName = parent.Name + "; " + locName
		}
	}

	row := []string{wo.ID, wo.Name, projName, wo.Status, wo.Assignee, wo.OwnerName, wo.Priority, getStringDate(wo.CreationDate), getStringDate(wo.InstallDate), locName}

	row = append(row, properties...)

	return row, nil
}

func getStringDate(t time.Time) string {
	y, m, d := t.Date()
	if y != 1 || m != time.January || d != 1 {
		return fmt.Sprintf("%d %v %d", d, m.String(), y)
	}
	return ""
}

func paramToWOFilterInput(params string) ([]*models.WorkOrderFilterInput, error) {
	var inputs []woFilterInput
	err := json.Unmarshal([]byte(params), &inputs)
	if err != nil {
		return nil, err
	}

	ret := make([]*models.WorkOrderFilterInput, 0, len(inputs))
	for _, f := range inputs {
		upperName := strings.ToUpper(f.Name.String())
		upperOp := strings.ToUpper(f.Operator.String())
		propertyValue := f.PropertyValue
		inp := models.WorkOrderFilterInput{
			FilterType:    models.WorkOrderFilterType(upperName),
			Operator:      models.FilterOperator(upperOp),
			StringValue:   pointer.ToString(f.StringValue),
			IDSet:         f.IDSet,
			StringSet:     f.StringSet,
			PropertyValue: &propertyValue,
			MaxDepth:      pointer.ToInt(5),
		}
		ret = append(ret, &inp)
	}
	return ret, nil
}
