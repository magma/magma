/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import axios from 'axios';
import {DEACTIVATED_PAGE_PATH} from '../components/DeactivatedPage';
import {Environment, Network, RecordSource, Store} from 'relay-runtime';

function handleDeactivatedUser(error) {
  const errorResponse = error?.response;
  if (
    errorResponse != null &&
    errorResponse.status === 403 &&
    typeof errorResponse.data === 'string' &&
    errorResponse.data.includes('user is deactivated')
  ) {
    window.location.replace(DEACTIVATED_PAGE_PATH);
  }

  throw error;
}

function fetchQuery(operation, variables) {
  return axios
    .post('/graph/query', {
      query: operation.text,
      variables,
    })
    .then(response => {
      return response.data;
    })
    .catch(handleDeactivatedUser);
}

const RelayEnvironment = new Environment({
  network: Network.create(fetchQuery),
  store: new Store(new RecordSource()),
});

export default RelayEnvironment;
