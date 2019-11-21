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
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';

import Actions from './Actions';
import ActionsList from './ActionsList';
import AppContent from '@fbcnms/ui/components/layout/AppContent';
import AppContext from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import ComputerIcon from '@material-ui/icons/Computer';
import NavListItem from '@fbcnms/ui/components/NavListItem';
import {Redirect, Route, Switch} from 'react-router-dom';

import {getProjectLinks} from '@fbcnms/magmalte/app/common/projects';
import {makeStyles} from '@material-ui/styles';
import {shouldShowSettings} from '@fbcnms/magmalte/app/components/Settings';
import {useContext} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
  },
}));

function NavItems() {
  const {relativeUrl} = useRouter();

  return (
    <>
      <NavListItem
        label="Actions"
        path={relativeUrl('/actions')}
        icon={<ComputerIcon />}
      />
    </>
  );
}

function NavRoutes() {
  const {relativeUrl} = useRouter();
  return (
    <Switch>
      <Route path={relativeUrl('/actions/list')} component={ActionsList} />
      <Route path={relativeUrl('/actions')} component={Actions} />
      <Redirect to={relativeUrl('/actions')} />
    </Switch>
  );
}

function Automation() {
  const classes = useStyles();
  const {tabs, user, ssoEnabled} = useContext(AppContext);
  return (
    <div className={classes.root}>
      <AppSideBar
        mainItems={<NavItems />}
        projects={getProjectLinks(tabs, user)}
        user={user}
        showSettings={shouldShowSettings({
          isSuperUser: user.isSuperUser,
          ssoEnabled,
        })}
      />
      <AppContent>
        <NavRoutes />
      </AppContent>
    </div>
  );
}

export default () => {
  return (
    <ApplicationMain>
      <AppContextProvider>
        <Automation />
      </AppContextProvider>
    </ApplicationMain>
  );
};
