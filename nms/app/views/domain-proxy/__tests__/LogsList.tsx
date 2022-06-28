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

import LogsList from '../LogsList';
import MomentUtils from '@date-io/moment';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiPickersUtilsProvider} from '@material-ui/pickers';

import MagmaAPI from '../../../../api/MagmaAPI';
import moment from 'moment';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {
  fireEvent,
  render,
  screen,
  waitFor,
  within,
} from '@testing-library/react';
import {mockAPI} from '../../../util/TestUtils';

const mockEnqueueSnackbar = jest.fn();
jest.mock('../../../hooks/useSnackbar', () => ({
  useEnqueueSnackbar: () => mockEnqueueSnackbar,
}));

const networkId = 'test-network';

const renderWithProviders = (jsx: React.ReactNode) => {
  return render(
    <MemoryRouter
      initialEntries={[`/nms/${networkId}/metrics/domain-proxy-logs`]}
      initialIndex={0}>
      <MuiPickersUtilsProvider utils={MomentUtils}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <Routes>
              <Route
                path="/nms/:networkId/metrics/domain-proxy-logs"
                element={jsx}
              />
            </Routes>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
        ,
      </MuiPickersUtilsProvider>
    </MemoryRouter>,
  );
};

describe('<LogsList />', () => {
  describe('Filtering', () => {
    let getLogsMock: jest.SpyInstance;

    const expectApiCallParam = async (param: string, value: unknown) => {
      await waitFor(() =>
        // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
        expect(getLogsMock.mock.calls[0][0][param]).toBe(value),
      );
    };

    const fillInput = (testId: string, value: unknown) => {
      fireEvent.change(screen.getByTestId(testId), {target: {value}});
    };

    // See https://stackoverflow.com/a/61491607
    const fillMuiSelect = (testId: string, optionText: string) => {
      const select = screen.getByTestId(testId);
      fireEvent.mouseDown(within(select).getByRole('button'));
      const listbox = within(screen.getByRole('listbox'));
      fireEvent.click(listbox.getByText(new RegExp(`^${optionText}`, 'i')));
    };

    const clickSearchButton = () => {
      fireEvent.click(screen.getByTestId('search-button'));
    };

    const filterValues = {
      serialNumber: 'test-serial',
      fccId: 'test-fcc-id',
      logDirectionFrom: 'CBSD',
      startDate: '2000/01/01 00:00',
      responseCode: 422,
      logName: 'test-log-name',
      logDirectionTo: 'DP',
      endDate: '2000/05/05 00:00',
    };

    beforeEach(async () => {
      getLogsMock = mockAPI(MagmaAPI.logs, 'dpNetworkIdLogsGet');
      renderWithProviders(<LogsList />);

      // Wait for initial request after mount and clear it
      // So we can test calls caused by clicking the search button
      await expectApiCallParam('offset', 0);
      getLogsMock.mockClear();
    });

    it('Sends serial number', async () => {
      fillInput('serial-number-input', filterValues.serialNumber);
      clickSearchButton();
      await expectApiCallParam('serialNumber', filterValues.serialNumber);
    });

    it('Sends fcc id', async () => {
      fillInput('fcc-id-input', filterValues.fccId);
      clickSearchButton();
      await expectApiCallParam('fccId', filterValues.fccId);
    });

    it('Sends logs direction from', async () => {
      fillMuiSelect('logs-direction-from-input', filterValues.logDirectionFrom);
      clickSearchButton();
      await expectApiCallParam('from', filterValues.logDirectionFrom);
    });

    it('Sends start date', async () => {
      fillInput('start-date-input', filterValues.startDate);
      clickSearchButton();
      await expectApiCallParam(
        'begin',
        moment(filterValues.startDate).toISOString(),
      );
    });

    it('Sends responseCode', async () => {
      fillInput('response-code-input', `${filterValues.responseCode}`);
      clickSearchButton();
      await expectApiCallParam('responseCode', filterValues.responseCode);
    });

    it('Sends log name', async () => {
      fillInput('log-name-input', filterValues.logName);
      clickSearchButton();
      await expectApiCallParam('type', filterValues.logName);
    });

    it('Sends logs direction to', async () => {
      fillMuiSelect('logs-direction-to-input', filterValues.logDirectionTo);
      clickSearchButton();
      await expectApiCallParam('to', filterValues.logDirectionTo);
    });

    it('Sends end date', async () => {
      fillInput('end-date-input', filterValues.endDate);
      clickSearchButton();
      await expectApiCallParam(
        'end',
        moment(filterValues.endDate).toISOString(),
      );
    });
  });
});
