/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  DeleteUsersGroupMutation,
  DeleteUsersGroupMutationResponse,
  DeleteUsersGroupMutationVariables,
} from './__generated__/DeleteUsersGroupMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

import RelayEnvironemnt from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

const mutation = graphql`
  mutation DeleteUsersGroupMutation($id: ID!) {
    deleteUsersGroup(id: $id)
  }
`;

export default (
  variables: DeleteUsersGroupMutationVariables,
  callbacks?: MutationCallbacks<DeleteUsersGroupMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<DeleteUsersGroupMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
