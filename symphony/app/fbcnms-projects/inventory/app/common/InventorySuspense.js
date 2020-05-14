/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import LoadingIndicator from './LoadingIndicator';
import {Suspense} from 'react';

type Props = $ReadOnly<{|
  children: React.Node,
|}>;

const InventorySuspense = ({children}: Props) => {
  return <Suspense fallback={<LoadingIndicator />}>{children}</Suspense>;
};

export default InventorySuspense;
