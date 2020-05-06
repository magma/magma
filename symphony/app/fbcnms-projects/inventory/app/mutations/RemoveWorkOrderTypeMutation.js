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
  RemoveWorkOrderTypeMutation,
  RemoveWorkOrderTypeMutationResponse,
  RemoveWorkOrderTypeMutationVariables,
} from './__generated__/RemoveWorkOrderTypeMutation.graphql';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation RemoveWorkOrderTypeMutation($id: ID!) {
    removeWorkOrderType(id: $id)
  }
`;

export default (
  variables: RemoveWorkOrderTypeMutationVariables,
  callbacks?: MutationCallbacks<RemoveWorkOrderTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveWorkOrderTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
