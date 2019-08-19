/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

import {addDecorator} from '@storybook/react/dist/client/preview';
import {BrowserRouter} from 'react-router-dom';
import {configure} from '@storybook/react';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
import defaultTheme from '@fbcnms/ui/theme/default';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';

// automatically import all files ending in *.stories.js
const req = require.context('../stories', true, /.stories.js$/);
function loadStories() {
  req.keys().forEach(filename => req(filename));
}

addDecorator(story => (
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
));

configure(loadStories, module);
