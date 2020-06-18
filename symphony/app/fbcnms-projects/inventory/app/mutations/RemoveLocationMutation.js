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

import {getGraphError} from '../common/EntUtils';
import {removeLocationFromStore} from './utils/LocationStoreUtils';

const mutation = graphql`
  mutation RemoveLocationMutation($id: ID!) {
    removeLocation(id: $id)
  }
`;

export const removeLocation = (
  locationId: string,
  parentLocationId: ?string,
): Promise<RemoveLocationMutationResponse> => {
  const variables: RemoveLocationMutationVariables = {
    id: locationId,
  };

  return new Promise((resolve, reject) => {
    const callbacks: MutationCallbacks<RemoveLocationMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          return reject(getGraphError(errors[0]));
        } else {
          resolve(response);
        }
      },
      onError: (error: Error) => reject(getGraphError(error)),
    };
    const updater = store => {
      removeLocationFromStore(store, locationId, parentLocationId);
    };
    CommitRemoveLocationMutation(variables, callbacks, updater);
  });
};

const CommitRemoveLocationMutation = (
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

export default CommitRemoveLocationMutation;
