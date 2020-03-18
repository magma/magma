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
import NavListItem from '@fbcnms/ui/components/NavListItem';
import React, {useContext} from 'react';
import RouterIcon from '@material-ui/icons/Router';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import nullthrows from '@fbcnms/util/nullthrows';
import {Redirect, Route, Switch} from 'react-router-dom';
import {getProjectLinks} from '@fbcnms/magmalte/app/common/projects';
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
      <Text>Service Name</Text>
      <TextInput />
      <br />
      <Button>Create Service</Button>
    </div>
  );
}

function Main() {
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
      <AppContent>
        <CreateServiceForm />
      </AppContent>
    </div>
  );
}

export default () => {
  const relativeUrl = useRelativeUrl();
  return (
    <ApplicationMain>
      <AppContextProvider>
        <Switch>
          <Redirect exact from="/hub" to={relativeUrl('/services')} />
          <Route path="/hub/services" component={Main} />
        </Switch>
      </AppContextProvider>
    </ApplicationMain>
  );
};
