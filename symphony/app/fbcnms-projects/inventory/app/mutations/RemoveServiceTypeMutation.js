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
  RemoveServiceTypeMutation,
  RemoveServiceTypeMutationResponse,
  RemoveServiceTypeMutationVariables,
} from './__generated__/RemoveServiceTypeMutation.graphql';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation RemoveServiceTypeMutation($id: ID!) {
    removeServiceType(id: $id)
  }
`;

export default (
  variables: RemoveServiceTypeMutationVariables,
  callbacks?: MutationCallbacks<RemoveServiceTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveServiceTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
