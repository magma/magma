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
import {DefaultLinkFactory} from '@projectstorm/react-diagrams';
import {LinkWithContextWidget} from './LinkWithContextWidget';

export class LinkWithContextFactory extends DefaultLinkFactory {
  generateReactWidget(diagramEngine, link) {
    return React.createElement(LinkWithContextWidget, {
      diagramEngine: diagramEngine,
      link: link,
    });
  }
}
