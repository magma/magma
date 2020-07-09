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

import AddEditEnodeButton from '../EnodebDetailConfigEdit';
import EnodebConfig from '../EnodebDetailConfig';
import MagmaAPIBindings from '@fbcnms/magma-api';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default.js';

import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {cleanup, fireEvent, render, wait} from '@testing-library/react';
import {useState} from 'react';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');
afterEach(cleanup);
const enqueueSnackbarMock = jest.fn();
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
  .mockReturnValue(enqueueSnackbarMock);

describe('<AddEditEnodeButton />', () => {
  const ran = {
    bandwidth_mhz: 20,
    tdd_config: {
      earfcndl: 44390,
      special_subframe_pattern: 7,
      subframe_assignment: 2,
    },
  };

  const enbInfo = {
    enb: {
      attached_gateway_id: '',
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
      name: 'testEnodeb0',
      serial: 'testEnodebSerial0',
      description: 'test enodeb description',
    },
    enb_state: {
      enodeb_configured: true,
      enodeb_connected: true,
      fsm_state: 'Completed provisioning eNB. Awaiting new Inform.',
      gps_connected: true,
      gps_latitude: '',
      gps_longitude: '',
      mme_connected: false,
      opstate_enabled: false,
      ptp_connected: false,
      reporting_gateway_id: '',
      rf_tx_desired: true,
      rf_tx_on: false,
    },
  };

  beforeEach(() => {
    MagmaAPIBindings.getLteByNetworkIdCellularRan.mockResolvedValue(ran);
  });

  afterEach(() => {
    MagmaAPIBindings.getLteByNetworkIdCellularRan.mockClear();
    MagmaAPIBindings.postLteByNetworkIdEnodebs.mockClear();
    MagmaAPIBindings.putLteByNetworkIdEnodebsByEnodebSerial.mockClear();
  });

  const AddWrapper = () => (
    <MemoryRouter initialEntries={['/nms/test/enode']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <Route
            path="/nms/:networkId/enode"
            render={props => (
              <AddEditEnodeButton
                {...props}
                title="Add Enodeb"
                isLink={false}
              />
            )}
          />
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  const DetailWrapper = () => {
    const [enbInf, setEnbInfo] = useState(enbInfo);
    return (
      <MemoryRouter
        initialEntries={['/nms/mynetwork/enodeb/testEnodebSerial0/overview']}
        initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <Route
              path="/nms/:networkId/enodeb/:enodebSerial/overview"
              render={props => (
                <EnodebConfig
                  {...props}
                  enbInfo={enbInf}
                  onSave={enb => {
                    setEnbInfo({...enbInf, enb: enb});
                  }}
                />
              )}
            />
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('Verify Enode Configs', async () => {
    const {getByTestId} = render(<DetailWrapper />);
    await wait();
    expect(MagmaAPIBindings.getLteByNetworkIdCellularRan).toHaveBeenCalledTimes(
      1,
    );
    const config = getByTestId('config');
    expect(config).toHaveTextContent('testEnodeb0');
    expect(config).toHaveTextContent('testEnodeb0Serial');
    expect(config).toHaveTextContent('test enodeb description');

    const ran = getByTestId('ran');
    expect(ran).toHaveTextContent('20');
    expect(ran).toHaveTextContent('TDD');
    expect(ran).toHaveTextContent('44290');
    expect(ran).toHaveTextContent('7');
    expect(ran).toHaveTextContent('2');
  });

  it('Verify Enode Add', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<AddWrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByText('Add Enodeb'));
    await wait();

    // check if only first tab (config) is active
    expect(queryByTestId('configEdit')).not.toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();

    let enbSerial = getByTestId('serial').firstChild;
    let enbName = getByTestId('name').firstChild;
    let enbDesc = getByTestId('description').firstChild;

    if (
      enbSerial instanceof HTMLInputElement &&
      enbName instanceof HTMLInputElement &&
      enbDesc instanceof HTMLTextAreaElement
    ) {
      fireEvent.change(enbSerial, {target: {value: 'TestEnodebSerial1'}});
      fireEvent.change(enbName, {target: {value: 'Test Enodeb 1'}});
      fireEvent.change(enbDesc, {
        target: {value: 'Enode1 Description'},
      });
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(MagmaAPIBindings.postLteByNetworkIdEnodebs).toHaveBeenCalledWith({
      enodeb: {
        config: {
          cell_id: 0,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          transmit_enabled: false,
        },
        description: 'Enode1 Description',
        name: 'Test Enodeb 1',
        serial: 'TestEnodebSerial1',
      },
      networkId: 'test',
    });

    // now tab should move to ran edit component
    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).not.toBeNull();

    // switch tab to config and verify editing of recently created enodeb
    fireEvent.click(getByTestId('configTab'));
    await wait();

    expect(queryByTestId('configEdit')).not.toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();

    enbSerial = getByTestId('serial').firstChild;
    enbName = getByTestId('name').firstChild;
    enbDesc = getByTestId('description').firstChild;

    if (
      enbSerial instanceof HTMLInputElement &&
      enbName instanceof HTMLInputElement &&
      enbDesc instanceof HTMLTextAreaElement
    ) {
      expect(enbSerial.value).toBe('TestEnodebSerial1');
      expect(enbName.value).toBe('Test Enodeb 1');
      expect(enbDesc.value).toBe('Enode1 Description');

      // enodeb serial shouldn't be editable
      expect(enbSerial.readOnly).toBe(true);
      fireEvent.change(enbDesc, {
        target: {value: 'Enode1 New Description'},
      });
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(
      MagmaAPIBindings.putLteByNetworkIdEnodebsByEnodebSerial,
    ).toHaveBeenCalledWith({
      enodeb: {
        config: {
          cell_id: 0,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          transmit_enabled: false,
        },
        description: 'Enode1 New Description',
        name: 'Test Enodeb 1',
        serial: 'TestEnodebSerial1',
      },
      enodebSerial: 'TestEnodebSerial1',
      networkId: 'test',
    });

    // clear mock info
    MagmaAPIBindings.putLteByNetworkIdEnodebsByEnodebSerial.mockClear();

    // now tab should move to ran edit component
    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).not.toBeNull();

    const earfcndl = getByTestId('earfcndl').firstChild;
    const pci = getByTestId('pci').firstChild;
    const cellId = getByTestId('cellId').firstChild;
    if (
      earfcndl instanceof HTMLElement &&
      pci instanceof HTMLElement &&
      cellId instanceof HTMLElement
    ) {
      fireEvent.change(earfcndl, {target: {value: '44000'}});
      fireEvent.change(pci, {target: {value: '8'}});
      fireEvent.change(cellId, {target: {value: '2'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Add eNodeB'));
    await wait();
    expect(
      MagmaAPIBindings.putLteByNetworkIdEnodebsByEnodebSerial,
    ).toHaveBeenCalledWith({
      enodeb: {
        config: {
          bandwidth_mhz: 20,
          cell_id: 2,
          earfcndl: 44000,
          pci: 8,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          transmit_enabled: false,
        },
        description: 'Enode1 New Description',
        name: 'Test Enodeb 1',
        serial: 'TestEnodebSerial1',
      },
      enodebSerial: 'TestEnodebSerial1',
      networkId: 'test',
    });
    expect(MagmaAPIBindings.getLteByNetworkIdCellularRan).toHaveBeenCalledTimes(
      1,
    );
  });

  it('Verify Enode Edit Config', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<DetailWrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('configEditButton'));
    await wait();

    // check if only first tab (config) is active
    expect(queryByTestId('configEdit')).not.toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();

    const enbSerial = getByTestId('serial').firstChild;
    const enbName = getByTestId('name').firstChild;
    const enbDesc = getByTestId('description').firstChild;

    if (
      enbSerial instanceof HTMLInputElement &&
      enbName instanceof HTMLInputElement &&
      enbDesc instanceof HTMLTextAreaElement
    ) {
      expect(enbSerial.value).toBe('testEnodebSerial0');
      expect(enbName.value).toBe('testEnodeb0');
      expect(enbDesc.value).toBe('test enodeb description');
      fireEvent.change(enbDesc, {
        target: {value: 'test enodeb new description'},
      });
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await wait();
    expect(
      MagmaAPIBindings.putLteByNetworkIdEnodebsByEnodebSerial,
    ).toHaveBeenCalledWith({
      enodeb: {
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
        attached_gateway_id: '',
        description: 'test enodeb new description',
        name: 'testEnodeb0',
        serial: 'testEnodebSerial0',
      },
      enodebSerial: 'testEnodebSerial0',
      networkId: 'mynetwork',
    });

    const config = getByTestId('config');
    expect(config).toHaveTextContent('test enodeb new description');
  });

  it('Verify Enode Edit Ran', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<DetailWrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('ranEditButton'));
    await wait();

    expect(MagmaAPIBindings.getLteByNetworkIdCellularRan).toHaveBeenCalledTimes(
      1,
    );
    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).not.toBeNull();

    const earfcndl = getByTestId('earfcndl').firstChild;
    const pci = getByTestId('pci').firstChild;
    const cellId = getByTestId('cellId').firstChild;
    if (
      earfcndl instanceof HTMLInputElement &&
      pci instanceof HTMLInputElement &&
      cellId instanceof HTMLInputElement
    ) {
      expect(earfcndl.value).toBe('44290');
      expect(pci.value).toBe('36');
      expect(cellId.value).toBe('1');
      fireEvent.change(earfcndl, {target: {value: '40000'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await wait();
    expect(
      MagmaAPIBindings.putLteByNetworkIdEnodebsByEnodebSerial,
    ).toHaveBeenCalledWith({
      enodeb: {
        attached_gateway_id: '',
        config: {
          bandwidth_mhz: 20,
          cell_id: 1,
          device_class: 'Baicells ID TDD/FDD',
          earfcndl: 40000,
          pci: 36,
          special_subframe_pattern: 7,
          subframe_assignment: 2,
          tac: 1,
          transmit_enabled: true,
        },
        description: 'test enodeb description',
        name: 'testEnodeb0',
        serial: 'testEnodebSerial0',
      },
      enodebSerial: 'testEnodebSerial0',
      networkId: 'mynetwork',
    });
  });
});
