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
import CbsdContext, {
  CbsdContextType,
} from '../../../components/context/CbsdContext';
import CbsdsTable from '../CbsdsTable';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';

import {MuiThemeProvider} from '@material-ui/core/styles';
import {fireEvent, render, within} from '@testing-library/react';

jest.mock('axios');

jest.mock('../../../hooks/useSnackbar');

const cbsds = [
  {
    capabilities: {
      antenna_gain: 0,
      max_power: 24,
      min_power: 0,
      number_of_antennas: 1,
      max_ibw_mhz: 150,
    },
    carrier_aggregation_enabled: false,
    grant_redundancy: true,
    cbsd_category: 'b',
    cbsd_id: '2AG32PBS31010/1202000291213VB0009',
    desired_state: 'registered',
    fcc_id: '2AG32PBS31010',
    frequency_preferences: {
      bandwidth_mhz: 20,
      frequencies_mhz: [3600],
    },
    installation_param: {
      antenna_gain: 4,
      height_m: 8,
      height_type: 'agl',
      indoor_deployment: true,
      latitude_deg: 40.019393,
    },
    id: 28,
    is_active: false,
    serial_number: '1202000291213VB0009',
    single_step_enabled: false,
    state: 'unregistered',
    user_id: 'SAS-Freedomfi',
  },
  {
    capabilities: {
      antenna_gain: 0,
      max_power: 0,
      min_power: 0,
      number_of_antennas: 1,
      max_ibw_mhz: 150,
    },
    carrier_aggregation_enabled: false,
    grant_redundancy: true,
    cbsd_category: 'b',
    desired_state: 'unregistered',
    fcc_id: 'test',
    frequency_preferences: {
      bandwidth_mhz: 5,
      frequencies_mhz: [3555],
    },
    id: 30,
    installation_param: {
      indoor_deployment: false,
    },
    is_active: false,
    serial_number: 'test-serial-number-2',
    single_step_enabled: false,
    state: 'unregistered',
    user_id: 'test',
  },
];

const cbsdState = {
  state: {
    isLoading: false,
    totalCount: 2,
    pageSize: 10,
    page: 0,
    cbsds,
  },
  setPaginationOptions: jest.fn(),
  refetch: jest.fn(),
  create: jest.fn(),
  update: jest.fn(),
  deregister: jest.fn(),
  remove: jest.fn(),
} as CbsdContextType;

const renderTable = () => {
  return render(
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <CbsdContext.Provider value={cbsdState}>
          <CbsdsTable />
        </CbsdContext.Provider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>,
  );
};

describe('<CbsdsTable /> with 2 cbsds', () => {
  it('Shows "Serial Number" column', async () => {
    const {findByText} = renderTable();
    await findByText('Serial Number');
  });

  it('Shows cbds serial numbers', async () => {
    const {findByText} = renderTable();
    await findByText(cbsds[0].serial_number);
    await findByText(cbsds[1].serial_number);
  });

  it('Opens edit modal when Edit button is clicked', async () => {
    const component = renderTable();

    const tables = await component.findAllByRole('table');
    const table = tables[0];

    const rows = await within(table).findAllByRole('row');
    // 0 is headings, 1 is the first row
    const row = rows[1];

    const menuButton = await within(row).findByRole(
      (role, element) => element?.getAttribute('title') === 'Actions',
    );

    fireEvent.click(menuButton);

    const menu = await component.findByTestId('actions-menu');

    const editButton = await within(menu).findByText('Edit');

    fireEvent.click(editButton);

    const dialog = await component.findByRole('dialog');

    expect(within(dialog).getByText('Edit CBSD')).toBeInTheDocument();
  });
});
