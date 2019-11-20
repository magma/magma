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
  AddEquipmentPortTypeMutationResponse,
  AddEquipmentPortTypeMutationVariables,
} from './__generated__/AddEquipmentPortTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';

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
