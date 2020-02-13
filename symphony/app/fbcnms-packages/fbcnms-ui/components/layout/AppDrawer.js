/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ProjectLink} from './AppDrawerProjectNavigation';
import type {TopBarContextType} from './TopBarContext';

import * as React from 'react';
import AppDrawerProjectNavigation from './AppDrawerProjectNavigation';
import ChevronLeftIcon from '@material-ui/icons/ChevronLeft';
import Collapse from '@material-ui/core/Collapse';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import IconButton from '@material-ui/core/IconButton';
import TopBarContext from './TopBarContext';

import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';

const DRAWER_WIDTH = 240;

const useStyles = makeStyles(theme => ({
  toolbarIcon: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar,
  },
  drawerPaper: {
    position: 'relative',
    whiteSpace: 'nowrap',
    width: DRAWER_WIDTH,
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
    backgroundColor: theme.palette.primary.dark,
    color: '#fff',
  },
  drawerPaperClose: {
    overflowX: 'hidden',
    transition: theme.transitions.create('width', {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    width: theme.spacing(7),
  },
}));

type Props = {
  children: React.Node,
  projects?: ProjectLink[],
};

export default function AppDrawer(props: Props) {
  const classes = useStyles();
  const topbarContext = React.useContext<TopBarContextType>(TopBarContext);
  const projects = props.projects || [];

  return (
    <Drawer
      variant="permanent"
      classes={{
        paper: classNames(
          classes.drawerPaper,
          !topbarContext.drawerOpen && classes.drawerPaperClose,
        ),
      }}
      open={topbarContext.drawerOpen}>
      <div>
        <div className={classes.toolbarIcon}>
          <IconButton onClick={() => topbarContext.closeDrawer()}>
            <ChevronLeftIcon />
          </IconButton>
        </div>
        {projects.length > 0 && (
          <Collapse in={topbarContext.drawerOpen}>
            <div>
              <AppDrawerProjectNavigation projects={projects} />
            </div>
          </Collapse>
        )}
      </div>
      <Divider />
      {props.children}
    </Drawer>
  );
}
