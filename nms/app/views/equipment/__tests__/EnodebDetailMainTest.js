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
import type {promql_return_object} from '../../../../generated/MagmaAPIBindings';

import * as hooks from '../../../components/context/RefreshContext';
// $FlowFixMe migrated to typescript
import EnodebContext from '../../../components/context/EnodebContext';
import EnodebDetail from '../EnodebDetailMain';
import MagmaAPIBindings from '../../../../generated/MagmaAPIBindings';
import MomentUtils from '@date-io/moment';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';

import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiPickersUtilsProvider} from '@material-ui/pickers';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {render, wait} from '@testing-library/react';

jest.mock('axios');
jest.mock('../../../../generated/MagmaAPIBindings.js');
jest.mock('../../../hooks/useSnackbar');

const mockThroughput: promql_return_object = {
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
    // eslint-disable-next-line max-len
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mockResolvedValue(
      mockThroughput,
    );
    MagmaAPIBindings.getNetworks.mockResolvedValue([]);
  });

  const enbInfo = {
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
    jest.spyOn(hooks, 'useRefreshingContext').mockImplementation(() => ({
      enbInfo: {testEnodebSerial0: enbInfo['testEnodebSerial0']},
    }));
    const Wrapper = () => (
      <MemoryRouter
        initialEntries={['/nms/mynetwork/enodeb/testEnodebSerial0/overview']}
        initialIndex={0}>
        <MuiPickersUtilsProvider utils={MomentUtils}>
          <MuiThemeProvider theme={defaultTheme}>
            <MuiStylesThemeProvider theme={defaultTheme}>
              <EnodebContext.Provider
                value={{
                  state: {enbInfo: enbInfo},
                  setState: async _ => {},
                }}>
                <Routes>
                  <Route
                    path="/nms/:networkId/enodeb/:enodebSerial/overview/*"
                    element={<EnodebDetail enbInfo={enbInfo} />}
                  />
                </Routes>
              </EnodebContext.Provider>
            </MuiStylesThemeProvider>
          </MuiThemeProvider>
        </MuiPickersUtilsProvider>
      </MemoryRouter>
    );
    const {getByTestId} = render(<Wrapper />);
    await wait();

    // TODO - commenting this out till we have per enodeb metric support
    // expect(
    //   MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange,
    // ).toHaveBeenCalledTimes(2);
    expect(getByTestId('eNodeB Serial Number')).toHaveTextContent(
      'testEnodebSerial0',
    );
    expect(getByTestId('eNodeB Externally Managed')).toHaveTextContent('False');
    expect(getByTestId('Health')).toHaveTextContent('Bad');
    expect(getByTestId('Transmit Enabled')).toHaveTextContent('Enabled');
    expect(getByTestId('Gateway ID')).toHaveTextContent('testGw1');
    expect(getByTestId('Mme Connected')).toHaveTextContent('Disconnected');
  });

  it('unManaged eNodeB', async () => {
    jest.spyOn(hooks, 'useRefreshingContext').mockImplementation(() => ({
      enbInfo: {testEnodebSerial1: enbInfo['testEnodebSerial1']},
    }));
    const Wrapper = () => (
      <MemoryRouter
        initialEntries={['/nms/mynetwork/enodeb/testEnodebSerial1/overview']}
        initialIndex={0}>
        <MuiPickersUtilsProvider utils={MomentUtils}>
          <MuiThemeProvider theme={defaultTheme}>
            <MuiStylesThemeProvider theme={defaultTheme}>
              <EnodebContext.Provider
                value={{
                  state: {enbInfo: enbInfo},
                  setState: async _ => {},
                }}>
                <Routes>
                  <Route
                    path="/nms/:networkId/enodeb/:enodebSerial/overview/*"
                    element={<EnodebDetail enbInfo={enbInfo} />}
                  />
                </Routes>
              </EnodebContext.Provider>
            </MuiStylesThemeProvider>
          </MuiThemeProvider>
        </MuiPickersUtilsProvider>
      </MemoryRouter>
    );
    const {getByTestId} = render(<Wrapper />);
    await wait();

    // TODO - commenting this out till we have per enodeb metric support
    // expect(
    //   MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange,
    // ).toHaveBeenCalledTimes(2);
    expect(getByTestId('eNodeB Serial Number')).toHaveTextContent(
      'testEnodebSerial1',
    );
    expect(getByTestId('eNodeB Externally Managed')).toHaveTextContent('True');
    expect(getByTestId('Health')).toHaveTextContent('-');
    expect(getByTestId('Transmit Enabled')).toHaveTextContent('Disabled');
    expect(getByTestId('Gateway ID')).toHaveTextContent('testGw2');
    expect(getByTestId('Mme Connected')).toHaveTextContent('Disconnected');
  });
});
