// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
)

type surveyTemplateQuestionResolver struct{}

func (surveyTemplateQuestionResolver) QuestionType(ctx context.Context, obj *ent.SurveyTemplateQuestion) (models.SurveyQuestionType, error) {
	return models.SurveyQuestionType(obj.QuestionType), nil
}
