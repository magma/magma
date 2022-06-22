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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import APMetrics from './APMetrics';
import AppBar from '@material-ui/core/AppBar';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AppContext from '../context/AppContext';
import CWFNetworkMetrics from './CWFNetworkMetrics';
import Grafana from '../Grafana';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import IMSIMetrics from './IMSIMetrics';
// $FlowFixMe migrated to typescript
import NestedRouteLink from '../NestedRouteLink';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import {
  Navigate,
  Route,
  Routes,
  useLocation,
  useResolvedPath,
} from 'react-router-dom';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../../theme/default';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';

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
  const resolvedPath = useResolvedPath('');
  const location = useLocation();

  const grafanaEnabled =
    useContext(AppContext).isFeatureEnabled('grafana_metrics') &&
    useContext(AppContext).user.isSuperUser;

  const tabNames = ['ap', 'network', 'subscribers'];
  if (grafanaEnabled) {
    tabNames.push('grafana');
  }

  const currentTab = findIndex(tabNames, route =>
    location.pathname.startsWith(resolvedPath.pathname + '/' + route),
  );

  return (
    <>
      <AppBar position="static" color="default" className={classes.bar}>
        <Tabs
          value={currentTab !== -1 ? currentTab : 0}
          indicatorColor="primary"
          textColor="inherit"
          className={classes.tabs}>
          <Tab component={NestedRouteLink} label="Access Points" to="ap" />
          <Tab component={NestedRouteLink} label="Network" to="network" />
          <Tab
            component={NestedRouteLink}
            label="Subscribers"
            to="subscribers"
          />
          {grafanaEnabled && (
            <Tab component={NestedRouteLink} label="Grafana" to="grafana" />
          )}
        </Tabs>
      </AppBar>
      <Routes>
        <Route path="/ap/*" element={<APMetrics />} />
        <Route path="/network" element={<CWFNetworkMetrics />} />
        <Route path="/subscribers/*" element={<IMSIMetrics />} />
        {grafanaEnabled && (
          <Route path="/grafana" element={<GrafanaDashboard />} />
        )}
        <Route index element={<Navigate to="ap" replace />} />
      </Routes>
    </>
  );
}
