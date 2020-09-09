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

import APMetrics from './APMetrics';
import AppBar from '@material-ui/core/AppBar';
import AppContext from '@fbcnms/ui/context/AppContext';
import CWFNetworkMetrics from './CWFNetworkMetrics';
import Grafana from '../Grafana';
import IMSIMetrics from './IMSIMetrics';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import {Redirect, Route, Switch} from 'react-router-dom';
import {colors} from '../../theme/default';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(_ => ({
  bar: {
    backgroundColor: colors.primary.brightGray,
  },
  tabs: {
    flex: 1,
    color: colors.primary.white,
  },
}));

function GrafanaDashboard() {
  return <Grafana grafanaURL="/grafana" />;
}

export default function () {
  const classes = useStyles();
  const {match, relativePath, relativeUrl, location} = useRouter();

  const grafanaEnabled =
    useContext(AppContext).isFeatureEnabled('grafana_metrics') &&
    useContext(AppContext).user.isSuperUser;

  const tabNames = ['ap', 'network', 'subscribers'];
  if (grafanaEnabled) {
    tabNames.push('grafana');
  }

  const currentTab = findIndex(tabNames, route =>
    location.pathname.startsWith(match.url + '/' + route),
  );

  return (
    <>
      <AppBar position="static" color="default" className={classes.bar}>
        <Tabs
          value={currentTab !== -1 ? currentTab : 0}
          indicatorColor="primary"
          textColor="inherit"
          className={classes.tabs}>
          <Tab component={NestedRouteLink} label="Access Points" to="/ap" />
          <Tab component={NestedRouteLink} label="Network" to="/network" />
          <Tab
            component={NestedRouteLink}
            label="Subscribers"
            to="/subscribers"
          />
          {grafanaEnabled && (
            <Tab component={NestedRouteLink} label="Grafana" to="/grafana" />
          )}
        </Tabs>
      </AppBar>
      <Switch>
        <Route path={relativePath('/ap')} component={APMetrics} />
        <Route path={relativePath('/network')} component={CWFNetworkMetrics} />
        <Route path={relativePath('/subscribers')} component={IMSIMetrics} />
        {grafanaEnabled && (
          <Route path={relativePath('/grafana')} component={GrafanaDashboard} />
        )}
        <Redirect to={relativeUrl('/ap')} />
      </Switch>
    </>
  );
}
