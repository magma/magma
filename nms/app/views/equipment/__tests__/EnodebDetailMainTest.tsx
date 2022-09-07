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
import type {PromqlReturnObject} from '../../../../generated';

import EnodebContext from '../../../context/EnodebContext';
import EnodebDetail from '../EnodebDetailMain';
import MagmaAPI from '../../../api/MagmaAPI';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {AdapterDateFns} from '@mui/x-date-pickers/AdapterDateFns';
import {EnodebInfo} from '../../../components/lte/EnodebUtils';
import {LocalizationProvider} from '@mui/x-date-pickers';
import {MemoryRouter, Route, Routes} from 'react-router-dom';

import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {mockAPI} from '../../../util/TestUtils';
import {render} from '@testing-library/react';

jest.mock('../../../hooks/useSnackbar');

const mockThroughput: PromqlReturnObject = {
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

describe('<Enodeb />', () => {
  beforeEach(() => {
    mockAPI(
      MagmaAPI.metrics,
      'networksNetworkIdPrometheusQueryRangeGet',
      mockThroughput,
    );
  });

  const enbInfo: Record<string, EnodebInfo> = {
    testEnodebSerial0: {
      enb: {
        attached_gateway_id: 'testGw1',
        config: {
          bandwidth_mhz: 20,
          cell_id: 1,
          device_class: 'Baicells ID TDD/FDD',
          earfcndl: 44290,
          pci: 36,
          special_subframe_pattern: 7,
          subframe_assignment: 2,
          tac: 1,
          transmit_enabled: true,
        },
        enodeb_config: {
          config_type: 'MANAGED',
          managed_config: {
            bandwidth_mhz: 20,
            cell_id: 1,
            device_class: 'Baicells ID TDD/FDD',
            earfcndl: 44290,
            pci: 36,
            special_subframe_pattern: 7,
            subframe_assignment: 2,
            tac: 1,
            transmit_enabled: true,
          },
        },
        name: 'testEnodeb0',
        serial: 'testEnodebSerial0',
      },
      enb_state: {
        enodeb_configured: true,
        enodeb_connected: true,
        fsm_state: 'Completed provisioning eNB. Awaiting new Inform.',
        gps_connected: true,
        gps_latitude: '41.799182',
        gps_longitude: '-88.097308',
        mme_connected: false,
        opstate_enabled: false,
        ptp_connected: false,
        reporting_gateway_id: 'testGw1',
        rf_tx_desired: true,
        rf_tx_on: false,
        time_reported: 0,
        ip_address: '192.168.1.254',
      },
    },
    testEnodebSerial1: {
      enb: {
        attached_gateway_id: 'testGw2',
        config: {
          cell_id: 0,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          transmit_enabled: false,
        },
        enodeb_config: {
          config_type: 'UNMANAGED',
          unmanaged_config: {
            cell_id: 1,
            ip_address: '1.1.1.1',
            tac: 1,
          },
        },
        name: 'testEnodeb0',
        serial: 'testEnodebSerial0',
      },
      enb_state: {
        enodeb_configured: true,
        enodeb_connected: true,
        fsm_state: 'Completed provisioning eNB. Awaiting new Inform.',
        gps_connected: true,
        gps_latitude: '41.799182',
        gps_longitude: '-88.097308',
        mme_connected: false,
        opstate_enabled: false,
        ptp_connected: false,
        reporting_gateway_id: 'testGw2',
        rf_tx_desired: true,
        rf_tx_on: false,
        time_reported: 0,
        ip_address: '192.168.1.254',
      },
    },
  };

  it('managed eNodeB', async () => {
    const Wrapper = () => (
      <MemoryRouter
        initialEntries={['/nms/mynetwork/enodeb/testEnodebSerial0/overview']}
        initialIndex={0}>
        <LocalizationProvider dateAdapter={AdapterDateFns}>
          <StyledEngineProvider injectFirst>
            <ThemeProvider theme={defaultTheme}>
              <ThemeProvider theme={defaultTheme}>
                <EnodebContext.Provider
                  value={{
                    state: {enbInfo: enbInfo},
                    setState: async () => {},
                    refetch: () => {},
                  }}>
                  <Routes>
                    <Route
                      path="/nms/:networkId/enodeb/:enodebSerial/overview/*"
                      element={<EnodebDetail />}
                    />
                  </Routes>
                </EnodebContext.Provider>
              </ThemeProvider>
            </ThemeProvider>
          </StyledEngineProvider>
        </LocalizationProvider>
      </MemoryRouter>
    );
    const {findByTestId, getByTestId} = render(<Wrapper />);
    expect(await findByTestId('eNodeB Serial Number')).toHaveTextContent(
      'testEnodebSerial0',
    );
    expect(getByTestId('eNodeB Externally Managed')).toHaveTextContent('False');
    expect(getByTestId('Health')).toHaveTextContent('Bad');
    expect(getByTestId('Transmit Enabled')).toHaveTextContent('Enabled');
    expect(getByTestId('Gateway ID')).toHaveTextContent('testGw1');
    expect(getByTestId('Mme Connected')).toHaveTextContent('Disconnected');
  });

  it('unManaged eNodeB', async () => {
    const Wrapper = () => (
      <MemoryRouter
        initialEntries={['/nms/mynetwork/enodeb/testEnodebSerial1/overview']}
        initialIndex={0}>
        <LocalizationProvider dateAdapter={AdapterDateFns}>
          <StyledEngineProvider injectFirst>
            <ThemeProvider theme={defaultTheme}>
              <ThemeProvider theme={defaultTheme}>
                <EnodebContext.Provider
                  value={{
                    state: {enbInfo: enbInfo},
                    setState: async () => {},
                    refetch: () => {},
                  }}>
                  <Routes>
                    <Route
                      path="/nms/:networkId/enodeb/:enodebSerial/overview/*"
                      element={<EnodebDetail />}
                    />
                  </Routes>
                </EnodebContext.Provider>
              </ThemeProvider>
            </ThemeProvider>
          </StyledEngineProvider>
        </LocalizationProvider>
      </MemoryRouter>
    );
    const {findByTestId, getByTestId} = render(<Wrapper />);

    expect(await findByTestId('eNodeB Serial Number')).toHaveTextContent(
      'testEnodebSerial1',
    );
    expect(getByTestId('eNodeB Externally Managed')).toHaveTextContent('True');
    expect(getByTestId('Health')).toHaveTextContent('-');
    expect(getByTestId('Transmit Enabled')).toHaveTextContent('Disabled');
    expect(getByTestId('Gateway ID')).toHaveTextContent('testGw2');
    expect(getByTestId('Mme Connected')).toHaveTextContent('Disconnected');
  });
});
