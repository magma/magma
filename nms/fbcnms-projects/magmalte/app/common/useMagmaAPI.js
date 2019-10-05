/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import MagmaV1API from './MagmaV1API';

import {useEffect, useState} from 'react';

export default function<TParams: {...}, TResponse>(
  func: TParams => Promise<TResponse>,
  params: TParams,
  cacheCounter?: string | number,
): {
  response: ?TResponse,
  // we can't really do better than this for now
  // eslint-disable-next-line flowtype/no-weak-types
  error: any,
  isLoading: boolean,
} {
  const [response, setResponse] = useState();
  const [error, setError] = useState();
  const [isLoading, setIsLoading] = useState(true);
  const jsonParams = JSON.stringify(params);

  useEffect(() => {
    func
      .bind(MagmaV1API)((JSON.parse(jsonParams): TParams))
      .then(res => {
        setResponse(res);
        setError(null);
        setIsLoading(false);
      })
      .catch(err => {
        setError(err);
        setResponse(null);
        setIsLoading(false);
      });
  }, [jsonParams, func, cacheCounter]);

  return {error, response, isLoading};
}
