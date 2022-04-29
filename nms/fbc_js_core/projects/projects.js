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

import type {Tab} from '../../fbc_js_core/types/tabs';

type ProjectLink = {
  id: Tab,
  name: string,
  secondary: string,
  url: string,
};

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

export function getProjectTabs(): {id: Tab, name: string}[] {
  return allTabs.map(tab => tab);
}
