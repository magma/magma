// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

type surveyQuestionResolver struct{}

func (surveyQuestionResolver) PhotoData(ctx context.Context, obj *ent.SurveyQuestion) (*ent.File, error) {
	return obj.QueryPhotoData().Only(ctx)
}

func (surveyQuestionResolver) QuestionFormat(ctx context.Context, obj *ent.SurveyQuestion) (*models.SurveyQuestionType, error) {
	typ := models.SurveyQuestionType(obj.QuestionFormat)
	return &typ, nil
}

func (surveyQuestionResolver) DateData(ctx context.Context, obj *ent.SurveyQuestion) (*int, error) {
	secs := int(obj.DateData.Unix())
	return &secs, nil
}

func (surveyQuestionResolver) WifiData(ctx context.Context, obj *ent.SurveyQuestion) ([]*ent.SurveyWiFiScan, error) {
	return obj.QueryWifiScan().All(ctx)
}

func (surveyQuestionResolver) CellData(ctx context.Context, obj *ent.SurveyQuestion) ([]*ent.SurveyCellScan, error) {
	return obj.QueryCellScan().All(ctx)
}

func (surveyQuestionResolver) Images(ctx context.Context, obj *ent.SurveyQuestion) ([]*ent.File, error) {
	return obj.QueryImages().All(ctx)
}
