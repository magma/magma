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
  AddLocationTypeMutation,
  AddLocationTypeMutationResponse,
  AddLocationTypeMutationVariables,
} from './__generated__/AddLocationTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation AddLocationTypeMutation($input: AddLocationTypeInput!) {
    addLocationType(input: $input) {
      id
      name
      ...LocationTypeItem_locationType
      ...AddEditLocationTypeCard_editingLocationType
    }
  }
`;

export default (
  variables: AddLocationTypeMutationVariables,
  callbacks?: MutationCallbacks<AddLocationTypeMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddLocationTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
