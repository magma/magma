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
  EditUserMutation,
  EditUserMutationResponse,
  EditUserMutationVariables,
} from './__generated__/EditUserMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation EditUserMutation($input: EditUserInput!) {
    editUser(input: $input) {
      id
      authID
      firstName
      lastName
      email
      status
      role
    }
  }
`;

export default (
  variables: EditUserMutationVariables,
  callbacks?: MutationCallbacks<EditUserMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditUserMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
