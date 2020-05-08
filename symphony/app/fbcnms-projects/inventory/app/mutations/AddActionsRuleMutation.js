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
  AddActionsRuleMutation,
  AddActionsRuleMutationResponse,
  AddActionsRuleMutationVariables,
} from './__generated__/AddActionsRuleMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation AddActionsRuleMutation($input: AddActionsRuleInput!) {
    addActionsRule(input: $input) {
      id
      name
    }
  }
`;

export default (
  variables: AddActionsRuleMutationVariables,
  callbacks?: MutationCallbacks<AddActionsRuleMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddActionsRuleMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
