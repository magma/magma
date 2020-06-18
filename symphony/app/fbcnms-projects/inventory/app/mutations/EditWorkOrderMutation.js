/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import RelayEnvironemnt from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';
import type {
  EditWorkOrderMutation,
  EditWorkOrderMutationResponse,
  EditWorkOrderMutationVariables,
} from './__generated__/EditWorkOrderMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation EditWorkOrderMutation($input: EditWorkOrderInput!) {
    editWorkOrder(input: $input) {
      id
      name
      description
      owner {
        id
        email
      }
      creationDate
      installDate
      status
      priority
      assignedTo {
        id
        email
      }
      ...WorkOrderDetails_workOrder
      ...WorkOrdersView_workOrder
    }
  }
`;

export default (
  variables: EditWorkOrderMutationVariables,
  callbacks?: MutationCallbacks<EditWorkOrderMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditWorkOrderMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
