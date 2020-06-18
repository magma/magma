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
  AddWorkOrderTypeMutation,
  AddWorkOrderTypeMutationResponse,
  AddWorkOrderTypeMutationVariables,
} from './__generated__/AddWorkOrderTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';
import type {WorkOrderType} from '../common/WorkOrder';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {ConnectionHandler} from 'relay-runtime';
import {commitMutation, graphql} from 'react-relay';
import {convertWorkOrderTypeToMutationInput} from '../common/WorkOrder';
import {getGraphError} from '../common/EntUtils';

const mutation = graphql`
  mutation AddWorkOrderTypeMutation($input: AddWorkOrderTypeInput!) {
    addWorkOrderType(input: $input) {
      id
      name
      description
      ...AddEditWorkOrderTypeCard_workOrderType
    }
  }
`;

export const addWorkOrderType = (
  workOrderType: WorkOrderType,
): Promise<AddWorkOrderTypeMutationResponse> => {
  const variables: AddWorkOrderTypeMutationVariables = {
    input: convertWorkOrderTypeToMutationInput(workOrderType),
  };

  return new Promise((resolve, reject) => {
    const callbacks: MutationCallbacks<AddWorkOrderTypeMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          return reject(getGraphError(errors[0]));
        } else {
          resolve(response);
        }
      },
      onError: (error: Error) => reject(getGraphError(error)),
    };
    const updater = store => {
      const rootQuery = store.getRoot();
      const newNode = store.getRootField('addWorkOrderType');
      if (!newNode) {
        return;
      }
      const types = ConnectionHandler.getConnection(
        rootQuery,
        'Configure_workOrderTypes',
      );
      if (types == null) {
        return;
      }
      const edge = ConnectionHandler.createEdge(
        store,
        types,
        newNode,
        'WorkOrderTypesEdge',
      );
      ConnectionHandler.insertEdgeBefore(types, edge);
    };
    CommitWorkOrderTypeMutation(variables, callbacks, updater);
  });
};

const CommitWorkOrderTypeMutation = (
  variables: AddWorkOrderTypeMutationVariables,
  callbacks?: MutationCallbacks<AddWorkOrderTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddWorkOrderTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};

export default CommitWorkOrderTypeMutation;
