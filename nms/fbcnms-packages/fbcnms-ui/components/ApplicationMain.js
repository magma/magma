/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import AppContext from '@fbcnms/ui/context/AppContext';
import CssBaseline from '@material-ui/core/CssBaseline';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import MuiThemeProvider from '@material-ui/core/styles/MuiThemeProvider';
import defaultTheme from '@fbcnms/ui/theme/default';
import {SnackbarProvider} from 'notistack';
import {TopBarContextProvider} from '@fbcnms/ui/components/layout/TopBarContext';

type Props = {
  appContext: any,
  children: React.Element<*>,
};

const ApplicationMain = (props: Props) => {
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
          <AppContext.Provider value={props.appContext}>
            <TopBarContextProvider>
              <CssBaseline />
              {props.children}
            </TopBarContextProvider>
          </AppContext.Provider>
        </SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  );
};

export default ApplicationMain;
