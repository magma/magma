/*
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
 */
import {findIndex} from 'lodash';

// Get current tab position within tabList from url specified
export function GetCurrentTabPos(url: string, tabItems: Array<string>): number {
  const tabPos = findIndex(tabItems, route =>
    location.pathname.startsWith(url + '/' + route),
  );
  return tabPos != -1 ? tabPos : 0;
}

export const DetailTabItems = ['overview', 'event', 'config'];
