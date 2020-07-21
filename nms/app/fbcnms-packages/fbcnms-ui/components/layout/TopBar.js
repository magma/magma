/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {TopBarContextType} from './TopBarContext';

import * as React from 'react';
import classNames from 'classnames';

import AppBar from '@material-ui/core/AppBar';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import Toolbar from '@material-ui/core/Toolbar';
import TopBarContext from './TopBarContext';
import Typography from '@material-ui/core/Typography';

import {makeStyles} from '@material-ui/styles';

const DRAWER_WIDTH = 240;

const useStyles = makeStyles(theme => ({
  appBarSpacer: theme.mixins.toolbar,
  toolbar: {
    paddingRight: 24, // keep right padding when drawer closed
  },
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  appBarShift: {
    marginLeft: DRAWER_WIDTH,
    width: `calc(100% - ${DRAWER_WIDTH}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  link: {
    textDecoration: 'none',
  },
  menuButton: {
    marginLeft: 12,
    marginRight: 36,
  },
  menuButtonHidden: {
    display: 'none',
  },
  title: {
    flexGrow: 1,
  },
}));

type Props = {
  children: React.Node,
  title?: string,
};

export default function TopBar(props: Props) {
  const classes = useStyles();
  const context = React.useContext<TopBarContextType>(TopBarContext);
  return (
    <>
      <AppBar
        position="absolute"
        color="primary"
        className={classNames(
          classes.appBar,
          context.drawerOpen && classes.appBarShift,
        )}>
        <Toolbar
          disableGutters={!context.drawerOpen}
          className={classes.toolbar}>
          <IconButton
            color="inherit"
            aria-label="Open drawer"
            onClick={() => context.openDrawer()}
            className={classNames(
              classes.menuButton,
              context.drawerOpen && classes.menuButtonHidden,
            )}>
            <MenuIcon />
          </IconButton>
          <Typography
            variant="h6"
            color="inherit"
            noWrap
            className={classes.title}>
            {props.title}
          </Typography>
          {props.children}
        </Toolbar>
      </AppBar>
      <div className={classes.appBarSpacer} />
    </>
  );
}
