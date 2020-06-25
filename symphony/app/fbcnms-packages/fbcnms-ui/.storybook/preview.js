/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

import {
  addDecorator,
  addParameters,
} from '@storybook/react/dist/client/preview';
import {BrowserRouter} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
import defaultTheme from '@fbcnms/ui/theme/default';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import {themes} from '@storybook/theming';
import Theme from '../theme/symphony';
import {makeStyles} from '@material-ui/styles';
import Story from '../stories/Story';
import {DocsPage, DocsContainer} from '@storybook/addon-docs/blocks';

const useStyles = makeStyles(() => ({
  '@global': {
    body: {
      margin: 0,
    },
  },
}));

addDecorator((story, config) => {
  return (
    <BrowserRouter>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <SnackbarProvider
            maxSnack={3}
            autoHideDuration={10000}
            anchorOrigin={{
              vertical: 'bottom',
              horizontal: 'right',
            }}>
            {story()}
          </SnackbarProvider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </BrowserRouter>
  );
});

addParameters({
  options: {
    name: 'FBC Design System',
    isFullscreen: false,
    showNav: true,
    showPanel: false,
    isToolshown: true,
    theme: {
      ...themes.light,
      appContentBg: Theme.palette.D10,
    },
    storySort: (a, b) => a[1].id.localeCompare(b[1].id),
    hierarchySeparator: /\//,
  },
  docs: {
    container: DocsContainer,
    page: DocsPage,
  },
});
