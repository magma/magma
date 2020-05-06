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
  RemoveEquipmentPortTypeMutation,
  RemoveEquipmentPortTypeMutationResponse,
  RemoveEquipmentPortTypeMutationVariables,
} from './__generated__/RemoveEquipmentPortTypeMutation.graphql';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation RemoveEquipmentPortTypeMutation($id: ID!) {
    removeEquipmentPortType(id: $id)
  }
`;

export default (
  variables: RemoveEquipmentPortTypeMutationVariables,
  callbacks?: MutationCallbacks<RemoveEquipmentPortTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveEquipmentPortTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
