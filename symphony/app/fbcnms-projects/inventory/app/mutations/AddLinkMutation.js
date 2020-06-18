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
  AddLinkMutation,
  AddLinkMutationResponse,
  AddLinkMutationVariables,
} from './__generated__/AddLinkMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation AddLinkMutation($input: AddLinkInput!) {
    addLink(input: $input) {
      ...EquipmentPortsTable_link @relay(mask: false)
    }
  }
`;

export default (
  variables: AddLinkMutationVariables,
  callbacks?: MutationCallbacks<AddLinkMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddLinkMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
