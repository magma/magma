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
import React from 'react';
import TrafficDashboard from '../TrafficOverview';
import defaultTheme from '../../../theme/default';
import {ApnContextProvider} from '../../../context/ApnContext';
import {PolicyProvider} from '../../../context/PolicyContext';

import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {PolicyQosProfile, PolicyRule, RatingGroup} from '../../../../generated';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {fireEvent, waitFor} from '@testing-library/react';
import {mockAPI} from '../../../util/TestUtils';
import {render} from '../../../util/TestingLibrary';

jest.mock('axios');
jest.mock('../../../../app/hooks/useSnackbar');
const apns = {
  apn_0: {
    apn_configuration: {
      ambr: {
        max_bandwidth_dl: 200000000,
        max_bandwidth_ul: 100000000,
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
        max_bandwidth_dl: 100000000,
        max_bandwidth_ul: 100000000,
      },
      qos_profile: {
        class_id: 6,
        preemption_capability: false,
        preemption_vulnerability: false,
        priority_level: 10,
      },
    },
    apn_name: 'apn_1',
  },
};

const policies: Record<string, PolicyRule> = {
  policy_0: {
    flow_list: [],
    id: 'policy_0',
    monitoring_key: '',
    priority: 1,
  },
  policy_1: {
    flow_list: [
      {
        action: 'PERMIT',
        match: {
          direction: 'UPLINK',
          ip_proto: 'IPPROTO_IP',
        },
      },
      {
        action: 'PERMIT',
        match: {
          direction: 'DOWNLINK',
          ip_proto: 'IPPROTO_IP',
        },
      },
    ],
    id: 'policy_1',
    monitoring_key: '',
    priority: 1,
  },
  policy_2: {
    flow_list: [],
    id: 'policy_2',
    monitoring_key: '',
    priority: 10,
  },
};

const qosProfiles: Record<string, PolicyQosProfile> = {
  profile_1: {
    id: 'profile_1',
    class_id: 1,
    max_req_bw_ul: 9,
    max_req_bw_dl: 9,
  },
  profile_2: {
    id: 'profile_2',
    class_id: 2,
    max_req_bw_ul: 10,
    max_req_bw_dl: 10,
  },
};
const ratingGroups: Record<string, RatingGroup> = {
  '0': {
    id: 0,
    limit_type: 'FINITE',
  },
  '1': {
    id: 1,
    limit_type: 'INFINITE_UNMETERED',
  },
};

