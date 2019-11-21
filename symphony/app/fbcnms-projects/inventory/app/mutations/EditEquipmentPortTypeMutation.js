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
  EditEquipmentPortTypeMutationResponse,
  EditEquipmentPortTypeMutationVariables,
} from './__generated__/EditEquipmentPortTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

export const mutation = graphql`
  mutation EditEquipmentPortTypeMutation($input: EditEquipmentPortTypeInput!) {
    editEquipmentPortType(input: $input) {
      id
      name
      ...EquipmentPortTypeItem_equipmentPortType
      ...AddEditEquipmentPortTypeCard_editingEquipmentPortType
    }
  }
`;

export default (
  variables: EditEquipmentPortTypeMutationVariables,
  callbacks?: MutationCallbacks<EditEquipmentPortTypeMutationResponse>,
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
