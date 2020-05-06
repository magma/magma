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
  AddServiceLinkMutation,
  AddServiceLinkMutationResponse,
  AddServiceLinkMutationVariables,
} from './__generated__/AddServiceLinkMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation AddServiceLinkMutation($id: ID!, $linkId: ID!) {
    addServiceLink(id: $id, linkId: $linkId) {
      ...ServiceCard_service
    }
  }
`;

export default (
  variables: AddServiceLinkMutationVariables,
  callbacks?: MutationCallbacks<AddServiceLinkMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddServiceLinkMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