describe('<TrafficDashboard />', () => {
  const networkId = 'test';

  beforeEach(() => {
    mockAPI(MagmaAPI.apns, 'lteNetworkIdApnsGet', apns);
    mockAPI(
      MagmaAPI.policies,
      'networksNetworkIdPoliciesRulesviewfullGet',
      policies,
    );
    mockAPI(MagmaAPI.policies, 'networksNetworkIdPoliciesBaseNamesGet', []);
    mockAPI(
      MagmaAPI.ratingGroups,
      'networksNetworkIdRatingGroupsGet',
      ratingGroups,
    );
    mockAPI(MagmaAPI.policies, 'lteNetworkIdPolicyQosProfilesGet', qosProfiles);
  });

  const Wrapper = () => (
    <MemoryRouter
      initialEntries={['/nms/test/traffic/policy']}
      initialIndex={0}>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={defaultTheme}>
          <PolicyProvider networkId={networkId}>
            <ApnContextProvider networkId={networkId}>
              <Routes>
                <Route
                  path="/nms/:networkId/traffic/*"
                  element={<TrafficDashboard />}
                />
              </Routes>
            </ApnContextProvider>
          </PolicyProvider>
        </ThemeProvider>
      </StyledEngineProvider>
    </MemoryRouter>
  );
  it('renders', async () => {
    const {
      findAllByRole,
      findByTestId,
      getByText,
      openActionsTableMenu,
    } = render(<Wrapper />);

    // Policy tables rows
    const rowItemsPolicy = await findAllByRole('row');
    // first row is the header
    expect(rowItemsPolicy[0]).toHaveTextContent('Policy ID');
    expect(rowItemsPolicy[0]).toHaveTextContent('Flows');
    expect(rowItemsPolicy[0]).toHaveTextContent('Priority');
    expect(rowItemsPolicy[0]).toHaveTextContent('Subscribers');
    expect(rowItemsPolicy[0]).toHaveTextContent('Monitoring Key');
    expect(rowItemsPolicy[0]).toHaveTextContent('Rating');
    expect(rowItemsPolicy[0]).toHaveTextContent('Tracking Type');
    expect(rowItemsPolicy[1]).toHaveTextContent('policy_0');
    expect(rowItemsPolicy[1]).toHaveTextContent('0');
    expect(rowItemsPolicy[1]).toHaveTextContent('1');
    expect(rowItemsPolicy[1]).toHaveTextContent('0');
    expect(rowItemsPolicy[1]).toHaveTextContent('Not Found');
    expect(rowItemsPolicy[1]).toHaveTextContent('NO_TRACKING');
    expect(rowItemsPolicy[2]).toHaveTextContent('policy_1');
    expect(rowItemsPolicy[2]).toHaveTextContent('2');
    expect(rowItemsPolicy[2]).toHaveTextContent('1');
    expect(rowItemsPolicy[2]).toHaveTextContent('0');
    expect(rowItemsPolicy[2]).toHaveTextContent('Not Found');
    expect(rowItemsPolicy[2]).toHaveTextContent('NO_TRACKING');
    expect(rowItemsPolicy[3]).toHaveTextContent('policy_2');
    expect(rowItemsPolicy[3]).toHaveTextContent('0');
    expect(rowItemsPolicy[3]).toHaveTextContent('10');
    expect(rowItemsPolicy[3]).toHaveTextContent('0');
    expect(rowItemsPolicy[3]).toHaveTextContent('Not Found');
    expect(rowItemsPolicy[3]).toHaveTextContent('NO_TRACKING');

    // click the actions button for policy 0
    await openActionsTableMenu(0);

    expect(await findByTestId('actions-menu')).toBeVisible();

    // Profiles tab
    fireEvent.click(getByText('Profiles'));
    const rowItemsProfile = await findAllByRole('row');
    // first row is the header
    expect(rowItemsProfile[0]).toHaveTextContent('Profile ID');
    expect(rowItemsProfile[0]).toHaveTextContent('Class ID');
    expect(rowItemsProfile[0]).toHaveTextContent('Uplink Bandwidth');
    expect(rowItemsProfile[0]).toHaveTextContent('Downlink Bandwidth');
    // profile_1
    expect(rowItemsProfile[1]).toHaveTextContent('profile_1');
    expect(rowItemsProfile[1]).toHaveTextContent('1');
    expect(rowItemsProfile[1]).toHaveTextContent('9');
    expect(rowItemsProfile[1]).toHaveTextContent('9');

    // profile_2
    expect(rowItemsProfile[2]).toHaveTextContent('profile_2');
    expect(rowItemsProfile[2]).toHaveTextContent('2');
    expect(rowItemsProfile[2]).toHaveTextContent('10');
    expect(rowItemsProfile[2]).toHaveTextContent('10');

    //Rating Groups Tab
    fireEvent.click(getByText('Rating Groups'));
    const rowItemsRatingGroups = await findAllByRole('row');
    // first row is the header
    expect(rowItemsRatingGroups[0]).toHaveTextContent('Rating Group ID');
    expect(rowItemsRatingGroups[0]).toHaveTextContent('Limit type');
    // Rating Group 0
    expect(rowItemsRatingGroups[1]).toHaveTextContent('0');
    expect(rowItemsRatingGroups[1]).toHaveTextContent('FINITE');
    // Rating Group 1
    expect(rowItemsRatingGroups[2]).toHaveTextContent('1');
    expect(rowItemsRatingGroups[2]).toHaveTextContent('INFINITE_UNMETERED');
    // click the actions button for rating group 0
    await openActionsTableMenu(0);
    expect(await findByTestId('actions-menu')).toBeVisible();
  });

  it('shows prompt when remove policy is clicked', async () => {
    const deleteMock = jest
      .spyOn(MagmaAPI.policies, 'networksNetworkIdPoliciesRulesRuleIdDelete')
      .mockImplementation();
    const {getByText, findByText, openActionsTableMenu} = render(<Wrapper />);

    // click remove action for policy 0
    await openActionsTableMenu(0);
    fireEvent.click(await findByText('Remove'));

    expect(
      await findByText('Are you sure you want to delete policy_0?'),
    ).toBeInTheDocument();
    // Confirm deletion
    fireEvent.click(getByText('Confirm'));
    await waitFor(() => {
      expect(deleteMock).toHaveBeenCalledWith({
        networkId: 'test',
        ruleId: 'policy_0',
      });
    });
  });
  it('shows prompt when remove profile is clicked', async () => {
    const deleteMock = jest
      .spyOn(MagmaAPI.policies, 'lteNetworkIdPolicyQosProfilesProfileIdDelete')
      .mockImplementation();
    const {getByText, findByText, openActionsTableMenu} = render(<Wrapper />);
    // Profiles tab
    fireEvent.click(await findByText('Profiles'));

    // click remove action for profile_1
    await openActionsTableMenu(0);
    fireEvent.click(await findByText('Remove'));

    expect(
      await findByText('Are you sure you want to delete profile_1?'),
    ).toBeInTheDocument();
    // Confirm deletion
    fireEvent.click(getByText('Confirm'));
    await waitFor(() => {
      expect(deleteMock).toHaveBeenCalledWith({
        networkId: 'test',
        profileId: 'profile_1',
      });
    });
  });
  it('shows prompt when remove rating group is clicked', async () => {
    const deleteMock = jest
      .spyOn(
        MagmaAPI.ratingGroups,
        'networksNetworkIdRatingGroupsRatingGroupIdDelete',
      )
      .mockImplementation();
    const {getByText, findByText, openActionsTableMenu} = render(<Wrapper />);

    // Rating Groups tab
    fireEvent.click(await findByText('Rating Groups'));

    // click remove action for rating group 0
    await openActionsTableMenu(0);
    fireEvent.click(await findByText('Remove'));
    expect(
      await findByText('Are you sure you want to delete Rating Group 0?'),
    ).toBeInTheDocument();
    // Confirm deletion
    fireEvent.click(getByText('Confirm'));
    await waitFor(() => {
      expect(deleteMock).toHaveBeenCalledWith({
        networkId: 'test',
        ratingGroupId: 0,
      });
    });
  });
});

