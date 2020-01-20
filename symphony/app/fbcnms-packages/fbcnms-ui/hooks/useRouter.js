/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {useContext, useEffect, useState} from 'react';

import {__RouterContext as RouterContext} from 'react-router';

const useRouter = () => {
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
    relativeUrl: (path: string) => `${routerContext.match.url}${path}`,
    relativePath: (path: string) => `${routerContext.match.path}${path}`,
  };
};

export default useRouter;
