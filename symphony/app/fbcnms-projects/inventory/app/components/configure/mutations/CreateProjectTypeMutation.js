/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import RelayEnvironment from '../../../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';
import type {
  CreateProjectTypeMutation,
  CreateProjectTypeMutationResponse,
  CreateProjectTypeMutationVariables,
} from './__generated__/CreateProjectTypeMutation.graphql';
import type {MutationCallbacks} from '../../../mutations/MutationCallbacks.js';
import type {StoreUpdater} from '../../../common/RelayEnvironment';

const mutation = graphql`
  mutation CreateProjectTypeMutation($input: AddProjectTypeInput!) {
    createProjectType(input: $input) {
      ...ProjectTypeCard_projectType
      ...AddEditProjectTypeCard_editingProjectType
    }
  }
`;

export default (
  variables: CreateProjectTypeMutationVariables,
  callbacks?: MutationCallbacks<CreateProjectTypeMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<CreateProjectTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
