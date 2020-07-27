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
 * @flow strict-local
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
