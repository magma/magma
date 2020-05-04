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
  RemoveEquipmentMutation,
  RemoveEquipmentMutationResponse,
  RemoveEquipmentMutationVariables,
} from './__generated__/RemoveEquipmentMutation.graphql';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation RemoveEquipmentMutation($id: ID!, $work_order_id: ID) {
    removeEquipment(id: $id, workOrderId: $work_order_id)
  }
`;

export default (
  variables: RemoveEquipmentMutationVariables,
  callbacks?: MutationCallbacks<RemoveEquipmentMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveEquipmentMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
