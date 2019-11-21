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
  EditWorkOrderTypeMutationResponse,
  EditWorkOrderTypeMutationVariables,
} from './__generated__/EditWorkOrderTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

export const mutation = graphql`
  mutation EditWorkOrderTypeMutation($input: EditWorkOrderTypeInput!) {
    editWorkOrderType(input: $input) {
      id
      name
      ...WorkOrderTypeItem_workOrderType
      ...AddEditWorkOrderTypeCard_editingWorkOrderType
    }
  }
`;

export default (
  variables: EditWorkOrderTypeMutationVariables,
  callbacks?: MutationCallbacks<EditWorkOrderTypeMutationResponse>,
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
