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

import './common/axiosConfig';
import './common/polyfill';

import LoginForm from './views/login/LoginForm';
import React from 'react';
import ReactDOM from 'react-dom';
import ThemeProvider from '@material-ui/styles/ThemeProvider';
import defaultTheme from './theme/default';
import nullthrows from '../shared/util/nullthrows';
import {AppContextProvider} from './components/context/AppContext';
import {BrowserRouter} from 'react-router-dom';

function LoginWrapper() {
  return (
    <LoginForm
      action="/user/login"
      title="Magma"
      ssoAction="/user/login/saml"
      ssoEnabled={window.CONFIG.appData.ssoEnabled}
      csrfToken={window.CONFIG.appData.csrfToken}
    />
  );
}

ReactDOM.render(
  <AppContextProvider>
    <BrowserRouter>
      <ThemeProvider theme={defaultTheme}>
        <LoginWrapper />
      </ThemeProvider>
    </BrowserRouter>
  </AppContextProvider>,
  nullthrows(document.getElementById('root')),
);
