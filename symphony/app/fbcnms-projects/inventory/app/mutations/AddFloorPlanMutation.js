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
  AddFloorPlanMutationResponse,
  AddFloorPlanMutationVariables,
} from './__generated__/AddFloorPlanMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';

const mutation = graphql`
  mutation AddFloorPlanMutation($input: AddFloorPlanInput!) {
    addFloorPlan(input: $input) {
      id
      name
    }
  }
`;

export default (
  variables: AddFloorPlanMutationVariables,
  callbacks?: MutationCallbacks<AddFloorPlanMutationResponse>,
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
