/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export type MapFeatureCollection = {
  type: 'FeatureCollection',
  features: MapFeature[],
};

export type MapFeature = {
  type: 'Feature',
  geometry: {
    type: 'Point',
    coordinates: [number, number],
  },
  properties: {
    id: string | number,
    name?: string,
    iconSize: IconSize,
    [key: string]: any,
  },
};

export type IconSize = [number, number];
