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
 * @flow
 * @format
 */

import 'jest-dom/extend-expect';
import * as React from 'react';
import AccountSettings from '../AccountSettings';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import defaultTheme from '../../theme/default';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
import {cleanup, fireEvent, render} from '@testing-library/react';

const Wrapper = (props: {children: React.Node}) => (
  <MemoryRouter initialEntries={['/nms/mynetwork/settings']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <SnackbarProvider>
          <Routes>
            <Route
              path="/nms/:networkId/settings"
              element={<>{props.children}</>}
            />
          </Routes>
        </SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

afterEach(cleanup);

describe('<AccountSettings />', () => {
  it('Save button is disabled if form is not filled-out', () => {
    const {getByRole, getByPlaceholderText} = render(
      <Wrapper>
        <AccountSettings />
      </Wrapper>,
    );

    const button = getByRole('button', {name: 'Save'});
    expect(button).toBeDisabled();

    fireEvent.change(getByPlaceholderText('Enter Current Password'), {
      target: {value: '1234'},
    });
    fireEvent.change(getByPlaceholderText('Enter New Password'), {
      target: {value: 'secret'},
    });
    expect(button).toBeDisabled();

    fireEvent.change(getByPlaceholderText('Confirm New Password'), {
      target: {value: 'secret'},
    });
    expect(button).not.toBeDisabled();
  });
});
