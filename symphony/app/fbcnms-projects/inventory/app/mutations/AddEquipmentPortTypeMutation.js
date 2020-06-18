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
  AddEquipmentPortTypeMutation,
  AddEquipmentPortTypeMutationResponse,
  AddEquipmentPortTypeMutationVariables,
} from './__generated__/AddEquipmentPortTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

const mutation = graphql`
  mutation AddEquipmentPortTypeMutation($input: AddEquipmentPortTypeInput!) {
    addEquipmentPortType(input: $input) {
      id
      name
      ...EquipmentPortTypeItem_equipmentPortType
    }
  }
`;

export default (
  variables: AddEquipmentPortTypeMutationVariables,
  callbacks?: MutationCallbacks<AddEquipmentPortTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddEquipmentPortTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
