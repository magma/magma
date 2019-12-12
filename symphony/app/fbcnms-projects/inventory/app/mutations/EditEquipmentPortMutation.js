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
  EditEquipmentPortMutation,
  EditEquipmentPortMutationResponse,
  EditEquipmentPortMutationVariables,
} from './__generated__/EditEquipmentPortMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {StoreUpdater} from '../common/RelayEnvironment';

export const mutation = graphql`
  mutation EditEquipmentPortMutation($input: EditEquipmentPortInput!) {
    editEquipmentPort(input: $input) {
      id
      ...EquipmentPortsTable_port @relay(mask: false)
    }
  }
`;

export default (
  variables: EditEquipmentPortMutationVariables,
  callbacks?: MutationCallbacks<EditEquipmentPortMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditEquipmentPortMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
