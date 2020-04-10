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
import {configure} from '@storybook/react';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
import defaultTheme from '@fbcnms/ui/theme/default';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import {themes} from '@storybook/theming';
import Theme from '../theme/symphony';
import {compareStoriesName} from '../stories/storybookUtils';
import {makeStyles} from '@material-ui/styles';
import Story from '../stories/Story';

// automatically import all files ending in *.stories.js
const req = require.context('../stories', true, /.stories.js$/);
function loadStories() {
  const designSystemStories = [
    './foundation/colors.stories.js',
    './foundation/shadows.stories.js',
    './foundation/typography.stories.js',
    './inputs/text-input.stories.js',
    './inputs/form-field.stories.js',
    './components/card.stories.js',
    './containers/viewHeader.stories.js',
    './containers/sideMenu.stories.js',
    './containers/navigatableViews.stories.js',
  ];

  designSystemStories.map(story => req(story));
  req
    .keys()
    .filter(story => !designSystemStories.includes(story))
    .forEach(filename => req(filename));
}

const useStyles = makeStyles(() => ({
  '@global': {
    body: {
      margin: 0,
    },
  },
}));

addDecorator((story, config) => {
  const _classes = useStyles();
  console.log('config', config);
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
            <Story name={config.name}>{story()}</Story>
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
    hierarchySeparator: /\//,
  },
});

configure(loadStories, module);
