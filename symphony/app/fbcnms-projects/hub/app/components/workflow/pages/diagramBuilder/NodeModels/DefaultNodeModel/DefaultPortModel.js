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
import {DefaultLinkModel} from '@projectstorm/react-diagrams';
import {DiagramEngine} from '@projectstorm/react-diagrams';
import {LinkModel} from '@projectstorm/react-diagrams';
import {PortModel} from '@projectstorm/react-diagrams';

export class DefaultPortModel extends PortModel {
  in: boolean;
  label: string;

  constructor(
    isInput: boolean,
    name: string,
    label: string = null,
    id?: string,
  ) {
    super(name, 'default', id);
    this.in = isInput;
    this.label = label || name;
  }

  deSerialize(object, engine: DiagramEngine) {
    super.deSerialize(object, engine);
    this.in = object.in;
    this.label = object.label;
  }

  serialize() {
    return _.merge(super.serialize(), {
      in: this.in,
      label: this.label,
    });
  }

  link(port: PortModel): LinkModel {
    const link = this.createLinkModel();
    link.setSourcePort(this);
    link.setTargetPort(port);
    return link;
  }

  canLinkToPort(port: PortModel): boolean {
    return !this.in && port.in;
  }

  createLinkModel(): LinkModel {
    const link = super.createLinkModel();
    return link || new DefaultLinkModel();
  }
}
