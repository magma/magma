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
  // $FlowFixMe (T62907961) Relay flow types
  RemoveEquipmentMutationMutationResponse,
  // $FlowFixMe (T62907961) Relay flow types
  RemoveEquipmentMutationMutationVariables,
} from './__generated__/RemoveEquipmentMutation.graphql';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation RemoveEquipmentMutation($id: ID!, $work_order_id: ID) {
    removeEquipment(id: $id, workOrderId: $work_order_id)
  }
`;

export default (
  variables: RemoveEquipmentMutationMutationVariables,
  callbacks?: MutationCallbacks<RemoveEquipmentMutationMutationResponse>,
  updater?: StoreUpdater,
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
