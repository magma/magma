/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AccountSettings from './AccountSettings';
import AppContext, {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import React, {useContext} from 'react';
import symphony from '@fbcnms/ui/theme/symphony';
import {getProjectLinks} from '@fbcnms/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {useMainContext} from '../../components/MainContext';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    background: symphony.palette.background,
  },
}));

function SettingsApp() {
  const classes = useStyles();
  const {tabs} = useContext(AppContext);
  const {integrationUserDefinition} = useMainContext();

  return (
    <div className={classes.root}>
      <AppSideBar
        showSettings={true}
        user={integrationUserDefinition}
        projects={getProjectLinks(tabs, integrationUserDefinition)}
        mainItems={[]}
      />
      <AccountSettings />
    </div>
  );
}

export default function Settings() {
  return (
    <ApplicationMain>
      <AppContextProvider>
        <SettingsApp />
      </AppContextProvider>
    </ApplicationMain>
  );
}
