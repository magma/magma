/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export type Node = {
  id: string,
};

export type TopologyLink = {
  source: Node,
  target: Node,
};

export type TopologyNetwork = {
  nodes: Array<Node>,
  links: Array<TopologyLink>,
};
