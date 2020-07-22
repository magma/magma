/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import AppContent from '../layout/AppContent';
import AppContext from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar';

import {getProjectLinks} from '@fbcnms/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {shouldShowSettings} from '../Settings';

const useStyles = makeStyles(_theme => ({
  root: {
    display: 'flex',
  },
}));

export default function ErrorLayout({children}: {children: React.Node}) {
  const classes = useStyles();
  const {user, tabs, ssoEnabled} = React.useContext(AppContext);

  return (
    <div className={classes.root}>
      <AppSideBar
        mainItems={[]}
        secondaryItems={[]}
        projects={getProjectLinks(tabs, user)}
        showSettings={shouldShowSettings({
          isSuperUser: user.isSuperUser,
          ssoEnabled,
        })}
        user={user}
      />
      <AppContent>{children}</AppContent>
    </div>
  );
}
