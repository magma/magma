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
import {CircleNodeStart} from './CircleNodeStart';
import {CircleStartNodeModel} from './CircleStartNodeModel';

export class CircleStartNodeFactory extends SRD.AbstractNodeFactory {
  constructor() {
    super('start');
  }

  generateReactWidget(
    diagramEngine: SRD.DiagramEngine,
    node: SRD.NodeModel,
  ): JSX.Element {
    return <CircleNodeStart node={node} diagramEngine={diagramEngine} />;
  }

  getNewInstance() {
    return new CircleStartNodeModel();
  }
}
