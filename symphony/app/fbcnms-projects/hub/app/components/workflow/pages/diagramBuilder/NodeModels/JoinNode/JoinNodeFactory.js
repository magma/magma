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
import {JoinNode} from './JoinNode';
import {JoinNodeModel} from './JoinNodeModel';

export class JoinNodeFactory extends SRD.AbstractNodeFactory {
  constructor() {
    super('join');
  }

  generateReactWidget(
    diagramEngine: SRD.DiagramEngine,
    node: SRD.NodeModel,
  ): JSX.Element {
    return <JoinNode node={node} diagramEngine={diagramEngine} />;
  }

  getNewInstance() {
    return new JoinNodeModel();
  }
}
