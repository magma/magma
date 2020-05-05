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
  DeleteFloorPlanMutation,
  DeleteFloorPlanMutationResponse,
  DeleteFloorPlanMutationVariables,
} from './__generated__/DeleteFloorPlanMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation DeleteFloorPlanMutation($id: ID!) {
    deleteFloorPlan(id: $id)
  }
`;

export default (
  variables: DeleteFloorPlanMutationVariables,
  callbacks?: MutationCallbacks<DeleteFloorPlanMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<DeleteFloorPlanMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
