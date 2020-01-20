/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {PropertyType} from './PropertyType';
import type {SurveyQuestionType} from '../components/configure/__generated__/AddEditLocationTypeCard_editingLocationType.graphql.js';

export type LocationType = {
  id: string,
  name: string,
  mapType: string,
  mapZoomLevel: string,
  propertyTypes: Array<PropertyType>,
  numberOfLocations: number,
  surveyTemplateCategories: SurveyTemplateCategory[],
  isSite: boolean,
  index?: number,
};

export type SurveyTemplateCategory = {
  id: string,
  categoryTitle: string,
  categoryDescription: string,
  surveyTemplateQuestions: SurveyTemplateQuestion[],
};

export type SurveyTemplateQuestion = {
  id: string,
  questionTitle: string,
  questionDescription: string,
  questionType: SurveyQuestionType,
  index: number,
};

export type LocationTypeIndex = {
  locationTypeID: string,
  index: number,
};
