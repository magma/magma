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
import type {promql_return_object} from '@fbcnms/magma-api';

import 'jest-dom/extend-expect';
import Enodeb from '../EquipmentEnodeb';
import EnodebContext from '../../../components/context/EnodebContext';
import MagmaAPIBindings from '@fbcnms/magma-api';
import MomentUtils from '@date-io/moment';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import axiosMock from 'axios';
import defaultTheme from '@fbcnms/ui/theme/default';

import * as hooks from '../../../components/context/RefreshContext';
import {MemoryRouter, Route} from 'react-router-dom';
import {MuiPickersUtilsProvider} from '@material-ui/pickers';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {cleanup, render, wait} from '@testing-library/react';
jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');
afterEach(cleanup);

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

const currTime = Date.now();

describe('<Enodeb />', () => {
  beforeAll(() => {
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mockResolvedValue(
      mockThroughput,
    );
    //   MagmaAPIBindings.getLteByNetworkIdEnodebs.mockResolvedValue(enbInfo);
  });

  afterEach(() => {
    axiosMock.get.mockClear();
  });

  const enbInfo0 = {
    enb: {
      attached_gateway_id: 'us_baltic_gw1',
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
      reporting_gateway_id: '',
      rf_tx_desired: true,
      rf_tx_on: false,
      time_reported: 0,
      ip_address: '192.168.1.254',
    },
  };

  const enbInfo1 = Object.assign({}, enbInfo0);
  enbInfo1.enb = {...enbInfo1.enb, name: 'testEnodeb1'};
  enbInfo1.enb_state = {
    ...enbInfo1.enb_state,
    fsm_state: 'initializing',
    time_reported: currTime,
    rf_tx_on: true,
  };
  const enbInfo = {
    testEnodebSerial0: enbInfo0,
    testEnodebSerial1: enbInfo1,
  };

  const enbCtx = {
    state: {enbInfo: enbInfo},
    setState: async _ => {},
  };

  jest
    .spyOn(hooks, 'useRefreshingContext')
    .mockImplementation(() => enbCtx.state);

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork/enodeb']} initialIndex={0}>
      <MuiPickersUtilsProvider utils={MomentUtils}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <EnodebContext.Provider value={enbCtx}>
              <Route path="/nms/:networkId/enodeb/" render={_ => <Enodeb />} />
            </EnodebContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MuiPickersUtilsProvider>
    </MemoryRouter>
  );

  it('renders', async () => {
    const {getAllByRole} = render(<Wrapper />);
    await wait();

    expect(
      MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange,
    ).toHaveBeenCalledTimes(1);

    const rowItems = await getAllByRole('row');

    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Name');
    expect(rowItems[0]).toHaveTextContent('Serial Number');
    expect(rowItems[0]).toHaveTextContent('Session State Name');
    expect(rowItems[0]).toHaveTextContent('Health');
    expect(rowItems[0]).toHaveTextContent('Reported Time');

    expect(rowItems[1]).toHaveTextContent('testEnodeb0');
    expect(rowItems[1]).toHaveTextContent('testEnodebSerial0');
    expect(rowItems[1]).toHaveTextContent(
      'Completed provisioning eNB. Awaiting new Inform.',
    );
    expect(rowItems[1]).toHaveTextContent('Bad');
    expect(rowItems[1]).toHaveTextContent(new Date(0).toLocaleDateString());

    expect(rowItems[2]).toHaveTextContent('testEnodeb1');
    expect(rowItems[2]).toHaveTextContent('testEnodebSerial1');
    expect(rowItems[2]).toHaveTextContent('initializing');
    expect(rowItems[2]).toHaveTextContent('Good');
    expect(rowItems[2]).toHaveTextContent(
      new Date(currTime).toLocaleDateString(),
    );
  });
});
