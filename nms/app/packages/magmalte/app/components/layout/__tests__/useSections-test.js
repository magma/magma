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
 * @flow
 * @format
 */

import MagmaAPIBindings from '@fbcnms/magma-api';
import NetworkContext from '../../context/NetworkContext';
import React from 'react';
import useSections from '../useSections';

import {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import {act, renderHook} from '@testing-library/react-hooks';

const enqueueSnackbarMock = jest.fn();
jest.mock('@fbcnms/magma-api');
jest.mock('mapbox-gl', () => {});
jest.mock('@fbcnms/ui/insights/map/MapView', () => {});
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
  .mockReturnValue(enqueueSnackbarMock);

import {CWF, FEG, FEG_LTE, LTE, WIFI, XWFM} from '@fbcnms/types/network';

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
      'alerts',
    ],
  },
  feg_lte: {
    default: 'map',
    sections: [
      'map',
      'metrics',
      'subscribers',
      'gateways',
      'enodebs',
      'configure',
      'alerts',
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
  feg: {
    default: 'gateways',
    sections: ['gateways', 'configure', 'alerts', 'metrics'],
  },
  carrier_wifi_network: {
    default: 'gateways',
    sections: ['gateways', 'configure', 'metrics', 'alerts'],
  },
  xwfm: {
    default: 'gateways',
    sections: ['gateways', 'configure', 'metrics', 'alerts'],
  },
};

const AllNetworkTypes = [CWF, FEG, LTE, FEG_LTE, WIFI, XWFM];
AllNetworkTypes.forEach(networkType => {
  const testCase = testCases[networkType];
  // XWF-M network selection in NMS creates a CWF network on the API just with
  // different config defaults
  const apiNetworkType = networkType === XWFM ? CWF : networkType;
  test('Should render ' + networkType, async () => {
    MagmaAPIBindings.getNetworksByNetworkIdType.mockResolvedValue(
      apiNetworkType,
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
