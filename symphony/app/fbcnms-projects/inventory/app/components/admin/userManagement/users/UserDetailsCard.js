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
import {MessageShowingContextProvider} from '@fbcnms/ui/components/design-system/Dialog/MessageShowingContext';
import {
  Redirect,
  Route,
  Switch,
  useHistory,
  useLocation,
  useRouteMatch,
} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useContext, useMemo} from 'react';

const useStyles = makeStyles(() => ({
  root: {
    height: '100%',
    maxHeight: '100%',
    overflow: 'hidden',
    display: 'flex',
    flexDirection: 'column',
  },
  tabsContainer: {
    paddingLeft: '24px',
  },
  viewContainer: {
    flexGrow: 1,
  },
}));

export type Props = $ReadOnly<{|
  user: User,
  onChange: User => void,
|}>;

type ViewTab = {|
  tab: TabProps,
  path: string,
  view: React.Node,
|};

function UserDetailsCard(props: Props) {
  const classes = useStyles();
  const {isFeatureEnabled} = useContext(AppContext);
  const userManagementDevMode = isFeatureEnabled('user_management_dev');
  const history = useHistory();

  const {user, onChange} = props;

  const userDetailParts: Array<ViewTab> = useMemo(() => {
    const parts = [
      {
        tab: {
          label: `${fbt('Profile', '')}`,
        },
        path: 'profile',
        view: (
          <MessageShowingContextProvider>
            <UserProfilePane user={user} onChange={onChange} />
          </MessageShowingContextProvider>
        ),
      },
      {
        tab: {
          label: `${fbt('Account', '')}`,
        },
        path: 'account',
        view: <UserAccountPane user={user} onChange={onChange} />,
      },
    ];
    if (userManagementDevMode) {
      parts.push({
        tab: {
          label: `${fbt('Permissions', '')}`,
        },
        path: 'permissions',
        view: <UserPermissionsPane user={user} />,
      });
    }
    return parts;
  }, [onChange, user, userManagementDevMode]);

  const match = useRouteMatch();
  const location = useLocation();

  const activePart = useMemo(() => {
    const activePartPath = location.pathname.slice(match.url.length + 1);
    return userDetailParts.findIndex(part => part.path === activePartPath);
  }, [location.pathname, match.url.length, userDetailParts]);

  return (
    <div className={classes.root}>
      <TabsBar
        className={classes.tabsContainer}
        tabs={userDetailParts.map(part => part.tab)}
        activeTabIndex={activePart}
        onChange={tabIndex => {
          history.push(`${match.url}/${userDetailParts[tabIndex].path}`);
        }}
        spread={false}
      />
      <Switch>
        {userDetailParts.map(part => (
          <Route path={`${match.path}/${part.path}`} children={part.view} />
        ))}
        <Redirect
          from={`${match.path}/`}
          to={`${match.path}/${userDetailParts[0].path}`}
        />
      </Switch>
    </div>
  );
}
export default UserDetailsCard;
