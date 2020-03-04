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
  RemoveLocationMutation,
  // $FlowFixMe (T62907961) Relay flow types
  RemoveLocationMutationMutationResponse,
  // $FlowFixMe (T62907961) Relay flow types
  RemoveLocationMutationMutationVariables,
} from './__generated__/RemoveLocationMutation.graphql';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation RemoveLocationMutation($id: ID!) {
    removeLocation(id: $id)
  }
`;

export default (
  variables: RemoveLocationMutationMutationVariables,
  callbacks?: MutationCallbacks<RemoveLocationMutationMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<RemoveLocationMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
