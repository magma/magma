/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Location} from './Location';
import type {Property} from './Property';
import type {PropertyType} from './PropertyType';
import type {WorkOrder} from './WorkOrder';

export type ProjectType = {
  id: string,
  name: string,
  propertyTypes: Array<PropertyType>,
};

export type Project = {
  id: string,
  type: ?ProjectType,
  name: string,
  description: ?string,
  location: ?Location,
  creator: ?string,
  properties: Array<Property>,
  workOrders: Array<WorkOrder>,
  numberOfWorkOrders: number,
};
