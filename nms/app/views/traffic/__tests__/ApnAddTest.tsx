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
import MagmaAPI from '../../../api/MagmaAPI';
import NetworkContext from '../../../context/NetworkContext';
import React from 'react';
import TrafficDashboard from '../TrafficOverview';
import defaultTheme from '../../../theme/default';

import {ApnContextProvider} from '../../../context/ApnContext';
import {LteNetworkContextProvider} from '../../../context/LteNetworkContext';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {fireEvent, render, waitFor} from '@testing-library/react';
import {mockAPI} from '../../../util/TestUtils';

jest.mock('axios');
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
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={defaultTheme}>
          <NetworkContext.Provider
            value={{
              networkId: 'test',
            }}>
            <LteNetworkContextProvider networkId={'test'}>
              <ApnContextProvider networkId={'test'}>
                <Routes>
                  <Route
                    path="/nms/:networkId/traffic/*"
                    element={<TrafficDashboard />}
                  />
                </Routes>
              </ApnContextProvider>
            </LteNetworkContextProvider>
          </NetworkContext.Provider>
        </ThemeProvider>
      </StyledEngineProvider>
    </MemoryRouter>
  );

  // verify apn add
  // verify apn edit

  it('verify apn add', async () => {
    const networkId = 'test';
    const {
      queryByTestId,
      getByTestId,
      findByTestId,
      getByText,
      findByText,
    } = render(<ApnWrapper />);

    await waitFor(() => {
      expect(MagmaAPI.lteNetworks.lteNetworkIdGet).toHaveBeenCalledWith({
        networkId,
      });
      expect(MagmaAPI.apns.lteNetworkIdApnsGet).toHaveBeenCalledWith({
        networkId,
      });
    });

    expect(queryByTestId('editDialog')).toBeNull();
    const newAPNButton = await findByTestId('newApnButton');

    fireEvent.click(newAPNButton);

    expect(await findByTestId('editDialog')).not.toBeNull();
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

    expect(await findByTestId('configEditError')).toHaveTextContent(
      'APN apn_0 already exists',
    );

    if (
      apnID instanceof HTMLInputElement &&
      classID instanceof HTMLInputElement &&
      apnPriority instanceof HTMLInputElement &&
      apnBandwidthUL instanceof HTMLInputElement &&
      apnBandwidthDL instanceof HTMLInputElement &&
      preemptionCapability instanceof HTMLInputElement &&
      preemptionVulnerability instanceof HTMLInputElement &&
      pdnType instanceof HTMLElement
    ) {
      fireEvent.change(apnID, {target: {value: 'apn_2'}});
      fireEvent.change(classID, {target: {value: 9}});
      fireEvent.change(apnPriority, {target: {value: 15}});
      fireEvent.change(apnBandwidthUL, {target: {value: 1000000}});
      fireEvent.change(apnBandwidthDL, {target: {value: 1000000}});
      fireEvent.click(preemptionCapability);
      fireEvent.click(preemptionVulnerability);
      fireEvent.mouseDown(pdnType);
      fireEvent.click(await findByText('IPv6'));
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));

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
    await waitFor(() => {
      expect(MagmaAPI.apns.lteNetworkIdApnsPost).toHaveBeenCalledWith({
        networkId,
        apn: newApn,
      });
    });
  });

  it('verify apn edit', async () => {
    const networkId = 'test';
    const {queryByTestId, getByTestId, getByText, findByText} = render(
      <ApnWrapper />,
    );

    await waitFor(() => {
      // verify if lte api calls are invoked
      expect(MagmaAPI.lteNetworks.lteNetworkIdGet).toHaveBeenCalledWith({
        networkId,
      });
      expect(MagmaAPI.apns.lteNetworkIdApnsGet).toHaveBeenCalledWith({
        networkId,
      });
    });

    expect(queryByTestId('editDialog')).toBeNull();

    // click on apns tab
    fireEvent.click(await findByText('apn_0'));

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

    await waitFor(() => {
      expect(MagmaAPI.apns.lteNetworkIdApnsApnNamePut).toHaveBeenCalledWith({
        networkId: 'test',
        apn: newApn,
        apnName: newApn.apn_name,
      });
    });
  });
});
