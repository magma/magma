/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
