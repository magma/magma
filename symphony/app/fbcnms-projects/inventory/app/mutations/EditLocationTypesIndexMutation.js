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
  EditLocationTypesIndexMutationResponse,
  EditLocationTypesIndexMutationVariables,
} from './__generated__/EditLocationTypesIndexMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

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
  // eslint-disable-next-line flowtype/no-weak-types
  updater?: (store: any) => void,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
