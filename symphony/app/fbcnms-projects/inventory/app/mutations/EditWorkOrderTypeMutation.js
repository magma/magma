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
import type {
  EditWorkOrderTypeMutation,
  EditWorkOrderTypeMutationResponse,
  EditWorkOrderTypeMutationVariables,
} from './__generated__/EditWorkOrderTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

export const mutation = graphql`
  mutation EditWorkOrderTypeMutation($input: EditWorkOrderTypeInput!) {
    editWorkOrderType(input: $input) {
      id
      name
      ...AddEditWorkOrderTypeCard_editingWorkOrderType
    }
  }
`;

export default (
  variables: EditWorkOrderTypeMutationVariables,
  callbacks?: MutationCallbacks<EditWorkOrderTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditWorkOrderTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
