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
  EditLocationTypesIndexMutation,
  EditLocationTypesIndexMutationResponse,
  EditLocationTypesIndexMutationVariables,
} from './__generated__/EditLocationTypesIndexMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {StoreUpdater} from '../common/RelayEnvironment';

export const mutation = graphql`
  mutation EditLocationTypesIndexMutation(
    $locationTypeIndex: [LocationTypeIndex]!
  ) {
    editLocationTypesIndex(locationTypesIndex: $locationTypeIndex) {
      id
      name
      index
      ...LocationTypeItem_locationType
      ...AddEditLocationTypeCard_editingLocationType
    }
  }
`;

export default (
  variables: EditLocationTypesIndexMutationVariables,
  callbacks?: MutationCallbacks<EditLocationTypesIndexMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditLocationTypesIndexMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
