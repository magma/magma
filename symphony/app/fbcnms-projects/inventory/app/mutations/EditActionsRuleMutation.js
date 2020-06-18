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
  EditActionsRuleMutation,
  EditActionsRuleMutationResponse,
  EditActionsRuleMutationVariables,
} from './__generated__/EditActionsRuleMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation EditActionsRuleMutation($id: ID!, $input: AddActionsRuleInput!) {
    editActionsRule(id: $id, input: $input) {
      ...ActionsListCard_actionsRule
    }
  }
`;

export default (
  variables: EditActionsRuleMutationVariables,
  callbacks?: MutationCallbacks<EditActionsRuleMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditActionsRuleMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
