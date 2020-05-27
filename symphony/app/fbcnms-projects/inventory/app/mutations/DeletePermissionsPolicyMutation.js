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
  DeletePermissionsPolicyMutation,
  DeletePermissionsPolicyMutationResponse,
  DeletePermissionsPolicyMutationVariables,
} from './__generated__/DeletePermissionsPolicyMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

import RelayEnvironemnt from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

const mutation = graphql`
  mutation DeletePermissionsPolicyMutation($id: ID!) {
    deletePermissionsPolicy(id: $id)
  }
`;

export default (
  variables: DeletePermissionsPolicyMutationVariables,
  callbacks?: MutationCallbacks<DeletePermissionsPolicyMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<DeletePermissionsPolicyMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
