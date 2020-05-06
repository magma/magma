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
  AddServiceEndpointMutation,
  AddServiceEndpointMutationResponse,
  AddServiceEndpointMutationVariables,
} from './__generated__/AddServiceEndpointMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation AddServiceEndpointMutation($input: AddServiceEndpointInput!) {
    addServiceEndpoint(input: $input) {
      ...ServiceCard_service
    }
  }
`;

export default (
  variables: AddServiceEndpointMutationVariables,
  callbacks?: MutationCallbacks<AddServiceEndpointMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddServiceEndpointMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
