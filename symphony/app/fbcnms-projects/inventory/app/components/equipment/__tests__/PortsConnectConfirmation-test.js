/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

jest.mock('../../../common/RelayEnvironment');

import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import PortsConnectConfirmation from '../PortsConnectConfirmation';
import React from 'react';
import {MemoryRouter} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';

import 'jest-dom/extend-expect';
import defaultTheme from '@fbcnms/ui/theme/default';

import {cleanup, render, wait} from '@testing-library/react';

const Wrapper = props => (
  <MemoryRouter initialEntries={['/inventory']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <SnackbarProvider>{props.children}</SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

afterEach(cleanup);

describe('<PortsConnectConfirmation />', () => {
  const eq = (name, tname) => {
    return {
      id: '3',
      name: name,
      equipmentType: {
        id: '1',
        name: tname,
        numberOfEquipment: 1,
        portDefinitions: [],
        positionDefinitions: [],
        propertyTypes: [],
      },
      workOrder: null,
      ports: [],
      futureState: 'INSTALL',
      device: null,
      locationHierarchy: [],
      positionHierarchy: [],
      parentPosition: null,
      parentLocation: null,
      positions: [],
      properties: [],
      services: [],
    };
  };
  const eqPort = (name, ename, tname) => {
    return {
      id: '1',
      link: {
        id: '2',
        futureState: 'INSTALL',
        ports: [],
        properties: [],
        workOrder: null,
        services: [],
      },
      definition: {
        id: '1',
        index: 1,
        name: name + 'name',
        portType: {
          name: name + 'type',
        },
      },
      parentEquipment: eq(ename, tname),
      properties: [],
      serviceEndpoints: [],
    };
  };

  it('renders PortsConnectConfirmation', async () => {
    const {getByText} = render(
      <Wrapper>
        <PortsConnectConfirmation
          // $FlowFixMe
          aSidePort={eqPort('a', 'a_eq_name', 'a_eq_type_name')}
          aSideEquipment={eq('a_eq_name', 'a_eq_type_name')}
          zSidePort={eqPort('z', 'z_eq_name', 'z_eq_type_name')}
          zSideEquipment={eq('z_eq_name', 'z_eq_type_name')}
          classes={{}}
        />
      </Wrapper>,
    );

    await wait();

    expect(
      getByText(
        'Are you sure you would like to connect port atype aname on a_eq_type_name a_eq_name to port ztype zname on z_eq_type_name z_eq_name?',
      ),
    ).toBeInTheDocument();
  });
});
