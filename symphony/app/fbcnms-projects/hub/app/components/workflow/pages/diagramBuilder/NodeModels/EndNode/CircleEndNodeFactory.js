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
import * as SRD from '@projectstorm/react-diagrams';
import {CircleEndNodeModel} from './CircleEndNodeModel';
import {CircleNodeEnd} from './CircleNodeEnd';

export class CircleEndNodeFactory extends SRD.AbstractNodeFactory {
  constructor() {
    super('end');
  }

  generateReactWidget(
    diagramEngine: SRD.DiagramEngine,
    node: SRD.NodeModel,
  ): JSX.Element {
    return <CircleNodeEnd node={node} diagramEngine={diagramEngine} />;
  }

  getNewInstance() {
    return new CircleEndNodeModel();
  }
}
