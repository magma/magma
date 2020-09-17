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
import 'jest-dom/extend-expect';

import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import SubscriberContext from '../../../components/context/SubscriberContext';
import SubscriberDashboard from '../SubscriberOverview';
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

describe('<SubscriberDashboard />', () => {
  const subscribers = {
    IMSI00000000001002: {
      name: 'subscriber0',
      active_apns: ['oai.ipv4'],
      id: 'IMSI722070171001002',
      lte: {
        auth_algo: 'MILENAGE',
        auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
        auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
        state: 'ACTIVE',
        sub_profile: 'default',
      },
      config: {
        lte: {
          auth_algo: 'MILENAGE',
          auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
          auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
          state: 'ACTIVE',
          sub_profile: 'default',
        },
      },
    },
    IMSI00000000001003: {
      name: 'subscriber1',
      active_apns: ['oai.ipv4'],
      id: 'IMSI722070171001003',
      lte: {
        auth_algo: 'MILENAGE',
        auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
        auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
        state: 'INACTIVE',
        sub_profile: 'default',
      },
      config: {
        lte: {
          auth_algo: 'MILENAGE',
          auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
          auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
          state: 'INACTIVE',
          sub_profile: 'default',
        },
      },
    },
  };

  const Wrapper = () => {
    const subscriberCtx = {
      state: subscribers,
      gwSubscriberMap: {},
    };

    return (
      <MemoryRouter
        initialEntries={['/nms/test/subscribers/overview']}
        initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <SubscriberContext.Provider value={subscriberCtx}>
              <Route
                path="/nms/:networkId/subscribers/overview"
                component={SubscriberDashboard}
              />
            </SubscriberContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('Verify Subscribers Dashboard', async () => {
    const {getAllByTitle, getByTestId, getAllByRole} = render(<Wrapper />);
    await wait();
    const rowItems = await getAllByRole('row');

    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Name');
    expect(rowItems[0]).toHaveTextContent('IMSI');
    expect(rowItems[0]).toHaveTextContent('Service');
    expect(rowItems[0]).toHaveTextContent('Current Usage');
    expect(rowItems[0]).toHaveTextContent('Daily Average');
    expect(rowItems[0]).toHaveTextContent('Last Reported Time');

    expect(rowItems[1]).toHaveTextContent('subscriber0');
    expect(rowItems[1]).toHaveTextContent('IMSI00000000001002');
    expect(rowItems[1]).toHaveTextContent('ACTIVE');
    expect(rowItems[1]).toHaveTextContent('0');

    expect(rowItems[2]).toHaveTextContent('subscriber1');
    expect(rowItems[2]).toHaveTextContent('IMSI00000000001003');
    expect(rowItems[2]).toHaveTextContent('INACTIVE');
    expect(rowItems[2]).toHaveTextContent('0');

    // click the actions button for subscriber0
    const actionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(actionList[0]);
    await wait();
    expect(getByTestId('actions-menu')).toBeVisible();
  });
});
