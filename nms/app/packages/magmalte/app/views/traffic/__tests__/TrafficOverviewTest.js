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
import 'jest-dom/extend-expect';
import ApnContext from '../../../components/context/ApnContext';
import MagmaAPIBindings from '@fbcnms/magma-api';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import PolicyContext from '../../../components/context/PolicyContext';
import React from 'react';
import TrafficDashboard from '../TrafficOverview';
import axiosMock from 'axios';
import defaultTheme from '@fbcnms/ui/theme/default';

import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SetApnState} from '../../../state/lte/ApnState';
import {SetPolicyState} from '../../../state/PolicyState';
import {cleanup, fireEvent, render, wait} from '@testing-library/react';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');
afterEach(cleanup);

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

describe('<TrafficDashboard />', () => {
  const networkId = 'test';
  const policyCtx = {
    state: policies,
    qosProfiles: {},
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
    <MemoryRouter
      initialEntries={['/nms/test/traffic/policy']}
      initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <PolicyContext.Provider value={policyCtx}>
            <ApnContext.Provider value={apnCtx}>
              <Route
                path="/nms/:networkId/traffic/policy"
                component={TrafficDashboard}
              />
            </ApnContext.Provider>
          </PolicyContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
  it('renders', async () => {
    const {getByTestId, getAllByRole, getAllByTitle, getByText} = render(
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
    // Apns tab
    fireEvent.click(getByText('APNs'));
    await wait();
    expect(getByTestId('title_APNs')).toHaveTextContent('APNs');
    // Apn tables rows
    const rowItemsApns = await getAllByRole('row');
    // first row is the header
    expect(rowItemsApns[0]).toHaveTextContent('Apn ID');
    expect(rowItemsApns[0]).toHaveTextContent('Description');
    expect(rowItemsApns[0]).toHaveTextContent('Qos Profile');
    expect(rowItemsApns[0]).toHaveTextContent('Added');
    expect(rowItemsApns[1]).toHaveTextContent('apn_0');
    expect(rowItemsApns[1]).toHaveTextContent('Test APN description');
    expect(rowItemsApns[1]).toHaveTextContent('1');
    expect(rowItemsApns[2]).toHaveTextContent('apn_1');
    expect(rowItemsApns[2]).toHaveTextContent('Test APN description');
    expect(rowItemsApns[2]).toHaveTextContent('1');
    // click the actions button for apn 0
    const apnActionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(apnActionList[0]);
    await wait();
    expect(getByTestId('actions-menu')).toBeVisible();
  });
  it('shows prompt when remove policy is clicked', async () => {
    MagmaAPIBindings.deleteNetworksByNetworkIdPoliciesRulesByRuleId.mockResolvedValueOnce(
      {},
    );
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
    expect(
      MagmaAPIBindings.deleteNetworksByNetworkIdPoliciesRulesByRuleId,
    ).toHaveBeenCalledWith({
      networkId: 'test',
      ruleId: 'policy_0',
    });
    axiosMock.delete.mockClear();
  });
  it('shows prompt when remove apn is clicked', async () => {
    MagmaAPIBindings.deleteLteByNetworkIdApnsByApnName.mockResolvedValueOnce(
      {},
    );
    const {getByText, getByTestId, getAllByTitle} = render(<Wrapper />);
    await wait();
    fireEvent.click(getByText('APNs'));
    await wait();
    // click remove action for policy 0
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
    expect(
      MagmaAPIBindings.deleteLteByNetworkIdApnsByApnName,
    ).toHaveBeenCalledWith({
      networkId: 'test',
      apnName: 'apn_0',
    });
    axiosMock.delete.mockClear();
  });
});
