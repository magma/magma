/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import RelayEnvironemnt from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {
  RemoveEquipmentFromPositionMutation,
  RemoveEquipmentFromPositionMutationResponse,
  RemoveEquipmentFromPositionMutationVariables,
} from './__generated__/RemoveEquipmentFromPositionMutation.graphql';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation RemoveEquipmentFromPositionMutation(
    $position_id: ID!
    $work_order_id: ID
  ) {
    removeEquipmentFromPosition(
      positionId: $position_id
      workOrderId: $work_order_id
    ) {
      ...EquipmentPropertiesCard_position
    }
  }
`;

export default (
  variables: RemoveEquipmentFromPositionMutationVariables,
  callbacks?: MutationCallbacks<RemoveEquipmentFromPositionMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveEquipmentFromPositionMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
