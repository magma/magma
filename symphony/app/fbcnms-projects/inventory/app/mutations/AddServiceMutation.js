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
  AddServiceMutation,
  AddServiceMutationResponse,
  AddServiceMutationVariables,
} from './__generated__/AddServiceMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation AddServiceMutation($data: ServiceCreateData!) {
    addService(data: $data) {
      id
      ...ServicesView_service
    }
  }
`;

export default (
  variables: AddServiceMutationVariables,
  callbacks?: MutationCallbacks<AddServiceMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddServiceMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
