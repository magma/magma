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
  AddLocationMutation,
  AddLocationMutationResponse,
  AddLocationMutationVariables,
} from './__generated__/AddLocationMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation AddLocationMutation($input: AddLocationInput!) {
    addLocation(input: $input) {
      ...LocationsTree_location @relay(mask: false)
      children {
        ...LocationsTree_location @relay(mask: false)
      }
    }
  }
`;

export default (
  variables: AddLocationMutationVariables,
  callbacks?: MutationCallbacks<AddLocationMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddLocationMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
