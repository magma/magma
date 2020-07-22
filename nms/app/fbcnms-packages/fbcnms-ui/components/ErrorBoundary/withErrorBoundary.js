/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import ErrorBoundary from './ErrorBoundary';

export default function withErrorBoundary<TComponent: React.ComponentType<*>>(
  Component: TComponent,
  onError: ?() => void = null,
): React.ComponentType<React.ElementConfig<TComponent>> {
  return class extends React.Component<React.ElementConfig<TComponent>> {
    render(): React.Node {
      return (
        <ErrorBoundary onError={onError}>
          <Component {...this.props} />
        </ErrorBoundary>
      );
    }
  };
}
