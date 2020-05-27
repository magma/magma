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
import {DefaultNodeFactory} from '@projectstorm/react-diagrams';
import {NodeWithContextWidget} from './NodeWithContextWidget';

export class NodeWithContextFactory extends DefaultNodeFactory {
  generateReactWidget(diagramEngine, node) {
    return React.createElement(NodeWithContextWidget, {
      node: node,
      diagramEngine: diagramEngine,
    });
  }
}
