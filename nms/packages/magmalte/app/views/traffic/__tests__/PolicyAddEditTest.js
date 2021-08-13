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
import MagmaAPIBindings from '@fbcnms/magma-api';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import NetworkContext from '../../../components/context/NetworkContext';
import React from 'react';
import TrafficDashboard from '../TrafficOverview';
import defaultTheme from '@fbcnms/ui/theme/default';

import {FEG_LTE, LTE} from '@fbcnms/types/network';
import {
  LteNetworkContextProvider,
  PolicyProvider,
} from '../../../components/lte/LteContext';

import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {cleanup, fireEvent, render, wait} from '@testing-library/react';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
const enqueueSnackbarMock = jest.fn();
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
  .mockReturnValue(enqueueSnackbarMock);
afterEach(cleanup);

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

const feg_network = {
  type: 'feg',
  description: '',
  dns: {
    enable_caching: false,
    local_ttl: 0,
  },
  federation: {},
  id: 'test_feg_network',
  name: 'Test Feg Network Description',
};

const feg_lte_network = {
  type: 'feg_lte',
  cellular: {
    epc: {
      default_rule_id: 'default_rule_1',
      gx_gy_relay_enabled: true,
      hss_relay_enabled: true,
      lte_auth_amf: 'gAA=',
      lte_auth_op: 'EREREREREREREREREREREQ==',
      mcc: '001',
      mnc: '01',
      mobility: {
        ip_allocation_mode: 'NAT',
        nat: {
          ip_blocks: ['192.168.0.0/16'],
        },
        reserved_addresses: ['192.168.0.1'],
        static: {
          ip_blocks_by_tac: {
            '1': ['192.168.0.0/16'],
            '2': ['172.10.0.0/16', '172.20.0.0/16'],
          },
        },
      },
      network_services: ['dpi', 'policy_enforcement'],
      sub_profiles: {
        additionalProp1: {
          max_dl_bit_rate: 20000000,
          max_ul_bit_rate: 100000000,
        },
        additionalProp2: {
          max_dl_bit_rate: 20000000,
          max_ul_bit_rate: 100000000,
        },
        additionalProp3: {
          max_dl_bit_rate: 20000000,
          max_ul_bit_rate: 100000000,
        },
      },
      tac: 1,
    },
    feg_network_id: 'test_feg_network',
    ran: {
      bandwidth_mhz: 20,
      tdd_config: {
        earfcndl: 44590,
        special_subframe_pattern: 7,
        subframe_assignment: 2,
      },
    },
  },
  description: 'Foobar network',
  dns: {
    enable_caching: false,
    local_ttl: 0,
  },

  federation: {
    feg_network_id: 'test_feg_network',
  },
  id: 'test_feg_lte_network',
  name: 'Test Feg LTE Network Description',
};

