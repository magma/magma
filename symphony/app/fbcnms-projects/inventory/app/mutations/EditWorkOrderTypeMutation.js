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
  EditWorkOrderTypeMutation,
  EditWorkOrderTypeMutationResponse,
  EditWorkOrderTypeMutationVariables,
} from './__generated__/EditWorkOrderTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';
import type {WorkOrderType} from '../common/WorkOrder';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';
import {convertWorkOrderTypeToMutationInput} from '../common/WorkOrder';
import {getGraphError} from '../common/EntUtils';

export const mutation = graphql`
  mutation EditWorkOrderTypeMutation($input: EditWorkOrderTypeInput!) {
    editWorkOrderType(input: $input) {
      id
      name
      description
      ...AddEditWorkOrderTypeCard_workOrderType
    }
  }
`;

export const editWorkOrderType = (
  workOrderType: WorkOrderType,
): Promise<EditWorkOrderTypeMutationResponse> => {
  const variables: EditWorkOrderTypeMutationVariables = {
    input: {
      id: workOrderType.id,
      ...convertWorkOrderTypeToMutationInput(workOrderType),
    },
  };

  return new Promise((resolve, reject) => {
    const callbacks: MutationCallbacks<EditWorkOrderTypeMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          return reject(getGraphError(errors[0]));
        } else {
          resolve(response);
        }
      },
      onError: (error: Error) => reject(getGraphError(error)),
    };
    CommitEditWorkOrderTypeMutation(variables, callbacks);
  });
};

const CommitEditWorkOrderTypeMutation = (
  variables: EditWorkOrderTypeMutationVariables,
  callbacks?: MutationCallbacks<EditWorkOrderTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditWorkOrderTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};

export default CommitEditWorkOrderTypeMutation;
