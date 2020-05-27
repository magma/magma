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
import {DiagramEngine} from '@projectstorm/react-diagrams';
import {JoinNodePortModel} from './JoinNodePortModel';
import {NodeModel} from '@projectstorm/react-diagrams';

export class JoinNodeModel extends NodeModel {
  name: string;
  color: string;
  inputs: {};

  constructor(
    name: string = 'Untitled',
    color: string = 'rgb(0,192,255)',
    inputs: {},
  ) {
    super('join');
    this.name = name;
    this.color = color;
    super.extras = {inputs: inputs};

    this.addPort(new JoinNodePortModel(true, 'left'));
    this.addPort(new JoinNodePortModel(false, 'right'));
  }

  deSerialize(object, engine: DiagramEngine) {
    super.deSerialize(object, engine);
    this.name = object.name;
    this.color = object.color;
  }

  serialize() {
    return _.merge(super.serialize(), {
      name: this.name,
      color: this.color,
    });
  }

  getInputs() {
    return this.inputs;
  }
}
