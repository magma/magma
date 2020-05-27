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
import {DecisionNodePortModel} from './DecisionNodePortModel';
import {DiagramEngine} from '@projectstorm/react-diagrams';
import {NodeModel} from '@projectstorm/react-diagrams';

export class DecisionNodeModel extends NodeModel {
  name: string;
  color: string;
  inputs: {};

  constructor(
    name: string = 'Untitled',
    color: string = 'rgb(0,192,255)',
    inputs: {},
  ) {
    super('decision');
    this.name = name;
    this.color = color;
    super.extras = {inputs: inputs};

    this.addPort(new DecisionNodePortModel(true, 'left', 'inputPort'));
    this.addPort(new DecisionNodePortModel(false, 'right', 'neutralPort'));
    this.addPort(new DecisionNodePortModel(false, 'bottom', 'failPort'));
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
