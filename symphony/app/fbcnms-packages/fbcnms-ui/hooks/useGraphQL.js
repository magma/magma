/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {Environment, fetchQuery} from 'relay-runtime';
import {useEffect, useState} from 'react';

export default function(
  env: Environment,
  query: any,
  variables: {[string]: mixed},
) {
  const [error, setError] = useState(null);
  const [response, setResponse] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const variablesJSON = JSON.stringify(variables);
  useEffect(() => {
    const variables = JSON.parse(variablesJSON);

    setError(null);
    setIsLoading(true);
    fetchQuery(env, query, variables)
      .then(response => {
        setResponse(response);
        setIsLoading(false);
      })
      .catch(error => {
        setError(error);
        setIsLoading(false);
      });
  }, [env, query, variablesJSON]);

  return {error, response, isLoading};
}