describe('<TrafficDashboard />', () => {
  afterEach(() => {
    MagmaAPIBindings.getFegLteByNetworkId.mockClear();
    MagmaAPIBindings.postNetworksByNetworkIdPoliciesRules.mockClear();
    MagmaAPIBindings.putFegLteByNetworkIdSubscriberConfig.mockClear();
    MagmaAPIBindings.putFegByNetworkIdSubscriberConfig.mockClear();
    MagmaAPIBindings.putLteByNetworkIdSubscriberConfig.mockClear();
  });
  beforeEach(() => {
    MagmaAPIBindings.getFegLteByNetworkId.mockResolvedValue(feg_lte_network);
    MagmaAPIBindings.getFegByNetworkId.mockResolvedValue(feg_network);
    MagmaAPIBindings.getLteByNetworkIdPolicyQosProfiles.mockResolvedValue({});
    MagmaAPIBindings.getNetworksByNetworkIdPoliciesRulesViewFull.mockResolvedValue(
      policies,
    );
    MagmaAPIBindings.getNetworksByNetworkIdPoliciesRulesByRuleId.mockResolvedValue(
      {
        app_name: undefined,
        app_service_type: undefined,
        assigned_subscribers: undefined,
        flow_list: [],
        id: 'test_policy_0',
        monitoring_key: '',
        priority: 1,
        qos_profile: undefined,
        rating_group: 0,
        redirect_information: {},
      },
    );
  });

  const PolicyWrapper = ({networkType}) => (
    <MemoryRouter
      initialEntries={['/nms/test/traffic/policy']}
      initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <NetworkContext.Provider
            value={{
              networkId: 'test',
              networkType: networkType,
            }}>
            <LteNetworkContextProvider
              networkId={'test'}
              networkType={networkType}>
              <PolicyProvider networkId={'test'} networkType={networkType}>
                <Route
                  path="/nms/:networkId/traffic/policy"
                  component={TrafficDashboard}
                />
              </PolicyProvider>
            </LteNetworkContextProvider>
          </NetworkContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  // verify feg_lte network wide policy add
  // verify lte network wide policy add
  // verify policy add
  // verify policy edit
  // verify policy edit to network wide
  // verify dpi services - app policy configuration
  it('verify feg_lte network wide policy add', async () => {
    const networkId = 'test';
    const fegNetworkId = 'test_feg_network';
    const {queryByTestId, getByTestId, getByText} = render(
      <PolicyWrapper networkType={FEG_LTE} />,
    );
    await wait();

    // verify if feg_lte and feg api calls are invoked
    expect(MagmaAPIBindings.getFegLteByNetworkId).toHaveBeenCalledWith({
      networkId,
    });
    expect(
      MagmaAPIBindings.getFegLteByNetworkIdSubscriberConfig,
    ).toHaveBeenCalledWith({networkId});
    expect(
      MagmaAPIBindings.getNetworksByNetworkIdPoliciesRulesViewFull,
    ).toHaveBeenCalledWith({networkId});
    expect(
      MagmaAPIBindings.getLteByNetworkIdPolicyQosProfiles,
    ).toHaveBeenCalledWith({networkId});
    expect(
      MagmaAPIBindings.getNetworksByNetworkIdPoliciesRulesViewFull,
    ).toHaveBeenCalledWith({networkId: fegNetworkId});

    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByText('Create New'));
    await wait();

    const newPolicyMenu = getByTestId('newPolicyMenuItem');
    fireEvent.click(newPolicyMenu);
    await wait();

    expect(queryByTestId('editDialog')).not.toBeNull();

    expect(queryByTestId('infoEdit')).not.toBeNull();
    expect(queryByTestId('flowEdit')).toBeNull();
    expect(queryByTestId('redirectEdit')).toBeNull();
    expect(queryByTestId('trackingEdit')).toBeNull();
    expect(queryByTestId('headerEnrichmentEdit')).toBeNull();
    expect(queryByTestId('appEdit')).toBeNull();

    const policyID = getByTestId('policyID').firstChild;
    const prio = getByTestId('policyPriority').firstChild;
    const networkWide = getByTestId('networkWide').firstChild;

    // test adding an existing policy
    if (policyID instanceof HTMLInputElement) {
      fireEvent.change(policyID, {target: {value: 'policy_0'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(getByTestId('configEditError')).toHaveTextContent(
      'Policy policy_0 already exists',
    );

    if (
      policyID instanceof HTMLInputElement &&
      prio instanceof HTMLInputElement &&
      networkWide instanceof HTMLElement
    ) {
      fireEvent.change(policyID, {target: {value: 'test_policy_0'}});
      fireEvent.change(prio, {target: {value: '1'}});
      if (networkWide.firstChild instanceof HTMLElement) {
        fireEvent.click(networkWide.firstChild);
      }
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(
      MagmaAPIBindings.postNetworksByNetworkIdPoliciesRules,
    ).toHaveBeenCalledWith({
      networkId: 'test',
      policyRule: {
        app_name: undefined,
        app_service_type: undefined,
        assigned_subscribers: undefined,
        flow_list: [],
        id: 'test_policy_0',
        monitoring_key: '',
        priority: 1,
        qos_profile: undefined,
        rating_group: 0,
        redirect_information: {},
      },
    });
    expect(
      MagmaAPIBindings.postNetworksByNetworkIdPoliciesRules,
    ).toHaveBeenCalledWith({
      networkId: 'test_feg_network',
      policyRule: {
        app_name: undefined,
        app_service_type: undefined,
        assigned_subscribers: undefined,
        flow_list: [],
        id: 'test_policy_0',
        monitoring_key: '',
        priority: 1,
        qos_profile: undefined,
        rating_group: 0,
        redirect_information: {},
      },
    });

    // verify if network's subscriber config is populated as well
    expect(
      MagmaAPIBindings.putFegLteByNetworkIdSubscriberConfig,
    ).toHaveBeenCalledWith({
      networkId: 'test',
      record: {
        network_wide_base_names: undefined,
        network_wide_rule_names: ['test_policy_0'],
      },
    });
    expect(
      MagmaAPIBindings.putFegByNetworkIdSubscriberConfig,
    ).toHaveBeenCalledWith({
      networkId: 'test_feg_network',
      record: {
        network_wide_base_names: undefined,
        network_wide_rule_names: ['test_policy_0'],
      },
    });
  });

  it('verify lte network wide policy add', async () => {
    const networkId = 'test';
    const {queryByTestId, getByTestId, getByText} = render(
      <PolicyWrapper networkType={LTE} />,
    );
    await wait();

    // verify if lte api calls are invoked
    expect(MagmaAPIBindings.getLteByNetworkId).toHaveBeenCalledWith({
      networkId,
    });
    expect(
      MagmaAPIBindings.getNetworksByNetworkIdPoliciesRulesViewFull,
    ).toHaveBeenCalledWith({networkId});
    expect(
      MagmaAPIBindings.getLteByNetworkIdPolicyQosProfiles,
    ).toHaveBeenCalledWith({networkId});

    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByText('Create New'));
    await wait();

    fireEvent.click(getByTestId('newPolicyMenuItem'));
    await wait();

    expect(queryByTestId('editDialog')).not.toBeNull();

    expect(queryByTestId('infoEdit')).not.toBeNull();
    expect(queryByTestId('flowEdit')).toBeNull();
    expect(queryByTestId('redirectEdit')).toBeNull();
    expect(queryByTestId('trackingEdit')).toBeNull();
    expect(queryByTestId('appEdit')).toBeNull();

    const policyID = getByTestId('policyID').firstChild;
    const prio = getByTestId('policyPriority').firstChild;
    const networkWide = getByTestId('networkWide').firstChild;

    // test adding an existing policy
    if (policyID instanceof HTMLInputElement) {
      fireEvent.change(policyID, {target: {value: 'policy_0'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(getByTestId('configEditError')).toHaveTextContent(
      'Policy policy_0 already exists',
    );

    if (
      policyID instanceof HTMLInputElement &&
      prio instanceof HTMLInputElement &&
      networkWide instanceof HTMLElement
    ) {
      fireEvent.change(policyID, {target: {value: 'test_policy_0'}});
      fireEvent.change(prio, {target: {value: '1'}});
      if (networkWide.firstChild instanceof HTMLElement) {
        fireEvent.click(networkWide.firstChild);
      }
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(
      MagmaAPIBindings.postNetworksByNetworkIdPoliciesRules,
    ).toHaveBeenCalledWith({
      networkId: 'test',
      policyRule: {
        app_name: undefined,
        app_service_type: undefined,
        assigned_subscribers: undefined,
        flow_list: [],
        id: 'test_policy_0',
        monitoring_key: '',
        priority: 1,
        qos_profile: undefined,
        rating_group: 0,
        redirect_information: {},
      },
    });

    // verify if network's subscriber config is populated as well
    expect(
      MagmaAPIBindings.putLteByNetworkIdSubscriberConfig,
    ).toHaveBeenCalledWith({
      networkId: 'test',
      record: {
        network_wide_base_names: undefined,
        network_wide_rule_names: ['test_policy_0'],
      },
    });
  });

  it('verify lte policy full add', async () => {
    const networkId = 'test';
    const {queryByTestId, getByTestId, getByText} = render(
      <PolicyWrapper networkType={LTE} />,
    );
    await wait();

    // verify if lte api calls are invoked
    expect(MagmaAPIBindings.getLteByNetworkId).toHaveBeenCalledWith({
      networkId,
    });
    expect(
      MagmaAPIBindings.getNetworksByNetworkIdPoliciesRulesViewFull,
    ).toHaveBeenCalledWith({networkId});
    expect(
      MagmaAPIBindings.getLteByNetworkIdPolicyQosProfiles,
    ).toHaveBeenCalledWith({networkId});

    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByText('Create New'));
    await wait();

    fireEvent.click(getByTestId('newPolicyMenuItem'));
    await wait();
    expect(queryByTestId('editDialog')).not.toBeNull();

    expect(queryByTestId('infoEdit')).not.toBeNull();
    expect(queryByTestId('flowEdit')).toBeNull();
    expect(queryByTestId('redirectEdit')).toBeNull();
    expect(queryByTestId('trackingEdit')).toBeNull();
    expect(queryByTestId('appEdit')).toBeNull();

    const policyID = getByTestId('policyID').firstChild;
    const prio = getByTestId('policyPriority').firstChild;

    if (
      policyID instanceof HTMLInputElement &&
      prio instanceof HTMLInputElement
    ) {
      fireEvent.change(policyID, {target: {value: 'test_policy_0'}});
      fireEvent.change(prio, {target: {value: '1'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();

    const newRule = {
      app_name: undefined,
      app_service_type: undefined,
      assigned_subscribers: undefined,
      flow_list: [],
      id: 'test_policy_0',
      monitoring_key: '',
      priority: 1,
      qos_profile: undefined,
      rating_group: 0,
      redirect_information: {},
    };

    expect(
      MagmaAPIBindings.postNetworksByNetworkIdPoliciesRules,
    ).toHaveBeenCalledWith({networkId, policyRule: newRule});

    // verify if we transition to flow tab
    expect(queryByTestId('infoEdit')).toBeNull();
    expect(queryByTestId('flowEdit')).not.toBeNull();
    expect(queryByTestId('redirectEdit')).toBeNull();
    expect(queryByTestId('trackingEdit')).toBeNull();
    expect(queryByTestId('headerEnrichmentEdit')).toBeNull();
    expect(queryByTestId('appEdit')).toBeNull();

    fireEvent.click(getByTestId('addFlowButton'));
    await wait();

    const ipSrc = getByTestId('ipSrc').firstChild;
    const ipDest = getByTestId('ipDest').firstChild;

    if (
      ipSrc instanceof HTMLInputElement &&
      ipDest instanceof HTMLInputElement
    ) {
      fireEvent.change(ipSrc, {target: {value: '1.1.1.1'}});
      fireEvent.change(ipDest, {target: {value: '1.1.1.2'}});
    } else {
      throw 'invalid type';
    }

    newRule['flow_list'] = [
      {
        action: 'DENY',
        match: {
          direction: 'UPLINK',
          ip_dst: {
            address: '1.1.1.2',
            version: 'IPv4',
          },
          ip_proto: 'IPPROTO_IP',
          ip_src: {
            address: '1.1.1.1',
            version: 'IPv4',
          },
        },
      },
    ];
    MagmaAPIBindings.getNetworksByNetworkIdPoliciesRulesByRuleId.mockResolvedValue(
      newRule,
    );
    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(
      MagmaAPIBindings.putNetworksByNetworkIdPoliciesRulesByRuleId,
    ).toHaveBeenCalledWith({
      networkId,
      ruleId: 'test_policy_0',
      policyRule: newRule,
    });

    // verify if we transition to tracking tab
    expect(queryByTestId('trackingEdit')).not.toBeNull();
    expect(queryByTestId('infoEdit')).toBeNull();
    expect(queryByTestId('flowEdit')).toBeNull();
    expect(queryByTestId('redirectEdit')).toBeNull();
    expect(queryByTestId('headerEnrichmentEdit')).toBeNull();

    const ratingGroup = getByTestId('ratingGroup').firstChild;

    if (ratingGroup instanceof HTMLInputElement) {
      fireEvent.change(ratingGroup, {target: {value: 1}});
    } else {
      throw 'invalid type';
    }

    newRule['rating_group'] = 1;
    MagmaAPIBindings.getNetworksByNetworkIdPoliciesRulesByRuleId.mockResolvedValue(
      newRule,
    );
    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(
      MagmaAPIBindings.putNetworksByNetworkIdPoliciesRulesByRuleId,
    ).toHaveBeenCalledWith({
      networkId,
      ruleId: 'test_policy_0',
      policyRule: newRule,
    });

    // verify if we transition to redirect tab
    expect(queryByTestId('trackingEdit')).toBeNull();
    expect(queryByTestId('infoEdit')).toBeNull();
    expect(queryByTestId('flowEdit')).toBeNull();
    expect(queryByTestId('redirectEdit')).not.toBeNull();
    expect(queryByTestId('headerEnrichmentEdit')).toBeNull();

    const serverAddress = getByTestId('serverAddress').firstChild;

    if (serverAddress instanceof HTMLInputElement) {
      fireEvent.change(serverAddress, {target: {value: '192.168.1.1'}});
    } else {
      throw 'invalid type';
    }
    (newRule['redirect_information'] = {
      server_address: '192.168.1.1',
    }),
      MagmaAPIBindings.getNetworksByNetworkIdPoliciesRulesByRuleId.mockResolvedValue(
        newRule,
      );
    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(
      MagmaAPIBindings.putNetworksByNetworkIdPoliciesRulesByRuleId,
    ).toHaveBeenCalledWith({
      networkId,
      ruleId: 'test_policy_0',
      policyRule: newRule,
    });

    // verify if we transition to header enrichment tab
    expect(queryByTestId('trackingEdit')).toBeNull();
    expect(queryByTestId('infoEdit')).toBeNull();
    expect(queryByTestId('flowEdit')).toBeNull();
    expect(queryByTestId('redirectEdit')).toBeNull();
    expect(queryByTestId('headerEnrichmentEdit')).not.toBeNull();

    const url = getByTestId('newUrl').firstChild;

    if (url instanceof HTMLInputElement) {
      fireEvent.change(url, {target: {value: 'http://example.com'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByTestId('addUrlButton'));
    await wait();

    fireEvent.click(getByText('Save And Continue'));

    await wait();

    expect(
      MagmaAPIBindings.putNetworksByNetworkIdPoliciesRulesByRuleId,
    ).toHaveBeenCalledWith({
      networkId,
      ruleId: 'test_policy_0',
      policyRule: {
        ...newRule,
        header_enrichment_targets: ['http://example.com'],
      },
    });
  });

  it('verify lte policy edit', async () => {
    const networkId = 'test';
    const {queryByTestId, getByTestId, getByText} = render(
      <PolicyWrapper networkType={LTE} />,
    );
    await wait();

    // verify if lte api calls are invoked
    expect(MagmaAPIBindings.getLteByNetworkId).toHaveBeenCalledWith({
      networkId,
    });
    expect(
      MagmaAPIBindings.getNetworksByNetworkIdPoliciesRulesViewFull,
    ).toHaveBeenCalledWith({networkId});
    expect(
      MagmaAPIBindings.getLteByNetworkIdPolicyQosProfiles,
    ).toHaveBeenCalledWith({networkId});

    expect(queryByTestId('editDialog')).toBeNull();
    fireEvent.click(getByText('policy_0'));
    await wait();
    expect(queryByTestId('editDialog')).not.toBeNull();

    expect(queryByTestId('infoEdit')).not.toBeNull();
    expect(queryByTestId('flowEdit')).toBeNull();
    expect(queryByTestId('redirectEdit')).toBeNull();
    expect(queryByTestId('trackingEdit')).toBeNull();
    expect(queryByTestId('headerEnrichmentEdit')).toBeNull();
    expect(queryByTestId('appEdit')).toBeNull();

    const prio = getByTestId('policyPriority').firstChild;

    if (prio instanceof HTMLInputElement) {
      fireEvent.change(prio, {target: {value: '2'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await wait();

    expect(
      MagmaAPIBindings.putNetworksByNetworkIdPoliciesRulesByRuleId,
    ).toHaveBeenCalledWith({
      networkId,
      ruleId: 'policy_0',
      policyRule: {
        flow_list: [],
        id: 'policy_0',
        monitoring_key: '',
        priority: 2,
      },
    });
  });
});
