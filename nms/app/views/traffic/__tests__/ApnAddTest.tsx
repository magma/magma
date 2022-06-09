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
import MagmaAPI from '../../../../api/MagmaAPI';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import NetworkContext from '../../../components/context/NetworkContext';
import React from 'react';
import TrafficDashboard from '../TrafficOverview';
import defaultTheme from '../../../theme/default';
import {LTE} from '../../../../shared/types/network';

import {
  ApnProvider,
  LteNetworkContextProvider,
} from '../../../components/lte/LteContext';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {fireEvent, render, wait} from '@testing-library/react';
import {mockAPI} from '../../../util/TestUtils';

jest.mock('axios');
jest.mock('../../../../generated/MagmaAPIBindings.js');
jest.mock('../../../hooks/useSnackbar');

const apns = {
  apn_0: {
    apn_configuration: {
      ambr: {
        max_bandwidth_dl: 1000000,
        max_bandwidth_ul: 1000000,
      },
      qos_profile: {
        class_id: 9,
        preemption_capability: true,
        preemption_vulnerability: false,
        priority_level: 15,
      },
    },
    apn_name: 'apn_0',
  },
  apn_1: {
    apn_configuration: {
      ambr: {
        max_bandwidth_dl: 1000000,
        max_bandwidth_ul: 1000000,
      },
      qos_profile: {
        class_id: 9,
        preemption_capability: true,
        preemption_vulnerability: false,
        priority_level: 15,
      },
    },
    apn_name: 'apn_1',
  },
};

