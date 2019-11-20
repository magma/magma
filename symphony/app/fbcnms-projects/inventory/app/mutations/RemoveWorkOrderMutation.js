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
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {
  RemoveWorkOrderMutationMutationResponse,
  RemoveWorkOrderMutationMutationVariables,
} from './__generated__/RemoveWorkOrderMutation.graphql';

const mutation = graphql`
  mutation RemoveWorkOrderMutation($id: ID!) {
    removeWorkOrder(id: $id)
  }
`;

export default (
  variables: RemoveWorkOrderMutationMutationVariables,
  callbacks?: MutationCallbacks<RemoveWorkOrderMutationMutationResponse>,
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
