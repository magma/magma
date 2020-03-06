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
  RemoveProjectMutation,
  // $FlowFixMe (T62907961) Relay flow types
  RemoveProjectMutationMutationResponse,
  // $FlowFixMe (T62907961) Relay flow types
  RemoveProjectMutationMutationVariables,
} from './__generated__/RemoveProjectMutation.graphql';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation RemoveProjectMutation($id: ID!) {
    deleteProject(id: $id)
  }
`;

export default (
  variables: RemoveProjectMutationMutationVariables,
  callbacks?: MutationCallbacks<RemoveProjectMutationMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveProjectMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
