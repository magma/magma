/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AppContent from '@fbcnms/ui/components/layout/AppContent';
import AppContext, {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import Button from '@fbcnms/ui/components/design-system/Button';
import CreateService from './CreateService';
import HubVersion from './HubVersion';
import NavListItem from '@fbcnms/ui/components/NavListItem';
import React, {useContext} from 'react';
import RouterIcon from '@material-ui/icons/Router';
import Shuffle from '@material-ui/icons/Shuffle';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import WorkflowApp from './workflow/App';
import nullthrows from '@fbcnms/util/nullthrows';
import {Redirect, Route, Switch} from 'react-router-dom';
import {getProjectLinks} from '@fbcnms/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {shouldShowSettings} from '@fbcnms/magmalte/app/components/Settings';
import {useRelativeUrl} from '@fbcnms/ui/hooks/useRouter';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  header: {
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing(3),
  },
}));

function NavBarItems() {
  return [
    <NavListItem
      key={1}
      label="Services"
      path="/hub/services"
      icon={<RouterIcon />}
      onClick={() => {}}
    />,
    <NavListItem
      key={2}
      label="Workflows"
      path="/hub/workflows"
      icon={<Shuffle />}
      onClick={() => {}}
    />,
  ];
}

function CreateServiceForm() {
  const classes = useStyles();
  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <Text variant="h5">Create a Service</Text>
      </div>
      <br />
      <HubVersion />
      <br />
      <Text>
        Create a simple service (only single port vlan is supported for now).
      </Text>
      <br />
      <br />
      <CreateService />
      <br />
      <br />
      <br />
      <Text>Service ID</Text>
      <TextInput />
      <Button>Delete Service</Button>
    </div>
  );
}

function Main(subApp) {
  const {user, tabs, ssoEnabled} = useContext(AppContext);
  const classes = useStyles();
  return (
    <div className={classes.root}>
      <AppSideBar
        mainItems={<NavBarItems />}
        secondaryItems={[]}
        projects={getProjectLinks(tabs, user)}
        showSettings={shouldShowSettings({
          isSuperUser: user.isSuperUser,
          ssoEnabled,
        })}
        user={nullthrows(user)}
      />
      <AppContent>{subApp}</AppContent>
    </div>
  );
}

export default () => {
  const relativeUrl = useRelativeUrl();
  const cs = () => Main(CreateServiceForm());
  const wf = () => Main(WorkflowApp());
  return (
    <ApplicationMain>
      <AppContextProvider>
        <Switch>
          <Route path={relativeUrl('/services')} component={cs} />
          <Route path={relativeUrl('/workflows')} component={wf} />
          <Redirect exact from="/" to={relativeUrl('/hub')} />
          <Redirect exact from="/hub" to={relativeUrl('/services')} />
        </Switch>
      </AppContextProvider>
    </ApplicationMain>
  );
};
