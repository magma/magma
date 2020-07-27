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

import * as React from 'react';
import emptyFunction from '@fbcnms/util/emptyFunction';

export type MenuContextValue = {
  onClose: () => void,
  shown: boolean,
};

const MenuContext = React.createContext<MenuContextValue>({
  onClose: emptyFunction,
  shown: false,
});

type Props = {
  value: MenuContextValue,
  children: React.Node,
};

const MenuContextProvider = ({children, value}: Props) => {
  return <MenuContext.Provider value={value}>{children}</MenuContext.Provider>;
};

export function useMenuContext() {
  return React.useContext(MenuContext);
}

export default MenuContextProvider;
