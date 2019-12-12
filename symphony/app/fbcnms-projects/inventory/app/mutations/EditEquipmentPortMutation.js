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
  EditEquipmentPortMutation,
  EditEquipmentPortMutationResponse,
  EditEquipmentPortMutationVariables,
} from './__generated__/EditEquipmentPortMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

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
  updater?: (store: any) => void,
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
