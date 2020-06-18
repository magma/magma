/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import * as _ from 'lodash';
import {
  DefaultLinkModel,
  DiagramEngine,
  LinkModel,
  PortModel,
} from '@projectstorm/react-diagrams';

export class ForkNodePortModel extends PortModel {
  position: string | 'top' | 'bottom' | 'left' | 'right';

  constructor(isInput: boolean, pos: string = 'left') {
    super(pos, 'fork');
    this.in = isInput;
    this.position = pos;
  }

  serialize() {
    return _.merge(super.serialize(), {
      position: this.position,
    });
  }

  link(port: PortModel): LinkModel {
    const link = this.createLinkModel();
    link.setSourcePort(this);
    link.setTargetPort(port);
    return link;
  }

  deSerialize(data, engine: DiagramEngine) {
    super.deSerialize(data, engine);
    this.position = data.position;
  }

  canLinkToPort(port: PortModel): boolean {
    return !this.in && port.in;
  }

  createLinkModel(): LinkModel {
    return new DefaultLinkModel();
  }
}
