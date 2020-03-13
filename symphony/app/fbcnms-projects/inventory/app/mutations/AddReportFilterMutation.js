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
  AddReportFilterMutation,
  AddReportFilterMutationResponse,
  AddReportFilterMutationVariables,
} from './__generated__/AddReportFilterMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {StoreUpdater} from '../common/RelayEnvironment';

const mutation = graphql`
  mutation AddReportFilterMutation($input: ReportFilterInput!) {
    addReportFilter(input: $input) {
      id
      name
      entity
      filters {
        filterType
        operator
        stringValue
        idSet
        stringSet
        boolValue
        propertyValue {
          id
          name
          type
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
  variables: AddReportFilterMutationVariables,
  callbacks?: MutationCallbacks<AddReportFilterMutationResponse>,
  updater?: StoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<AddReportFilterMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
