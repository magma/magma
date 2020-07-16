/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import 'jest-dom/extend-expect';

import MagmaAPIBindings from '@fbcnms/magma-api';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import NetworkDashboard from '../NetworkDashboard';
import React from 'react';
import axiosMock from 'axios';
import defaultTheme from '../../../theme/default.js';

import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {cleanup, fireEvent, render, wait} from '@testing-library/react';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');
afterEach(cleanup);
const enqueueSnackbarMock = jest.fn();
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
  .mockReturnValue(enqueueSnackbarMock);

describe('<NetworkDashboard />', () => {
  const testNetwork = {
    description: 'Test Network Description',
    dns: {
      enable_caching: true,
      local_ttl: 60,
    },
    features: {
      features: {
        networkType: 'lte',
      },
    },
    id: 'test_network',
    name: 'Test Network',
    type: 'lte',
  };

  const epc = {
    default_rule_id: 'default_rule_1',
    lte_auth_amf: 'gAA=',
    lte_auth_op: 'EREREREREREREREREREREQ==',
    mcc: '001',
    mnc: '01',
    network_services: ['dpi', 'policy_enforcement'],
    relay_enabled: false,
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
  };

  const ran = {
    bandwidth_mhz: 20,
    tdd_config: {
      earfcndl: 44390,
      special_subframe_pattern: 7,
      subframe_assignment: 2,
    },
  };

  const gateways = {
    test_gateway1: {
      cellular: {
        epc: {
          ip_block: '192.168.128.0/24',
          nat_enabled: true,
        },
        ran: {
          pci: 260,
          transmit_enabled: false,
        },
      },
      connected_enodeb_serials: null,
      description: '',
      device: null,
      id: '',
      magmad: null,
      name: '',
      tier: '',
    },
  };

  const enodebs = {
    '120200020718CJP0013': {
      attached_gateway_id: 'mpk_dogfooding_tiplab_1',
      config: {
        bandwidth_mhz: 10,
        cell_id: 6553601,
        device_class: 'Baicells Nova-233 G2 OD FDD',
        earfcndl: 9410,
        pci: 60,
        tac: 6,
        transmit_enabled: true,
      },
      name: '',
      serial: '120200020718CJP0013',
    },
  };

  const rules = ['test1', 'test2'];

  const subscribers = {
    IMSI00000000001002: {
      active_apns: ['oai.ipv4'],
      id: 'IMSI722070171001002',
      lte: {
        auth_algo: 'MILENAGE',
        auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
        auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
        state: 'ACTIVE',
        sub_profile: 'default',
      },
    },
  };

  const apns = {
    internet: {
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
      apn_name: 'internet',
    },
    'oai.ipv4': {
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
      apn_name: 'oai.ipv4',
    },
  };

  beforeEach(() => {
    MagmaAPIBindings.getNetworksByNetworkId.mockResolvedValue(testNetwork);
    MagmaAPIBindings.getLteByNetworkIdCellularEpc.mockResolvedValue(epc);
    MagmaAPIBindings.getLteByNetworkIdCellularRan.mockResolvedValue(ran);
    MagmaAPIBindings.getLteByNetworkIdGateways.mockResolvedValue(gateways);
    MagmaAPIBindings.getLteByNetworkIdEnodebs.mockResolvedValue(enodebs);
    // eslint-disable-next-line max-len
    MagmaAPIBindings.getNetworksByNetworkIdPoliciesRules.mockResolvedValue(
      rules,
    );
    // eslint-disable-next-line max-len
    MagmaAPIBindings.getLteByNetworkIdSubscribers.mockResolvedValue(
      subscribers,
    );
    MagmaAPIBindings.getLteByNetworkIdApns.mockResolvedValue(apns);

    axiosMock.post.mockImplementation(() =>
      Promise.resolve({data: {success: true}}),
    );
    MagmaAPIBindings.putNetworksByNetworkId.mockImplementation(() =>
      Promise.resolve({data: {success: true}}),
    );
    MagmaAPIBindings.putLteByNetworkIdCellularEpc.mockImplementation(() =>
      Promise.resolve({data: {success: true}}),
    );
    MagmaAPIBindings.putLteByNetworkIdCellularRan.mockImplementation(() =>
      Promise.resolve({data: {success: true}}),
    );
  });

  afterEach(() => {
    axiosMock.get.mockClear();
    MagmaAPIBindings.getNetworksByNetworkId.mockClear();
    MagmaAPIBindings.putNetworksByNetworkId.mockClear();
    MagmaAPIBindings.getLteByNetworkIdCellularEpc.mockClear();
    MagmaAPIBindings.putLteByNetworkIdCellularEpc.mockClear();
    MagmaAPIBindings.getLteByNetworkIdCellularRan.mockClear();
    MagmaAPIBindings.putLteByNetworkIdCellularRan.mockClear();
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/test/network']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <Route path="/nms/:networkId/network" component={NetworkDashboard} />
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  it('Verify Network Dashboard', async () => {
    const {getByTestId, getByLabelText} = render(<Wrapper />);
    await wait();

    expect(MagmaAPIBindings.getNetworksByNetworkId).toHaveBeenCalledTimes(1);
    // eslint-disable-next-line max-len
    expect(MagmaAPIBindings.getLteByNetworkIdCellularEpc).toHaveBeenCalledTimes(
      1,
    );
    // eslint-disable-next-line max-len
    expect(MagmaAPIBindings.getLteByNetworkIdCellularRan).toHaveBeenCalledTimes(
      1,
    );
    expect(MagmaAPIBindings.getLteByNetworkIdGateways).toHaveBeenCalledTimes(1);
    expect(MagmaAPIBindings.getLteByNetworkIdEnodebs).toHaveBeenCalledTimes(1);
    // eslint-disable-next-line max-len
    expect(
      MagmaAPIBindings.getNetworksByNetworkIdPoliciesRules,
    ).toHaveBeenCalledTimes(1);
    // eslint-disable-next-line max-len
    expect(MagmaAPIBindings.getLteByNetworkIdSubscribers).toHaveBeenCalledTimes(
      1,
    );
    expect(MagmaAPIBindings.getLteByNetworkIdApns).toHaveBeenCalledTimes(1);

    const info = getByTestId('info');
    expect(info).toHaveTextContent('Test Network');
    expect(info).toHaveTextContent('test_network');
    expect(info).toHaveTextContent('Test Network Description');
    expect(info).toHaveTextContent('lte');

    const ran = getByTestId('ran');
    expect(ran).toHaveTextContent('20');
    expect(ran).toHaveTextContent('TDD');
    expect(ran).toHaveTextContent('44390');
    expect(ran).toHaveTextContent('7');
    expect(ran).toHaveTextContent('2');

    const epc = getByTestId('epc');
    expect(epc).toHaveTextContent('Enabled');
    expect(epc).toHaveTextContent('001');
    expect(epc).toHaveTextContent('01');
    expect(epc).toHaveTextContent('1');

    let epcPasswordInputElement = getByTestId('epcPassword').firstChild;
    if (
      epcPasswordInputElement instanceof HTMLInputElement &&
      epcPasswordInputElement.value &&
      epcPasswordInputElement.type
    ) {
      expect(epcPasswordInputElement.value).toBe('gAA=');
      expect(epcPasswordInputElement.type).toBe('password');
    } else {
      throw 'unexpected types';
    }

    fireEvent.click(getByLabelText('toggle password visibility'));
    await wait();
    epcPasswordInputElement = getByTestId('epcPassword').firstChild;
    if (
      epcPasswordInputElement instanceof HTMLInputElement &&
      epcPasswordInputElement.value &&
      epcPasswordInputElement.type
    ) {
      expect(epcPasswordInputElement.type).toBe('text');
    }

    // verify KPI tray
    expect(getByTestId('Gateways')).toHaveTextContent('1');
    expect(getByTestId('eNodeBs')).toHaveTextContent('1');
    expect(getByTestId('Subscribers')).toHaveTextContent('1');
    expect(getByTestId('Policies')).toHaveTextContent('2');
    expect(getByTestId('APNs')).toHaveTextContent('2');
  });

  it('Verify Network Add', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<Wrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByText('Add Network'));
    await wait();

    // check if only first tab (network) is active
    expect(queryByTestId('networkInfoEdit')).not.toBeNull();
    expect(queryByTestId('networkEpcEdit')).toBeNull();
    expect(queryByTestId('networkRanEdit')).toBeNull();

    let netIdField = getByTestId('networkID').firstChild;
    let netNameField = getByTestId('networkName').firstChild;
    let netDescField = getByTestId('networkDescription').firstChild;

    if (
      netIdField instanceof HTMLInputElement &&
      netNameField instanceof HTMLInputElement &&
      netDescField instanceof HTMLTextAreaElement
    ) {
      fireEvent.change(netIdField, {target: {value: 'testNetworkID'}});
      fireEvent.change(netNameField, {target: {value: 'Test LTE Network'}});
      fireEvent.change(netDescField, {
        target: {value: 'LTE test network description'},
      });
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(axiosMock.post).toHaveBeenCalledWith('/nms/network/create', {
      networkID: 'testNetworkID',
      data: {
        name: 'Test LTE Network',
        description: 'LTE test network description',
        networkType: 'lte',
        fegNetworkID: '',
        servedNetworkIDs: '',
      },
    });

    // now tab should move to epc edit component
    expect(queryByTestId('networkInfoEdit')).toBeNull();
    expect(queryByTestId('networkEpcEdit')).not.toBeNull();
    expect(queryByTestId('networkRanEdit')).toBeNull();

    // switch tab to network and verify editing of recently created network
    fireEvent.click(getByTestId('networkTab'));
    await wait();

    expect(queryByTestId('networkInfoEdit')).not.toBeNull();
    expect(queryByTestId('networkEpcEdit')).toBeNull();
    expect(queryByTestId('networkRanEdit')).toBeNull();

    netIdField = getByTestId('networkID').firstChild;
    netNameField = getByTestId('networkName').firstChild;
    netDescField = getByTestId('networkDescription').firstChild;

    if (
      netIdField instanceof HTMLInputElement &&
      netNameField instanceof HTMLInputElement &&
      netDescField instanceof HTMLTextAreaElement
    ) {
      expect(netIdField.value).toBe('testNetworkID');
      // networkID shouldn't be editable
      expect(netIdField.readOnly).toBe(true);
      expect(netNameField.value).toBe('Test LTE Network');
      expect(netDescField.value).toBe('LTE test network description');

      fireEvent.change(netDescField, {
        target: {value: 'New LTE test network description'},
      });
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(MagmaAPIBindings.putNetworksByNetworkId).toHaveBeenCalledWith({
      networkId: 'testNetworkID',
      network: {
        name: 'Test LTE Network',
        description: 'New LTE test network description',
        type: 'lte',
        dns: {
          enable_caching: false,
          local_ttl: 0,
          records: [],
        },
        id: 'testNetworkID',
      },
    });

    // verify adding EPC parameters
    const mncField = getByTestId('mnc').firstChild;
    const mccField = getByTestId('mcc').firstChild;
    const tacField = getByTestId('tac').firstChild;
    if (
      mncField instanceof HTMLInputElement &&
      mccField instanceof HTMLInputElement &&
      tacField instanceof HTMLInputElement
    ) {
      fireEvent.change(mncField, {target: {value: '02'}});
      fireEvent.change(mccField, {target: {value: '003'}});
      fireEvent.change(tacField, {target: {value: '1'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(MagmaAPIBindings.putLteByNetworkIdCellularEpc).toHaveBeenCalledWith({
      config: {
        cloud_subscriberdb_enabled: false,
        default_rule_id: 'default_rule_1',
        lte_auth_amf: 'gAA=',
        lte_auth_op: 'EREREREREREREREREREREQ==',
        mcc: '003',
        mnc: '02',
        network_services: ['policy_enforcement'],
        relay_enabled: false,
        sub_profiles: {},
        tac: 1,
      },
      networkId: 'testNetworkID',
    });

    // now save and continue should move to Ran component
    expect(queryByTestId('networkInfoEdit')).toBeNull();
    expect(queryByTestId('networkEpcEdit')).toBeNull();
    expect(queryByTestId('networkRanEdit')).not.toBeNull();

    const earfcndl = getByTestId('earfcndl').firstChild;
    const specialSubframePattern = getByTestId('specialSubframePattern')
      .firstChild;
    const subframeAssignment = getByTestId('subframeAssignment').firstChild;
    if (
      earfcndl instanceof HTMLElement &&
      subframeAssignment instanceof HTMLElement &&
      specialSubframePattern instanceof HTMLElement
    ) {
      fireEvent.change(earfcndl, {target: {value: '44000'}});
      fireEvent.change(specialSubframePattern, {target: {value: '8'}});
      fireEvent.change(subframeAssignment, {target: {value: '2'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save And Add Network'));
    await wait();
    expect(MagmaAPIBindings.putLteByNetworkIdCellularRan).toHaveBeenCalledWith({
      config: {
        bandwidth_mhz: 20,
        fdd_config: undefined,
        tdd_config: {
          earfcndl: 44000,
          special_subframe_pattern: 8,
          subframe_assignment: 2,
        },
      },
      networkId: 'testNetworkID',
    });
  });

  it('Verify Network Edit Info', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<Wrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('infoEditButton'));
    await wait();

    // check if first tab (network) is active
    expect(queryByTestId('networkInfoEdit')).not.toBeNull();
    expect(queryByTestId('networkEpcEdit')).toBeNull();
    expect(queryByTestId('networkRanEdit')).toBeNull();

    const netIdField = getByTestId('networkID').firstChild;
    const netNameField = getByTestId('networkName').firstChild;
    const netDescField = getByTestId('networkDescription').firstChild;

    if (
      netIdField instanceof HTMLInputElement &&
      netNameField instanceof HTMLInputElement &&
      netDescField instanceof HTMLTextAreaElement
    ) {
      expect(netIdField.value).toBe('test_network');

      // networkID shouldn't be editable
      expect(netIdField.readOnly).toBe(true);
      expect(netNameField.value).toBe('Test Network');
      expect(netDescField.value).toBe('Test Network Description');

      fireEvent.change(netDescField, {
        target: {value: 'Edit LTE test network description'},
      });
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await wait();
    expect(MagmaAPIBindings.putNetworksByNetworkId).toHaveBeenCalledWith({
      networkId: 'test_network',
      network: {
        ...testNetwork,
        description: 'Edit LTE test network description',
      },
    });

    // verify that info component is updated with edited info
    expect(getByTestId('info')).toHaveTextContent(
      'Edit LTE test network description',
    );
  });

  it('Verify Network Edit EPC', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<Wrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('epcEditButton'));
    await wait();

    expect(queryByTestId('networkInfoEdit')).toBeNull();
    expect(queryByTestId('networkEpcEdit')).not.toBeNull();
    expect(queryByTestId('networkRanEdit')).toBeNull();

    const mncField = getByTestId('mnc').firstChild;
    if (mncField instanceof HTMLInputElement) {
      fireEvent.change(mncField, {target: {value: '03'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save'));
    await wait();

    expect(MagmaAPIBindings.putLteByNetworkIdCellularEpc).toHaveBeenCalledWith({
      config: {...epc, mnc: '03'},
      networkId: 'test_network',
    });

    // verify epc component is updated with edited epc
    expect(getByTestId('epc')).toHaveTextContent('03');
  });

  it('Verify Network Edit Ran', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<Wrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('ranEditButton'));
    await wait();

    expect(queryByTestId('networkInfoEdit')).toBeNull();
    expect(queryByTestId('networkEpcEdit')).toBeNull();
    expect(queryByTestId('networkRanEdit')).not.toBeNull();

    const earfcndl = getByTestId('earfcndl').firstChild;
    if (earfcndl instanceof HTMLElement) {
      fireEvent.change(earfcndl, {target: {value: '40000'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save'));
    await wait();

    expect(MagmaAPIBindings.putLteByNetworkIdCellularRan).toHaveBeenCalledWith({
      config: {
        ...ran,
        tdd_config: {
          ...ran.tdd_config,
          earfcndl: 40000,
        },
      },
      networkId: 'test_network',
    });

    // verify ran component is updated with edited ran info
    expect(getByTestId('ran')).toHaveTextContent('40000');
  });
});