describe('<TrafficDashboard APNs/>', () => {
  const {location} = window;
  const networkId = 'test';

  beforeEach((): void => {
    window.location = {
      pathname: '/nms/test/traffic/apn',
    } as Location;

    mockAPI(MagmaAPI.apns, 'lteNetworkIdApnsGet', apns);

    mockAPI(
      MagmaAPI.policies,
      'networksNetworkIdPoliciesRulesviewfullGet',
      policies,
    );
    mockAPI(MagmaAPI.policies, 'networksNetworkIdPoliciesBaseNamesGet', []);
    mockAPI(
      MagmaAPI.ratingGroups,
      'networksNetworkIdRatingGroupsGet',
      ratingGroups,
    );
    mockAPI(MagmaAPI.policies, 'lteNetworkIdPolicyQosProfilesGet', qosProfiles);
  });

  afterEach((): void => {
    window.location = location;
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/test/traffic/apn']} initialIndex={0}>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={defaultTheme}>
          <PolicyProvider networkId={networkId}>
            <ApnContextProvider networkId={networkId}>
              <Routes>
                <Route
                  path="/nms/:networkId/traffic/*"
                  element={<TrafficDashboard />}
                />
              </Routes>
            </ApnContextProvider>
          </PolicyProvider>
        </ThemeProvider>
      </StyledEngineProvider>
    </MemoryRouter>
  );
  it('renders', async () => {
    const {
      findAllByText,
      getByTestId,
      getAllByRole,
      openActionsTableMenu,
    } = render(<Wrapper />);

    const apnTitles = await findAllByText('APNs');
    expect(apnTitles.length).toBe(2);

    // Apn tables rows
    const rowItemsApns = getAllByRole('row');
    // first row is the header
    expect(rowItemsApns[0]).toHaveTextContent('Apn ID');
    expect(rowItemsApns[0]).toHaveTextContent('Class ID');
    expect(rowItemsApns[0]).toHaveTextContent('Priority Level');
    expect(rowItemsApns[0]).toHaveTextContent('Max Reqd UL Bw');
    expect(rowItemsApns[0]).toHaveTextContent('Max Reqd DL Bw');
    expect(rowItemsApns[0]).toHaveTextContent('Pre-emption Capability');
    expect(rowItemsApns[0]).toHaveTextContent('Pre-emption Vulnerability');

    // check first data row
    expect(rowItemsApns[1]).toHaveTextContent('apn_0');
    expect(rowItemsApns[1]).toHaveTextContent('9');
    expect(rowItemsApns[1]).toHaveTextContent('15');
    expect(rowItemsApns[1]).toHaveTextContent('100000000');
    expect(rowItemsApns[1]).toHaveTextContent('200000000');
    expect(rowItemsApns[1]).toHaveTextContent('true');
    expect(rowItemsApns[1]).toHaveTextContent('false');

    // check second data row
    expect(rowItemsApns[2]).toHaveTextContent('apn_1');
    expect(rowItemsApns[2]).toHaveTextContent('6');
    expect(rowItemsApns[2]).toHaveTextContent('10');
    expect(rowItemsApns[2]).toHaveTextContent('100000000');
    expect(rowItemsApns[2]).toHaveTextContent('100000000');
    expect(rowItemsApns[2]).toHaveTextContent('false');
    expect(rowItemsApns[2]).toHaveTextContent('false');

    // click the actions button for apn 0
    await openActionsTableMenu(0);
    await waitFor(() => {
      expect(getByTestId('actions-menu')).toBeVisible();
    });
  });

  it('shows prompt when remove apn is clicked', async () => {
    const deleteMock = jest
      .spyOn(MagmaAPI.apns, 'lteNetworkIdApnsApnNameDelete')
      .mockImplementation();
    const {getByText, findByText, openActionsTableMenu} = render(<Wrapper />);

    await openActionsTableMenu(0);
    fireEvent.click(await findByText('Remove'));
    expect(
      await findByText('Are you sure you want to delete apn_0?'),
    ).toBeInTheDocument();
    // Confirm deletion
    fireEvent.click(getByText('Confirm'));
    await waitFor(() => {
      expect(deleteMock).toHaveBeenCalledWith({
        networkId: 'test',
        apnName: 'apn_0',
      });
    });
  });
});
