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
import MagmaAPI from '../../api/MagmaAPI';
import React from 'react';
import defaultTheme from '../../theme/default';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {mockAPI} from '../../util/TestUtils';
import {render, waitFor} from '@testing-library/react';
import {subDays, subHours} from 'date-fns';
import type {PromqlReturnObject} from '../../../generated';

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

const testCases = [
  {
    startDate: subHours(new Date(), 4),
    endDate: new Date(),
    step: '15m',
  },
  {
    startDate: subDays(new Date(), 10),
    endDate: new Date(),
    step: '24h',
  },
  {
    startDate: new Date(),
    endDate: subDays(new Date(), 10),
    step: '5m',
  },
];

const Wrapper = (props: {startDate: Date; endDate: Date}) => (
  <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
    <StyledEngineProvider injectFirst>
      <ThemeProvider theme={defaultTheme}>
        <ThemeProvider theme={defaultTheme}>
          <Routes>
            <Route
              path="/nms/:networkId"
              element={
                <EventAlertChart startEnd={[props.startDate, props.endDate]} />
              }
            />
          </Routes>
        </ThemeProvider>
      </ThemeProvider>
    </StyledEngineProvider>
  </MemoryRouter>
);

jest.mock('axios');
jest.mock('../../../app/hooks/useSnackbar');

// chart component was failing here so mocking this out
// this shouldn't affect the prop verification part in the React chart component
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

  it.each(testCases)('renders', async ({startDate, endDate, step}) => {
    render(<Wrapper startDate={startDate} endDate={endDate} />);
    await waitFor(() =>
      expect(
        MagmaAPI.metrics.networksNetworkIdPrometheusQueryRangeGet,
      ).toHaveBeenCalledTimes(1),
    );

    expect(
      MagmaAPI.metrics.networksNetworkIdPrometheusQueryRangeGet,
    ).toBeCalledWith({
      start: startDate.toISOString(),
      end: endDate.toISOString(),
      step: step,
      networkId: 'mynetwork',
      query: 'sum(ALERTS)',
    });
  });
});
