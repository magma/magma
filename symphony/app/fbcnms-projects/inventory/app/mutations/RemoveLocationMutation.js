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
  RemoveLocationMutationResponse,
  RemoveLocationMutationVariables,
} from './__generated__/RemoveLocationMutation.graphql';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation RemoveLocationMutation($id: ID!) {
    removeLocation(id: $id)
  }
`;

export default (
  variables: RemoveLocationMutationVariables,
  callbacks?: MutationCallbacks<RemoveLocationMutationResponse>,
  updater?: SelectorStoreUpdater,
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
