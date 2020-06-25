/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  MoveLocationMutation,
  MoveLocationMutationResponse,
  MoveLocationMutationVariables,
} from './__generated__/MoveLocationMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {
  addLocationToStore,
  removeLocationFromStore,
} from './utils/LocationStoreUtils';
import {commitMutation, graphql} from 'react-relay';
import {getGraphError} from '../common/EntUtils';

const mutation = graphql`
  mutation MoveLocationMutation($locationID: ID!, $parentLocationID: ID) {
    moveLocation(locationID: $locationID, parentLocationID: $parentLocationID) {
      ...LocationsTree_location @relay(mask: false)
    }
  }
`;

export const moveLocation = (
  locationId: string,
  parentLocationId: ?string,
  targetLocationId: ?string,
): Promise<MoveLocationMutationResponse> => {
  const variables: MoveLocationMutationVariables = {
    locationID: locationId,
    parentLocationID: targetLocationId,
  };

  return new Promise((resolve, reject) => {
    const callbacks: MutationCallbacks<MoveLocationMutationResponse> = {
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
      const newNode = store.getRootField('moveLocation');
      if (newNode == null) {
        return;
      }
      removeLocationFromStore(store, locationId, parentLocationId);
      addLocationToStore(store, newNode, targetLocationId);
    };
    CommitMoveLocationMutation(variables, callbacks, updater);
  });
};

const CommitMoveLocationMutation = (
  variables: MoveLocationMutationVariables,
  callbacks?: MutationCallbacks<MoveLocationMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<MoveLocationMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};

export default CommitMoveLocationMutation;
