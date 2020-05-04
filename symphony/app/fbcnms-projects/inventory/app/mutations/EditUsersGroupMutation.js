/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import RelayEnvironemnt from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';
import type {
  EditUsersGroupMutation,
  EditUsersGroupMutationResponse,
  EditUsersGroupMutationVariables,
} from './__generated__/EditUsersGroupMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation EditUsersGroupMutation($input: EditUsersGroupInput!) {
    editUsersGroup(input: $input) {
      id
      name
      description
      status
      members {
        id
        authID
      }
    }
  }
`;

export default (
  variables: EditUsersGroupMutationVariables,
  callbacks?: MutationCallbacks<EditUsersGroupMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditUsersGroupMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
