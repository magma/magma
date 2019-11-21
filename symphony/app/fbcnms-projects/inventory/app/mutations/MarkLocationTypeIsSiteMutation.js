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
  MarkLocationTypeIsSiteMutationMutationResponse,
  MarkLocationTypeIsSiteMutationVariables,
} from './__generated__/MarkLocationTypeIsSiteMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

const mutation = graphql`
  mutation MarkLocationTypeIsSiteMutation($id: ID!, $isSite: Boolean!) {
    markLocationTypeIsSite(id: $id, isSite: $isSite) {
      id
      isSite
    }
  }
`;

export default (
  variables: MarkLocationTypeIsSiteMutationVariables,
  callbacks?: MutationCallbacks<MarkLocationTypeIsSiteMutationMutationResponse>,
  updater?: (store: any) => void,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
