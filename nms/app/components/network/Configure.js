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
 * @flow
 * @format
 */

import type {ComponentType} from 'react';

import AppBar from '@material-ui/core/AppBar';
// $FlowFixMe migrated to typescript
import NestedRouteLink from '../NestedRouteLink';
import Paper from '@material-ui/core/Paper';
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
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  paper: {
    margin: theme.spacing(3),
  },
  tabs: {
    flex: 1,
  },
}));

type Props = {
  tabRoutes: TabRoute[],
};

type TabRoute = {
  component: ComponentType<any>,
  label: string,
  path: string,
};

export default function Configure(props: Props) {
  const classes = useStyles();
  const location = useLocation();
  const resolvedPath = useResolvedPath('');
  const {tabRoutes} = props;

  const initialTab = findIndex(tabRoutes, route =>
    location.pathname.startsWith(resolvedPath.pathname + '/' + route.path),
  );
  const [currentTab, setCurrentTab] = useState(
    initialTab !== -1 ? initialTab : 0,
  );

  if (location.pathname.endsWith('/configure')) {
    return <Navigate to={tabRoutes[0].path} replace />;
  }

  return (
    <Paper className={classes.paper} elevation={2}>
      <AppBar position="static" color="default">
        <Tabs
          value={currentTab}
          indicatorColor="primary"
          textColor="primary"
          onChange={(_, tab) => setCurrentTab(tab)}
          className={classes.tabs}>
          {tabRoutes.map((route, i) => (
            <Tab
              key={i}
              component={NestedRouteLink}
              label={route.label}
              to={route.path}
            />
          ))}
        </Tabs>
      </AppBar>
      <Routes>
        {tabRoutes.map((route, i) => (
          <Route
            key={i}
            path={`${route.path}/*`}
            element={<route.component />}
          />
        ))}
      </Routes>
    </Paper>
  );
}
