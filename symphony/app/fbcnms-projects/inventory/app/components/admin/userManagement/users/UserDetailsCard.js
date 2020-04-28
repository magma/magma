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
import type {User} from '../utils/UserManagementUtils';

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import TabsBar from '@fbcnms/ui/components/design-system/Tabs/TabsBar';
import UserAccountPane from './UserAccountPane';
import UserPermissionsPane from './UserPermissionsPane';
import UserProfilePane from './UserProfilePane';
import fbt from 'fbt';
import {makeStyles} from '@material-ui/styles';
import {useContext, useMemo, useState} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    maxHeight: '100%',
    overflow: 'hidden',
    display: 'flex',
    flexDirection: 'column',
  },
  tabsContainer: {
    paddingLeft: '24px',
  },
  viewContainer: {
    padding: '24px',
    overflowY: 'auto',
  },
}));

export type Props = {
  user: User,
  onChange: User => void,
};

type ViewTab = {|
  tab: TabProps,
  view: React.Node,
|};

export default function UserDetailsCard(props: Props) {
  const classes = useStyles();
  const {isFeatureEnabled} = useContext(AppContext);
  const userManagementDevMode = isFeatureEnabled('user_management_dev');

  const {user, onChange} = props;
  const userDetailParts: Array<ViewTab> = useMemo(() => {
    const parts = [
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
        view: <UserAccountPane user={user} onChange={onChange} />,
      },
    ];
    if (userManagementDevMode) {
      parts.push({
        tab: {
          label: `${fbt('Permissions', '')}`,
        },
        view: <UserPermissionsPane user={user} />,
      });
    }
    return parts;
  }, [onChange, user, userManagementDevMode]);
  const [activePart, setActivePart] = useState(0);

  return (
    <div className={classes.root}>
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
