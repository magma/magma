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

import MagmaAPI from '../../api/MagmaAPI';
import NetworkContext from '../context/NetworkContext';
import NetworkSelector from '../NetworkSelector';
import React from 'react';
import {AppContextProvider} from '../context/AppContext';
import {LTE} from '../../../shared/types/network';
import {MemoryRouter} from 'react-router-dom';
import {SnackbarProvider} from 'notistack';
import {fireEvent, render, waitFor} from '@testing-library/react';
import {mockAPI} from '../../util/TestUtils';
import type {EmbeddedData} from '../../../shared/types/embeddedData';

const Wrapper = (props: {
  currentNetworkId?: string;
  children: React.ReactNode;
  isSuperUser: boolean;
}) => {
  window.CONFIG = {
    appData: {
      user: {
        isSuperUser: props.isSuperUser,
      },
    } as EmbeddedData,
  };

  return (
    <MemoryRouter initialEntries={['/nms']} initialIndex={0}>
      <SnackbarProvider>
        <AppContextProvider>
          <NetworkContext.Provider
            value={{
              networkId: props.currentNetworkId,
            }}>
            {props.children}
          </NetworkContext.Provider>
        </AppContextProvider>
      </SnackbarProvider>
    </MemoryRouter>
  );
};

describe('NetworkSelector', () => {
  it('renders nothing without network for regular user', () => {
    mockAPI(MagmaAPI.networks, 'networksGet', []);
    const {container} = render(
      <Wrapper isSuperUser={false}>
        <NetworkSelector />
      </Wrapper>,
    );
    expect(container).toBeEmpty();
  });

  it('renders text with single network for regular user', () => {
    mockAPI(MagmaAPI.networks, 'networksGet', ['test']);
    mockAPI(MagmaAPI.networks, 'networksNetworkIdTypeGet', LTE);

    const {queryByRole, getByText} = render(
      <Wrapper isSuperUser={false} currentNetworkId="test">
        <NetworkSelector />
      </Wrapper>,
    );
    expect(getByText('Network: test')).toBeInTheDocument();
    expect(queryByRole('button')).not.toBeInTheDocument();
  });

  it('renders menu with network links for regular user', async () => {
    mockAPI(MagmaAPI.networks, 'networksGet', ['test', 'other']);
    mockAPI(MagmaAPI.networks, 'networksNetworkIdTypeGet', LTE);

    const {getByRole, queryAllByRole} = render(
      <Wrapper isSuperUser={false} currentNetworkId="test">
        <NetworkSelector />
      </Wrapper>,
    );

    await waitFor(() =>
      expect(getByRole('button')).toHaveTextContent('Network: test'),
    );
    fireEvent.click(getByRole('button'));

    await waitFor(() =>
      expect(queryAllByRole('menuitem').map(link => link.textContent)).toEqual([
        'test',
        'other',
      ]),
    );
  });

  it('renders menu with network links and extra entries for super user', async () => {
    mockAPI(MagmaAPI.networks, 'networksGet', ['test', 'other']);
    mockAPI(MagmaAPI.networks, 'networksNetworkIdTypeGet', LTE);

    const {getByRole, queryAllByRole} = render(
      <Wrapper isSuperUser={true} currentNetworkId="test">
        <NetworkSelector />
      </Wrapper>,
    );
    await waitFor(() =>
      expect(getByRole('button')).toHaveTextContent('Network: test'),
    );
    fireEvent.click(getByRole('button'));

    await waitFor(() =>
      expect(queryAllByRole('menuitem').map(link => link.textContent)).toEqual([
        'Create Network',
        'Manage Networks',
        'test',
        'other',
      ]),
    );
  });
});
