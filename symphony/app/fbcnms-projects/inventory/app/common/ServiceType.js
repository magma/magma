/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {DiscoveryMethod} from '../components/configure/__generated__/AddEditServiceTypeCard_editingServiceType.graphql';
import type {PropertyType} from './PropertyType';

type ServiceEquipmentType = {
  id: string,
  name: string,
};

export type ServiceType = {
  id: string,
  name: string,
  propertyTypes: Array<PropertyType>,
  discoveryMethod: ?DiscoveryMethod,
  numberOfServices: number,
  endpointDefinitions: Array<ServiceEndpointDefinition>,
};

export type ServiceEndpointDefinition = {
  id: string,
  name: string,
  role: ?string,
  index: number,
  equipmentType: ?ServiceEquipmentType,
};
