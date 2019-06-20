/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AxiosXHRConfig, $AxiosXHR} from 'axios';

import {useEffect, useState} from 'react';
import axios from 'axios';
import {merge} from 'lodash';

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

  const stringConfig = JSON.stringify(config);

  useEffect(() => {
    const source = axios.CancelToken.source();
    const configWithCancelToken = merge({}, config, {
      cancelToken: source.token,
    });
    setIsLoading(true);
    setError(null);
    axios
      .request<T, R>(configWithCancelToken)
      .then(res => {
        setIsLoading(false);
        setResponse(res);
        config.onResponse && config.onResponse(res);
        setLoadedUrl(config.url);
      })
      .catch(error => {
        if (!axios.isCancel(error)) {
          setIsLoading(false);
          setError(error);
          setLoadedUrl(config.url);
        }
      });
    return () => {
      source.cancel();
      setIsLoading(false);
      setLoadedUrl(config.url);
    };
  }, [stringConfig]);
  return {
    error,
    isLoading,
    response,
    loadedUrl,
  };
}
