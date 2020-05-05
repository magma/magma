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
  UpdateUsersGroupMembersMutation,
  UpdateUsersGroupMembersMutationResponse,
  UpdateUsersGroupMembersMutationVariables,
} from './__generated__/UpdateUsersGroupMembersMutation.graphql';

const mutation = graphql`
  mutation UpdateUsersGroupMembersMutation(
    $input: UpdateUsersGroupMembersInput!
  ) {
    updateUsersGroupMembers(input: $input) {
      ...UserManagementUtils_group @relay(mask: false)
    }
  }
`;

export default (
  variables: UpdateUsersGroupMembersMutationVariables,
  callbacks?: MutationCallbacks<UpdateUsersGroupMembersMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<UpdateUsersGroupMembersMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
