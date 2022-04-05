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

import {useCallback, useContext, useEffect, useState} from 'react';
import {useRouteMatch} from 'react-router-dom';

// eslint-disable-next-line no-warning-comments
// $FlowFixMe - use react-router hooks
import {__RouterContext as RouterContext} from 'react-router-dom';

export const useRelativeUrl = () => {
  const {url} = useRouteMatch();
  return useCallback((path: string) => `${url}${path}`, [url]);
};

export const useRelativePath = () => {
  const match = useRouteMatch();
  return useCallback((path: string) => `${match.path}${path}`, [match.path]);
};

const useRouter = () => {
  const relativeUrl = useRelativeUrl();
  const relativePath = useRelativePath();
  const [, setUpdateCount] = useState(0);
  const routerContext = useContext(RouterContext);
  if (!routerContext) {
    throw new Error('You must be operating in a react-router context.');
  }

  const forceUpdate = () =>
    setUpdateCount(updateCount => (updateCount + 1) % 99999999999);

  useEffect(() => routerContext.history.listen(forceUpdate), [routerContext]);
  return {
    ...routerContext,
    relativeUrl,
    relativePath,
  };
};

export default useRouter;
