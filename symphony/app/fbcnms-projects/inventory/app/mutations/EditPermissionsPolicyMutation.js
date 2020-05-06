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
  EditPermissionsPolicyMutation,
  EditPermissionsPolicyMutationResponse,
  EditPermissionsPolicyMutationVariables,
} from './__generated__/EditPermissionsPolicyMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation EditPermissionsPolicyMutation($input: EditPermissionsPolicyInput!) {
    editPermissionsPolicy(input: $input) {
      ...UserManagementUtils_policies @relay(mask: false)
    }
  }
`;

export default (
  variables: EditPermissionsPolicyMutationVariables,
  callbacks?: MutationCallbacks<EditPermissionsPolicyMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditPermissionsPolicyMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
