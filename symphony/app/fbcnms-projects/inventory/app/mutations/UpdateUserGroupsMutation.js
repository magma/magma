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
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';
import type {
  UpdateUserGroupsMutation,
  UpdateUserGroupsMutationResponse,
  UpdateUserGroupsMutationVariables,
} from './__generated__/UpdateUserGroupsMutation.graphql';

const mutation = graphql`
  mutation UpdateUserGroupsMutation($input: UpdateUserGroupsInput!) {
    updateUserGroups(input: $input) {
      ...UserManagementUtils_user @relay(mask: false)
    }
  }
`;

export default (
  variables: UpdateUserGroupsMutationVariables,
  callbacks?: MutationCallbacks<UpdateUserGroupsMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<UpdateUserGroupsMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