describe('<TrafficDashboard />', () => {
  beforeEach(() => {
    mockAPI(MagmaAPI.lteNetworks, 'lteNetworkIdGet');
    mockAPI(MagmaAPI.apns, 'lteNetworkIdApnsGet', apns);
    mockAPI(MagmaAPI.networks, 'networksGet', []);
    mockAPI(MagmaAPI.networks, 'networksNetworkIdTypeGet');
    mockAPI(MagmaAPI.apns, 'lteNetworkIdApnsPost');
    mockAPI(MagmaAPI.apns, 'lteNetworkIdApnsApnNamePut');
  });

  const {location} = window;
  beforeAll((): void => {
    // @ts-ignore
    delete window.location;
    window.location = {
      pathname: '/nms/test/traffic/apn',
    } as Location;
  });

  afterAll((): void => {
    window.location = location;
  });

  const ApnWrapper = () => (
    <MemoryRouter initialEntries={['/nms/test/traffic/apn']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <NetworkContext.Provider
            value={{
              networkId: 'test',
            }}>
            <LteNetworkContextProvider networkId={'test'} networkType={LTE}>
              <ApnProvider networkId={'test'} networkType={LTE}>
                <Routes>
                  <Route
                    path="/nms/:networkId/traffic/*"
                    element={<TrafficDashboard />}
                  />
                </Routes>
              </ApnProvider>
            </LteNetworkContextProvider>
          </NetworkContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  // verify apn add
  // verify apn edit

  it('verify apn add', async () => {
    jest.setTimeout(30000);
    const networkId = 'test';
    const {queryByTestId, getByTestId, getByText} = render(<ApnWrapper />);
    await wait();

    expect(MagmaAPI.lteNetworks.lteNetworkIdGet).toHaveBeenCalledWith({
      networkId,
    });
    expect(MagmaAPI.apns.lteNetworkIdApnsGet).toHaveBeenCalledWith({
      networkId,
    });

    expect(queryByTestId('editDialog')).toBeNull();
    await wait();

    const newAPNButton = queryByTestId('newApnButton');
    expect(newAPNButton).not.toBeNull();

    if (newAPNButton) {
      fireEvent.click(newAPNButton);
      await wait();
    }
    expect(queryByTestId('editDialog')).not.toBeNull();

    expect(queryByTestId('apnEditDialog')).not.toBeNull();

    const apnID = getByTestId('apnID').firstChild;
    const classID = getByTestId('classID').firstChild;
    const apnPriority = getByTestId('apnPriority').firstChild;
    const apnBandwidthUL = getByTestId('apnBandwidthUL').firstChild;
    const apnBandwidthDL = getByTestId('apnBandwidthDL').firstChild;
    const preemptionCapability = getByTestId('preemptionCapability').firstChild;
    const preemptionVulnerability = getByTestId('preemptionVulnerability')
      .firstChild;
    const pdnType = getByTestId('pdnType').firstChild;

    // test adding an existing apn
    if (apnID instanceof HTMLInputElement) {
      fireEvent.change(apnID, {target: {value: 'apn_0'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await wait();

    expect(getByTestId('configEditError')).toHaveTextContent(
      'APN apn_0 already exists',
    );

    if (
      apnID instanceof HTMLInputElement &&
      classID instanceof HTMLInputElement &&
      apnPriority instanceof HTMLInputElement &&
      apnBandwidthUL instanceof HTMLInputElement &&
      apnBandwidthDL instanceof HTMLInputElement &&
      pdnType instanceof HTMLElement
    ) {
      fireEvent.change(apnID, {target: {value: 'apn_2'}});
      fireEvent.change(classID, {target: {value: 9}});
      fireEvent.change(apnPriority, {target: {value: 15}});
      fireEvent.change(apnBandwidthUL, {target: {value: 1000000}});
      fireEvent.change(apnBandwidthDL, {target: {value: 1000000}});
      if (preemptionCapability?.firstChild instanceof HTMLElement) {
        fireEvent.click(preemptionCapability.firstChild);
      }
      if (preemptionVulnerability?.firstChild instanceof HTMLElement) {
        fireEvent.click(preemptionVulnerability.firstChild);
      }
      fireEvent.mouseDown(pdnType);
      await wait();
      fireEvent.click(getByText('IPv6'));
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await wait();

    const newApn = {
      apn_configuration: {
        ambr: {max_bandwidth_dl: 1000000, max_bandwidth_ul: 1000000},
        pdn_type: 1,
        qos_profile: {
          class_id: 9,
          preemption_capability: true,
          preemption_vulnerability: true,
          priority_level: 15,
        },
      },
      apn_name: 'apn_2',
    };

    expect(MagmaAPI.apns.lteNetworkIdApnsPost).toHaveBeenCalledWith({
      networkId,
      apn: newApn,
    });
  });

  it('verify apn edit', async () => {
    const networkId = 'test';
    const {queryByTestId, getByTestId, getByText} = render(<ApnWrapper />);
    await wait();

    // verify if lte api calls are invoked
    expect(MagmaAPI.lteNetworks.lteNetworkIdGet).toHaveBeenCalledWith({
      networkId,
    });
    expect(MagmaAPI.apns.lteNetworkIdApnsGet).toHaveBeenCalledWith({
      networkId,
    });

    expect(queryByTestId('editDialog')).toBeNull();

    // click on apns tab
    fireEvent.click(getByText('apn_0'));
    await wait();
    expect(queryByTestId('editDialog')).not.toBeNull();
    expect(queryByTestId('apnEditDialog')).not.toBeNull();

    const classID = getByTestId('classID').firstChild;
    const apnPriority = getByTestId('apnPriority').firstChild;

    if (
      classID instanceof HTMLInputElement &&
      apnPriority instanceof HTMLInputElement
    ) {
      fireEvent.change(classID, {target: {value: 8}});
      fireEvent.change(apnPriority, {target: {value: 10}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await wait();

    const newApn = {
      apn_configuration: {
        ambr: {max_bandwidth_dl: 1000000, max_bandwidth_ul: 1000000},
        pdn_type: 0,
        qos_profile: {
          class_id: 8,
          preemption_capability: true,
          preemption_vulnerability: false,
          priority_level: 10,
        },
      },
      apn_name: 'apn_0',
    };

    expect(MagmaAPI.apns.lteNetworkIdApnsApnNamePut).toHaveBeenCalledWith({
      networkId: 'test',
      apn: newApn,
      apnName: newApn.apn_name,
    });
  });
});
