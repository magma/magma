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
import {DefaultLinkWidget} from '@projectstorm/react-diagrams';
import {LinkContextMenu, LinkMenuProvider} from './ContextMenu';

export class LinkWithContextWidget extends DefaultLinkWidget {
  render() {
    return (
      <g>
        <LinkMenuProvider link={this.props.link}>
          {super.render()}
        </LinkMenuProvider>
        <LinkContextMenu
          link={this.props.link}
          diagramEngine={this.props.diagramEngine}
        />
      </g>
    );
  }
}
