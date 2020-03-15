/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ProjectLink} from '@fbcnms/ui/components/layout/AppDrawerProjectNavigation';
import type {Tab} from '@fbcnms/types/tabs';

const allTabs: $ReadOnlyArray<ProjectLink> = [
  {
    id: 'inventory',
    name: 'Inventory',
    secondary: 'Inventory Management',
    url: '/inventory',
  },
  {
    id: 'workorders',
    name: 'Work Orders',
    secondary: 'Workforce Management',
    url: '/workorders',
  },
  {
    id: 'nms',
    name: 'NMS',
    secondary: 'Network Management',
    url: '/nms',
  },
  {
    id: 'automation',
    name: 'Automation',
    secondary: 'Automation Management',
    url: '/automation',
  },
];

const ADMIN: ProjectLink = {
  id: 'admin',
  name: 'Admin',
  secondary: 'Administrative Tools',
  url: '/admin',
};

export function getProjectLinks(
  enabledTabs: $ReadOnlyArray<Tab>,
  user: ?{isSuperUser: boolean},
): ProjectLink[] {
  const links = allTabs.filter(tab => enabledTabs.includes(tab.id));
  if (user && user.isSuperUser) {
    links.push(ADMIN);
  }
  return links;
}

export function getProjectTabs(): {id: Tab, name: string}[] {
  return allTabs.map(tab => tab);
}
