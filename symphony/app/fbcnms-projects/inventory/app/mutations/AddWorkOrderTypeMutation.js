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
  AddWorkOrderTypeMutation,
  AddWorkOrderTypeMutationResponse,
  AddWorkOrderTypeMutationVariables,
} from './__generated__/AddWorkOrderTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation AddWorkOrderTypeMutation($input: AddWorkOrderTypeInput!) {
    addWorkOrderType(input: $input) {
      id
      name
      ...AddEditWorkOrderTypeCard_editingWorkOrderType
    }
  }
`;

export default (
  variables: AddWorkOrderTypeMutationVariables,
  callbacks?: MutationCallbacks<AddWorkOrderTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddWorkOrderTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
