/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MutationCallbacks} from './MutationCallbacks.js';
import type {
  RemoveServiceLinkMutationResponse,
  RemoveServiceLinkMutationVariables,
} from './__generated__/RemoveServiceLinkMutation.graphql';

import RelayEnvironemnt from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

const mutation = graphql`
  mutation RemoveServiceLinkMutation($id: ID!, $linkId: ID!) {
    removeServiceLink(id: $id, linkId: $linkId) {
      ...ServiceCard_service
    }
  }
`;

export default (
  variables: RemoveServiceLinkMutationVariables,
  callbacks?: MutationCallbacks<RemoveServiceLinkMutationResponse>,
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
