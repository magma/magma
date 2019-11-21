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
  AddLocationMutationResponse,
  AddLocationMutationVariables,
} from './__generated__/AddLocationMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';

import RelayEnvironemnt from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

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
  updater?: (store: any) => void,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
