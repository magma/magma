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
// $FlowFixMe migrated to typescript
import ApnContext from '../../../components/context/ApnContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import MagmaAPI from '../../../../api/MagmaAPI';
import MagmaAPIBindings from '../../../../generated/MagmaAPIBindings';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
// $FlowFixMe migrated to typescript
import PolicyContext from '../../../components/context/PolicyContext';
import React from 'react';
import TrafficDashboard from '../TrafficOverview';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../../theme/default';

import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
// $FlowFixMe migrated to typescript
import {SetApnState} from '../../../state/lte/ApnState';
import {
  SetPolicyState,
  SetQosProfileState,
  SetRatingGroupState,
  // $FlowFixMe migrated to typescript
} from '../../../state/PolicyState';
import {fireEvent, render, wait} from '@testing-library/react';

jest.mock('axios');
jest.mock('../../../../generated/MagmaAPIBindings.js');
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

const policies = {
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

const qosProfiles = {
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
const ratingGroups = {
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
  const policyCtx = {
    state: policies,
    baseNames: {},
    qosProfiles: qosProfiles,
    ratingGroups: ratingGroups,
    setBaseNames: (_key, _value?) => Promise.resolve(),
    setRatingGroups: (key, value?) => {
      return SetRatingGroupState({
        ratingGroups,
        setRatingGroups: () => {},
        networkId,
        key,
        value,
      });
    },
    setQosProfiles: (key, value?) => {
      return SetQosProfileState({
        qosProfiles,
        setQosProfiles: () => {},
        networkId,
        key,
        value,
      });
    },
    setState: (key, value?) => {
      return SetPolicyState({
        policies,
        setPolicies: () => {},
        networkId,
        key,
        value,
      });
    },
  };

  const apnCtx = {
    state: apns,
    setState: (key, value?) => {
      return SetApnState({
        apns,
        setApns: () => {},
        networkId,
        key,
        value,
      });
    },
  };

  beforeEach(() => {
    MagmaAPIBindings.getNetworks.mockResolvedValue([]);
  });

  const Wrapper = () => (
    <MemoryRouter
      initialEntries={['/nms/test/traffic/policy']}
      initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <PolicyContext.Provider value={policyCtx}>
            <ApnContext.Provider value={apnCtx}>
              <Routes>
                <Route
                  path="/nms/:networkId/traffic/*"
                  element={<TrafficDashboard />}
                />
              </Routes>
            </ApnContext.Provider>
          </PolicyContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
  it('renders', async () => {
    jest.setTimeout(30000);
    const {getByTestId, getAllByRole, getByText, getAllByTitle} = render(
      <Wrapper />,
    );
    await wait();
    // Policy tables rows
    const rowItemsPolicy = await getAllByRole('row');
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
    const policyActionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(policyActionList[0]);
    await wait();
    expect(getByTestId('actions-menu')).toBeVisible();

    // Profiles tab
    fireEvent.click(getByText('Profiles'));
    await wait();
    const rowItemsProfile = await getAllByRole('row');
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
    await wait();
    const rowItemsRatingGroups = await getAllByRole('row');
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
    const ratingGroupActionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(ratingGroupActionList[0]);
    await wait();
    expect(getByTestId('actions-menu')).toBeVisible();
  });

  it('shows prompt when remove policy is clicked', async () => {
    const deleteMock = jest
      .spyOn(MagmaAPI.policies, 'networksNetworkIdPoliciesRulesRuleIdDelete')
      .mockResolvedValueOnce({});
    const {getByText, getByTestId, getAllByTitle} = render(<Wrapper />);
    await wait();
    // click remove action for policy 0
    const policyActionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(policyActionList[0]);
    await wait();
    fireEvent.click(getByText('Remove'));
    await wait();
    expect(
      getByText('Are you sure you want to delete policy_0?'),
    ).toBeInTheDocument();
    // Confirm deletion
    fireEvent.click(getByText('Confirm'));
    await wait();
    expect(deleteMock).toHaveBeenCalledWith({
      networkId: 'test',
      ruleId: 'policy_0',
    });
  });
  it('shows prompt when remove profile is clicked', async () => {
    const deleteMock = jest
      .spyOn(MagmaAPI.policies, 'lteNetworkIdPolicyQosProfilesProfileIdDelete')
      .mockResolvedValueOnce({});
    const {getByText, getByTestId, getAllByTitle} = render(<Wrapper />);
    await wait();
    // Profiles tab
    fireEvent.click(getByText('Profiles'));
    await wait();
    // click remove action for profile_1
    const profileActionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(profileActionList[0]);
    await wait();
    fireEvent.click(getByText('Remove'));
    await wait();
    expect(
      getByText('Are you sure you want to delete profile_1?'),
    ).toBeInTheDocument();
    // Confirm deletion
    fireEvent.click(getByText('Confirm'));
    await wait();
    expect(deleteMock).toHaveBeenCalledWith({
      networkId: 'test',
      profileId: 'profile_1',
    });
  });
  it('shows prompt when remove rating group is clicked', async () => {
    const deleteMock = jest
      .spyOn(
        MagmaAPI.ratingGroups,
        'networksNetworkIdRatingGroupsRatingGroupIdDelete',
      )
      .mockResolvedValueOnce({});
    const {getByText, getByTestId, getAllByTitle} = render(<Wrapper />);
    await wait();
    // Rating Groups tab
    fireEvent.click(getByText('Rating Groups'));
    await wait();
    // click remove action for rating group 0
    const ratingGroupActionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(ratingGroupActionList[0]);
    await wait();
    fireEvent.click(getByText('Remove'));
    await wait();
    expect(
      getByText('Are you sure you want to delete Rating Group 0?'),
    ).toBeInTheDocument();
    // Confirm deletion
    fireEvent.click(getByText('Confirm'));
    await wait();
    expect(deleteMock).toHaveBeenCalledWith({
      networkId: 'test',
      ratingGroupId: 0,
    });
  });
});

describe('<TrafficDashboard APNs/>', () => {
  const {location} = window;
  beforeEach(() => {
    MagmaAPIBindings.getNetworks.mockResolvedValue([]);
  });

  beforeAll((): void => {
    delete window.location;
    window.location = {
      pathname: '/nms/test/traffic/apn',
    };
  });

  afterAll((): void => {
    window.location = location;
  });

  const networkId = 'test';
  const policyCtx = {
    state: policies,
    baseNames: {},
    qosProfiles: {},
    ratingGroups: {},
    setBaseNames: async () => {},
    setRatingGroups: async () => {},
    setQosProfiles: async () => {},
    setState: (key, value?) => {
      return SetPolicyState({
        policies,
        setPolicies: () => {},
        networkId,
        key,
        value,
      });
    },
  };
  const apnCtx = {
    state: apns,
    setState: (key, value?) => {
      return SetApnState({
        apns,
        setApns: () => {},
        networkId,
        key,
        value,
      });
    },
  };
  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/test/traffic/apn']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <PolicyContext.Provider value={policyCtx}>
            <ApnContext.Provider value={apnCtx}>
              <Routes>
                <Route
                  path="/nms/:networkId/traffic/*"
                  element={<TrafficDashboard />}
                />
              </Routes>
            </ApnContext.Provider>
          </PolicyContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
  it('renders', async () => {
    const {getAllByText, getByTestId, getAllByRole, getAllByTitle} = render(
      <Wrapper />,
    );
    await wait();

    const apnTitles = getAllByText('APNs');
    expect(apnTitles.length).toBe(2);

    // Apn tables rows
    const rowItemsApns = await getAllByRole('row');
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
    const apnActionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(apnActionList[0]);
    await wait();
    expect(getByTestId('actions-menu')).toBeVisible();
  });

  it('shows prompt when remove apn is clicked', async () => {
    MagmaAPIBindings.deleteLteByNetworkIdApnsByApnName.mockResolvedValueOnce(
      {},
    );
    const deleteMock = jest
      .spyOn(MagmaAPI.apns, 'lteNetworkIdApnsApnNameDelete')
      .mockResolvedValueOnce({});
    const {getByText, getByTestId, getAllByTitle} = render(<Wrapper />);
    await wait();

    const apnActionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(apnActionList[0]);
    await wait();
    fireEvent.click(getByText('Remove'));
    await wait();
    expect(
      getByText('Are you sure you want to delete apn_0?'),
    ).toBeInTheDocument();
    // Confirm deletion
    fireEvent.click(getByText('Confirm'));
    await wait();
    expect(deleteMock).toHaveBeenCalledWith({
      networkId: 'test',
      apnName: 'apn_0',
    });
  });
});
