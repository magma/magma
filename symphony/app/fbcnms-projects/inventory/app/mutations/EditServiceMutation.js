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
  EditServiceMutation,
  EditServiceMutationResponse,
  EditServiceMutationVariables,
} from './__generated__/EditServiceMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation EditServiceMutation($data: ServiceEditData!) {
    editService(data: $data) {
      ...ServiceCard_service
    }
  }
`;

export default (
  variables: EditServiceMutationVariables,
  callbacks?: MutationCallbacks<EditServiceMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditServiceMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
