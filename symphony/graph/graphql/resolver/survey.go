// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
)

type surveyResolver struct{}

func (surveyResolver) CreationTimestamp(_ context.Context, obj *ent.Survey) (*int, error) {
	timestamp := int(obj.CreationTimestamp.Unix())
	if timestamp < 0 {
		return nil, nil
	}
	return &timestamp, nil
}

func (surveyResolver) CompletionTimestamp(_ context.Context, obj *ent.Survey) (int, error) {
	return int(obj.CompletionTimestamp.Unix()), nil
}

func (surveyResolver) LocationID(ctx context.Context, obj *ent.Survey) (int, error) {
	return obj.QueryLocation().OnlyID(ctx)
}

func (surveyResolver) SurveyResponses(ctx context.Context, obj *ent.Survey) ([]*ent.SurveyQuestion, error) {
	return obj.QueryQuestions().All(ctx)
}

func (surveyResolver) SourceFile(ctx context.Context, obj *ent.Survey) (*ent.File, error) {
	survey, err := obj.QuerySourceFile().Only(ctx)
	return survey, ent.MaskNotFound(err)
}

type surveyCellScanResolver struct{}

func (surveyCellScanResolver) NetworkType(_ context.Context, obj *ent.SurveyCellScan) (models.CellularNetworkType, error) {
	return models.CellularNetworkType(obj.NetworkType), nil
}

func (surveyCellScanResolver) Timestamp(_ context.Context, obj *ent.SurveyCellScan) (*int, error) {
	timestamp := int(obj.Timestamp.Unix())
	if timestamp < 0 {
		return nil, nil
	}
	return &timestamp, nil
}

type surveyWiFiScanResolver struct{}

func (surveyWiFiScanResolver) Timestamp(_ context.Context, obj *ent.SurveyWiFiScan) (int, error) {
	return int(obj.Timestamp.Unix()), nil
}
