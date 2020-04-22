/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {PropertyTypeInfo} from './form/PropertyTypeTable';

import fbt from 'fbt';

export const PropertyTypeLabels: {[string]: PropertyTypeInfo} = {
  date: {label: fbt('Date', '')},
  datetime_local: {label: fbt('Date & Time', '')},
  int: {label: fbt('Number', '')},
  float: {label: fbt('Float', '')},
  string: {label: fbt('Text', '')},
  email: {label: fbt('Email', '')},
  gps_location: {label: fbt('Coordinates', '')},
  bool: {label: fbt('True or False', '')},
  range: {label: fbt('Range', '')},
  enum: {label: fbt('Multiple choice', '')},
  equipment: {label: fbt('Equipment', ''), isNode: true},
  location: {label: fbt('Location', ''), isNode: true},
  service: {
    label: fbt('Service', ''),
    featureFlag: 'services',
    isNode: true,
  },
  work_order: {
    label: fbt('Work Order', ''),
    isNode: true,
  },
};
