/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';
import type {
  EditLocationTypeSurveyTemplateCategoriesMutation,
  EditLocationTypeSurveyTemplateCategoriesMutationResponse,
  EditLocationTypeSurveyTemplateCategoriesMutationVariables,
} from './__generated__/EditLocationTypeSurveyTemplateCategoriesMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

export const mutation = graphql`
  mutation EditLocationTypeSurveyTemplateCategoriesMutation(
    $id: ID!
    $surveyTemplateCategories: [SurveyTemplateCategoryInput!]!
  ) {
    editLocationTypeSurveyTemplateCategories(
      id: $id
      surveyTemplateCategories: $surveyTemplateCategories
    ) {
      id
    }
  }
`;

export default (
  variables: EditLocationTypeSurveyTemplateCategoriesMutationVariables,
  callbacks?: MutationCallbacks<EditLocationTypeSurveyTemplateCategoriesMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditLocationTypeSurveyTemplateCategoriesMutation>(
    RelayEnvironment,
    {
      mutation,
      variables,
      updater,
      onCompleted,
      onError,
    },
  );
};
