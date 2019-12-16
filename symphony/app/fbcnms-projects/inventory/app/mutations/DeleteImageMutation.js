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
  DeleteImageMutation,
  DeleteImageMutationResponse,
  DeleteImageMutationVariables,
} from './__generated__/DeleteImageMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation DeleteImageMutation(
    $entityType: ImageEntity!
    $entityId: ID!
    $id: ID!
  ) {
    deleteImage(entityType: $entityType, entityId: $entityId, id: $id) {
      ...EntityDocumentsTable_files @relay(mask: false)
      ...FileAttachment_file
    }
  }
`;

export default (
  variables: DeleteImageMutationVariables,
  callbacks?: MutationCallbacks<DeleteImageMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<DeleteImageMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
