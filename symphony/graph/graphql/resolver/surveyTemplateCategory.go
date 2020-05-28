// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/ent"
)

type surveyTemplateCategoryResolver struct{}

func (surveyTemplateCategoryResolver) SurveyTemplateQuestions(ctx context.Context, obj *ent.SurveyTemplateCategory) ([]*ent.SurveyTemplateQuestion, error) {
	return obj.QuerySurveyTemplateQuestions().All(ctx)
}
