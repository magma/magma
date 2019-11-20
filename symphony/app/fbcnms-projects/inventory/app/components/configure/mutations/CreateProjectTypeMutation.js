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
  CreateProjectTypeMutationResponse,
  CreateProjectTypeMutationVariables,
} from './__generated__/CreateProjectTypeMutation.graphql';
import type {MutationCallbacks} from '../../../mutations/MutationCallbacks.js';

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
