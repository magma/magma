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
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import OrganizationEdit from '../OrganizationEdit';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Organizations from '../Organizations';
import axios from 'axios';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../../theme/default';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AppContextProvider} from '../../../components/context/AppContext';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
// $FlowFixMe[missing-export]
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
    networkIDs: ['test1', 'test', 'test3'],
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
global.CONFIG = {
  appData: {
    enabledFeatures: [],
    tabs: ['nms'],
    user: {
      isSuperUser: true,
    },
  },
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
                <Route
                  path="organizations/detail/:name"
                  element={<OrganizationEdit />}
                />
              </Routes>
            </AppContextProvider>
          </SnackbarProvider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
};

const WrappedOrganizationsEdit = () => {
  return (
    <MemoryRouter
      initialEntries={['/organizations/detail/magma-test']}
      initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <SnackbarProvider>
            <AppContextProvider>
              <Routes>
                <Route
                  path="organizations/detail/:name"
                  element={<OrganizationEdit />}
                />
              </Routes>
            </AppContextProvider>
          </SnackbarProvider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
};

describe('<OrganizationEdit />', () => {
  it('Navigate to Organization Edit', async () => {
    const responses = {
      '/host/organization/async': {data: {organizations: organizationsMock}},
      '/host/networks/async': {
        data: ['test', 'test1', 'test2', 'test3', 'test4'],
      },
      '/host/organization/async/host': {
        data: {organization: {...organizationsMock[0]}},
      },
    };
    mockUseAxios(responses);

    axios.get.mockResolvedValueOnce({
      data: usersMock,
    });
    axios.get.mockResolvedValue({
      data: hostUserMock,
    });

    const {getByTestId, getByText, queryByTestId, getAllByTitle} = render(
      <WrappedOrganizations />,
    );

    await waitFor(() => {
      expect(getByTestId('organizationTitle')).toBeInTheDocument();
    });

    //Onboarding Modal
    expect(getByTestId('onboardingDialog')).not.toBeNull();
    fireEvent.click(getByText('Get Started'));
    expect(queryByTestId('onboardingDialog')).toBeNull();

    // Open menu to go to organization detail page
    const actionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(actionList[0]);
    expect(getByTestId('actions-menu')).toBeVisible();
    await waitFor(() => {
      fireEvent.click(getByText('View'));
    });
    expect(getByTestId('actions-menu')).not.toBeVisible();

    // organizationEdit
    expect(getByTestId('organizationEditTitle')).toHaveTextContent('host');
  });

  it('Updates one organization', async () => {
    const responses = {
      '/host/networks/async': {
        data: ['test', 'test1', 'test2', 'test3', 'test4'],
      },
      '/host/organization/async/magma-test': {
        data: {organization: {...organizationsMock[1]}},
      },
    };
    mockUseAxios(responses);
    axios.put.mockImplementationOnce(() =>
      Promise.resolve({
        status: 200,
        data: null,
      }),
    );
    axios.get.mockResolvedValue({
      data: hostUserMock,
    });

    const {getByTestId, getByText, getAllByText} = render(
      <WrappedOrganizationsEdit />,
    );
    await waitFor(() => {
      expect(getByTestId('organizationEditTitle')).toHaveTextContent(
        'magma-test',
      );
    });

    expect(getByText('test, test1, test2, test3, test4')).toBeVisible();

    // open edit organization dialog
    fireEvent.click(getAllByText('Edit')[0]);
    expect(getByTestId('OrganizationDialog')).toBeVisible();
    expect(getByText('Edit Organization')).toBeVisible();
    const organizationName = getByTestId('name').firstChild;
    expect(organizationName).toBeDisabled();
    fireEvent.click(getByText('Advanced Settings'));
    expect(getByTestId('organizationNetworks')).toBeVisible();
    const organizationNetworks = getByTestId('organizationNetworks').firstChild;

    if (organizationNetworks instanceof HTMLElement) {
      await waitFor(() => {
        fireEvent.mouseDown(organizationNetworks);
      });
      // remove accessible networks
      fireEvent.click(getByText('test1'));
      fireEvent.click(getByText('test2'));
      fireEvent.click(getByText('test3'));
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save'));
    expect(axios.put).toHaveBeenCalledWith(
      '/host/organization/async/magma-test',
      {
        csvCharset: '',
        customDomains: [],
        id: 2,
        name: 'magma-test',
        networkIDs: ['test', 'test4'],
        ssoCert: '',
        ssoEntrypoint: '',
        ssoIssuer: '',
        ssoOidcClientID: '',
        ssoOidcClientSecret: '',
        ssoOidcConfigurationURL: '',
        ssoSelectedType: 'none',
      },
    );

    await waitFor(() => {
      expect(getByText('test, test4')).toBeVisible();
    });

    // add user to organization
    const addUser = getAllByText('Add User');
    fireEvent.click(addUser[0]);
    expect(getByTestId('OrganizationDialog')).toBeVisible();
  });
});
