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

import AccountSettings from '../AccountSettings';
import AppContent from '../layout/AppContent';
import AppSideBar from '../AppSideBar';
import ApplicationMain from '../../components/ApplicationMain';
import AssignmentIcon from '@material-ui/icons/Assignment';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CloudMetrics from '../../views/metrics/CloudMetrics';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Features from '../../views/features/Features';
import FlagIcon from '@material-ui/icons/Flag';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import OrganizationEdit from '../../views/organizations/OrganizationEdit';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Organizations from '../../views/organizations/Organizations';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import ShowChartIcon from '@material-ui/icons/ShowChart';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import UsersSettings from '../UsersSettings';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AppContextProvider} from '../context/AppContext';
import {Navigate, Outlet, Route, Routes} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
  },
}));

function Frame() {
  const classes = useStyles();

  const sidebarItems = [
    {
      label: 'Organizations',
      path: '/host/organizations',
      icon: <AssignmentIcon />,
    },
    {
      label: 'Features',
      path: '/host/features',
      icon: <FlagIcon />,
    },
    {label: 'Metrics', path: '/host/metrics', icon: <ShowChartIcon />},
    {
      label: 'Users',
      path: '/host/users',
      icon: <PeopleIcon />,
    },
  ];

  return (
    <div className={classes.root}>
      <AppSideBar items={sidebarItems} />
      <AppContent>
        <Outlet />
      </AppContent>
    </div>
  );
}

const Index = () => {
  return (
    <ApplicationMain>
      <AppContextProvider isOrganizations={true}>
        <Routes>
          <Route path="/host" element={<Frame />}>
            <Route
              path="organizations/detail/:name"
              element={<OrganizationEdit />}
            />
            <Route path="organizations/*" element={<Organizations />} />
            <Route path="features/*" element={<Features />} />
            <Route path="metrics" element={<CloudMetrics />} />
            <Route path="users" element={<UsersSettings />} />
            <Route path="settings" element={<AccountSettings />} />
            <Route index element={<Navigate to="organizations" replace />} />
          </Route>
        </Routes>
      </AppContextProvider>
    </ApplicationMain>
  );
};

export default Index;
