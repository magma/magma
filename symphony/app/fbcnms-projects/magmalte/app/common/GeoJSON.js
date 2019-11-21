/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {magmad_gateway} from '@fbcnms/magma-api';

export type MagmaFeatureCollection = {
  type: 'FeatureCollection',
  features: MagmaGatewayFeature[],
};

export type MagmaGatewayFeature = {
  type: 'Feature',
  geometry: {
    type: 'Point',
    coordinates: [number, number],
  },
  properties: {
    id: string | number,
    name?: string,
    iconSize: IconSize,
    gateway?: magmad_gateway,
    [key: string]: any,
  },
};

export type MagmaConnectionFeature = {
  type: 'Feature',
  geometry: {
    type: 'LineString',
    coordinates: Array<[number, number]>,
  },
  properties: {
    id: string | number,
    name: string,

    [string]: any,
  },
};

export type IconSize = [number, number];
