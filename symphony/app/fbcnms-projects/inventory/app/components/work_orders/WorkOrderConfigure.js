/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import SideNavigationPanel from '@fbcnms/ui/components/SideNavigationPanel';
import WorkOrderProjectTypes from '../configure/WorkOrderProjectTypes';
import WorkOrderTypes from '../configure/WorkOrderTypes';
import {LogEvents, ServerLogger} from '../../common/LoggingUtils';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
    height: '100vh',
    transform: 'translateZ(0)',
  },
}));

export default function WorkOrderConfigure() {
  const {location, history, relativeUrl} = useRouter();
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <SideNavigationPanel
        title="Templates"
        items={[
          {
            key: 'work_order_types',
            label: 'Work Orders',
          },
          {
            key: 'project_types',
            label: 'Projects',
          },
        ]}
        selectedItemId={location.pathname.match(/([^\/]*)\/*$/)[1]}
        onItemClicked={item => {
          ServerLogger.info(
            LogEvents.WORK_ORDERS_CONFIGURE_TAB_NAVIGATION_CLICKED,
            {
              tab: item.key,
            },
          );
          history.push(`/workorders/configure/${item.key}`);
        }}
      />
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
