/**
 * Copyright 2022 The Magma Authors.
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
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import Organizations from '../Organizations';
import axios from 'axios';
import defaultTheme from '../../../theme/default';

import {AppContextProvider} from '../../../components/context/AppContext';
import {EmbeddedData} from '../../../../shared/types/embeddedData';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
import {fireEvent, render, waitFor} from '@testing-library/react';
import {mockUseAxios} from '../useAxiosTestHelper';

jest.mock('axios');
jest.mock('../../../hooks/useAxios');

const organizationsMock = [
  {
    customDomains: [],
    id: 1,
    name: 'host',
    tabs: ['admin'],
    csvCharset: '',
    networkIDs: [],
    ssoSelectedType: 'none',
    ssoCert: '',
    ssoEntrypoint: '',
    ssoIssuer: '',
    ssoOidcClientID: '',
    ssoOidcClientSecret: '',
    ssoOidcConfigurationURL: '',
  },
  {
    customDomains: [],
    id: 2,
    name: 'magma-test',
    tabs: ['nms'],
    csvCharset: '',
    networkIDs: ['test', 'test1', 'test2', 'test3', 'test4'],
    ssoSelectedType: 'none',
    ssoCert: '',
    ssoEntrypoint: '',
    ssoIssuer: '',
    ssoOidcClientID: '',
    ssoOidcClientSecret: '',
    ssoOidcConfigurationURL: '',
  },
];
const usersMock = [
  {
    networkIDs: [],
    tabs: ['nms'],
    id: 1,
    email: 'admin@magma.test',
    organization: 'magma-test',
    role: 3,
    createdAt: '2022-05-10T08:45:04.474Z',
    updatedAt: '2022-05-10T08:45:04.474Z',
  },
  {
    networkIDs: ['test1', 'test', 'test3', 'test5', 'test6'],
    tabs: [],
    id: 4,
    email: 'popo@gmail.com',
    organization: 'magma-test',
    role: 0,
    createdAt: '2022-05-10T08:48:28.979Z',
    updatedAt: '2022-05-11T13:52:54.572Z',
  },
  {
    networkIDs: [],
    tabs: [],
    id: 8,
    email: 'testi@gmail.com',
    organization: 'magma-test',
    role: 0,
    createdAt: '2022-05-11T14:38:08.969Z',
    updatedAt: '2022-05-11T14:38:08.969Z',
  },
];

const hostUserMock = [
  {
    networkIDs: [],
    tabs: ['nms'],
    id: 2,
    email: 'admin@magma.test',
    organization: 'host',
    role: 3,
    createdAt: '2022-05-10T08:45:08.134Z',
    updatedAt: '2022-05-10T08:45:08.134Z',
  },
  {
    networkIDs: [],
    tabs: [],
    id: 3,
    email: 'test@gmail.com',
    organization: 'host',
    role: 0,
    createdAt: '2022-05-10T08:46:38.237Z',
    updatedAt: '2022-05-10T08:46:38.237Z',
  },
];
window.CONFIG = {
  appData: ({
    enabledFeatures: [],
    tabs: ['nms'],
    user: {
      isSuperUser: true,
    },
  } as unknown) as EmbeddedData,
};
const WrappedOrganizations = () => {
  return (
    <MemoryRouter initialEntries={['/organizations']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <SnackbarProvider>
            <AppContextProvider>
              <Routes>
                <Route path="organizations/*" element={<Organizations />} />
              </Routes>
            </AppContextProvider>
          </SnackbarProvider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
};

describe('<Organizations />', () => {
  it('renders with no organizations', async () => {
    const responses = {
      '/host/organization/async': {data: {organizations: []}},
    };
    mockUseAxios(responses);
    const {getByTestId, getByText, queryByTestId} = render(
      <WrappedOrganizations />,
    );
    await waitFor(() => {
      expect(getByTestId('organizationTitle')).toBeInTheDocument();
    });

    //Onboarding Modal
    expect(getByTestId('onboardingDialog')).not.toBeNull();
    fireEvent.click(getByText('Get Started'));
    await waitFor(() => {
      expect(queryByTestId('onboardingDialog')).toBeNull();
    });
  });

  it('renders with 2 organizations', async () => {
    const responses = {
      '/host/organization/async': {data: {organizations: organizationsMock}},
      '/host/networks/async': {data: []},
    };
    mockUseAxios(responses);
    (axios.get as jest.Mock).mockResolvedValueOnce({
      data: hostUserMock,
    });
    (axios.get as jest.Mock).mockResolvedValueOnce({
      data: usersMock,
    });
    const {getByTestId, getByText, queryByTestId, getAllByRole} = render(
      <WrappedOrganizations />,
    );

    await waitFor(() => {
      expect(getByTestId('organizationTitle')).toBeInTheDocument();
    });

    //Onboarding Modal
    expect(getByTestId('onboardingDialog')).not.toBeNull();
    fireEvent.click(getByText('Get Started'));
    await waitFor(() => {
      expect(queryByTestId('onboardingDialog')).toBeNull();
    });
    const rowItems = getAllByRole('row');

    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Name');
    expect(rowItems[0]).toHaveTextContent('Accessible Networks');
    expect(rowItems[0]).toHaveTextContent('Link to Organization Portal');
    expect(rowItems[0]).toHaveTextContent('Number of Users');

    expect(rowItems[1]).toHaveTextContent('host');
    expect(rowItems[1]).toHaveTextContent('-');
    expect(rowItems[1]).toHaveTextContent('Visit host Organization Portal');
    expect(rowItems[1]).toHaveTextContent('2');

    expect(rowItems[2]).toHaveTextContent('magma-test');
    expect(rowItems[2]).toHaveTextContent('test, test1, test2 + 2 more');
    expect(rowItems[2]).toHaveTextContent(
      'Visit magma-test Organization Portal',
    );
    expect(rowItems[2]).toHaveTextContent('3');
  });
});
