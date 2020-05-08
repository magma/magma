/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Equipment, EquipmentPort, Link} from './Equipment';
import type {Property} from './Property';
import type {ServiceEndpointDefinition, ServiceType} from './ServiceType';

export type ServiceStatus =
  | 'PENDING'
  | 'IN_SERVICE'
  | 'MAINTENANCE'
  | 'DISCONNECTED';

export const serviceStatusToVisibleNames = {
  PENDING: 'Pending',
  IN_SERVICE: 'In Service',
  MAINTENANCE: 'Maintenance',
  DISCONNECTED: 'Disconnected',
};

export const serviceStatusToColor = {
  PENDING: 'orange',
  IN_SERVICE: 'green',
  MAINTENANCE: 'orange',
  DISCONNECTED: 'gray',
};

export type Customer = {
  id: string,
  name: string,
};

export type ServiceEndpoint = {
  id: string,
  port: EquipmentPort,
  definition: ServiceEndpointDefinition,
  service: Service,
  equipment: Equipment,
};

export type Service = {
  id: string,
  name: string,
  externalId: ?string,
  status: ServiceStatus,
  customer: ?Customer,
  serviceType: ServiceType,
  upstream: Array<Service>,
  downstream: Array<Service>,
  properties: Array<Property>,
  endpoints: Array<ServiceEndpoint>,
  links: Array<Link>,
};
