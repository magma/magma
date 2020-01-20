/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

'use strict';

function getGridCoordinate(
  coordinate: number,
  collapseToGrid: boolean,
  granularity: number,
): number {
  return collapseToGrid
    ? Math.floor(coordinate / granularity) * granularity
    : coordinate;
}

function generateLatLngKey(latitude: number, longitude: number): string {
  return `${latitude.toFixed(5)},${longitude.toFixed(5)}`;
}

export {generateLatLngKey, getGridCoordinate};
