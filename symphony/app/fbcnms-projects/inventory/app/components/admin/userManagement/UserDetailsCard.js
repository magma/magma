/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {TabProps} from '@fbcnms/ui/components/design-system/Tabs/TabsBar';
import type {User} from './TempTypes';

import * as React from 'react';
import TabsBar from '@fbcnms/ui/components/design-system/Tabs/TabsBar';
import UserAccountPane from './UserAccountPane';
import UserPermissionsPane from './UserPermissionsPane';
import UserProfilePane from './UserProfilePane';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';
import {useMemo, useState} from 'react';

const useStyles = makeStyles(() => ({
  tabsContainer: {
    paddingLeft: '24px',
  },
  viewContainer: {
    padding: '24px',
  },
}));

export type Props = {
  user: User,
  onChange?: ?(User) => void,
};

type ViewTab = {|
  tab: TabProps,
  view: React.Node,
|};

export default function UserDetailsCard(props: Props) {
  const classes = useStyles();
  const {user, onChange} = props;
  const userDetailParts: Array<ViewTab> = useMemo(
    () => [
      {
        tab: {
          label: `${fbt('Profile', '')}`,
        },
        view: <UserProfilePane user={user} onChange={onChange} />,
      },
      {
        tab: {
          label: `${fbt('Account', '')}`,
        },
        view: <UserAccountPane user={user} />,
      },
      {
        tab: {
          label: `${fbt('Permissions', '')}`,
        },
        view: <UserPermissionsPane user={user} />,
      },
    ],
    [onChange, user],
  );
  const [activePart, setActivePart] = useState(0);

  return (
    <div>
      <TabsBar
        className={classes.tabsContainer}
        tabs={userDetailParts.map(part => part.tab)}
        activeTabIndex={activePart}
        onChange={setActivePart}
        spread={false}
      />
      <div className={classes.viewContainer}>
        {userDetailParts[activePart].view}
      </div>
    </div>
  );
}
