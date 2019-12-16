/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Equipment} from './Equipment';

export type TopologyLink = {
  source: string,
  target: string,
};

export type TopologyNetwork = {
  nodes: Array<Equipment>,
  links: Array<TopologyLink>,
};
