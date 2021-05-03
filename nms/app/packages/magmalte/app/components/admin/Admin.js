/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import AdminContextProvider from './AdminContextProvider';
import AdminMain from './AdminMain';
import ApplicationMain from '../ApplicationMain';
import AssignmentIcon from '@material-ui/icons/Assignment';
import AuditLog from './AuditLog';
import NavListItem from '@fbcnms/ui/components/NavListItem';
import Networks from './Networks';
import Paper from '@material-ui/core/Paper';
import PeopleIcon from '@material-ui/icons/People';
import SecuritySettings from '../SecuritySettings';
import SignalCellularAlt from '@material-ui/icons/SignalCellularAlt';
import UsersSettings from '../UsersSettings';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  paper: {
    margin: theme.spacing(3),
    padding: theme.spacing(),
  },
}));

function NavItems() {
  const {relativeUrl} = useRouter();

  return (
    <>
      <NavListItem
        label="Users"
        path={relativeUrl('/users')}
        icon={<PeopleIcon />}
      />
      <NavListItem
        label="Audit Log"
        path={relativeUrl('/audit_log')}
        icon={<AssignmentIcon />}
      />
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
  const {relativeUrl} = useRouter();
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
