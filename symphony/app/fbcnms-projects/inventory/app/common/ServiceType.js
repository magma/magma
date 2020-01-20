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

export type ServiceType = {
  id: string,
  name: string,
  propertyTypes: Array<PropertyType>,
  numberOfServices: number,
};
