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
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {
  RemoveSiteSurveyMutationResponse,
  RemoveSiteSurveyMutationVariables,
} from './__generated__/RemoveSiteSurveyMutation.graphql';

const mutation = graphql`
  mutation RemoveSiteSurveyMutation($id: ID!) {
    removeSiteSurvey(id: $id)
  }
`;

export default (
  variables: RemoveSiteSurveyMutationVariables,
  callbacks?: MutationCallbacks<RemoveSiteSurveyMutationResponse>,
  updater?: (store: any) => void,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
