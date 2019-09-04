/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AppContext from '@fbcnms/ui/context/AppContext';
import NetworkCreate from './NetworkCreate';
import React, {useContext} from 'react';

import {Redirect, Route, Switch} from 'react-router-dom';
import {useRouter} from '@fbcnms/ui/hooks';

export default function Network() {
  const appContext = useContext(AppContext);
  const {relativePath} = useRouter();
  const path = relativePath('/create');
  return (
    <Switch>
      {appContext.user.isSuperUser ? (
        <Route path={path} component={NetworkCreate} />
      ) : null}
      <Redirect to={`/nms/`} />
    </Switch>
  );
}
