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
  AddProjectMutation,
  AddProjectMutationResponse,
  AddProjectMutationVariables,
} from './__generated__/AddProjectMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation AddProjectMutation($input: AddProjectInput!) {
    createProject(input: $input) {
      ...ProjectsTableView_projects
    }
  }
`;

export default (
  variables: AddProjectMutationVariables,
  callbacks?: MutationCallbacks<AddProjectMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddProjectMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
