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
import AppSideBar from '../../../fbc_js_core/ui/components/layout/AppSideBar';
import ApplicationMain from '../../components/ApplicationMain';
import AssignmentIcon from '@material-ui/icons/Assignment';
import CloudMetrics from '../../views/metrics/CloudMetrics';
import Features from '../../../fbc_js_core/ui/host/Features';
import FlagIcon from '@material-ui/icons/Flag';
import OrganizationEdit from '../../../fbc_js_core/ui/host/OrganizationEdit';
import Organizations from '../../../fbc_js_core/ui/host/Organizations';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import UsersSettings from '../admin/userManagement/UsersSettings';
import {AppContextProvider} from '../../../fbc_js_core/ui/context/AppContext';
import {Redirect, Route, Switch} from 'react-router-dom';
import {getProjectTabs as getAllProjectTabs} from '../../../fbc_js_core/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {useRelativeUrl} from '../../../fbc_js_core/ui/hooks/useRouter';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
  },
}));

const accessibleTabs = ['NMS'];

function Host() {
  const classes = useStyles();
  const relativePath = useRelativeUrl();

  const sidebarItems = [
    {
      label: 'Organizations',
      path: '/organizations',
      icon: <AssignmentIcon />,
    },
    {
      label: 'Features',
      path: '/features',
      icon: <FlagIcon />,
    },
    {label: 'Metrics', path: '/metrics', icon: <ShowChartIcon />},
    {
      label: 'Users',
      path: '/users',
      icon: <PeopleIcon />,
    },
  ];

  return (
    <div className={classes.root}>
      <AppSideBar items={sidebarItems} />
      <AppContent>
        <Switch>
          <Route
            path={relativePath('/organizations/detail/:name')}
            render={() => (
              <OrganizationEdit
                getProjectTabs={() =>
                  getAllProjectTabs().filter(tab =>
                    accessibleTabs.includes(tab.name),
                  )
                }
              />
            )}
          />
          <Route
            path={relativePath('/organizations')}
            component={Organizations}
          />
          <Route path={relativePath('/features')} component={Features} />
          <Route path={relativePath('/metrics')} component={CloudMetrics} />
          <Route path={relativePath('/users')} component={UsersSettings} />
          <Route path={relativePath('/settings')} component={AccountSettings} />
          <Redirect to={relativePath('/organizations')} />
        </Switch>
      </AppContent>
    </div>
  );
}

const Index = () => {
  return (
    <ApplicationMain>
      <AppContextProvider isOrganizations={true}>
        <Host />
      </AppContextProvider>
    </ApplicationMain>
  );
};

export default () => <Route path="/host" component={Index} />;
