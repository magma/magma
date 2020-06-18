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
  RemoveActionsRuleMutation,
  RemoveActionsRuleMutationResponse,
  RemoveActionsRuleMutationVariables,
} from './__generated__/RemoveActionsRuleMutation.graphql';
import type {SelectorStoreUpdater} from 'relay-runtime';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

const mutation = graphql`
  mutation RemoveActionsRuleMutation($id: ID!) {
    removeActionsRule(id: $id)
  }
`;

export default (
  variables: RemoveActionsRuleMutationVariables,
  callbacks?: MutationCallbacks<RemoveActionsRuleMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveActionsRuleMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
