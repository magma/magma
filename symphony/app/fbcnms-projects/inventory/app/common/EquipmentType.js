/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {PropertyType} from './PropertyType';

export type PositionDefinition = {
  id: string,
  name: string,
  index?: ?number,
  visibleLabel?: ?string,
};

export type EquipmentPortType = {
  id: string,
  name: string,
  propertyTypes: Array<PropertyType>,
  linkPropertyTypes: Array<PropertyType>,
  numberOfPortDefinitions: number,
};

export type PortDefinition = {
  id: string,
  name: string,
  index: number,
  visibleLabel?: ?string,
  portType: ?EquipmentPortType,
  bandwidth?: ?string,
};

export type EquipmentType = {
  id: string,
  name: string,
  positionDefinitions: Array<PositionDefinition>,
  portDefinitions: Array<PortDefinition>,
  propertyTypes: Array<PropertyType>,
  numberOfEquipment: number,
};
