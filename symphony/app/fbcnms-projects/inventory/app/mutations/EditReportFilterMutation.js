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
  EditReportFilterMutation,
  EditReportFilterMutationResponse,
  EditReportFilterMutationVariables,
} from './__generated__/EditReportFilterMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

import RelayEnvironemnt from '../common/RelayEnvironment.js';
import {commitMutation, graphql} from 'react-relay';

const mutation = graphql`
  mutation EditReportFilterMutation($input: EditReportFilterInput!) {
    editReportFilter(input: $input) {
      id
      name
      entity
      filters {
        filterType
        key
        operator
        stringValue
        idSet
        stringSet
        boolValue
        propertyValue {
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
  variables: EditReportFilterMutationVariables,
  callbacks?: MutationCallbacks<EditReportFilterMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<EditReportFilterMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
