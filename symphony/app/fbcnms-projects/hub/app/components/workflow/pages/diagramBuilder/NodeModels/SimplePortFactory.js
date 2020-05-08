/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import {AbstractPortFactory, PortModel} from '@projectstorm/react-diagrams';

export class SimplePortFactory extends AbstractPortFactory {
  cb: initialConfig => PortModel;

  constructor(type: string, cb: initialConfig => PortModel) {
    super(type);
    this.cb = cb;
  }

  getNewInstance(initialConfig): PortModel {
    return this.cb(initialConfig);
  }
}
