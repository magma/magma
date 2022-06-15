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

import EventAlertChart from '../EventAlertChart';
import MagmaAPI from '../../../api/MagmaAPI';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../theme/default';
import moment from 'moment';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {mockAPI} from '../../util/TestUtils';
import {render, wait} from '@testing-library/react';
import type {PromqlReturnObject} from '../../../generated-ts';

const mockMetricSt: PromqlReturnObject = {
  status: 'success',
  data: {
    resultType: 'matrix',
    result: [
      {
        metric: {},
        values: [['1588898968.042', '6']],
      },
    ],
  },
};

jest.mock('axios');
jest.mock('../../../app/hooks/useSnackbar');

// chart component was failing here so mocking this out
// this shouldn't affect the prop verification part in the react
// chart component
// @ts-ignore
window.HTMLCanvasElement.prototype.getContext = () => {};

describe('<EventAlertChart/>', () => {
  beforeEach(() => {
    mockAPI(
      MagmaAPI.metrics,
      'networksNetworkIdPrometheusQueryRangeGet',
      mockMetricSt,
    );
    mockAPI(MagmaAPI.events, 'eventsNetworkIdAboutCountGet');
  });

  const testCases = [
    {
      startDate: moment().subtract(2, 'hours'),
      endDate: moment(),
      step: '15m',
      valid: true,
    },
    {
      startDate: moment().subtract(10, 'day'),
      endDate: moment(),
      step: '24h',
      valid: true,
    },
    {
      startDate: moment(),
      endDate: moment().subtract(10, 'day'),
      step: '24h',
      valid: false,
    },
  ];

  it.each(testCases)('renders', async tc => {
    // const endDate = moment();
    // const startDate = moment().subtract(3, 'hours');
    const Wrapper = () => (
      <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <Routes>
              <Route
                path="/nms/:networkId"
                element={
                  <EventAlertChart startEnd={[tc.startDate, tc.endDate]} />
                }
              />
            </Routes>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );

    render(<Wrapper />);
    await wait();

    const currentStep = tc.valid ? tc.step : '5m';
    expect(
      MagmaAPI.metrics.networksNetworkIdPrometheusQueryRangeGet,
    ).toHaveBeenCalledTimes(1);
    expect(
      MagmaAPI.metrics.networksNetworkIdPrometheusQueryRangeGet,
    ).toBeCalledWith({
      start: tc.startDate.toISOString(),
      end: tc.endDate.toISOString(),
      step: currentStep,
      networkId: 'mynetwork',
      query: 'sum(ALERTS)',
    });
  });
});
