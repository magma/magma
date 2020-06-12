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
  EditLocationTypesIndexMutation,
  EditLocationTypesIndexMutationResponse,
  EditLocationTypesIndexMutationVariables,
} from './__generated__/EditLocationTypesIndexMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

import RelayEnvironment from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';
import {getGraphError} from '../common/EntUtils';

export const mutation = graphql`
  mutation EditLocationTypesIndexMutation(
    $locationTypeIndex: [LocationTypeIndex]!
  ) {
    editLocationTypesIndex(locationTypesIndex: $locationTypeIndex) {
      id
      name
      index
      ...LocationTypeItem_locationType
      ...AddEditLocationTypeCard_editingLocationType
    }
  }
`;

export const saveLocationTypeIndexes = (
  input: $ReadOnlyArray<{|
    locationTypeID: string,
    index: number,
  |}>,
): Promise<EditLocationTypesIndexMutationResponse> => {
  const variables: EditLocationTypesIndexMutationVariables = {
    locationTypeIndex: input,
  };

  return new Promise((resolve, reject) => {
    const callbacks: MutationCallbacks<EditLocationTypesIndexMutationResponse> = {
      onCompleted: (response, errors) => {
        if (errors && errors[0]) {
          return reject(getGraphError(errors[0]));
        } else {
          resolve(response);
        }
      },
      onError: (error: Error) => reject(getGraphError(error)),
    };
    CommitEditLocationTypesIndexMutation(variables, callbacks);
  });
};

const CommitEditLocationTypesIndexMutation = (
  variables: EditLocationTypesIndexMutationVariables,
  callbacks?: MutationCallbacks<EditLocationTypesIndexMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditLocationTypesIndexMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};

export default CommitEditLocationTypesIndexMutation;
