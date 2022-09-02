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
 */

import * as React from 'react';
import AccountSettings from '../AccountSettings';
import defaultTheme from '../../theme/default';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {SnackbarProvider} from 'notistack';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {fireEvent, render} from '@testing-library/react';

const Wrapper = (props: {children: React.ReactNode}) => (
  <MemoryRouter initialEntries={['/nms/mynetwork/settings']} initialIndex={0}>
    <StyledEngineProvider injectFirst>
      <ThemeProvider theme={defaultTheme}>
        <SnackbarProvider>
          <Routes>
            <Route
              path="/nms/:networkId/settings"
              element={<>{props.children}</>}
            />
          </Routes>
        </SnackbarProvider>
      </ThemeProvider>
    </StyledEngineProvider>
  </MemoryRouter>
);

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
