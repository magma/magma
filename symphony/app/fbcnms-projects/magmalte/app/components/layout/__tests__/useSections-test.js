/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import MagmaAPIBindings from '@fbcnms/magma-api';
import NetworkContext from '../../context/NetworkContext';
import React from 'react';
import useSections from '../useSections';
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import {act, renderHook} from '@testing-library/react-hooks';
jest.mock('@fbcnms/magma-api');
jest.mock('mapbox-gl', () => {});
jest.mock('../../MapView', () => {});

import {AllNetworkTypes} from '@fbcnms/types/network';

global.CONFIG = {
  appData: {
    enabledFeatures: [],
  },
};

const wrapper = ({children}) => (
  <AppContextProvider networkIDs={['network1']}>
    <NetworkContext.Provider value={{networkId: 'network1'}}>
      {children}
    </NetworkContext.Provider>
  </AppContextProvider>
);

type TestCase = {
  default: string,
  sections: string[],
};

const testCases: {[string]: TestCase} = {
  lte: {
    default: 'map',
    sections: [
      'map',
      'metrics',
      'subscribers',
      'gateways',
      'enodebs',
      'configure',
    ],
  },
  mesh: {
    default: 'map',
    sections: [],
  },
  wifi_network: {
    default: 'map',
    sections: ['map', 'metrics', 'devices', 'configure'],
  },
  third_party: {
    default: 'devices',
    sections: ['devices', 'metrics', 'agents'],
  },
  symphony: {
    default: 'devices',
    sections: ['devices', 'metrics', 'agents'],
  },
  rhino: {
    default: 'metrics',
    sections: ['metrics'],
  },
  feg: {
    default: 'gateways',
    sections: ['gateways', 'configure'],
  },
  carrier_wifi_network: {
    default: 'gateways',
    sections: ['gateways', 'configure', 'metrics'],
  },
};

AllNetworkTypes.forEach(networkType => {
  const testCase = testCases[networkType];
  test('Should render ' + networkType, async () => {
    MagmaAPIBindings.getNetworksByNetworkIdType.mockResolvedValueOnce(
      networkType,
    );

    const {result, waitForNextUpdate} = renderHook(() => useSections(), {
      wrapper,
    });

    await act(async () => {
      // State is updated when we wait for the update, so we need this wrapped
      // in act
      await waitForNextUpdate();
    });

    expect(result.current[0]).toBe(testCase.default);

    const paths = result.current[1].map(r => r.path);
    expect(paths).toStrictEqual(testCase.sections);

    MagmaAPIBindings.getNetworksByNetworkIdType.mockClear();
  });
});
