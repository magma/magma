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
import InventoryErrorBoundary from './InventoryErrorBoundary';

export default function withInventoryErrorBoundary<
  TComponent: React.ComponentType<*>,
>(Component: TComponent): React.ComponentType<React.ElementConfig<TComponent>> {
  return class extends React.Component<React.ElementConfig<TComponent>> {
    render(): React.Node {
      return (
        <InventoryErrorBoundary>
          <Component {...this.props} />
        </InventoryErrorBoundary>
      );
    }
  };
}
