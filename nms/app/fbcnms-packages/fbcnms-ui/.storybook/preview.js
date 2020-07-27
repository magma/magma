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
    name: 'Symphony Design System',
    isFullscreen: false,
    showNav: true,
    showPanel: false,
    isToolshown: true,
    selectedPanel: 'storybook/docs/panel',
    theme: {
      ...themes.light,
      colorPrimary: Theme.palette.D900,
      colorSecondary: Theme.palette.B600,
      appBg: Theme.palette.D10,
      appContentBg: Theme.palette.D10,
      appBorderColor: Theme.palette.D300,
      appBorderRadius: 4,
      fontBase: '"Roboto"',
      textColor: Theme.palette.D900,
      textInverseColor: Theme.palette.white,
      barTextColor: Theme.palette.D400,
      barSelectedColor: Theme.palette.primary,
      barBg: Theme.palette.white,
      inputBg: Theme.palette.D50,
      inputBorder: Theme.palette.D100,
      inputTextColor: Theme.palette.D900,
      inputBorderRadius: 4,
      brandTitle: 'Symphony Design System',
    },
    storySort: (a, b) => a[1].id.localeCompare(b[1].id),
    hierarchySeparator: /\//,
  },
  docs: {
    container: DocsContainer,
    page: DocsPage,
  },
});
