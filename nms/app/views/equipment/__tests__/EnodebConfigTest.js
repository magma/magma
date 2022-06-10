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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AddEditEnodeButton from '../EnodebDetailConfigEdit';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import EnodebConfig from '../EnodebDetailConfig';
// $FlowFixMe migrated to typescript
import EnodebContext from '../../../components/context/EnodebContext';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../../theme/default';

import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
// $FlowFixMe migrated to typescript
import MagmaAPI from '../../../../api/MagmaAPI';
// $FlowFixMe migrated to typescript
import {SetEnodebState} from '../../../state/lte/EquipmentState';
import {fireEvent, render, wait} from '@testing-library/react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';
import {useState} from 'react';

jest.mock('axios');
jest.mock('../../../hooks/useSnackbar');

describe('<AddEditEnodeButton />', () => {
  beforeEach(() => {
    jest
      .spyOn(MagmaAPI.enodebs, 'lteNetworkIdEnodebsPost')
      .mockImplementation();
    jest
      .spyOn(MagmaAPI.enodebs, 'lteNetworkIdEnodebsEnodebSerialPut')
      .mockImplementation();
  });

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
      ip_address: '192.168.1.254',
    },
  };

  beforeEach(() => {
    (useEnqueueSnackbar: JestMockFn<
      Array<empty>,
      $Call<typeof useEnqueueSnackbar>,
    >).mockReturnValue(jest.fn());
  });

  const AddWrapper = () => {
    const [enbInf, setEnbInfo] = useState({testEnodebSerial0: enbInfo});
    return (
      <MemoryRouter initialEntries={['/nms/test/enode']} initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <EnodebContext.Provider
              value={{
                state: {enbInfo: enbInf},
                lteRanConfigs: ran,
                setState: async (key, value?) =>
                  SetEnodebState({
                    enbInfo: enbInf,
                    setEnbInfo: setEnbInfo,
                    networkId: 'test',
                    key: key,
                    value: value,
                  }),
              }}>
              <Routes>
                <Route
                  path="/nms/:networkId/enode"
                  element={
                    <AddEditEnodeButton title="Add Enodeb" isLink={false} />
                  }
                />
              </Routes>
            </EnodebContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  const DetailWrapper = () => {
    const [enbInf, setEnbInfo] = useState({testEnodebSerial0: enbInfo});
    return (
      <MemoryRouter
        initialEntries={['/nms/mynetwork/enodeb/testEnodebSerial0/overview']}
        initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <EnodebContext.Provider
              value={{
                state: {enbInfo: enbInf},
                lteRanConfigs: ran,
                setState: async (key, value?) =>
                  SetEnodebState({
                    enbInfo: enbInf,
                    setEnbInfo: setEnbInfo,
                    networkId: 'mynetwork',
                    key: key,
                    value: value,
                  }),
              }}>
              <Routes>
                <Route
                  path="/nms/:networkId/enodeb/:enodebSerial/overview"
                  element={<EnodebConfig />}
                />
              </Routes>
            </EnodebContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('Verify Enode Configs', async () => {
    const {getByTestId} = render(<DetailWrapper />);
    await wait();

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

  it('Verify Enode unManaged eNodeBs', async () => {
    const unmanagedEnbInfo = {
      enb: {
        attached_gateway_id: '',
        config: {
          cell_id: 0,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          transmit_enabled: false,
        },
        enodeb_config: {
          config_type: 'UNMANAGED',
          unmanaged_config: {
            cell_id: 111,
            ip_address: '1.1.1.2',
            tac: 1,
          },
        },
        name: 'testEnodeb1',
        serial: 'testEnodebSerial1',
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
        ip_address: '192.168.1.254',
      },
    };
    const UnmanagedEnodeWrapper = () => {
      const enbInf = {
        testEnodebSerial1: unmanagedEnbInfo,
      };
      return (
        <MemoryRouter
          initialEntries={['/nms/mynetwork/enodeb/testEnodebSerial1/overview']}
          initialIndex={0}>
          <MuiThemeProvider theme={defaultTheme}>
            <MuiStylesThemeProvider theme={defaultTheme}>
              <EnodebContext.Provider
                value={{
                  state: {enbInfo: enbInf},
                  setState: async () => {},
                }}>
                <Routes>
                  <Route
                    path="/nms/:networkId/enodeb/:enodebSerial/overview"
                    element={<EnodebConfig />}
                  />
                </Routes>
              </EnodebContext.Provider>
            </MuiStylesThemeProvider>
          </MuiThemeProvider>
        </MemoryRouter>
      );
    };

    const {getByTestId} = render(<UnmanagedEnodeWrapper />);
    await wait();
    const config = getByTestId('config');
    expect(config).toHaveTextContent('testEnodeb1');
    expect(config).toHaveTextContent('testEnodebSerial1');
    expect(config).toHaveTextContent('test enodeb description');

    const ran = getByTestId('ran');
    expect(ran).toHaveTextContent('1');
    expect(ran).toHaveTextContent('1.1.1.2');
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

    // test adding an existing enodeb
    if (enbSerial instanceof HTMLInputElement) {
      fireEvent.change(enbSerial, {target: {value: 'testEnodebSerial0'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(getByTestId('configEditError')).toHaveTextContent(
      'eNodeB testEnodebSerial0 already exists',
    );

    // test adding new enodeb
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
    expect(MagmaAPI.enodebs.lteNetworkIdEnodebsPost).toHaveBeenCalledWith({
      enodeb: {
        config: {
          cell_id: 0,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          transmit_enabled: false,
        },
        enodeb_config: {
          config_type: 'MANAGED',
          managed_config: {
            cell_id: 0,
            device_class: 'Baicells Nova-233 G2 OD FDD',
            transmit_enabled: false,
          },
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
      MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialPut,
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
        enodeb_config: {
          config_type: 'MANAGED',
          managed_config: {
            cell_id: 0,
            device_class: 'Baicells Nova-233 G2 OD FDD',
            transmit_enabled: false,
          },
        },
      },
      enodebSerial: 'TestEnodebSerial1',
      networkId: 'test',
    });

    // clear mock info
    MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialPut.mockClear();

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
      MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialPut,
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
        enodeb_config: {
          config_type: 'MANAGED',
          managed_config: {
            bandwidth_mhz: 20,
            cell_id: 2,
            earfcndl: 44000,
            pci: 8,
            device_class: 'Baicells Nova-233 G2 OD FDD',
            transmit_enabled: false,
          },
        },
        description: 'Enode1 New Description',
        name: 'Test Enodeb 1',
        serial: 'TestEnodebSerial1',
      },
      enodebSerial: 'TestEnodebSerial1',
      networkId: 'test',
    });
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
      MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialPut,
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
          unmanaged_config: undefined,
        },
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
      MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialPut,
    ).toHaveBeenCalledWith({
      enodeb: {
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
        enodeb_config: {
          config_type: 'MANAGED',
          unmanaged_config: undefined,
          managed_config: {
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
        },
        description: 'test enodeb description',
        name: 'testEnodeb0',
        serial: 'testEnodebSerial0',
      },
      enodebSerial: 'testEnodebSerial0',
      networkId: 'mynetwork',
    });
  });

  it('Verify Enode Edit unmanaged eNodeB', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<DetailWrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('ranEditButton'));
    await wait();

    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).not.toBeNull();

    const enbConfigType = getByTestId('enodeb_config').firstChild;
    if (
      enbConfigType instanceof HTMLElement &&
      enbConfigType.firstChild instanceof HTMLElement
    ) {
      fireEvent.click(enbConfigType.firstChild);
    } else {
      throw 'invalid type';
    }
    await wait();

    const ipAddress = getByTestId('ipAddress').firstChild;
    const cellId = getByTestId('cellId').firstChild;
    const tac = getByTestId('tac').firstChild;
    if (
      ipAddress instanceof HTMLInputElement &&
      cellId instanceof HTMLInputElement &&
      tac instanceof HTMLInputElement
    ) {
      fireEvent.change(ipAddress, {target: {value: '1.1.1.1'}});
      fireEvent.change(cellId, {target: {value: '1'}});
      fireEvent.change(tac, {target: {value: '1'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await wait();
    expect(
      MagmaAPI.enodebs.lteNetworkIdEnodebsEnodebSerialPut,
    ).toHaveBeenCalledWith({
      enodeb: {
        config: {
          cell_id: 0,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          transmit_enabled: false,
        },
        enodeb_config: {
          config_type: 'UNMANAGED',
          managed_config: undefined,
          unmanaged_config: {
            cell_id: 1,
            ip_address: '1.1.1.1',
            tac: 1,
          },
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
