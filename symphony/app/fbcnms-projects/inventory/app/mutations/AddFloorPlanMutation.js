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
  AddFloorPlanMutation,
  AddFloorPlanMutationResponse,
  AddFloorPlanMutationVariables,
} from './__generated__/AddFloorPlanMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation AddFloorPlanMutation($input: AddFloorPlanInput!) {
    addFloorPlan(input: $input) {
      id
      name
      image {
        ...FileAttachment_file
      }
    }
  }
`;

export default (
  variables: AddFloorPlanMutationVariables,
  callbacks?: MutationCallbacks<AddFloorPlanMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddFloorPlanMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
