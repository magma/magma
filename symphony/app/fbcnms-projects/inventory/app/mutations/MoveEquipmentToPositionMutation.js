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
  MoveEquipmentToPositionMutation,
  MoveEquipmentToPositionMutationResponse,
  MoveEquipmentToPositionMutationVariables,
} from './__generated__/MoveEquipmentToPositionMutation.graphql.js';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation MoveEquipmentToPositionMutation(
    $parent_equipment_id: ID!
    $position_definition_id: ID!
    $equipment_id: ID!
  ) {
    moveEquipmentToPosition(
      parentEquipmentId: $parent_equipment_id
      positionDefinitionId: $position_definition_id
      equipmentId: $equipment_id
    ) {
      ...EquipmentPropertiesCard_position
    }
  }
`;

export default (
  variables: MoveEquipmentToPositionMutationVariables,
  callbacks?: MutationCallbacks<MoveEquipmentToPositionMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<MoveEquipmentToPositionMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
