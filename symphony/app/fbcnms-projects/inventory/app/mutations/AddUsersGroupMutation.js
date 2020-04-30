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
  AddUsersGroupMutation,
  AddUsersGroupMutationResponse,
  AddUsersGroupMutationVariables,
} from './__generated__/AddUsersGroupMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation AddUsersGroupMutation($input: AddUsersGroupInput!) {
    addUsersGroup(input: $input) {
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
  variables: AddUsersGroupMutationVariables,
  callbacks?: MutationCallbacks<AddUsersGroupMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddUsersGroupMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
