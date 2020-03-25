/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import AppContext from '@fbcnms/ui/context/AppContext';
import React from 'react';
import Settings from '@fbcnms/magmalte/app/components/Settings';

import {useContext} from 'react';

export default function() {
  const {user} = useContext(AppContext);
  return <Settings isSuperUser={user?.isSuperUser} />;
}
