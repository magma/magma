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

import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {AddEditCbsdButton, CbsdAddEditDialog} from '../CbsdEdit';

import MagmaAPI from '../../../../api/MagmaAPI';
import {CbsdContextProvider} from '../../../components/lte/LteContext';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {fireEvent, render, waitFor, within} from '@testing-library/react';
import {mockAPI, mockAPIError} from '../../../util/TestUtils';
import type {Cbsd, MutableCbsd} from '../../../../generated-ts';

const mockEnqueueSnackbar = jest.fn();
jest.mock('../../../hooks/useSnackbar', () => ({
  useEnqueueSnackbar: () => mockEnqueueSnackbar,
}));

const mockCbsd: Cbsd = {
  capabilities: {
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
  id: 28,
  installation_param: {
    antenna_gain: 0,
  },
  is_active: false,
  serial_number: '1202000291213VB0009',
  single_step_enabled: false,
  state: 'unregistered',
  user_id: 'SAS-Freedomfi',
};

const networkId = 'test-network';

const convertCbsdToMutableCbsd = (cbsdToConvert: Cbsd): MutableCbsd => {
  const {is_active, id, state, cbsd_id, ...payload} = cbsdToConvert;
  return payload;
};

const renderWithProviders = (jsx: React.ReactNode) => {
  return render(
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <CbsdContextProvider networkId={networkId} networkType="test">
          {jsx}
        </CbsdContextProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>,
  );
};

describe('<AddEditCbsdButton />', () => {
  beforeEach(() => {
    mockAPI(MagmaAPI.cbsds, 'dpNetworkIdCbsdsGet', {
      cbsds: [mockCbsd],
      total_count: 1,
    });
  });

  it('Shows Add new CBSD dialog when clicked', async () => {
    const {findByRole, findByText} = renderWithProviders(
      <AddEditCbsdButton title="test" />,
    );

    const button = await findByRole('button');

    fireEvent.click(button);

    await findByText('Add New CBSD');
  });
});

describe('<CbsdAddEditDialog />', () => {
  beforeEach(() => {
    mockAPI(MagmaAPI.cbsds, 'dpNetworkIdCbsdsGet', {
      cbsds: [mockCbsd],
      total_count: 1,
    });
    mockAPI(MagmaAPI.cbsds, 'dpNetworkIdCbsdsCbsdIdPut');
    mockAPI(MagmaAPI.cbsds, 'dpNetworkIdCbsdsPost');
  });

  it('Shows "Add New CBSD" text when rendered without cbsd', async () => {
    const {findByText} = renderWithProviders(
      <CbsdAddEditDialog open={true} onClose={() => {}} />,
    );
    await findByText('Add New CBSD');
  });

  it('Shows "Edit CBSD" text when rendered with a cbsd', async () => {
    const {findByText} = renderWithProviders(
      <CbsdAddEditDialog open={true} onClose={() => {}} cbsd={mockCbsd} />,
    );
    await findByText('Edit CBSD');
  });

  it('Calls putDpByNetworkIdCbsdsByCbsdId() and shows success snackbar when CBSD is edited', async () => {
    const {findByTestId} = renderWithProviders(
      <CbsdAddEditDialog open={true} onClose={() => {}} cbsd={mockCbsd} />,
    );

    const updatedSerial = 'test-changed-serial';

    const input = await findByTestId('serial-number-input');

    fireEvent.change(input, {target: {value: updatedSerial}});

    const button = await findByTestId('save-cbsd-button');

    fireEvent.click(button);

    await waitFor(() =>
      expect(
        mockEnqueueSnackbar,
      ).toHaveBeenCalledWith('CBSD saved successfully', {variant: 'success'}),
    );

    const expectedCbsdPayload = convertCbsdToMutableCbsd({
      ...mockCbsd,
      serial_number: updatedSerial,
    });

    await waitFor(() =>
      expect(MagmaAPI.cbsds.dpNetworkIdCbsdsCbsdIdPut).toHaveBeenCalledWith({
        cbsd: expectedCbsdPayload,
        networkId,
        cbsdId: mockCbsd.id,
      }),
    );
  });

  it('Shows error snackbar when putDpByNetworkIdCbsdsByCbsdId() throws error saving CBSD', async () => {
    const {findByTestId} = renderWithProviders(
      <CbsdAddEditDialog open={true} onClose={() => {}} cbsd={mockCbsd} />,
    );

    const response = {
      code: 422,
      errors: [
        {
          code: 602,
          in: 'body',
          message: 'serial_number in body is required',
          name: 'serial_number',
          value: '',
          values: null,
        },
      ],
      message: 'validation failure list',
    };

    mockAPIError(MagmaAPI.cbsds, 'dpNetworkIdCbsdsCbsdIdPut', response);

    const button = await findByTestId('save-cbsd-button');

    fireEvent.click(button);

    await waitFor(() =>
      expect(mockEnqueueSnackbar).toHaveBeenCalledWith(
        'failed to update CBSD',
        {variant: 'error'},
      ),
    );
  });

  it('Calls postDpByNetworkIdCbsds() and shows success snackbar when CBSD is created', async () => {
    const {getByTestId, findByTestId, getByRole} = renderWithProviders(
      <CbsdAddEditDialog open={true} onClose={() => {}} />,
    );

    const fillInput = (testId: string, value: unknown) => {
      fireEvent.change(getByTestId(testId), {target: {value}});
    };

    // See https://stackoverflow.com/a/61491607
    const fillMuiSelect = (testId: string, optionText: string | number) => {
      const select = getByTestId(testId);
      fireEvent.mouseDown(within(select).getByRole('button'));
      const listbox = within(getByRole('listbox'));
      fireEvent.click(listbox.getByText(new RegExp(`^${optionText}`, 'i')));
    };

    fillInput('serial-number-input', mockCbsd.serial_number);
    fillInput('fcc-id-input', mockCbsd.fcc_id);
    fillInput('user-id-input', mockCbsd.user_id);
    fillInput('min-power-input', mockCbsd.capabilities.min_power);
    fillInput('max-power-input', mockCbsd.capabilities.max_power);
    fillInput('max-power-input', mockCbsd.capabilities.max_power);
    fillInput(
      'number-of-antennas-input',
      mockCbsd.capabilities.number_of_antennas,
    );
    fillInput('antenna-gain-input', mockCbsd.installation_param.antenna_gain!);
    fillMuiSelect('desired-state-input', mockCbsd.desired_state);
    fillMuiSelect('cbsd-category-input', mockCbsd.cbsd_category);
    fillInput('single-step-enabled-input', mockCbsd.single_step_enabled);
    fillMuiSelect(
      'bandwidth-input',
      mockCbsd.frequency_preferences.bandwidth_mhz,
    );
    fillInput(
      'frequencies-input',
      mockCbsd.frequency_preferences.frequencies_mhz,
    );

    const button = await findByTestId('save-cbsd-button');

    fireEvent.click(button);

    await waitFor(() =>
      expect(
        mockEnqueueSnackbar,
      ).toHaveBeenCalledWith('CBSD saved successfully', {variant: 'success'}),
    );

    const expectedCbsdPayload = convertCbsdToMutableCbsd(mockCbsd);

    await waitFor(() =>
      expect(MagmaAPI.cbsds.dpNetworkIdCbsdsPost).toHaveBeenCalledWith({
        cbsd: expectedCbsdPayload,
        networkId,
      }),
    );
  });
});
