/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import NetworkCreate from './NetworkCreate';
import React from 'react';

import {Route, Switch} from 'react-router-dom';
import {useRouter} from '@fbcnms/ui/hooks';

export default function Network() {
  const {relativePath} = useRouter();
  return (
    <Switch>
      <Route path={relativePath('/create')} component={NetworkCreate} />
    </Switch>
  );
}
