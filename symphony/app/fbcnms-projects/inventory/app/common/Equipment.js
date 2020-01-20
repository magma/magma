/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  EquipmentType,
  PortDefinition,
  PositionDefinition,
} from './EquipmentType';
import type {FutureState, WorkOrder} from './WorkOrder';
import type {Location} from './Location';
import type {Property} from './Property';
import type {Service, ServiceEndpoint} from './Service';

import {getInitialPropertyFromType} from './PropertyType';

export type Equipment = {
  id: string,
  name: string,
  equipmentType: EquipmentType,
  parentLocation: ?Location,
  parentPosition: ?EquipmentPosition,
  positions: Array<EquipmentPosition>,
  ports: Array<EquipmentPort>,
  properties: Array<Property>,
  futureState: ?FutureState,
  workOrder: ?WorkOrder,
  device: ?Device,
  locationHierarchy: Array<Location>,
  positionHierarchy: Array<EquipmentPosition>,
  services: Array<Service>,
};

export type EquipmentPosition = {
  id: string,
  definition: PositionDefinition,
  parentEquipment: Equipment,
  attachedEquipment: ?Equipment,
};

export type EquipmentPort = {
  id: string,
  definition: PortDefinition,
  parentEquipment: Equipment,
  link: ?Link,
  properties: Array<Property>,
  serviceEndpoints: Array<ServiceEndpoint>,
};

export type Link = {
  id: string,
  ports: Array<EquipmentPort>,
  futureState: ?FutureState,
  workOrder: ?WorkOrder,
  properties: Array<Property>,
  services: Array<Service>,
};

type Device = {
  id: string,
  up: ?boolean,
  name: string,
};

export const getInitialPortFromDefinition = (
  parentEquipment: Equipment,
  definition: PortDefinition,
): EquipmentPort => ({
  id: `EquipmentPort${parentEquipment.id}@tmp-${definition.id}`,
  definition: definition,
  parentEquipment: parentEquipment,
  properties: definition.portType
    ? definition.portType.propertyTypes.map(getInitialPropertyFromType)
    : [],
  link: null,
  serviceEndpoints: [],
});

export const getNonInstancePositionDefinitions = (
  positions: Array<EquipmentPosition>,
  positionDefinitions: Array<PositionDefinition>,
): Array<PositionDefinition> => {
  const definitionIds = positions.map(p => p.definition.id);
  return positionDefinitions.filter(def => !definitionIds.includes(def.id));
};

export const getNonInstancePortsDefinitions = (
  ports: Array<EquipmentPort>,
  portsDefinitions: Array<PortDefinition>,
): Array<PortDefinition> => {
  const definitionIds = ports.map(p => p.definition.id);
  return portsDefinitions.filter(def => !definitionIds.includes(def.id));
};
