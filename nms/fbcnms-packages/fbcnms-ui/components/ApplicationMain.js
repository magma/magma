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
import CssBaseline from '@material-ui/core/CssBaseline';
import ThemeProvider from '@material-ui/styles/ThemeProvider';
import defaultTheme from '@fbcnms/ui/theme/default';
import {SnackbarProvider} from 'notistack';
import {TopBarContextProvider} from '@fbcnms/ui/components/layout/TopBarContext';

type Props = {
  children: React.Element<*>,
};

const ApplicationMain = (props: Props) => {
  return (
    <ThemeProvider theme={defaultTheme}>
      <SnackbarProvider
        maxSnack={3}
        autoHideDuration={10000}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}>
        <TopBarContextProvider>
          <CssBaseline />
          {props.children}
        </TopBarContextProvider>
      </SnackbarProvider>
    </ThemeProvider>
  );
};

export default ApplicationMain;
