/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment, Link} from './Equipment';
import type {Property} from './Property';
import type {ServiceType} from './ServiceType';

export type Customer = {
  id: string,
  name: string,
};

export type Service = {
  id: string,
  name: string,
  externalId: ?string,
  customer: ?Customer,
  serviceType: ServiceType,
  upstream: Array<Service>,
  downstream: Array<Service>,
  properties: Array<Property>,
  terminationPoints: Array<Equipment>,
  links: Array<Link>,
};
