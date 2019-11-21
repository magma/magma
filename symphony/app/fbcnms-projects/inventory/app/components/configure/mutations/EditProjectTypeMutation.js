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
  EditProjectTypeMutationResponse,
  EditProjectTypeMutationVariables,
} from './__generated__/EditProjectTypeMutation.graphql';
import type {MutationCallbacks} from '../../../mutations/MutationCallbacks.js';

const mutation = graphql`
  mutation EditProjectTypeMutation($input: EditProjectTypeInput!) {
    editProjectType(input: $input) {
      ...ProjectTypeCard_projectType
      ...AddEditProjectTypeCard_editingProjectType
    }
  }
`;

export default (
  variables: EditProjectTypeMutationVariables,
  callbacks?: MutationCallbacks<EditProjectTypeMutationResponse>,
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
