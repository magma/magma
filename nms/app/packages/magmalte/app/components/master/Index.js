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

import AppContent from '../layout/AppContent';
import AppContext, {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar';
import ApplicationMain from '../../components/ApplicationMain';
import AssignmentIcon from '@material-ui/icons/Assignment';
import CloudMetrics from '@fbcnms/ui/master/CloudMetrics';
import Features from '@fbcnms/ui/master/Features';
import FlagIcon from '@material-ui/icons/Flag';
import NavListItem from '@fbcnms/ui/components/NavListItem';
import OrganizationEdit from '@fbcnms/ui/master/OrganizationEdit';
import Organizations from '@fbcnms/ui/master/Organizations';
import Paper from '@material-ui/core/Paper';
import PeopleIcon from '@material-ui/icons/People';
import React, {useContext} from 'react';
import SecuritySettings from '@fbcnms/magmalte/app/components/SecuritySettings';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import UsersSettings from '@fbcnms/magmalte/app/components/admin/userManagement/UsersSettings';
import nullthrows from '@fbcnms/util/nullthrows';
import {Redirect, Route, Switch} from 'react-router-dom';
import {getProjectTabs as getAllProjectTabs} from '@fbcnms/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {useRelativeUrl} from '@fbcnms/ui/hooks/useRouter';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  paper: {
    margin: theme.spacing(3),
    padding: theme.spacing(),
  },
}));

const accessibleTabs = ['NMS'];

function NavItems() {
  const relativeUrl = useRelativeUrl();
  return (
    <>
      <NavListItem
        label="Organizations"
        path={relativeUrl('/organizations')}
        icon={<AssignmentIcon />}
      />
      <NavListItem
        label="Features"
        path={relativeUrl('/features')}
        icon={<FlagIcon />}
      />
      <NavListItem
        label="Metrics"
        path={relativeUrl('/metrics')}
        icon={<ShowChartIcon />}
      />
      <NavListItem
        label="Users"
        path={relativeUrl('/users')}
        icon={<PeopleIcon />}
      />
    </>
  );
}

function Master() {
  const classes = useStyles();
  const {user, ssoEnabled} = useContext(AppContext);
  const relativeUrl = useRelativeUrl();

  return (
    <div className={classes.root}>
      <AppSideBar
        mainItems={<NavItems />}
        user={nullthrows(user)}
        showSettings={!ssoEnabled}
      />
      <AppContent>
        <Switch>
          <Route
            path={relativeUrl('/organizations/detail/:name')}
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
            path={relativeUrl('/organizations')}
            component={Organizations}
          />
          <Route path={relativeUrl('/features')} component={Features} />
          <Route path={relativeUrl('/metrics')} component={CloudMetrics} />
          <Route path={relativeUrl('/users')} component={UsersSettings} />
          <Route
            path={relativeUrl('/settings')}
            render={() => (
              <Paper className={classes.paper}>
                <SecuritySettings />
              </Paper>
            )}
          />
          <Redirect to={relativeUrl('/organizations')} />
        </Switch>
      </AppContent>
    </div>
  );
}

const Index = () => {
  return (
    <ApplicationMain>
      <AppContextProvider>
        <Master />
      </AppContextProvider>
    </ApplicationMain>
  );
};

export default () => <Route path="/master" component={Index} />;
