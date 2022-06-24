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
 */

import ApplicationMain from '../ApplicationMain';
import MagmaAPI from '../../../api/MagmaAPI';
import Main, {NO_NETWORK_MESSAGE} from '../Main';
import React from 'react';
import {AppContextProvider} from '../context/AppContext';
import {EmbeddedData} from '../../../shared/types/embeddedData';
import {MemoryRouter} from 'react-router-dom';
import {mockAPI} from '../../util/TestUtils';
import {render, wait} from '@testing-library/react';

// eslint-disable-next-line @typescript-eslint/no-unsafe-return
jest.mock('../main/Index', () => ({
  __esModule: true,
  ...jest.requireActual('../main/Index'),
  default: () => <div>Index</div>,
}));

jest.mock('../IndexWithoutNetwork', () => ({
  __esModule: true,
  default: () => <div>IndexWithoutNetwork</div>,
}));

const Wrapper = (props: {path: string; children: React.ReactNode}) => (
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
`(
  'renders $path with networks $networks',
  ({
    path,
    text,
    networks,
  }: {
    path: string;
    text: string;
    networks: Array<string>;
  }) => {
    beforeEach(() => {
      mockAPI(MagmaAPI.networks, 'networksGet', networks);
    });

    it(`renders for ${path} path`, async () => {
      window.CONFIG = {
        appData: ({
          enabledFeatures: [],
          user: {
            isSuperUser: false,
          },
        } as unknown) as EmbeddedData,
      };

      const {getByText} = render(
        <Wrapper path={path}>
          <Main />
        </Wrapper>,
      );

      await wait();

      expect(getByText(text)).toBeInTheDocument();
    });
  },
);
