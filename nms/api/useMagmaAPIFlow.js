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
 * @flow strict-local
 * @format
 */

import MagmaV1API from '../generated/WebClient';

import {useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../app/hooks/useSnackbar';

export default function <TParams: {...}, TResponse>(
  func: TParams => Promise<TResponse>,
  params: TParams,
  onResponse?: TResponse => void,
  cacheCounter?: string | number,
): {
  response: ?TResponse,
  // we can't really do better than this for now
  // eslint-disable-next-line flowtype/no-weak-types
  // $FlowFixMe[unclear-type] TODO(andreilee): migrated from fbcnms-ui
  error: any,
  isLoading: boolean,
} {
  const [response, setResponse] = useState();
  const [error, setError] = useState();
  const [isLoading, setIsLoading] = useState(true);
  const jsonParams = JSON.stringify(params);
  const enqueueSnackbar = useEnqueueSnackbar();

  useEffect(() => {
    func
      .bind(MagmaV1API)((JSON.parse(jsonParams): TParams))
      .then(res => {
        setResponse(res);
        setError(null);
        setIsLoading(false);
        onResponse && onResponse(res);
      })
      .catch(err => {
        setError(err);
        setResponse(null);
        setIsLoading(false);
        if (err?.response?.status === 503) {
          enqueueSnackbar(
            'There was a problem connecting to the Orchestrator server',
            {variant: 'error'},
          );
        }
      });
  }, [jsonParams, func, cacheCounter, onResponse, enqueueSnackbar]);

  return {error, response, isLoading};
}
