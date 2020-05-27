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
  EditProjectMutation,
  EditProjectMutationResponse,
  EditProjectMutationVariables,
} from './__generated__/EditProjectMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation EditProjectMutation($input: EditProjectInput!) {
    editProject(input: $input) {
      id
      name
      description
      createdBy {
        id
        email
      }
      properties {
        stringValue
        intValue
        floatValue
        booleanValue
        latitudeValue
        longitudeValue
        rangeFromValue
        rangeToValue
        propertyType {
          id
          name
          type
          nodeType
          isEditable
          isInstanceProperty
          stringValue
          intValue
          floatValue
          booleanValue
          latitudeValue
          longitudeValue
          rangeFromValue
          rangeToValue
        }
      }
    }
  }
`;

export default (
  variables: EditProjectMutationVariables,
  callbacks?: MutationCallbacks<EditProjectMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditProjectMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
