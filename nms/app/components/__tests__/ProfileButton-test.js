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

import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import ProfileButton from '../ProfileButton';
import React, {useState} from 'react';
import defaultTheme from '../../theme/default';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AppContextProvider} from '../context/AppContext';
import {MemoryRouter} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
// $FlowFixMe[missing-export]
import {fireEvent, render, waitFor} from '@testing-library/react';

type Props = {
  expanded: boolean,
  path: string,
  isOrganizations: boolean,
};

const WrappedProfileButton = (props: Props) => {
  const [isMenuOpen, setMenuOpen] = useState(false);
  return (
    <MemoryRouter initialEntries={[props.path]} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <SnackbarProvider>
            <AppContextProvider isOrganizations={props.isOrganizations}>
              <ProfileButton
                isMenuOpen={isMenuOpen}
                setMenuOpen={setMenuOpen}
                expanded={props.expanded}
              />
            </AppContextProvider>
          </SnackbarProvider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
};

describe('<ProfileButton />', () => {
  it.each([true, false])('respects expanded=%s', expanded => {
    window.CONFIG = {
      appData: {
        user: {},
        ssoEnabled: false,
        enabledFeatures: [],
      },
    };

    const {queryByText} = render(
      <WrappedProfileButton
        path="/admin"
        isOrganizations={false}
        expanded={expanded}
      />,
    );
    if (expanded) {
      expect(queryByText('Account & Settings')).toBeInTheDocument();
    } else {
      expect(queryByText('Account & Settings')).not.toBeInTheDocument();
    }
  });

  async function getRenderedLinks({
    isOrganizations = false,
    isSuperUser = false,
    ssoEnabled = false,
    enabledFeatures = [],
  }: {
    isSuperUser?: boolean,
    isOrganizations?: boolean,
    ssoEnabled?: boolean,
    enabledFeatures?: Array<string>,
  }) {
    window.CONFIG = {
      appData: {
        user: {isSuperUser},
        ssoEnabled,
        enabledFeatures,
      },
    };

    const {getByRole, getByTestId} = render(
      <WrappedProfileButton
        path="/admin"
        isOrganizations={isOrganizations}
        expanded={false}
      />,
    );

    const button = getByTestId('profileButton');
    fireEvent.click(button);

    const links = await waitFor(() =>
      getByRole('navigation').querySelectorAll('a'),
    );
    return Array.from(links).map(t => t.textContent);
  }

  it.each([true, false])(
    'renders Account Settings depending on ssoEnabled',
    async ssoEnabled => {
      const links = await getRenderedLinks({ssoEnabled});
      if (ssoEnabled) {
        expect(links).not.toContain('Account Settings');
      } else {
        expect(links).toContain('Account Settings');
      }
    },
  );

  it.each([true, false])(
    'renders Administration depending on isSuperUser',
    async isSuperUser => {
      const links = await getRenderedLinks({isSuperUser});
      if (isSuperUser) {
        expect(links).toContain('Administration');
      } else {
        expect(links).not.toContain('Administration');
      }
    },
  );

  it('does not render Administration on organizations page', async () => {
    const links = await getRenderedLinks({
      isOrganizations: true,
      isSuperUser: true,
    });
    expect(links).not.toContain('Administration');
  });

  it.each([true, false])(
    'renders Documentation depending on feature flag',
    async isFeatureEnabled => {
      const links = await getRenderedLinks({
        enabledFeatures: isFeatureEnabled ? ['documents_site'] : [],
      });
      if (isFeatureEnabled) {
        expect(links).toContain('Documentation');
      } else {
        expect(links).not.toContain('Documentation');
      }
    },
  );
});
