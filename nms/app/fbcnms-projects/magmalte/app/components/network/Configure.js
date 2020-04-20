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
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import {Redirect, Route, Switch} from 'react-router-dom';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
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
  const {match, location, relativeUrl} = useRouter();
  const {tabRoutes} = props;

  const initialTab = findIndex(tabRoutes, route =>
    location.pathname.startsWith(match.url + '/' + route.path),
  );
  const [currentTab, setCurrentTab] = useState(
    initialTab !== -1 ? initialTab : 0,
  );

  if (location.pathname.endsWith('/configure')) {
    return <Redirect to={relativeUrl(`/${tabRoutes[0].path}`)} />;
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
      <Switch>
        {tabRoutes.map((route, i) => (
          <Route
            key={i}
            path={`${match.path}/${route.path}`}
            component={route.component}
          />
        ))}
      </Switch>
    </Paper>
  );
}
