/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {$AxiosXHR, AxiosXHRConfig} from 'axios';

import axios from 'axios';
import {merge} from 'lodash';
import {useEffect, useState} from 'react';

type AxiosResponse<T, R> = {
  error: any,
  isLoading: boolean,
  response: ?$AxiosXHR<T, R>,
  loadedUrl: ?string,
};

export default function useAxios<T, R>(
  config: {onResponse?: ($AxiosXHR<T, R>) => void} & AxiosXHRConfig<T, R>,
): AxiosResponse<T, R> {
  const [error, setError] = useState(null);
  const [response, setResponse] = useState(null);
  const [isLoading, setIsLoading] = useState(false);
  const [loadedUrl, setLoadedUrl] = useState(null);

  // implicitly filters out functions, e.g. onResponse
  const stringConfigs = JSON.stringify(config);
  const onResponse = config.onResponse;

  useEffect(() => {
    const requestConfigs = JSON.parse(stringConfigs);
    const source = axios.CancelToken.source();
    const configWithCancelToken = merge({}, requestConfigs, {
      cancelToken: source.token,
    });
    setIsLoading(true);
    setError(null);
    axios
      .request<T, R>(configWithCancelToken)
      .then(res => {
        setIsLoading(false);
        setResponse(res);
        onResponse && onResponse(res);
        setLoadedUrl(requestConfigs.url);
      })
      .catch(error => {
        if (!axios.isCancel(error)) {
          setIsLoading(false);
          setError(error);
          setLoadedUrl(requestConfigs.url);
        }
      });
    return () => {
      source.cancel();
      setIsLoading(false);
      setLoadedUrl(requestConfigs.url);
    };
  }, [onResponse, stringConfigs]);
  return {
    error,
    isLoading,
    response,
    loadedUrl,
  };
}
