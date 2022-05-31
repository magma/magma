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
 */

import {AxiosError, AxiosPromise} from 'axios';
import {BASE_API} from './MagmaAPI';
import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../app/hooks/useSnackbar';

export default function <TParams, TResponse>(
  func: (params: TParams) => AxiosPromise<TResponse>,
  params: TParams,
  onResponse?: (response: TResponse) => void,
  cacheCounter?: string | number,
): {
  response: TResponse | undefined;
  error?: Error;
  isLoading: boolean;
} {
  const [response, setResponse] = useState<TResponse>();
  const [error, setError] = useState<Error>();
  const [isLoading, setIsLoading] = useState(true);
  const jsonParams = JSON.stringify(params);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    func
      .call(BASE_API, JSON.parse(jsonParams) as TParams)
      .then(({data}) => {
        setResponse(data);
        setError(undefined);
        setIsLoading(false);
        onResponse && onResponse(data);
      })
      .catch((err: Error | AxiosError) => {
        setError(err);
        setResponse(undefined);
        setIsLoading(false);
        if ('response' in err && err.response?.status === 503) {
          enqueueSnackbar(
            'There was a problem connecting to the Orchestrator server',
            {variant: 'error'},
          );
        }
      });
  }, [jsonParams, func, cacheCounter, onResponse, enqueueSnackbar]);

  return {error, response, isLoading};
}
