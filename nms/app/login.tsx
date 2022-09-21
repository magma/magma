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

import './util/axiosConfig';
import './util/polyfill';

import CssBaseline from '@mui/material/CssBaseline';
import LoginForm from './views/login/LoginForm';
import React from 'react';
import ReactDOM from 'react-dom';
import defaultTheme from './theme/default';
import nullthrows from '../shared/util/nullthrows';
import {AppContextProvider} from './context/AppContext';
import {BrowserRouter} from 'react-router-dom';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';

const LOGIN_ERROR_MESSAGE = 'Invalid email or password';

function LoginWrapper() {
  const params = new URLSearchParams(window.location.search);
  const loginInvalid = params.get('invalid');
  return (
    <LoginForm
      action="/user/login"
      title="Magma"
      ssoAction="/user/login/saml"
      ssoEnabled={window.CONFIG.appData.ssoEnabled}
      csrfToken={window.CONFIG.appData.csrfToken}
      error={loginInvalid ? LOGIN_ERROR_MESSAGE : undefined}
    />
  );
}

ReactDOM.render(
  <AppContextProvider>
    <BrowserRouter>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={defaultTheme}>
          <CssBaseline />
          <LoginWrapper />
        </ThemeProvider>
      </StyledEngineProvider>
    </BrowserRouter>
  </AppContextProvider>,
  nullthrows(document.getElementById('root')),
);
