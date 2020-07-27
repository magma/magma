/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
  {
    id: 'hub',
    name: 'Hub',
    secondary: 'Network Hub',
    url: '/hub',
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
