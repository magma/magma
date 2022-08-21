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
 */

import Paper from '@mui/material/Paper';
import React from 'react';
import TopBar from '../../components/TopBar';
import {Navigate, Route, Routes, useLocation} from 'react-router-dom';
import {Theme} from '@mui/material/styles';
import {makeStyles} from '@mui/styles';
import type {ComponentType} from 'react';

const useStyles = makeStyles<Theme>(theme => ({
  paper: {
    margin: theme.spacing(3),
  },
}));

type Props = {
  tabRoutes: Array<TabRoute>;
};

type TabRoute = {
  component: ComponentType<any>;
  label: string;
  path: string;
};

export default function Configure(props: Props) {
  const classes = useStyles();
  const location = useLocation();
  const {tabRoutes} = props;

  if (location.pathname.endsWith('/configure')) {
    return <Navigate to={tabRoutes[0].path} replace />;
  }

  return (
    <>
      <TopBar
        header="Configure"
        tabs={tabRoutes.map(route => ({to: route.path, label: route.label}))}
      />
      <Paper className={classes.paper} elevation={2}>
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
    </>
  );
}
