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
  DeleteReportFilterMutation,
  DeleteReportFilterMutationResponse,
  DeleteReportFilterMutationVariables,
} from './__generated__/DeleteReportFilterMutation.graphql';
import type {MutationCallbacks} from './MutationCallbacks.js';
import type {SelectorStoreUpdater} from 'relay-runtime';

const mutation = graphql`
  mutation DeleteReportFilterMutation($id: ID!) {
    deleteReportFilter(id: $id)
  }
`;

export default (
  variables: DeleteReportFilterMutationVariables,
  callbacks?: MutationCallbacks<DeleteReportFilterMutationResponse>,
  updater?: SelectorStoreUpdater,
) => {
  const {onCompleted, onError} = callbacks ? callbacks : {};
  commitMutation<DeleteReportFilterMutation>(RelayEnvironemnt, {
    mutation,
    variables,
    updater,
    onCompleted,
    onError,
  });
};
