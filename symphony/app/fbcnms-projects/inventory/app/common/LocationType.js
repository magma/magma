/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {LocationTypeNodesQuery} from './__generated__/LocationTypeNodesQuery.graphql';
import type {NamedNode} from './EntUtils';
import type {PropertyType} from './PropertyType';
import type {SurveyQuestionType} from '../components/configure/__generated__/AddEditLocationTypeCard_editingLocationType.graphql.js';

import {graphql} from 'relay-runtime';
import {useLazyLoadQuery} from 'react-relay/hooks';

export type LocationType = {
  ...NamedNode,
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

const locationTypeNodesQuery = graphql`
  query LocationTypeNodesQuery {
    locationTypes {
      edges {
        node {
          id
          name
        }
      }
    }
  }
`;

export type LocationTypeNode = $Exact<NamedNode>;

export function useLocationTypeNodes(): $ReadOnlyArray<LocationTypeNode> {
  const response = useLazyLoadQuery<LocationTypeNodesQuery>(
    locationTypeNodesQuery,
  );
  const locationTypesData = response.locationTypes?.edges || [];
  const locationTypes = locationTypesData.map(p => p.node).filter(Boolean);
  return locationTypes;
}
