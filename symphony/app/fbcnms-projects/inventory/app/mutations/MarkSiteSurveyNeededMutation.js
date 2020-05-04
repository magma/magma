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
  MarkSiteSurveyNeededMutation,
  MarkSiteSurveyNeededMutationResponse,
  MarkSiteSurveyNeededMutationVariables,
} from './__generated__/MarkSiteSurveyNeededMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation MarkSiteSurveyNeededMutation($locationId: ID!, $needed: Boolean!) {
    markSiteSurveyNeeded(locationId: $locationId, needed: $needed) {
      ...LocationsTree_location @relay(mask: false)
    }
  }
`;

export default (
  variables: MarkSiteSurveyNeededMutationVariables,
  callbacks?: MutationCallbacks<MarkSiteSurveyNeededMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<MarkSiteSurveyNeededMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
