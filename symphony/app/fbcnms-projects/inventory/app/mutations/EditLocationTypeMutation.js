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
  EditLocationTypeMutation,
  EditLocationTypeMutationResponse,
  EditLocationTypeMutationVariables,
} from './__generated__/EditLocationTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

export const mutation = graphql`
  mutation EditLocationTypeMutation($input: EditLocationTypeInput!) {
    editLocationType(input: $input) {
      id
      name
      ...LocationTypeItem_locationType
      ...AddEditLocationTypeCard_editingLocationType
    }
  }
`;

export default (
  variables: EditLocationTypeMutationVariables,
  callbacks?: MutationCallbacks<EditLocationTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditLocationTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
