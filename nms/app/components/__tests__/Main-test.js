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

import ApplicationMain from '../ApplicationMain';
import MagmaAPIBindings from '../../../generated/MagmaAPIBindings';
import Main, {NO_NETWORK_MESSAGE} from '../Main';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AppContextProvider} from '../context/AppContext';
import {MemoryRouter} from 'react-router-dom';
import {render, wait} from '@testing-library/react';

jest.mock('../../../generated/MagmaAPIBindings');

jest.mock('../main/Index', () => ({
  __esModule: true,
  ...jest.requireActual('../main/Index'),
  default: () => <div>Index</div>,
}));

jest.mock('../IndexWithoutNetwork', () => ({
  __esModule: true,
  default: () => <div>IndexWithoutNetwork</div>,
}));

const Wrapper = props => (
  <MemoryRouter initialEntries={[props.path]} initialIndex={0}>
    <AppContextProvider>
      <ApplicationMain>{props.children}</ApplicationMain>
    </AppContextProvider>
  </MemoryRouter>
);

describe.each`
  path                | text                     | networks
  ${'/nms/mynetwork'} | ${'Index'}               | ${['mynetwork']}
  ${'/nms'}           | ${'Index'}               | ${['mynetwork']}
  ${'/nms'}           | ${NO_NETWORK_MESSAGE}    | ${[]}
  ${'/admin'}         | ${'IndexWithoutNetwork'} | ${[]}
  ${'/admin'}         | ${'Index'}               | ${['mynetwork']}
  ${'/settings'}      | ${'IndexWithoutNetwork'} | ${[]}
  ${'/settings'}      | ${'Index'}               | ${['mynetwork']}
`('renders $path with networks $networks', ({path, text, networks}) => {
  beforeEach(() => {
    MagmaAPIBindings.getNetworks.mockResolvedValueOnce(networks);
  });

  it(`renders for ${path} path`, async () => {
    window.CONFIG = {
      appData: {
        enabledFeatures: [],
        user: {
          isSuperUser: false,
        },
      },
    };

    const {getByText} = render(
      <Wrapper path={path}>
        <Main />
      </Wrapper>,
    );

    await wait();

    expect(getByText(text)).toBeInTheDocument();
  });
});
