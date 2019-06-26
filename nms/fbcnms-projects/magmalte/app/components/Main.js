/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import AppContent from '@fbcnms/ui/components/layout/AppContent';
import AppContext from './context/AppContext';
import AppDrawer from '@fbcnms/ui/components/layout/AppDrawer';
import CssBaseline from '@material-ui/core/CssBaseline';
import defaultTheme from '@fbcnms/ui/theme/default';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import {MuiThemeProvider} from '@material-ui/core/styles';
import NetworkContext from './context/NetworkContext';
import {Redirect, Route, Switch} from 'react-router-dom';
import SectionLinks from './layout/SectionLinks';
import SectionRoutes from './layout/SectionRoutes';
import {SnackbarProvider} from 'notistack';
import {TopBarContextProvider} from '@fbcnms/ui/components/layout/TopBarContext';
import VersionTooltip from './VersionTooltip';

import {MagmaAPIUrls} from '../common/MagmaAPI';
import {hot} from 'react-hot-loader';
import {makeStyles} from '@material-ui/styles';
import {sortBy} from 'lodash';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  toolbarIcon: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar,
  },
}));

// These won't be considered networkIds
const ROOT_PATHS = new Set(['network']);

function Index() {
  const classes = useStyles();
  const {match} = useRouter();
  const networkId = ROOT_PATHS.has(match.params.networkId)
    ? null
    : match.params.networkId;

  return (
    <NetworkContext.Provider value={{networkId}}>
      <div className={classes.root}>
        <AppDrawer>
          <SectionLinks />
          <VersionTooltip />
        </AppDrawer>
        <AppContent>
          <SectionRoutes />
        </AppContent>
      </div>
    </NetworkContext.Provider>
  );
}

function Main() {
  const {match} = useRouter();
  const {response, error} = useAxios({
    method: 'get',
    url: MagmaAPIUrls.networks(),
  });

  const networkIds = sortBy(response?.data) || ['mpk_test'];
  const appContext = {
    ...window.CONFIG.appData,
    networkIds,
  };

  if (networkIds.length > 0 && !match.params.networkId) {
    return <Redirect to={`/nms/${networkIds[0]}/map/`} />;
  }

  if (
    response &&
    !error &&
    networkIds.length === 0 &&
    match.params.networkId !== 'network'
  ) {
    return <Redirect to="/nms/network/create" />;
  }

  return (
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <SnackbarProvider
          maxSnack={3}
          autoHideDuration={10000}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'right',
          }}>
          <AppContext.Provider value={appContext}>
            <TopBarContextProvider>
              <CssBaseline />
              <Index />
            </TopBarContextProvider>
          </AppContext.Provider>
        </SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  );
}

/* eslint-disable-next-line no-undef */
export default hot(module)(() => (
  <Switch>
    <Route path="/nms/:networkId" component={Main} />
    <Route path="/nms" component={Main} />
  </Switch>
));
