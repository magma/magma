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
  EditServiceTypeMutation,
  EditServiceTypeMutationResponse,
  EditServiceTypeMutationVariables,
} from './__generated__/EditServiceTypeMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation EditServiceTypeMutation($data: ServiceTypeEditData!) {
    editServiceType(data: $data) {
      id
      name
      propertyTypes {
        ...PropertyTypeFormField_propertyType @relay(mask: false)
      }
    }
  }
`;

export default (
  variables: EditServiceTypeMutationVariables,
  callbacks?: MutationCallbacks<EditServiceTypeMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditServiceTypeMutation>(RelayEnvironment, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
