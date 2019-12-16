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
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import {StylesProvider} from '@material-ui/styles';
import {createMuiTheme} from '@material-ui/core/styles';

function MaterialTheme({children}: {children: React.Node}) {
  const theme = createMuiTheme({});
  return (
    <StylesProvider>
      <MuiStylesThemeProvider theme={theme}>{children}</MuiStylesThemeProvider>
    </StylesProvider>
  );
}

export function TestWrapper({children, ...props}: {children: React.Node}) {
  return <MaterialTheme {...props} children={children} />;
}
