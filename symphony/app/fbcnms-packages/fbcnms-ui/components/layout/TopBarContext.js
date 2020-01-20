/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React, {useState} from 'react';

export type TopBarContextType = {
  drawerOpen: boolean,
  openDrawer: () => void,
  closeDrawer: () => void,
};

const TopBarContext = React.createContext<TopBarContextType>({
  drawerOpen: false,
  closeDrawer: () => {},
  openDrawer: () => {},
});

type Props = {
  children: any,
};

export function TopBarContextProvider(props: Props) {
  const [drawerOpen, setDrawerOpen] = useState<boolean>(false);
  return (
    <TopBarContext.Provider
      value={{
        drawerOpen,
        openDrawer: () => setDrawerOpen(true),
        closeDrawer: () => setDrawerOpen(false),
      }}>
      {props.children}
    </TopBarContext.Provider>
  );
}

export default TopBarContext;
