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

import AppBar from '@material-ui/core/AppBar';
import AppContext from '@fbcnms/ui/context/AppContext';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React, {useContext} from 'react';
import SecuritySettings from './SecuritySettings';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import {Redirect, Route, Switch} from 'react-router-dom';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  tabs: {
    flex: 1,
  },
  paper: {
    margin: theme.spacing(3),
  },
}));

type Config = {ssoEnabled: boolean, isSuperUser: boolean};

export function shouldShowSettings(config: Config) {
  return getTabs(config).length > 0;
}

function getTabs(config: Config) {
  const tabs = [];
  if (!config.ssoEnabled) {
    tabs.push(
      <Tab
        component={NestedRouteLink}
        label="Security"
        to="/security/"
        key="1"
      />,
    );
  }

  if (config.isSuperUser) {
    tabs.push(
      <Tab component={NestedRouteLink} label="Users" to="/users/" key="2" />,
    );
  }

  return tabs;
}

export default function Settings(_props: {isSuperUser?: boolean}) {
  const classes = useStyles();
  const {match, relativePath, relativeUrl, location} = useRouter();
  const {user, ssoEnabled} = useContext(AppContext);

  const currentTab = findIndex(['security', 'users'], route =>
    location.pathname.startsWith(match.url + '/' + route),
  );

  return (
    <Paper className={classes.paper} elevation={2}>
      <AppBar position="static" color="default">
        <Tabs
          value={currentTab !== -1 ? currentTab : 0}
          indicatorColor="primary"
          textColor="primary"
          className={classes.tabs}>
          {getTabs({isSuperUser: user.isSuperUser, ssoEnabled})}
        </Tabs>
      </AppBar>
      <Switch>
        <Route path={relativePath('/security')} component={SecuritySettings} />
        {user.isSuperUser && (
          <Route
            path={relativePath('/users')}
            render={() => <Redirect to="/admin/users" />}
          />
        )}
        {!ssoEnabled && <Redirect to={relativeUrl('/security')} />}
      </Switch>
    </Paper>
  );
}
