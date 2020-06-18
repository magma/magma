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
import RelayEnvironment from './RelayEnvironment';
import {RelayEnvironmentProvider} from 'react-relay/hooks';
import {Suspense} from 'react';

type Props = $ReadOnly<{|
  children: React.Node,
  isTopLevel?: ?boolean,
|}>;

const InventorySuspense = (props: Props) => {
  const {children, isTopLevel} = props;
  const suspense = (
    <Suspense fallback={<LoadingIndicator />}>{children}</Suspense>
  );
  if (isTopLevel) {
    return (
      <RelayEnvironmentProvider environment={RelayEnvironment}>
        {suspense}
      </RelayEnvironmentProvider>
    );
  }
  return suspense;
};

export default InventorySuspense;
