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
  date: {
    label: fbt('Date', ''),
    kind: 'date',
  },
  datetime_local: {
    label: fbt('Date & Time', ''),
    kind: 'datetime_local',
  },
  int: {
    label: fbt('Number', ''),
    kind: 'int',
  },
  float: {
    label: fbt('Float', ''),
    kind: 'float',
  },
  string: {
    label: fbt('Text', ''),
    kind: 'string',
  },
  email: {
    label: fbt('Email', ''),
    kind: 'email',
  },
  gps_location: {
    label: fbt('Coordinates', ''),
    kind: 'gps_location',
  },
  bool: {
    label: fbt('True or False', ''),
    kind: 'bool',
  },
  range: {
    label: fbt('Range', ''),
    kind: 'range',
  },
  enum: {
    label: fbt('Multiple choice', ''),
    kind: 'enum',
  },
  equipment: {
    label: fbt('Equipment', ''),
    kind: 'node',
  },
  location: {
    label: fbt('Location', ''),
    kind: 'node',
  },
  service: {
    label: fbt('Service', ''),
    featureFlag: 'services',
    kind: 'node',
  },
  work_order: {
    label: fbt('Work Order', ''),
    kind: 'node',
  },
  user: {
    label: fbt('User', ''),
    kind: 'node',
  },
};
