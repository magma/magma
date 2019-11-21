/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import WorkOrderProjectTypes from '../configure/WorkOrderProjectTypes';
import WorkOrderTypes from '../configure/WorkOrderTypes';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    height: '100vh',
    transform: 'translateZ(0)',
  },
  tabs: {
    backgroundColor: 'white',
    borderBottom: `1px ${theme.palette.grey[200]} solid`,
    minHeight: '60px',
    overflow: 'visible',
  },
  tabContainer: {
    width: '250px',
  },
  tabsRoot: {
    top: 0,
    left: 0,
    right: 0,
    height: '60px',
  },
}));

function getTabURI(newValue: number) {
  switch (newValue) {
    case 0:
      return 'work_order_types';
    case 1:
      return 'project_types';
    default:
      return '';
  }
}

function getTabValue(tabURI: string) {
  switch (tabURI) {
    case 'work_order_types':
      return 0;
    case 'project_types':
      return 1;
    default:
      return 0;
  }
}

export default function WorkOrderConfigure() {
  const {location, history, relativeUrl} = useRouter();
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <Tabs
        className={classes.tabs}
        classes={{flexContainer: classes.tabsRoot}}
        value={getTabValue(location.pathname.match(/([^\/]*)\/*$/)[1])}
        onChange={(e, newValue) => {
          const tab = getTabURI(newValue);
          ServerLogger.info(
            LogEvents.WORK_ORDERS_CONFIGURE_TAB_NAVIGATION_CLICKED,
            {
              tab,
            },
          );
          history.push(`/workorders/configure/${tab}`);
        }}
        indicatorColor="primary"
        textColor="primary">
        <Tab
          classes={{root: classes.tabContainer}}
          label="Work Order Templates"
        />
        <Tab classes={{root: classes.tabContainer}} label="Project Templates" />
      </Tabs>
      <Switch>
        <Route
          path={relativeUrl('/work_order_types')}
          component={WorkOrderTypes}
        />
        <Route
          path={relativeUrl('/project_types')}
          component={WorkOrderProjectTypes}
        />
        <Redirect
          from={relativeUrl('/')}
          to={relativeUrl('/work_order_types')}
        />
      </Switch>
    </div>
  );
}
