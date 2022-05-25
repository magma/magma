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

import MagmaAPIBindings from '../../../../generated/MagmaAPIBindings';
import MomentUtils from '@date-io/moment';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
// $FlowFixMe migrated to typescript
import NetworkContext from '../../../components/context/NetworkContext';
import React from 'react';
import SubscriberChart from '../SubscriberChart';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../../theme/default';

import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiPickersUtilsProvider} from '@material-ui/pickers';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {render, wait} from '@testing-library/react';

jest.mock('axios');
jest.mock('../../../../generated/MagmaAPIBindings.js');
jest.mock('../../../hooks/useSnackbar');

const mockAvgCurDataUsage = {
  status: 'success',
  data: {
    resultType: 'vector',
    result: [
      {
        metric: {},
        value: [1627325883.103, '12108000.691521779536'],
      },
    ],
  },
};

const mockAvgDailyDataUsage = {
  status: 'success',
  data: {
    resultType: 'vector',
    result: [
      {
        metric: {},
        value: [1627325883.103, '2210400.691521779536'],
      },
    ],
  },
};

const mockAvgMonthlyDataUsage = {
  status: 'success',
  data: {
    resultType: 'vector',
    result: [
      {
        metric: {},
        value: [1627325883.103, '52108.691521779536'],
      },
    ],
  },
};

const mockEmptyDataset = {
  status: 'success',
  data: {
    metric: {},
    resultType: 'vector',
    result: [],
  },
};

describe('<SubscriberChart />', () => {
  beforeEach(() => {
    // Called by the chart component
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mockResolvedValue(
      mockEmptyDataset,
    );
    // avg data usage in the past 1 hour
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQuery
      .mockReturnValueOnce(mockAvgCurDataUsage)
      // avg data usage in the past 1 day
      .mockReturnValueOnce(mockAvgDailyDataUsage)
      // avg data usage in the past 1 month
      .mockReturnValueOnce(mockAvgMonthlyDataUsage)
      // avg data usage in the past 1 year
      .mockReturnValueOnce(mockEmptyDataset);
  });

  const Wrapper = () => {
    return (
      <MemoryRouter
        initialEntries={[
          '/nms/test/subscribers/overview/config/IMSI001011234560000/overview',
        ]}
        initialIndex={0}>
        <MuiPickersUtilsProvider utils={MomentUtils}>
          <MuiThemeProvider theme={defaultTheme}>
            <MuiStylesThemeProvider theme={defaultTheme}>
              <NetworkContext.Provider
                value={{
                  networkId: 'test',
                }}>
                <Routes>
                  <Route
                    path="/nms/:networkId/subscribers/overview/config/:subscriberId/overview"
                    element={<SubscriberChart />}
                  />
                </Routes>
              </NetworkContext.Provider>
            </MuiStylesThemeProvider>
          </MuiThemeProvider>
        </MuiPickersUtilsProvider>
      </MemoryRouter>
    );
  };

  it('Verify Subscribers Data KPI', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();
    expect(getByTestId('Hourly Usage MB/s')).toHaveTextContent('12.11');
    expect(getByTestId('Daily Avg MB/s')).toHaveTextContent('2.21');
    expect(getByTestId('Monthly Avg Mb/s')).toHaveTextContent('0.05');
    expect(getByTestId('Yearly Avg Mb/s')).toHaveTextContent('0.00');
  });
});
