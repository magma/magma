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

import MagmaAPIBindings from '../../../../generated/MagmaAPIBindings';
import NetworkContext from '../../context/NetworkContext';
import React from 'react';
import useSections from '../useSections';

import {AppContextProvider} from '../../context/AppContext';
import {act, renderHook} from '@testing-library/react-hooks';

import {CWF, FEG, FEG_LTE, LTE, XWFM} from '../../../../shared/types/network';
import {EmbeddedData} from '../../../../shared/types/embeddedData';

jest.mock('../../../../generated/MagmaAPIBindings.js');

window.CONFIG = {
  appData: ({
    enabledFeatures: [],
  } as unknown) as EmbeddedData,
};

const wrapper = ({children}: {children: React.ReactNode}) => (
  <AppContextProvider networkIDs={['network1']}>
    <NetworkContext.Provider value={{networkId: 'network1'}}>
      {children}
    </NetworkContext.Provider>
  </AppContextProvider>
);

const testCases = {
  lte: {
    default: 'dashboard',
    sections: [
      'dashboard',
      'equipment',
      'network',
      'subscribers',
      'traffic',
      'tracing',
      'metrics',
    ],
  },
  feg_lte: {
    default: 'dashboard',
    sections: [
      'dashboard',
      'equipment',
      'network',
      'subscribers',
      'traffic',
      'tracing',
      'metrics',
    ],
  },
  feg: {
    default: 'dashboard',
    sections: [
      'dashboard',
      'gateways',
      'network',
      'equipment',
      'configure',
      'alerts',
      'metrics',
    ],
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

describe.each([CWF, FEG, LTE, FEG_LTE, XWFM])('Should render', networkType => {
  const testCase = testCases[networkType];
  // XWF-M network selection in NMS creates a CWF network on the API just with
  // different config defaults
  const apiNetworkType = networkType === XWFM ? CWF : networkType;
  it(networkType, async () => {
    (MagmaAPIBindings.getNetworksByNetworkIdType as jest.Mock).mockResolvedValue(
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
  });
});
