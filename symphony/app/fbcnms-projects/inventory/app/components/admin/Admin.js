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
import AdminContextProvider from '@fbcnms/magmalte/app/components/admin/AdminContextProvider';
import AdminMain from '@fbcnms/magmalte/app/components/admin/AdminMain';
import AppContext from '@fbcnms/ui/context/AppContext';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import AssignmentIcon from '@material-ui/icons/Assignment';
import AuditLog from '@fbcnms/magmalte/app/components/admin/AuditLog';
import NavListItem from '@fbcnms/ui/components/NavListItem';
import Networks from '@fbcnms/magmalte/app/components/admin/Networks';
import Paper from '@material-ui/core/Paper';
import PeopleIcon from '@material-ui/icons/People';
import SecuritySettings from '@fbcnms/magmalte/app/components/SecuritySettings';
import SignalCellularAlt from '@material-ui/icons/SignalCellularAlt';
import UsersSettings from '@fbcnms/magmalte/app/components/UsersSettings';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
import {useRelativeUrl} from '@fbcnms/ui/hooks/useRouter';

const useStyles = makeStyles(theme => ({
  paper: {
    margin: theme.spacing(3),
    padding: theme.spacing(),
  },
}));

function NavItems() {
  const relativeUrl = useRelativeUrl();
  const {isFeatureEnabled} = useContext(AppContext);
  const auditLogEnabled = isFeatureEnabled('audit_log_view');

  return (
    <>
      <NavListItem
        label="Users"
        path={relativeUrl('/users')}
        icon={<PeopleIcon />}
      />
      {auditLogEnabled && (
        <NavListItem
          label="Audit Log"
          path={relativeUrl('/audit_log')}
          icon={<AssignmentIcon />}
        />
      )}
      <NavListItem
        label="Networks"
        path={relativeUrl('/networks')}
        icon={<SignalCellularAlt />}
      />
    </>
  );
}

function NavRoutes() {
  const classes = useStyles();
  const relativeUrl = useRelativeUrl();
  return (
    <Switch>
      <Route path={relativeUrl('/users')} component={UsersSettings} />
      <Route path={relativeUrl('/audit_log')} component={AuditLog} />
      <Route path={relativeUrl('/networks')} component={Networks} />
      <Route
        path={relativeUrl('/settings')}
        render={() => (
          <Paper className={classes.paper}>
            <SecuritySettings />
          </Paper>
        )}
      />
      <Redirect to={relativeUrl('/users')} />
    </Switch>
  );
}

export default () => (
  <ApplicationMain>
    <AdminContextProvider>
      <AdminMain
        navRoutes={() => <NavRoutes />}
        navItems={() => <NavItems />}
      />
    </AdminContextProvider>
  </ApplicationMain>
);
