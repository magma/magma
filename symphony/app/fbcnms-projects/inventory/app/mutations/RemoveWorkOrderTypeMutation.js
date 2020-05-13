/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MutationCallbacks} from './MutationCallbacks.js';
import type {
  RemoveWorkOrderTypeMutation,
  RemoveWorkOrderTypeMutationResponse,
  RemoveWorkOrderTypeMutationVariables,
} from './__generated__/RemoveWorkOrderTypeMutation.graphql';
import type {SelectorStoreUpdater} from 'relay-runtime';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {ConnectionHandler} from 'relay-runtime';
import {commitMutation, graphql} from 'react-relay';

const mutation = graphql`
  mutation RemoveWorkOrderTypeMutation($id: ID!) {
    removeWorkOrderType(id: $id)
  }
`;

export const deleteWorkOrderType = (workOrderTypeId: string): Promise<void> => {
  return new Promise((resolve, reject) => {
    CommitRemoveWorkOrderTypeMutation(
      {
        id: workOrderTypeId,
      },
      {
        onCompleted: (response, errors) => {
          if (errors && errors[0]) {
            return reject(errors[0]);
          }
          resolve();
        },
        onError: reject,
      },
      store => {
        const rootQuery = store.getRoot();
        const workOrderTypes = ConnectionHandler.getConnection(
          rootQuery,
          'Configure_workOrderTypes',
        );
        if (workOrderTypes != null) {
          ConnectionHandler.deleteNode(workOrderTypes, workOrderTypeId);
        }
        store.delete(workOrderTypeId);
      },
    );
  });
};

const CommitRemoveWorkOrderTypeMutation = (
  variables: RemoveWorkOrderTypeMutationVariables,
  callbacks?: MutationCallbacks<RemoveWorkOrderTypeMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveWorkOrderTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};

export default CommitRemoveWorkOrderTypeMutation;
