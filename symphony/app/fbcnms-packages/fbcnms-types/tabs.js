/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export type Tab =
  | 'automation'
  | 'admin'
  | 'inventory'
  | 'nms'
  | 'workorders'
  | 'hub';

export const TABS: {[string]: Tab} = Object.freeze({
  admin: 'admin',
  automation: 'automation',
  inventory: 'inventory',
  nms: 'nms',
  workorders: 'workorders',
  hub: 'hub',
});

export function coerceToTab(tab: string): Tab {
  if (TABS[tab]) {
    return TABS[tab];
  }
  throw new Error('Invalid tab: ' + tab);
}
