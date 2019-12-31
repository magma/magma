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
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {
  RemoveServiceEndpointMutation,
  RemoveServiceEndpointMutationResponse,
  RemoveServiceEndpointMutationVariables,
} from './__generated__/RemoveServiceEndpointMutation.graphql';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation RemoveServiceEndpointMutation($serviceEndpointId: ID!) {
    removeServiceEndpoint(serviceEndpointId: $serviceEndpointId) {
      ...ServiceCard_service
    }
  }
`;

export default (
  variables: RemoveServiceEndpointMutationVariables,
  callbacks?: MutationCallbacks<RemoveServiceEndpointMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveServiceEndpointMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
