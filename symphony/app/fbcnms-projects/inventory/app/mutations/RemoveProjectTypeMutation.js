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
  RemoveProjectTypeMutation,
  RemoveProjectTypeMutationResponse,
  RemoveProjectTypeMutationVariables,
} from './__generated__/RemoveProjectTypeMutation.graphql';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation RemoveProjectTypeMutation($id: ID!) {
    deleteProjectType(id: $id)
  }
`;

export default (
  variables: RemoveProjectTypeMutationVariables,
  callbacks?: MutationCallbacks<RemoveProjectTypeMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveProjectTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
