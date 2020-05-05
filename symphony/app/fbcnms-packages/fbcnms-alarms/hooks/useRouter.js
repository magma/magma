/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {useCallback, useContext, useEffect, useState} from 'react';
import {useRouteMatch} from 'react-router';

// eslint-disable-next-line no-warning-comments
// $FlowFixMe - use react-router hooks
import {__RouterContext as RouterContext} from 'react-router';

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
