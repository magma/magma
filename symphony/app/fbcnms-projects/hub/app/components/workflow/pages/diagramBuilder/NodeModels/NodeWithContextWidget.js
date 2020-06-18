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
import {DefaultNodeWidget} from '@projectstorm/react-diagrams';
import {NodeContextMenu, NodeMenuProvider} from './ContextMenu';

export class NodeWithContextWidget extends DefaultNodeWidget {
  render() {
    return (
      <div {...this.getProps()}>
        <NodeMenuProvider node={this.props.node}>
          {super.render()}
        </NodeMenuProvider>
        <NodeContextMenu
          node={this.props.node}
          diagramEngine={this.props.diagramEngine}
        />
      </div>
    );
  }
}
