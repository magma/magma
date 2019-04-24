/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ProjectLink} from './AppDrawerProjectNavigation';

import {makeStyles} from '@material-ui/styles';
import AppSideBarProjectNavigation from './AppSideBarProjectNavigation.react';
import ProfileButton from '../ProfileButton.react';
import React from 'react';

const useStyles = makeStyles(theme => ({
  root: {
    alignItems: 'center',
    backgroundColor: theme.palette.blueGrayDark,
    boxShadow: '1px 0px 0px 0px rgba(0, 0, 0, 0.1)',
    display: 'flex',
    flexDirection: 'column',
    height: '100vh',
    width: '82px',
    padding: '60px 0px 20px 0px',
  },
  mainItems: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    width: '100%',
    flexGrow: 1,
  },
  secondaryItems: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    width: '100%',
  },
}));

type Props = {
  mainItems: any,
  secondaryItems?: any,
  user: {
    email: string,
    isSuperUser: boolean,
  },
  projects?: ProjectLink[],
};

export default function AppDrawer(props: Props) {
  const {user} = props;
  const classes = useStyles();
  const projects = props.projects || [];

  return (
    <div className={classes.root}>
      <div className={classes.mainItems}>{props.mainItems}</div>
      <div className={classes.secondaryItems}>
        {props.secondaryItems}
        <ProfileButton user={user} />
      </div>
      <AppSideBarProjectNavigation projects={projects} />
    </div>
  );
}
