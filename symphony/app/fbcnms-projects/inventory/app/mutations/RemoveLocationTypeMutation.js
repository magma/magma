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
  RemoveLocationTypeMutation,
  RemoveLocationTypeMutationResponse,
  RemoveLocationTypeMutationVariables,
} from './__generated__/RemoveLocationTypeMutation.graphql';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation RemoveLocationTypeMutation($id: ID!) {
    removeLocationType(id: $id)
  }
`;

export default (
  variables: RemoveLocationTypeMutationVariables,
  callbacks?: MutationCallbacks<RemoveLocationTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveLocationTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
