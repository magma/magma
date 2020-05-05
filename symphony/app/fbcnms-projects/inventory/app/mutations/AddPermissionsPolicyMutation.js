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
  AddPermissionsPolicyMutation,
  AddPermissionsPolicyMutationResponse,
  AddPermissionsPolicyMutationVariables,
} from './__generated__/AddPermissionsPolicyMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation AddPermissionsPolicyMutation($input: AddPermissionsPolicyInput!) {
    addPermissionsPolicy(input: $input) {
      ...UserManagementUtils_policies @relay(mask: false)
    }
  }
`;

export default (
  variables: AddPermissionsPolicyMutationVariables,
  callbacks?: MutationCallbacks<AddPermissionsPolicyMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddPermissionsPolicyMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
