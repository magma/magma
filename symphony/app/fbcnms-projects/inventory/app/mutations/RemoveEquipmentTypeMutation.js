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
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {
  RemoveEquipmentTypeMutation,
  RemoveEquipmentTypeMutationResponse,
  RemoveEquipmentTypeMutationVariables,
} from './__generated__/RemoveEquipmentTypeMutation.graphql';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation RemoveEquipmentTypeMutation($id: ID!) {
    removeEquipmentType(id: $id)
  }
`;

export default (
  variables: RemoveEquipmentTypeMutationVariables,
  callbacks?: MutationCallbacks<RemoveEquipmentTypeMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveEquipmentTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
