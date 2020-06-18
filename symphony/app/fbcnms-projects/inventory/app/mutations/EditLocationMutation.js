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
  EditLocationMutation,
  EditLocationMutationResponse,
  EditLocationMutationVariables,
} from './__generated__/EditLocationMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation EditLocationMutation($input: EditLocationInput!) {
    editLocation(input: $input) {
      ...LocationsTree_location @relay(mask: false)
    }
  }
`;

export default (
  variables: EditLocationMutationVariables,
  callbacks?: MutationCallbacks<EditLocationMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditLocationMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
