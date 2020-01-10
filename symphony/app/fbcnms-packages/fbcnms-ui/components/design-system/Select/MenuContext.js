/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
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
