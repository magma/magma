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
import {NodeContextMenu, NodeMenuProvider} from '../ContextMenu';
import {PortWidget} from '@projectstorm/react-diagrams';

export class CircleNodeEnd extends React.Component {
  render() {
    return (
      <div className={'srd-circle-node'}>
        <NodeMenuProvider node={this.props.node}>
          <svg width="60" height="60">
            <g>
              <circle cx="30" cy="30" r="30" fill="white" />
              <text x="17" y="35">
                End
              </text>
            </g>
          </svg>
        </NodeMenuProvider>
        <div style={{position: 'absolute', zIndex: 10, left: -10, top: 21}}>
          <PortWidget name="left" node={this.props.node} />
        </div>
        <NodeContextMenu
          node={this.props.node}
          diagramEngine={this.props.diagramEngine}
        />
      </div>
    );
  }
}
