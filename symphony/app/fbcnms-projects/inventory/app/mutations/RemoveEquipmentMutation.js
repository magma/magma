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
  RemoveEquipmentMutationMutationResponse,
  RemoveEquipmentMutationMutationVariables,
} from './__generated__/RemoveEquipmentMutation.graphql';

const mutation = graphql`
  mutation RemoveEquipmentMutation($id: ID!, $work_order_id: ID) {
    removeEquipment(id: $id, workOrderId: $work_order_id)
  }
`;

export default (
  variables: RemoveEquipmentMutationMutationVariables,
  callbacks?: MutationCallbacks<RemoveEquipmentMutationMutationResponse>,
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
